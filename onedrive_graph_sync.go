package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

func (a *runningApp) serveGraphAuthHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if !a.auth(r) {
		http.Error(w, "No autorizado", 403)
		return
	}
	respond := func(code int, v any) { w.WriteHeader(code); _ = json.NewEncoder(w).Encode(v) }
	if a.graph == nil {
		respond(503, map[string]string{"error": "Microsoft Graph no está disponible"})
		return
	}
	switch r.URL.Path {
	case "/api/onedrive-desktop/status":
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", 405)
			return
		}
		respond(200, a.graph.status())
	case "/api/onedrive-desktop/start":
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", 405)
			return
		}
		if a.cloud.status().Enabled {
			respond(409, map[string]string{"error": "Ya hay una carpeta configurada. Para cambiar a Graph, exporta primero una copia JSON y desactiva explícitamente la configuración anterior."})
			return
		}
		if a.graph.active() {
			respond(409, map[string]string{"error": "Ya hay una cuenta Microsoft conectada. Desconéctala explícitamente antes de cambiar de cuenta."})
			return
		}
		a.mu.Lock()
		a.syncPreviewToken = ""
		a.mu.Unlock()
		data, err := a.graph.begin()
		if err != nil {
			respond(409, map[string]string{"error": err.Error()})
			return
		}
		if err = openDesktopMicrosoftBrowser(stringValue(data, "login_url")); err != nil {
			data["browser_warning"] = err.Error()
		}
		respond(200, data)
	case "/api/onedrive-desktop/disconnect":
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", 405)
			return
		}
		if err := a.graph.logout(); err != nil {
			respond(409, map[string]string{"error": err.Error()})
			return
		}
		a.mu.Lock()
		a.syncPreviewToken = ""
		a.mu.Unlock()
		respond(200, a.graph.status())
	default:
		http.NotFound(w, r)
	}
}

