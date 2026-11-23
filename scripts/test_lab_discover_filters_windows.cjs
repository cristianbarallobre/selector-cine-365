'use strict';
const fs=require('node:fs');
const assert=require('node:assert/strict');
const vm=require('node:vm');
const ui=fs.readFileSync('lab_discover_windows.html','utf8');
const go=fs.readFileSync('lab_discover_windows.go','utf8');
const part=ui.split('<!-- LAB_DISCOVERY_SCRIPT -->');
assert.equal(part.length,2);
const js=part[1];
new vm.Script('(()=>{\n'+js+'\n})();',{filename:'lab_discover_windows.js'});

for(const id of ['labDiscoverGenre','labDiscoverDuration','labDiscoverFilterHint']) {
 assert.ok(ui.includes('id="'+id+'"'),'Falta el control '+id);
}
for(const value of ['any','short','medium','long'])assert.ok(ui.includes('value="'+value+'"'),'Falta duración '+value);
for(const id of ['878','35','10752','18','9648','27','16','28']){
 assert.ok(ui.includes('value="'+id+'"'),'Falta género TMDB '+id);
 assert.match(go,new RegExp('(?:,|\\s)'+id+':true'),'El servidor debe validar género '+id);
}
assert.match(go,/endpoint="\/discover\/movie"/,'TMDB debe filtrar en el servidor por categoría');
assert.match(go,/params\.Set\("with_genres",strconv\.Itoa\(genreID\)\)/);
assert.match(go,/params\.Set\("with_runtime\.lte","90"\)/);
assert.match(go,/params\.Set\("with_runtime\.gte","91"\)/);
assert.match(go,/params\.Set\("with_runtime\.lte","120"\)/);
assert.match(go,/params\.Set\("with_runtime\.gte","121"\)/);
assert.match(go,/if err!=nil\|\|!labDiscoverGenreIDs\[n\]/,'No aceptar IDs TMDB inventados');
assert.match(go,/default:fail\(400,"Duración no válida"\)/);
assert.match(js,/function labDiscoverFiltered\(\)/);
assert.match(js,/function labSyncDiscoverFilters\(\)/);
assert.match(js,/genre\.disabled=search;duration\.disabled=search/,'No presentar filtros engañosos en búsqueda por título');
assert.match(js,/Con filtros, Tendencias se ordena por popularidad/,'Explicar la diferencia frente a ranking semanal');
assert.match(js,/labResetBrowsePages\(\);labSyncDiscoverFilters\(\);void labBrowse\(\)/,'Cambiar filtros reinicia la paginación');
assert.match(js,/const cacheKey=path\+'\|'\+current/,'No mezclar páginas TMDB de filtros distintos');
assert.match(js,/Math\.max\(3,Math\.ceil/,'Reservar al menos tres filas visibles');
assert.match(js,/grid\.parentElement\?\.clientWidth/,'No fiarse del ancho del grid vacío');

const ids={labDiscoverGenre:{value:'any',disabled:false},labDiscoverDuration:{value:'any',disabled:false},
 labDiscoverFilterHint:{hidden:true,textContent:''}};
const functions=js.slice(js.indexOf('function labDiscoverFiltered(){'),js.indexOf('function labKindSelect(kind){'));
assert.ok(functions.startsWith('function labDiscoverFiltered(){'));
let reset=0,browses=0;const handlers={};
for(const id of ['labDiscoverGenre','labDiscoverDuration']){
 ids[id].addEventListener=(event,fn)=>{handlers[id]=fn};
}
const ctx={
 $:id=>{assert.ok(ids[id],id);return ids[id]},
 labQuery:'',labKind:'trending',
 labResetBrowsePages:()=>{reset++},
 labBrowse:()=>{browses++}
};
vm.createContext(ctx);vm.runInContext(functions,ctx);
ctx.labSyncDiscoverFilters();
assert.equal(ids.labDiscoverFilterHint.hidden,true);
ids.labDiscoverGenre.value='878';handlers.labDiscoverGenre();
assert.equal(reset,1);assert.equal(browses,1);
assert.equal(ids.labDiscoverFilterHint.hidden,false,'Tendencias filtradas muestra nota explicativa');
assert.equal(ctx.labDiscoverFiltered(),true);
ids.labDiscoverDuration.value='short';handlers.labDiscoverDuration();
assert.equal(reset,2);assert.equal(browses,2);
ctx.labQuery='Matrix';ctx.labSyncDiscoverFilters();
assert.equal(ctx.labDiscoverFiltered(),false,'Los filtros no se envían en búsquedas por título');
assert.equal(ids.labDiscoverDuration.disabled,true);
assert.equal(ids.labDiscoverGenre.disabled,true);
assert.match(ids.labDiscoverFilterHint.textContent,/búsqueda por título no los aplica/);
ctx.labQuery='';ctx.labKind='popular';ctx.labSyncDiscoverFilters();
assert.equal(ids.labDiscoverDuration.disabled,false);
assert.equal(ids.labDiscoverFilterHint.hidden,true);
console.log('OK: Descubrir con filtros reales TMDB, búsqueda clara, validación y rejilla de tres filas');
