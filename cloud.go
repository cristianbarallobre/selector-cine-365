package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// A separate backup copy, not a replacement for the user's AppData storage.
// No API keys or poster cache are copied to cloud storage: only collection data.
type cloudSettings struct {
	Folder string `json:"folder"`
}
type cloudStatus struct {
	Enabled    bool   `json:"enabled"`
	Folder     string `json:"folder"`
	Target     string `json:"target"`
	LastBackup string `json:"last_backup,omitempty"`
	LastError  string `json:"last_error,omitempty"`
	SyncPaired bool `json:"sync_paired"`
	NeedsReview bool `json:"needs_review"`
}
type cloudBackups struct {
	mu           sync.Mutex
	sourceDir    string
	settingsFile string
	folder       string
	lastBackup   string
	lastError    string
}

const cloudDirectoryName = "Selector_Cine_365_Copias"

func newCloudBackups(dir string) *cloudBackups {
	c := &cloudBackups{sourceDir: dir, settingsFile: filepath.Join(dir, "copia_nube_ajustes.json")}
	if b, e := os.ReadFile(c.settingsFile); e == nil {
		var s cloudSettings
		if json.Unmarshal(b, &s) == nil && s.Folder != "" {
			c.folder = s.Folder
			if st, err := os.Stat(filepath.Join(cloudSyncRoot(c.folder), "actual", "colecciones.json")); err == nil {
				c.lastBackup = st.ModTime().Format(time.RFC3339)
			}
		}
	}
	return c
}
func (c *cloudBackups) statusLocked() cloudStatus {
	return cloudStatus{Enabled: c.folder != "", Folder: c.folder, Target: func() string {
		if c.folder == "" {
			return ""
		}
		return cloudSyncRoot(c.folder)
	}(), LastBackup: c.lastBackup, LastError: c.lastError}
}
func (c *cloudBackups) status() cloudStatus {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.statusLocked()
}
func (c *cloudBackups) configure(folder string) error {
	folder = strings.Trim(strings.TrimSpace(folder), "\"")
	if folder != "" {
		if !filepath.IsAbs(folder) {
			return errors.New("Selecciona una ruta completa a tu carpeta sincronizada")
		}
		folder = filepath.Clean(folder)
		if strings.EqualFold(folder, filepath.Clean(c.sourceDir)) {
			return errors.New("La carpeta de la copia debe ser distinta de AppData")
		}
		info, err := os.Stat(folder)
		if err != nil || !info.IsDir() {
			return errors.New("La carpeta elegida no existe o no está disponible")
		}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	b, _ := json.MarshalIndent(cloudSettings{Folder: folder}, "", "  ")
	if err := os.MkdirAll(c.sourceDir, 0700); err != nil {
		return err
	}
	if err := os.WriteFile(c.settingsFile+".tmp", b, 0600); err != nil {
		return err
	}
	if err := os.Rename(c.settingsFile+".tmp", c.settingsFile); err != nil {
		// Replacing an existing file can fail on some Windows setups.
		if e := os.WriteFile(c.settingsFile, b, 0600); e != nil {
			return e
		}
		_ = os.Remove(c.settingsFile + ".tmp")
	}
	c.folder = folder
	c.lastError = ""
	if folder == "" {
		c.lastBackup = ""
	}
	return nil
}
func cloudWriteFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		if err2 := os.WriteFile(path, data, 0600); err2 != nil {
			return err2
		}
		_ = os.Remove(tmp)
	}
	return nil
}
func (c *cloudBackups) backup() cloudStatus {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.folder == "" {
		return c.statusLocked()
	}
	// If the provider has been disconnected, never silently create its missing
	// parent: that could make the user believe the files were synced to cloud.
	info, err := os.Stat(c.folder)
	if err != nil || !info.IsDir() {
		c.lastError = "La carpeta de nube no está disponible"
		return c.statusLocked()
	}
	dest := cloudSyncRoot(c.folder)
	if err = os.MkdirAll(filepath.Join(dest, "actual"), 0700); err != nil {
		c.lastError = err.Error()
		return c.statusLocked()
	}
	files := []string{"colecciones.json", "mis_365_peliculas.txt", "identificaciones_tmdb.json", "identidades_cine.json"}
	copied := 0
	for _, filename := range files {
		b, readErr := os.ReadFile(filepath.Join(c.sourceDir, filename))
		if errors.Is(readErr, os.ErrNotExist) {
			continue
		}
		if readErr != nil {
			c.lastError = readErr.Error()
			return c.statusLocked()
		}
		if writeErr := cloudWriteFile(filepath.Join(dest, "actual", filename), b); writeErr != nil {
			c.lastError = writeErr.Error()
			return c.statusLocked()
		}
		// One immutable restore point per calendar day (the first save of the
		// day). Later saves update "actual" without erasing that restore point.
		snap := filepath.Join(dest, "historico", time.Now().Format("2006-01-02"), filename)
		if _, e := os.Stat(snap); errors.Is(e, os.ErrNotExist) {
			if writeErr := cloudWriteFile(snap, b); writeErr != nil {
				c.lastError = writeErr.Error()
				return c.statusLocked()
			}
		}
		copied++
	}
	if copied == 0 {
		c.lastError = "Todavía no hay ninguna lista guardada"
		return c.statusLocked()
	}
	c.lastBackup = time.Now().Format(time.RFC3339)
	c.lastError = ""
	return c.statusLocked()
}
// Detectar dispositivos conectados antes de publicar una copia antigua.
func cloudHasDeviceSnapshots(folder string) bool {
 if folder=="" {return false}
 base:=filepath.Join(cloudSyncRoot(folder),"sincronizacion","dispositivos")
 entries,err:=os.ReadDir(base);if err!=nil{return false}
 for _,entry:=range entries{if !entry.IsDir() && strings.HasSuffix(entry.Name(),".json"){return true}}
 return false
}
func (a *runningApp) cloudDisplayStatus() cloudStatus {
 s:=a.cloud.status()
 s.SyncPaired=a.syncStore!=nil&&a.syncStore.isPaired(s.Folder)
 s.NeedsReview=s.Enabled&&!s.SyncPaired&&cloudHasDeviceSnapshots(s.Folder)
 return s
}
func (a *runningApp) serveCloudHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if !a.auth(r) {
		http.Error(w, "No autorizado", 403)
		return
	}
	respond := func(status int, v any) { w.WriteHeader(status); _ = json.NewEncoder(w).Encode(v) }
	switch r.URL.Path {
	case "/api/cloud/status":
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", 405)
			return
		}
		respond(200, a.cloudDisplayStatus())
	case "/api/cloud/browse":
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", 405)
			return
		}
		folder, err := chooseCloudFolder()
		if err != nil {
			respond(500, map[string]string{"error": err.Error()})
			return
		}
		respond(200, map[string]string{"folder": folder})
	case "/api/cloud/config":
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", 405)
			return
		}
		var payload cloudSettings
		if err := json.NewDecoder(io.LimitReader(r.Body, 8<<10)).Decode(&payload); err != nil {
			respond(400, map[string]string{"error": "Ruta incorrecta"})
			return
		}
		if err := a.cloud.configure(payload.Folder); err != nil {
			respond(400, map[string]string{"error": err.Error()})
			return
		}
		a.mu.Lock();a.syncPreviewToken="";a.mu.Unlock()
        // Si hay un Android registrado, primero revisar: no publicar la
        // biblioteca local como una copia compartida por el simple hecho de conectar.
        if payload.Folder!="" && a.syncStore!=nil && !a.syncStore.isPaired(payload.Folder) && cloudHasDeviceSnapshots(payload.Folder) {
            respond(200,a.cloudDisplayStatus())
            return
        }
		if payload.Folder != "" {
			a.mu.Lock()
			saveErr := saveLibrary(a.libraryPath, a.savePath, a.library)
			a.mu.Unlock()
			if saveErr != nil {
				respond(500, map[string]string{"error": saveErr.Error()})
				return
			}
		}
		respond(200, a.cloud.backup())
	case "/api/cloud/backup":
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", 405)
			return
		}
		folder:=a.cloud.status().Folder
        if a.syncStore!=nil && !a.syncStore.isPaired(folder) && cloudHasDeviceSnapshots(folder) {
            respond(409,map[string]string{"error":"Hay otros dispositivos: revisa la sincronización antes de publicar la copia"})
            return
        }
		s := a.cloud.backup()
		if !s.Enabled {
			respond(400, map[string]string{"error": "Configura antes una carpeta de nube"})
			return
		}
		if s.LastError != "" {
			respond(500, map[string]string{"error": s.LastError})
			return
		}
		if err:=a.publishDeviceSnapshot();err!=nil {
            respond(500,map[string]string{"error":err.Error()});return
        }
		respond(200, a.cloudDisplayStatus())
    case "/api/cloud/preview":
        if r.Method!=http.MethodPost {http.Error(w,"Método no permitido",405);return}
        data,err:=a.cloudSyncPreview()
        if err!=nil{respond(409,map[string]string{"error":err.Error()});return}
        respond(200,data)
    case "/api/cloud/sync":
        if r.Method!=http.MethodPost {http.Error(w,"Método no permitido",405);return}
        var payload struct {
            Token string `json:"token"`
            Decisions map[string]bool `json:"decisions"`
        }
        if err:=json.NewDecoder(io.LimitReader(r.Body,100<<10)).Decode(&payload);err!=nil {
            respond(400,map[string]string{"error":"Solicitud incorrecta"});return
        }
        data,err:=a.cloudSyncConfirm(payload.Token,payload.Decisions)
        if err!=nil{respond(409,map[string]string{"error":err.Error()});return}
        respond(200,data)
	default:
		http.NotFound(w, r)
	}
}
func (c *cloudBackups) String() string {
	s := c.status()
	return fmt.Sprintf("%t %s", s.Enabled, s.Folder)
}
