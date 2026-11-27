//go:build windows

package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

func chooseBackupFolder() (string, error) {
	// The folder path is selected by Windows, never interpolated into a shell command.
	const script = `$ErrorActionPreference='Stop'; Add-Type -AssemblyName System.Windows.Forms; [Console]::OutputEncoding = New-Object System.Text.UTF8Encoding($false); $d=New-Object System.Windows.Forms.FolderBrowserDialog; $d.Description='Elige la carpeta para la copia adicional'; $d.ShowNewFolderButton=$true; if($d.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK){[Console]::WriteLine($d.SelectedPath)}`
	cmd := exec.Command("powershell.exe", "-NoProfile", "-STA", "-Command", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("No se pudo abrir el selector de carpetas: %s", strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(string(out)), nil
}
