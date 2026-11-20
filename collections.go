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
}
type libraryState struct {
	Version       int              `json:"version"`
	ActiveID      string           `json:"active_id"`
	Lists         []filmCollection `json:"lists"`
	ComparisonIDs []string         `json:"comparison_ids,omitempty"`
}

var listIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,60}$`)

func initialLibrary(legacy string) libraryState {
	return libraryState{Version: 1, ActiveID: "atlas", Lists: []filmCollection{{ID: "atlas", Name: "365 películas", Text: legacy}}}
}
func validateLibrary(l libraryState) error {
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
	// Once a day, retain a copy of the previous good collection state.
	backupDir := filepath.Join(filepath.Dir(path), "copias_automaticas")
	if before, e := os.ReadFile(path); e == nil && len(before) > 0 {
		_ = os.MkdirAll(backupDir, 0700)
		backup := filepath.Join(backupDir, time.Now().Format("2006-01-02")+".json")
		if _, e = os.Stat(backup); os.IsNotExist(e) {
			_ = os.WriteFile(backup, before, 0600)
		}
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, bytes, 0600); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		// Handles older Windows installations where replacing an open file fails.
		if err2 := os.WriteFile(path, bytes, 0600); err2 != nil {
			return fmt.Errorf("no se pudo guardar: %w", err2)
		}
		_ = os.Remove(tmp)
	}
	for _, it := range l.Lists {
		if it.ID == "atlas" {
			if err := os.WriteFile(legacyPath, []byte(it.Text), 0600); err != nil {
				return err
			}
			break
		}
	}
	return nil
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
		a.mu.Unlock()
		_ = json.NewEncoder(w).Encode(state)
	case http.MethodPost:
		var state libraryState
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 3<<20)).Decode(&state); err != nil {
			http.Error(w, "JSON inválido", 400)
			return
		}
		if err := validateLibrary(state); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		a.mu.Lock()
		err := saveLibrary(a.libraryPath, a.savePath, state)
		if err == nil {
			a.library = state
		}
		a.mu.Unlock()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if a.cloud != nil {
			cloudState := a.cloud.backup()
			if cloudState.Enabled {
				if cloudState.LastError != "" {
					w.Header().Set("X-Cine-Cloud", "error")
				} else {
					w.Header().Set("X-Cine-Cloud", "ok")
				}
			}
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Método no permitido", 405)
	}
}

func readSafeBytes(r io.Reader, max int64) ([]byte, error) { return io.ReadAll(io.LimitReader(r, max)) }
