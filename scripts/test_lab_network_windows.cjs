'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');

const source=fs.readFileSync('lab_network_windows.js','utf8');
const injector=fs.readFileSync('lab_network_inject_desktop.go','utf8');
const backend=fs.readFileSync('lab_network_desktop.go','utf8');
const other=fs.readFileSync('lab_network_other.go','utf8');
const main=fs.readFileSync('main.go','utf8');
const shared=fs.readFileSync('plantilla.html','utf8');
new vm.Script(source,{filename:'lab_network_windows.js'});
assert.match(injector,/\/\/go:build windows/);
assert.match(backend,/\/\/go:build windows/);
assert.match(other,/\/\/go:build !windows/);
assert.doesNotMatch(injector,/if !appIsLab/,'También disponible en Windows estable');
assert.doesNotMatch(backend,/if !appIsLab/,'También disponible en Windows estable');
assert.match(main,/out = injectLabNetwork\(out\)/);
assert.match(main,/a\.serveLabNetworkLog\(w,\s*r\)/);
assert.match(injector,/CINE_NETWORK_INTERCEPTOR/,'Insertar el interceptor en su marcador explícito');
assert.ok(shared.indexOf('/* CINE_NETWORK_INTERCEPTOR */')>0&&shared.indexOf('/* CINE_NETWORK_INTERCEPTOR */')<shared.indexOf('void initLibrary()'));
assert.match(shared,/Guardado no confirmado/,'No declarar que un guardado incierto ha fallado');
assert.match(shared,/Guardado: /,'Distinguir errores externos al leer el estado de la copia adicional');
assert.match(backend,/labNetMaxLogBytes=256\*1024/,'Limitar tamaño del registro');
assert.match(backend,/DisallowUnknownFields/,'Validar campos del registro');
assert.ok(!shared.includes('labOriginalFetch'),'La plantilla estable debe permanecer intacta');
assert.ok(!source.includes('console.log'),'No imprimir tokens o títulos en consola');
assert.ok(!source.includes('JSON.stringify(options)'),'No registrar cuerpos de solicitudes');

