//go:build windows

package main

// Ventana nativa Windows + Microsoft Edge WebView2. El proceso que crea la
// ventana es este mismo EXE: la barra de tareas no ve Chrome ni msedge.exe.
// Solo se utilizan Win32/COM y la biblioteca estándar de Go (sin CGO).

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

const (
	wmDestroy               = 0x0002
	wmClose                 = 0x0010
	wmSize                  = 0x0005
	wmMove                  = 0x0003
	wmSetFocus              = 0x0007
	wmAppNavigate           = 0x8021
	wmAppCloseResult        = 0x8022
	wsOverlappedWindow      = 0x00CF0000
	wsClipChildren          = 0x02000000
	wsClipSiblings          = 0x04000000
	swShow                  = 5
	swShowMaximized         = 3 // Win32 SW_SHOWMAXIMIZED, conserva la barra de título y restaurar
	coinitApartmentThreaded = 2
	sOK                     = 0
	eNoInterface            = uintptr(0x80004002)
	eFail                   = uintptr(0x80004005)
)

type winRect struct{ Left, Top, Right, Bottom int32 }
type winPoint struct{ X, Y int32 }
type winMsg struct {
	Hwnd     uintptr
	Message  uint32
	_        uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       winPoint
	LPrivate uint32
}
type winClass struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   uintptr
	Icon       uintptr
	Cursor     uintptr
	Background uintptr
	MenuName   *uint16
	ClassName  *uint16
	IconSmall  uintptr
}

type winGuid struct {
	Data1        uint32
	Data2, Data3 uint16
	Data4        [8]byte
}

var (
	iidIUnknown           = winGuid{0, 0, 0, [8]byte{0xC0, 0, 0, 0, 0, 0, 0, 0x46}}
	iidEnvironmentHandler = winGuid{0x4e8a3389, 0xc9d8, 0x4bd2, [8]byte{0xb6, 0xb5, 0x12, 0x4f, 0xee, 0x6c, 0xc1, 0x4d}}
	iidControllerHandler  = winGuid{0x6c4819f3, 0xc9b7, 0x4260, [8]byte{0x81, 0x27, 0xc9, 0xf5, 0xbd, 0xe7, 0xf6, 0x8c}}
	iidNavStartHandler    = winGuid{0x9adbe429, 0xf36d, 0x432b, [8]byte{0x9d, 0xdc, 0xf8, 0x88, 0x1f, 0xbd, 0x76, 0xe3}}
	iidNavDoneHandler     = winGuid{0xd33a35bf, 0x1c49, 0x4f98, [8]byte{0x93, 0xab, 0x00, 0x6e, 0x05, 0x33, 0xfe, 0x1c}}
	iidProcessFailHandler = winGuid{0x79e0aea4, 0x990b, 0x42d9, [8]byte{0xaa, 0x1d, 0x0f, 0xcc, 0x2e, 0x5b, 0xc7, 0xf1}}
	iidWebMessageHandler  = winGuid{0x57213f19, 0x00e6, 0x49fa, [8]byte{0x8e, 0x07, 0x89, 0x8e, 0xa0, 0x1e, 0xcb, 0xd2}}
)

type callbackVTable struct{ QueryInterface, AddRef, Release, Invoke uintptr }
type callbackObject struct {
	VTable   *callbackVTable
	RefCount int32
	Kind     uint32
}

type nativeWindow struct {
	hwnd                uintptr
	environment         unsafe.Pointer
	controller          unsafe.Pointer
	webview             unsafe.Pointer
	url                 string
	profile             string
	envCallback         callbackObject
	ctrlCallback        callbackObject
	navStartCallback    callbackObject
	navDoneCallback     callbackObject
	processFailCallback callbackObject
	webMessageCallback  callbackObject
	navStartVTable      callbackVTable
	navDoneVTable       callbackVTable
	processFailVTable   callbackVTable
	webMessageVTable    callbackVTable
	navStartToken       uint64
	navDoneToken        uint64
	processFailToken    uint64
	webMessageToken     uint64
	ready               bool
	closing             bool
	closeAttempt        uintptr
	closeDetail         string
	envVTable           callbackVTable
	ctrlVTable          callbackVTable
	loader              *syscall.DLL
	failure             error
}

var activeNativeWindow *nativeWindow // una instancia por proceso

