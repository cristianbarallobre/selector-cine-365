package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Recovery backups stay on this computer. No cloud file or arbitrary path is
// accepted as a restore source, and restoration always creates a rescue point.
var safeRecoveryName = regexp.MustCompile(`^[0-9A-Za-z._-]{5,95}\.json$`)
var syncRecoveryName = regexp.MustCompile(`^\d{8}T\d{6}\.\d{9}Z?(?:_antes_restaurar|_cambios_pendientes)?\.json$`)

const pendingRecoverySuffix = "_cambios_pendientes"

func (a *runningApp) saveSyncRecovery(state libraryState, suffix string) error {
	dir := filepath.Join(filepath.Dir(a.libraryPath), "sync_recuperacion")
	now := time.Now()
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	if len(data) > libraryFileLimit {
		return errors.New("La copia de recuperación supera el límite de tamaño")
	}
	name := now.UTC().Format("20060102T150405.000000000Z") + suffix + ".json"
	if err = os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	if err = replaceLibraryFile(filepath.Join(dir, name), data); err != nil {
		return err
	}
	names, err := labBackupCandidates(dir, syncRecoveryName)
	if err != nil {
		return nil
	}
	var bytes int64
	kept := 0
	for i := len(names) - 1; i >= 0; i-- {
		file := filepath.Join(dir, names[i])
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		if names[i] != name && (kept >= 20 || bytes+info.Size() > 256<<20 || now.Sub(info.ModTime()) > 30*24*time.Hour) {
			_ = os.Remove(file)
			continue
		}
		kept++
		bytes += info.Size()
	}
	return nil
}

// preservePendingSave runs with the library mutex held. A failed backup
// never changes the original save result.
func (a *runningApp) preservePendingSave(w http.ResponseWriter, state libraryState) {
	if a.saveSyncRecovery(state, pendingRecoverySuffix) == nil {
		w.Header().Set("X-Cine-Recovery", "saved")
	}
}

func readRecoveryFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !stat.Mode().IsRegular() || stat.Size() == 0 || stat.Size() > libraryFileLimit {
		return nil, errors.New("Tamaño de copia de recuperación no válido")
	}
	data, err := io.ReadAll(io.LimitReader(f, libraryFileLimit+1))
	if err == nil && (len(data) == 0 || len(data) > libraryFileLimit) {
		return nil, errors.New("Tamaño de copia de recuperación no válido")
	}
	return data, err
}

type recoveryInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	WhenMS int64  `json:"when_ms"`
	Bytes  int64  `json:"bytes"`
}

func (a *runningApp) localRecoveryList() ([]recoveryInfo, error) {
	root := filepath.Dir(a.libraryPath)
	all := []recoveryInfo{}
	directories := []struct{ prefix, path string }{
		{"rec", filepath.Join(root, "sync_recuperacion")},
		{"day", filepath.Join(root, "copias_automaticas")},
	}
	if protectedLibraryEnabled {
		directories = append(directories, struct{ prefix, path string }{"quick", filepath.Join(root, "copias_rapidas_lab")})
	}
	for _, d := range directories {
		entries, err := os.ReadDir(d.path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if entry.IsDir() || !safeRecoveryName.MatchString(entry.Name()) {
				continue
			}
			info, err := entry.Info()
			if err != nil || info.Size() > libraryFileLimit {
				continue
			}
			all = append(all, recoveryInfo{
				ID: d.prefix + ":" + entry.Name(), Name: entry.Name(),
				WhenMS: info.ModTime().UnixMilli(), Bytes: info.Size(),
			})
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].WhenMS > all[j].WhenMS })
	if len(all) > 50 {
		all = all[:50]
	}
	return all, nil
}

