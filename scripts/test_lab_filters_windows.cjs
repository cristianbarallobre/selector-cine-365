'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');

const ui=fs.readFileSync('lab_filters_windows.html','utf8');
const injector=fs.readFileSync('lab_filters_inject_windows.go','utf8');
const stub=fs.readFileSync('lab_filters_other.go','utf8');
const base=fs.readFileSync('plantilla.html','utf8');
const core=fs.readFileSync('main.go','utf8');
const parts=ui.split('<!-- LAB_FILTERS_HTML -->');
assert.equal(parts.length,2);
const more=parts[1].split('<!-- LAB_FILTERS_SCRIPT -->');
assert.equal(more.length,2);
const js=more[1];
new vm.Script('(()=>{\n'+js+'\n})();',{filename:'lab_filters_windows.js'});
assert.match(injector,/\/\/go:build windows/);
assert.match(stub,/\/\/go:build !windows/);
assert.match(injector,/if !appIsLab/);
assert.match(core,/out = injectLabLibraryFilters\(out\)/);
assert.match(injector,/function renderList\(\)\{labFiltersRefreshChoices\(\)/);
assert.match(injector,/global\?filmMembershipIndex\(\)\.rows:parsed\.entries/,'Todas las listas sin texto de búsqueda');
assert.match(injector,/if\(selectedListPlatforms\.size\)/,'Plataformas también entre todas las listas');
assert.match(injector,/return labFiltersMatches\(f\)/);
assert.match(injector,/matches=labFiltersSort\(matches\)/);
assert.match(injector,/labFiltersUpdateSummary\(matches.length\)/);
assert.match(injector,/labFiltersCanReorder\(\)/,'Prohibir arrastrar títulos mientras está activada ordenación visual');
assert.match(injector,/labStatsHeroValues\(seen,known,minutes\)/,'Resumen de Estadísticas mejorado');
assert.match(injector,/data-view="instructions" title="Instrucciones" aria-label="Instrucciones">\?<\/button>/,'Solo ? en ayuda');
for(const id of ['labLibraryToolbar','labFilterOpen','labFiltersOverlay','labFiltersDialog','labFilterScopeSlot','labFilterStatusSlot','labFilterPlatformSlot','labFilterDecade','labFilterGenre','labFilterDuration','labFilterOrder','labFilterResetAll','labFilterDone']){
 assert.ok(ui.includes('id="'+id+'"'),'Falta '+id);
}
assert.ok(!base.includes('labFiltersOverlay'),'No editar la plantilla estable');
assert.match(ui,/id="labFilterAllPlatforms" type="checkbox" checked/,'Todas las plataformas seleccionada por defecto');
assert.match(ui,/id="labFilterVisiblePlatforms"/,'Lista de plataformas visible en el panel');
assert.match(ui,/id="labFilterResetAll"/,'Limpiar filtros disponible sin desplegar el panel');
assert.ok(ui.indexOf('id="labFilterResetAll"')<ui.indexOf('id="labFiltersOverlay"'),
 'El botón Limpiar filtros debe quedar en la barra superior, no dentro del diálogo');
assert.match(ui,/#labFiltersDialog \.library-platform-row\{display:none!important\}/,
 'No debe visualizarse el dropdown anterior de plataformas');
assert.match(js,/\$\('labFilterVisiblePlatforms'\)\.appendChild\(\$\('listPlatformOptions'\)\)/,
 'Se reutilizan los checkboxes originales sin nuevos desplegables');
assert.match(js,/function labFiltersUpdatePlatformSelection\(\)/);
assert.match(js,/labFiltersUpdatePlatformSelection\(\);renderList\(\);/,'Actualizar resultados al elegir Todas');
assert.match(ui,/#view-library \.lab-library-reset-button/);
assert.match(ui,/\.lab-help-nav/);
assert.match(ui,/#view-statistics \.stats-hero\{grid-template-columns:repeat\(5/,'Cinco datos útiles');
assert.match(js,/document\.body\.appendChild\(labFiltersOverlay\)/,'El panel flotante no se recorta por el panel de biblioteca');
assert.match(js,/appendChild\(\$\('view-library'\)\.querySelector\('\.filters'\)\)/,'Reutiliza los filtros antiguos y sus eventos');
assert.match(js,/appendChild\(\$\('view-library'\)\.querySelector\('\.search-scope'\)\)/,'Reutiliza selector Lista actual/Todas las listas');
assert.match(js,/appendChild\(\$\('view-library'\)\.querySelector\('\.library-platform-row'\)\)/,'Reutiliza plataformas existentes');

// La opción Todas representa un filtro vacío y cambia con los checkboxes existentes.
const allCheckbox={checked:false},chosenCaption={textContent:''};
const platformSelection=new Set();
const platformCtx={
 selectedListPlatforms:platformSelection,
 $:id=>id==='labFilterAllPlatforms'?allCheckbox:id==='labFilterPlatformSelected'?chosenCaption:null
};
vm.createContext(platformCtx);
const selectionBlock=ui.split('<!-- LAB_FILTERS_SCRIPT -->')[1];
const fnCode=selectionBlock.match(/function labFiltersUpdatePlatformSelection\(\)\{[\s\S]*?\n\}/)?.[0];
assert.ok(fnCode,'Falta función de sincronización de Todas');
vm.runInContext(fnCode,platformCtx);
platformCtx.labFiltersUpdatePlatformSelection();
assert.equal(allCheckbox.checked,true,'Todas marcada cuando no se filtra');
platformSelection.add('Netflix');platformCtx.labFiltersUpdatePlatformSelection();
assert.equal(allCheckbox.checked,false,'Todas se desmarca al elegir Netflix');
assert.equal(chosenCaption.textContent,'1 seleccionada');
platformSelection.clear();platformCtx.labFiltersUpdatePlatformSelection();
assert.equal(allCheckbox.checked,true,'Limpiar vuelve a seleccionar Todas');

function between(a,b){const i=js.indexOf(a),j=js.indexOf(b,i+a.length);assert.ok(i>=0&&j>i,a);return js.slice(i,j);}
const data=between('function labFilterNormalize(value){','function labFiltersAreDefault(){');
const dctx={
 labFilterDecade:{value:'any'},labFilterGenre:{value:'any'},labFilterDuration:{value:'any'},labFilterOrder:{value:'original'},
 labSmartMetadata:f=>f.meta,
 filmMembershipIndex:()=>({rows:[]}),
 parsed:{entries:[]},searchScope:'current',document:{createElement:()=>({})},
 Number,Math,Map,Set,String
};
vm.createContext(dctx);vm.runInContext(data,dctx);
const items=[
 {title:'B película',year:'1968',checked:true,meta:{runtime:90,genres:['Bélica'],vote_count:120,vote_average:7.2}},
 {title:'A película',year:'1972',checked:false,meta:{runtime:95,genres:['Comedia'],vote_count:200,vote_average:8.3}},
 {title:'C película',year:'1964',checked:true,meta:{runtime:130,genres:['Bélica'],vote_count:0,vote_average:0}},
 {title:'D película',year:'',checked:false,meta:{runtime:0,genres:[],vote_count:0,vote_average:0}}
];
const names=ls=>Array.from(ls,x=>x.title);
assert.equal(items.filter(dctx.labFiltersMatches).length,4,'Sin filtros deben verse incluso títulos sin metadatos');
dctx.labFilterDecade.value='1960s';assert.deepEqual(names(items.filter(dctx.labFiltersMatches)),['B película','C película']);
dctx.labFilterDecade.value='any';dctx.labFilterGenre.value='belica';
assert.deepEqual(names(items.filter(dctx.labFiltersMatches)),['B película','C película'],'Género sin acentos');
dctx.labFilterGenre.value='any';dctx.labFilterDuration.value='short';
assert.deepEqual(names(items.filter(dctx.labFiltersMatches)),['B película'],'Hasta 90 excluye 95, 130 y duración desconocida');
dctx.labFilterDuration.value='medium';
assert.deepEqual(names(items.filter(dctx.labFiltersMatches)),['A película']);
dctx.labFilterDuration.value='long';
assert.deepEqual(names(items.filter(dctx.labFiltersMatches)),['C película']);
dctx.labFilterDuration.value='any';
dctx.labFilterOrder.value='title-asc';assert.deepEqual(names(dctx.labFiltersSort([...items])),['A película','B película','C película','D película']);
dctx.labFilterOrder.value='title-desc';assert.deepEqual(names(dctx.labFiltersSort([...items])),['D película','C película','B película','A película']);
dctx.labFilterOrder.value='year-new';assert.deepEqual(names(dctx.labFiltersSort([...items])),['A película','B película','C película','D película']);
dctx.labFilterOrder.value='year-old';assert.deepEqual(names(dctx.labFiltersSort([...items])),['C película','B película','A película','D película']);
dctx.labFilterOrder.value='rating-desc';assert.deepEqual(names(dctx.labFiltersSort([...items])),['A película','B película','C película','D película']);
dctx.labFilterOrder.value='runtime-asc';assert.deepEqual(names(dctx.labFiltersSort([...items])),['B película','A película','C película','D película']);
dctx.labFilterOrder.value='runtime-desc';assert.deepEqual(names(dctx.labFiltersSort([...items])),['C película','A película','B película','D película']);
dctx.labFilterOrder.value='watched-first';assert.deepEqual(names(dctx.labFiltersSort([...items])).slice(0,2),['B película','C película']);
dctx.labFilterOrder.value='pending-first';assert.deepEqual(names(dctx.labFiltersSort([...items])).slice(0,2),['A película','D película']);
dctx.labFilterOrder.value='original';assert.deepEqual(names(dctx.labFiltersSort([...items])),names(items),'No alterar el orden original');
const heroCode=js.slice(js.indexOf('function labStatsHeroValues('));
vm.runInContext(heroCode,dctx);
const watched=items.filter(x=>x.checked);
const report=Array.from(dctx.labStatsHeroValues(watched,watched.map(x=>x.meta),220),x=>Array.from(x));
assert.deepEqual(report,[['Películas vistas',2],['Horas registradas','3 h 40 min'],['Nota media TMDB','★ 7.2'],['Década más vista','1960s'],['Género más visto','Bélica']]);
assert.ok(!report.some(([title])=>title==='Pendientes'||title==='Fichas TMDB'),'Eliminar métricas tachadas');
console.log('OK: panel filtros completo, listas globales, plataforma, duración, ordenación y nuevas estadísticas Windows LAB');
