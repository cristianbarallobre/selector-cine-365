package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitialMovieCounts(t *testing.T) {
	s, e := getInitialState(filepath.Join(t.TempDir(), "absent.txt"))
	if e != nil {
		t.Fatal(e)
	}
	if n := countFilms(s); n != 365 {
		t.Fatalf("expected 365 movies; got %d", n)
	}
	if n := strings.Count(s, "☑"); n != 0 {
		t.Fatalf("expected 0 seen in public default list; got %d", n)
	}
}
func TestSaveInPlace(t *testing.T) {
	path := filepath.Join(t.TempDir(), "AtlasDelCine", "mis_365_peliculas.txt")
	app := &runningApp{savePath: path, token: "test-token", stopCh: make(chan struct{})}
	initial, _ := getInitialState(path)
	if err := app.save(initial); err != nil {
		t.Fatal(err)
	}
	result, err := getInitialState(path)
	if err != nil {
		t.Fatal(err)
	}
	if result != initial {
		t.Fatal("file did not roundtrip")
	}
	// Reject a cross-site request without the private token.
	r := httptest.NewRequest(http.MethodPost, "/api/save", strings.NewReader("☐ 001. Otra película"))
	w := httptest.NewRecorder()
	app.serveHTTP(w, r)
	if w.Code != 403 {
		t.Fatalf("expected 403, got %d", w.Code)
	}
	r = httptest.NewRequest(http.MethodPost, "/api/save", strings.NewReader("☐ 001. Otra película"))
	r.Header.Set("X-Atlas-Token", "test-token")
	w = httptest.NewRecorder()
	app.serveHTTP(w, r)
	if w.Code != 204 {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	b, _ := os.ReadFile(path)
	if string(b) != "☐ 001. Otra película" {
		t.Fatal("no updated text")
	}
}
func TestPageDisplaysUserData(t *testing.T) {
	page, err := makeHTML("☑ 001. Ciudadano Kane (1941)\n☐ 002. La película <especial> (2020)", "sample123")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(page, "&lt;especial&gt;") {
		t.Fatal("html not escaped")
	}
	if !strings.Contains(page, "sample123") {
		t.Fatal("token not set")
	}
	if !strings.Contains(page, "--gold:#f3c66d") {
		t.Fatal("visual theme was lost")
	}
}
