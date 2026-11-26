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

{
 const html=fs.readFileSync('plantilla.html','utf8');
 const reactions=fs.readFileSync('lab_reactions_windows.html','utf8');
 const take=(s,from,to)=>s.slice(s.indexOf(from),s.indexOf(to,s.indexOf(from)));
 const state=new Function([
  "const libraryState={lists:[{text:'☐ Matrix (1999)\\n☐ The Lord of the Rings (2001)\\n☐ Otra (2020)'}],reactions:{'matrix|1999':{title:'Matrix',year:'1999',favorite:true},'el señor de los anillos|2001':{title:'El señor de los anillos',year:'2001',tmdb_id:120,favorite:true}}};",
  "let labReactionsByTMDB=null,cineRevisionsLoaded=true;const cineRecentIntents=new Map(),cineRevisionCache=new Map();",
  "const filmIdentity=f=>f.title.toLowerCase(),filmKey=f=>filmIdentity(f)+'|'+f.year,knownTMDB=f=>f.title==='The Lord of the Rings'?120:0;",
  "const parse=text=>{const lines=text.split('\\n');return {lines,entries:lines.map((line,idx)=>({idx,checked:line[0]==='☑',title:line.slice(2,-7),year:line.slice(-5,-1)}))};};",
  "const setMarkOnLine=(line,checked)=>(checked?'☑':'☐')+line.slice(1);",
  take(html,'function buildIdentityGroups(){','// INICIO FUNCIONES DE BÚSQUEDA GLOBAL'),
  take(html,'function cineRevisionKey(f){','// La identificación se efectúa'),
  take(reactions,'function labReactionMovieID(f){','function labFavoritesRecords(){'),
  "assert.equal(reconcileWatched(),2);return libraryState.lists[0].text;"
 ].join('\n').replace('assert.equal(reconcileWatched(),2);','if(reconcileWatched()!==2)throw Error("No se heredaron los dos favoritos");'))();
 assert.equal(state,'☑ Matrix (1999)\n☑ The Lord of the Rings (2001)\n☐ Otra (2020)');
}