var (
	user32DLL            = syscall.NewLazyDLL("user32.dll")
	kernel32DLL          = syscall.NewLazyDLL("kernel32.dll")
	ole32DLL             = syscall.NewLazyDLL("ole32.dll")
	shell32DLL           = syscall.NewLazyDLL("shell32.dll")
	gdi32DLL             = syscall.NewLazyDLL("gdi32.dll")
	procRegisterClass    = user32DLL.NewProc("RegisterClassExW")
	procCreateWindow     = user32DLL.NewProc("CreateWindowExW")
	procDefWindow        = user32DLL.NewProc("DefWindowProcW")
	procDestroyWindow    = user32DLL.NewProc("DestroyWindow")
	procShowWindow       = user32DLL.NewProc("ShowWindow")
	procUpdateWindow     = user32DLL.NewProc("UpdateWindow")
	procGetMessage       = user32DLL.NewProc("GetMessageW")
	procTranslateMessage = user32DLL.NewProc("TranslateMessage")
	procDispatchMessage  = user32DLL.NewProc("DispatchMessageW")
	procPostQuit         = user32DLL.NewProc("PostQuitMessage")
	procPostMessage      = user32DLL.NewProc("PostMessageW")
	procGetClientRect    = user32DLL.NewProc("GetClientRect")
	procLoadCursor       = user32DLL.NewProc("LoadCursorW")
	procGetModuleHandle  = kernel32DLL.NewProc("GetModuleHandleW")
	procCoInitialize     = ole32DLL.NewProc("CoInitializeEx")
	procCoUninitialize   = ole32DLL.NewProc("CoUninitialize")
	procCoTaskMemFree    = ole32DLL.NewProc("CoTaskMemFree")
	procSetAppID         = shell32DLL.NewProc("SetCurrentProcessExplicitAppUserModelID")
	procExtractIcon      = shell32DLL.NewProc("ExtractIconExW")
	procBlackBrush       = gdi32DLL.NewProc("GetStockObject")
)

