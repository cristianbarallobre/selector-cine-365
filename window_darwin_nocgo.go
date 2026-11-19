//go:build darwin && !cgo

package main

import (
	"errors"
	"fmt"
	"os"
)

func launchVisualWindow(url, _ string, onClose func(), _ func() bool) error {
	if os.Getenv("ATLAS_TEST_NO_BROWSER") == "1" {
		fmt.Println("ATLAS_URL", url)
		return nil
	}
	return errors.New("la app de macOS requiere compilarse con CGO_ENABLED=1 y Xcode Command Line Tools")
}
func maybeNativeChild() bool { return false }
func notifyError(title, body string) { fmt.Fprintln(os.Stderr, title+":", body) }
func traceHTTP(path, userAgent string) {}
func logDiagnostic(message string) {}
