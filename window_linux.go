//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func launchVisualWindow(url, profile string, onClose func(), heartbeat func() bool) error {
	if os.Getenv("ATLAS_TEST_NO_BROWSER") == "1" {
		fmt.Println("ATLAS_URL", url)
		return nil
	}
	cmd := exec.Command("chromium", "--headless", "--no-sandbox", "--disable-gpu", "--disable-dev-shm-usage", "--remote-debugging-port=0", "--user-data-dir="+profile, url)
	return cmd.Start()
}
func notifyError(title, body string) { fmt.Fprintln(os.Stderr, title, body) }

func maybeNativeChild() bool { return false }

func traceHTTP(path, userAgent string) {}

func logDiagnostic(message string) {}
