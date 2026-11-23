'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');

const lab=fs.readFileSync('lab_discover_windows.html','utf8');
const base=fs.readFileSync('plantilla.html','utf8');
const backend=fs.readFileSync('lab_discover_cache_windows.go','utf8');
const separator='<!-- LAB_DISCOVERY_SCRIPT -->';
const [style,js]=lab.split(separator);
assert.ok(style&&js);
new vm.Script('(()=>{\n'+js+'\n})();',{filename:'lab_discover_windows.html'});
assert.match(style,/<div class="lab-detail-cast" id="labDetailCast"><\/div>/,'Los créditos son un contenedor estructurado');
assert.match(style,/\.lab-cast-actor\{[^}]*font-weight:770/,'Diferenciar nombres de intérpretes');
assert.match(style,/\.lab-cast-role\{[^}]*color:#d4bb94/,'Diferenciar personajes por color y peso');
assert.match(style,/\.lab-cast-role::before\{content:"Personaje: "/);
assert.match(backend,/"name":p.Name,"character":p.Character/,'Conservar actores y personajes TMDB por separado');
assert.match(js,/labRenderCredits\(m\)/);
assert.ok(!js.includes("names.join(' · ')"),'El reparto no debe convertirse en una única línea ambigua');

const start=js.indexOf('function labRenderCredits(movie){');
const end=js.indexOf('// Clasificaciones locales:',start);
assert.ok(start>=0&&end>start,'Falta la creación de fichas estructuradas');
const fn=js.slice(start,end);
class FakeElement{
 constructor(tag){this.tagName=tag;this.className='';this.textContent='';this.children=[];}
 replaceChildren(...nodes){this.children=[...nodes];}
 append(...nodes){this.children.push(...nodes);}
 appendChild(node){this.children.push(node);}
}
const target=new FakeElement('div');
const ctx={
 $:id=>{assert.equal(id,'labDetailCast');return target},
 document:{createElement:tag=>new FakeElement(tag)},
 labHideActorPhoto:()=>{}
};
vm.createContext(ctx);vm.runInContext(fn,ctx);
ctx.labRenderCredits({director:'Kitao Sakurai',
 cast:[{name:'Andrew Koji',character:'Ryu'},{name:'Noah Centineo',character:'Ken Masters'},
       {name:'Actor sin personaje',character:''},{name:'',character:'No se conoce'}]});
assert.equal(target.children.length,3,'Dirección, título Reparto y cuadrícula');
assert.equal(target.children[0].className,'lab-cast-director');
assert.equal(target.children[0].children[1].textContent,'Kitao Sakurai');
const cards=target.children[2];
assert.equal(cards.className,'lab-cast-grid');
assert.equal(cards.children.length,3,'Rechazar entradas sin nombre, mantener roles opcionales');
const [ryu,ken,unknown]=cards.children;
assert.equal(ryu.children[0].className,'lab-cast-actor');
assert.equal(ryu.children[0].textContent,'Andrew Koji');
assert.equal(ryu.children[1].className,'lab-cast-role');
assert.equal(ryu.children[1].textContent,'Ryu');
assert.equal(ken.children[0].textContent,'Noah Centineo');
assert.equal(ken.children[1].textContent,'Ken Masters');
assert.equal(unknown.children.length,1);
ctx.labRenderCredits({director:'',cast:[]});
assert.equal(target.children.length,1,'Mantener Dirección sin reparto y sin datos');
assert.equal(target.children[0].children[1].textContent,'Sin datos');

// El scroll está en #view-collections.view como en Descubrir, con la tarjeta sin scroll.
// Las dos vistas mantienen el mismo estilo y dejan fija la navegación global.
assert.match(style,/#view-collections\.view\{[^}]*overflow-y:auto!important/);
assert.match(style,/#view-collections > \.panel\{[^}]*overflow:visible!important/);
assert.match(style,/#view-collections\.view\[hidden\]\{display:none!important\}/);
assert.match(style,/#view-collections\.view\{[^}]*scrollbar-color:#51667f transparent/);
assert.match(style,/\.lab-discover\.view\{[^}]*overflow-y:auto/);
assert.ok(base.includes('#view-collections .panel{min-height:0;overflow-y:auto'),
 'El cambio debe superar el estilo anterior del panel sin tocar la plantilla base');
assert.ok(!base.includes('.lab-cast-actor'),'No tocar el canal estable');
console.log('OK: actores y personajes separados, Dirección clara y Mis listas con scroll como Descubrir');
