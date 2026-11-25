package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCollectionsRoundTripAndLegacyCompatibility(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "colecciones.json")
	legacy := filepath.Join(dir, "mis_365_peliculas.txt")
	l := initialLibrary("☑ 001. Ciudadano Kane (1941)\n☐ 002. Vértigo (1958)")
	l.Lists = append(l.Lists, filmCollection{ID: "terror", Name: "Terror", Text: "☐ Nosferatu (1922)"})
	l.ActiveID = "terror"
	if err := saveLibrary(p, legacy, l); err != nil {
		t.Fatal(err)
	}
	got := loadLibrary(p, "other")
	if got.ActiveID != "terror" || len(got.Lists) != 2 {
		t.Fatalf("unexpected: %+v", got)
	}
	txt, err := os.ReadFile(legacy)
	if err != nil || string(txt) != l.Lists[0].Text {
		t.Fatal("legacy atlas list not maintained")
	}
	got.Lists[1].Text = "☑ Nosferatu (1922)"
	if err := saveLibrary(p, legacy, got); err != nil {
		t.Fatal(err)
	}
	if len(loadLibrary(p, "").Lists) != 2 {
		t.Fatal("lost collections")
	}
}
func TestInvalidCollectionsRejected(t *testing.T) {
	l := initialLibrary("abc")
	l.Lists = append(l.Lists, filmCollection{ID: "atlas", Name: "Duplicate", Text: "foo"})
	if validateLibrary(l) == nil {
		t.Fatal("duplicate ID accepted")
	}
}

func TestLibraryAPIAuthenticatesAndWrites(t *testing.T) {
	dir := t.TempDir()
	state := initialLibrary("☑ Una (2020)\n☐ Dos (2021)")
	a := &runningApp{libraryPath: filepath.Join(dir, "colecciones.json"), savePath: filepath.Join(dir, "mis_365_peliculas.txt"), library: state, token: "secret", stopCh: make(chan struct{})}
	r := httptest.NewRequest(http.MethodGet, "/api/library", nil)
	w := httptest.NewRecorder()
	a.serveHTTP(w, r)
	if w.Code != 403 {
		t.Fatalf("expected 403, got %d", w.Code)
	}
	r = httptest.NewRequest(http.MethodGet, "/api/library", nil)
	r.Header.Set("X-Atlas-Token", "secret")
	w = httptest.NewRecorder()
	a.serveHTTP(w, r)
	if w.Code != 200 || !strings.Contains(w.Body.String(), "Una") {
		t.Fatalf("bad library response: %d %s", w.Code, w.Body.String())
	}
	etag := w.Header().Get("ETag")
	state.Lists = append(state.Lists, filmCollection{ID: "lista2", Name: "Otra", Text: "☐ Tres (2022)"})
	payload, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	r = httptest.NewRequest(http.MethodPost, "/api/library", bytes.NewReader(payload))
	r.Header.Set("X-Atlas-Token", "secret")
	r.Header.Set("If-Match", etag)
	w = httptest.NewRecorder()
	a.serveHTTP(w, r)
	if w.Code != 204 {
		t.Fatalf("expected 204, got %d (%s)", w.Code, w.Body.String())
	}
	got := loadLibrary(a.libraryPath, "")
	if len(got.Lists) != 2 {
		t.Fatal("library not saved")
	}
}

func TestComparisonIDsPersistAndValidate(t *testing.T) {
	l := initialLibrary("☑ Ciudadano Kane (1941)")
	l.Lists = append(l.Lists, filmCollection{ID: "extra", Name: "Favoritas", Text: "☐ Ciudadano Kane (1941)"})
	l.ComparisonIDs = []string{"atlas", "extra"}
	if err := validateLibrary(l); err != nil {
		t.Fatalf("valid comparison rejected: %v", err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "colecciones.json")
	if err := saveLibrary(path, filepath.Join(dir, "legacy.txt"), l); err != nil {
		t.Fatal(err)
	}
	reloaded := loadLibrary(path, "")
	if len(reloaded.ComparisonIDs) != 2 || reloaded.ComparisonIDs[1] != "extra" {
		t.Fatalf("comparison lost: %+v", reloaded.ComparisonIDs)
	}
	l.ComparisonIDs = []string{"atlas", "atlas"}
	if validateLibrary(l) == nil {
		t.Fatal("duplicate comparison list accepted")
	}
	l.ComparisonIDs = []string{"atlas", "missing"}
	if validateLibrary(l) == nil {
		t.Fatal("unknown comparison list accepted")
	}
	l.ComparisonIDs = []string{"atlas"}
	if validateLibrary(l) == nil {
		t.Fatal("one-list comparison accepted")
	}
}
