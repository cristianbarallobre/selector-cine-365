package main

import (
  "encoding/json"
  "net/http"
  "net/http/httptest"
  "testing"
)

// Regresión observada en Beta 13: la interfaz consultaba
// /api/platforms/cache y el router principal respondía 404.
// La clave de TMDB estaba configurada, pero la pantalla quedaba en
// «Conexión local TMDB no disponible. Reintentando…».
func TestDesktopTMDBStatusAndCacheRoutes(t *testing.T) {
  app := &runningApp{
    token: "test-local-token",
    platforms: newAvailabilityServer(t.TempDir()),
  }
  app.platforms.cache["casablanca|1942"] = providerResult{}
  for _, path := range []string{"/api/platforms/status", "/api/platforms/cache"} {
    t.Run(path, func(t *testing.T) {
      req := httptest.NewRequest(http.MethodGet, path, nil)
      req.Header.Set("X-Atlas-Token", "test-local-token")
      recorder := httptest.NewRecorder()
      app.serveHTTP(recorder, req)
      if recorder.Code != http.StatusOK {
        t.Fatalf("%s devuelve HTTP %d en lugar de 200: %s", path, recorder.Code, recorder.Body.String())
      }
      var body map[string]json.RawMessage
      if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
        t.Fatalf("%s devuelve JSON inválido: %v", path, err)
      }
      if path == "/api/platforms/status" && body["configured"] == nil {
        t.Fatal("Falta estado de configuración TMDB")
      }
      if path == "/api/platforms/cache" && body["entries"] == nil {
        t.Fatal("Falta objeto de entradas de caché")
      }
      if path == "/api/platforms/cache" {
        var entries map[string]json.RawMessage
        if err := json.Unmarshal(body["entries"], &entries); err != nil {t.Fatal(err)}
        if _,ok := entries["casablanca|1942"]; !ok { t.Fatal("La caché almacenada no se devuelve") }
      }
    })
  }
}

func TestDesktopTMDBCacheRouteRejectsUnauthorizedRequests(t *testing.T) {
  app := &runningApp{token: "test-local-token", platforms: newAvailabilityServer(t.TempDir())}
  req := httptest.NewRequest(http.MethodGet, "/api/platforms/cache", nil)
  recorder := httptest.NewRecorder()
  app.serveHTTP(recorder, req)
  if recorder.Code != http.StatusForbidden {
    t.Fatalf("La caché debe requerir token local; respuesta %d", recorder.Code)
  }
}
