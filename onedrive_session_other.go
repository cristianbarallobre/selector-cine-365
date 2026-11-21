//go:build !windows

package main

// Native secure persistence is intentionally Windows-specific for now.
// macOS and Linux must never fall back to plaintext refresh token files.
func windowsProtectedSessionStore(dir string) graphSessionStore {return nil}
