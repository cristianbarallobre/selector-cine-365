//go:build windows && lab

package main

import (
 "strings"
 "testing"
)

func TestLabReactionsAreIntegratedOnlyAsVirtualList(t *testing.T){
 page,err:=makeHTML("☑ Una película (2024)\n☐ Otra (2025)","lab-reactions-test")
 if err!=nil{t.Fatal(err)}
 for _,part:=range []string{
  `id="labRatingOverlay"`,
  `id="labFavoritesCollection" hidden`,
  `id="labFavoritesOpen"`,
  `id="labDetailFavorite"`,
  "function labReactionControls(f){",
  "function labFavoriteMembershipButton(f){",
  "function labJumpToFavoriteFilm(f){",
  "const favoriteBadge=labFavoriteMembershipButton(f);if(favoriteBadge)membership.appendChild(favoriteBadge);",
  "button.title='Ver en Favoritos';",
  "labJumpToFavoriteFilm(f);",
  "labListJumpHighlight={listID:labFavoriteID,rawIndex:target.idx};",
  "function labFavoritesFilms(){",
  "function labOpenFavorites(){",
  "function labOpenRating(f){",
  "labDiscoverFavoriteButton(film)",
  "b.append(togg,labFilmRowActions(f,pill,labCreateDetailButton(f),fix));",
  "const source=labFavoritesMode&&!global?labFavoritesFilms():",
  "const readonly=labFavoritesMode||virtualMode&&!global;",
  "function switchCollection(id){if(id===labFavoriteID)",
  "sel.value=labFavoritesMode?labFavoriteID:",
 }{
  if !strings.Contains(page,part){t.Errorf("Falta una parte de reacciones LAB: %s",part)}
 }
 if strings.Count(page,`id="labRatingOverlay"`)!=1{t.Error("El diálogo de valoración está duplicado")}
 if strings.Contains(page,"libraryState.lists.push({id:labFavoriteID"){t.Error("La lista calculada no debe persistirse como una lista duplicada")}
 if strings.Contains(page,"button.title='En Favoritos'"){t.Error("El corazón pequeño debe navegar, no editar")}

}
