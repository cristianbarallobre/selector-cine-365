'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');

for(const path of [process.argv[2]||'plantilla.html']){
 const html=fs.readFileSync(path,'utf8');
 const start=html.indexOf('function importTxtAsNewCollection('),end=html.indexOf('function replaceList(',start);
 assert.ok(start>=0&&end>start,'TXT import helper missing in '+path);
 assert.match(html,/\.platform-select\{cursor:pointer;/,'sort platform selector should show hand');
 assert.match(html,/const imported=importTxtAsNewCollection\(txt,file\.name\)/,'TXT file handler creates a collection');
 assert.doesNotMatch(html,/notify\(replaceList\(txt\)/,'TXT file handler must never replace active collection');
 const scope=new Function(`
 const listSizeLimit=900000;
 let rawText='☐ Antiguo (1980)',virtualMode=false;
 const libraryState={lists:[{id:'atlas',name:'Selector Cine 365 · 365',text:rawText}],active_id:'atlas'};
 const calls={restart:0,save:0,view:''};
 const activeCollection=()=>libraryState.lists.find(l=>l.id===libraryState.active_id);
 const parse=text=>({entries:String(text).split('\\n').filter(line=>/^[☐☑] /.test(line))});
 const syncActiveText=()=>{activeCollection().text=rawText;};
 const reconcileWatched=()=>{};
 const restartList=()=>{calls.restart++;};
 const scheduleSave=()=>{calls.save++;};
 const showView=v=>{calls.view=v;};
 `+html.slice(start,end)+`
 return {importTxtAsNewCollection,libraryState,calls,getRawText:()=>rawText};
 `)();
 const r=scope.importTxtAsNewCollection('☑ Uno (2020)\n☐ Dos (2021)','Mis películas favoritas.txt');
 assert.equal(r.ok,true);
 assert.equal(r.count,2);
 assert.equal(r.name,'Mis películas favoritas');
 assert.equal(scope.libraryState.lists.length,2);
 assert.equal(scope.libraryState.lists[0].text,'☐ Antiguo (1980)','original list must be untouched');
 assert.equal(scope.libraryState.active_id,scope.libraryState.lists[1].id);
 assert.equal(scope.calls.view,'library');
 assert.equal(scope.calls.save,1);
 assert.equal(scope.calls.restart,1);
 assert.equal(scope.getRawText(),'☑ Uno (2020)\n☐ Dos (2021)');
 const second=scope.importTxtAsNewCollection('☐ Tres (2022)','Mis películas favoritas.TXT');
 assert.equal(second.name,'Mis películas favoritas (2)');
 assert.equal(scope.libraryState.lists.length,3);
 const before=scope.libraryState.lists.length;
 assert.equal(scope.importTxtAsNewCollection('Texto sin registros','Inválido.txt').ok,false);
 assert.equal(scope.libraryState.lists.length,before);
 while(scope.libraryState.lists.length<40)scope.libraryState.lists.push({id:'x'+scope.libraryState.lists.length,name:'Extra',text:''});
 assert.equal(scope.importTxtAsNewCollection('☐ Cuatro (2023)','Límite.txt').ok,false);
 assert.equal(scope.libraryState.lists.length,40);
 for(const script of html.matchAll(/<script>([\s\S]*?)<\/script>/g))assert.doesNotThrow(()=>new Function(script[1]));
 console.log('OK: TXT import and platform selector behavior '+path);
}
