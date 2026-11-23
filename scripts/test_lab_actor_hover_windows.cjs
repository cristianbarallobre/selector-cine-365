'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');

const html=fs.readFileSync('lab_discover_windows.html','utf8');
const js=html.split('<!-- LAB_DISCOVERY_SCRIPT -->')[1];
const backend=fs.readFileSync('lab_discover_windows.go','utf8');
const cache=fs.readFileSync('lab_discover_cache_windows.go','utf8');
new vm.Script('(()=>{\n'+js+'\n})();',{filename:'lab_discover_windows.js'});
assert.match(backend,/ProfilePath string/,'Leer foto de los créditos TMDB');
assert.match(cache,/"profile_path":photo/,'Conservar foto del actor en la caché');
assert.match(cache,/p\["profile_path"\]==nil/,'Actualizar fichas anteriores que carezcan de fotografía');
assert.match(html,/id="labActorPhotoTooltip" role="tooltip" hidden/);
assert.match(html,/#labActorPhotoTooltip\{position:fixed;z-index:150/,'Evitar que la ventana recorte la foto');
assert.match(html,/#labActorPhotoTooltip\[hidden\]\{display:none!important\}/);
assert.match(html,/@media\(hover:none\)/,'No mostrar tooltips en dispositivos táctiles');
assert.match(js,/actor\.addEventListener\('mouseenter',/);
assert.match(js,/actor\.addEventListener\('mouseleave',labHideActorPhoto\)/);
assert.match(js,/actor\.addEventListener\('focus',/);
assert.match(js,/actor\.addEventListener\('blur',labHideActorPhoto\)/);
assert.match(js,/actor\.addEventListener\('click',\(\)=>\{labHideActorPhoto\(\);void labOpenPerson\(personID\);\}\)/);

class FakeElement {
 constructor(tag){this.tagName=tag;this.children=[];this.hidden=true;this.className='';this.textContent='';this.isConnected=true;this.events={};this.style={};}
 replaceChildren(...items){this.children=[...items];}
 appendChild(item){this.children.push(item);}
 append(...items){this.children.push(...items);}
 addEventListener(type,fn){this.events[type]=fn;}
 setAttribute(key,val){this[key]=val;}
 getBoundingClientRect(){return this.tagName==='tooltip'?{width:174,height:250}:{left:500,right:610,top:300,bottom:320};}
 focus(){this.events.focus?.();}
}
const tooltip=new FakeElement('tooltip'),credits=new FakeElement('root');
const dom={labActorPhotoTooltip:tooltip,labDetailCast:credits};
const overlay=new FakeElement('overlay');overlay.hidden=false;
const timers=new Map();
let nextTimer=0,lastDelay=0;
const ctx={
 $:id=>{assert.ok(id in dom,id);return dom[id];},
 document:{createElement:tag=>new FakeElement(tag),documentElement:{clientWidth:1000,clientHeight:720}},
 labDetailOverlay:overlay,
 labPosterSrc:p=>'/api/discover/poster?path='+encodeURIComponent(p),
 labOpenPerson:()=>{},
 window:{innerWidth:1000,innerHeight:720,addEventListener:()=>{}},
 setTimeout:(fn,ms)=>{const id=++nextTimer;timers.set(id,fn);lastDelay=ms;return id;},
 clearTimeout:id=>timers.delete(id)
};
vm.createContext(ctx);
const start=js.indexOf('let labActorPhotoTimer=null,labActorPhotoTarget=null;');
const end=js.indexOf('// Clasificaciones locales:',start);
assert.ok(start>0&&end>start);
vm.runInContext(js.slice(start,end),ctx);
assert.equal(ctx.labActorPhotoPath({profile_path:'/safe_profile.jpg'}),'/safe_profile.jpg');
assert.equal(ctx.labActorPhotoPath({profile_path:'https://evil.example/a.jpg'}),'');
assert.equal(ctx.labActorPhotoPath({profile_path:'/bad?url=.jpg'}),'');
assert.equal(ctx.labActorPhotoPath({profile_path:null}),'');
assert.equal(tooltip.hidden,true);

ctx.labRenderCredits({director:'Directora',cast:[
 {id:42,name:'Intérprete uno',character:'Papel principal',profile_path:'/portrait_a.jpg'},
 {id:43,name:'Intérprete sin retrato',profile_path:''},
 {name:'Sin ID',profile_path:'/portrait_b.jpg'}
]});
const cards=credits.children[2].children;
const button=cards[0].children[0];
assert.equal(button.tagName,'button');
assert.equal(tooltip.children.length,0,'No cargar ninguna foto al dibujar la ficha');
button.events.mouseenter();
assert.equal(lastDelay,200,'Retrasar la carga para ignorar movimientos rápidos del ratón');
assert.equal(timers.size,1);
button.events.mouseleave();
assert.equal(timers.size,0,'Cancelar petición antes de mostrar tooltip');
assert.equal(tooltip.hidden,true);
button.events.mouseenter();
assert.equal(timers.size,1);
for(const run of timers.values())run();timers.clear();
assert.equal(tooltip.hidden,false);
assert.equal(tooltip.children[0].tagName,'img');
assert.match(tooltip.children[0].src,/portrait_a\.jpg/);
assert.equal(tooltip.children[1].textContent,'Intérprete uno');
assert.equal(tooltip.style.left,'500px');
assert.equal(tooltip.style.top,'41px');
button.events.mouseleave();
assert.equal(tooltip.hidden,true,'No dejar retratos visibles al retirar el cursor');
const noPortrait=cards[1].children[0];
noPortrait.events.mouseenter();
assert.equal(timers.size,0,'Sin foto TMDB no solicitar información adicional');
assert.equal(tooltip.hidden,true);
button.events.focus();
assert.equal(lastDelay,0,'Con teclado mostrar el retrato sin demora');
assert.equal(timers.size,1);
overlay.hidden=true;
for(const run of timers.values())run();timers.clear();
assert.equal(tooltip.hidden,true,'Si la ficha se cierra mientras carga, no mostrar el retrato');
console.log('OK: retratos al hover/foco, carga diferida, sin descargas de golpe, sin fotos falsas ni ventanas cortadas');
