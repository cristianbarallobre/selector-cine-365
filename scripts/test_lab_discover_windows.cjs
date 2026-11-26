'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');
const ui=fs.readFileSync('lab_discover_windows.html','utf8');
const windows=fs.readFileSync('lab_discover_desktop.go','utf8');
const other=fs.readFileSync('lab_discover_other.go','utf8');
const core=fs.readFileSync('main.go','utf8');
const parts=ui.split('<!-- LAB_DISCOVERY_SCRIPT -->');
assert.equal(parts.length,2,'La pantalla y su JS deben inyectarse por separado');
new vm.Script('(()=>{\n'+parts[1]+'\n})();',{filename:'lab_discover_windows.js'});
for(const id of ['labDiscoverView','labDiscoverNav','labSearchForm','labMoviesGrid','labDetailOverlay','labDetailAdd','labAddList'])
 assert.ok((ui+windows).includes(id),'Falta '+id);
assert.match(windows,/\/\/go:build windows/);
assert.match(other,/\/\/go:build !windows/);
assert.match(core,/out = injectLabDiscovery\(out\)/);
assert.doesNotMatch(windows,/if !appIsLab/,'También disponible en Windows estable');
assert.match(ui,/syncActiveText\(\)/,'Debe guardar el texto activo antes de añadir');
assert.match(ui,/if\(existing\)/,'Debe evitar duplicados');
assert.match(ui,/scheduleSave\(\)/,'Debe guardar mediante el flujo normal');
assert.match(ui,/labDetailOverlay\.hidden=true/,'Debe cerrar la ficha');
assert.match(ui,/#labDetailPlaceholder\[hidden\]\s*\{\s*display\s*:\s*none\s*!important\s*\}/,'El marcador de carátula debe desaparecer al activar hidden');
assert.match(ui,/^\s*const labNav=/m,'Debe compartir el estado del script existente');
assert.match(windows,/b\.append\(togg,pill,labCreateDetailButton\(f\),fix\);frag\.appendChild\(b\);/,'La biblioteca muestra Pendiente/Vista, Detalles y Lápiz en ese orden');
assert.match(ui,/function labCreateDetailButton\(f\)/,'Debe existir el botón de detalles en la Biblioteca');
assert.match(ui,/button\.textContent='Detalles'/,'El botón de Biblioteca ha de indicar Detalles');
assert.match(ui,/details\.textContent='Detalles'/,'Cada tarjeta de Descubrir ha de incluir Detalles');
assert.match(ui,/details\.addEventListener\('click'/,'El botón Detalles debe abrir fichas en Descubrir');
assert.doesNotMatch(ui,/card\.addEventListener\('click'/,'La tarjeta completa de Descubrir no debe ser clicable');
assert.doesNotMatch(ui,/\$\('filmList'\)\.addEventListener\('click'/,'Clicar la fila de Biblioteca no abre la ficha');
assert.match(ui,/#labDetailPlaceholder\[hidden\]\{display:none!important\}/,'El cartel de reserva debe ocultarse si se cargó la carátula');
assert.ok(fs.readFileSync('plantilla.html','utf8').includes('titleLine.appendChild(title);'),'El canal oficial mantiene el código original');
assert.ok(!fs.readFileSync('android/app/src/main/assets/index.html','utf8').includes('labDiscoverView'),'Android no debe heredar Descubrir');

const cache=fs.readFileSync('lab_discover_cache_desktop.go','utf8');
assert.match(cache,/\/\/go:build windows/,'La caché sólo debe compilarse en Windows');
assert.match(cache,/tmdb_fichas_lab\.json/,'La caché debe estar en un archivo LAB independiente');
assert.match(cache,/labMetaTTL\s*=\s*30\*24\*time\.Hour/,'Metadatos refrescados cada 30 días');
assert.match(cache,/labExtraTTL\s*=\s*7\*24\*time\.Hour/,'Notas y tráiler refrescados cada 7 días');
assert.match(cache,/labProvidersTTL\s*=\s*24\*time\.Hour/,'Plataformas refrescadas cada 24 horas');
assert.match(cache,/labPrefetchDelay\s*=\s*900\*time\.Millisecond/,'Descargas limitadas para respetar TMDB');
assert.match(cache,/func \(c \*labCatalogCache\) runWorker/,'Debe descargar fichas de forma secuencial');
assert.match(cache,/func \(c \*labCatalogCache\) fileSaveLocked/,'Debe persistir sin cambiar las listas');
assert.match(cache,/if prior,ok:=c\.record\(id\);ok/,'La apertura debe consultar primero la caché');
assert.match(cache,/if labNeedsRefresh\(prior\)/,'La actualización obsoleta debe ser diferida');
assert.match(ui,/function labCollectFilms/,'Precargar biblioteca completa');
assert.match(ui,/for\(const list of ordered\)/,'Recorrer todas las listas');
assert.match(ui,/function labQueuePrefetchSoon/,'Programar precarga no bloqueante');
assert.match(windows,/id="labPrefetchPause"/,'Debe haber botón de pausa en la UI Windows inyectada');
assert.match(windows,/b\.append\(togg,pill,labCreateDetailButton\(f\),fix\);frag\.appendChild\(b\);/,'Fila sin ↗ y con orden correcto');
const template=fs.readFileSync('plantilla.html','utf8');
assert.match(template,/await initPlatforms\(\);\s*void checkCloudIncoming\(\);\s*\/\* CINE_DESKTOP_READY \*\//,'El arranque solo comprueba la nube y ofrece el punto de precarga');
assert.match(windows,/"\/\* CINE_DESKTOP_READY \*\/", "labQueuePrefetchSoon\(\);"/,'Lanzar la precarga en el punto de arranque');
assert.match(windows,/function scheduleSave\(\)\{labQueuePrefetchSoon\(\)/,'Preparar películas agregadas tras guardar');
assert.match(core,/\/api\/discover\/prefetch/,'Ruta autenticada para precarga');

assert.match(ui,/\.lab-discover\.view\{[^}]*overflow-y:auto/,'Descubrir debe tener scroll vertical propio');
assert.match(ui,/\.lab-discover \.panel\{height:auto/,'Los paneles no pueden ocupar cada uno el 100% de la altura y recortar las tarjetas');
assert.match(ui,/\.lab-discover \[hidden\]\{display:none!important\}/,'Respetar elementos ocultos de la paginación');
assert.doesNotMatch(ui,/id="labTopPager"|labTopNext|labTopPrev/,'La paginación superior debe haberse eliminado');
assert.match(ui,/id="labDiscoverPager"/,'La paginación inferior debe mantenerse');
assert.match(ui,/Math\.min\(500,Number\(data\.total_pages\)/,'Mostrar más allá de las cinco primeras páginas');
assert.match(ui,/labTMDBTotal\.toLocaleString\('es-ES'\)/,'Mostrar número total de resultados');
assert.match(ui,/function labGoPage\(direction\)/,'La paginación inferior debe navegar');
assert.match(windows,/n>500/,'El servidor acepta páginas más allá de 5');
assert.match(windows,/TotalResults int/,'El servidor incluye el total de resultados de TMDB');
assert.match(ui,/margin:7px 13px 7px 5px/,'Detalles y lápiz deben tener espacio visible');
assert.match(ui,/addEventListener\('click',labLeave,true\)/,'Descubrir debe ocultarse en fase capture antes del renderizado de Biblioteca');
assert.match(windows,/function showView\(name\)\{labLeave\(\);if\(name==='library'\)/,'También debe ocultarse Descubrir en navegación programática entre listas y pestañas');
assert.match(ui,/ResizeObserver/,'La lista debe recuperar la paginación tras cambios de tamaño');
assert.match(ui,/if\(height<160\|\|Math\.abs\(height-labLastListHeight\)<4\)return/,'No debe recalcularse con altura colapsada');
assert.match(ui,/details\.addEventListener\('click'/,'Solo Detalles abrirá la ficha desde Descubrir');
assert.doesNotMatch(ui,/card\.tabIndex|card\.addEventListener\('keydown'/,'Las tarjetas completas no son objetivos de teclado');



// La altura nueva solo existe en LAB Windows: plantilla y builds estables quedan intactos.
assert.match(ui,/#view-library \.film-item\{\s*height:var\(--lab-film-row-height,96px\)/,'Las filas LAB deben usar altura adaptable');
assert.match(ui,/#view-library \.film-toggle\{[^}]*align-items:center/,'Los controles de fila deben quedar centrados verticalmente');
assert.match(ui,/#view-library \.film-poster,#view-library \.poster-placeholder\{width:42px;height:60px;align-self:center\}/,'Carátulas LAB centradas');
assert.match(windows,/Math\.floor\(\(labUsableHeight\+7\)\/103\)/,'La paginación debe conservar la capacidad anterior');
assert.match(windows,/Math\.min\(110,Math\.max\(96,Math\.floor/,'Las filas deben crecer solo donde sobre espacio');
assert.match(windows,/parent\.style\.setProperty\('--lab-film-row-height',labRowHeight\+'px'\)/,'La altura adaptativa debe aplicarse en cada renderizado');
for(const areaHeight of [275,320,370,405,460,488,520,590,640,720]){
 const labUsableHeight=Math.max(0,areaHeight-2);
 const pageSize=Math.max(1,Math.floor((labUsableHeight+7)/103));
 const rowHeight=Math.min(110,Math.max(96,Math.floor((labUsableHeight-7*(pageSize-1))/pageSize)));
 assert.ok(rowHeight>=96 && rowHeight<=110,'La altura de fila debe quedar entre 96 y 110px');
 assert.ok(pageSize*rowHeight+7*(pageSize-1)<=labUsableHeight,'Las filas no pueden invadir el paginador');
 assert.equal(pageSize,Math.max(1,Math.floor((labUsableHeight+7)/103)),'No perder capacidad por crecer las filas');
}
assert.ok(!fs.readFileSync('plantilla.html','utf8').includes('--lab-film-row-height'),'No alterar la plantilla común');


// Ficha de película: mostrar pertenencias reales antes de intentar añadir.
// Comprobar coincidencias por ID TMDB, título/año, listas múltiples y Favoritos.
for(const id of ['labDetailMemberships','labAddList','labDetailAdd'])assert.match(ui,new RegExp('id="'+id+'"'));
assert.match(ui,/EN TUS LISTAS/);
assert.match(ui,/function labMovieMemberships\(movie\)/);
assert.match(ui,/function labRenderDetailMemberships\(movie\)/);
assert.match(ui,/function labUpdateDetailAddState\(\)/);
assert.match(ui,/labRenderDetailMemberships\(movie\);labUpdateDetailAddState\(\);/,'Actualizar inmediatamente tras añadir');
assert.match(ui,/\$\('labAddList'\)\.addEventListener\('change',labUpdateDetailAddState\)/,'Reaccionar al cambio de la lista de destino');
assert.match(ui,/labRenderDetailMemberships\(null\);/,'Vaciar pertenencias durante la carga');
assert.match(ui,/labRenderDetailMemberships\(m\);/,'Mostrar pertenencias en una ficha nueva');
assert.match(ui,/\.lab-detail-membership-chips\{[^}]*flex-wrap:wrap/,'Listas con nombres largos deben poder saltar de línea');

class MembershipElement{
 constructor(tag='div'){
  this.tagName=tag;this.children=[];this.textContent='';this.title='';
  this.disabled=false;this.value='';this.className='';this.events={};
 }
 replaceChildren(...els){this.children=[...els];}
 append(...els){this.children.push(...els);}
 appendChild(el){this.children.push(el);}
 get childElementCount(){return this.children.length;}
 addEventListener(type,fn){this.events[type]=fn;}
}
const holders={
 labDetailMemberships:new MembershipElement('div'),
 labAddList:new MembershipElement('select'),
 labDetailAdd:new MembershipElement('button')
};
const movie={id:77,title:'Vengadores: Doomsday',original_title:'Avengers: Doomsday',release_date:'2026-12-18'};
const memberships=[
 {id:'listB',n:2,name:'Cine de acción',title:'Avengers: Doomsday',year:'2026'},
 {id:'listA',n:1,name:'Quiero ver',title:'Vengadores: Doomsday',year:'2026'},
 {id:'listX',n:4,name:'Estrenos de 2025',title:'Vengadores: Doomsday',year:'2025'},
 {id:'listC',n:3,name:'Especial de superhéroes',title:'Título alternativo',year:'2026'}
];
const knownIDs=new Map([['titulo alternativo|2026',77],['otra película|2026',88]]);
const membershipCtx={
 $:id=>{assert.ok(holders[id],id);return holders[id];},
 document:{createElement:tag=>new MembershipElement(tag)},
 labFold:s=>String(s||'').normalize('NFD').replace(/[\u0300-\u036f]/g,'').toLowerCase().trim(),
 knownTMDB:f=>knownIDs.get(String(f.title).normalize('NFD').replace(/[\u0300-\u036f]/g,'').toLowerCase()+'|'+f.year)||0,
 filmKey:f=>String(f.title).toLowerCase()+'|'+f.year,
 filmMembershipIndex:()=>({rows:memberships.map(m=>({_lists:[m]}))}),
 libraryState:{tmdb_overrides:{}},
 labReactionFor:()=>null,
 labCurrentMovie:movie,
 Number,String,Set,Map
};
vm.createContext(membershipCtx);
const js=parts[1],start=js.indexOf('function labMovieMemberships(movie){'),end=js.indexOf('async function labShowDetail(',start);
assert.ok(start>0&&end>start,'Funciones de pertenencias deben estar disponibles antes de la ficha');
vm.runInContext(js.slice(start,end),membershipCtx);
const matches=membershipCtx.labMovieMemberships(movie);
assert.equal(matches.length,3,'Dos nombres de lista y alias con ID TMDB, sin falsos positivos por año');
assert.deepEqual(Array.from(matches,x=>x.id),['listA','listB','listC'],'Orden de listas conservado');
membershipCtx.labRenderDetailMemberships(movie);
const chips=holders.labDetailMemberships.children;
assert.equal(chips.length,3);
assert.equal(chips[0].children[0].textContent,'1 ·');
assert.equal(chips[0].children[1].textContent,'Quiero ver');
assert.equal(chips[1].children[1].textContent,'Cine de acción');
assert.equal(chips[2].children[1].textContent,'Especial de superhéroes');
holders.labAddList.value='listA';
membershipCtx.labUpdateDetailAddState();
assert.equal(holders.labDetailAdd.disabled,true,'No permitir añadir de nuevo a la lista donde está');
assert.match(holders.labDetailAdd.textContent,/Ya está/);
holders.labAddList.value='listNew';
holders.labAddList.events.change();
assert.equal(holders.labDetailAdd.disabled,false,'Se permite seleccionar una lista diferente');
assert.match(holders.labDetailAdd.textContent,/Añadir pendiente/);
membershipCtx.labReactionFor=()=>({favorite:true});
membershipCtx.labRenderDetailMemberships(movie);
assert.equal(holders.labDetailMemberships.children.length,4,'Favoritos se añade como vista virtual sin número');
assert.equal(holders.labDetailMemberships.children.at(-1).textContent,'♥ Favoritos');
membershipCtx.labReactionFor=()=>null;
membershipCtx.filmMembershipIndex=()=>({rows:[]});
membershipCtx.labRenderDetailMemberships(movie);
assert.equal(holders.labDetailMemberships.children.length,1);
assert.match(holders.labDetailMemberships.children[0].textContent,/Todavía no/);
membershipCtx.labRenderDetailMemberships(null);
assert.equal(holders.labDetailMemberships.children.length,1);
assert.match(holders.labDetailMemberships.children[0].textContent,/Comprobando/);
console.log('OK: ficha muestra pertenencias y Favoritos, detecta alias TMDB, bloquea duplicados y actualiza lista elegida');

console.log('OK: sintaxis, aislamiento Windows y seguridad básica de Descubrir LAB');