func utf16Ptr(s string) *uint16    { return syscall.StringToUTF16Ptr(s) }
func ptr(p unsafe.Pointer) uintptr { return uintptr(p) }
func badHR(hr uintptr) bool        { return uint32(hr)&0x80000000 != 0 }
func comMethod(obj unsafe.Pointer, index uintptr) uintptr {
	table := *(*unsafe.Pointer)(obj)
	return *(*uintptr)(unsafe.Add(table, index*unsafe.Sizeof(uintptr(0))))
}
func invokeCOM(obj unsafe.Pointer, index uintptr, args ...uintptr) uintptr {
	argv := make([]uintptr, 0, len(args)+1)
	argv = append(argv, uintptr(obj))
	argv = append(argv, args...)
	result, _, _ := syscall.SyscallN(comMethod(obj, index), argv...)
	runtime.KeepAlive(obj)
	return result
}
func resizeNativeWebview(w *nativeWindow) {
	if w == nil || w.controller == nil || w.hwnd == 0 {
		return
	}
	var rc winRect
	procGetClientRect.Call(w.hwnd, ptr(unsafe.Pointer(&rc)))
	invokeCOM(w.controller, 6, ptr(unsafe.Pointer(&rc)))
}
func nativeWndProc(hwnd, msg, wparam, lparam uintptr) uintptr {
	w := activeNativeWindow
	switch msg {
	case wmSize:
		resizeNativeWebview(w)
		return 0
	case wmMove:
		if w != nil && w.controller != nil {
			invokeCOM(w.controller, 23)
		}
		return 0
	case wmAppNavigate:
		if w != nil && w.webview != nil {
			logDiagnostic("WM_APP: iniciando Navigate, fuera del callback COM")
			addr := utf16Ptr(w.url)
			hr := invokeCOM(w.webview, 5, ptr(unsafe.Pointer(addr)))
			runtime.KeepAlive(addr)
			logDiagnostic(fmt.Sprintf("Navigate: HRESULT 0x%08X", uint32(hr)))
			if badHR(hr) {
				w.failure = fmt.Errorf("Navigate WebView2 falló: 0x%08X", uint32(hr))
				procDestroyWindow.Call(w.hwnd)
			}
		}
		return 0
	case wmSetFocus:
		if w != nil && w.controller != nil {
			invokeCOM(w.controller, 12, 0)
		}
		return 0
	case wmClose:
		if w != nil && w.ready && w.webview != nil {
			if w.closing {
				return 0
			}
			w.closing = true
			w.closeAttempt++
			w.closeDetail = ""
			attempt := w.closeAttempt
			message := utf16Ptr(fmt.Sprintf("cine-prepare-close:%d", attempt))
			hr := invokeCOM(w.webview, 33, ptr(unsafe.Pointer(message))) // PostWebMessageAsString
			runtime.KeepAlive(message)
			if badHR(hr) {
				procPostMessage.Call(hwnd, wmAppCloseResult, 0, attempt)
			} else {
				time.AfterFunc(20*time.Second, func() { procPostMessage.Call(hwnd, wmAppCloseResult, 0, attempt) })
			}
			return 0
		}
		procDestroyWindow.Call(hwnd)
		return 0
	case wmAppCloseResult:
		if w != nil && w.closing && lparam == w.closeAttempt {
			if wparam == 1 {
				procDestroyWindow.Call(hwnd)
			} else {
				// Invalidate late replies before entering the native modal message loop.
				w.closeAttempt++
				logDiagnostic("Cierre: guardado no confirmado; se ofrece cerrar de todas formas")
				if confirmCloseWithoutSave(hwnd, w.closeDetail) {
					diagnosticEvent("WARN", "cierre", "usuario confirma cierre sin guardado confirmado")
					procDestroyWindow.Call(hwnd)
				} else {
					diagnosticEvent("INFO", "cierre", "usuario cancela cierre")
					w.closing = false
				}
			}
		}
		return 0
	case wmDestroy:
		if w != nil {
			if w.controller != nil {
				invokeCOM(w.controller, 24) // Close
			}
			if w.webview != nil {
				invokeCOM(w.webview, 2) // Release AddRef
				w.webview = nil
			}
			if w.controller != nil {
				invokeCOM(w.controller, 2) // Release AddRef
				w.controller = nil
			}
			if w.environment != nil {
				invokeCOM(w.environment, 2)
				w.environment = nil
			}
			w.hwnd = 0
		}
		procPostQuit.Call(0)
		return 0
	}
	r, _, _ := procDefWindow.Call(hwnd, msg, wparam, lparam)
	return r
}
func callbackQueryInterface(cb *callbackObject, riid *winGuid, out **callbackObject) uintptr {
	if riid == nil || out == nil {
		return eFail
	}
	iid := *riid
	if iid == iidIUnknown || (cb.Kind == 0 && iid == iidEnvironmentHandler) || (cb.Kind == 1 && iid == iidControllerHandler) ||
		(cb.Kind == 2 && iid == iidNavStartHandler) ||
		(cb.Kind == 3 && iid == iidNavDoneHandler) ||
		(cb.Kind == 4 && iid == iidProcessFailHandler) ||
		(cb.Kind == 5 && iid == iidWebMessageHandler) {
		*out = cb
		atomic.AddInt32(&cb.RefCount, 1)
		return sOK
	}
	*out = nil
	return eNoInterface
}
func callbackAddRef(cb *callbackObject) uintptr {
	return uintptr(atomic.AddInt32(&cb.RefCount, 1))
}
func callbackRelease(cb *callbackObject) uintptr {
	return uintptr(atomic.AddInt32(&cb.RefCount, -1))
}
func navigationStarting(_, _, _ unsafe.Pointer) uintptr {
	if w := activeNativeWindow; w != nil {
		w.ready = false
	}
	logDiagnostic("NavigationStarting: WebView2 ha iniciado navegación")
	return sOK
}
func navigationCompleted(_ unsafe.Pointer, sender, args unsafe.Pointer) uintptr {
	if args == nil {
		logDiagnostic("NavigationCompleted: sin argumentos")
		return sOK
	}
	var success int32
	var errorCode int32
	hr := invokeCOM(args, 3, ptr(unsafe.Pointer(&success)))
	if badHR(hr) {
		logDiagnostic(fmt.Sprintf("NavigationCompleted: get_IsSuccess falló 0x%08X", uint32(hr)))
		return sOK
	}
	hr = invokeCOM(args, 4, ptr(unsafe.Pointer(&errorCode)))
	if badHR(hr) {
		logDiagnostic(fmt.Sprintf("NavigationCompleted: get_WebErrorStatus falló 0x%08X", uint32(hr)))
		return sOK
	}
	logDiagnostic(fmt.Sprintf("NavigationCompleted: éxito=%v, WebErrorStatus=%d", success != 0, errorCode))
	if w := activeNativeWindow; w != nil {
		w.ready = success != 0 && webviewString(sender, 4, len(w.url)) == w.url
	}
	return sOK
}
func webviewString(object unsafe.Pointer, method uintptr, maxLength int) string {
	if object == nil {
		return ""
	}
	var value *uint16
	if badHR(invokeCOM(object, method, ptr(unsafe.Pointer(&value)))) || value == nil {
		return ""
	}
	defer procCoTaskMemFree.Call(ptr(unsafe.Pointer(value)))
	for i, character := range unsafe.Slice(value, maxLength+1) {
		if character == 0 {
			return syscall.UTF16ToString(unsafe.Slice(value, i))
		}
	}
	return ""
}
func webMessageReceived(_ unsafe.Pointer, sender, args unsafe.Pointer) uintptr {
	w := activeNativeWindow
	if w == nil || args == nil || sender != w.webview {
		return sOK
	}
	if webviewString(args, 3, len(w.url)) != w.url {
		return sOK
	}
	message := webviewString(args, 5, uiDraftLimit+256)
	if strings.HasPrefix(message,"cine-draft:"){
		ack:=nativeDraftMessage(filepath.Dir(getSavePath()),message)
		if ack!=""{
			value:=utf16Ptr(ack)
			invokeCOM(w.webview,33,ptr(unsafe.Pointer(value)))
			runtime.KeepAlive(value)
		}
		return sOK
	}
	if len(message)>1400{return sOK}
	if strings.HasPrefix(message, "cine-diagnostic:") {
		diagnosticUIEvent(strings.TrimPrefix(message, "cine-diagnostic:"))
		return sOK
	}
	if !w.closing { return sOK }
	ready := fmt.Sprintf("cine-close-ready:%d", w.closeAttempt)
	pending := fmt.Sprintf("cine-close-pending:%d", w.closeAttempt)
	switch {
	case message == ready:
		procPostMessage.Call(w.hwnd, wmAppCloseResult, 1, w.closeAttempt)
	case message == pending || strings.HasPrefix(message, pending+":"):
		w.closeDetail = strings.TrimPrefix(strings.TrimPrefix(message, pending), ":")
		procPostMessage.Call(w.hwnd, wmAppCloseResult, 0, w.closeAttempt)
	}
	return sOK
}
func processFailed(_, _, args unsafe.Pointer) uintptr {
	var kind int32
	if args != nil {
		if hr := invokeCOM(args, 3, ptr(unsafe.Pointer(&kind))); badHR(hr) {
			logDiagnostic(fmt.Sprintf("ProcessFailed: get_ProcessFailedKind 0x%08X", uint32(hr)))
		}
	}
	logDiagnostic(fmt.Sprintf("ProcessFailed: kind=%d", kind))
	return sOK
}
func environmentReady(_ unsafe.Pointer, hresult uintptr, env unsafe.Pointer) uintptr {
	w := activeNativeWindow
	if w == nil || w.hwnd == 0 {
		return sOK
	}
	if badHR(hresult) || env == nil {
		logDiagnostic(fmt.Sprintf("Environment failed: 0x%08X", uint32(hresult)))
		w.failure = fmt.Errorf("WebView2 no pudo arrancar (HRESULT 0x%08X). Comprueba que Microsoft Edge WebView2 Runtime esté instalado", uint32(hresult))
		procDestroyWindow.Call(w.hwnd)
		return sOK
	}
	logDiagnostic("WebView2 Environment creado")
	// COM pointers from asynchronous completion handlers are borrowed.
	// Retain them explicitly so browser initialization survives the callback return.
	invokeCOM(env, 1) // AddRef; balanced on WM_DESTROY
	w.environment = env
	hr := invokeCOM(env, 3, w.hwnd, ptr(unsafe.Pointer(&w.ctrlCallback)))
	if badHR(hr) {
		w.failure = fmt.Errorf("Error al crear el control WebView2: 0x%08X", uint32(hr))
		procDestroyWindow.Call(w.hwnd)
	}
	return sOK
}
func controllerReady(_ unsafe.Pointer, hresult uintptr, controller unsafe.Pointer) uintptr {
	w := activeNativeWindow
	if w == nil || w.hwnd == 0 {
		return sOK
	}
	if badHR(hresult) || controller == nil {
		logDiagnostic(fmt.Sprintf("Controller failed: 0x%08X", uint32(hresult)))
		w.failure = fmt.Errorf("WebView2 no pudo crear la ventana (0x%08X)", uint32(hresult))
		procDestroyWindow.Call(w.hwnd)
		return sOK
	}
	logDiagnostic("WebView2 Controller creado")
	invokeCOM(controller, 1) // AddRef; keep the controller alive after Invoke returns
	w.controller = controller
	resizeNativeWebview(w)
	invokeCOM(controller, 4, 1) // mostrar WebView2
	var webview unsafe.Pointer
	hr := invokeCOM(controller, 25, ptr(unsafe.Pointer(&webview)))
	if badHR(hr) || webview == nil {
		w.failure = fmt.Errorf("Error al iniciar la vista WebView2: 0x%08X", uint32(hr))
		procDestroyWindow.Call(w.hwnd)
		return sOK
	}
	logDiagnostic("WebView2 CoreWebView2 creado; preparando navegación")
	invokeCOM(webview, 1) // AddRef; balanced on WM_DESTROY
	w.webview = webview
	for _, ev := range []struct {
		name   string
		method uintptr
		cb     *callbackObject
		token  *uint64
	}{
		{"NavigationStarting", 7, &w.navStartCallback, &w.navStartToken},
		{"NavigationCompleted", 15, &w.navDoneCallback, &w.navDoneToken},
		{"ProcessFailed", 25, &w.processFailCallback, &w.processFailToken},
		{"WebMessageReceived", 34, &w.webMessageCallback, &w.webMessageToken},
	} {
		hr = invokeCOM(webview, ev.method, ptr(unsafe.Pointer(ev.cb)), ptr(unsafe.Pointer(ev.token)))
		logDiagnostic(fmt.Sprintf("Evento %s: registro 0x%08X", ev.name, uint32(hr)))
	}
	// Post a message: don't invoke Navigate reentrantly from a COM completion callback.
	posted, _, postErr := procPostMessage.Call(w.hwnd, wmAppNavigate, 0, 0)
	logDiagnostic(fmt.Sprintf("Navegación diferida: PostMessage=%d, error=%v", posted, postErr))
	if posted == 0 {
		w.failure = fmt.Errorf("No se pudo programar Navigate: %v", postErr)
		procDestroyWindow.Call(w.hwnd)
	}
	return sOK
}

