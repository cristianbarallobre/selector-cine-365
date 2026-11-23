'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');

const ui=fs.readFileSync('lab_discover_windows.html','utf8');
const backend=fs.readFileSync('lab_discover_windows.go','utf8');
const cache=fs.readFileSync('lab_discover_cache_windows.go','utf8');
const main=fs.readFileSync('main.go','utf8');
const js=ui.split('<!-- LAB_DISCOVERY_SCRIPT -->')[1];
assert.ok(js,'El módulo LAB debe mantener el separador');
new vm.Script('(()=>{\n'+js+'\n})();',{filename:'lab_discover_windows.html'});
assert.match(main,/\/api\/discover\/person/,'La ruta de intérprete está registrada');
assert.match(backend,/func \(a \*runningApp\) labPersonDetail/);
assert.match(backend,/append_to_response":{"movie_credits"}/,'La filmografía procede de TMDB');
assert.match(cache,/"id":p.ID/,'Las fichas de película conservan ID del reparto');
assert.match(cache,/labCastWithoutIDs/,'Las fichas antiguas actualizan créditos sin borrar la caché');
assert.match(ui,/id="labPersonOverlay"/);
assert.ok(!ui.includes('id="labDiscoverDuration"'));

class FakeElement {
 constructor(tag='div'){this.tagName=tag;this.className='';this.textContent='';this.children=[];this.hidden=false;this.events={};this.isConnected=true;}
 get childElementCount(){return this.children.length;}
 replaceChildren(...children){this.children=[...children];this.textContent='';}
 append(...children){this.children.push(...children);}
 appendChild(child){this.children.push(child);}
 addEventListener(name,callback){this.events[name]=callback;}
 setAttribute(name,value){this[name]=value;}
 removeAttribute(name){delete this[name];}
 focus(){this.focused=true;}
 click(){this.events.click?.();}
}
const credits=new FakeElement();
let actorID=0;
const ctx={
 $:id=>{assert.equal(id,'labDetailCast');return credits;},
 document:{createElement:tag=>new FakeElement(tag)},
 labOpenPerson:async id=>{actorID=id;}
};
vm.createContext(ctx);
const render=js.slice(js.indexOf('function labRenderCredits(movie){'),js.indexOf('async function labOpenPerson(id){'));
assert.ok(render.startsWith('function labRenderCredits(movie){'));
vm.runInContext(render,ctx);
ctx.labRenderCredits({director:'Director de prueba',cast:[
 {id:1234,name:'Actriz de prueba',character:'Protagonista'},
 {name:'Intérprete antiguo',character:'Personaje sin ID'}
]});
const cards=credits.children[2].children;
assert.equal(cards[0].children[0].tagName,'button','Actores con ID son botones accesibles');
assert.equal(cards[0].children[0].type,'button');
assert.equal(cards[0].children[0]['aria-label'],'Ver ficha de Actriz de prueba');
cards[0].children[0].click();
assert.equal(actorID,1234);
assert.equal(cards[1].children[0].tagName,'strong','No abrir identidades dudosas sin ID');

const elements={};
const get=id=>elements[id]??(elements[id]=new FakeElement());
const overlay=new FakeElement();
let request='',selectedMovie=0,closeCount=0;
const person={
 id:1234,name:'Actriz de prueba',biography:'Una biografía de TMDB.',
 birthday:'1988-01-01',deathday:'',place_of_birth:'Ciudad de prueba',
 known_for_department:'Acting',profile_path:'/persona.jpg',
 movies:[{id:77,title:'Película A',release_date:'2026-10-15',character:'Protagonista'},
         {id:88,title:'Película B',release_date:'2025-01-01',character:'Amiga'}],
 tmdb_url:'https://www.themoviedb.org/person/1234'
};
const detailCtx={
 $:get,document:{activeElement:new FakeElement('button'),createElement:tag=>new FakeElement(tag)},
 labPersonSerial:0,labPersonReturnFocus:null,labPersonOverlay:overlay,
 labPosterSrc:path=>'/api/discover/poster?path='+encodeURIComponent(path),
 labGet:async url=>{request=url;return person;},
 labClosePerson:()=>{closeCount++;overlay.hidden=true;},
 labShowDetail:id=>{selectedMovie=id;}
};
vm.createContext(detailCtx);
const open=js.slice(js.indexOf('async function labOpenPerson(id){'),js.indexOf('async function labShowDetail(id){'));
assert.ok(open.startsWith('async function labOpenPerson(id){'));
vm.runInContext(open,detailCtx);

(async()=>{
 await detailCtx.labOpenPerson(1234);
 assert.equal(request,'/api/discover/person?id=1234','La aplicación consulta su backend TMDB');
 assert.equal(overlay.hidden,false);
 assert.equal(get('labPersonName').textContent,'Actriz de prueba');
 assert.equal(get('labPersonBio').textContent,'Una biografía de TMDB.');
 assert.ok(get('labPersonPortrait').src.includes('persona.jpg'));
 assert.equal(get('labPersonFilms').children.length,2);
 assert.ok(get('labPersonFilms').children[0].textContent==='');
 assert.equal(get('labPersonFilms').children[0].children[0].textContent,'Película A (2026)');
 assert.equal(get('labPersonLinks').children[0].href,'https://www.themoviedb.org/person/1234');
 get('labPersonFilms').children[0].click();
 assert.equal(closeCount,1,'Salir de la ficha del actor al elegir película');
 assert.equal(selectedMovie,77,'La filmografía abre la ficha de la película');
 console.log('OK: reparto accesible por ID, ficha personal TMDB, filmografía navegable y duración ausente');
})().catch(err=>{console.error(err);process.exitCode=1;});
