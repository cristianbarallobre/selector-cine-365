'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');

const ui=fs.readFileSync('lab_discover_windows.html','utf8');
const go=fs.readFileSync('lab_discover_desktop.go','utf8');
const part=ui.split('<!-- LAB_DISCOVERY_SCRIPT -->');
assert.equal(part.length,2);
const js=part[1];
new vm.Script('(()=>{\n'+js+'\n})();',{filename:'lab_discover_windows.js'});

assert.ok(ui.includes('id="labDiscoverGenre"'),'Debe conservarse el filtro por género');
assert.ok(ui.includes('id="labDiscoverFilterHint"'),'Conservar la explicación de filtros');
assert.ok(!ui.includes('id="labDiscoverDuration"'),'Descubrir no presenta filtro de duración');
assert.ok(!js.includes('labDiscoverDuration'),'Descubrir no consulta filtros de duración');
for(const id of ['878','35','10752','18','9648','27','16','28']){
 assert.ok(ui.includes('value="'+id+'"'),'Falta género TMDB '+id);
 assert.match(go,new RegExp('(?:,|\\s)'+id+':true'),'El servidor debe validar género '+id);
}
assert.match(go,/endpoint="\/discover\/movie"/,'TMDB debe filtrar en el servidor por categoría');
assert.match(go,/params\.Set\("with_genres",strconv\.Itoa\(genreID\)\)/);
assert.doesNotMatch(go,/with_runtime\./,'Descubrir no debe filtrar por duración en el servidor');
assert.match(go,/if err!=nil\|\|!labDiscoverGenreIDs\[n\]/,'No aceptar IDs TMDB inventados');
assert.match(js,/function labDiscoverFiltered\(\)/);
assert.match(js,/function labSyncDiscoverFilters\(\)/);
assert.match(js,/genre\.disabled=search/,'No presentar filtros engañosos en búsqueda por título');
assert.match(js,/Con filtros, Tendencias se ordena por popularidad/,'Explicar la diferencia frente al ranking semanal');
assert.match(js,/labResetBrowsePages\(\);labSyncDiscoverFilters\(\);void labBrowse\(\)/,'Cambiar filtros reinicia la paginación');
assert.match(js,/cacheKey=path\+'\|'\+page/,'No mezclar páginas TMDB de filtros distintos');
assert.match(js,/Math\.max\(3,Math\.ceil/,'Reservar al menos tres filas visibles');
assert.match(js,/grid\.parentElement\?\.clientWidth/,'No fiarse del ancho del grid vacío');

const ids={
 labDiscoverGenre:{value:'any',disabled:false},
 labDiscoverFilterHint:{hidden:true,textContent:''}
};
const functions=js.slice(js.indexOf('function labDiscoverFiltered(){'),js.indexOf('function labKindSelect(kind){'));
assert.ok(functions.startsWith('function labDiscoverFiltered(){'));
let reset=0,browses=0;
const handlers={};
ids.labDiscoverGenre.addEventListener=(event,fn)=>{handlers[event]=fn};
const ctx={
 $:id=>{assert.ok(ids[id],id);return ids[id]},
 labQuery:'',labKind:'trending',
 labResetBrowsePages:()=>{reset++},
 labBrowse:()=>{browses++}
};
vm.createContext(ctx);vm.runInContext(functions,ctx);
ctx.labSyncDiscoverFilters();
assert.equal(ids.labDiscoverFilterHint.hidden,true);
ids.labDiscoverGenre.value='878';handlers.change();
assert.equal(reset,1);assert.equal(browses,1);
assert.equal(ids.labDiscoverFilterHint.hidden,false,'Tendencias filtradas muestra nota explicativa');
assert.equal(ctx.labDiscoverFiltered(),true);
ctx.labQuery='Matrix';ctx.labSyncDiscoverFilters();
assert.equal(ctx.labDiscoverFiltered(),false,'Los filtros no se envían en búsquedas por título');
assert.equal(ids.labDiscoverGenre.disabled,true);
assert.match(ids.labDiscoverFilterHint.textContent,/búsqueda por título no lo aplica/);
ctx.labQuery='';ctx.labKind='popular';ctx.labSyncDiscoverFilters();
assert.equal(ids.labDiscoverGenre.disabled,false);
assert.equal(ids.labDiscoverFilterHint.hidden,true);
console.log('OK: Descubrir filtra solo por género, conserva validación y rejilla de tres filas');
