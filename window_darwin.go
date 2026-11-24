//go:build darwin && cgo

package main

/*
#cgo LDFLAGS: -framework Cocoa -framework WebKit
#include <stdlib.h>
int cine_macos_run(const char* url);
*/
import "C"

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
	"unsafe"
)

// Cocoa exige que NSApplication viva en el hilo principal del proceso.
func init() { runtime.LockOSThread() }

// WKWebView utiliza el mismo servidor local y la misma interfaz que Windows.
// Esta función bloquea hasta que la persona cierra la ventana de macOS.
func launchVisualWindow(url, _ string, onClose func(), _ func() bool) error {
	if os.Getenv("ATLAS_TEST_NO_BROWSER") == "1" {
		fmt.Println("ATLAS_URL", url)
		return nil
	}
	u := C.CString(url)
	defer C.free(unsafe.Pointer(u))
	logDiagnostic("Iniciando ventana macOS WKWebView")
	result := int(C.cine_macos_run(u))
	logDiagnostic("Ventana macOS finalizada")
	onClose()
	if result != 0 {
		return fmt.Errorf("no se pudo inicializar la ventana de macOS (código %d)", result)
	}
	return nil
}
func maybeNativeChild() bool { return false }
func notifyError(title, body string) { fmt.Fprintln(os.Stderr, title+":", body) }
var firstWebKitHeartbeat sync.Once
func traceHTTP(path, _ string) {
	if path == "/api/ping" { firstWebKitHeartbeat.Do(func(){ logDiagnostic("WKWebView: interfaz iniciada") }) }
}
func logDiagnostic(message string) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return
	}
	path := filepath.Join(dir, "AtlasDelCine", "diagnostico_macos.txt")
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err == nil {
		defer f.Close()
		fmt.Fprintf(f, "%s %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
	}
}
