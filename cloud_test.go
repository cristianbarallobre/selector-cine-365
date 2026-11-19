package main

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCloudBackupsAndHistory(t *testing.T) {
	dir := t.TempDir()
	data := filepath.Join(dir, "AppData")
	sync := filepath.Join(dir, "OneDrive")
	if e := os.MkdirAll(sync, 0700); e != nil {
		t.Fatal(e)
	}
	state := initialLibrary("☑ 001. Ciudadano Kane (1941)\n☐ 002. Casablanca (1942)")
	legacy := filepath.Join(data, "mis_365_peliculas.txt")
	if e := saveLibrary(filepath.Join(data, "colecciones.json"), legacy, state); e != nil {
		t.Fatal(e)
	}
	c := newCloudBackups(data)
	if e := c.configure(sync); e != nil {
		t.Fatal(e)
	}
	result := c.backup()
	if !result.Enabled || result.LastError != "" {
		t.Fatalf("failed initial backup: %+v", result)
	}
	last := filepath.Join(sync, cloudDirectoryName, "actual", "colecciones.json")
	d, e := os.ReadFile(last)
	if e != nil {
		t.Fatal(e)
	}
	var back libraryState
	if e := json.Unmarshal(d, &back); e != nil || len(back.Lists) != 1 || back.Lists[0].Text != state.Lists[0].Text {
		t.Fatal("cloud did not copy collection")
	}
	state.Lists[0].Text = "☑ 001. Ciudadano Kane (1941)\n☑ 002. Casablanca (1942)"
	if e := saveLibrary(filepath.Join(data, "colecciones.json"), legacy, state); e != nil {
		t.Fatal(e)
	}
	result = c.backup()
	if result.LastError != "" {
		t.Fatal(result.LastError)
	}
	d, _ = os.ReadFile(last)
	if !strings.Contains(string(d), "☑ 002.") {
		t.Fatal("cloud copy not updated")
	}
	entries, e := os.ReadDir(filepath.Join(sync, cloudDirectoryName, "historico"))
	if e != nil || len(entries) != 1 {
		t.Fatalf("no daily version: %v %v", entries, e)
	}
	daydata, e := os.ReadFile(filepath.Join(sync, cloudDirectoryName, "historico", entries[0].Name(), "colecciones.json"))
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(string(daydata), "☐ 002.") {
		t.Fatal("daily restore point overwritten")
	}
	again := newCloudBackups(data)
	if !again.status().Enabled {
		t.Fatal("folder setting did not persist")
	}
	if e := again.configure(""); e != nil {
		t.Fatal(e)
	}
	if again.status().Enabled {
		t.Fatal("cloud still enabled")
	}
}
func TestCloudFailureDoesNotTouchLocalData(t *testing.T) {
	dir := t.TempDir()
	local := filepath.Join(dir, "local")
	target := filepath.Join(dir, "cloud")
	_ = os.MkdirAll(target, 0700)
	c := newCloudBackups(local)
	_ = c.configure(target)
	_ = os.RemoveAll(target)
	s := c.backup()
	if s.LastError == "" {
		t.Fatal("expected backup error")
	}
	// Local save remains independent and succeeds even with failed backup.
	state := initialLibrary("☐ 001. Metrópolis (1927)")
	if e := saveLibrary(filepath.Join(local, "colecciones.json"), filepath.Join(local, "mis_365_peliculas.txt"), state); e != nil {
		t.Fatal(e)
	}
	if _, e := os.Stat(filepath.Join(local, "colecciones.json")); e != nil {
		t.Fatal(e)
	}
}

func TestCloudAPIAndLocalSave(t *testing.T) {
	dir := t.TempDir()
	root := filepath.Join(dir, "AppData")
	dest := filepath.Join(dir, "Sync")
	_ = os.MkdirAll(dest, 0700)
	state := initialLibrary("☐ 001. Casablanca (1942)")
	a := &runningApp{libraryPath: filepath.Join(root, "colecciones.json"), savePath: filepath.Join(root, "mis_365_peliculas.txt"), library: state, cloud: newCloudBackups(root), token: "demo"}
	request := func(path, method, body string, auth bool) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		if auth {
			req.Header.Set("X-Atlas-Token", "demo")
		}
		rec := httptest.NewRecorder()
		a.serveHTTP(rec, req)
		return rec
	}
	if got := request("/api/cloud/status", "GET", "", false); got.Code != 403 {
		t.Fatalf("unauth status %v", got.Code)
	}
	b, _ := json.Marshal(map[string]string{"folder": dest})
	if got := request("/api/cloud/config", "POST", string(b), true); got.Code != 200 {
		t.Fatalf("config: %d %s", got.Code, got.Body.String())
	}
	state.Lists[0].Text = "☑ 001. Casablanca (1942)"
	b, _ = json.Marshal(state)
	got := request("/api/library", "POST", string(b), true)
	if got.Code != 204 || got.Header().Get("X-Cine-Cloud") != "ok" {
		t.Fatalf("save response %d %v", got.Code, got.Header())
	}
	copyPath := filepath.Join(dest, cloudDirectoryName, "actual", "colecciones.json")
	bytes, e := os.ReadFile(copyPath)
	if e != nil || !strings.Contains(string(bytes), "☑ 001.") {
		t.Fatalf("cloud save failed %v", e)
	}
	_ = os.RemoveAll(dest)
	state.Lists[0].Text = "☐ 001. Casablanca (1942)"
	b, _ = json.Marshal(state)
	got = request("/api/library", "POST", string(b), true)
	if got.Code != 204 || got.Header().Get("X-Cine-Cloud") != "error" {
		t.Fatalf("local save should succeed even when cloud fails: %d %v", got.Code, got.Header())
	}
	if file, e := os.ReadFile(a.libraryPath); e != nil || !strings.Contains(string(file), "☐ 001.") {
		t.Fatalf("local data gone on cloud failure: %v", e)
	}
}
