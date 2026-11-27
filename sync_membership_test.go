package main

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSyncMembershipRemovalSurvivesOldSnapshotAndReadd(t *testing.T) {
	before := initialLibrary("☐ Maridos (1970)\n☐ Insomnio (2002)")
	before.Lists = append(before.Lists, filmCollection{ID: "other", Name: "Otra lista", Text: "☐ Maridos (1970)"})
	removed := before
	removed.Lists = append([]filmCollection(nil), before.Lists...)
	removed.Lists[0].Text = "☐ Insomnio (2002)"
	recordMembershipEdits(before, &removed)

	// Another device has an old copy of Maridos and a legitimate new film.
	peer := before
	peer.Lists = append([]filmCollection(nil), before.Lists...)
	peer.Lists[0].Text += "\n☐ Origen (2010)"
	merged, _, _, _, err := syncMergeVersioned(removed, nil, nil, syncSnapshot{Protocol: 2, Library: peer}, true)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(merged.Lists[0].Text, "Maridos") || !strings.Contains(merged.Lists[0].Text, "Origen") ||
		!strings.Contains(merged.Lists[1].Text, "Maridos") {
		t.Fatal("La fusión resucitó la pertenencia retirada o perdió otros cambios", merged.Lists)
	}

	// A device on the updated version can receive the removal in the opposite direction.
	received, _, _, _, err := syncMergeVersioned(peer, nil, nil, syncSnapshot{Protocol: 2, Library: merged}, true)
	if err != nil || strings.Contains(received.Lists[0].Text, "Maridos") {
		t.Fatal("La retirada no llegó al otro equipo", err)
	}

	readded := merged
	readded.Lists = append([]filmCollection(nil), merged.Lists...)
	readded.Lists[0].Text += "\n☐ Maridos (1970)"
	recordMembershipEdits(merged, &readded)
	received, _, _, _, err = syncMergeVersioned(readded, nil, nil, syncSnapshot{Protocol: 2, Library: merged}, true)
	if err != nil || !strings.Contains(received.Lists[0].Text, "Maridos") {
		t.Fatal("Una retirada vieja sobrescribió la nueva incorporación", err)
	}
}

func TestDesktopMembershipSaveAndManualSync(t *testing.T) {
	if !desktopFeatures {
		t.Skip("Solo escritorio")
	}
	base := t.TempDir()
	localDir, folder := filepath.Join(base, "local"), filepath.Join(base, "OneDrive")
	if err := os.MkdirAll(folder, 0700); err != nil {
		t.Fatal(err)
	}
	initial := initialLibrary("☐ Maridos (1970)\n☐ Insomnio (2002)")
	initial.Lists = append(initial.Lists, filmCollection{ID: "other", Name: "Otra", Text: "☐ Maridos (1970)"})
	store, err := newDesktopSyncStore(localDir)
	if err != nil {
		t.Fatal(err)
	}
	cloud := newCloudBackups(localDir)
	if err := cloud.configure(folder); err != nil {
		t.Fatal(err)
	}
	app := &runningApp{library: initial, token: "test", cloud: cloud, syncStore: store,
		libraryPath: filepath.Join(localDir, "colecciones.json"), savePath: filepath.Join(localDir, "legacy.txt")}
	preview, err := app.cloudSyncPreview()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := app.cloudSyncConfirm(preview["token"].(string), nil); err != nil {
		t.Fatal(err)
	}
	peer := syncSnapshot{Protocol: 2, DeviceID: "peer-old", Library: initial}
	body, err := json.Marshal(peer)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(cloudSyncRoot(folder), "sincronizacion", "dispositivos", "peer-old.json")
	if err := cloudWriteFile(path, body); err != nil {
		t.Fatal(err)
	}
	local := initial
	local.Lists = append([]filmCollection(nil), initial.Lists...)
	local.Lists[0].Text = "☐ Insomnio (2002)"
	if w := postTestLibrary(t, app, local); w.Code != 204 {
		t.Fatal(w.Body.String())
	}
	// Simulate a reopen before synchronizing, including a legacy UI payload on the next save.
	app.library = loadLibrary(app.libraryPath, "")
	if rev := app.library.MembershipChanges["atlas"]["maridos|1970"]; rev.Present || rev.ChangedAt == 0 {
		t.Fatal("No se persistió la retirada")
	}
	if w := postTestLibrary(t, app, local); w.Code != 204 {
		t.Fatal(w.Body.String())
	}
	req := httptest.NewRequest("POST", "/api/cloud/auto", strings.NewReader(`{"force":true}`))
	req.Header.Set("X-Atlas-Token", app.token)
	w := httptest.NewRecorder()
	app.serveHTTP(w, req)
	if w.Code != 200 || strings.Contains(w.Body.String(), "warning") {
		t.Fatal(w.Body.String())
	}
	saved := loadLibrary(app.libraryPath, "")
	if strings.Contains(saved.Lists[0].Text, "Maridos") || !strings.Contains(saved.Lists[1].Text, "Maridos") {
		t.Fatal("La sincronización cambió una pertenencia incorrecta", saved.Lists)
	}
}
