//go:build windows || darwin

package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDiscoverLibraryCacheReusesIdentifiedDetail(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("LOCALAPPDATA", dir)
	app := &runningApp{savePath: filepath.Join(dir, "legacy.txt"), libraryPath: filepath.Join(dir, "colecciones.json"),
		platforms: newAvailabilityServer(dir), token: "test"}
	app.labCatalog().remember(101, labDetailRecord{
		ProvidersAt: time.Now().Unix(),
		Data: map[string]any{"title": "Buena suerte", "release_date": "2025-12-01",
			"poster_path": "/cover.png", "subscriptions": []string{"Filmin"}, "vote_count": 20},
	}, "")
	req := httptest.NewRequest(http.MethodPost, "/api/discover/library-cache", strings.NewReader(`{"id":101}`))
	req.Header.Set("X-Atlas-Token", "test")
	w := httptest.NewRecorder()
	app.serveLabDiscovery(w, req)
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	key := movieCacheKey("Buena suerte", 2025)
	if err := app.platforms.importManualOverrides(map[string]int{key: 101}); err != nil {
		t.Fatal(err)
	}
	reopened := newAvailabilityServer(dir)
	data, ok := reopened.getCached(key, false)
	if !ok || data.TMDBID != 101 || data.PosterPath != "/cover.png" ||
		len(data.Subscriptions) != 1 || data.Subscriptions[0] != "Filmin" {
		t.Fatal("La Biblioteca perdió la ficha de Descubrir al confirmar el mismo ID", data)
	}
	if err := reopened.importManualOverrides(map[string]int{key: 202}); err != nil {
		t.Fatal(err)
	}
	if _, found := reopened.getCached(key, false); found {
		t.Fatal("Se conservó el cartel de otra identificación")
	}
}
