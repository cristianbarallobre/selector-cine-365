package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func postTestLibrary(t *testing.T, a *runningApp, state libraryState) *httptest.ResponseRecorder {
	t.Helper()
	body, err := json.Marshal(state)
	if err != nil || len(body) > libraryFileLimit {
		t.Fatalf("invalid test request: bytes=%d err=%v", len(body), err)
	}
	r := httptest.NewRequest("POST", "/api/library", bytes.NewReader(body))
	r.Header.Set("X-Atlas-Token", a.token)
	w := httptest.NewRecorder()
	a.serveHTTP(w, r)
	return w
}

func TestLargeLibraryCanBeSavedAndRestored(t *testing.T) {
	dir := t.TempDir()
	a := &runningApp{libraryPath: filepath.Join(dir, "colecciones.json"), savePath: filepath.Join(dir, "legacy.txt"), token: "test"}
	state := initialLibrary(strings.Repeat("☐ Película (2000)\n", 40500))
	for i := 1; i < 5; i++ {
		state.Lists = append(state.Lists, filmCollection{ID: fmt.Sprint(i), Name: "Lista", Text: state.Lists[0].Text})
	}
	if payload, _ := json.Marshal(state); len(payload) <= 3<<20 {
		t.Fatal("fixture must exceed the previous HTTP limit")
	}
	if w := postTestLibrary(t, a, state); w.Code != 204 {
		t.Fatalf("save: %d %s", w.Code, w.Body.String())
	}
	if err := a.saveSyncRecovery(state, ""); err != nil {
		t.Fatal(err)
	}
	copies, err := a.localRecoveryList()
	if err != nil || len(copies) != 1 {
		t.Fatalf("list recovery: %v %v", copies, err)
	}
	if w := postTestLibrary(t, a, initialLibrary("Cambiada")); w.Code != 204 {
		t.Fatalf("replace: %d %s", w.Code, w.Body.String())
	}
	if _, err := a.restoreLocalRecovery(copies[0].ID); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(loadLibrary(a.libraryPath, ""), state) {
		t.Fatal("large library changed on recovery")
	}
}

func TestLibraryReloadsLargeIndentedFileAndBackup(t *testing.T) {
	if !protectedLibraryEnabled {
		t.Skip("LAB recovery")
	}
	dir := t.TempDir()
	a := &runningApp{libraryPath: filepath.Join(dir, "colecciones.json"), savePath: filepath.Join(dir, "legacy.txt"), token: "test"}
	state := initialLibrary("Atlas")
	state.Reactions = make(map[string]filmReaction)
	for i := 0; i < 15000; i++ {
		state.Reactions[fmt.Sprintf("film%d|2020", i)] = filmReaction{Title: strings.Repeat("Z", 75), Year: "2020", Rating: 4, ChangedAt: 1}
	}
	// Seed the exact on-disk format written by LAB-02.15.
	old, err := json.MarshalIndent(state, "", "  ")
	if err != nil || len(old) <= 3<<20 {
		t.Fatalf("fixture: bytes=%d err=%v", len(old), err)
	}
	if err := os.WriteFile(a.libraryPath, old, 0600); err != nil {
		t.Fatal(err)
	}
	loaded, recovered, err := labLoadLibraryProtected(a.libraryPath, "")
	if err != nil || recovered || !reflect.DeepEqual(state, loaded) {
		t.Fatalf("opening old file: recovered=%t err=%v", recovered, err)
	}
	state.Lists[0].Text = "Updated"
	if w := postTestLibrary(t, a, state); w.Code != 204 {
		t.Fatalf("save: %d %s", w.Code, w.Body.String())
	}
	updated, recovered, err := labLoadLibraryProtected(a.libraryPath, "")
	if err != nil || recovered || !reflect.DeepEqual(state, updated) {
		t.Fatalf("reopening after save: recovered=%t err=%v", recovered, err)
	}
	if err := os.WriteFile(a.libraryPath, []byte("{broken"), 0600); err != nil {
		t.Fatal(err)
	}
	restored, recovered, err := labLoadLibraryProtected(a.libraryPath, "")
	if err != nil || !recovered || !reflect.DeepEqual(loaded, restored) {
		t.Fatalf("large backup: recovered=%t err=%v", recovered, err)
	}
}

func TestLibraryCommitWithUnavailableMirror(t *testing.T) {
	dir := t.TempDir()
	a := &runningApp{library: initialLibrary("Before"), libraryPath: filepath.Join(dir, "colecciones.json"), savePath: filepath.Join(dir, "legacy.txt"), token: "test"}
	if err := os.Mkdir(a.savePath, 0700); err != nil {
		t.Fatal(err)
	}
	state := initialLibrary("After")
	w := postTestLibrary(t, a, state)
	if w.Code != 204 || w.Header().Get("X-Cine-Legacy") != "pending" || !reflect.DeepEqual(state, a.library) || !reflect.DeepEqual(state, loadLibrary(a.libraryPath, "")) {
		t.Fatalf("inconsistent commit: HTTP %d headers=%v", w.Code, w.Header())
	}
	if err := os.Remove(a.savePath); err != nil {
		t.Fatal(err)
	}
	w = postTestLibrary(t, a, state)
	data, err := os.ReadFile(a.savePath)
	if w.Code != 204 || w.Header().Get("X-Cine-Legacy") != "" || err != nil || string(data) != "After" {
		t.Fatalf("mirror retry: HTTP %d TXT=%q err=%v", w.Code, data, err)
	}
}

func TestLibraryFailedCommitPreservesMirrorAndMemory(t *testing.T) {
	dir := t.TempDir()
	state := initialLibrary("Before")
	a := &runningApp{library: state, libraryPath: dir, savePath: filepath.Join(dir, "legacy.txt"), token: "test"}
	if err := os.WriteFile(a.savePath, []byte("Before"), 0600); err != nil {
		t.Fatal(err)
	}
	w := postTestLibrary(t, a, initialLibrary("After"))
	data, err := os.ReadFile(a.savePath)
	if w.Code != 500 || !reflect.DeepEqual(state, a.library) || err != nil || string(data) != "Before" {
		t.Fatalf("failed commit: HTTP %d TXT=%q err=%v", w.Code, data, err)
	}
}
