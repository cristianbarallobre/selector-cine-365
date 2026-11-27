'use strict';
const fs=require('node:fs'),vm=require('node:vm'),assert=require('node:assert/strict');
const page=fs.readFileSync('plantilla.html','utf8');
const start=page.indexOf('globalThis.cineLocalFetch='),end=page.indexOf('// Abrir TXT',start);
assert.ok(start>0&&end>start);
const listeners=[],messages=[],requests=[],nodes=new Map();
const $=id=>{if(!nodes.has(id))nodes.set(id,{});return nodes.get(id);};
let mode='ok',revisionsPending=new Promise(()=>{});
const ctx={
 $,Blob,AbortController,Uint8Array,crypto:require('node:crypto').webcrypto,
 setTimeout,clearTimeout,window:{chrome:{webview:{addEventListener(_,fn){listeners.push(fn)},postMessage(message){
  messages.push(message);
  const [_,session,revision,action]=message.split(':');
  if(action==='put')queueMicrotask(()=>listeners.forEach(fn=>fn({data:'cine-draft-ack:'+session+':'+revision+':ok'})));
 }}}},
 token:'test',saveNeeded:false,editRevision:0,saveEpoch:0,savePump:null,saveQueue:Promise.resolve(),libraryETag:'"library-0"',
 recoverySavedRevision:-1,restoringBackup:false,libraryBase:null,membershipEditRevision:-1,membershipCache:null,
 libraryState:{version:1,active_id:'a',lists:[{id:'a',name:'Cine',text:'☐ Origen (2010)'}]},
 syncActiveText(){},cineLoadRevisions:()=>revisionsPending,cineRecentIntents:new Map(),cloudRefresh(){},cineRefreshRecovery(){},
 fetch:async(url,options)=>{
  if(url==='/api/library'){
   requests.push({url,options});
   if(mode==='fail')throw new Error('motor caído');
   return {ok:true,status:204,headers:{get:key=>key==='ETag'?'"library-'+requests.length+'"':null}};
  }
  if(url==='/never')return new Promise((resolve,reject)=>options.signal.addEventListener('abort',()=>reject(new Error('abort')),{once:true}));
  throw new Error('unexpected '+url);
 }
};
vm.createContext(ctx);
vm.runInContext(page.slice(start,end),ctx);
(async()=>{
 ctx.scheduleSave();
 await ctx.savePump;
 assert.equal(ctx.saveNeeded,false,'a stalled revisions GET must not block the save queue');
 assert.ok(messages[0].includes(':1:put:'));
 assert.ok(messages.some(m=>m.includes(':1:clear:')));
 assert.ok(requests[0].options.headers['X-Cine-Draft'].endsWith(':1'));
 ctx.libraryState.lists[0].text='☑ Origen (2010)';
 mode='fail';ctx.scheduleSave();await ctx.savePump;
 assert.equal(ctx.saveNeeded,true);
 assert.equal(ctx.cineDraftIsConfirmed(2),true);
 assert.match($('saveState').textContent,/Borrador local recuperable/);
 assert.ok(!messages.some(m=>m.includes(':2:clear:')),'never clear a failed draft');
 mode='ok';assert.equal(await ctx.flushLibrary(),true);
 assert.ok(messages.some(m=>m.includes(':2:clear:')),'clear only after acknowledgement');
 await assert.rejects(ctx.cineLocalFetch('/never',{},10),/tiempo de espera/);
 console.log('OK: independent native draft, failed save recovery, unblocked queue and bounded timeout');
})().catch(err=>{console.error(err);process.exitCode=1;});
