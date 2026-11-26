'use strict';
const fs=require('node:fs');
const vm=require('node:vm');
const assert=require('node:assert/strict');
const page=fs.readFileSync('plantilla.html','utf8');
const parseSource=page.slice(page.indexOf('const pendingSymbols='),page.indexOf('function update(){'));
const closeSource=page.slice(page.indexOf('function describePendingChanges(){'),page.indexOf("window.addEventListener('beforeunload'"));
const nodes=new Map();
const $=id=>{
 if(!nodes.has(id))nodes.set(id,{textContent:'',dataset:{},hidden:true,disabled:false,handlers:{},addEventListener(event,fn){this.handlers[event]=fn;}});
 return nodes.get(id);
};
const before={version:1,active_id:'a',lists:[{id:'a',name:'Pendientes',text:'☐ 1. Vértigo (1958)\n☐ 2. Origen (2010)'}]};
const after={version:1,active_id:'a',lists:[{id:'a',name:'Pendientes',text:'☑ 1. Vértigo (1958)\n☐ 2. Alien (1979)'}],reactions:{alien:{title:'Alien',favorite:true}}};
let exported,posted;
const webview={addEventListener(event,fn){this.handler=fn;},postMessage(msg){posted=msg;}};
const ctx={$,window:{chrome:{webview}},document:{body:{inert:false}},MutationObserver:class{observe(){}},
 saveNeeded:true,libraryReady:true,libraryBase:JSON.stringify(before),libraryState:after,
 syncActiveText(){},savePump:null,recoverySavedRevision:-1,editRevision:1,restoringBackup:false,cloudManualBusy:false,
 flushLibrary:async()=>false,fetch:async()=>{throw new Error('Servidor inaccesible');},token:'test',
 setTimeout,clearTimeout,AbortController,saveDownload(name,body,mime){exported={name,body,mime};}};
vm.createContext(ctx);
vm.runInContext(parseSource+closeSource,ctx);
(async()=>{
 const changes=ctx.describePendingChanges();
 assert.match(changes,/Pendientes/);
 assert.match(changes,/añadidas: Alien/);
 assert.match(changes,/retiradas: Origen/);
 assert.match(changes,/Vértigo \(1958\) \(vista\)/);
 assert.match(changes,/Favoritos o valoraciones: Alien/);
 assert.equal(await ctx.window.cinePrepareClose(),false);
 assert.equal(ctx.document.body.inert,false,'Se debe desbloquear la interfaz cuando falla el guardado');
 assert.match(ctx.window.cineCloseFailureDetail,/Servidor inaccesible/);
 assert.equal($('saveRecovery').hidden,false);
 assert.equal($('saveState').dataset.warning,'true');
 $('emergencyBackup').handlers.click();
 assert.match(exported.name,/^Cine365_emergencia_.*\.json$/);
 assert.equal(JSON.parse(exported.body).lists[0].text,after.lists[0].text);
 await webview.handler({data:'cine-prepare-close:7'});
 assert.match(posted,/^cine-close-pending:7:.*Cambios sin confirmar/s);
 ctx.fetch=async()=>({status:204});
 assert.equal(await ctx.window.cinePrepareClose(),true);
 assert.equal(ctx.saveNeeded,false);
 assert.equal(ctx.recoverySavedRevision,1);
 assert.equal(await ctx.window.cinePrepareClose(),true,'Sin cambios pendientes se puede cerrar');
 ctx.saveNeeded=true;ctx.editRevision=2;ctx.libraryBase=null;
 assert.match(ctx.describePendingChanges(),/copia de referencia/);
 ctx.libraryBase=JSON.stringify(before);
 ctx.flushLibrary=async()=>{ctx.saveNeeded=false;return true;};
 await $('retryLocalSave').handlers.click();
 assert.equal($('retryLocalSave').disabled,false);
 assert.equal($('saveRecovery').hidden,true);
 await assert.rejects(ctx.awaitCloseStep(new Promise(()=>{}),Date.now()+1),/no respondió a tiempo/);
 const main=fs.readFileSync('main.go','utf8'),win=fs.readFileSync('window_windows.go','utf8'),mac=fs.readFileSync('macos_view_darwin.m','utf8');
 assert.match(main,/!nativeDesktopWindowActive.Load\(\) && time.Since/);
 assert.match(win,/nativeDesktopWindowActive.Store\(true\)/);
 assert.match(win,/cmd.Wait\(\)\s+nativeDesktopWindowActive.Store\(false\)/);
 assert.match(win,/if confirmCloseWithoutSave\(hwnd, w.closeDetail\) \{\s+procDestroyWindow.Call\(hwnd\)/);
 assert.match(mac,/Cerrar de todas formas/);
 assert.match(mac,/reply\(choice == NSAlertSecondButtonReturn\)/);
 assert.match(mac,/dispatch_after\(dispatch_time\(DISPATCH_TIME_NOW, 20 \* NSEC_PER_SEC\)/);
 console.log('OK: cambios sin confirmar, recuperación, exportación local, reintento, tiempo límite y cierre nativo');
})().catch(err=>{console.error(err);process.exitCode=1;});
