//go:build darwin && cgo

package main

import (
	"path/filepath"
	"testing"
)

func TestMacCacheIsolation(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Dir(getSavePath())
	wantData := filepath.Join(home, "Library", "Application Support", appDataDirectoryName)
	wantCache := filepath.Join(home, "Library", "Caches", appDataDirectoryName)
	if dir != wantData || appCacheRoot(dir) != wantCache {
		t.Fatalf("data %s cache %s", dir, appCacheRoot(dir))
	}
	server := newAvailabilityServer(dir, appCacheRoot(dir))
	if filepath.Dir(server.settingsPath) != dir || filepath.Dir(server.overridesPath) != dir || filepath.Dir(server.cachePath) != wantCache || filepath.Dir(server.posterDir) != wantCache {
		t.Fatal("rutas de caché y datos mezcladas")
	}
}

