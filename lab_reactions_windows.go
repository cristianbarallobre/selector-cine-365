//go:build windows

package main

import (
 "embed"
 "strings"
)

//go:embed lab_reactions_windows.html
var labReactionFiles embed.FS

// Solo en Windows LAB: nunca modificar plantilla.html ni la versión estable.
func injectLabReactions(page string) string {
 if !appIsLab{return page}
 raw,err:=labReactionFiles.ReadFile("lab_reactions_windows.html")
 if err!=nil{return page}
 fragments:=strings.SplitN(strings.ReplaceAll(string(raw),"\r\n","\n"),"<!-- LAB_REACTIONS_HTML -->",2)
 if len(fragments)!=2{return page}
 sections:=strings.SplitN(fragments[1],"<!-- LAB_REACTIONS_SCRIPT -->",2)
 if len(sections)!=2{return page}
 expected:=[]string{
  "function switchCollection(id){",
  "function renderCollectionControls(){",
  "function renderList(){",
  "function filmOrderEditable(){",
  "function selectSearchScope(scope){",
  "const source=global?filmMembershipIndex().rows:parsed.entries;",
  "const readonly=virtualMode&&!global;",
  "b.append(togg,pill,labCreateDetailButton(f),fix);frag.appendChild(b);",
  "sel.value=virtualMode?virtualID:libraryState.active_id;",
  "if(!global&&!virtualMode){",
  "updateListPlatformOptions();update();renderCollectionControls();updatePlatformProgress();void initLibrary().then(",
 }
 for _,needle:=range expected{if !strings.Contains(page,needle){return page}}
 page=strings.Replace(page,"</head>",fragments[0]+"\n</head>",1)
 page=strings.Replace(page,`<div id="collectionRows" class="collection-rows"></div>`,
  `<div id="labFavoritesCollection" hidden>
    <div><div class="lab-favorites-title">♥ Favoritos</div><div class="lab-favorites-info" id="labFavoritesCount">Lista automática</div></div>
    <button type="button" id="labFavoritesOpen">Ver favoritos</button>
   </div><div id="collectionRows" class="collection-rows"></div>`,1)
 page=strings.Replace(page,"</main>",sections[0]+"\n</main>",1)
 page=strings.Replace(page,"function switchCollection(id){",
  "function switchCollection(id){if(id===labFavoriteID){labOpenFavorites();return;}const labWasFavorite=labFavoritesMode;labFavoritesMode=false;",1)
 page=strings.Replace(page,"(!virtualMode&&id===libraryState.active_id))return;",
  "(!virtualMode&&id===libraryState.active_id&&!labWasFavorite))return;",1)
 page=strings.Replace(page,"function renderCollectionControls(){",
  "function renderCollectionControls(){labUpdateFavoritesCard();",1)
 page=strings.Replace(page,`if(hasComparison()){const o=document.createElement('option');o.value=virtualID;`,
  `if(labFavoriteCount()){const o=document.createElement('option');o.value=labFavoriteID;o.textContent='♥ Favoritos ('+labFavoriteCount()+')';sel.appendChild(o);}if(hasComparison()){const o=document.createElement('option');o.value=virtualID;`,1)
 page=strings.Replace(page,"sel.value=virtualMode?virtualID:libraryState.active_id;",
  "sel.value=labFavoritesMode?labFavoriteID:virtualMode?virtualID:libraryState.active_id;",1)
 // Si un dispositivo remoto retira el último corazón, no dejar la vista
 // virtual abierta sin opción válida en el selector.
 page=strings.Replace(page,"function renderList(){labFiltersRefreshChoices();",
  "function renderList(){labFiltersRefreshChoices();if(labFavoritesMode&&!labFavoriteCount())labFavoritesMode=false;",1)
 page=strings.Replace(page,"function filmOrderEditable(){",
  "function filmOrderEditable(){if(labFavoritesMode)return false;",1)
 page=strings.Replace(page,"function selectSearchScope(scope){",
  "function selectSearchScope(scope){if(scope==='all')labFavoritesMode=false;",1)
 page=strings.Replace(page,"const source=global?filmMembershipIndex().rows:parsed.entries;",
  "const source=labFavoritesMode&&!global?labFavoritesFilms():global?filmMembershipIndex().rows:parsed.entries;",1)
 page=strings.Replace(page,"const readonly=virtualMode&&!global;",
  "const readonly=labFavoritesMode||virtualMode&&!global;",1)
 page=strings.Replace(page,"if(!global&&!virtualMode){",
  "if(!global&&!virtualMode&&!labFavoritesMode){",1)
 page=strings.Replace(page,"b.append(togg,pill,labCreateDetailButton(f),fix);frag.appendChild(b);",
  "b.append(togg,pill,labReactionControls(f),labCreateDetailButton(f),fix);frag.appendChild(b);",1)
 page=strings.Replace(page,"$('libraryCount').textContent=global?matches.length+' resultados':parsed.entries.length+' películas';",
  "$('libraryCount').textContent=labFavoritesMode?matches.length+' favoritos':global?matches.length+' resultados':parsed.entries.length+' películas';",1)
 page=strings.Replace(page,
  "$('listPlatformCount').textContent=(global?'Todas las listas · ':'Lista actual · ')+matches.length+' resultados';",
  "$('listPlatformCount').textContent=(labFavoritesMode?'Favoritos · ':global?'Todas las listas · ':'Lista actual · ')+matches.length+' resultados';",1)
 boot:="updateListPlatformOptions();update();renderCollectionControls();updatePlatformProgress();void initLibrary().then("
 page=strings.Replace(page,boot,sections[1]+"\n"+boot,1)
 return page
}
