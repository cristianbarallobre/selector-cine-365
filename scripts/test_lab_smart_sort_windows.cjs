'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');

const ui=fs.readFileSync('lab_smart_sort_windows.html','utf8');
const injector=fs.readFileSync('lab_smart_sort_inject_windows.go','utf8');
const windows=fs.readFileSync('lab_smart_sort_windows.go','utf8');
const other=fs.readFileSync('lab_smart_sort_other.go','utf8');
const main=fs.readFileSync('main.go','utf8');
const discover=fs.readFileSync('lab_discover_windows.go','utf8');
const base=fs.readFileSync('plantilla.html','utf8');
const parts=ui.split('<!-- LAB_SMART_HTML -->');
assert.equal(parts.length,2);
const components=parts[1].split('<!-- LAB_SMART_SCRIPT -->');
assert.equal(components.length,2);
const script=components[1];
new vm.Script('(()=>{\n'+script+'\n})();',{filename:'lab_smart_sort_windows.js'});
assert.match(injector,/\/\/go:build windows/);
assert.match(other,/\/\/go:build !windows/);
assert.match(windows,/\/\/go:build windows/);
assert.match(injector,/if !appIsLab/);
assert.match(windows,/if !appIsLab/);
assert.match(main,/out = injectLabSmartSort\(out\)/);
assert.match(main,/a\.serveLabSmartCache\(w,r\)/);
assert.match(discover,/LAB-02\.6 WINDOWS/);
assert.match(windows,/!a\.auth\(r\)/,'Exigir autenticación local');
assert.match(windows,/c\.mu\.Lock\(\)/,'Leer la caché protegida');
assert.match(windows,/len\(reply\.Items\)>=2500/,'Acotar el tamaño de respuesta');
assert.ok(!base.includes('labSmartEnabled'),'La plantilla compartida no puede modificarse');
for(const id of ['labSmartBox','labSmartEnabled','labSmartMood','labSmartDuration','labSmartGenre','labSmartStatus']){
 assert.ok(ui.includes('id="'+id+'"'),'Falta '+id);
}
for(const v of ['balanced','quick','top','surprise'])assert.match(ui,new RegExp('value="'+v+'"'));
for(const v of ['any','short','medium','long'])assert.match(ui,new RegExp('value="'+v+'"'));
assert.match(script,/eligibleFilms\(\)/,'Debe respetar el filtro de plataformas y las pendientes');
assert.match(script,/\.checked/,'Debe distinguir pendientes de vistas');
assert.match(script,/labSmartHistory=.*\.slice\(0,8\)/,'No repetir últimas 8 elecciones si existen alternativas');
assert.match(script,/e\.stopImmediatePropagation\(\)/,'Interceptar el botón solo en modo inteligente');
assert.match(script,/if\(!labSmartEnable\.checked\)return/,'El sorteo clásico no debe bloquearse');
assert.match(script,/labSmartPickWeighted/,'Selección ponderada y aleatoria');
assert.match(script,/labSmartMetadata/,'Datos de TMDB');
assert.match(script,/getRandomValues/,'Fuente aleatoria criptográfica cuando exista');
assert.match(injector,/queueMicrotask\(labSmartUpdateCount\)/,'Actualizar filtros al cambiar de lista y plataformas');
assert.match(injector,/labSmartReason.*hidden/,'Borrar la explicación cuando el resultado se reinicia');
assert.match(injector,/Todas las películas pendientes.*Todas las plataformas/,'La etiqueta visual del selector debe decir Todas las plataformas');
assert.match(injector,/id="labSmartResultDetails"/,'La tarjeta debe recibir su botón Detalles');
assert.match(injector,/function updateChoicePlatforms\(\)\{labSmartSyncDetailsButton\(\);/,'Mostrar Detalles también tras el sorteo clásico');
assert.match(injector,/function resetChoice\(\)\{chosen=null;labSmartSyncDetailsButton\(\);/,'Ocultar Detalles cuando no haya película elegida');
assert.match(script,/labOpenFilmDetails\(chosen\)/,'Detalles debe reutilizar la ficha real de Descubrir/Biblioteca');
assert.match(ui,/lab-smart-details-btn\[hidden\]/,'No presentar el botón hasta que haya película elegida');
assert.ok(!base.includes('labSmartResultDetails'),'El botón es exclusivo Windows LAB');

function extract(start,end){
 const i=script.indexOf(start),j=script.indexOf(end,i+start.length);
 assert.ok(i>=0&&j>i,'No se pudo extraer '+start);
 return script.slice(i,j);
}
const candidatesCode=extract('function labSmartCandidates(){','function labSmartGenres(){');
const movies=[
 {title:'Corta',metadata:{runtime:80,genres:['Comedia']}},
 {title:'Límite',metadata:{runtime:90,genres:['Comedia']}},
 {title:'Casi corta',metadata:{runtime:91,genres:['Comedia']}},
 {title:'Lo que el viento se llevó',metadata:{runtime:238,genres:['Drama']}},
 {title:'Desconocida',metadata:{runtime:0,genres:[]}}
];
const ctx={
 eligibleFilms:()=>movies,
 labSmartMetadata:f=>f.metadata,
 labSmartGenre:{value:'any'},
 labSmartDuration:{value:'short'},
 labSmartMood:{value:'balanced'},
 labSmartNormalize:s=>s.toLowerCase()
};
vm.createContext(ctx);
vm.runInContext(candidatesCode,ctx);
assert.deepEqual(Array.from(ctx.labSmartCandidates(),x=>x.title),['Corta','Límite'],'El límite de 90 minutos incluye 90 y descarta 91 y 238');
ctx.labSmartDuration.value='any';ctx.labSmartGenre.value='drama';
assert.deepEqual(Array.from(ctx.labSmartCandidates(),x=>x.title),['Lo que el viento se llevó'],'Género debe respetarse');
ctx.labSmartGenre.value='any';
assert.equal(ctx.labSmartCandidates().length,5,'Sin filtros no excluir títulos sin datos TMDB');
ctx.labSmartMood.value='quick';
assert.deepEqual(Array.from(ctx.labSmartCandidates(),x=>x.title),['Corta','Límite'],
 'Algo corto debe descartar incluso una película de 238 minutos cuando Duración indica Cualquiera');
ctx.labSmartGenre.value='drama';
assert.equal(ctx.labSmartCandidates().length,0,'No usar películas largas para rellenar un filtro corto vacío');
ctx.labSmartGenre.value='any';ctx.labSmartMood.value='balanced';

const listeners={};
const eventsCode=extract("labSmartMood.addEventListener('change',","setInterval(");
const eventCtx={
 labSmartMood:{value:'balanced',addEventListener:(type,fn)=>{listeners.mood=fn;}},
 labSmartDuration:{value:'any',addEventListener:(type,fn)=>{listeners.duration=fn;}},
 labSmartGenre:{value:'any',addEventListener:(type,fn)=>{listeners.genre=fn;}},
 resetChoice:()=>{},labSmartUpdateCount:()=>{}
};
vm.createContext(eventCtx);vm.runInContext(eventsCode,eventCtx);
eventCtx.labSmartMood.value='quick';listeners.mood();
assert.equal(eventCtx.labSmartDuration.value,'short','Algo corto debe mostrar automáticamente Hasta 90 minutos');
eventCtx.labSmartDuration.value='any';listeners.duration();
assert.equal(eventCtx.labSmartMood.value,'balanced','Cambiar duración a Cualquiera no debe dejar seleccionado Algo corto');

const detailsCode=extract('function labSmartSyncDetailsButton(){',"const labSmartEnable=");
const selected={title:'Película elegida',year:2020};
let clicked=null,visibleClass=false;
const detailsButton={hidden:true,disabled:true,setAttribute(k,v){this[k]=v;},addEventListener(type,fn){this.clickHandler=fn;}};
const detailCtx={
 chosen:null,
 $:id=>id==='labSmartResultDetails'?detailsButton:id==='resultPoster'?{closest:()=>({classList:{toggle:(key,show)=>{visibleClass=show;}}})}:null,
 labOpenFilmDetails:film=>{clicked=film;}
};
vm.createContext(detailCtx);vm.runInContext(detailsCode,detailCtx);
detailCtx.labSmartSyncDetailsButton();
assert.equal(detailsButton.hidden,true,'Sin película no hay botón Detalles');
detailCtx.chosen=selected;detailCtx.labSmartSyncDetailsButton();
assert.equal(detailsButton.hidden,false,'El botón aparece tras un sorteo');
assert.equal(detailsButton.disabled,false,'El botón está activo con película elegida');
assert.equal(visibleClass,true,'La tarjeta reserva espacio para el botón');
detailsButton.clickHandler();assert.equal(clicked,selected,'Detalles debe abrir la película sorteada');
detailCtx.chosen=null;detailCtx.labSmartSyncDetailsButton();
assert.equal(detailsButton.hidden,true,'Reset elimina el botón para evitar detalles anteriores');

const picker=extract('function labSmartPickWeighted(items,weights){','function labSmartGenreExposure(){');
const fakeCrypto={getRandomValues:v=>{v[0]=0}};
const randomCtx={window:{crypto:fakeCrypto},crypto:fakeCrypto,Uint32Array,Math};
vm.createContext(randomCtx);
vm.runInContext(picker,randomCtx);
assert.equal(randomCtx.labSmartPickWeighted(['A','B'],[1,9]),'A','El primer tramo elige la primera');
randomCtx.window.crypto.getRandomValues=v=>{v[0]=0xffffffff};
assert.equal(randomCtx.labSmartPickWeighted(['A','B'],[1,9]),'B','El último tramo elige la segunda');

const weights=extract('function labSmartWeight(f,exposure){','function labSmartChoose(){');
const weightContext={
 labSmartMetadata:f=>f.metadata,
 labSmartMood:{value:'top'},
 labSmartNormalize:s=>s.toLowerCase()
};
vm.createContext(weightContext);
vm.runInContext(weights,weightContext);
const votes=new Map();
const low={metadata:{vote_count:800,vote_average:5.5,runtime:100,genres:['Drama']}};
const high={metadata:{vote_count:800,vote_average:8.4,runtime:100,genres:['Drama']}};
assert.ok(weightContext.labSmartWeight(high,votes)>weightContext.labSmartWeight(low,votes),'Las buenas notas deben tener más peso');
weightContext.labSmartMood.value='quick';
const short={metadata:{vote_count:800,vote_average:7.5,runtime:85,genres:['Drama']}};
const long={metadata:{vote_count:800,vote_average:7.5,runtime:165,genres:['Drama']}};
assert.ok(weightContext.labSmartWeight(short,votes)>weightContext.labSmartWeight(long,votes),'El modo corto debe favorecer duración corta');
console.log('OK: LAB-02.6; limite 90 min, sincronizacion controles, plataformas y Detalles en ambos sorteos');
