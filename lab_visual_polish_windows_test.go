//go:build windows && lab

package main

import (
 "strings"
 "testing"
)

func TestLabVisualPolishWindows(t *testing.T) {
 page,err:=makeHTML("☐ Prueba (2026)","token")
 if err!=nil{t.Fatal(err)}
 required:=[]string{
  "LAB-02.13 WINDOWS",
  ".header-side .header-tag{order:1}",
  ".header-side .creator-credit{order:2}",
  ".header-side .saved{order:3}",
  `class="page-btn" id="labPrevPage" type="button">‹ Anterior</button>`,
  `class="page-btn" id="labNextPage" type="button">Siguiente ›</button>`,
  `class="page-label" id="labPageLabel"`,
  `class="page-btn" id="labDetailClose" type="button" aria-label="Cerrar ficha" title="Cerrar ficha">✕</button>`,
  "#labDiscoverView select:not(:disabled)",
  ".lab-detail-actions select{",
  "cursor:pointer",
  "window.fetch=labResilientFetch;",
 }
 for _,fragment:=range required {
  if !strings.Contains(page,fragment){t.Errorf("falta cambio Windows LAB: %s",fragment)}
 }
 if strings.Contains(page,"✕ Cerrar"){t.Fatal("el cierre de Detalles debe mostrar solo la cruz")}
 if strings.Count(page,`id="labDetailClose"`)!=1{t.Fatal("debe existir un solo botón de cierre de Detalles")}
 if strings.Index(page,"const labOriginalFetch=")>strings.Index(page,"void initLibrary()"){t.Fatal("instalar fetch antes de arrancar biblioteca")}
}
