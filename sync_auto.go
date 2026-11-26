package main

import (
	"errors"
	"os"
	"reflect"
	"time"
)

// cloudAutoSync imports only conflict-free changes after the user has explicitly
// paired a folder in the manual preview/confirmation flow. No silent initial pairing.
func (a *runningApp) cloudAutoSync() (map[string]any, error) {
	a.syncPlanMu.Lock()
	defer a.syncPlanMu.Unlock()
	if a.cloud == nil || a.syncStore == nil {
		return map[string]any{"changed": false, "skipped": "disabled"}, nil
	}
	folder := a.cloud.status().Folder
	if folder == "" || !a.syncStore.isPaired(folder) {
		return map[string]any{"changed": false, "skipped": "not_paired", "needs_review": folder != ""}, nil
	}
	a.mu.Lock()
	local, revision := a.library, a.libraryRevision
	a.mu.Unlock()
	plan, err := a.cloudPlanFor(local, folder)
	if err != nil {
		return nil, err
	}
	a.mu.Lock()
	if revision != a.libraryRevision || folder != a.cloud.status().Folder {
		a.mu.Unlock()
		return map[string]any{"changed": false, "skipped": "local_changed"}, nil
	}
	a.syncLastCheck.Store(time.Now().UnixMilli())
	if len(plan.Conflicts) != 0 {
		a.syncPendingConflict.Store(true)
		a.mu.Unlock()
		return map[string]any{"changed": false, "needs_review": true, "conflicts": len(plan.Conflicts)}, nil
	}
	a.syncPendingConflict.Store(false)
	if plan.Sources == 0 {
		a.mu.Unlock()
		return map[string]any{"changed": false, "skipped": "no_remote_sources"}, nil
	}
	_, beforeRevisions, beforeKnown := a.syncStore.parts()
	libChanged := !reflect.DeepEqual(a.library, plan.Library)
	metadataChanged := !reflect.DeepEqual(beforeRevisions, plan.Revisions) || !reflect.DeepEqual(beforeKnown, plan.Known)
	if !libChanged && !metadataChanged {
		a.mu.Unlock()
		return map[string]any{"changed": false}, nil
	}
	if err = a.syncStore.checkShareable(plan, nil); err != nil {
		a.mu.Unlock()
		return nil, err
	}
	if libChanged {
		if err = a.saveSyncRecovery(a.library, ""); err != nil {
			a.mu.Unlock()
			return nil, err
		}
		if err = saveLibrary(a.libraryPath, a.savePath, plan.Library); err != nil {
			a.mu.Unlock()
			return nil, err
		}
		a.library = plan.Library
		a.libraryRevision++
		a.syncLastReceived.Store(time.Now().UnixMilli())
	}
	metadataErr := a.syncStore.accept(plan, nil, folder)
	a.syncPreviewToken = ""
	a.mu.Unlock()
	result := syncResult(plan, libChanged, metadataErr, a.publishFolderSync)
	result["metadata_updated"] = metadataChanged && metadataErr == nil
	return result, nil
}

func (a *runningApp) markManualCloudUpload() error {
	if a.syncStore == nil || a.cloud == nil || !(a.cloud.status().Enabled || a.graph != nil && a.graph.active()) {
		return nil
	}
	return a.syncStore.markGraphUploadPending()
}

// Read other devices without applying a plan or publishing a snapshot.
func (a *runningApp) cloudIncomingPending() (bool, error) {
	if a.syncStore == nil || a.cloud == nil {
		return false, errors.New("Nube sin configurar")
	}
	if !a.syncPlanMu.TryLock() {
		return false, errors.New("Hay otra comprobación en curso")
	}
	defer a.syncPlanMu.Unlock()
	a.mu.Lock()
	local, revision := a.library, a.libraryRevision
	a.mu.Unlock()

	var plan syncPlan
	var err error
	if a.graph != nil && a.graph.active() {
		op, cancel := a.graph.operation().bounded()
		defer cancel()
		plan, err = a.graphPlanFor(local, op)
		if err == nil && !op.current() {
			return false, errGraphSessionChanged
		}
	} else {
		folder := a.cloud.status().Folder
		if info, err := os.Stat(folder); err != nil || !info.IsDir() {
			return false, errors.New("Carpeta de nube no disponible")
		}
		plan, err = a.cloudPlanFor(local, folder)
		if err == nil && folder != a.cloud.status().Folder {
			return false, errors.New("La carpeta ha cambiado")
		}
	}
	if err != nil {
		return false, err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if revision != a.libraryRevision {
		return false, errors.New("La biblioteca ha cambiado")
	}
	return len(plan.Conflicts) != 0 || !reflect.DeepEqual(local, plan.Library), nil
}