func (g graphOperation) existingBase() (string, error) {
	root := g.pairingID()
	if root == "" {
		return "", errors.New("Conecta primero OneDrive")
	}
	id := strings.TrimPrefix(root, "onedrive://item/")
	for _, name := range []string{"actual", "sincronizacion"} {
		item, exists, err := g.child(id, name)
		if err != nil {
			return "", err
		}
		if exists && item.Folder != nil {
			return id, nil
		}
	}
	item, exists, err := g.child(id, cloudDirectoryName)
	if err != nil {
		return "", err
	}
	if exists {
		if item.Folder == nil {
			return "", errors.New("La carpeta de copias de OneDrive está ocupada por un archivo")
		}
		return item.ID, nil
	}
	return "", nil
}
func (g graphOperation) ensureBase() (string, error) {
	base, err := g.existingBase()
	if err != nil || base != "" {
		return base, err
	}
	root := g.root
	return g.ensureDir(root, cloudDirectoryName)
}
func (g graphOperation) childFolder(parent, name string) (string, error) {
	if parent == "" {
		return "", nil
	}
	item, exists, err := g.child(parent, name)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", nil
	}
	if item.Folder == nil {
		return "", fmt.Errorf("%s no es una carpeta de OneDrive", name)
	}
	return item.ID, nil
}
func (a *runningApp) graphPlanFor(local libraryState, op graphOperation) (syncPlan, error) {
	if op.root == "" {
		return syncPlan{}, errors.New("Conecta primero tu cuenta de OneDrive")
	}
	if a.syncStore == nil {
		return syncPlan{}, errors.New("Sincronización no inicializada")
	}
	if err := validateLibrary(local); err != nil {
		return syncPlan{}, err
	}
	pairing := op.pairingID()
	id, revision, known := a.syncStore.parts()
	strict := a.syncStore.isPaired(pairing)
	p := syncPlan{Library: local, Revisions: revision, Known: known, Conflicts: []syncConflict{}}
	h := sha256.New()
	// The preview belongs to this exact Microsoft drive, never another account.
	fmt.Fprintf(h, "graph-pairing:%s:%d|", pairing, op.generation)
	localBytes, _ := json.Marshal(local)
	h.Write(localBytes)
	process := func(name string, data []byte, legacy bool) error {
		sum := sha256.Sum256(data)
		fmt.Fprintf(h, "%s:%x|", name, sum)
		var peer syncSnapshot
		if legacy {
			peer.Protocol = 1
			peer.DeviceID = "windows-legacy"
			if err := json.Unmarshal(data, &peer.Library); err != nil {
				return fmt.Errorf("biblioteca antigua corrupta: %w", err)
			}
			p.Legacy = true
		} else {
			if err := json.Unmarshal(data, &peer); err != nil {
				return fmt.Errorf("snapshot remoto corrupto (%s): %w", name, err)
			}
			if peer.DeviceID == id {
				return nil
			}
		}
		merged, rev, k, issues, err := syncMergeVersioned(p.Library, p.Revisions, p.Known, peer, strict)
		if err != nil {
			return err
		}
		p.Library = merged
		p.Revisions = rev
		p.Known = k
		p.Sources++
		seen := map[string]bool{}
		for _, issue := range p.Conflicts {
			seen[issue.Key] = true
		}
		for _, issue := range issues {
			if !seen[issue.Key] {
				p.Conflicts = append(p.Conflicts, issue)
				seen[issue.Key] = true
			}
		}
		return nil
	}
	base, err := op.existingBase()
	if err != nil {
		return p, err
	}
	if base != "" {
		if !strict {
			actual, err := op.childFolder(base, "actual")
			if err != nil {
				return p, err
			}
			if actual != "" {
				item, exists, err := op.child(actual, "colecciones.json")
				if err != nil {
					return p, err
				}
				if exists {
					raw, err := op.readItem(item.ID)
					if err != nil {
						return p, err
					}
					if err = process("actual/colecciones.json", raw, true); err != nil {
						return p, err
					}
				}
			}
		}
		syncDir, err := op.childFolder(base, "sincronizacion")
		if err != nil {
			return p, err
		}
		devices, err := op.childFolder(syncDir, "dispositivos")
		if err != nil {
			return p, err
		}
		if devices != "" {
			list, err := op.list("me/drive/items/" + url.PathEscape(devices) + "/children?$select=id,name,folder,file&$top=200")
			if err != nil {
				return p, err
			}
			if len(list) > 100 {
				return p, errors.New("Demasiadas copias de dispositivo; no se sincronizarán datos parciales")
			}
			for _, item := range list {
				if item.File == nil || !strings.HasSuffix(item.Name, ".json") {
					continue
				}
				raw, err := op.readItem(item.ID)
				if err != nil {
					return p, err
				}
				if err = process("sincronizacion/dispositivos/"+item.Name, raw, false); err != nil {
					return p, err
				}
			}
		}
	}
	p.Token = hex.EncodeToString(h.Sum(nil))
	return p, nil
}
func (a *runningApp) graphSyncPreview() (map[string]any, error) {
	if a.graph == nil {
		return nil, errors.New("Conecta primero tu cuenta de OneDrive")
	}
	a.syncPlanMu.Lock()
	defer a.syncPlanMu.Unlock()
	a.mu.Lock()
	local, revision := a.library, a.libraryRevision
	a.mu.Unlock()
	op, cancel := a.graph.operation().bounded()
	defer cancel()
	p, err := a.graphPlanFor(local, op)
	if err != nil {
		return nil, err
	}
	if err = a.syncStore.checkShareable(p, nil); err != nil {
		return nil, err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if revision != a.libraryRevision || !op.current() {
		return nil, errors.New("La biblioteca o la cuenta ha cambiado. Vuelve a revisar la sincronización")
	}
	a.syncPreviewToken = p.Token
	return map[string]any{
		"token": p.Token, "sources": p.Sources, "legacy": p.Legacy,
		"local_lists": len(a.library.Lists), "local_seen": syncWatched(a.library),
		"result_lists": len(p.Library.Lists), "result_seen": syncWatched(p.Library),
		"conflicts": p.Conflicts, "requires_confirmation": true,
	}, nil
}
func (a *runningApp) graphSyncConfirm(token string, decisions map[string]bool) (map[string]any, error) {
	if a.graph == nil {
		return nil, errors.New("Conecta primero tu cuenta de OneDrive")
	}
	a.syncPlanMu.Lock()
	defer a.syncPlanMu.Unlock()
	a.mu.Lock()
	if token == "" || token != a.syncPreviewToken {
		a.mu.Unlock()
		return nil, errors.New("Revisa primero los datos de OneDrive")
	}
	local, revision := a.library, a.libraryRevision
	a.mu.Unlock()
	op, cancel := a.graph.operation().bounded()
	defer cancel()
	p, err := a.graphPlanFor(local, op)
	if err != nil {
		return nil, err
	}
	a.mu.Lock()
	if token != p.Token || token != a.syncPreviewToken || revision != a.libraryRevision || !op.current() {
		a.mu.Unlock()
		return nil, errors.New("OneDrive o la biblioteca ha cambiado desde la vista previa. Revisa otra vez")
	}
	for _, issue := range p.Conflicts {
		if _, ok := decisions[issue.Key]; !ok {
			a.mu.Unlock()
			return nil, fmt.Errorf("Falta resolver el conflicto %s", issue.Key)
		}
	}
	p.Library = syncApplyDecisions(p.Library, decisions)
	if err = a.syncStore.checkShareable(p, decisions); err != nil {
		a.mu.Unlock()
		return nil, err
	}
	if err = a.saveSyncRecovery(a.library, ""); err != nil {
		a.mu.Unlock()
		return nil, err
	}
	if err = saveLibrary(a.libraryPath, a.savePath, p.Library); err != nil {
		a.mu.Unlock()
		return nil, err
	}
	a.library = p.Library
	a.libraryRevision++
	metadataErr := a.syncStore.accept(p, decisions, op.pairingID())
	a.syncPreviewToken = ""
	a.syncPendingConflict.Store(false)
	a.syncLastCheck.Store(time.Now().UnixMilli())
	a.syncLastReceived.Store(time.Now().UnixMilli())
	a.mu.Unlock()
	return syncResult(p, true, metadataErr, func() error { return a.publishGraphOperation(op) }), nil
}
func (a *runningApp) graphPublishSnapshot() error {
	if a.graph == nil {
		return errors.New("OneDrive no está conectado")
	}
	op, cancel := a.graph.operation().bounded()
	defer cancel()
	return a.publishGraphOperation(op)
}
func (a *runningApp) publishGraphOperation(op graphOperation) (err error) {
	g := op.graph
	g.publishMu.Lock()
	defer g.publishMu.Unlock()
	if !op.current() {
		return errGraphSessionChanged
	}
	pairing := op.pairingID()
	if pairing == "" {
		return errors.New("OneDrive no está conectado")
	}
	if a.syncStore == nil || !a.syncStore.isPaired(pairing) {
		return errors.New("Revisa y confirma primero la sincronización")
	}
	a.mu.Lock()
	err = a.syncStore.markGraphUploadPending()
	snapshot := a.syncStore.snapshot(a.library)
	revision := a.libraryRevision
	version := a.syncStore.graphPublicationVersion()
	a.mu.Unlock()
	defer func() {
		if err == nil {
			a.mu.Lock()
			if !op.current() {
				err = errGraphSessionChanged
			} else if revision == a.libraryRevision {
				err = a.syncStore.acknowledgeGraphUpload(version)
			}
			a.mu.Unlock()
		}
		a.graph.mu.Lock()
		defer a.graph.mu.Unlock()
		if g.generation != op.generation || g.rootID != op.root {
			return
		}
		if err == nil {
			a.graph.lastError = ""
			a.graph.uploadRetryDelay = 0
			a.graph.uploadRetryAt = time.Time{}
			return
		}
		a.graph.lastError = err.Error()
		a.graph.uploadRetryDelay = min(5*time.Minute, max(15*time.Second, 2*a.graph.uploadRetryDelay))
		a.graph.uploadRetryAt = time.Now().Add(a.graph.uploadRetryDelay)
	}()
	if err != nil {
		return err
	}
	bytes, err := encodeSyncSnapshot(snapshot)
	if err != nil {
		return err
	}
	base, err := op.ensureBase()
	if err != nil {
		return err
	}
	syncDir, err := op.ensureDir(base, "sincronizacion")
	if err != nil {
		return err
	}
	devices, err := op.ensureDir(syncDir, "dispositivos")
	if err != nil {
		return err
	}
	name := snapshot.DeviceID + ".json"
	old, exists, err := op.child(devices, name)
	if err != nil {
		return err
	}
	// The snapshot contains a timestamp, but an identical payload still needs no PUT.
	if exists {
		previous, err := op.readItem(old.ID)
		if err != nil {
			return err
		}
		var previousSnapshot syncSnapshot
		if json.Unmarshal(previous, &previousSnapshot) == nil {
			currentSnapshot := snapshot
			previousSnapshot.UpdatedAt = 0
			currentSnapshot.UpdatedAt = 0
			if reflect.DeepEqual(previousSnapshot, currentSnapshot) {
				return nil
			}
		}
	}
	if err = op.writeItem(devices, name, bytes); err != nil {
		return err
	}
	actual, exists, err := op.child(devices, name)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("OneDrive no confirmó la copia del dispositivo")
	}
	confirm, err := op.readItem(actual.ID)
	if err != nil {
		return err
	}
	if sha256.Sum256(confirm) != sha256.Sum256(bytes) {
		return errors.New("OneDrive no confirmó el contenido de la copia")
	}
	return nil
}
func (a *runningApp) graphAutoSync(force bool) (map[string]any, error) {
	a.syncPlanMu.Lock()
	defer a.syncPlanMu.Unlock()
	if a.graph == nil || !a.graph.active() || a.syncStore == nil {
		return map[string]any{"changed": false, "skipped": "disconnected"}, nil
	}
	op, cancel := a.graph.operation().bounded()
	defer cancel()
	pairing := op.pairingID()
	if !a.syncStore.isPaired(pairing) {
		return map[string]any{"changed": false, "skipped": "not_paired", "needs_review": true}, nil
	}
	a.graph.mu.Lock()
	if !force && time.Since(a.graph.lastCheck) < 90*time.Second {
		a.graph.mu.Unlock()
		return map[string]any{"changed": false, "skipped": "cooldown"}, nil
	}
	a.graph.lastCheck = time.Now()
	a.graph.mu.Unlock()
	a.mu.Lock()
	local, revision := a.library, a.libraryRevision
	a.mu.Unlock()
	p, err := a.graphPlanFor(local, op)
	if err != nil {
		return nil, err
	}
	a.mu.Lock()
	if revision != a.libraryRevision || !op.current() {
		a.mu.Unlock()
		return map[string]any{"changed": false, "skipped": "local_changed"}, nil
	}
	a.syncLastCheck.Store(time.Now().UnixMilli())
	if len(p.Conflicts) > 0 {
		a.syncPendingConflict.Store(true)
		a.mu.Unlock()
		return map[string]any{"changed": false, "needs_review": true, "conflicts": len(p.Conflicts)}, nil
	}
	a.syncPendingConflict.Store(false)
	if p.Sources == 0 {
		a.mu.Unlock()
		result := map[string]any{"changed": false, "skipped": "no_remote_sources"}
		a.retryGraphPublication(op, force, result)
		return result, nil
	}
	_, oldRev, oldKnown := a.syncStore.parts()
	libChanged := !reflect.DeepEqual(a.library, p.Library)
	metadataChanged := !reflect.DeepEqual(oldRev, p.Revisions) || !reflect.DeepEqual(oldKnown, p.Known)
	if !libChanged && !metadataChanged {
		a.mu.Unlock()
		result := map[string]any{"changed": false}
		a.retryGraphPublication(op, force, result)
		return result, nil
	}
	if err = a.syncStore.checkShareable(p, nil); err != nil {
		a.mu.Unlock()
		return nil, err
	}
	if libChanged {
		if err = a.saveSyncRecovery(a.library, ""); err != nil {
			a.mu.Unlock()
			return nil, err
		}
		if err = saveLibrary(a.libraryPath, a.savePath, p.Library); err != nil {
			a.mu.Unlock()
			return nil, err
		}
		a.library = p.Library
		a.libraryRevision++
		a.syncLastReceived.Store(time.Now().UnixMilli())
	}
	metadataErr := a.syncStore.accept(p, nil, pairing)
	a.syncPreviewToken = ""
	a.mu.Unlock()
	result := syncResult(p, libChanged, metadataErr, func() error { return a.publishGraphOperation(op) })
	result["metadata_updated"] = metadataChanged && metadataErr == nil
	return result, nil
}

