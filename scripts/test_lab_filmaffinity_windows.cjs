'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');
const lab=fs.readFileSync('lab_discover_windows.html','utf8');
const js=lab.split('<!-- LAB_DISCOVERY_SCRIPT -->')[1];
assert.ok(js);
new vm.Script('(()=>{\n'+js+'\n})();');
assert.match(js,/labDetailLink\('↗ Buscar en FilmAffinity España',labFilmAffinitySearchURL\(m\)\)/);
assert.doesNotMatch(fs.readFileSync('plantilla.html','utf8'),/labFilmAffinitySearchURL/,'No alterar versión estable');
assert.doesNotMatch(js,/filmaffinity\.com\/es\/film\d+\.html/,'No inventar IDs propios de FilmAffinity');

const start=js.indexOf('function labFilmAffinitySearchURL(movie){');
const end=js.indexOf('function labTextChips(',start);
assert.ok(start>0&&end>start,'Falta el generador y enlace FilmAffinity');
const links=[];
const ctx={
 $:id=>{assert.equal(id,'labDetailLinks');return{appendChild:node=>links.push(node)}},
 document:{createElement:tag=>{assert.equal(tag,'a');return{}}}
};
vm.createContext(ctx);vm.runInContext(js.slice(start,end),ctx);
assert.equal(
 ctx.labFilmAffinitySearchURL({id:157336,title:'La Odisea',release_date:'2026-07-01'}),
 'https://www.filmaffinity.com/es/search.php?stype=title&stext=La%20Odisea'
);
assert.equal(
 ctx.labFilmAffinitySearchURL({title:'¿Qué fue de Ana & Luis?'}),
 'https://www.filmaffinity.com/es/search.php?stype=title&stext=%C2%BFQu%C3%A9%20fue%20de%20Ana%20%26%20Luis%3F'
);
assert.equal(ctx.labFilmAffinitySearchURL({title:' ',original_title:'The Matrix'}), // título vacío: no emitir consultas inventadas
 '','No inventar un enlace de título vacío');
assert.equal(ctx.labFilmAffinitySearchURL({}), '');
ctx.labDetailLink('↗ Buscar en FilmAffinity España',ctx.labFilmAffinitySearchURL({title:'Interstellar'}));
assert.equal(links.length,1);
assert.equal(links[0].target,'_blank');
assert.equal(links[0].rel,'noopener noreferrer');
assert.match(links[0].href,/filmaffinity\.com\/es\/search\.php/);
assert.equal(links[0].textContent,'↗ Buscar en FilmAffinity España');
console.log('OK: enlace FilmAffinity España por título, URL codificada, sin API ni coincidencias falsas');
