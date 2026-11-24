'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');
const {performance}=require('node:perf_hooks');

const source=fs.readFileSync('plantilla.html','utf8');
const perf=fs.readFileSync('lab_search_perf_windows.js','utf8');
const injector=fs.readFileSync('lab_search_perf_desktop.go','utf8');
const main=fs.readFileSync('main.go','utf8');
new vm.Script('(()=>{\n'+perf+'\n})();',{filename:'lab_search_perf_windows.js'});
assert.match(injector,/function filmMembershipIndex\(\)\{return labGlobalMembershipIndex\(\);/);
assert.match(injector,/for\(const f of labIndexedEntries\(list\)\)/);
assert.match(injector,/labShouldRefreshFilterChoices\(\)/);
assert.match(injector,/labScheduleLibrarySearch\(e.target.value\)/);
assert.match(injector,/labCancelLibrarySearch\(\)/);
assert.match(injector,/labSmartCache=data;labMetadataVersion\+\+;/);
assert.match(main,/out = injectLabSearchPerformance\(out\)/);
assert.match(perf,/const labSearchDelayMs=150/);
assert.doesNotMatch(source,/labGlobalMembershipIndex\(/,'Plantilla estable sin modificar');

let scans=0,renders=0,clock=1000;
const parse=text=>{
 scans++;
 const entries=[];
 let idx=0;
 for(const line of text.split('\n')){
  const mark=/^([☐☑]) ([^\n]+?)(?: \((\d{4})\))?$/.exec(line);
  if(mark)entries.push({idx,checked:mark[1]==='☑',title:mark[2],name:mark[2],year:mark[3]||''});
  idx++;
 }
 return {entries};
};
const identityMap=new Map(),platformData=new Map();
const lists=[];
for(let n=0;n<40;n++){
 const rows=[];
 for(let i=0;i<250;i++){
  const number=n*250+i;
  rows.push((i%5===0?'☑ ':'☐ ')+'Película '+number+' ('+(2010+number%16)+')');
 }
 lists.push({id:'lista-'+n,name:'Lista '+(n+1),text:rows.join('\n')});
}
const ctx={
 Map,Set,String,Number,Math,Array,Object,
 parse,identityMap,platformData,libraryState:{lists,active_id:'lista-0'},
 providerKey:f=>f.title.toLowerCase()+'|'+(Number(f.year)||0),
 filmIdentity:f=>String(f.title).normalize('NFD').replace(/[\u0300-\u036f]/g,'').toLowerCase().replace(/[^a-z0-9]+/g,' ').trim(),
 membershipCache:null,searchScope:'all',parsed:{entries:[]},virtualMode:false,
 performance, setTimeout:()=>0,clearTimeout:()=>{},
 search:'',filmPage:3,renderList:()=>{renders++;}
};
vm.createContext(ctx);
const start=source.indexOf('function filmKey(f){');
const end=source.indexOf('function numberedListName(l){',start);
assert.ok(start>0&&end>start,'No se encuentra el índice real en plantilla.html');
let code=source.slice(start,end);
code=code.replace('function filmMembershipIndex(){','function filmMembershipIndex(){return labGlobalMembershipIndex();');
code=code.replace('for(const f of parse(list.text).entries){const key=filmKey(f);','for(const f of labIndexedEntries(list)){const key=filmKey(f);');
assert.match(code,/for\(const f of labIndexedEntries\(list\)\)/);
vm.runInContext(code+'\n'+perf,ctx,{timeout:10000});
const run=code=>vm.runInContext(code,ctx,{timeout:10000});
const memory0=process.memoryUsage().heapUsed;
const first=performance.now();
const index=run('filmMembershipIndex()');
const coldMs=performance.now()-first;
assert.equal(index.rows.length,10000);
assert.equal(scans,40,'Solo analizar 40 listas, no 40 por cada búsqueda');
assert.equal(index.byKey.size,10000);
assert.equal(index.listNames.size,40);
const buildVersion=run('labGlobalIndexGeneration');
const warmStart=performance.now();
for(let i=0;i<100;i++)assert.equal(run('filmMembershipIndex()'),index);
const warmMs=performance.now()-warmStart;
assert.equal(scans,40,'Buscar varias veces no vuelve a analizar ningún TXT');
assert.equal(run('labGlobalIndexGeneration'),buildVersion,'El índice se reutiliza incluso después de 100 consultas');
ctx.membershipCache=null;
assert.equal(run('filmMembershipIndex()'),index,'No reconstruir si scheduleSave invalida una caché antigua');
assert.equal(scans,40);

const changed=lists[12];
changed.text+='\n☑ Película nueva (2026)';
let updated=run('filmMembershipIndex()');
assert.equal(scans,41,'Editar una lista solo debe analizar su TXT');
assert.equal(updated.rows.length,10001);
assert.notEqual(updated,index,'El agregado sí debe reflejar las modificaciones reales');
assert.ok(updated.byKey.has('pelicula nueva|2026'));
assert.equal(run('labGlobalIndexGeneration'),buildVersion+1);

changed.name='Lista con nombre nuevo';
updated=run('filmMembershipIndex()');
assert.equal(scans,41,'Renombrar no debe reanalizar películas');
assert.equal(updated.byKey.get('pelicula nueva|2026')._lists[0].name,'Lista con nombre nuevo');
assert.equal(updated.byKey.get('pelicula nueva|2026')._lists[0].n,13);

const other=lists[13];
other.text+='\n☐ Título alternativo (2026)';
updated=run('filmMembershipIndex()');
assert.equal(scans,42);
assert.equal(updated.rows.length,10002);
ctx.identityMap.set('titulo alternativo|2026',777);
ctx.identityMap.set('pelicula nueva|2026',777);
const combined=run('filmMembershipIndex()');
assert.equal(scans,42,'Correcciones de TMDB no reparsean TXT');
assert.equal(combined.rows.length,10001,'Dos títulos con mismo ID TMDB son una sola película');
const merged=combined.byKey.get(combined.groups.get('pelicula nueva|2026'));
assert.equal(merged._lists.length,2,'Conservar las listas de títulos equivalentes');
assert.equal(merged.checked,true,'Si una copia está vista, la global también');
assert.equal(run('labShouldRefreshFilterChoices()'),true);
assert.equal(run('labShouldRefreshFilterChoices()'),false);
ctx.platformData.set('x',{tmdb_id:0,genres:['Comedia']});
assert.equal(run('labShouldRefreshFilterChoices()'),true,'Metadatos nuevos actualizan selector de géneros');
assert.equal(run('labShouldRefreshFilterChoices()'),false);
ctx.identityMap.delete('titulo alternativo|2026');
const unmerged=run('filmMembershipIndex()');
assert.equal(unmerged.rows.length,10002,'Deshacer un ID TMDB desagrupa los registros');
assert.equal(scans,42);

const moved=lists.splice(13,1)[0];lists.splice(0,0,moved);
updated=run('filmMembershipIndex()');
assert.equal(scans,42,'Cambiar de orden las listas no debe reanalizar sus películas');
assert.equal(updated.byKey.get('titulo alternativo|2026')._lists[0].n,1,'Números de lista coherentes tras reordenar');
lists.splice(0,1);
updated=run('filmMembershipIndex()');
assert.equal(updated.rows.length,9751,'Borrar una lista retira sus 251 películas sin mantener resultados obsoletos');
assert.equal(run('labListParseCache.size'),39,'La caché tampoco mantiene TXT de listas eliminadas');
assert.equal(scans,42);

const memoryGrowth=process.memoryUsage().heapUsed-memory0;
assert.ok(coldMs<10000,'Construcción de 10.000 películas demasiado lenta: '+coldMs+' ms');
assert.ok(warmMs<2000,'Cien búsquedas repetidas demasiado lentas: '+warmMs+' ms');
assert.ok(memoryGrowth<160*1024*1024,'Índice global ocupa demasiada RAM: '+memoryGrowth);

// Simular pulsaciones rápidas: solo la última provoca render tras 150 ms.
// Cancelar al cambiar lista o ámbito impide renders tardíos.
const callbacks=new Map();
let timer=0;
ctx.setTimeout=(fn,delay)=>{const id=++timer;callbacks.set(id,{fn,delay});return id;};
ctx.clearTimeout=id=>callbacks.delete(id);
run("labScheduleLibrarySearch('i')");
run("labScheduleLibrarySearch('in')");
run("labScheduleLibrarySearch('inte')");
assert.equal(ctx.search,'inte','Actualizar valor de búsqueda inmediatamente');
assert.equal(ctx.filmPage,0);
assert.equal(callbacks.size,1,'Solo mantener un temporizador');
assert.equal([...callbacks.values()][0].delay,150);
assert.equal(renders,0,'No renderizar con cada tecla');
const fn=[...callbacks.values()][0].fn;callbacks.clear();fn();
assert.equal(renders,1,'Una única búsqueda tras finalizar la escritura');
run("labScheduleLibrarySearch('otra')");
run('labCancelLibrarySearch()');
assert.equal(callbacks.size,0,'Descartar búsqueda al cambiar de lista');

console.log(JSON.stringify({
 result:'OK',listsTested:40,filmOccurrences:10000,
 buildMs:Math.round(coldMs),repeat100Ms:Math.round(warmMs),
 memoryGrowthMiB:+(memoryGrowth/1024/1024).toFixed(1),
 parsedInitial:40,parsedAfterTwoEdits:42,debounceMs:150
}));
