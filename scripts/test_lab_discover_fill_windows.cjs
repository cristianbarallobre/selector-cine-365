'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');

const html=fs.readFileSync('lab_discover_windows.html','utf8');
const script=html.split('<!-- LAB_DISCOVERY_SCRIPT -->')[1];
assert.ok(script);
new vm.Script('(()=>{\n'+script+'\n})();',{filename:'lab_discover_windows.js'});
assert.match(script,/function labDiscoverFillTarget\(\)/);
assert.match(script,/async function labBrowse\(\)/);
assert.match(script,/labMoviePool\.slice\(startIndex,required\)/,'No mostrar películas sobrantes de TMDB en la fila incompleta');
assert.match(script,/labMoviePool\.push\(movie\)/,'Conservar títulos sobrantes para la página siguiente');
assert.match(script,/labTMDBNextPage=page\+1/,'Avanzar TMDB sin repetir sus páginas');
assert.match(script,/labMovieResponseCache\.set\(cacheKey,data\)/,'Conservar resultados en la sesión');
assert.match(script,/window\.addEventListener\('resize',labDiscoverOnResize\)/);
assert.match(script,/labPage=Math\.floor\(oldIndex\/target\)\+1/,'Reagrupar tras cambiar dimensiones');
assert.match(script,/columns\*Math\.max\(1,Math\.min\(rows/,'Las páginas deben ser múltiplos de columnas');

const start=script.indexOf('function labDiscoverFillTarget(){');
const end=script.indexOf('function labCard(film){',start);
assert.ok(start>=0&&end>start);
const code=script.slice(start,end);

const grid={
 items:[],clientWidth:1200,
 get childElementCount(){return this.items.length;},
 replaceChildren(){this.items=[];},
 appendChild(card){this.items.push(card);},
 getBoundingClientRect(){return{top:250,bottom:250}}
};
const el={
 labMoviesGrid:grid,
 labDiscoverPager:{hidden:true},
 labResultsHeading:{textContent:''},
 labDiscoverStatus:{textContent:''},
 labResultsCount:{textContent:''},
 labPageLabel:{textContent:''},
 labPrevPage:{disabled:false},
 labNextPage:{disabled:false},
 labDiscoverGenre:{value:'any',disabled:false}
};
let calls=[],shouldError=false,itemsPerTMDBPage=20,pages=2;
const context={
 $:id=>{assert.ok(id in el,'Falta elemento '+id);return el[id];},
 window:{
  getComputedStyle:()=>({gridTemplateColumns:'none',columnGap:'14px'}),
  addEventListener:()=>{}
 },
 labView:{getBoundingClientRect:()=>({bottom:900}),hidden:false},
 labStatus:s=>{el.labDiscoverStatus.textContent=s;},
 labSyncDiscoverFilters:()=>{},
 labCard:x=>x,
 labKind:'upcoming',labQuery:'',labPage:1,labPages:1,
 labMoviePool:[],labMovieIDs:new Set(),labTMDBNextPage:1,labTMDBTotal:0,
 labMovieResponseCache:new Map(),labLastFillTarget:0,labResizeTimer:null,labRequestSerial:0,
 labGet:async url=>{
  const match=url.match(/[?&]page=(\d+)/);
  assert.ok(match,url);
  const page=Number(match[1]);calls.push(page);
  if(shouldError){shouldError=false;throw new Error('La conexión falló');}
  return {
   page,total_pages:pages,total_results:pages*itemsPerTMDBPage,
   results:Array.from({length:itemsPerTMDBPage},(_,i)=>({
    id:page*100+i+1,title:'Película '+page+' #'+(i+1)
   }))
  };
 },
 setTimeout:()=>{},clearTimeout:()=>{},
 activeView:'discover',Number,Math,Map,Set
};
vm.createContext(context);
vm.runInContext(code,context);

async function test(){
 const target=context.labDiscoverFillTarget();
 assert.equal(target,18,'Con seis columnas y tres filas visibles caben 18 películas');
 await context.labBrowse();
 assert.deepEqual(calls,[1],'TMDB devuelve 20 títulos, solo hace falta la primera petición');
 assert.equal(grid.items.length,18,'La página visual debe tener 3 filas COMPLETAS, no 20 títulos (6+6+6+2)');
 assert.equal(context.labMoviePool.length,20,'Guardar las dos películas que sobran para la página siguiente');
 assert.deepEqual(grid.items.map(x=>x.id),Array.from({length:18},(_,i)=>101+i));
 assert.equal(el.labNextPage.disabled,false,'Quedan películas en el buffer y en TMDB');
 assert.match(el.labPageLabel.textContent,/Películas 1–18 de 40/);
 const firstIDs=grid.items.map(x=>x.id);

 context.labPage=2;
 await context.labBrowse();
 assert.deepEqual(calls,[1,2],'La segunda página debe continuar donde acabó la primera');
 assert.equal(grid.items.length,18,'Completar la segunda página visual también');
 assert.deepEqual(grid.items.slice(0,2).map(x=>x.id),[119,120],
   'No perder las dos películas sobrantes de la primera página TMDB');
 assert.ok(grid.items.every(x=>!firstIDs.includes(x.id)),'No repetir películas entre páginas visuales');
 const secondIDs=grid.items.map(x=>x.id);
 assert.match(el.labPageLabel.textContent,/Películas 19–36 de 40/);

 context.labPage=3;
 await context.labBrowse();
 assert.equal(grid.items.length,4,'Al agotar TMDB pueden quedar menos de una fila, pero no películas ocultas');
 assert.equal(el.labNextPage.disabled,true,'No inventar una página siguiente');
 assert.match(el.labDiscoverStatus.textContent,/últimos resultados disponibles/);
 const all=[...firstIDs,...secondIDs,...grid.items.map(x=>x.id)];
 assert.equal(new Set(all).size,40,'No repetir ni perder ninguna de las 40 películas');
 assert.deepEqual(calls,[1,2],'Las tres páginas visuales usan solo dos solicitudes TMDB');

 context.labPage=1;
 await context.labBrowse();
 assert.deepEqual(calls,[1,2],'Volver atrás reutiliza las películas de memoria');
 assert.deepEqual(grid.items.map(x=>x.id),firstIDs);
 assert.equal(el.labPrevPage.disabled,true);

 // Reproducción del error antiguo de 8 títulos por página TMDB.
 // Debe completar la primera página con 8+8+8, SIN enseñar las 24.
 context.labResetBrowsePages=undefined;
 context.labPage=1;context.labPages=1;context.labMoviePool=[];
 context.labMovieIDs.clear();context.labTMDBNextPage=1;context.labTMDBTotal=0;
 context.labMovieResponseCache.clear();
 calls=[];itemsPerTMDBPage=8;pages=500;
 await context.labBrowse();
 assert.deepEqual(calls,[1,2,3]);
 assert.equal(grid.items.length,18,'Mostrar tres filas, no 24 y otra fila con cuatro huecos');
 assert.equal(context.labMoviePool.length,24,'Mantener seis títulos para siguiente página');

 // Error temporal: no perder caché previa; si se inicia desde cero, dar aviso.
 context.labPage=1;context.labPages=1;context.labMoviePool=[];
 context.labMovieIDs.clear();context.labTMDBNextPage=1;context.labTMDBTotal=0;
 context.labMovieResponseCache.clear();calls=[];shouldError=true;
 await context.labBrowse();
 assert.equal(grid.items.length,0);
 assert.match(el.labDiscoverStatus.textContent,/conexión/i);
 assert.equal(el.labResultsCount.textContent,'Página 1');

 console.log('OK: Descubrir pagina por filas completas, guarda sobrantes, navega sin pérdidas y trata errores');
}
test().catch(err=>{console.error(err);process.exitCode=1;});
