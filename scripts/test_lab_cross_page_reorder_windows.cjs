'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');
const html=fs.readFileSync('lab_discover_windows.html','utf8');
const js=html.split('<!-- LAB_DISCOVERY_SCRIPT -->')[1];
const injection=fs.readFileSync('lab_discover_windows.go','utf8');
const base=fs.readFileSync('plantilla.html','utf8');
new vm.Script('(()=>{\n'+js+'\n})();',{filename:'LAB HTML'});
assert.match(injection,/handle\.draggable=false/);
assert.match(base,/handle\.draggable=true/,'La versión estable conserva arrastre nativo');
assert.match(js,/document\.body\.setPointerCapture/);
assert.match(js,/document\.elementFromPoint/);
assert.match(js,/filmPage=page;\s*state\.to=null;state\.hoverPage=null;state\.lastNavigatedButton=button;\s*renderList\(\)/);
assert.match(js,/moveFilmInCurrentList\(state\.from,target\)/);
assert.match(js,/setTimeout\(\(\)=>\{/);
assert.match(html,/#filmPaging \.lab-crosspage-nav-target/);

class Cl {
 constructor(){this.names=new Set();}
 add(v){this.names.add(v);}
 remove(v){this.names.delete(v);}
 contains(v){return this.names.has(v);}
}
class El {
 constructor(id='',idx=null){this.id=id;this.dataset=idx===null?{}:{idx:String(idx)};this.classList=new Cl();this.style={};this.children=[];this.disabled=false;this.textContent='';this.clientHeight=0;this.events={};this.parent=null;}
 append(...els){this.children.push(...els);}
 appendChild(el){this.children.push(el);el.parent=this;}
 remove(){this.removed=true;}
 setAttribute(){}
 closest(selector){
  if(selector==='#filmList .film-reorder-handle')return this.id==='grip'?this:null;
  if(selector==='.film-item')return this.id==='grip'?this.parent:null;
  if(selector==='#filmList .film-item')return this.id.startsWith('row')?this:null;
  if(selector==='#filmPaging button')return this.id==='filmNext'||this.id==='filmPrev'||this.id==='number'?this:null;
  if(selector==='#filmPageNumbers')return this.id==='number'?this:null;
  return null;
 }
}
const entries=Array.from({length:13},(_,i)=>({idx:i,title:'Película '+i}));
const visible=[];
const filmList=new El('filmList');filmList.clientHeight=412;
const next=new El('filmNext'),prev=new El('filmPrev'),number=new El('number');
number.classList.add('page-number');number.textContent='4';
const pageTotal=new El('filmPageTotal');pageTotal.textContent='/ 4';
const filmPaging=new El('filmPaging');
const body=new El('body');body.capture=null;
body.setPointerCapture=id=>{body.capture=id;};
body.hasPointerCapture=id=>body.capture===id;
body.releasePointerCapture=id=>{if(body.capture===id)body.capture=null;};
const docEvents={},winEvents={},timers=new Map();
let nextTimer=0,view='row',offset=0,enabled=true,renderCount=0,saveCount=0;
const dom={filmList,filmPaging,filmPageTotal:pageTotal};
const doc={
 body,documentElement:{clientWidth:1024,clientHeight:720},
 createElement:tag=>new El(tag),
 querySelectorAll:sel=>{
  if(sel==='#filmList .lab-crosspage-row-target')return visible.filter(el=>el.classList.contains('lab-crosspage-row-target'));
  if(sel==='#filmPaging .lab-crosspage-nav-target')return [next,prev,number].filter(el=>el.classList.contains('lab-crosspage-nav-target'));
  return [];
 },
 elementFromPoint:(_x,y)=>{
  if(y>=500)return view==='number'?number:next;
  return visible[offset]||null;
 },
 addEventListener:(name,cb)=>{docEvents[name]=cb;}
};
const win={innerWidth:1024,innerHeight:720,addEventListener:(name,cb)=>{winEvents[name]=cb;}};
const ctx={
 document:doc,window:win,$:id=>{assert.ok(dom[id],id);return dom[id]},
 libraryState:{active_id:'list-1'},filmPage:0,parsed:{entries},
 filmOrderEditable:()=>enabled,
 renderList:()=>{
  renderCount++;visible.splice(0,visible.length,...entries.slice(ctx.filmPage*4,ctx.filmPage*4+4).map(e=>new El('row-'+e.idx,e.idx)));
  next.disabled=ctx.filmPage>=3;prev.disabled=ctx.filmPage<=0;
 },
 moveFilmInCurrentList:(from,to)=>{saveCount++;ctx.moves.push({from,to,page:ctx.filmPage});ctx.renderList();return true;},
 moves:[],
 setTimeout:(cb,ms)=>{const id=++nextTimer;timers.set(id,{cb,ms});return id;},
 clearTimeout:id=>timers.delete(id)
};
vm.createContext(ctx);
const start=js.indexOf('let labCrossPageDrag=null;');
const end=js.indexOf('// Abrir una lista desde su número:',start);
assert.ok(start>=0&&end>start,'No se encontró el módulo de arrastre de LAB');
vm.runInContext(js.slice(start,end),ctx);
ctx.renderList();
let row=visible[1],grip=new El('grip');grip.parent=row;
docEvents.pointerdown({button:0,target:grip,pointerId:7,clientX:30,clientY:100,preventDefault(){}});
assert.equal(body.capture,7,'El origen se conserva en BODY al cambiar de página');
assert.equal(vm.runInContext('labCrossPageDrag.from',ctx),1);
const move=(x,y)=>docEvents.pointermove({pointerId:7,clientX:x,clientY:y,cancelable:true,preventDefault(){}});
move(39,115);
assert.equal(body.classList.contains('lab-crosspage-dragging'),true);
move(400,550);
assert.equal(timers.size,1);
assert.equal([...timers.values()][0].ms,460,'Cambio de página retardado intencionadamente');
for(const {cb} of [...timers.values()])cb();timers.clear();
assert.equal(ctx.filmPage,1);
assert.ok(renderCount>=2,'Ha cambiado la página sin perder el gesto');
offset=2;move(100,120);
assert.equal(visible[2].classList.contains('lab-crosspage-row-target'),true);
docEvents.pointerup({pointerId:7,clientX:100,clientY:120});
assert.deepEqual(ctx.moves,[{from:1,to:6,page:1}],'Reordenar desde página 1 hasta página 2');
assert.equal(body.capture,null,'Se libera captura al soltar');
assert.equal(body.classList.contains('lab-crosspage-dragging'),false);
assert.equal(saveCount,1);

// Arrastrar al número 4 y soltar directamente sobre el botón.
ctx.filmPage=0;ctx.renderList();offset=0;view='number';
row=visible[0];grip=new El('grip');grip.parent=row;
docEvents.pointerdown({button:0,target:grip,pointerId:8,clientX:30,clientY:100,preventDefault(){}});
docEvents.pointermove({pointerId:8,clientX:300,clientY:540,cancelable:true,preventDefault(){}});
docEvents.pointerup({pointerId:8,clientX:300,clientY:540});
assert.deepEqual(ctx.moves[1],{from:0,to:12,page:3},'Soltar sobre número coloca la película al inicio de esa página');
assert.equal(saveCount,2);

ctx.filmPage=0;ctx.renderList();view='row';enabled=true;
row=visible[0];grip=new El('grip');grip.parent=row;
docEvents.pointerdown({button:0,target:grip,pointerId:9,clientX:30,clientY:100,preventDefault(){}});
docEvents.pointermove({pointerId:9,clientX:60,clientY:120,cancelable:true,preventDefault(){}});
winEvents.keydown({key:'Escape',preventDefault(){},stopImmediatePropagation(){}});
assert.equal(body.capture,null);
assert.equal(saveCount,2,'Escape cancela sin modificar la lista');
enabled=false;
docEvents.pointerdown({button:0,target:grip,pointerId:10,clientX:30,clientY:100,preventDefault(){}});
assert.equal(body.capture,null,'No iniciar arrastres cuando hay filtros activos');

// Si la página cambió al mantener el cursor sobre «Siguiente»,
// soltar en ese mismo botón debe llevar al INICIO de la página mostrada,
// no producir otro salto accidental.
enabled=true;ctx.filmPage=0;ctx.renderList();view='row';
row=visible[2];grip=new El('grip');grip.parent=row;
docEvents.pointerdown({button:0,target:grip,pointerId:11,clientX:30,clientY:100,preventDefault(){}});
view='next';
docEvents.pointermove({pointerId:11,clientX:400,clientY:550,cancelable:true,preventDefault(){}});
assert.equal(timers.size,1);
for(const {cb} of [...timers.values()])cb();timers.clear();
assert.equal(ctx.filmPage,1);
docEvents.pointerup({pointerId:11,clientX:400,clientY:550});
assert.deepEqual(ctx.moves[2],{from:2,to:4,page:1},'No saltar una página extra al soltar sobre Siguiente');
assert.equal(saveCount,3);
console.log('OK: arrastre entre páginas, soltar en números o Siguiente, Escape y límites de filtros');
