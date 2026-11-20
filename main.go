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
	return filepath.Join(config, "AtlasDelCine", "mis_365_peliculas.txt")
}

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
	out := strings.ReplaceAll(string(template), "__ATLAS_LISTA_INICIAL__", html.EscapeString(initial))
	out = strings.ReplaceAll(out, "__ATLAS_TOKEN__", token)
	out = strings.ReplaceAll(out, "__APP_AUTHOR__", html.EscapeString(appAuthor))
	out = strings.ReplaceAll(out, "__APP_COPYRIGHT__", html.EscapeString(appCopyright))
	return out, nil
}

type runningApp struct {
	mu            sync.Mutex
	savePath      string
	libraryPath   string
	library       libraryState
	platforms     *availabilityServer
	cloud         *cloudBackups
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
	case "/api/sync/history", "/api/sync/revisions", "/api/sync/recovery/list", "/api/sync/recovery/restore":
		a.serveSyncSupport(w,r)
	case "/api/cloud/status", "/api/cloud/config", "/api/cloud/backup", "/api/cloud/browse", "/api/cloud/preview", "/api/cloud/sync", "/api/cloud/auto":
		a.serveCloudHTTP(w, r)
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
	syncStore, syncErr := newDesktopSyncStore(filepath.Dir(savePath))
 if syncErr!=nil {
  // Un fallo en el estado de sincronización no debe impedir usar datos locales.
  logDiagnostic("Sincronización no disponible: "+syncErr.Error())
 }
 app := &runningApp{html: []byte(page), savePath: savePath, libraryPath: libPath, library: lib, platforms: newAvailabilityServer(filepath.Dir(savePath)), cloud: newCloudBackups(filepath.Dir(savePath)), syncStore:syncStore, token: token, stopCh: make(chan struct{})}
	app.touch()
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("no se ha podido iniciar la aplicación local: %w", err)
	}
	server := &http.Server{Handler: http.HandlerFunc(app.serveHTTP), ReadHeaderTimeout: 6 * time.Second}
	go func() { _ = server.Serve(listener) }()
	url := "http://" + listener.Addr().String() + "/"
	logDiagnostic("Servidor local activo en " + url)
	if err := launchVisualWindow(url, filepath.Join(filepath.Dir(savePath), "perfil_webview2_v7_3"), app.stop, func() bool { return app.pageHeartbeat.Load() }); err != nil {
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
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			return server.Shutdown(ctx)
		case <-ticker.C:
			if time.Since(time.Unix(app.lastActivity.Load(), 0)) > 12*time.Minute {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				return server.Shutdown(ctx)
			}
		}
	}
}

func main() {
	if maybeNativeChild() {
		return
	}
	if err := serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		notifyError("Selector Cine 365", err.Error())
	}
}
