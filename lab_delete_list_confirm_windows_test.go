//go:build windows && lab

package main

import (
 "strings"
 "testing"
)

// Asegura que el ejecutable LAB incorpora el diálogo sin modificar la estable.
func TestLabDeleteListUsesCustomConfirmation(t *testing.T){
 page,err:=makeHTML("☐ Película de prueba (2026)","token-prueba")
 if err!=nil{t.Fatal(err)}
 for _,needle:=range []string{
  `id="labDeleteListOverlay" hidden`,
  `role="alertdialog" aria-modal="true"`,
  `id="labDeleteListCancel"`,
  `id="labDeleteListConfirm"`,
  "async function deleteList(id){",
  "if(!await labConfirmDeleteList(l))return;",
  "libraryState.deleted_lists[id]=true;",
 }{
  if !strings.Contains(page,needle){t.Errorf("Falta diálogo LAB: %s",needle)}
 }
 if strings.Contains(page,"if(!confirm('¿Eliminar «'+l.name+'»?"){t.Error("Sigue usando confirm() de navegador al borrar lista")}
 if strings.Count(page,`id="labDeleteListOverlay"`)!=1{t.Error("La confirmación debe existir una sola vez")}
}
