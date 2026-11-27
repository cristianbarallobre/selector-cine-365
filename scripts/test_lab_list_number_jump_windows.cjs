'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');

const html=fs.readFileSync('lab_discover_windows.html','utf8');
const source=html.split('<!-- LAB_DISCOVERY_SCRIPT -->')[1];
const go=fs.readFileSync('lab_discover_desktop.go','utf8');
const reactions=fs.readFileSync('lab_reactions_desktop.go','utf8');
const stable=fs.readFileSync('plantilla.html','utf8');
new vm.Script('(()=>{\n'+source+'\n})();');
assert.match(html,/#view-library \.film-item\.lab-list-jump-target/);
assert.match(source,/function labJumpToFilmList\(m\)/);
assert.match(source,/if\(m\.id!==libraryState\.active_id\|\|virtualMode\|\|labFavoritesMode\)switchCollection\(m\.id\)/,'El número debe salir de Favoritos aunque apunte a la lista real ya activa');
assert.match(reactions,/const labWasFavorite=labFavoritesMode;labFavoritesMode=false;/,'El cambio de lista desactiva Favoritos');
assert.match(reactions,/!virtualMode&&id===libraryState.active_id&&!labWasFavorite/,'El cambio a la misma lista desde Favoritos no se omite');
assert.match(source,/search='';\$\('search'\)\.value='';filter='all'/);
assert.match(source,/labFiltersResetForList\(\)/);
assert.match(source,/filmPage=Math\.floor\(index\/capacity\)/);
assert.match(source,/setTimeout\(\(\)=>\{[\s\S]*?\},4200\)/);
assert.match(go,/function openFilmInList\(m\)\{labJumpToFilmList\(m\);\}/);
assert.match(go,/b\.classList\.add\('lab-list-jump-target'\)/);
assert.ok(stable.includes("selectSearchScope('current');search=m.title;$('search').value=search;"),'La plantilla estable queda intacta');

class Row{
 constructor(item){
  this.dataset={idx:String(item.idx)};
  this.styles=new Set();
  this.classList={
   add:key=>this.styles.add(key),
   remove:key=>this.styles.delete(key),
   contains:key=>this.styles.has(key)
  };
 }
 scrollIntoView(options){this.scrollOptions=options;}
}
const listA={id:'one',name:'Uno',items:Array.from({length:8},(_,i)=>({idx:i,title:'A '+i,year:'2024'}))};
const listB={id:'two',name:'Dos',items:Array.from({length:15},(_,i)=>({idx:i+30,title:'B '+i,year:'2026'}))};
const lists=[listA,listB];
const searchEl={value:'Antes'},filterButtons=['all','seen','pending'].map(filter=>({dataset:{filter},active:filter==='seen',classList:{toggle(key,enabled){assert.equal(key,'selected');this.selected=enabled;}}}));
const visible=[];
const filmEl={
 querySelectorAll(query){
  if(query==='.film-item')return visible;
  throw Error('Unexpected film query '+query);
 }
};
const doc={
 querySelectorAll(selector){
  if(selector==='.filter')return filterButtons;
  if(selector==='#filmList .lab-list-jump-target')return visible.filter(row=>row.classList.contains('lab-list-jump-target'));
  throw Error('Unexpected document query '+selector);
 }
};
const $=name=>({search:searchEl,filmList:filmEl})[name]||(()=>{throw Error('Unknown element '+name)})();
const pendingTimers=new Map(),frames=[];
let timerCounter=0,loaded=null,view='',switches=0,filterResets=0,renders=0;
const ctx={
 libraryState:{lists:[{id:'one',name:'Uno'},{id:'two',name:'Dos'}],active_id:'one'},
 parsed:{entries:listA.items},search:'Antes',filter:'seen',searchScope:'all',filmPage:1,virtualMode:false,labFavoritesMode:false,
 selectedListPlatforms:new Set(['Netflix']),
 labListJumpHighlight:null,labListJumpTimer:null,labListJumpSerial:0,
 document:doc,$,
 closeListPlatformMenu:()=>{},updateListPlatformOptions:()=>{},
 labFiltersResetForList:()=>{filterResets++;ctx.labDecade='any';ctx.labGenre='any';ctx.labDuration='any';ctx.labOrder='original';},
 labDecade:'1980s',labGenre:'Drama',labDuration:'long',labOrder:'rating',
 switchCollection:id=>{switches++;ctx.labFavoritesMode=false;ctx.libraryState.active_id=id;ctx.parsed={entries:lists.find(l=>l.id===id).items};ctx.filmPage=0;view='library';},
 showView:name=>{view=name;ctx.renderList();},
 selectSearchScope:scope=>{ctx.searchScope=scope;ctx.filmPage=0;ctx.renderList();},
 renderList:()=>{
  renders++;
  const list=ctx.parsed.entries;
  const source=list.filter(f=>(!ctx.search||f.title.toLowerCase().includes(ctx.search.toLowerCase()))&&(ctx.filter==='all'||f.checked));
  const start=ctx.filmPage*4;
  visible.splice(0,visible.length,...source.slice(start,start+4).map(movie=>{
   const row=new Row(movie);
   const marker=vm.runInContext('labListJumpHighlight',ctx);
   if(marker?.listID===ctx.libraryState.active_id&&marker.rawIndex===movie.idx)row.classList.add('lab-list-jump-target');
   return row;
  }));
 },
 requestAnimationFrame:fn=>{frames.push(fn);},
 setTimeout:(fn,time)=>{const id=++timerCounter;pendingTimers.set(id,{fn,time});return id;},
 clearTimeout:id=>pendingTimers.delete(id),
 Array,Math,String
};
vm.createContext(ctx);
const start=source.indexOf('let labListJumpHighlight=null,labListJumpTimer=null,labListJumpSerial=0;');
const end=source.indexOf('// Confirmación visual exclusiva de LAB;',start);
assert.ok(start>=0&&end>start,'Encontrar el módulo propio sin abrir elementos de otras pantallas');
vm.runInContext(source.slice(start,end),ctx);

ctx.labJumpToFilmList({id:'two',title:'B 12',year:'2026'});
assert.equal(ctx.libraryState.active_id,'two','Cambia a la lista correspondiente');
assert.equal(ctx.search,'','Vacía la búsqueda');
assert.equal(searchEl.value,'','El campo visual de búsqueda queda vacío');
assert.equal(ctx.searchScope,'current','Sale de la búsqueda global');
assert.equal(ctx.filter,'all','Muestra vistas y pendientes');
assert.equal(ctx.selectedListPlatforms.size,0,'Limpia las plataformas');
assert.equal(ctx.labDecade,'any');assert.equal(ctx.labDuration,'any');assert.equal(ctx.labGenre,'any');assert.equal(ctx.labOrder,'original');
assert.equal(ctx.parsed.entries.length,15,'No altera el contenido de la lista');
assert.equal(ctx.filmPage,3,'Busca en la cuarta página visual');
assert.equal(visible.length,3);
assert.equal(visible[0].dataset.idx,'42');
assert.equal(visible[0].classList.contains('lab-list-jump-target'),true,'Resalta la fila correcta');
assert.equal(pendingTimers.size,1);
assert.equal([...pendingTimers.values()][0].time,4200);
frames.shift()();
assert.equal(visible[0].scrollOptions.block,'center','Desplaza hasta la película');
const old=visible[0];
ctx.renderList();
assert.equal(visible[0].classList.contains('lab-list-jump-target'),true,'Resiste refrescos durante la precarga');
for(const timer of [...pendingTimers.values()])timer.fn();
pendingTimers.clear();
assert.equal(visible[0].classList.contains('lab-list-jump-target'),false,'El resaltado desaparece solo');
assert.equal(ctx.labListJumpHighlight,null);
assert.equal(old.classList.contains('lab-list-jump-target'),true,'La fila sustituida no afecta la selección nueva');

ctx.search='Otra';searchEl.value='Otra';ctx.filter='pending';
ctx.labJumpToFilmList({id:'two',title:'B 3',year:'2026'});
assert.equal(switches,1,'No vuelve a cambiar de colección si ya está en ella');
assert.equal(ctx.filmPage,0);
assert.equal(visible[3].classList.contains('lab-list-jump-target'),true);
ctx.labJumpToFilmList({id:'missing',title:'B 3',year:'2026'});
assert.equal(ctx.libraryState.active_id,'two','Ignora listas inexistentes');
assert.equal(ctx.search,'');
assert.ok(renders>=5);
assert.ok(filterResets>=2);

// Regresión real del error: Favoritos conserva «two» como lista activa.
// Pulsar su número 2 debe SALIR de Favoritos, aunque active_id ya sea 'two'.
const baselineSwitches=switches;
ctx.labFavoritesMode=true;
ctx.search='Película';searchEl.value='Película';ctx.filter='seen';
ctx.labJumpToFilmList({id:'two',title:'B 12',year:'2026'});
assert.equal(switches,baselineSwitches+1,'Salir de Favoritos obliga a cambiar la vista aunque el ID ya coincida');
assert.equal(ctx.labFavoritesMode,false,'Se desactiva la vista automática');
assert.equal(ctx.libraryState.active_id,'two');
assert.equal(ctx.filmPage,3,'Salta a la página correcta de la lista real');
assert.equal(visible[0].classList.contains('lab-list-jump-target'),true,'Resalta la película de la lista real');
assert.equal(searchEl.value,'','No debe dejar una búsqueda que oculte la película');
assert.equal(ctx.filter,'all');

// Cambiar a otra lista real desde Favoritos también debe funcionar.
ctx.labFavoritesMode=true;
ctx.labJumpToFilmList({id:'one',title:'A 7',year:'2024'});
assert.equal(ctx.labFavoritesMode,false);
assert.equal(ctx.libraryState.active_id,'one');
assert.equal(ctx.filmPage,1);
assert.equal(visible[3].dataset.idx,'7');
assert.equal(visible[3].classList.contains('lab-list-jump-target'),true);

// Entre listas normales la navegación mantiene su comportamiento previo.
const beforeNormal=switches;
ctx.labJumpToFilmList({id:'one',title:'A 0',year:'2024'});
assert.equal(switches,beforeNormal,'No cambia la lista cuando ya está en ella');
assert.equal(ctx.filmPage,0);
assert.equal(visible[0].classList.contains('lab-list-jump-target'),true);
console.log('OK: números de lista desde Favoritos (misma lista u otra), paginación y resaltado sin filtrar la película');
