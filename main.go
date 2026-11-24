package main

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"
)

//go:embed plantilla.html lista_inicial.txt icon.png
var resources embed.FS

func countFilms(data string) int {
	n := 0
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || (strings.HasPrefix(strings.ToUpper(line), "ATLAS DEL CINE") || strings.HasPrefix(strings.ToUpper(line), "SELECTOR CINE 365")) ||
			strings.HasPrefix(line, "#") || line == "---" {
			continue
		}
		n++
	}
	return n
}

func getSavePath() string {
	config, err := os.UserConfigDir()
	if err != nil {
		config = os.TempDir()
	}
	return filepath.Join(config, appDataDirectoryName, "mis_365_peliculas.txt")
}

// Las cachés desechables de Windows LAB se guardan en LocalAppData;
// los TXT, las listas, la nube y las identificaciones permanecen en Roaming.
func appCacheRoot(dataDir string) string {
 if runtime.GOOS == "darwin" {
  if root,err:=os.UserCacheDir();err==nil && root!="" {return filepath.Join(root,appDataDirectoryName)}
  return dataDir
 }
 if !appIsLab || runtime.GOOS!="windows"{return dataDir}
 root,err:=os.UserCacheDir()
 if err!=nil||root=="" {return dataDir} // no bloquear la aplicación por un error de entorno
 return filepath.Join(root,appDataDirectoryName,"cache_lab")
}

const desktopFeatures = runtime.GOOS == "windows" || runtime.GOOS == "darwin"
const protectedLibraryEnabled = appIsLab || desktopFeatures

func getInitialState(savePath string) (string, error) {
	defaultList, err := resources.ReadFile("lista_inicial.txt")
	if err != nil {
		return "", err
	}
	if saved, err := os.ReadFile(savePath); err == nil && countFilms(string(saved)) > 0 {
		return strings.TrimPrefix(string(saved), "\ufeff"), nil
	}
	return string(defaultList), nil
}

func makeHTML(initial, token string) (string, error) {
	template, err := resources.ReadFile("plantilla.html")
	if err != nil {
		return "", err
	}
	out := string(template)
	out = strings.ReplaceAll(out, "__ATLAS_TOKEN__", token)
	out = strings.ReplaceAll(out, "__APP_AUTHOR__", html.EscapeString(appAuthor))
	out = strings.ReplaceAll(out, "__APP_COPYRIGHT__", html.EscapeString(appCopyright))
	if appIsLab {
		out = strings.Replace(out, "<title>Selector Cine 365", "<title>Selector Cine 365 LAB", 1)
		out = strings.Replace(out, `<div class="eyebrow">Selector Cine 365</div>`, `<div class="eyebrow">Selector Cine 365 · LAB</div>`, 1)
		out = strings.Replace(out, "🎞️ SELECTOR CINE 365 ·", "🧪 SELECTOR CINE 365 LAB ·", 1)
	}
	out = injectLabDiscovery(out)
	out = injectLabSmartSort(out)
	out = injectLabNetwork(out)
	out = injectLabLibraryFilters(out)
	out = injectLabReactions(out)
	out = injectLabSearchPerformance(out)
 out = injectLabMaintenance(out)
	if desktopFeatures && !appIsLab {
        out = strings.NewReplacer(
            "<title>Selector Cine 365 · Beta 14 Preview</title>", "<title>Selector Cine 365 · 9.2.0</title>",
            "🎞️ SELECTOR CINE 365 · v9 Beta 14 Preview", "🎞️ SELECTOR CINE 365 · 9.2.0",
            ` <span class="lab-nav-tag">LAB</span>`, "",
            `<span class="lab-smart-caption">LAB</span>`, "",
            "SELECTOR CINE 365 · LAB WINDOWS", "SELECTOR CINE 365",
            "SELECTOR CINE 365 · LAB</p>", "SELECTOR CINE 365</p>",
            "Mantenimiento de Windows LAB", "Mantenimiento",
        ).Replace(out)
    }
    if runtime.GOOS == "darwin" {
        out = strings.NewReplacer(
            "WebView2 administra su propia caché y no se elimina mientras la aplicación está abierta.", "WebKit gestiona los datos temporales de esta ventana.",
            "entre Windows y Android", "entre Mac, Windows y Android",
            "En Windows la autorización se guarda cifrada para este usuario", "En el Mac la autorización se guarda en tu llavero",
            "actualizar el EXE", "actualizar la aplicación",
            "desde Windows", "desde el llavero del Mac",
            "guardada cifrada en Windows", "guardada en el llavero del Mac",
            "En Windows se recordará", "En el Mac se recordará",
            "instalada en Windows", "instalada en el Mac",
            "dos versiones de Windows", "dos versiones de la aplicación",
            "este PC", "este Mac",
        ).Replace(out)
    }
    return strings.ReplaceAll(out, "__ATLAS_LISTA_INICIAL__", html.EscapeString(initial)), nil
}