// Solo información resumida de la copia elegida, sin exponer los TXT.
func (a *runningApp) previewLocalRecovery(id string) (map[string]any, error) {
	prefix, name, ok := strings.Cut(id, ":")
	if !ok || !safeRecoveryName.MatchString(name) {
		return nil, errors.New("copia no válida")
	}
	root := filepath.Dir(a.libraryPath)
	dir := ""
	switch prefix {
	case "rec":
		dir = filepath.Join(root, "sync_recuperacion")
	case "day":
		dir = filepath.Join(root, "copias_automaticas")
	case "quick":
		if protectedLibraryEnabled {
			dir = filepath.Join(root, "copias_rapidas_lab")
		}
	}
	if dir == "" {
		return nil, errors.New("origen no permitido")
	}
	data, err := readRecoveryFile(filepath.Join(dir, name))
	if err != nil {
		return nil, err
	}
	var state libraryState
	if json.Unmarshal(data, &state) != nil || validateLibrary(state) != nil {
		return nil, errors.New("La copia no es válida. No es seguro restaurarla")
	}
	favorites := 0
	for _, rec := range state.Reactions {
		if rec.Favorite {
			favorites++
		}
	}
	return map[string]any{
		"lists": len(state.Lists), "seen": syncWatched(state), "favorites": favorites,
		"valid": true, "message": "La copia ha superado la validación de la biblioteca",
	}, nil
}

var errLibraryChanged = errors.New("La biblioteca cambió. Recarga antes de restaurar una copia")

