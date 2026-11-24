//go:build windows

package main

import (
 "embed"
 "strings"
)

//go:embed lab_search_perf_windows.js
var labPerformanceScript embed.FS

// Añade optimizaciones solo a Windows LAB. La plantilla estable no se toca.
func injectLabSearchPerformance(page string) string {
 if !appIsLab{return page}
 raw,err:=labPerformanceScript.ReadFile("lab_search_perf_windows.js")
 if err!=nil{return page}
 const oldInput="$('search').addEventListener('input',e=>{search=e.target.value;filmPage=0;renderList();});"
 const oldGroup="for(const f of parse(list.text).entries){const key=filmKey(f);"
 const boot="updateListPlatformOptions();update();renderCollectionControls();updatePlatformProgress();void initLibrary().then("
 required:=[]string{
  "function filmMembershipIndex(){",
  "function buildIdentityGroups(){",
  oldGroup,
  oldInput,
  "function labFiltersRefreshChoices(){",
  "function selectSearchScope(scope){",
  "function restartList(){",
  boot,
 }
 for _,needle:=range required {if !strings.Contains(page,needle){return page}}
 page=strings.Replace(page,oldGroup,
  "for(const f of labIndexedEntries(list)){const key=filmKey(f);",1)
 // El índice alternativo devuelve la misma estructura y propiedades que el
 // original; los llamadores de búsquedas, listas y sincronización no cambian.
 page=strings.Replace(page,"function filmMembershipIndex(){",
  "function filmMembershipIndex(){return labGlobalMembershipIndex();",1)
 page=strings.Replace(page,"function labFiltersRefreshChoices(){",
  "function labFiltersRefreshChoices(){if(labFiltersOverlay.hidden||!labShouldRefreshFilterChoices())return;",1)
 page=strings.Replace(page,oldInput,
  "$('search').addEventListener('input',e=>labScheduleLibrarySearch(e.target.value));",1)
 page=strings.Replace(page,"function restartList(){",
  "function restartList(){labCancelLibrarySearch();",1)
 page=strings.Replace(page,"function selectSearchScope(scope){",
  "function selectSearchScope(scope){labCancelLibrarySearch();",1)
 // Al cambiar los metadatos secundarios también se actualizan los géneros
 // del selector; estos proceden de la caché de fichas TMDB.
 page=strings.Replace(page,"labSmartCache=data;",
  "labSmartCache=data;labMetadataVersion++;",1)
 page=strings.Replace(page,"function labJumpToFilmList(m){",
  "function labJumpToFilmList(m){labCancelLibrarySearch();",1)
 // Mantenerse dentro de la IIFE principal, y ANTES de la carga inicial.
 page=strings.Replace(page,boot,string(raw)+"\n"+boot,1)
 return page
}