// El SDK de Microsoft distribuye la DLL cargadora por separado del runtime.
// Instalamos únicamente esa DLL desde NuGet oficial, una sola vez por usuario.
const webview2PackageURL = "https://api.nuget.org/v3-flatcontainer/microsoft.web.webview2/1.0.3537.50/microsoft.web.webview2.1.0.3537.50.nupkg"
const webview2LoaderSHA256 = "2f965e10aed3b356a408978a0e6d74eb86e3e722dd008fa9ad39f68884479e85"

func verifiedWebview2Loader(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	hash := sha256.New()
	size, err := io.Copy(hash, io.LimitReader(file, (5<<20)+1))
	return err == nil && size <= 5<<20 && fmt.Sprintf("%x", hash.Sum(nil)) == webview2LoaderSHA256
}

func obtainWebview2Loader() (string, error) {
	config, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dest := filepath.Join(config, "AtlasDelCine", "webview2", "WebView2Loader.dll")
	if verifiedWebview2Loader(dest) {
		return dest, nil
	}
	// Primero intentamos con un DLL al lado del ejecutable, para instalaciones sin red.
	if exe, err := os.Executable(); err == nil {
		adjacent := filepath.Join(filepath.Dir(exe), "WebView2Loader.dll")
		if verifiedWebview2Loader(adjacent) {
			return adjacent, nil
		}
	}
	client := http.Client{Timeout: 40 * time.Second}
	response, err := client.Get(webview2PackageURL)
	if err != nil {
		return "", fmt.Errorf("no pude obtener el componente oficial WebView2Loader (necesita internet la primera vez): %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", fmt.Errorf("descarga de WebView2Loader: respuesta HTTP %d", response.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, 32<<20))
	if err != nil {
		return "", err
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("paquete NuGet inválido: %w", err)
	}
	var content []byte
	for _, f := range zr.File {
		if strings.EqualFold(f.Name, "runtimes/win-x64/native/WebView2Loader.dll") || strings.HasSuffix(strings.ToLower(f.Name), "/x64/webview2loader.dll") {
			if f.UncompressedSize64 > 5<<20 {
				return "", errors.New("DLL WebView2 demasiado grande")
			}
			fd, err := f.Open()
			if err != nil {
				return "", err
			}
			content, err = io.ReadAll(io.LimitReader(fd, 5<<20))
			fd.Close()
			if err != nil {
				return "", err
			}
			break
		}
	}
	if fmt.Sprintf("%x", sha256.Sum256(content)) != webview2LoaderSHA256 {
		return "", errors.New("no coincide la integridad de WebView2Loader.dll x64")
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0700); err != nil {
		return "", err
	}
	tmp := dest + ".tmp"
	if err := os.WriteFile(tmp, content, 0600); err != nil {
		return "", err
	}
	if err := os.Rename(tmp, dest); err != nil {
		os.Remove(tmp)
		return "", err
	}
	return dest, nil
}
func nativeWindowLoop(url, profile string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := os.MkdirAll(profile, 0755); err != nil {
		return err
	}
	// La interfaz es idéntica a la v6, pero los recursos de ventana los crea este EXE.
	logDiagnostic("Buscando cargador WebView2")
	loaderFile, err := obtainWebview2Loader()
	if err != nil {
		return err
	}
	logDiagnostic("DLL de WebView2 localizada")
	dll, err := syscall.LoadDLL(loaderFile)
	if err != nil {
		return fmt.Errorf("no se pudo cargar WebView2Loader.dll: %w", err)
	}
	defer dll.Release()
	createEnv, err := dll.FindProc("CreateCoreWebView2EnvironmentWithOptions")
	if err != nil {
		return fmt.Errorf("la DLL de Microsoft no ofrece la entrada WebView2: %w", err)
	}
	hr, _, _ := procCoInitialize.Call(0, coinitApartmentThreaded)
	if badHR(hr) {
		return fmt.Errorf("CoInitializeEx falló: 0x%08X", uint32(hr))
	}
	defer procCoUninitialize.Call()
	procSetAppID.Call(ptr(unsafe.Pointer(utf16Ptr("AtlasDelCine.SelectorDeCine365"))))
	w := &nativeWindow{url: url, profile: profile, loader: dll}
	activeNativeWindow = w
	defer func() { activeNativeWindow = nil }()
	qi := syscall.NewCallback(callbackQueryInterface)
	add := syscall.NewCallback(callbackAddRef)
	release := syscall.NewCallback(callbackRelease)
	w.envVTable = callbackVTable{qi, add, release, syscall.NewCallback(environmentReady)}
	w.ctrlVTable = callbackVTable{qi, add, release, syscall.NewCallback(controllerReady)}
	w.navStartVTable = callbackVTable{qi, add, release, syscall.NewCallback(navigationStarting)}
	w.navDoneVTable = callbackVTable{qi, add, release, syscall.NewCallback(navigationCompleted)}
	w.processFailVTable = callbackVTable{qi, add, release, syscall.NewCallback(processFailed)}
	w.webMessageVTable = callbackVTable{qi, add, release, syscall.NewCallback(webMessageReceived)}
	w.navStartCallback = callbackObject{VTable: &w.navStartVTable, RefCount: 1, Kind: 2}
	w.navDoneCallback = callbackObject{VTable: &w.navDoneVTable, RefCount: 1, Kind: 3}
	w.processFailCallback = callbackObject{VTable: &w.processFailVTable, RefCount: 1, Kind: 4}
	w.webMessageCallback = callbackObject{VTable: &w.webMessageVTable, RefCount: 1, Kind: 5}
	w.envCallback = callbackObject{VTable: &w.envVTable, RefCount: 1, Kind: 0}
	w.ctrlCallback = callbackObject{VTable: &w.ctrlVTable, RefCount: 1, Kind: 1}
	instance, _, _ := procGetModuleHandle.Call(0)
	iconFile, _ := os.Executable()
	var icoLarge, icoSmall uintptr
	procExtractIcon.Call(ptr(unsafe.Pointer(utf16Ptr(iconFile))), 0, ptr(unsafe.Pointer(&icoLarge)), ptr(unsafe.Pointer(&icoSmall)), 1)
	classname := utf16Ptr("SelectorCine365WebView2Window")
	cursor, _, _ := procLoadCursor.Call(0, 32512)
	blackBrush, _, _ := procBlackBrush.Call(4)
	wc := winClass{Size: uint32(unsafe.Sizeof(winClass{})), Style: 3, WndProc: syscall.NewCallback(nativeWndProc), Instance: instance, Icon: icoLarge, IconSmall: icoSmall, Cursor: cursor, Background: blackBrush, ClassName: classname}
	cls, _, e := procRegisterClass.Call(ptr(unsafe.Pointer(&wc)))
	if cls == 0 {
		return fmt.Errorf("RegisterClassExW: %w", e)
	}
	title := utf16Ptr("Selector de Cine 365")
	hwnd, _, e := procCreateWindow.Call(0, ptr(unsafe.Pointer(classname)), ptr(unsafe.Pointer(title)), wsOverlappedWindow|wsClipChildren|wsClipSiblings, 80, 70, 1240, 850, 0, 0, instance, 0)
	if hwnd == 0 {
		return fmt.Errorf("CreateWindowExW: %w", e)
	}
	logDiagnostic("Ventana Win32 creada")
	w.hwnd = hwnd
	// Windows LAB abre maximizado, no en modo quiosco: permite restaurar y cerrar.
	showMode := uintptr(swShow)
	showMode = uintptr(swShowMaximized)
	procShowWindow.Call(hwnd, showMode)
	procUpdateWindow.Call(hwnd)
	logDiagnostic("Solicitando crear entorno WebView2")
	hr, _, _ = createEnv.Call(0, ptr(unsafe.Pointer(utf16Ptr(profile))), 0, ptr(unsafe.Pointer(&w.envCallback)))
	if badHR(hr) {
		procDestroyWindow.Call(hwnd)
		return fmt.Errorf("WebView2 no puede iniciarse (HRESULT 0x%08X). Instala Microsoft Edge WebView2 Runtime desde https://developer.microsoft.com/microsoft-edge/webview2/", uint32(hr))
	}
	var msg winMsg
	for {
		ret, _, _ := procGetMessage.Call(ptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 || int32(ret) == -1 {
			break
		}
		procTranslateMessage.Call(ptr(unsafe.Pointer(&msg)))
		procDispatchMessage.Call(ptr(unsafe.Pointer(&msg)))
	}
	runtime.KeepAlive(w)
	logDiagnostic("Bucle de mensajes terminado")
	return w.failure
}

// Cada inicio nativo se ejecuta en un proceso separado. Si el SDK provoca
// una excepción grave, el servidor y los datos permanecen vivos, y podemos
// ofrecer automáticamente el modo compatible, sin dejar al usuario sin app.
func launchVisualWindow(url, profile string, onClose func(), heartbeat func() bool) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	logDiagnostic("Ventana WebView2: arrancando proceso independiente")
	cmd := exec.Command(exe, "--webview2-host", url, profile)
	// Do not pass SW_HIDE via STARTUPINFO to a process whose first ShowWindow
	// must display a native window; -H=windowsgui already hides its console.
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}
	nativeDesktopWindowActive.Store(true)
	if err := cmd.Start(); err != nil {
		nativeDesktopWindowActive.Store(false)
		diagnosticError("nativo", "no se pudo arrancar WebView2", err)
		return launchCompatibleBrowser(url, filepath.Join(filepath.Dir(profile), "perfil_compatible_v7_3"))
	}
	// Impide que WebView2 mantenga una ventana negra indefinidamente.
	// La interfaz llama /api/ping al cargarse, y después cada 12 segundos.
	var finished atomic.Bool
	watchdog := time.AfterFunc(26*time.Second, func() {
		if !finished.Load() && !heartbeat() {
			logDiagnostic("WebView2: no ha llegado /api/ping antes del tiempo límite; cerrando host y activando navegador compatible")
			_ = cmd.Process.Kill()
		}
	})
	go func() {
		err := cmd.Wait()
		nativeDesktopWindowActive.Store(false)
		finished.Store(true)
		watchdog.Stop()
		if err != nil {
			diagnosticError("nativo", "Host WebView2 falló", err)
			fbErr := launchCompatibleBrowser(url, filepath.Join(filepath.Dir(profile), "perfil_compatible_v7_3"))
			if fbErr != nil {
				notifyError("Selector de Cine 365", "WebView2 no pudo iniciarse ni fue posible abrir el navegador: "+fbErr.Error()+"\nArchivo de errores: "+diagnosticPath())
				onClose()
			}
			return
		}
		if !heartbeat() {
			logDiagnostic("El usuario cerró la ventana nativa antes del primer heartbeat")
		}

		logDiagnostic("Host WebView2 cerrado normalmente")
		time.Sleep(300 * time.Millisecond)
		onClose()
	}()
	return nil
}
func launchCompatibleBrowser(url, profile string) error {
	if err := os.MkdirAll(profile, 0755); err != nil {
		return err
	}
	var candidates []string
	if p, err := exec.LookPath("msedge.exe"); err == nil {
		candidates = append(candidates, p)
	}
	for _, base := range []string{os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramFiles"), os.Getenv("LocalAppData")} {
		if base == "" {
			continue
		}
		candidates = append(candidates,
			filepath.Join(base, "Microsoft", "Edge", "Application", "msedge.exe"),
			filepath.Join(base, "Google", "Chrome", "Application", "chrome.exe"))
	}
	if p, err := exec.LookPath("chrome.exe"); err == nil {
		candidates = append(candidates, p)
	}
	for _, path := range candidates {
		if fi, err := os.Stat(path); err != nil || fi.IsDir() {
			continue
		}
		args := []string{"--app=" + url, "--new-window", "--user-data-dir=" + profile, "--no-first-run", "--no-default-browser-check", "--disable-background-mode"}
		args = append(args, "--start-maximized")
		cmd := exec.Command(path, args...)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if err := cmd.Start(); err == nil {
			logDiagnostic("Modo compatible abierto con navegador dedicado")
			go func() { _ = cmd.Wait() }()
			return nil
		}
	}
	// No hay Edge/Chrome en las rutas habituales; intentar navegador predeterminado.
	cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", url)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() { _ = cmd.Wait() }()
	logDiagnostic("Modo compatible abierto con navegador predeterminado")
	return nil
}
// Se llama antes del servidor HTTP en el proceso secundario.
func maybeNativeChild() bool {
	if len(os.Args) < 2 || os.Args[1] != "--webview2-host" {
		return false
	}
	if len(os.Args) != 4 {
		logDiagnostic("Host: argumentos incompletos")
		os.Exit(2)
	}
	logDiagnostic("Host: inicializando WebView2")
	err := nativeWindowLoop(os.Args[2], os.Args[3])
	if err != nil {
		diagnosticError("nativo", "Host: error de WebView2", err)
		os.Exit(3)
	}
	logDiagnostic("Host: cerrado sin errores")
	os.Exit(0)
	return true
}
// The native dialog is independent of the renderer and HTTP server.
func confirmCloseWithoutSave(hwnd uintptr, detail string) bool {
	if detail == "" {
		detail = "La interfaz no respondió; no se puede identificar qué cambios quedaron pendientes."
	}
	proc := user32DLL.NewProc("MessageBoxW")
	body := utf16Ptr("No se pudo confirmar el guardado local.\n\n" + detail + "\n\nSi cierras ahora, puedes perder cambios sin guardar.\n\n¿Cerrar de todas formas?\nSí: cerrar ahora.\nNo: volver para reintentar o exportar una copia de emergencia.")
	title := utf16Ptr("Selector Cine 365 · Guardado no confirmado")
	result, _, _ := proc.Call(hwnd, ptr(unsafe.Pointer(body)), ptr(unsafe.Pointer(title)), 0x4|0x30|0x100)
	runtime.KeepAlive(body)
	runtime.KeepAlive(title)
	return result == 6 // IDYES; default button is No.
}

func notifyError(title, body string) {
	proc := user32DLL.NewProc("MessageBoxW")
	proc.Call(0, ptr(unsafe.Pointer(utf16Ptr(body))), ptr(unsafe.Pointer(utf16Ptr(title))), 0x10)
}

// Record startup and resumed heartbeats rather than every 12-second ping.
var lastWindowsHeartbeat atomic.Int64
func traceHTTP(path, _ string) {
	if path == "/" { diagnosticEvent("INFO", "motor", "interfaz solicitada"); return }
	if path != "/api/ping" { return }
	now := time.Now().Unix()
	previous := lastWindowsHeartbeat.Swap(now)
	if previous == 0 { diagnosticEvent("INFO", "motor", "primer heartbeat") } else if now-previous > 45 {
		diagnosticEvent("WARN", "motor", fmt.Sprintf("heartbeat reanudado tras %d s", now-previous))
	}
}
