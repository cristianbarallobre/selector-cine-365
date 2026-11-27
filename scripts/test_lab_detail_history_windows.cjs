'use strict';
const fs=require('node:fs');
const assert=require('node:assert/strict');
const vm=require('node:vm');

const html=fs.readFileSync('lab_discover_windows.html','utf8');
const js=html.split('<!-- LAB_DISCOVERY_SCRIPT -->')[1];
assert.ok(js);
assert.match(html,/id="labDetailBack"[^>]*hidden/,'La primera película no muestra Atrás');
assert.match(html,/id="labPersonBack"/,'Los intérpretes tienen Atrás');
assert.match(html,/id="labPersonClose" aria-label="Cerrar todas las fichas"/,'La X del intérprete cierra todo');
assert.match(js,/labEnterDetail\('movie',id,record\)/,'Cada película registra su visita');
assert.match(js,/labEnterDetail\('person',id,record,role\)/,'Cada persona registra el tipo de filmografía');
assert.match(js,/labDetailOverlay\.scrollTop=labDetailCurrent\(\)\?\.kind==='movie'/,'Restaurar scroll de película');
assert.match(js,/labPersonOverlay\.scrollTop=labDetailCurrent\(\)\?\.kind==='person'/,'Restaurar scroll de intérprete');
new vm.Script('(()=>{\n'+js+'\n})();');

const start=js.indexOf("const labPersonOverlay=$('labPersonOverlay');");
const end=js.indexOf('function labDetailLink(label,url){',start);
assert.ok(start>=0&&end>start,'No se encuentran las funciones de navegación');
const handlers={};
class FakeElement{
 constructor(hidden=true){this.hidden=hidden;this.scrollTop=0;this.events={};}
 addEventListener(ev,cb){this.events[ev]=cb;}
 focus(){}
 click(){this.events.click?.();}
}
const overlay=new FakeElement(true),person=new FakeElement(true);
const buttons=Object.fromEntries(['labDetailBack','labPersonBack','labDetailClose','labPersonClose'].map(k=>[k,new FakeElement(true)]));
const $=id=>id==='labPersonOverlay'?person:buttons[id];
const ctx={
 $,
 labDetailOverlay:overlay,labDetailSerial:0,labCurrentMovie:null,
 labHideActorPhoto:()=>{},
 window:{addEventListener:(type,fn)=>{handlers[type]=fn;}}
};
ctx.labShowDetail=(id,record=true)=>{
 ctx.labEnterDetail('movie',id,record);
 ctx.labHidePerson();overlay.hidden=false;
};
let currentRole='';
ctx.labOpenPerson=(id,record=true,role='')=>{
 currentRole=role;ctx.labEnterDetail('person',id,record,role);person.hidden=false;
};
vm.createContext(ctx);
vm.runInContext(js.slice(start,end),ctx);
const trail=()=>Array.from(vm.runInContext('labDetailTrail',ctx),x=>x.kind+':'+x.id);
const openMovie=ctx.labShowDetail,openPerson=ctx.labOpenPerson;

openMovie(101);
assert.deepEqual(trail(),['movie:101']);
assert.equal(buttons.labDetailBack.hidden,true);
openPerson(201);
assert.deepEqual(trail(),['movie:101','person:201']);
assert.equal(buttons.labPersonBack.hidden,false);
person.scrollTop=245;
openMovie(102);
assert.equal(person.hidden,true);
assert.deepEqual(trail(),['movie:101','person:201','movie:102']);
overlay.scrollTop=130;
openPerson(202,true,'director');
openMovie(103);
assert.deepEqual(trail(),['movie:101','person:201','movie:102','person:202','movie:103']);
buttons.labDetailBack.click();
assert.deepEqual(trail(),['movie:101','person:201','movie:102','person:202']);
assert.equal(person.hidden,false);
assert.equal(currentRole,'director','Atrás conserva la filmografía de dirección');
buttons.labPersonBack.click();
assert.deepEqual(trail(),['movie:101','person:201','movie:102']);
buttons.labDetailBack.click();
assert.deepEqual(trail(),['movie:101','person:201']);
assert.equal(vm.runInContext('labDetailCurrent().scrollTop',ctx),245,'Se conserva el desplazamiento de la primera ficha de actor');
buttons.labPersonBack.click();
assert.deepEqual(trail(),['movie:101'],'Debe llegar a la primera película');
assert.equal(buttons.labDetailBack.hidden,true,'Sin Atrás al llegar al origen');
assert.equal(overlay.hidden,false,'No cierra la primera película involuntariamente');

// La X debe limpiar toda la cadena incluso estando en un intérprete.
openPerson(204);
buttons.labPersonClose.click();
assert.deepEqual(trail(),[]);
assert.equal(person.hidden,true);
assert.equal(overlay.hidden,true);
openMovie(300);
assert.deepEqual(trail(),['movie:300'],'Una apertura nueva no hereda el historial antiguo');
openPerson(400);
let escaped=false;
handlers.keydown({key:'Escape',preventDefault:()=>{escaped=true;}});
assert.equal(escaped,true);
assert.deepEqual(trail(),['movie:300'],'Escape retrocede en fichas anidadas');
handlers.keydown({key:'Escape',preventDefault:()=>{}});
assert.deepEqual(trail(),[],'Escape desde el origen cierra la ficha');
assert.equal(overlay.hidden,true);
console.log('OK: historial LAB completo de películas y actores, scroll, Atrás, X y Escape');
