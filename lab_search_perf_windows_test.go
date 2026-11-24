//go:build windows && lab

package main

import (
 "strings"
 "testing"
)

func TestLabSearchPerformanceWiredIntoWindowsHTML(t *testing.T){
 page,err:=makeHTML("☐ Ejemplo (2024)\n☑ Otra (2025)","lab-search-perf-token")
 if err!=nil{t.Fatal(err)}
 for _,needle:=range []string{
  "const labSearchDelayMs=150;",
  "function labGlobalMembershipIndex(){",
  "function labIndexedEntries(list){",
  "function labShouldRefreshFilterChoices(){",
  "function labScheduleLibrarySearch(value){",
  "function filmMembershipIndex(){return labGlobalMembershipIndex();",
  "for(const f of labIndexedEntries(list)){const key=filmKey(f);",
  "function labFiltersRefreshChoices(){if(labFiltersOverlay.hidden||!labShouldRefreshFilterChoices())return;",
  "$('search').addEventListener('input',e=>labScheduleLibrarySearch(e.target.value));",
  "function selectSearchScope(scope){labCancelLibrarySearch();",
  "function restartList(){labCancelLibrarySearch();",
  "labSmartCache=data;labMetadataVersion++;",
 } {
  if !strings.Contains(page,needle){t.Errorf("Falta optimización LAB: %q",needle)}
 }
 if strings.Count(page,"function labGlobalMembershipIndex(){")!=1 {
  t.Error("El índice optimizado solo debe inyectarse una vez")
 }
 if strings.Contains(page,"$('search').addEventListener('input',e=>{search=e.target.value;filmPage=0;renderList();});") {
  t.Error("Sigue dibujando resultados en cada tecla")
 }
 if strings.Count(page,"function labFiltersRefreshChoices(){")!=1{
  t.Error("El selector de filtros no puede duplicarse")
 }
}

func TestLabSearchPerformanceDoesNotTouchStable(t *testing.T){
 original,err:=resources.ReadFile("plantilla.html")
 if err!=nil{t.Fatal(err)}
 if strings.Contains(string(original),"labGlobalMembershipIndex"){
  t.Error("El índice LAB ha modificado la plantilla estable")
 }
 if strings.Contains(string(original),"labScheduleLibrarySearch"){
  t.Error("El debounce de LAB ha modificado la plantilla estable")
 }
}
