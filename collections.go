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
    // Updated only on explicit drag-and-drop reordering. Older clients omit it.
    OrderChangedAt int64 `json:"order_changed_at,omitempty"`
}
// Valoración propia y favoritos: datos del usuario, nunca caché TMDB.
// Los registros desactivados se conservan como tombstones para que una
// sincronización antigua no reactive corazones o puntuaciones eliminadas.
type filmReaction struct {
 Favorite bool `json:"favorite"`
 Rating int `json:"rating,omitempty"`
 TMDBID int `json:"tmdb_id,omitempty"`
 Title string `json:"title"`
 Year string `json:"year,omitempty"`
 ChangedAt int64 `json:"changed_at"`
}
type libraryState struct {
	Version       int              `json:"version"`
	ActiveID      string           `json:"active_id"`
	Lists         []filmCollection `json:"lists"`
	ComparisonIDs []string         `json:"comparison_ids,omitempty"`
	// Tombstones compartidas: una lista borrada NO puede volver de una copia atrasada.
	DeletedLists map[string]bool   `json:"deleted_lists,omitempty"`
	// Identificaciones TMDB compartibles, nunca la clave API privada.
	TMDBOverrides map[string]int   `json:"tmdb_overrides,omitempty"`
 Reactions map[string]filmReaction `json:"reactions,omitempty"`
}

var listIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,60}$`)

func initialLibrary(legacy string) libraryState {
	return libraryState{Version: 1, ActiveID: "atlas", Lists: []filmCollection{{ID: "atlas", Name: "Selector Cine 365 · 365", Text: legacy}}}
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
	for id,deleted:=range l.DeletedLists {
		if !listIDPattern.MatchString(id)||!deleted||ids[id] {return errors.New("Registro de listas borradas incompatible")}
	}
	if len(l.DeletedLists)>1000 {return errors.New("Demasiadas eliminaciones registradas")}
	for _,it:=range l.Lists {if it.OrderChangedAt<0||it.OrderChangedAt>4102444800000{return errors.New("Marca temporal del orden no válida")}}
	if len(l.TMDBOverrides)>10000 {return errors.New("Demasiadas identificaciones TMDB")}
	for key,id:=range l.TMDBOverrides {
		if len(key)>275||!utf8.ValidString(key)||!strings.Contains(key,"|")||id<0||id>999999999 {return errors.New("Identificación TMDB inválida")}
	}
 if len(l.Reactions)>15000{return errors.New("Demasiadas valoraciones o favoritos")}
 for key,rec:=range l.Reactions {
  if !utf8.ValidString(key)||len(key)>275||!strings.Contains(key,"|")||
   rec.Rating<0||rec.Rating>5||rec.TMDBID<0||rec.TMDBID>999999999||!utf8.ValidString(rec.Title)||
   len([]rune(rec.Title))>200||strings.TrimSpace(rec.Title)==""||
   len(rec.Year)>4||rec.ChangedAt<0||rec.ChangedAt>4102444800000 {
   return errors.New("Valoración o favorito inválido")
  }
  if rec.Year!=""{
   for _,r:=range rec.Year{if r<'0'||r>'9'{return errors.New("Año de favorito no válido")}}
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
	// Si otro proceso dañó la biblioteca mientras la app estaba abierta,
// no sobrescribir el original en LAB. El siguiente inicio intentará
// recuperar una copia válida y conservará el archivo defectuoso.
if appIsLab {
 if prior,e:=os.ReadFile(path);e==nil&&!labValidBackup(prior){
  return errors.New("El archivo de biblioteca ha cambiado o está dañado. No se sobrescribió; reinicia Windows LAB para recuperarlo de una copia válida")
 }else if e!=nil&&!os.IsNotExist(e){return e}
}
// Once a day, retain a copy of the previous good collection state.
	backupDir := filepath.Join(filepath.Dir(path), "copias_automaticas")
	if before, e := os.ReadFile(path); e == nil && len(before) > 0 {
        // Versiones rápidas solo de un estado anterior válido, con un número
        // máximo de copias. No respaldar accidentalmente archivos corruptos.
        labQuickBackup(path,before)
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
	labPruneLibraryBackups(path)
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
        before:=a.library
		err := saveLibrary(a.libraryPath, a.savePath, state)
		if err == nil {
			a.library = state
            // Registrar cambios del usuario, no solo el snapshot final.
            if a.syncStore!=nil {
                if syncErr:=a.syncStore.observe(before,state);syncErr!=nil{
                    w.Header().Set("X-Cine-Cloud","error")
                }
            }
            a.syncPreviewToken="" // una vista previa anterior ya no es válida
		}
		a.mu.Unlock()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
        if a.graph!=nil&&a.graph.active(){
            if a.syncStore!=nil&&a.syncStore.isPaired(a.graph.pairingID()){
                a.queueGraphUpload()
            }
            // Pending means local save succeeded but the cloud write is asynchronous.
            w.Header().Set("X-Cine-Cloud","pending")
            w.WriteHeader(http.StatusNoContent)
            return
        }
		if a.cloud != nil {
            folder:=a.cloud.status().Folder
            if a.syncStore!=nil&&!a.syncStore.isPaired(folder)&&cloudHasDeviceSnapshots(folder){
                // Un PC nuevo no sobrescribe el respaldo de otro dispositivo.
                w.Header().Set("X-Cine-Cloud","pending")
            }else{
			    cloudState := a.cloud.backup()
			    if cloudState.Enabled {
				    if cloudState.LastError != "" {
					    w.Header().Set("X-Cine-Cloud", "error")
				    } else if pubErr:=a.publishDeviceSnapshot();pubErr!=nil {
                        w.Header().Set("X-Cine-Cloud","error")
                    }else{
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

func readSafeBytes(r io.Reader, max int64) ([]byte, error) { return io.ReadAll(io.LimitReader(r, max)) }
