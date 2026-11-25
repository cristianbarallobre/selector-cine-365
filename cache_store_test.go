package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAvailabilityCacheCoalescesUpdatesAndKeepsTheBudget(t *testing.T) {
	base := t.TempDir()
	cache := newAvailabilityServer(base)
	for i := 0; i < 3250; i++ {
		cache.putCached(fmt.Sprint(i), providerResult{State: "ok", Date: "2000-01-01T00:00:00Z", MetadataVersion: 3})
	}
	cache.putCached("current", providerResult{State: "ok", Date: "2026-01-01T00:00:00Z", MetadataVersion: 3, TMDBID: 42})
	if _, err := os.Stat(cache.cachePath); !os.IsNotExist(err) {
		t.Fatal("the cache was written one film at a time")
	}
	if _, ok := cache.getCached("current", false); !ok {
		t.Fatal("new entry is not immediately visible")
	}
	cache.flushProviderCache()
	loaded := newAvailabilityServer(base)
	if len(loaded.cache) != 3200 {
		t.Fatalf("incorrect cache budget: %d", len(loaded.cache))
	}
	if id, ok := loaded.knownIdentity("current"); !ok || id != 42 {
		t.Fatal("last identification was not persisted")
	}
}

func TestCacheFileMigrationPreservesOriginalAndNewer(t *testing.T) {
	legacy := filepath.Join(t.TempDir(), "Roaming", "plataformas_cache.json")
	dest := filepath.Join(t.TempDir(), "Local", "plataformas_cache.json")
	if err := os.MkdirAll(filepath.Dir(legacy), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(legacy, []byte("cache antigua"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := copyLegacyCacheIfMissing(legacy, dest, 1024); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(dest)
	if err != nil || string(b) != "cache antigua" {
		t.Fatalf("copia no realizada: %q %v", b, err)
	}
	if err := os.WriteFile(dest, []byte("datos nuevos"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := copyLegacyCacheIfMissing(legacy, dest, 1024); err != nil {
		t.Fatal(err)
	}
	b, err = os.ReadFile(dest)
	if err != nil || string(b) != "datos nuevos" {
		t.Fatalf("copia nueva sobrescrita: %q %v", b, err)
	}
	b, err = os.ReadFile(legacy)
	if err != nil || string(b) != "cache antigua" {
		t.Fatal("La fuente Roaming debe seguir intacta")
	}
}
func TestCacheFileMigrationRejectsOversized(t *testing.T) {
	old := filepath.Join(t.TempDir(), "antigua.json")
	dest := filepath.Join(t.TempDir(), "nueva.json")
	if err := os.WriteFile(old, []byte(strings.Repeat("x", 128)), 0600); err != nil {
		t.Fatal(err)
	}
	if err := copyLegacyCacheIfMissing(old, dest, 32); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		t.Fatalf("copió un archivo demasiado grande: %v", err)
	}
}
func TestPosterPolicyLimits(t *testing.T) {
	if appPosterMaxFiles != 2500 {
		t.Fatalf("Máximo de carteles inesperado: %d", appPosterMaxFiles)
	}
	if appPosterMaxBytes != int64(150<<20) {
		t.Fatalf("Presupuesto de carteles inesperado: %d", appPosterMaxBytes)
	}
	if appPosterMaxAge != 120*24*time.Hour {
		t.Fatalf("No cambiar caducidad de carteles: %s", appPosterMaxAge)
	}
}

func TestPosterBudgetDeletesOldestButNotUserFiles(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()
	for i, name := range []string{"abcde.jpg", "abcdf.jpg", "abcdg.jpg", "abcdh.jpg", "abcdj.jpg"} {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(strings.Repeat("x", 30)), 0600); err != nil {
			t.Fatal(err)
		}
		stamp := now.Add(-time.Duration(6-i) * 24 * time.Hour)
		if err := os.Chtimes(p, stamp, stamp); err != nil {
			t.Fatal(err)
		}
	}
	user := filepath.Join(dir, "colecciones.json")
	if err := os.WriteFile(user, []byte("important"), 0600); err != nil {
		t.Fatal(err)
	}
	removed, err := prunePosterCache(dir, 2, 70, 30*24*time.Hour)
	if err != nil || removed != 3 {
		t.Fatalf("limpieza incorrecta: %d %v", removed, err)
	}
	for _, name := range []string{"abcdh.jpg", "abcdj.jpg", "colecciones.json"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatalf("Archivo activo eliminado: %s", name)
		}
	}
	for _, name := range []string{"abcde.jpg", "abcdf.jpg", "abcdg.jpg"} {
		if _, err := os.Stat(filepath.Join(dir, name)); !os.IsNotExist(err) {
			t.Fatalf("Cartel viejo conservado %s", name)
		}
	}
}
func TestPosterExpirationWithoutTouchingOtherFiles(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "abcdef.png")
	if err := os.WriteFile(p, []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}
	stale := time.Now().Add(-200 * 24 * time.Hour)
	if err := os.Chtimes(p, stale, stale); err != nil {
		t.Fatal(err)
	}
	if n, err := prunePosterCache(dir, 100, 100000, 120*24*time.Hour); err != nil || n != 1 {
		t.Fatalf("No caducó el cartel: %d %v", n, err)
	}
}
func TestPlatformCacheLocalWhileCredentialsRemainSeparate(t *testing.T) {
	base := t.TempDir()
	local := t.TempDir()
	if err := os.WriteFile(filepath.Join(base, "plataformas_cache.json"), []byte(`{"entries":{"film|2026":{"state":"ok","title":"Película"}}}`), 0600); err != nil {
		t.Fatal(err)
	}
	p := newAvailabilityServer(base, local)
	if p.cachePath != filepath.Join(local, "plataformas_cache.json") {
		t.Fatalf("Caché no local: %q", p.cachePath)
	}
	if p.posterDir != filepath.Join(local, "carteles") {
		t.Fatal("Carteles no locales")
	}
	if p.settingsPath != filepath.Join(base, "plataformas_ajustes.json") {
		t.Fatal("Configuración fuera de Roaming")
	}
	if p.overridesPath != filepath.Join(base, "identificaciones_tmdb.json") {
		t.Fatal("Identificaciones manuales fuera de Roaming")
	}
	if p.cache["film|2026"].Title != "Película" {
		t.Fatal("La copia de la caché existente no funciona")
	}
	if _, err := os.Stat(filepath.Join(base, "plataformas_cache.json")); err != nil {
		t.Fatal("Se eliminó la caché antigua")
	}
}
