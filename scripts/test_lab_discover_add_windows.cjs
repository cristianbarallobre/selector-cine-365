'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');
const ui=fs.readFileSync('lab_discover_windows.html','utf8');
const script=ui.split('<!-- LAB_DISCOVERY_SCRIPT -->')[1];
assert.ok(script);
const start=script.indexOf('function labHydrateAddedMovie(movie,film){');
const end=script.indexOf('function labAddCurrentMovie(){',start);
assert.ok(start>=0&&end>start,'Debe existir la transferencia inmediata de caché TMDB a Biblioteca');
const body=script.slice(start,end);
new vm.Script(body,{filename:'lab_hydrate_discovered.js'});
assert.match(body,/platformData\.set\(key,data\)/);
assert.match(body,/runtime:Number\(movie\.runtime\)/);
assert.match(body,/poster_path:movie\.poster_path/);
assert.match(body,/tmdb_id:Number\(movie\.id\)/);
assert.match(body,/vote_average:Number\(movie\.vote_average\)/);
assert.match(body,/subscriptions:Array\.isArray\(movie\.subscriptions\)/);
assert.match(script,/labHydrateAddedMovie\(movie,\{title,year\}\);\s*update\(\);/,
 'La lista debe recuperar metadatos ANTES de renderizar, no esperar a otra consulta');
assert.match(body,/platformApi\('\/api\/platforms\/lookup',\{title:film\.title,year:Number\(film\.year\)\|\|0,force:false\}\)/,
 'Sincronizar la caché de plataformas en segundo plano');
assert.match(body,/Number\(info\.tmdb_id\)===Number\(movie\.id\)/,
 'No sustituir una ficha correcta por otro resultado que coincida solo por título');

(async()=>{
 const data=new Map();
 let backend=null,reads=0,renders=0;
 const src={
  platformData:data,
  providerKey:f=>f.title.toLowerCase()+'|'+Number(f.year),
  platformApi:async()=>{
   reads++;
   return await new Promise(resolve=>{backend=resolve;});
  },
  scheduledPlatformRender:()=>{renders++},
  Number,Array,Promise
 };
 vm.createContext(src);vm.runInContext(body,src);
 const film={title:'La película recién añadida',year:'2026'};
 const movie={
  id:998877,title:film.title,poster_path:'/poster.png',vote_average:8.6,
  vote_count:218,runtime:105,genres:['Drama','Misterio'],
  subscriptions:['Filmin'],free:[],ads:[],rent:['Apple TV'],buy:[],
  watch_link:'https://example.com/streaming'
 };
 src.labHydrateAddedMovie(movie,film);
 const cache=data.get('la película recién añadida|2026');
 assert.ok(cache,'La ficha debe quedar disponible sin esperar a la API');
 assert.equal(cache.state,'ok');
 assert.equal(cache.tmdb_id,998877);
 assert.equal(cache.runtime,105);
 assert.equal(cache.poster_path,'/poster.png');
 assert.equal(cache.vote_average,8.6);
 assert.deepEqual(Array.from(cache.subscriptions),['Filmin']);
 assert.equal(reads,1,'No se debe hacer doble consulta inicial');
 assert.equal(renders,0,'La actualización en segundo plano todavía no terminó');
 backend({state:'ok',tmdb_id:998877,film_url:'https://example.com/streaming',runtime:105});
 await Promise.resolve();await Promise.resolve();
 assert.equal(renders,1,'Una consulta coincidente puede refrescar la vista');
 assert.equal(data.get('la película recién añadida|2026').tmdb_id,998877);
 backend=null;
 src.labHydrateAddedMovie(movie,film);
 backend({state:'ok',tmdb_id:12345,poster_path:'/wrong.png'});
 await Promise.resolve();await Promise.resolve();
 assert.equal(data.get('la película recién añadida|2026').tmdb_id,998877,'No sustituir por un resultado distinto');
 assert.equal(data.get('la película recién añadida|2026').poster_path,'/poster.png');
 backend=null;src.labHydrateAddedMovie(movie,film);
 backend({state:'ok',tmdb_id:998877,poster_path:'',runtime:0,genres:[],vote_count:0,vote_average:0});
 await Promise.resolve();await Promise.resolve();
 const retained=data.get('la película recién añadida|2026');
 assert.equal(retained.poster_path,'/poster.png','Un backend parcial no borra la carátula');
 assert.equal(retained.runtime,105,'Un backend parcial no borra la duración');
 assert.equal(retained.vote_average,8.6,'Un backend parcial no borra la nota');
 assert.deepEqual(Array.from(retained.genres),['Drama','Misterio'],'Un backend parcial no borra géneros');
 console.log('OK: Descubrir rellena la ficha de Biblioteca al instante y protege identidad y metadatos TMDB');
})().catch(error=>{console.error(error);process.exitCode=1});
