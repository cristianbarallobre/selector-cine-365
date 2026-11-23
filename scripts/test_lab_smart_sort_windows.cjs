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
assert.match(discover,/LAB-02\.5 WINDOWS/);
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

function extract(start,end){
 const i=script.indexOf(start),j=script.indexOf(end,i+start.length);
 assert.ok(i>=0&&j>i,'No se pudo extraer '+start);
 return script.slice(i,j);
}
const candidatesCode=extract('function labSmartCandidates(){','function labSmartGenres(){');
const movies=[
 {title:'Corta',metadata:{runtime:80,genres:['Comedia']}},
 {title:'Larga',metadata:{runtime:155,genres:['Drama']}},
 {title:'Desconocida',metadata:{runtime:0,genres:[]}}
];
const ctx={
 eligibleFilms:()=>movies,
 labSmartMetadata:f=>f.metadata,
 labSmartGenre:{value:'any'},
 labSmartDuration:{value:'short'},
 labSmartNormalize:s=>s.toLowerCase()
};
vm.createContext(ctx);
vm.runInContext(candidatesCode,ctx);
assert.deepEqual(Array.from(ctx.labSmartCandidates(),x=>x.title),['Corta'],'El filtro corto exige duración conocida');
ctx.labSmartDuration.value='any';ctx.labSmartGenre.value='drama';
assert.deepEqual(Array.from(ctx.labSmartCandidates(),x=>x.title),['Larga'],'Género debe respetarse');
ctx.labSmartGenre.value='any';
assert.equal(ctx.labSmartCandidates().length,3,'Sin filtros no excluir títulos sin datos TMDB');

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
console.log('OK: sorteo inteligente Windows LAB-02.5; algoritmos, filtros, seguridad, sintaxis y aislamiento');
