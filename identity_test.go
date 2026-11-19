package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestIdentityCacheCarriesAcrossVersionsAndManualOverrideWins(t *testing.T) {
	dir := t.TempDir()
	oldCache := platformCache{Entries: map[string]providerResult{
		movieCacheKey("Vértigo. De entre los muertos", 1958): {State: "ok", TMDBID: 426, MetadataVersion: 3},
		movieCacheKey("Titán", 2000):                         {State: "ok", TMDBID: 10, MetadataVersion: 3},
	}}
	b, _ := json.Marshal(oldCache)
	if err := os.WriteFile(filepath.Join(dir, "plataformas_cache.json"), b, 0600); err != nil {
		t.Fatal(err)
	}
	override := map[string]int{movieCacheKey("Vertigo", 1958): 426}
	b, _ = json.Marshal(override)
	if err := os.WriteFile(filepath.Join(dir, "identificaciones_tmdb.json"), b, 0600); err != nil {
		t.Fatal(err)
	}
	p := newAvailabilityServer(dir)
	got := p.knownIdentities()
	if got[movieCacheKey("Vértigo. De entre los muertos", 1958)] != 426 || got[movieCacheKey("Vertigo", 1958)] != 426 {
		t.Fatalf("IDs previos no migraron: %+v", got)
	}
	if err := p.saveOverride("Titán", 2000, 15); err != nil {
		t.Fatal(err)
	}
	if p.knownIdentities()[movieCacheKey("Titán", 2000)] != 15 {
		t.Fatal("corrección manual debe tener prioridad")
	}
	p2 := newAvailabilityServer(dir)
	if p2.knownIdentities()[movieCacheKey("Titán", 2000)] != 15 {
		t.Fatal("el ID corregido no persiste")
	}
}
func TestIdentityNegativeCacheExpiresAndDoesNotRequireAPIForKnown(t *testing.T) {
	p := newAvailabilityServer(t.TempDir())
	key := movieCacheKey("Desconocida", 2022)
	p.recordIdentity(key, 0)
	if id, ok := p.knownIdentity(key); id != 0 || !ok {
		t.Fatal("no se almacenó el resultado negativo")
	}
	p.mu.Lock()
	entry := p.identities[key]
	entry.Updated = time.Now().Add(-72 * time.Hour).Format(time.RFC3339)
	p.identities[key] = entry
	p.mu.Unlock()
	if _, ok := p.knownIdentity(key); ok {
		t.Fatal("la negativa expirada debe reconsultarse")
	}
	p.recordIdentity(key, 5678)
	if id, ok := p.knownIdentity(key); id != 5678 || !ok {
		t.Fatal("ID positivo no disponible")
	}
}
func TestIdentityHTTPReportsSavedIDs(t *testing.T) {
	p := newAvailabilityServer(t.TempDir())
	p.recordIdentity(movieCacheKey("Vertigo", 1958), 426)
	app := &runningApp{token: "test", platforms: p}
	req := httptest.NewRequest("GET", "/api/identities", nil)
	req.Header.Set("X-Atlas-Token", "test")
	w := httptest.NewRecorder()
	app.serveIdentitiesHTTP(w, req)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"vertigo|1958":426`) {
		t.Fatalf("IDs no devueltos: %d %s", w.Code, w.Body.String())
	}
	req2 := httptest.NewRequest("POST", "/api/identities/resolve", strings.NewReader(`{"title":"Vertigo","year":1958}`))
	req2.Header.Set("X-Atlas-Token", "test")
	w2 := httptest.NewRecorder()
	app.serveIdentitiesHTTP(w2, req2)
	if w2.Code != http.StatusOK || !strings.Contains(w2.Body.String(), `"tmdb_id":426`) {
		t.Fatalf("ID cacheado no se resuelve: %d %s", w2.Code, w2.Body.String())
	}
}
