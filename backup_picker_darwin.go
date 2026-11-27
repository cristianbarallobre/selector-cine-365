//go:build darwin

package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// Un diálogo nativo de macOS para seleccionar iCloud Drive, Dropbox, etc.
// Script constante: ningún dato del usuario se evalúa como AppleScript.
func chooseBackupFolder() (string, error) {
	const script = `POSIX path of (choose folder with prompt "Elige la carpeta para guardar tus copias de Selector de Cine 365")`
	cmd := exec.Command("/usr/bin/osascript", "-e", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("No se ha seleccionado una carpeta: %s", strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}
