package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// Lists are stored separately from the legacy 365-film TXT. The legacy TXT is
// mirrored for backwards compatibility with Selector v4/v5.
type filmCollection struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Text string `json:"text"`
	// Updated only on explicit drag-and-drop reordering. Older clients omit it.
	OrderChangedAt int64 `json:"order_changed_at,omitempty"`
}

// Valoración propia y favoritos: datos del usuario, nunca caché TMDB.
// Los registros desactivados se conservan como tombstones para que una
// sincronización antigua no reactive corazones o puntuaciones eliminadas.
type filmReaction struct {
	Favorite  bool   `json:"favorite"`
	Rating    int    `json:"rating,omitempty"`
	TMDBID    int    `json:"tmdb_id,omitempty"`
	Title     string `json:"title"`
	Year      string `json:"year,omitempty"`
	ChangedAt int64  `json:"changed_at"`
}
type favoritesOrder struct {
	Keys      []string `json:"keys"`
	ChangedAt int64    `json:"changed_at"`
}
type libraryState struct {
	Version           int                   `json:"version"`
	ActiveID          string                `json:"active_id"`
	Lists             []filmCollection      `json:"lists"`
	ComparisonIDs     []string              `json:"comparison_ids,omitempty"`
	MembershipChanges syncMembershipChanges `json:"membership_changes,omitempty"`
	// Tombstones compartidas: una lista borrada NO puede volver de una copia atrasada.
	DeletedLists map[string]bool `json:"deleted_lists,omitempty"`
	// Identificaciones TMDB compartibles, nunca la clave API privada.
	TMDBOverrides  map[string]int          `json:"tmdb_overrides,omitempty"`
	FavoritesOrder *favoritesOrder         `json:"favorites_order,omitempty"`
	Reactions      map[string]filmReaction `json:"reactions,omitempty"`
}

var listIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,60}$`)

func initialLibrary(legacy string) libraryState {
	return libraryState{Version: 1, ActiveID: "atlas", Lists: []filmCollection{{ID: "atlas", Name: "Selector Cine 365 · 365", Text: legacy}}}
}
func validateLibrary(l libraryState) error {
	if err := validateMembershipChanges(l.MembershipChanges); err != nil {
		return err
	}
	if l.Version != 1 || len(l.Lists) < 1 || len(l.Lists) > 40 {
		return errors.New("Formato o cantidad de listas no válido")
	}
	ids := map[string]bool{}
	active := false
	for _, it := range l.Lists {
		if !listIDPattern.MatchString(it.ID) || ids[it.ID] || len([]rune(strings.TrimSpace(it.Name))) < 1 || len([]rune(it.Name)) > 100 || !utf8.ValidString(it.Text) || len(it.Text) > 900000 {
			return errors.New("Lista no válida")
		}
		ids[it.ID] = true
		if it.ID == l.ActiveID {
			active = true
		}
	}
	if !active {
		return errors.New("La lista activa no existe")
	}
	for id, deleted := range l.DeletedLists {
		if !listIDPattern.MatchString(id) || !deleted || ids[id] {
			return errors.New("Registro de listas borradas incompatible")
		}
	}
	if len(l.DeletedLists) > 1000 {
		return errors.New("Demasiadas eliminaciones registradas")
	}
	for _, it := range l.Lists {
		if it.OrderChangedAt < 0 || it.OrderChangedAt > 4102444800000 {
			return errors.New("Marca temporal del orden no válida")
		}
	}
	if len(l.TMDBOverrides) > 10000 {
		return errors.New("Demasiadas identificaciones TMDB")
	}
	for key, id := range l.TMDBOverrides {
		if len(key) > 275 || !utf8.ValidString(key) || !strings.Contains(key, "|") || id < 0 || id > 999999999 {
			return errors.New("Identificación TMDB inválida")
		}
	}
	if order := l.FavoritesOrder; order != nil {
		if order.ChangedAt < 0 || order.ChangedAt > 4102444800000 || len(order.Keys) > 15000 {
			return errors.New("Orden de Favoritos inválido")
		}
		seen := make(map[string]bool, len(order.Keys))
		for _, key := range order.Keys {
			if len(key) > 275 || !utf8.ValidString(key) || !strings.Contains(key, "|") || seen[key] {
				return errors.New("Título inválido en el orden de Favoritos")
			}
			seen[key] = true
		}
	}
	if len(l.Reactions) > 15000 {
		return errors.New("Demasiadas valoraciones o favoritos")
	}
	for key, rec := range l.Reactions {
		if !utf8.ValidString(key) || len(key) > 275 || !strings.Contains(key, "|") ||
			rec.Rating < 0 || rec.Rating > 5 || rec.TMDBID < 0 || rec.TMDBID > 999999999 || !utf8.ValidString(rec.Title) ||
			len([]rune(rec.Title)) > 200 || strings.TrimSpace(rec.Title) == "" ||
			len(rec.Year) > 4 || rec.ChangedAt < 0 || rec.ChangedAt > 4102444800000 {
			return errors.New("Valoración o favorito inválido")
		}
		if rec.Year != "" {
			for _, r := range rec.Year {
				if r < '0' || r > '9' {
					return errors.New("Año de favorito no válido")
				}
			}
		}
	}
	if len(l.ComparisonIDs) == 1 || len(l.ComparisonIDs) > len(l.Lists) {
		return errors.New("La comparación requiere al menos dos listas")
	}
	chosen := map[string]bool{}
	for _, id := range l.ComparisonIDs {
		if !ids[id] || chosen[id] {
			return errors.New("La comparación contiene listas no válidas")
		}
		chosen[id] = true
	}
	return nil
}
func loadLibrary(path, legacy string) libraryState {
	contents, err := os.ReadFile(path)
	if err == nil {
		var state libraryState
		if json.Unmarshal(contents, &state) == nil && validateLibrary(state) == nil {
			return state
		}
	}
	return initialLibrary(legacy)
}
func saveLibrary(path, legacyPath string, l libraryState) error {
	_, err := commitLibrary(path, legacyPath, l)
	return err
}

// A successful JSON replacement commits the library for every caller.
// The next save retries a failed compatibility TXT using the latest state.
func commitLibrary(path, legacyPath string, l libraryState) (mirrorPending bool, err error) {
	if err := saveLibraryJSON(path, l); err != nil {
		return false, err
	}
	for _, list := range l.Lists {
		if list.ID == "atlas" {
			if err := replaceLibraryFile(legacyPath, []byte(list.Text)); err != nil {
				logDiagnostic("Biblioteca guardada; copia TXT pendiente: " + err.Error())
				return true, nil
			}
			break
		}
	}
	return false, nil
}

func saveLibraryJSON(path string, l libraryState) error {
	if err := validateLibrary(l); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	if len(bytes) > libraryFileLimit {
		return fmt.Errorf("la biblioteca serializada supera el límite de %d bytes", libraryFileLimit)
	}
	// Si otro proceso dañó la biblioteca mientras la app estaba abierta,
	// no sobrescribir el original en LAB. El siguiente inicio intentará
	// recuperar una copia válida y conservará el archivo defectuoso.
	before, readErr := os.ReadFile(path)
	if protectedLibraryEnabled {
		if readErr == nil && !labValidBackup(before) {
			return errors.New("El archivo de biblioteca ha cambiado o está dañado. No se sobrescribió; reinicia la aplicación para recuperarlo de una copia válida")
		} else if readErr != nil && !os.IsNotExist(readErr) {
			return readErr
		}
	}
	// Once a day, retain a copy of the previous good collection state.
	backupDir := filepath.Join(filepath.Dir(path), "copias_automaticas")
	if readErr == nil && len(before) > 0 {
		// Versiones rápidas solo de un estado anterior válido, con un número
		// máximo de copias. No respaldar accidentalmente archivos corruptos.
		labQuickBackup(path, before)
		_ = os.MkdirAll(backupDir, 0700)
		backup := filepath.Join(backupDir, time.Now().Format("2006-01-02")+".json")
		if _, e := os.Stat(backup); os.IsNotExist(e) {
			_ = os.WriteFile(backup, before, 0600)
		}
	}
	if err := replaceLibraryFile(path, bytes); err != nil {
		return err
	}
	labPruneLibraryBackups(path)
	return nil
}

func replaceLibraryFile(path string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".cine-guardar-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	_, err = tmp.Write(data)
	if err == nil {
		err = tmp.Sync()
	}
	closeErr := tmp.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Rename(tmp.Name(), path)
}
func (a *runningApp) serveLibraryHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if !a.auth(r) {
		http.Error(w, "No autorizado", 403)
		return
	}
	switch r.Method {
	case http.MethodGet:
		a.mu.Lock()
		state := a.library
		w.Header().Set("ETag", a.libraryETag())
		a.mu.Unlock()
		_ = json.NewEncoder(w).Encode(state)
	case http.MethodPost:
		var state libraryState
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, libraryFileLimit)).Decode(&state); err != nil {
			var sizeErr *http.MaxBytesError
			if errors.As(err, &sizeErr) {
				http.Error(w, "La biblioteca supera el tamaño permitido", 413)
			} else {
				http.Error(w, "JSON inválido", 400)
			}
			return
		}
		if err := validateLibrary(state); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		a.mu.Lock()
		if r.Header.Get("If-Match") != a.libraryETag() {
			retry := state
			retry.MembershipChanges = a.library.MembershipChanges
			if r.Header.Get("If-Match") != "" && reflect.DeepEqual(retry, a.library) {
				w.Header().Set("ETag", a.libraryETag())
				a.mu.Unlock()
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if desktopFeatures && r.Header.Get("If-Match") != "" {
				a.preservePendingSave(w, state)
			}
			a.mu.Unlock()
			http.Error(w, "La biblioteca cambió desde que se abrió esta ventana", http.StatusPreconditionFailed)
			return
		}
		before := a.library
		if desktopFeatures {
			recordMembershipEdits(before, &state)
		}
		mirrorPending, err := commitLibrary(a.libraryPath, a.savePath, state)
		diagnosticError("guardado", "commit de biblioteca local", err)
		if mirrorPending {
			w.Header().Set("X-Cine-Legacy", "pending")
		}
		if err == nil {
			a.library = state
			a.libraryRevision++
			w.Header().Set("ETag", a.libraryETag())
			// Registrar cambios del usuario, no solo el snapshot final.
			if a.syncStore != nil {
				if syncErr := a.syncStore.observe(before, state); syncErr != nil {
					diagnosticError("sincronizacion", "registrar cambios locales", syncErr)
					w.Header().Set("X-Cine-Cloud", "error")
				}
			}
			if desktopFeatures {
				if err := a.markManualCloudUpload(); err != nil {
					diagnosticError("sincronizacion", "marcar subida pendiente", err)
					w.Header().Set("X-Cine-Cloud", "error")
				}
			}
			a.syncPreviewToken = "" // una vista previa anterior ya no es válida
		} else if desktopFeatures {
			a.preservePendingSave(w, state)
		}
		a.mu.Unlock()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if desktopFeatures || w.Header().Get("X-Cine-Cloud") == "error" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if a.graph != nil && a.graph.active() {
			if a.syncStore != nil && a.syncStore.isPaired(a.graph.pairingID()) {
				a.queueGraphUpload()
				w.Header().Set("X-Cine-Cloud", "pending")
			} else {
				w.Header().Set("X-Cine-Cloud", "review")
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if a.cloud != nil {
			folder := a.cloud.status().Folder
			if a.syncStore != nil && !a.syncStore.isPaired(folder) && cloudHasDeviceSnapshots(folder) {
				// Un PC nuevo no sobrescribe el respaldo de otro dispositivo.
				w.Header().Set("X-Cine-Cloud", "review")
			} else {
				cloudState := a.cloud.backup()
				if cloudState.Enabled {
					if cloudState.LastError != "" {
						w.Header().Set("X-Cine-Cloud", "error")
					} else if pubErr := a.publishDeviceSnapshot(); pubErr != nil {
						w.Header().Set("X-Cine-Cloud", "error")
					} else {
						w.Header().Set("X-Cine-Cloud", "ok")
					}
				}
			}
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Método no permitido", 405)
	}
}

// libraryETag is read with a.mu held.
func (a *runningApp) libraryETag() string {
	return fmt.Sprintf(`"library-%d"`, a.libraryRevision)
}

func readSafeBytes(r io.Reader, max int64) ([]byte, error) { return io.ReadAll(io.LimitReader(r, max)) }