type runningApp struct {
	mu            sync.Mutex
	savePath      string
	libraryPath   string
	library       libraryState
 libraryRecovered bool
	platforms     *availabilityServer
	cloud         *cloudBackups
	graph         *desktopGraph
	graphUploadCh chan struct{}
	syncStore     *desktopSyncStore
	syncPreviewToken string
    syncLastCheck atomic.Int64
    syncLastReceived atomic.Int64
    syncPendingConflict atomic.Bool
	html          []byte
	token         string
	lastActivity  atomic.Int64
	pageHeartbeat atomic.Bool
	stopCh        chan struct{}
	stopOnce      sync.Once
}

func (a *runningApp) touch() { a.lastActivity.Store(time.Now().Unix()) }
func (a *runningApp) stop()  { a.stopOnce.Do(func() { close(a.stopCh) }) }
func (a *runningApp) auth(r *http.Request) bool {
	return r.Header.Get("X-Atlas-Token") == a.token || (r.URL.Query().Get("token") == a.token && a.token != "")
}
func (a *runningApp) serveHTTP(w http.ResponseWriter, r *http.Request) {
	a.touch()
	if r.URL.Path=="/" && r.Method==http.MethodGet && (r.URL.Query().Has("code")||r.URL.Query().Has("error")) {
		a.serveGraphCallback(w,r);return
	}
	if r.URL.Path == "/" || r.URL.Path == "/api/ping" {
		traceHTTP(r.URL.Path, r.UserAgent())
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	// Metadata integrada en el ejecutable para identificar el origen del programa.
	w.Header().Set("X-Selector-Creator", appAuthor)
	w.Header().Set("X-Selector-Source", appSourceURL)
	switch r.URL.Path {
	case "/":
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", 405)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		a.mu.Lock()
		page := append([]byte(nil), a.html...)
		a.mu.Unlock()
		w.Write(page)
	case "/icon.png":
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", 405)
			return
		}
		data, err := resources.ReadFile("icon.png")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Write(data)
	case "/api/library":
		a.serveLibraryHTTP(w, r)
	case "/api/sync/history", "/api/sync/revisions", "/api/sync/recovery/list", "/api/sync/recovery/preview", "/api/sync/recovery/restore":
		a.serveSyncSupport(w,r)
	case "/api/cloud/status", "/api/cloud/config", "/api/cloud/backup", "/api/cloud/browse", "/api/cloud/preview", "/api/cloud/sync", "/api/cloud/auto", "/api/cloud/auto-folder", "/api/cloud/mode":
		a.serveCloudHTTP(w, r)
	case "/api/onedrive-desktop/status", "/api/onedrive-desktop/start", "/api/onedrive-desktop/disconnect":
        a.serveGraphAuthHTTP(w,r)
	case "/api/lab/maintenance":
  a.serveLabMaintenance(w,r)
	case "/api/lab/network-log":
		a.serveLabNetworkLog(w,r)
	case "/api/discover/smart-cache":
		a.serveLabSmartCache(w,r)
	case "/api/discover/list", "/api/discover/search", "/api/discover/detail", "/api/discover/person", "/api/discover/resolve", "/api/discover/poster", "/api/discover/prefetch", "/api/discover/prefetch/pause":
		a.serveLabDiscovery(w, r)
	case "/api/poster":
		a.servePosterHTTP(w, r)
	case "/api/identities", "/api/identities/resolve":
		a.serveIdentitiesHTTP(w, r)
	case "/api/platforms/status", "/api/platforms/cache", "/api/platforms/config", "/api/platforms/lookup", "/api/platforms/match", "/api/platforms/overrides", "/api/platforms/overrides/import":
		a.servePlatformHTTP(w, r)
	case "/api/ping":
		a.pageHeartbeat.Store(true)
		w.WriteHeader(204)
	case "/api/save", "/api/beacon-save":
		if r.Method != http.MethodPost || !a.auth(r) {
			http.Error(w, "No autorizado", 403)
			return
		}
		contents, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1024*1024))
		if err != nil {
			http.Error(w, "Lista demasiado grande", 413)
			return
		}
		data := strings.TrimPrefix(string(contents), "\ufeff")
		if !utf8.ValidString(data) || countFilms(data) == 0 {
			http.Error(w, "No se han encontrado películas", 400)
			return
		}
		if err := a.save(data); err != nil {
			http.Error(w, "Error al guardar: "+err.Error(), 500)
			return
		}
		w.WriteHeader(204)
	case "/api/quit":
		if r.Method != http.MethodPost || !a.auth(r) {
			http.Error(w, "No autorizado", 403)
			return
		}
		w.WriteHeader(204)
		// Allow a pending save request time to finish before stopping the server.
		time.AfterFunc(2*time.Second, a.stop)
	default:
		http.NotFound(w, r)
	}
}
func (a *runningApp) save(data string) error {
	page, err := makeHTML(data, a.token)
	if err != nil {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := os.MkdirAll(filepath.Dir(a.savePath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(a.savePath, []byte(data), 0644); err != nil {
		return err
	}
	a.html = []byte(page)
	if a.cloud != nil {
		_ = a.cloud.backup()
	}
	return nil
}

func randomToken() (string, error) {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf[:]), nil
}

func serve() error {
	token, err := randomToken()
	if err != nil {
		return err
	}
	savePath := getSavePath()
	state, err := getInitialState(savePath)
	if err != nil {
		return err
	}
	page, err := makeHTML(state, token)
	if err != nil {
		return err
	}
	libPath := filepath.Join(filepath.Dir(savePath), "colecciones.json")
	lib := loadLibrary(libPath, state)
 recovered:=false
 if protectedLibraryEnabled {
  var recoveryErr error
  lib,recovered,recoveryErr=labLoadLibraryProtected(libPath,state)
  if recoveryErr!=nil{return recoveryErr}
  if recovered {logDiagnostic("Biblioteca recuperada de copia local válida; original dañado conservado en cuarentena")}
 }
	syncStore, syncErr := newDesktopSyncStore(filepath.Dir(savePath))
 if syncErr!=nil {
  // Un fallo en el estado de sincronización no debe impedir usar datos locales.
  logDiagnostic("Sincronización no disponible: "+syncErr.Error())
 }
 app := &runningApp{libraryRecovered:recovered,html: []byte(page), savePath: savePath, libraryPath: libPath, library: lib, platforms: newAvailabilityServer(filepath.Dir(savePath),appCacheRoot(filepath.Dir(savePath))), cloud: newCloudBackups(filepath.Dir(savePath)), graph:newDesktopGraph(filepath.Dir(savePath)), graphUploadCh:make(chan struct{},1), syncStore:syncStore, token: token, stopCh: make(chan struct{})}
	app.touch()
	go app.graphPublisher()
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("no se ha podido iniciar la aplicación local: %w", err)
	}
	_,port,_:=net.SplitHostPort(listener.Addr().String())
	app.graph.setRedirect("http://localhost:"+port+"/")
	server := &http.Server{Handler: http.HandlerFunc(app.serveHTTP), ReadHeaderTimeout: 6 * time.Second}
	go func() { _ = server.Serve(listener) }()
	url := "http://" + listener.Addr().String() + "/"
	logDiagnostic("Servidor local activo en " + url)
	if err := launchVisualWindow(url, filepath.Join(appCacheRoot(filepath.Dir(savePath)), "perfil_webview2_v7_3"), app.stop, func() bool { return app.pageHeartbeat.Load() }); err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		return fmt.Errorf("no se pudo abrir la ventana visual: %w", err)
	}
	// A forced browser shutdown does not always fire pagehide. Reclaim the helper
	// after a prolonged absence of heartbeat traffic (including when the app is closed).
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-app.stopCh:
			return closeLocalServer(server)
		case <-ticker.C:
			if time.Since(time.Unix(app.lastActivity.Load(), 0)) > 12*time.Minute {
				return closeLocalServer(server)
			}
		}
	}
}

// Closing the UI can race with a slow OneDrive HTTP request. A shutdown
// deadline is not an application error and must not open a native Windows
// 'context deadline exceeded' dialog after the library was already saved.
func closeLocalServer(server *http.Server) error {
 ctx,cancel:=context.WithTimeout(context.Background(),3*time.Second)
 defer cancel()
 err:=server.Shutdown(ctx)
 if errors.Is(err,context.DeadlineExceeded) {
  logDiagnostic("Cierre del servidor local: operación pendiente excedió el tiempo; se cierra la conexión sin mostrar error de aplicación")
  _=server.Close()
  return nil
 }
 return err
}

func main() {
	if maybeNativeChild() {
		return
	}
	if err := serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		notifyError("Selector Cine 365", err.Error())
	}
}
