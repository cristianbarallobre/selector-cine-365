//go:build darwin && cgo

package main

import (
	"errors"
	"os"
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

func TestMacKeychainSession(t *testing.T) {
	dir := t.TempDir()
	s := protectedSessionStore(dir)
	defer s.clear()
	value := graphStoredSession{Version: graphSessionVersion, ClientID: desktopMicrosoftAppID, RootID: "fixture", Refresh: "synthetic-1"}
	for _, token := range []string{"synthetic-1", "synthetic-2"} {
		value.Refresh = token
		if err := s.save(value); err != nil {
			t.Fatal(err)
		}
		got, err := protectedSessionStore(dir).load()
		if err != nil || got != value {
			t.Fatalf("recuperación de sesión: %v", err)
		}
	}
	if err := s.clear(); err != nil {
		t.Fatal(err)
	}
	if _, err := s.load(); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("desconexión: %v", err)
	}
}
