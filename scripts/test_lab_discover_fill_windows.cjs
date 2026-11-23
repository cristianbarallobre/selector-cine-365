'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');

const html=fs.readFileSync('lab_discover_windows.html','utf8');
const script=html.split('<!-- LAB_DISCOVERY_SCRIPT -->')[1];
assert.ok(script);
new vm.Script('(()=>{\n'+script+'\n})();',{filename:'lab_discover_windows.js'});
assert.match(script,/async function labBrowse\(\)/);
assert.match(script,/function labDiscoverFillTarget\(\)/,'El número de resultados se deriva de la geometría');
assert.match(script,/while\(requested<8&&movies.length<target/,'Se limita el número de consultas externas por pantalla');
assert.match(script,/labPageStarts\[labPage\]=endPage\+1/,'No repetir en la siguiente página los resultados usados para rellenar');
assert.match(script,/labMovieResponseCache\.set\(cacheKey,data\)/,'Conservar las páginas ya descargadas en memoria');
assert.match(script,/window\.addEventListener\('resize',labDiscoverOnResize\)/,'Ajustar filas al tamaño de ventana');
assert.match(script,/labPageStarts=labPageStarts\.slice\(0,labPage\)/,'Redimensionar sin repetir páginas anteriores');
assert.match(script,/labResetBrowsePages\(\)/,'Reiniciar la agrupación al cambiar categoría o búsqueda');

const start=script.indexOf('function labDiscoverFillTarget(){');
const end=script.indexOf('function labCard(film){',start);
assert.ok(start>=0&&end>start);
const code=script.slice(start,end);

const grid={
 items:[],clientWidth:1200,
 get childElementCount(){return this.items.length;},
 replaceChildren(){this.items=[]},
 appendChild(card){this.items.push(card)},
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
 labNextPage:{disabled:false}
};
let calls=[],nextError=false;
const context={
 $:id=>{assert.ok(id in el,'Falta elemento '+id);return el[id];},
 window:{
  // Cuando auto-fill aún no tiene tarjetas, CSS devuelve 'none' en algunas versiones de Edge.
  getComputedStyle:()=>({gridTemplateColumns:'none',columnGap:'14px'}),
  addEventListener:()=>{}
 },
 labView:{getBoundingClientRect:()=>({bottom:900}),hidden:false},
 labStatus:s=>{el.labDiscoverStatus.textContent=s},
 labCard:x=>x,
 labKind:'upcoming',labQuery:'',labPage:1,labPages:1,
 labPageStarts:[1],labPageEnds:[],labMovieResponseCache:new Map(),
 labLastFillTarget:0,labResizeTimer:null,labRequestSerial:0,
 labResetBrowsePages:()=>{},
 labGet:async url=>{
  const match=url.match(/[?&]page=(\d+)/);
  assert.ok(match);
  const page=Number(match[1]);calls.push(page);
  if(nextError){nextError=false;throw new Error('La conexión falló');}
  return {
   page,total_pages:500,total_results:4000,
   results:Array.from({length:8},(_,i)=>({id:page*100+i,title:'Página '+page+' #'+i}))
  };
 },
 setTimeout:()=>{},clearTimeout:()=>{},
 activeView:'discover',Number,Math,Map,Set
};
vm.createContext(context);
vm.runInContext(code,context);
(async()=>{
 const expectedTarget=context.labDiscoverFillTarget();
 assert.equal(expectedTarget,12,'Seis columnas requieren doce tarjetas para completar dos filas');
 assert.match(script,/availableWidth=Math\.max\(180,grid\.getBoundingClientRect\(\)\.width\|\|grid\.clientWidth\)/,
  'El número de columnas debe calcularse por ancho, sin depender de grid-template-columns none');
 await context.labBrowse();
 assert.deepEqual(calls,[1,2],'Cargar una segunda página TMDB si la primera solo devuelve 8 películas');
 assert.equal(grid.items.length,16,'Rellenar el espacio con resultados adicionales');
 assert.equal(context.labPageStarts[1],3,'La siguiente página visual empieza en TMDB 3, sin repetir TMDB 2');
 assert.equal(el.labNextPage.disabled,false,'La navegación debe permanecer disponible');
 assert.match(el.labPageLabel.textContent,/TMDB 1–2 \/ 500/);
 const ids1=grid.items.map(x=>x.id);
 context.labPage=2;await context.labBrowse();
 assert.deepEqual(calls,[1,2,3,4]);
 assert.equal(grid.items.length,16);
 assert.ok(grid.items.every(x=>!ids1.includes(x.id)),'La página 2 no contiene películas de la página 1');
 assert.equal(context.labPageStarts[2],5);
 context.labPage=1;await context.labBrowse();
 assert.deepEqual(calls,[1,2,3,4],'Volver a una página anterior reutiliza caché de sesión');
 assert.deepEqual(grid.items.map(x=>x.id),ids1);
 assert.equal(el.labPrevPage.disabled,true);

 const labelBefore=el.labPageLabel.textContent;
 nextError=true;context.labMovieResponseCache.clear();
 await context.labBrowse();
 assert.match(el.labDiscoverStatus.textContent,/No se encontraron|conexión/i,'Mostrar error de red sin promesa rechazada');
 assert.equal(el.labResultsCount.textContent,'Página 1');
 assert.equal(grid.items.length,0,'Ante error completo no deben quedar tarjetas anteriores');
 assert.ok(labelBefore.includes('TMDB'));
 console.log('OK: Descubrir rellena 8+8 tarjetas, navega sin repetición, reutiliza caché y maneja errores de red');
})().catch(err=>{console.error(err);process.exitCode=1});
