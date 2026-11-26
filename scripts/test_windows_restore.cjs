'use strict';
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');

const page = fs.readFileSync('plantilla.html', 'utf8');
const start = page.indexOf('let restoringBackup=false;');
const writeStart = page.indexOf('async function writeLibrary(payload)');
const writeEnd = page.indexOf('async function flushLibrary()', writeStart);
const end = page.indexOf('// Copia adicional', start);
assert.ok(start >= 0 && end > start);

async function restore(backup, accept) {
  const original = {version: 1, active_id: 'old', lists: [{id:'old', name:'Anterior', text:'Anterior'}]};
  let handler, sent, warning;
  const status = {};
  const context = {
    document: {body:{}}, Blob, token:'prueba', saveTimer:null, saveEpoch:0, libraryETag:'"library-1"',
    saveQueue:Promise.resolve(), libraryState:original, confirm:()=>true,
    librarySizeLimit:256*1024*1024, flushLibrary:async()=>true,
    alert:message=>{warning=message;}, validLibrary:()=>true,
    $:id=>id==='backupFile'?{addEventListener:(_,fn)=>{handler=fn;}}:id==='reloadLibraryAfterConflict'?{addEventListener(){}}:status,
    activeCollection:()=>context.libraryState.lists[0], sanitizeComparisonIDs(){},
    restartList(){}, showView(){}, cineLoadRevisions:async()=>{},
    fetch:async(url,opts)=>{
      assert.equal(url,'/api/library');
      assert.equal(opts.headers['If-Match'],'"library-1"');
      sent=JSON.parse(opts.body);
      return {ok:accept, status:accept?204:412, text:async()=>'Biblioteca inválida', headers:{get:name=>name==='ETag'?'"library-2"':null}};
    }
  };
  vm.runInNewContext(page.slice(writeStart,writeEnd)+page.slice(start,end),context);
  const contents=JSON.stringify({...backup,backup_created_at:'2026-11-24'});
  await handler({target:{files:[{size:Buffer.byteLength(contents),text:async()=>contents}]}});
  assert.deepEqual(sent,backup);
  assert.deepEqual(JSON.parse(JSON.stringify(context.libraryState)),accept?backup:original);
  assert.equal(Boolean(warning),!accept);
  assert.equal(context.libraryETag,accept?'"library-2"':'"library-1"');
  assert.equal(context.document.body.inert,false);
}

(async()=>{
  const old={version:1,active_id:'atlas',lists:[{id:'atlas',name:'Cine',text:'☑ Vértigo (1958)'}]};
  const full={...old,lists:[{...old.lists[0],order_changed_at:100},{id:'extra',name:'Otra',text:'☐ Origen (2010)'}],
    comparison_ids:['atlas','extra'],deleted_lists:{retirada:true},tmdb_overrides:{'vertigo|1958':426},
    reactions:{'vertigo|1958':{favorite:true,rating:5,title:'Vértigo',year:'1958',changed_at:200}}};
  await restore(old,true);
  await restore(full,true);
  await restore(full,false);
})().catch(err=>{console.error(err);process.exitCode=1;});

async function testSavePump() {
  const from = page.indexOf('function scheduleSave()');
  const to = page.indexOf('// Abrir TXT', from);
  const status = {};
  const requests = [];
  const context = {
    saveNeeded:false, editRevision:0, saveEpoch:0, savePump:null,
    saveQueue:Promise.resolve(), recoverySavedRevision:-1, restoringBackup:false,
    libraryETag:'"library-0"', libraryState:{text:'Primera marca'}, token:'test',
    membershipCache:null, syncActiveText(){}, cineLoadRevisions:async()=>{},
    cineRecentIntents:new Map(), cloudRefresh(){}, setTimeout(){},
    cineRefreshRecovery(){},
    $:()=>status,
    fetch:(url,options)=>new Promise(resolve=>requests.push({
      text:JSON.parse(options.body).text, etag:options.headers['If-Match'], resolve
    }))
  };
  const tick=()=>new Promise(resolve=>setImmediate(resolve));
  const response=(code,etag,recovered=false)=>({
    ok:code===204,status:code,
    headers:{get:name=>name==='ETag'?etag:name==='X-Cine-Recovery'&&recovered?'saved':null}
  });
  vm.runInNewContext(page.slice(from,to),context);
  context.scheduleSave();
  await tick();
  assert.equal(requests.length,1,'El primer cambio se envía sin esperar un temporizador');
  context.libraryState.text='Segunda marca';
  context.scheduleSave();
  context.libraryState.text='Tercera marca';
  context.scheduleSave();
  requests[0].resolve(response(204,'"library-1"'));
  await tick();
  assert.equal(requests.length,2);
  assert.equal(requests[1].text,'Tercera marca');
  assert.equal(requests[1].etag,'"library-1"');
  requests[1].resolve(response(412,null,true));
  await context.savePump;
  assert.equal(context.saveNeeded,true);
  assert.equal(context.recoverySavedRevision,3);
  assert.match(status.textContent,/copia local de recuperación/);
}

testSavePump().catch(err=>{console.error(err);process.exitCode=1;});
