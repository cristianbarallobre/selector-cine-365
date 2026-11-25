//go:build windows || darwin

package main

import (
 "embed"
 "strings"
)

//go:embed lab_reactions_windows.html
var labReactionFiles embed.FS

// Solo en Windows LAB: nunca modificar plantilla.html ni la versión estable.
func injectLabReactions(page string) string {
 raw,err:=labReactionFiles.ReadFile("lab_reactions_windows.html")
 if err!=nil{return page}
 fragments:=strings.SplitN(strings.ReplaceAll(string(raw),"\r\n","\n"),"<!-- LAB_REACTIONS_HTML -->",2)
 if len(fragments)!=2{return page}
 sections:=strings.SplitN(fragments[1],"<!-- LAB_REACTIONS_SCRIPT -->",2)
 if len(sections)!=2{return page}
 expected:=[]string{
  "function switchCollection(id,openLibrary=true){",
  "function renderCollectionControls(){",
  "function renderList(){",
  "function filmOrderEditable(){",
  "function selectSearchScope(scope){",
  "const source=global?filmMembershipIndex().rows:parsed.entries;",
  "const readonly=virtualMode&&!global;",
  "b.append(togg,pill,labCreateDetailButton(f),fix);frag.appendChild(b);",
  "membership.append(caption,nums,edit);meta.appendChild(membership);",
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
 page=strings.Replace(page,"function switchCollection(id,openLibrary=true){",
  "function switchCollection(id,openLibrary=true){if(id===labFavoriteID){labOpenFavorites(openLibrary);return;}const labWasFavorite=labFavoritesMode;labFavoritesMode=false;",1)
 page=strings.Replace(page,"(!virtualMode&&id===libraryState.active_id))return;",
  "(!virtualMode&&id===libraryState.active_id&&!labWasFavorite))return;",1)
 page=strings.Replace(page,"function renderCollectionControls(){",
  "function renderCollectionControls(){labUpdateFavoritesCard();",1)
 page=strings.Replace(page,`if(hasComparison()){const o=document.createElement('option');o.value=virtualID;`,
  `if(labFavoriteCount()){const o=document.createElement('option');o.value=labFavoriteID;o.textContent='♥ Favoritos ('+labFavoriteCount()+')';sel.appendChild(o);}if(hasComparison()){const o=document.createElement('option');o.value=virtualID;`,1)
 page=strings.Replace(page,"sel.value=virtualMode?virtualID:libraryState.active_id;",
  "sel.value=labFavoritesMode?labFavoriteID:virtualMode?virtualID:libraryState.active_id;",1)
 // El selector de Favoritos es una vista calculada, pero el panel de
 // «Tu próximo visionado», el sorteo y la pestaña Estadísticas deben usar
 // ese mismo conjunto; nunca modificar parsed ni el TXT real al visitarla.
 page=strings.Replace(page,
  "function update(){membershipCache=null;parsed=parse(rawText);const n=parsed.entries.length;",
  "function update(){membershipCache=null;parsed=parse(rawText);const labStats=labStatsDataset();const n=labStats.entries.length;",1)
 page=strings.Replace(page,
  "$('total').textContent=n;$('watched').textContent=parsed.seen;$('remaining').textContent=parsed.pending.length;",
  "$('total').textContent=n;$('watched').textContent=labStats.seen;$('remaining').textContent=labStats.pending.length;",1)
 page=strings.Replace(page,
  "const pct=n?Math.round(100*parsed.seen/n):0;",
  "const pct=n?Math.round(100*labStats.seen/n):0;",1)
 page=strings.Replace(page,
  "if(!selectedSortPlatforms.size)return parsed.pending;",
  "if(!selectedSortPlatforms.size)return labStatsDataset().pending;",1)
 page=strings.Replace(page,
  "return parsed.pending.filter(f=>{",
  "return labStatsDataset().pending.filter(f=>{",1)
 page=strings.Replace(page,
  "const count=eligibleFilms().length;const total=parsed.pending.length;",
  "const count=eligibleFilms().length;const total=labStatsDataset().pending.length;",1)
 page=strings.Replace(page,
  "if(!parsed)return;const seen=parsed.entries.filter(f=>f.checked),hero=$('statsHero')",
  "if(!parsed)return;const seen=labStatsDataset().entries.filter(f=>f.checked),hero=$('statsHero')",1)
 page=strings.Replace(page,
  "$('statsTitle').textContent=virtualMode?'Películas coincidentes · solo lectura':activeCollection()?.name||'';",
  "$('statsTitle').textContent=labFavoritesMode?'♥ Favoritos':virtualMode?'Películas coincidentes · solo lectura':activeCollection()?.name||'';",1)
 // La fila de Favoritos es de solo lectura en cuanto al estado visto:
 // un sorteo de esa lista NO debe marcar por error la posición equivalente
 // de la última lista física que quedó activa.
 page=strings.Replace(page,"function toggle(idx){\n if(virtualMode)return;",
  "function toggle(idx){\n if(virtualMode||labFavoritesMode)return;",1)
 page=strings.ReplaceAll(page,
  "$('markWatched').disabled=virtualMode;",
  "$('markWatched').disabled=virtualMode||labFavoritesMode;")
 page=strings.Replace(page,
  "$('markWatched').addEventListener('click',()=>{if(chosen)toggle(chosen.idx);});",
  "$('markWatched').addEventListener('click',()=>{if(chosen&&!labFavoritesMode)toggle(chosen.idx);});",1)
 page=strings.Replace(page,"const all=parsed.entries.map(f=>({...f}));","const all=labStatsDataset().entries.map(f=>({...f}));",1)
 page=strings.Replace(page,"year:Number(f.year)||0,force:!!force","year:Number(f.year)||0,tmdb_id:f._favoriteTMDBID||0,force:!!force",1)
 page=strings.Replace(page,"all.filter(f=>!platformData.has(providerKey(f)))","all.filter(f=>{const cached=platformData.get(providerKey(f));return !cached||(f._favoriteTMDBID&&cached.state!=='error'&&cached.tmdb_id!==f._favoriteTMDBID);})",1)
 // El sorteo inteligente también debe usar Favoritos, tanto para
 // seleccionar candidatas como para ponderar géneros ya vistos.
 page=strings.Replace(page,
  "const films=parsed.entries.filter(f=>status==='both'||(status==='seen'?f.checked:!f.checked));",
  "const films=labStatsDataset().entries.filter(f=>status==='both'||(status==='seen'?f.checked:!f.checked));",1)
 page=strings.Replace(page,
  "function labSmartGenreExposure(){\n const exposure=new Map();\n for(const film of parsed.entries){",
  "function labSmartGenreExposure(){\n const exposure=new Map();\n for(const film of labStatsDataset().entries){",1)
 // La ventana de filtros utiliza las mismas películas de Favoritos.
 page=strings.Replace(page,
  "const all=searchScope==='all'?filmMembershipIndex().rows:parsed.entries;",
  "const all=searchScope==='all'?filmMembershipIndex().rows:labStatsDataset().entries;",1)
 // Si un dispositivo remoto retira el último corazón, no dejar la vista
 // virtual abierta sin opción válida en el selector.
 page=strings.Replace(page,"function renderList(){labFiltersRefreshChoices();",
  "function renderList(){labReactionsByTMDB=null;labFiltersRefreshChoices();if(labFavoritesMode&&!labFavoriteCount()){labFavoritesMode=false;labRefreshStatsPanels();}",1)
 page=strings.Replace(page,"const locations=parsed.entries.map(f=>f.idx);",
  "if(labFavoritesMode)return labMoveFavorite(fromRawIndex,toRawIndex);\n const locations=parsed.entries.map(f=>f.idx);",1)
 page=strings.Replace(page,"function selectSearchScope(scope){",
  "function selectSearchScope(scope){if(scope==='all'&&labFavoritesMode){labFavoritesMode=false;labRefreshStatsPanels();}",1)
 page=strings.Replace(page,"const source=global?filmMembershipIndex().rows:parsed.entries;",
  "const source=labFavoritesMode&&!global?labFavoritesFilms():global?filmMembershipIndex().rows:parsed.entries;",1)
 page=strings.Replace(page,"const readonly=virtualMode&&!global;",
  "const readonly=labFavoritesMode||virtualMode&&!global;",1)
 // La lista automática solo aparece representada por un corazón
 // cuando la película está marcada. No cambia ni los números ni el lápiz.
 page=strings.Replace(page,
  "membership.append(caption,nums,edit);meta.appendChild(membership);",
  "membership.append(caption,nums);const favoriteBadge=labFavoriteMembershipButton(f);if(favoriteBadge)membership.appendChild(favoriteBadge);if(!labFavoritesMode)membership.appendChild(edit);meta.appendChild(membership);",
  1)
 page=strings.Replace(page,"b.append(togg,pill,labCreateDetailButton(f),fix);frag.appendChild(b);",
  "b.append(togg,labFilmRowActions(f,pill,labCreateDetailButton(f),fix));frag.appendChild(b);",1)
 page=strings.Replace(page,"$('libraryCount').textContent=global?matches.length+' resultados':parsed.entries.length+' películas';",
  "$('libraryCount').textContent=labFavoritesMode?matches.length+' favoritos':global?matches.length+' resultados':parsed.entries.length+' películas';",1)
 page=strings.Replace(page,
  "$('listPlatformCount').textContent=(global?'Todas las listas · ':'Lista actual · ')+matches.length+' resultados';",
  "$('listPlatformCount').textContent=(labFavoritesMode?'Favoritos · ':global?'Todas las listas · ':'Lista actual · ')+matches.length+' resultados';",1)
 boot:="updateListPlatformOptions();update();renderCollectionControls();updatePlatformProgress();void initLibrary().then("
 page=strings.Replace(page,boot,sections[1]+"\n"+boot,1)
 return page
}
