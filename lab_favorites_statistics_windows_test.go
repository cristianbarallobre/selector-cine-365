//go:build windows && lab

package main

import (
 "strings"
 "testing"
)

func TestLabFavoritesStatisticsUseVirtualList(t *testing.T) {
 page,err:=makeHTML("☑ Vista física (2024)\n☐ Pendiente física (2025)","lab-favorites-statistics")
 if err!=nil{t.Fatal(err)}
 for _,needle:=range []string{
  "function labStatsDataset(){",
  "function labRefreshStatsPanels(){",
  "const labStats=labStatsDataset();const n=labStats.entries.length;",
  "$('watched').textContent=labStats.seen;",
  "$('remaining').textContent=labStats.pending.length;",
  "const pct=n?Math.round(100*labStats.seen/n):0;",
  "if(!selectedSortPlatforms.size)return labStatsDataset().pending;",
  "return labStatsDataset().pending.filter(f=>{",
  "const count=eligibleFilms().length;const total=labStatsDataset().pending.length;",
  "const seen=labStatsDataset().entries.filter(f=>f.checked)",
  "$('statsTitle').textContent=labFavoritesMode?'♥ Favoritos':",
  "$('markWatched').disabled=virtualMode||labFavoritesMode;",
  "if(chosen&&!labFavoritesMode)toggle(chosen.idx);",
  "function toggle(idx){\n if(virtualMode||labFavoritesMode)return;",
  "const all=searchScope==='all'?filmMembershipIndex().rows:labStatsDataset().entries;",
  "const films=labStatsDataset().entries.filter(f=>status==='both'||(status==='seen'?f.checked:!f.checked));",
  "for(const film of labStatsDataset().entries){",

  "if(labFavoritesMode&&!labFavoriteCount()){labFavoritesMode=false;labRefreshStatsPanels();}",
 }{
  if !strings.Contains(page,needle){t.Errorf("Estadística de Favoritos no integrada: %q",needle)}
 }
 if strings.Contains(page,"$('markWatched').disabled=virtualMode;"){
  t.Error("Un sorteo en Favoritos no debe permitir marcar por índice de la lista anterior")
 }
 if strings.Count(page,"function labStatsDataset(){")!=1{
  t.Error("El origen de datos estadísticos no puede estar duplicado")
 }
}
