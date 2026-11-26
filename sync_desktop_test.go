package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSyncMergeRemoteWatchedAndLocalPending(t *testing.T) {
	local := initialLibrary("☐ 001. Ciudadano Kane (1941)\n☐ 002. Casablanca (1942)")
	remote := initialLibrary("☑ 001. Ciudadano Kane (1941)\n☐ 002. Casablanca (1942)")
	peer := syncSnapshot{Protocol: 2, DeviceID: "android-test", Library: remote, Known: map[string]int64{}, Revisions: map[string]syncRevision{}}
	merged, _, _, conflicts, err := syncMerge(local, map[string]syncRevision{}, map[string]int64{}, peer)
	if err != nil {
		t.Fatal(err)
	}
	if len(conflicts) != 0 || syncWatched(merged) != 1 {
		t.Fatalf("merge: %v, watched %d", conflicts, syncWatched(merged))
	}
	if !strings.Contains(merged.Lists[0].Text, "☑ 001.") {
		t.Fatal("no se importó la marca de Android")
	}
}
func TestSyncMergeCausalUnwatch(t *testing.T) {
	l := initialLibrary("☑ 001. Casablanca (1942)")
	remote := initialLibrary("☐ 001. Casablanca (1942)")
	key := "casablanca|1942"
	localRev := syncRevision{Checked: true, Device: "android", Seq: 1, Clock: map[string]int64{"android": 1}}
	remoteRev := syncRevision{Checked: false, Device: "desktop", Seq: 2, Clock: map[string]int64{"android": 1, "desktop": 2}}
	peer := syncSnapshot{Protocol: 2, DeviceID: "desktop", Library: remote, Known: map[string]int64{"android": 1, "desktop": 2}, Revisions: map[string]syncRevision{key: remoteRev}}
	result, rev, _, conflicts, err := syncMerge(l, map[string]syncRevision{key: localRev}, map[string]int64{"android": 1}, peer)
	if err != nil {
		t.Fatal(err)
	}
	if len(conflicts) != 0 || syncWatched(result) != 0 || rev[key].Checked {
		t.Fatalf("causal newer unwatch lost: %v, %+v", conflicts, rev)
	}
}
func TestSyncMergeConcurrentOpposite(t *testing.T) {
	l := initialLibrary("☑ 001. Casablanca (1942)")
	remote := initialLibrary("☐ 001. Casablanca (1942)")
	key := "casablanca|1942"
	left := syncRevision{Checked: true, Device: "desktop", Seq: 1, Clock: map[string]int64{"desktop": 1}}
	right := syncRevision{Checked: false, Device: "android", Seq: 1, Clock: map[string]int64{"android": 1}}
	p := syncSnapshot{Protocol: 2, DeviceID: "android", Library: remote, Known: map[string]int64{"android": 1}, Revisions: map[string]syncRevision{key: right}}
	_, _, _, issues, err := syncMerge(l, map[string]syncRevision{key: left}, map[string]int64{"desktop": 1}, p)
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 || issues[0].Key != key {
		t.Fatalf("missing conflict: %+v", issues)
	}
}
func TestSyncPreviewAndPublishPreserveAndroidSnapshot(t *testing.T) {
	root := t.TempDir()
	localDir := filepath.Join(root, "desktop")
	cloudDir := filepath.Join(root, "shared")
	if err := os.MkdirAll(cloudDir, 0700); err != nil {
		t.Fatal(err)
	}
	local := initialLibrary("☐ 001. Ciudadano Kane (1941)")
	path := filepath.Join(localDir, "colecciones.json")
	if err := saveLibrary(path, filepath.Join(localDir, "mis_365_peliculas.txt"), local); err != nil {
		t.Fatal(err)
	}
	store, err := newDesktopSyncStore(localDir)
	if err != nil {
		t.Fatal(err)
	}
	cloud := newCloudBackups(localDir)
	if err := cloud.configure(cloudDir); err != nil {
		t.Fatal(err)
	}
	android := syncSnapshot{Protocol: 2, DeviceID: "android-test", Library: initialLibrary("☑ 001. Ciudadano Kane (1941)"), Known: map[string]int64{"android-test": 1}, Revisions: map[string]syncRevision{"ciudadano kane|1941": {Checked: true, Device: "android-test", Seq: 1, Clock: map[string]int64{"android-test": 1}}}}
	b, _ := json.Marshal(android)
	target := filepath.Join(cloudDir, cloudDirectoryName, "sincronizacion", "dispositivos", "android-test.json")
	if err := cloudWriteFile(target, b); err != nil {
		t.Fatal(err)
	}
	app := &runningApp{cloud: cloud, syncStore: store, library: local, libraryPath: path, savePath: filepath.Join(localDir, "mis_365_peliculas.txt")}
	preview, err := app.cloudSyncPreview()
	if err != nil {
		t.Fatal(err)
	}
	if preview["result_seen"] != 1 {
		t.Fatalf("expected one seen, got %v", preview)
	}
	if _, err = app.cloudSyncConfirm("wrong-token", map[string]bool{}); err == nil {
		t.Fatal("accepted unreviewed sync")
	}
	if _, err = app.cloudSyncConfirm(preview["token"].(string), map[string]bool{}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(b) {
		t.Fatal("Windows overwrote Android snapshot")
	}
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(saved), "☑ 001.") {
		t.Fatal("desktop did not import Android changes")
	}
	if !store.isPaired(cloudDir) {
		t.Fatal("desktop not paired")
	}
	id, _, _ := store.parts()
	if _, err := os.Stat(filepath.Join(cloudDir, cloudDirectoryName, "sincronizacion", "dispositivos", id+".json")); err != nil {
		t.Fatalf("desktop snapshot missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(localDir, "sync_recuperacion")); err != nil {
		t.Fatal("missing rescue copies")
	}
}

func TestDesktopSaveDoesNotOverwriteUnpairedRemote(t *testing.T) {
	root := t.TempDir()
	localDir := filepath.Join(root, "pc")
	shared := filepath.Join(root, "nube")
	if err := os.MkdirAll(shared, 0700); err != nil {
		t.Fatal(err)
	}
	before := initialLibrary("☐ 001. Casablanca (1942)")
	if err := saveLibrary(filepath.Join(localDir, "colecciones.json"), filepath.Join(localDir, "mis_365_peliculas.txt"), before); err != nil {
		t.Fatal(err)
	}
	store, err := newDesktopSyncStore(localDir)
	if err != nil {
		t.Fatal(err)
	}
	cloud := newCloudBackups(localDir)
	if err := cloud.configure(shared); err != nil {
		t.Fatal(err)
	}
	rootCloud := filepath.Join(shared, cloudDirectoryName)
	peerFile := filepath.Join(rootCloud, "sincronizacion", "dispositivos", "android-other.json")
	if err := cloudWriteFile(peerFile, []byte("{\"device_id\":\"android-other\",\"protocol\":2}")); err != nil {
		t.Fatal(err)
	}
	actual := filepath.Join(rootCloud, "actual", "colecciones.json")
	remote := []byte("{\"remote\":\"no tocar\"}")
	if err := cloudWriteFile(actual, remote); err != nil {
		t.Fatal(err)
	}
	app := &runningApp{libraryPath: filepath.Join(localDir, "colecciones.json"), savePath: filepath.Join(localDir, "mis_365_peliculas.txt"), library: before, cloud: cloud, syncStore: store, token: "test"}
	after := initialLibrary("☑ 001. Casablanca (1942)")
	payload, _ := json.Marshal(after)
	req := httptest.NewRequest(http.MethodPost, "/api/library", strings.NewReader(string(payload)))
	req.Header.Set("X-Atlas-Token", "test")
	req.Header.Set("If-Match", app.libraryETag())
	w := httptest.NewRecorder()
	app.serveHTTP(w, req)
	if w.Code != 204 {
		t.Fatalf("local save failed: %d %s", w.Code, w.Body.String())
	}
	if desktopFeatures && !store.graphUploadPending() {
		t.Fatal("No quedó registrada la subida pendiente")
	}
	if !desktopFeatures && w.Header().Get("X-Cine-Cloud") != "review" {
		t.Fatalf("expected pending review, got %s", w.Header().Get("X-Cine-Cloud"))
	}
	saved, _ := os.ReadFile(actual)
	if string(saved) != string(remote) {
		t.Fatal("unpaired Windows replaced shared legacy data")
	}
	savedPeer, _ := os.ReadFile(peerFile)
	if !strings.Contains(string(savedPeer), "android-other") {
		t.Fatal("remote snapshot changed")
	}
	savedLocal, _ := os.ReadFile(app.libraryPath)
	if !strings.Contains(string(savedLocal), "☑ 001.") {
		t.Fatal("local save lost")
	}
}
