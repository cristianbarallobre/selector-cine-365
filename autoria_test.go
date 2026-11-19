package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAutorVisibleYMetadatos(t *testing.T) {
	page, err := makeHTML("☐ 001. Prueba de identidad (2026)", "sample123")
	if err != nil {
		t.Fatal(err)
	}
	for _, must := range []string{appAuthor, "Hecho por", appCopyright, "<meta name=\"author\"", "Selector Cine 365"} {
		if !strings.Contains(page, must) {
			t.Errorf("marca de autor ausente: %s", must)
		}
	}
	if strings.Contains(page, "__APP_AUTHOR__") || strings.Contains(page, "__APP_COPYRIGHT__") {
		t.Fatal("falta sustituir los marcadores de autoría")
	}
	app := &runningApp{html: []byte(page)}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	result := httptest.NewRecorder()
	app.serveHTTP(result, req)
	if result.Header().Get("X-Selector-Creator") != appAuthor {
		t.Fatal("falta identificador del autor en cabeceras HTTP")
	}
	if result.Header().Get("X-Selector-Source") != appSourceURL {
		t.Fatal("falta repo canónico en cabeceras HTTP")
	}
}
