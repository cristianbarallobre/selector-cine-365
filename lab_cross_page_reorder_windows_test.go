//go:build windows && lab

package main

import (
 "strings"
 "testing"
)

func TestLabCrossPageReorderInjected(t *testing.T) {
 page,err:=makeHTML("☐ Ejemplo A (2026)\n☐ Ejemplo B (2026)","test-token")
 if err!=nil{t.Fatal(err)}
 for _,needle:=range []string{
  "handle.draggable=false;",
  "function labCrossPageMove(e){",
  "function labCrossPageFinish(e,cancel=false){",
  "document.body.setPointerCapture",
  "document.elementFromPoint(e.clientX,e.clientY)",
  "moveFilmInCurrentList(state.from,target)",
  "const page=labReorderPageTarget(element)",
  "filmPage=page;",
  "lab-crosspage-dragging",
 } {
  if !strings.Contains(page,needle){t.Errorf("Falta funcionalidad de arrastre LAB: %s",needle)}
 }
 if strings.Contains(page,"handle.draggable=true;"){
  t.Error("El HTML5 drag nativo debe estar desactivado en LAB porque pierde el origen al cambiar de página")
 }
 if strings.Count(page,"function labCrossPageMove(e){")!=1{
  t.Error("Solo debe registrarse un controlador de arrastre entre páginas")
 }
}