func (a *runningApp) graphUploadDue() bool {
	if !a.syncStore.graphUploadPending() {
		return false
	}
	a.graph.mu.Lock()
	defer a.graph.mu.Unlock()
	return !time.Now().Before(a.graph.uploadRetryAt)
}

func (a *runningApp) retryGraphPublication(op graphOperation, force bool, result map[string]any) {
	if !force && !a.graphUploadDue() {
		return
	}
	if err := a.publishGraphOperation(op); err != nil {
		result["warning"] = "Subida pendiente: " + err.Error()
	}
}

func (a *runningApp) queueGraphUpload() {
	if err := a.syncStore.markGraphUploadPending(); err != nil {
		diagnosticError("sincronizacion", "estado de subida pendiente", err)
	}
	if a.graphUploadCh == nil || a.cloud == nil || !a.cloud.autoMode() {
		return
	}
	select {
	case a.graphUploadCh <- struct{}{}:
	default:
	}
}

// One small batch after edits: serialize uploads, no per-keystroke network,
// never publish before the user has reviewed the initial pairing.
func (a *runningApp) graphPublisher() {
	retry := time.NewTicker(15 * time.Second)
	defer retry.Stop()
	for {
		select {
		case <-a.stopCh:
			return
		case <-retry.C:
		case <-a.graphUploadCh:
			timer := time.NewTimer(3 * time.Second)
			select {
			case <-timer.C:
			case <-a.stopCh:
				timer.Stop()
				return
			}
			select {
			case <-a.graphUploadCh:
			default:
			}
		}
		if a.cloud == nil || !a.cloud.autoMode() || a.graph == nil || !a.graph.active() ||
			a.syncStore == nil || !a.syncStore.isPaired(a.graph.pairingID()) || !a.graphUploadDue() {
			continue
		}
		_ = a.graphPublishSnapshot()
	}
}

func (a *runningApp) serveGraphCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if a.graph == nil || r.Host != a.graph.callbackHost() {
		http.Error(w, "Retorno de Microsoft no reconocido", 400)
		return
	}
	query := r.URL.Query()
	err := a.graph.finish(query.Get("state"), query.Get("code"), query.Get("error"))
	if err != nil {
		w.WriteHeader(400)
		_, _ = fmt.Fprintf(w, "<!doctype html><html lang=es><meta charset=utf-8><title>Selector Cine 365</title><body><h2>No se pudo conectar OneDrive</h2><p>%s</p><p>Vuelve a Selector Cine 365.</p></body></html>", html.EscapeString(err.Error()))
		return
	}
	a.mu.Lock()
	a.syncPreviewToken = ""
	a.mu.Unlock()
	_, _ = w.Write([]byte("<!doctype html><html lang=es><meta charset=utf-8><title>Selector Cine 365</title><body><h2>Cuenta Microsoft conectada</h2><p>Regresa a Selector Cine 365 para revisar y confirmar la primera sincronización.</p></body></html>"))
}
