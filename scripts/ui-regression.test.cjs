'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const html=fs.readFileSync('plantilla.html','utf8');
const script=html.match(/<script>([\s\S]*?)<\/script>/)?.[1];
assert.ok(script,'La plantilla debe incluir JavaScript');
assert.doesNotThrow(()=>new Function(script),'El JavaScript de la aplicación debe compilar');
const begin=script.indexOf('// BEGIN MEMBERSHIP_HELPERS');
const end=script.indexOf('// END MEMBERSHIP_HELPERS');
assert.ok(begin>=0&&end>begin);
const helpers=new Function(script.slice(begin,end)+';return {buildMembershipRecords,changeFilmMembership};')();
function parse(text){
 const lines=text.split('\n');
 const entries=lines.map((line,idx)=>{
  const m=line.match(/^([☐☑]) (.*?)(?: \((\d{4})\))?$/);
  return m?{idx,checked:m[1]==='☑',title:m[2],name:m[2],year:m[3]||''}:null;
 }).filter(Boolean);
 return {lines,entries};
}
const group=f=>f.title.toLowerCase()+'|'+f.year;
const lists=[
 {id:'first',name:'Primera',text:'☑ Interstellar (2014)'},
 {id:'second',name:'Segunda',text:'☐ Interstellar (2014)\n☐ Dune (2021)'},
 {id:'third',name:'Tercera',text:'☐ Another film (2020)'}
];
const global=helpers.buildMembershipRecords(lists,parse,group,new Map());
assert.equal(global.rows.length,3,'Un único resultado por película entre listas');
const interstellar=global.rows.find(x=>x.title==='Interstellar');
assert.equal(interstellar.checked,true,'Vista si alguna copia está vista');
assert.deepEqual(interstellar._lists.map(x=>x.number),[1,2],'Números de listas por orden');
assert.deepEqual(interstellar._lists.map(x=>x.id),['first','second'],'IDs sin cambiar');
assert.equal(helpers.changeFilmMembership(lists,interstellar,new Set(['second','third']),parse,group,new Map()),2,'Añadir a tercera y quitar de primera');
assert.ok(!lists[0].text.includes('Interstellar'),'Quitada de la primera');
assert.equal((lists[1].text.match(/Interstellar/g)||[]).length,1,'No duplica una película existente');
assert.ok(lists[2].text.includes('☑ Interstellar (2014)'),'Nuevas apariciones conservan el visto');
const merged=helpers.buildMembershipRecords([
 {id:'es',name:'Español',text:'☐ El viaje (2020)'},
 {id:'en',name:'Inglés',text:'☑ The Trip (2020)'}
],parse,()=>12,new Map());
assert.equal(merged.rows.length,1,'Identidades equivalentes por TMDB comparten resultado');
assert.deepEqual(merged.rows[0]._searchTerms,['El viaje','The Trip'],'Busca también los nombres alternativos');
assert.ok(html.includes("className='stat-body'"),'Las filas de Estadísticas deben tener contenedor propio');
assert.ok(html.includes('#view-statistics .stat-pager{flex:0 0 30px'),'El paginador de Estadísticas debe ocupar el pie fijo');
assert.ok(html.includes('data-search-scope="all"'),'Incluye el selector de búsqueda global');
console.log('OK: sintaxis, búsqueda global, pertenencia a listas, numeración y pie de estadísticas');
