//go:build windows || darwin

package main

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDiscoveryImageReuseAfterReopening(t *testing.T) {
	var picture bytes.Buffer
	if err := png.Encode(&picture, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	t.Setenv("LOCALAPPDATA", root)
	t.Setenv("HOME", root)
	oldClient := tmdbClient
	t.Cleanup(func() { tmdbClient = oldClient })
	calls := 0
	tmdbClient = &http.Client{Transport: roundTripMock(func(r *http.Request) (*http.Response, error) {
		calls++
		var body io.Reader = bytes.NewReader(picture.Bytes())
		if r.URL.Path == "/t/p/w342/invalid.png" {
			body = bytes.NewBufferString("<html>Respuesta incorrecta</html>")
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(body)}, nil
	})}
	app := func() *runningApp {
		return &runningApp{libraryPath: filepath.Join(root, "colecciones.json"), token: "test"}
	}
	get := func(a *runningApp, path string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/api/discover/poster?path="+path, nil)
		req.Header.Set("X-Atlas-Token", "test")
		w := httptest.NewRecorder()
		a.serveLabDiscovery(w, req)
		return w
	}
	for i := 0; i < 2; i++ {
		w := get(app(), "/portrait.png")
		if w.Code != 200 || w.Header().Get("Content-Type") != "image/png" || !bytes.Equal(w.Body.Bytes(), picture.Bytes()) {
			t.Fatalf("Image request %d: status %d, type %q", i, w.Code, w.Header().Get("Content-Type"))
		}
	}
	if calls != 1 {
		t.Fatalf("Expected one TMDB download across two app instances, got %d", calls)
	}
	if w := get(app(), "/invalid.png"); w.Code != 502 {
		t.Fatalf("Unexpected status on invalid image: %d", w.Code)
	}
	dir := discoveryImageDir(app().libraryPath)
	if count, size := posterDiskUsage(dir); count != 1 || size != int64(picture.Len()) {
		t.Fatalf("Unexpected cache usage: %d files, %d bytes", count, size)
	}
	if n, err := app().discoveryImages.maintain(dir, true); n != 1 || err != nil {
		t.Fatalf("Could not clear the downloaded portrait: %d, %v", n, err)
	}
}

func TestDiscoveryImageCacheKeepsRecentlyUsedImages(t *testing.T) {
	dir := t.TempDir()
	var picture bytes.Buffer
	if err := png.Encode(&picture, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"recent.png", "expired.png", "older.png"} {
		if err := os.WriteFile(filepath.Join(dir, name), picture.Bytes(), 0600); err != nil {
			t.Fatal(err)
		}
	}
	now := time.Now()
	for name, age := range map[string]time.Duration{"recent.png": 36 * time.Hour, "expired.png": discoveryImageMaxAge + time.Hour, "older.png": 24 * time.Hour} {
		date := now.Add(-age)
		if err := os.Chtimes(filepath.Join(dir, name), date, date); err != nil {
			t.Fatal(err)
		}
	}
	var cache discoveryImageCache
	if !bytes.Equal(cache.get(dir, "recent.png"), picture.Bytes()) || cache.get(dir, "expired.png") != nil {
		t.Fatal("Did not reuse the recent image or kept the expired one")
	}
	// Introduce a budget-sized older file to force eviction before the next insertion.
	large := filepath.Join(dir, "older.png")
	if err := os.Truncate(large, discoveryImageMaxBytes-int64(picture.Len())); err != nil {
		t.Fatal(err)
	}
	older := now.Add(-24 * time.Hour)
	if err := os.Chtimes(large, older, older); err != nil {
		t.Fatal(err)
	}
	if _, err := cache.maintain(dir, false); err != nil {
		t.Fatal(err)
	}
	cache.put(dir, "newest.png", picture.Bytes())
	if count, size := posterDiskUsage(dir); count != 2 || size != 2*int64(picture.Len()) {
		t.Fatalf("Unexpected cache size after eviction: %d files, %d bytes", count, size)
	}
	if cache.get(dir, "recent.png") == nil || cache.get(dir, "newest.png") == nil {
		t.Fatal("Eviction removed a recently used image")
	}
}
