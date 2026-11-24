//go:build windows && lab

package main

import (
 "strings"
 "testing"
)

func TestLabOpenListNumberWithoutSearch(t *testing.T){
 page,err:=makeHTML("☐ Película de prueba (2026)","lab-token")
 if err!=nil{t.Fatal(err)}
 for _,needle:=range []string{
  "function openFilmInList(m){labJumpToFilmList(m);}",
  "function labJumpToFilmList(m){",
  "b.classList.add('lab-list-jump-target')",
  "filmPage=Math.floor(index/capacity)",
  "search='';$('search').value='';filter='all'",
  "labFiltersResetForList();",
  "if(m.id!==libraryState.active_id||virtualMode||labFavoritesMode)switchCollection(m.id);",
  "const labWasFavorite=labFavoritesMode;labFavoritesMode=false;",
  "(!virtualMode&&id===libraryState.active_id&&!labWasFavorite))return;",

 }{
  if !strings.Contains(page,needle){t.Errorf("falta navegación por número LAB: %s",needle)}
 }
 if strings.Contains(page,"selectSearchScope('current');search=m.title;") {
  t.Error("El número de lista sigue activando una búsqueda por título")
 }
 if strings.Count(page,"function openFilmInList(m){")!=1{
  t.Error("Solo debe existir una implementación del salto por número")
 }
}