function http(status){
 return {status,ok:status>=200&&status<300,body:{cancel:async()=>{}}};
}
function build(fetchHandler){
 const calls=[],delays=[];
 const win={
  location:{href:'http://127.0.0.1:38111/',origin:'http://127.0.0.1:38111'},
  fetch:async(url,options)=>{
   calls.push({url,options});
   if(url==='/api/lab/network-log')return http(204);
   return fetchHandler(url,options,calls.filter(x=>x.url!=='/api/lab/network-log').length);
  }
 };
 const ctx={
  window:win, URL,Date,Set,Map,Promise,
  token:'private-session-token',
  setTimeout:(callback,ms)=>{delays.push(ms);callback();}
 };
 vm.createContext(ctx);
 vm.runInContext(source,ctx,{timeout:1500});
 return {fetch:win.fetch,calls,delays};
}
async function flush(){await Promise.resolve();await Promise.resolve();await Promise.resolve();}
(async()=>{
 let test=build((url,opt,n)=>{
  if(n<=2)throw new TypeError('Failed to fetch');
  return http(200);
 });
 let r=await test.fetch('/api/discover/detail?token=PRIVATE',{
  headers:{'X-Atlas-Token':'private-session-token'},cache:'no-store'
 });
 assert.equal(r.status,200);
 assert.deepEqual(test.delays,[300,900],'Dos reintentos con espera progresiva');
 assert.equal(test.calls.length,3,'Tres lecturas como máximo');

 test=build((url,opt,n)=>http(n<=2?503:200));
 r=await test.fetch('/api/extra-backup/status');
 assert.equal(r.status,200,'Repetir HTTP 503 solo cuando es consulta GET');
 assert.equal(test.calls.length,3);

 test=build(()=>http(502));
 r=await test.fetch('/api/discover/detail?q=CLAVE_PRIVADA');
 assert.equal(r.status,502,'HTTP fallido se devuelve al código existente, no se oculta');
 await flush();
 assert.equal(test.calls.filter(x=>x.url==='/api/discover/detail?q=CLAVE_PRIVADA').length,3);
 const reports=test.calls.filter(x=>x.url==='/api/lab/network-log');
 assert.equal(reports.length,1,'Registrar solo el error final');
 const report=JSON.parse(reports[0].options.body);
 assert.deepEqual(Object.keys(report).sort(),['attempts','kind','method','route','status']);
 assert.equal(report.route,'/api/discover/detail');
 assert.equal(report.attempts,3);
 assert.equal(report.kind,'http');
 assert.equal(report.status,502);
 assert.ok(!reports[0].options.body.includes('CLAVE_PRIVADA'),'Excluir querystrings y credenciales del registro');
 assert.ok(!reports[0].options.body.includes('private-session-token'),'No filtrar token de sesión');

 test=build(()=>http(429));
 r=await test.fetch('/api/platforms/cache');
 assert.equal(r.status,429,'No insistir si TMDB impone límite de peticiones');
 assert.equal(test.calls.filter(x=>x.url==='/api/platforms/cache').length,1);
 test=build(()=>http(403));
 r=await test.fetch('/api/library');
 assert.equal(r.status,403,'No reintentar error de autorización');
 assert.equal(test.calls.length,1,'Una petición sin reintentos');

 test=build(()=>{throw new TypeError('Failed to fetch')});
 await assert.rejects(test.fetch('/api/library',{method:'POST',body:'LISTA_SENSIBLE'}),/No se ha repetido la operación/);
 await flush();
 assert.equal(test.calls.filter(x=>x.url==='/api/library').length,1,'Nunca duplicar escrituras POST');
 const postReport=JSON.parse(test.calls.find(x=>x.url==='/api/lab/network-log').options.body);
 assert.deepEqual(postReport,{route:'/api/library',method:'POST',kind:'network',status:0,attempts:1});
 assert.ok(!JSON.stringify(postReport).includes('LISTA_SENSIBLE'),'No registrar el cuerpo guardado');

 test=build(()=>http(503));
 r=await test.fetch('/api/extra-backup/now',{method:'POST',body:'TOKEN_PRIVADO'});
 assert.equal(r.status,503,'La app debe recibir el HTTP de copia adicional');
 assert.equal(test.calls.filter(x=>x.url==='/api/extra-backup/now').length,1,'No repetir copias');
 await flush();
 assert.equal(JSON.parse(test.calls.find(x=>x.url==='/api/lab/network-log').options.body).attempts,1);

 test=build(()=>{throw new TypeError('Failed to fetch')});
 await assert.rejects(test.fetch('/api/ping'),/Failed to fetch/);
 assert.equal(test.calls.length,1,'No repetir ni registrar pings periódicos');
 test=build(()=>{throw new TypeError('Failed to fetch')});
 await assert.rejects(test.fetch('https://example.org/api/external'),/Failed to fetch/);
 assert.equal(test.calls.length,1,'No tocar peticiones a otros dominios');

 test=build(()=>{const error=new Error('Usuario canceló');error.name='AbortError';throw error});
 await assert.rejects(test.fetch('/api/discover/list'),/Usuario canceló/);
 assert.equal(test.calls.length,1,'No reintentar cancelaciones del usuario');

 test=build(()=>{throw new TypeError('Failed to fetch')});
 await assert.rejects(test.fetch('/api/discover/search'),/después de 3 intentos/);
 assert.equal(test.calls.filter(x=>x.url==='/api/discover/search').length,3);
 await flush();
 assert.equal(test.calls.filter(x=>x.url==='/api/lab/network-log').length,1,'Registrar final de una lectura fallida');

 console.log('OK: LAB-02.14 reintentos GET, errores HTTP/red, no duplicar escrituras, registro privado y aislamiento');
})().catch(e=>{console.error(e);process.exitCode=1});
