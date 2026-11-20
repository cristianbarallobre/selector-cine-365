'use strict';
const fs=require('node:fs'),assert=require('node:assert/strict'),vm=require('node:vm');
const html=fs.readFileSync('android/app/src/main/assets/index.html','utf8');
const main=fs.readFileSync('android/app/src/main/java/es/cine365/app/MainActivity.java','utf8');
const backend=fs.readFileSync('android/app/src/main/java/es/cine365/app/Backend.java','utf8');
const gradle=fs.readFileSync('android/app/build.gradle','utf8');
assert.equal((main.match(/backend = new Backend\(this\);/g)||[]).length,1,'No construir Backend dos veces');
for(const str of ['webView.onPause();webView.pauseTimers();','webView.onResume();webView.resumeTimers();',
 'appForeground=false;','appForeground=true;','window.cineBackgroundFlush','backend.saveOnExit(obj)']){
 assert.ok(main.includes(str),'Ciclo de vida: '+str);
}
for(const str of ['sync.observe(validLibrary(previous)?previous:null,state);',
 'if(!activity.isAppForeground())return ok(new JSONObject()','"skipped","background"',
 'case "/api/cloud/sync"','case "/api/cloud/backup"']){
 assert.ok(backend.includes(str),'Conservar cambios y seguridad: '+str);
}
for(const str of ['Date.now()-cloudAutoCheckedAt<55000','},60000);','if(document.hidden)return;',
 'window.cineBackgroundFlush=()=>','void saveLocal();',"await initLibrary();restartList();"]){
 assert.ok(html.includes(str),'Optimización sin pérdida de sincronización: '+str);
}
assert.ok(!html.includes("setInterval(()=>{fetch('/api/ping'"),'No ping Android cada 12 segundos');
assert.ok(gradle.includes('versionCode 9016')&&gradle.includes("versionName '9.0.0-beta.13.2'"));
assert.ok(gradle.includes("applicationId 'es.cine365.app'"),'Mismo paquete para instalar encima');
for(const match of html.matchAll(/<script>([\s\S]*?)<\/script>/g))new Function(match[1]);

// Ejercitar la función REAL de comprobación OneDrive con reloj y backend simulados.
const start=html.indexOf('let cloudAutoBusy=false,cloudAutoCheckedAt=0;');
const end=html.indexOf('window.cineSyncResume=',start);
assert.ok(start>=0&&end>start);
let now=100000,requests=0,status='no_changes';
const feedback={textContent:''};
const document={hidden:true,activeElement:null,querySelector:()=>null};
const ctx={document,Date:{now:()=>now},saveEpoch:0,saveTimer:null,
 cloudSyncBusy:false,cloudUserRequests:0,pendingSyncToken:null,
 saveQueue:Promise.resolve(),$:id=>id==='mobileAutoStatus'?feedback:{hidden:true},
 cloudBusyStart:()=>{},cloudBusyEnd:()=>{},
 cloudRequest:async()=>{requests++;return {skipped:status};},
 cloudRefresh:async()=>{},initLibrary:async()=>{},restartList:()=>{},console};
vm.createContext(ctx);
vm.runInContext(html.slice(start,end)+';this.check=cloudAutoReceive;',ctx);
(async()=>{
 await ctx.check();assert.equal(requests,0,'no consultar con ventana oculta');
 document.hidden=false;await ctx.check();assert.equal(requests,1,'consultar al abrir');
 now+=30000;await ctx.check();assert.equal(requests,1,'no consultar a 30s');
 now+=26000;await ctx.check();assert.equal(requests,2,'consultar a 56s');
 now+=55000;status='background';await ctx.check();assert.equal(requests,3);
 status='no_changes';await ctx.check();assert.equal(requests,4,'reintentar tras pausa');
 document.hidden=true;now+=60000;await ctx.check();assert.equal(requests,4,'sin consulta en segundo plano');
 console.log('OK: una inicialización, WebView pausa/reanudación, OneDrive 60s, sin pings, guardado y Beta 13.2');
})().catch(e=>{console.error(e);process.exitCode=1;});
