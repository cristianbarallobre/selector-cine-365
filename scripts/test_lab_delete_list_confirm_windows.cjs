'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');

const lab=fs.readFileSync('lab_discover_windows.html','utf8');
const base=fs.readFileSync('plantilla.html','utf8');
const injector=fs.readFileSync('lab_discover_windows.go','utf8');
const stable=base.slice(base.indexOf('function deleteList(id){'),base.indexOf('function showView(name){'));
assert.ok(stable.startsWith('function deleteList(id){'),'La función original debe seguir en la plantilla estable');
assert.match(stable,/confirm\('¿Eliminar «'/,'La plantilla estable no se modifica');
assert.match(injector,/async function deleteList\(id\)/,'LAB debe esperar respuesta asíncrona del diálogo');
assert.match(injector,/await labConfirmDeleteList\(l\)/,'LAB debe verificar confirmación antes de borrar');
assert.match(lab,/role="alertdialog" aria-modal="true"/);
assert.match(lab,/id="labDeleteListOverlay" hidden/);
assert.match(lab,/id="labDeleteListName"/);
assert.match(lab,/id="labDeleteListCancel"/);
assert.match(lab,/id="labDeleteListConfirm"/);
assert.match(lab,/#labDeleteListOverlay\[hidden\]\{display:none!important\}/);
const script=lab.split('<!-- LAB_DISCOVERY_SCRIPT -->')[1];
new vm.Script('(()=>{\n'+script+'\n})();');

const original="if(!confirm('¿Eliminar «'+l.name+'»? Esta lista dejará de aparecer en la aplicación.'))return;";
assert.ok(stable.includes(original));
const experimental=stable.replace('function deleteList(id){','async function deleteList(id){')
 .replace(original,'if(!await labConfirmDeleteList(l))return; if(libraryState.lists.length<2||!libraryState.lists.some(item=>item.id===id))return;');
assert.ok(!experimental.includes("!confirm("));
const start=script.indexOf('let labDeleteListResolver=null');
assert.ok(start>0,'No se encuentra el módulo del diálogo');
const confirmations=script.slice(start);
class Element {
 constructor(id){this.id=id;this.hidden=id==='labDeleteListOverlay';this.events={};this.isConnected=true;this.textContent='';}
 addEventListener(name,fn){this.events[name]=fn;}
 focus(){doc.activeElement=this;}
 click(){this.events.click?.({target:this});}
}
const ids=['labDeleteListOverlay','labDeleteListName','labDeleteListCancel','labDeleteListConfirm'];
const elements=Object.fromEntries(ids.map(id=>[id,new Element(id)]));
const doc={activeElement:new Element('trigger')};
const handlers={};
let saved=0,restarts=0;
const collection=()=>({lists:[{id:'a',name:'Lista <script>',text:'☐ película A'},{id:'b',name:'Otra lista',text:'☐ película B'}],active_id:'a'});
const ctx={
 $:id=>{assert.ok(elements[id],id);return elements[id];},
 document:doc,
 window:{addEventListener:(event,fn,capture)=>{assert.equal(capture,true);handlers[event]=fn;}},
 libraryState:collection(),
 virtualMode:false,rawText:'',
 syncActiveText:()=>{},
 sanitizeComparisonIDs:()=>{},
 activeCollection:()=>ctx.libraryState.lists.find(x=>x.id===ctx.libraryState.active_id),
 restartList:()=>{restarts++;},
 renderCollectionControls:()=>{},
 refreshVirtualContent:()=>{},
 scheduleSave:()=>{saved++;}
};
vm.createContext(ctx);
vm.runInContext(confirmations+'\n'+experimental,ctx);

const overlay=elements.labDeleteListOverlay, cancel=elements.labDeleteListCancel, okay=elements.labDeleteListConfirm;
async function test(){
 let p=ctx.deleteList('a');
 assert.equal(overlay.hidden,false,'El diálogo se abre sin borrar todavía');
 assert.equal(elements.labDeleteListName.textContent,'Lista <script>','El nombre se inserta como texto, no HTML');
 assert.equal(doc.activeElement,cancel,'Foco inicial seguro: Cancelar');
 assert.equal(saved,0);
 cancel.click();await p;
 assert.equal(overlay.hidden,true);
 assert.equal(ctx.libraryState.lists.length,2,'Cancelar no debe eliminar');
 assert.equal(saved,0,'Cancelar no debe guardar ni sincronizar');

 p=ctx.deleteList('a');
 const other=await ctx.labConfirmDeleteList({name:'Otra'});
 assert.equal(other,false,'No apilar varios diálogos');
 overlay.events.click({target:overlay});await p;
 assert.equal(saved,0,'El clic exterior cancela');

 p=ctx.deleteList('a');
 let prevented=0,stopped=0;
 handlers.keydown({key:'Escape',preventDefault:()=>{prevented++;},stopImmediatePropagation:()=>{stopped++;}});
 await p;
 assert.equal(prevented,1);
 assert.equal(stopped,1,'Escape no alcanza otros manejadores del navegador');
 assert.equal(saved,0,'Escape cancela el borrado');

 p=ctx.deleteList('a');
 handlers.keydown({key:'Tab',shiftKey:true,preventDefault:()=>{prevented++;}});
 assert.equal(doc.activeElement,okay,'Tab inverso se mantiene en la ventana');
 handlers.keydown({key:'Tab',shiftKey:false,preventDefault:()=>{prevented++;}});
 assert.equal(doc.activeElement,cancel,'Tab no sale del diálogo');
 okay.click();await p;
 assert.equal(ctx.libraryState.lists.length,1);
 assert.equal(ctx.libraryState.lists[0].id,'b');
 assert.equal(ctx.libraryState.deleted_lists.a,true,'Conservar tombstone de sincronización');
 assert.equal(ctx.libraryState.active_id,'b','Seleccionar otra lista si se borra la activa');
 assert.equal(saved,1,'El guardado se realiza una sola vez tras confirmar');
 assert.equal(restarts,1);
 await ctx.deleteList('b');
 assert.equal(saved,1,'No se permite borrar la última lista');
 assert.equal(overlay.hidden,true);
 console.log('OK: confirmación LAB integrada, cancelación con Escape y clic exterior, foco, borrado y guardado originales');
}
test().catch(err=>{console.error(err);process.exitCode=1;});
