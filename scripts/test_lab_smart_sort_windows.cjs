'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');

const ui=fs.readFileSync('lab_smart_sort_windows.html','utf8');
const injector=fs.readFileSync('lab_smart_sort_inject_windows.go','utf8');
const windows=fs.readFileSync('lab_smart_sort_windows.go','utf8');
const other=fs.readFileSync('lab_smart_sort_other.go','utf8');
const winShell=fs.readFileSync('window_windows.go','utf8');
const main=fs.readFileSync('main.go','utf8');
const discover=fs.readFileSync('lab_discover_windows.go','utf8');
const base=fs.readFileSync('plantilla.html','utf8');
const parts=ui.split('<!-- LAB_SMART_HTML -->');
assert.equal(parts.length,2,'HTML debe tener marcador de controles');
const pieces=parts[1].split('<!-- LAB_SMART_SCRIPT -->');
assert.equal(pieces.length,2,'JS debe tener marcador propio');
const script=pieces[1];
new vm.Script('(()=>{\n'+script+'\n})();',{filename:'lab_smart_sort_windows.js'});
assert.match(injector,/\/\/go:build windows/);
assert.match(other,/\/\/go:build !windows/);
assert.match(windows,/\/\/go:build windows/);
assert.match(injector,/if !appIsLab/);
assert.match(windows,/if !appIsLab/);
assert.match(main,/out = injectLabSmartSort\(out\)/);
assert.match(main,/a\.serveLabSmartCache\(w,r\)/);
assert.match(discover,/LAB-02\.13 WINDOWS/);
assert.match(windows,/!a\.auth\(r\)/,'La caché de metadatos requiere token');
assert.match(windows,/c\.mu\.Lock\(\)/,'Caché local protegida');
assert.match(windows,/len\(reply\.Items\)>=2500/,'Acotar respuesta de caché');
assert.ok(!base.includes('labSmartEnabled'),'La plantilla compartida no debe cambiarse');
assert.ok(!base.includes('labSmartSeen'),'El filtro de vistas solo existe en Windows LAB');
for(const id of ['labSmartEnabled','labSmartSeen','labSmartDuration','labSmartGenre','labSmartStatus','labSmartBox']){
 assert.ok(ui.includes('id="'+id+'"'),'Falta '+id);
}
for(const option of ['pending','seen','both'])assert.match(ui,new RegExp('value="'+option+'"'));
for(const option of ['any','short','medium','long'])assert.match(ui,new RegExp('value="'+option+'"'));
assert.ok(!ui.includes('labSmartMood'),'Se debe quitar ¿Qué te apetece? y toda su lógica');
assert.ok(!ui.includes('labSmartInfo'),'Se deben quitar las aclaraciones de los filtros');
assert.ok(!ui.includes('lab-smart-explain'),'No debe quedar espacio de aclaraciones');
assert.match(script,/function labSmartPool\(\)/,'Filtrar por estado antes de género y duración');
assert.match(script,/status==='both'/,'Incluir películas vistas y pendientes si ambas');
assert.match(script,/status==='seen'\?f\.checked:!f\.checked/,'Separar vistas/no vistas');
assert.match(script,/selectedSortPlatform==='all'/,'Mantener plataforma seleccionada');
assert.match(script,/subscriptions\(r\)/,'Solo streaming permitido cuando se elige plataforma');
assert.match(script,/labSmartHistory=.*\.slice\(0,8\)/,'No repetir últimas 8 elecciones');
assert.match(script,/e\.stopImmediatePropagation\(\)/,'No disparar sorteo clásico al usar el inteligente');
assert.match(script,/if\(!labSmartEnable\.checked\)return/,'Sorteo clásico intacto al desactivar modo');
assert.match(script,/labSmartPickWeighted/,'Sorteo ponderado');
assert.match(script,/getRandomValues/,'Aleatoriedad criptográfica si está disponible');
assert.match(script,/labOpenFilmDetails\(chosen\)/,'Ficha de Detalles existente');
assert.match(injector,/labSmartSyncSeenAction\(\)/,'Actualiza el botón vista/pendiente');
assert.match(injector,/function resetChoice\(\)\{chosen=null;labSmartSyncDetailsButton\(\);labSmartSyncSeenAction\(\)/,'Reinicia botón al limpiar selección');
assert.match(injector,/Todas las películas pendientes.*Todas las plataformas/,'Etiqueta de plataforma');
assert.match(winShell,/swShowMaximized\s*=\s*3/,'La ventana puede mostrarse maximizada');
assert.match(winShell,/if appIsLab \{ showMode = uintptr\(swShowMaximized\) \}/,'Solo maximiza Windows LAB');
assert.match(winShell,/procShowWindow\.Call\(hwnd, showMode\)/,'Aplica el modo al mostrar WebView2');
assert.match(winShell,/if appIsLab \{\s*args = append\(args, "--start-maximized"\)/,'El navegador de respaldo también maximiza');
assert.match(winShell,/else \{\s*args = append\(args, "--window-size=1200,820"\)/,'El Windows estable mantiene sus medidas');

function extract(start,end){
 const i=script.indexOf(start),j=script.indexOf(end,i+start.length);
 assert.ok(i>=0&&j>i,'No se pudo extraer '+start);
 return script.slice(i,j);
}
// Los filtros deben actuar sobre el estado de visionado, además de duración,
// género y plataforma. La selección clásica usa la función original.
const candidatesCode=extract('function labSmartPool(){','function labSmartGenres(){');
const films=[
 {title:'Corta pendiente',checked:false,metadata:{runtime:80,genres:['Comedia']}},
 {title:'Corta vista',checked:true,metadata:{runtime:85,genres:['Ciencia ficción']}},
 {title:'Límite pendiente',checked:false,metadata:{runtime:90,genres:['Ciencia ficción']}},
 {title:'Casi corta vista',checked:true,metadata:{runtime:91,genres:['Drama']}},
 {title:'Muy larga vista',checked:true,metadata:{runtime:238,genres:['Ciencia ficción']}},
 {title:'Sin duración pendiente',checked:false,metadata:{runtime:0,genres:[]}}
];
const providers=new Map([
 ['Corta pendiente',{state:'ok',subscriptions:['Netflix']}],
 ['Corta vista',{state:'ok',subscriptions:['Netflix']}],
 ['Límite pendiente',{state:'ok',subscriptions:['Filmin']}],
 ['Casi corta vista',{state:'ok',subscriptions:['Filmin']}],
 ['Muy larga vista',{state:'ok',subscriptions:[]}],
]);
const context={
 parsed:{entries:films},
 platformData:providers,
 selectedSortPlatform:'all',
 providerKey:f=>f.title,
 subscriptions:r=>r.subscriptions,
 labSmartMetadata:f=>f.metadata,
 labSmartSeen:{value:'pending'},
 labSmartDuration:{value:'any'},
 labSmartGenre:{value:'any'},
 labSmartNormalize:s=>String(s).normalize('NFD').replace(/[\u0300-\u036f]/g,'').toLowerCase()
};
vm.createContext(context);
vm.runInContext(candidatesCode,context);
const names=()=>Array.from(context.labSmartCandidates(),f=>f.title);
assert.deepEqual(names(),['Corta pendiente','Límite pendiente','Sin duración pendiente'],'Por defecto solo no vistas');
context.labSmartSeen.value='seen';
assert.deepEqual(names(),['Corta vista','Casi corta vista','Muy larga vista'],'El modo Vistas incluye también películas de 238 minutos sin filtro');
context.labSmartSeen.value='both';
assert.equal(names().length,6,'Ambas permite todas las películas de la lista');
context.labSmartDuration.value='short';
assert.deepEqual(names(),['Corta pendiente','Corta vista','Límite pendiente'],'Hasta 90 min excluye 91, 238 y sin duración');
context.labSmartDuration.value='medium';
assert.deepEqual(names(),['Casi corta vista'],'De 91 a 120 min incluye 91');
context.labSmartDuration.value='long';
assert.deepEqual(names(),['Muy larga vista'],'Más de 120 min');
context.labSmartDuration.value='any';context.labSmartGenre.value='ciencia ficcion';
assert.deepEqual(names(),['Corta vista','Límite pendiente','Muy larga vista'],'Filtrado de género en ambas');
context.labSmartSeen.value='seen';context.labSmartDuration.value='short';
assert.deepEqual(names(),['Corta vista'],'Vistas + ciencia ficción + hasta 90');
context.labSmartGenre.value='any';context.labSmartDuration.value='any';
context.selectedSortPlatform='Netflix';
assert.deepEqual(names(),['Corta vista'],'Vistas + Netflix');
context.labSmartSeen.value='both';
assert.deepEqual(names(),['Corta pendiente','Corta vista'],'Ambas + Netflix');
context.selectedSortPlatform='any';
assert.equal(names().length,4,'Solo streaming: excluir películas sin ofertas registradas');
context.selectedSortPlatform='all';
assert.equal(names().length,6,'Todas las plataformas incluye títulos sin streaming');
assert.match(script,/status==='seen'\?'vistas':status==='both'\?'películas':'pendientes'/,'El contador se adapta al estado');

const controlEvents=extract("for(const id of ['labSmartSeen','labSmartDuration','labSmartGenre']){","setInterval(");
const handlers={},elements={};
for(const id of ['labSmartSeen','labSmartDuration','labSmartGenre']){
 elements[id]={addEventListener:(name,fn)=>{handlers[id]=fn;}};
}
let resets=0,updates=0;
const eventCtx={
 $:name=>elements[name],
 resetChoice:()=>{resets++;},
 labSmartUpdateCount:()=>{updates++;}
};
vm.createContext(eventCtx);vm.runInContext(controlEvents,eventCtx);
for(const id of Object.keys(elements))handlers[id]();
assert.equal(resets,3,'Cambiar estado, duración o género reinicia la elección');
assert.equal(updates,3,'Los tres filtros actualizan los resultados');

const detailsCode=extract('function labSmartSyncSeenAction(){','const labSmartEnable=');
const resultButton={hidden:true,disabled:true,setAttribute(name,value){this[name]=value;},addEventListener(type,fn){this.clickHandler=fn;}};
const watchedButton={textContent:''};
let clicked=null,visibleClass=false;
const detailCtx={
 chosen:null,
 $:id=>id==='markWatched'?watchedButton:id==='labSmartResultDetails'?resultButton:id==='resultPoster'?{closest:()=>({classList:{toggle:(key,show)=>{visibleClass=show;}}})}:null,
 labOpenFilmDetails:film=>{clicked=film;}
};
vm.createContext(detailCtx);vm.runInContext(detailsCode,detailCtx);
detailCtx.labSmartSyncSeenAction();detailCtx.labSmartSyncDetailsButton();
assert.equal(watchedButton.textContent,'✓ Marcar como vista','Valor inicial');
assert.equal(resultButton.hidden,true,'Detalles oculto sin título');
const seen={title:'Película vista',checked:true};
detailCtx.chosen=seen;detailCtx.labSmartSyncSeenAction();detailCtx.labSmartSyncDetailsButton();
assert.equal(watchedButton.textContent,'↩ Marcar como pendiente','No marcar otra vez como vista una película ya vista');
assert.equal(resultButton.hidden,false,'Detalles visible tras sorteo');
assert.equal(visibleClass,true,'Dejar espacio al botón en la tarjeta');
resultButton.clickHandler();assert.equal(clicked,seen,'Detalles abre el título elegido');
detailCtx.chosen={title:'Película pendiente',checked:false};
detailCtx.labSmartSyncSeenAction();
assert.equal(watchedButton.textContent,'✓ Marcar como vista');
detailCtx.chosen=null;detailCtx.labSmartSyncDetailsButton();
assert.equal(resultButton.hidden,true,'Limpiar elección también oculta Detalles');

const pickCode=extract('function labSmartPickWeighted(items,weights){','function labSmartGenreExposure(){');
const fakeCrypto={getRandomValues:v=>{v[0]=0}};
const randCtx={window:{crypto:fakeCrypto},crypto:fakeCrypto,Uint32Array,Math};
vm.createContext(randCtx);vm.runInContext(pickCode,randCtx);
assert.equal(randCtx.labSmartPickWeighted(['A','B'],[1,9]),'A');
fakeCrypto.getRandomValues=v=>{v[0]=0xffffffff};
assert.equal(randCtx.labSmartPickWeighted(['A','B'],[1,9]),'B');

const weights=extract('function labSmartWeight(f,exposure){','function labSmartChoose(){');
const weightCtx={labSmartMetadata:f=>f.metadata,labSmartNormalize:s=>s.toLowerCase()};
vm.createContext(weightCtx);vm.runInContext(weights,weightCtx);
const low={metadata:{vote_count:800,vote_average:5.3,genres:['Drama']}};
const high={metadata:{vote_count:800,vote_average:8.4,genres:['Drama']}};
assert.ok(weightCtx.labSmartWeight(high,new Map())>weightCtx.labSmartWeight(low,new Map()),'Mantener ponderación por calidad');

console.log('OK: LAB-02.13 Windows; estados pendientes/vistas/ambas, duración, género, plataforma, Detalles, maximizado e aislamiento');
