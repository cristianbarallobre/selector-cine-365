//go:build windows || darwin

package main

import (
 "embed"
 "strings"
)

//go:embed lab_filters_windows.html
var labFiltersAssets embed.FS

func injectLabLibraryFilters(page string) string {
 raw,err:=labFiltersAssets.ReadFile("lab_filters_windows.html")
 if err!=nil{return page}
 parts:=strings.SplitN(string(raw),"<!-- LAB_FILTERS_HTML -->",2)
 if len(parts)!=2{return page}
 fragments:=strings.SplitN(parts[1],"<!-- LAB_FILTERS_SCRIPT -->",2)
 if len(fragments)!=2{return page}
 required:=[]string{
  "</head>",
  `<div class="search-scope" role="group"`,
  `const source=global?(needle?filmMembershipIndex().rows:[]):parsed.entries;`,
  `const matches=source.filter(f=>{`,
  `if(!global&&selectedListPlatforms.size)`,
  `const areaHeight=parent.clientHeight||370;`,
  `function renderList(){`,
  `function filmOrderEditable(){`,
  `function selectSearchScope(scope){`,
  `$('listPlatformButton').disabled=searchScope==='all';`,
  `['Películas vistas',seen.length],['Pendientes',parsed.pending.length]`,
  `<button class="nav-btn" type="button" data-view="instructions">❔ Instrucciones</button>`,
 }
 for _,check:=range required{if !strings.Contains(page,check){return page}}

 page=strings.Replace(page,"</head>",parts[0]+"\n</head>",1)
 page=strings.Replace(page,
  `<div class="search-scope" role="group"`,
  fragments[0]+`<div class="search-scope" role="group"`,1)
 page=strings.Replace(page,
  `<button class="nav-btn" type="button" data-view="instructions">❔ Instrucciones</button>`,
  `<button class="nav-btn lab-help-nav" type="button" data-view="instructions" title="Instrucciones" aria-label="Instrucciones">?</button>`,1)
 page=strings.Replace(page,
  `function renderList(){`,
  `function renderList(){labFiltersRefreshChoices();`,1)
 page=strings.Replace(page,
  `const source=global?(needle?filmMembershipIndex().rows:[]):parsed.entries;`,
  `const source=global?filmMembershipIndex().rows:parsed.entries;`,1)
 page=strings.Replace(page,
  `const matches=source.filter(f=>{`,
  `let matches=source.filter(f=>{`,1)
 page=strings.Replace(page,
  `if(!global&&selectedListPlatforms.size)`,
  `if(selectedListPlatforms.size)`,1)
 // Buscar el último "return true" dentro del predicado de renderList,
 // porque otras inyecciones LAB alteran el cuerpo antes de llegar a esta función.
 filterStart:=strings.Index(page,"let matches=source.filter(f=>{")
 heightStart:=strings.Index(page," const areaHeight=parent.clientHeight||370;")
 if filterStart<0||heightStart<=filterStart{return page}
 matchEnd:=strings.LastIndex(page[filterStart:heightStart],"return true;")
 if matchEnd<0{return page}
 pos:=filterStart+matchEnd
 page=page[:pos]+"return labFiltersMatches(f);"+page[pos+len("return true;"):]
 page=strings.Replace(page,
  " const areaHeight=parent.clientHeight||370;",
  " matches=labFiltersSort(matches);\n const areaHeight=parent.clientHeight||370;",1)
 page=strings.Replace(page,
  `global&&!needle?'Escribe para buscar en todas las listas':`,
  `global&&!source.length?'No hay películas en las listas':`,1)
 page=strings.Replace(page,
  `$('listPlatformCount').textContent=global?'Búsqueda global · '+matches.length+' resultados':'Filtradas: '+matches.length+' de '+parsed.entries.length;`,
  `$('listPlatformCount').textContent=(global?'Todas las listas · ':'Lista actual · ')+matches.length+' resultados';`,1)
 page=strings.Replace(page,
  `$('filmPrev').disabled=filmPage===0; $('filmNext').disabled=filmPage>=pageCount-1;`,
  `$('filmPrev').disabled=filmPage===0; $('filmNext').disabled=filmPage>=pageCount-1;labFiltersUpdateSummary(matches.length);`,1)
 page=strings.Replace(page,
  "function restartList(){",
  "function restartList(){labFiltersResetForList();",1)
 page=strings.Replace(page,
  `return !virtualMode&&searchScope==='current'&&filter==='all'&&!search.trim()&&!selectedListPlatforms.size;`,
  `return !virtualMode&&searchScope==='current'&&filter==='all'&&!search.trim()&&!selectedListPlatforms.size&&labFiltersCanReorder();`,1)
 page=strings.Replace(page,
  `if(searchScope==='all'){selectedListPlatforms.clear();closeListPlatformMenu();updateListPlatformOptions();}
 $('listPlatformButton').disabled=searchScope==='all';`,
  `closeListPlatformMenu();
 $('listPlatformButton').disabled=false;`,1)
 page=strings.Replace(page,
  `for(const [label,val] of [['Películas vistas',seen.length],['Pendientes',parsed.pending.length],['Horas registradas',minutes?Math.floor(minutes/60)+' h '+(minutes%60)+' min':'—'],['Fichas TMDB',known.filter(x=>x.metadata_version>=2).length]]){`,
  `for(const [label,val] of labStatsHeroValues(seen,known,minutes)){`,1)
 // La inicialización original llamaba a renderList() ANTES de que const
 // labFilterDecade/Genre/... hubiera sido inicializado. Eso provoca una
 // ReferenceError (temporal dead zone) e impide arrancar la aplicación.
 // Mover TODO el arranque al final, después de Descubrir, Sorteo, Red y Filtros.
 // Descubrir inserta labQueuePrefetchSoon() en esa misma línea.
 // Buscar el prefijo estable y mover la línea REAL completa, no una copia
 // antigua que no coincide después de las demás inyecciones LAB.
 bootPrefix:="updateListPlatformOptions();update();renderCollectionControls();updatePlatformProgress();void initLibrary().then("
 if strings.Count(page,bootPrefix)!=1{return page}
 bootStart:=strings.Index(page,bootPrefix)
 bootEnd:=strings.IndexByte(page[bootStart:],'\n')
 if bootEnd<0{return page}
 boot:=page[bootStart:bootStart+bootEnd]
 page=page[:bootStart]+page[bootStart+bootEnd:]
 end:=strings.LastIndex(page,"})();")
 script:=strings.LastIndex(page,"</script>")
 if end<0||script<0||end>=script{return page}
 return page[:end]+"\n"+fragments[1]+"\n"+boot+"\n"+page[end:]
}
