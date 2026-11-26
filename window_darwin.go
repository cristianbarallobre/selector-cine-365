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
	"runtime"
	"sync"
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

//export cine_log_ui_event
func cine_log_ui_event(event *C.char) {
	if event != nil { diagnosticUIEvent(C.GoString(event)) }
}
