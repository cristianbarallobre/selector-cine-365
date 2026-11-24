//go:build windows && lab

package main

import (
 "strings"
 "testing"
)

func TestLabLibraryAdvancedFiltersInjected(t *testing.T){
 page,err:=makeHTML("☐ Película (1970)","token")
 if err!=nil{t.Fatal(err)}
 checks:=[]string{
  "id=\"labLibraryToolbar\"",
  "id=\"labFiltersOverlay\"",
  "id=\"labFilterOpen\"",
  "id=\"labFilterDecade\"",
  "id=\"labFilterGenre\"",
  "id=\"labFilterDuration\"",
  "id=\"labFilterOrder\"",
  "id=\"labFilterResetAll\"",
  "function labFiltersMatches(f)",
  "function labFiltersSort(matches)",
  "function labStatsHeroValues(seen,known,minutes)",
  "labFiltersUpdateSummary(matches.length)",
  "const source=labFavoritesMode&&!global?labFavoritesFilms():global?filmMembershipIndex().rows:parsed.entries;",
  "if(selectedListPlatforms.size)",
  "return labFiltersMatches(f);",
  "matches=labFiltersSort(matches);",
  "labFiltersCanReorder();",
  "for(const [label,val] of labStatsHeroValues(seen,known,minutes)){",
  "data-view=\"instructions\" title=\"Instrucciones\" aria-label=\"Instrucciones\">?</button>",
 }
 for _,s:=range checks {if !strings.Contains(page,s){t.Errorf("no se ha inyectado el filtro LAB: %s",s)}}
 for _,s:=range []string{"['Pendientes',parsed.pending.length]","['Fichas TMDB',known.filter","Escribe para buscar en todas las listas"} {
  if strings.Contains(page,s){t.Errorf("debe desaparecer elemento antiguo: %s",s)}
 }
 if strings.Count(page,"id=\"labFiltersOverlay\"")!=1{t.Fatal("debe existir un único panel de filtros")}
 if strings.Count(page,"id=\"labFilterOpen\"")!=1{t.Fatal("debe existir un único botón Filtros")}
 if !strings.Contains(page,"id=\"labSmartEnabled\""){t.Fatal("el sorteo inteligente debe conservarse")}
 if !strings.Contains(page,"id=\"labDiscoverView\""){t.Fatal("Descubrir no debe desaparecer")}
}
