'use strict';
const fs=require('node:fs'),assert=require('node:assert/strict'),vm=require('node:vm');
const html=fs.readFileSync('android/app/src/main/assets/index.html','utf8');
const backend=fs.readFileSync('android/app/src/main/java/es/cine365/app/Backend.java','utf8');
const cloud=fs.readFileSync('android/app/src/main/java/es/cine365/app/CloudBackup.java','utf8');
const local=fs.readFileSync('android/app/src/main/java/es/cine365/app/LocalRecovery.java','utf8');
const activity=fs.readFileSync('android/app/src/main/java/es/cine365/app/MainActivity.java','utf8');
const gradle=fs.readFileSync('android/app/build.gradle','utf8');
function contains(s,needle,msg){assert.ok(s.includes(needle),msg||needle)}
function notContains(s,needle,msg){assert.ok(!s.includes(needle),msg||needle)}
contains(gradle,'versionCode 9025');
contains(gradle,"versionName '9.0.0-beta.13.11'");
contains(gradle,"applicationId 'es.cine365.app'");
contains(html,'Android v9 Beta 13.10');
contains(html,'id="cloudSyncQuick"');
contains(html,'aria-label="Sincronizar con OneDrive: enviar y recibir cambios"');
contains(html,"$('cloudSyncQuick').addEventListener('click'");
contains(html,"showView('settings')");
contains(html,"void syncDevices()");
contains(html,"pendingSyncToken=data.token");
contains(html,"cloudRequest('/api/cloud/sync','POST',{token:pendingSyncToken,decisions})");
contains(html,"Revisa listas, vistas y conflictos. Al confirmar");
notContains(html,'cloudAutoReceive(');
notContains(html,'setInterval(()=>{if(!document.hidden)');
notContains(html,"window.addEventListener('focus',()=>void cloud");
notContains(html,"document.addEventListener('visibilitychange',()=>{if(!document.hidden)");
notContains(html,'setInterval(()=>{fetch(');
notContains(html,"if(platformConfigured)void loadPlatforms(false)");
contains(activity,'webView.onPause();webView.pauseTimers()');
contains(activity,'webView.onResume();webView.resumeTimers()');
notContains(activity,'if(window.cineSyncResume)window.cineSyncResume()');
contains(cloud,'.put("auto_sync",false)');
contains(backend,'"skipped","manual_only"');
const libraryHandler=backend.slice(backend.indexOf('case "/api/library":'),backend.indexOf('case "/api/identities":'));
notContains(libraryHandler,'cloud.backup()', 'Saving local state must not upload OneDrive');
const matchHandler=backend.slice(backend.indexOf('case "/api/platforms/match"'),backend.indexOf('case "/api/onedrive/status"'));
notContains(matchHandler,'cloud.backup()','Pencil corrections must not upload OneDrive');
notContains(local,'backend.cloud.backup()','Restoring an old backup must not silently upload it');
for(const m of html.matchAll(/<script>([\s\S]*?)<\/script>/g))new Function(m[1]);

// Exercise real saveLocal() JavaScript: every write is a local /api/library POST only.
const start=html.indexOf('let lastStoredSnapshot=null;');
const end=html.indexOf('function replaceList(text){',start);
assert.ok(start>0&&end>start,'real saveLocal block found');
const state={version:1,active_id:'a',lists:[
 {id:'a',name:'Lista A',text:'☐ Casablanca (1942)'},
 {id:'b',name:'Lista B',text:'☐ Vertigo (1958)'}]};
const requests=[],saveLabel={textContent:''};
const ctx={
 token:'test-local-token',
 libraryState:state,
 membershipCache:null,
 saveTimer:null,
 saveQueue:Promise.resolve(),
 saveEpoch:0,
 syncActiveText:()=>{},
 $:()=>saveLabel,
 fetch:async(url,opt)=>{
  requests.push({url,method:opt.method,body:JSON.parse(opt.body)});
  return {ok:true,status:204};
 },
 cineLoadRevisions:async()=>{},
 cineRecentIntents:{clear:()=>{}},
 console,Promise,setTimeout:()=>0,clearTimeout:()=>{},
};
vm.createContext(ctx);
vm.runInContext(html.slice(start,end)+';this.commit=saveLocal;this.quietSave=scheduleSave;',ctx);
(async()=>{
 await ctx.commit();
 assert.equal(requests.length,1,'First local save');
 assert.equal(requests[0].url,'/api/library');
 await ctx.commit(true);
 assert.equal(requests.length,1,'Identical state: no disk writes');
 ctx.libraryState.active_id='b';
 ctx.quietSave(true);
 assert.notEqual(saveLabel.textContent,'Guardando…','No saving label on list switch');
 await ctx.commit(true);
 assert.equal(requests.length,2,'Selection may be stored locally');
 await ctx.commit(true);
 assert.equal(requests.length,2,'Repeated list: no redundant save');
 ctx.libraryState.lists[1].text='☑ Vertigo (1958)';
 await ctx.commit();
 assert.equal(requests.length,3,'Actual watched change is saved locally');
 assert.equal(saveLabel.textContent,'✓ En el móvil');
 assert.ok(requests.every(r=>r.url==='/api/library'),'No OneDrive request from saveLocal');

 // Exercise ONLY read-only remote preview at cold launch.
 const a=html.indexOf('let cloudAutoBusy=false;'),b=html.indexOf('void mobileRefreshHistory();',a);
 assert.ok(a>0&&b>a);
 const label={textContent:''},called=[];
 const cloudCtx={
  $:id=>id==='mobileAutoStatus'?label:{},
  cloudRequest:async(path)=>{called.push(path);return {enabled:true,activated:true};},
  cloudShow:()=>{},
  platformApi:async(path,body)=>{called.push(path);return {sources:1,conflicts:[]};},
  cloudRefresh:async()=>{},
  saveQueue:Promise.resolve(),saveTimer:null,
  saveLocal:async()=>{throw Error('Startup should not write without pending edits');},
  document:{hidden:false},window:{},console,Promise
 };
 vm.createContext(cloudCtx);
 vm.runInContext(html.slice(a,b)+';this.checkOnOpen=cloudCheckOnOpen;',cloudCtx);
 await cloudCtx.checkOnOpen();
 assert.deepEqual(called,['/api/cloud/status','/api/cloud/preview'],'Cold start only status and preview (no uploads)');
 cloudCtx.window.cineSyncResume();
 assert.equal(called.length,2,'Returning to app does not call OneDrive');
 assert.ok(label.textContent.includes('sin importar ni enviar'),'User sees read-only notice');
 console.log('OK: local save, no redundant write, startup read-only OneDrive, manual confirm, 9025');
})().catch(e=>{console.error(e);process.exitCode=1;});