func (a *runningApp) restoreLocalRecovery(id, etag string) (map[string]any, error) {
	prefix, name, ok := strings.Cut(id, ":")
	if !ok || !safeRecoveryName.MatchString(name) {
		return nil, errors.New("copia de recuperación no válida")
	}
	root := filepath.Dir(a.libraryPath)
	dir := ""
	switch prefix {
	case "rec":
		dir = filepath.Join(root, "sync_recuperacion")
	case "day":
		dir = filepath.Join(root, "copias_automaticas")
	case "quick":
		if protectedLibraryEnabled {
			dir = filepath.Join(root, "copias_rapidas_lab")
		}
	default:
		return nil, errors.New("origen de recuperación no permitido")
	}
	filename := filepath.Join(dir, name)
	data, err := readRecoveryFile(filename)
	if err != nil {
		return nil, err
	}
	var restored libraryState
	if err = json.Unmarshal(data, &restored); err != nil {
		return nil, errors.New("la copia JSON no es válida")
	}
	if err = validateLibrary(restored); err != nil {
		return nil, fmt.Errorf("no se ha restaurado: %w", err)
	}
	a.mu.Lock()
	if etag != a.libraryETag() {
		a.mu.Unlock()
		return nil, errLibraryChanged
	}
	before := a.library
	if err = a.saveSyncRecovery(before, "_antes_restaurar"); err != nil {
		a.mu.Unlock()
		return nil, err
	}
	mirrorPending, err := commitLibrary(a.libraryPath, a.savePath, restored)
	if err != nil {
		a.mu.Unlock()
		return nil, err
	}
	a.library = restored
	a.libraryRevision++
	a.syncPreviewToken = ""
	var warnings []string
	if mirrorPending {
		warnings = append(warnings, "Copia TXT pendiente")
	}
	syncReady := true
	if a.syncStore != nil {
		if err = a.syncStore.observe(before, restored); err != nil {
			syncReady = false
			warnings = append(warnings, "Estado de sincronización pendiente de guardar")
			logDiagnostic("Biblioteca restaurada, pero falló el guardado de sus revisiones: " + err.Error())
		}
	}
	if desktopFeatures && a.markManualCloudUpload() != nil {
		warnings = append(warnings, "No se pudo guardar el estado de subida pendiente")
	}
	a.mu.Unlock()
	if !desktopFeatures && syncReady && a.syncStore != nil && a.graph != nil && a.graph.active() {
		if a.syncStore.isPaired(a.graph.pairingID()) {
			a.queueGraphUpload()
		}
	} else if !desktopFeatures && syncReady && a.cloud != nil && a.syncStore != nil {
		folder := a.cloud.status().Folder
		if a.syncStore.isPaired(folder) {
			saved := a.cloud.backup()
			if saved.LastError != "" {
				warnings = append(warnings, saved.LastError)
			} else if err = a.publishDeviceSnapshot(); err != nil {
				warnings = append(warnings, err.Error())
			}
		}
	}
	return map[string]any{
		"restored": true, "lists": len(restored.Lists), "seen": syncWatched(restored), "warning": strings.Join(warnings, " · "),
	}, nil
}
func (a *runningApp) serveSyncSupport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if !a.auth(r) {
		http.Error(w, "No autorizado", 403)
		return
	}
	respond := func(code int, obj any) { w.WriteHeader(code); _ = json.NewEncoder(w).Encode(obj) }
	switch r.URL.Path {
	case "/api/sync/revisions":
		if r.Method != http.MethodGet {
			respond(405, map[string]string{"error": "Método incorrecto"})
			return
		}
		if a.syncStore == nil {
			respond(200, map[string]any{"revisions": map[string]syncRevision{}})
			return
		}
		_, revisions, _ := a.syncStore.parts()
		respond(200, map[string]any{"revisions": revisions})
	case "/api/sync/history":
		if r.Method != http.MethodGet {
			respond(405, map[string]string{"error": "Método incorrecto"})
			return
		}
		if a.syncStore == nil {
			respond(200, map[string]any{"events": []syncHistoryEntry{}})
			return
		}
		respond(200, map[string]any{"events": a.syncStore.history()})
	case "/api/sync/recovery/pending":
		if !desktopFeatures {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPost {
			respond(http.StatusMethodNotAllowed, map[string]string{"error": "Método incorrecto"})
			return
		}
		var state libraryState
		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, libraryFileLimit))
		if err != nil || json.Unmarshal(body, &state) != nil || validateLibrary(state) != nil {
			respond(http.StatusBadRequest, map[string]string{"error": "Copia de recuperación no válida"})
			return
		}
		a.mu.Lock()
		err = a.saveSyncRecovery(state, pendingRecoverySuffix)
		a.mu.Unlock()
		if err != nil {
			respond(http.StatusInternalServerError, map[string]string{"error": "No se pudo guardar la copia de recuperación"})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case "/api/sync/recovery/list":
		if r.Method != http.MethodGet {
			respond(405, map[string]string{"error": "Método incorrecto"})
			return
		}
		copies, err := a.localRecoveryList()
		if err != nil {
			respond(500, map[string]string{"error": err.Error()})
			return
		}
		respond(200, map[string]any{"copies": copies})
	case "/api/sync/recovery/preview":
		if !protectedLibraryEnabled {
			respond(404, map[string]string{"error": "No disponible"})
			return
		}
		if r.Method != http.MethodGet {
			respond(405, map[string]string{"error": "Método incorrecto"})
			return
		}
		info, err := a.previewLocalRecovery(r.URL.Query().Get("id"))
		if err != nil {
			respond(409, map[string]string{"error": err.Error()})
			return
		}
		respond(200, info)
	case "/api/sync/recovery/restore":
		if r.Method != http.MethodPost {
			respond(405, map[string]string{"error": "Método incorrecto"})
			return
		}
		var request struct {
			ID string `json:"id"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&request); err != nil {
			respond(400, map[string]string{"error": "Solicitud incorrecta"})
			return
		}
		restored, err := a.restoreLocalRecovery(request.ID, r.Header.Get("If-Match"))
		if errors.Is(err, errLibraryChanged) {
			respond(http.StatusPreconditionFailed, map[string]string{"error": err.Error()})
			return
		}
		if err != nil {
			respond(409, map[string]string{"error": err.Error()})
			return
		}
		respond(200, restored)
	default:
		http.NotFound(w, r)
	}
}
