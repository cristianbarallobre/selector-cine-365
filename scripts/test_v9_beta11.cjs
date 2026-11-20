'use strict';
const assert=require('node:assert/strict'),fs=require('node:fs');
function source(p){
 const s=fs.readFileSync(p,'utf8');
 const script=s.match(/<script>([\s\S]*?)<\/script>/)?.[1];
 assert.ok(script,'JS existente en '+p);
 assert.doesNotThrow(()=>new Function(script),'JS debe compilar '+p);
 for(const needle of ['cloudAutoReceive','cineRecentIntents','cineLoadRevisions','film-check-action','togg.append(tick,poster,main)','data-search-scope="all"','filmMembershipIndex','applyFilmMembershipLists','membershipOverlay','Conectar TMDB']){
  assert.ok(s.includes(needle),'Falta en '+p+': '+needle);
 }
 return {html:s,script};
}
const desktop=source('plantilla.html'),android=source('android/app/src/main/assets/index.html');
assert.ok(desktop.html.includes('data-view="instructions"'),'Windows debe conservar Instrucciones');
assert.ok(android.html.includes('data-view="help"'),'Android debe conservar Ayuda');
for(const {html,script} of [desktop,android]){
 assert.ok(!html.includes('main.inert=true'),'No bloquear toda la interfaz');
 assert.ok(script.includes('const readonly=virtualMode&&!global'),'Las coincidencias son solo lectura, la búsqueda global editable');
 assert.ok(script.includes('if(f._global)tick.dataset.globalGroup=f._group'),'Marcar visto con tick global');
 assert.ok(script.includes('else tick.dataset.idx=f.idx'),'Marcar visto con tick local');
 assert.ok(script.includes('if(!global&&selectedListPlatforms.size)'),'Filtro streaming no oculta búsqueda global');
 assert.ok(script.includes("num.addEventListener('click',()=>openFilmInList(m))"),'Números navegables');
}
assert.ok(desktop.script.includes('body.appendChild(line)'),'Estadísticas con cuerpo independiente');
assert.ok(desktop.html.includes('#view-statistics .stat-pager{flex:0 0 32px'),'Pie fijo de estadísticas Windows');
assert.ok(android.html.includes('Android: estadísticas completas sin paginación'),'Mantener estadísticas Android sin paginar');
assert.ok(android.html.includes('window.CineAndroid?Math.min(2000,'),'Mantener carga incremental Android');
const gradle=fs.readFileSync('android/app/build.gradle','utf8');
assert.ok(gradle.includes("applicationId 'es.cine365.app'"),'Identificador estable Android');
assert.ok(!gradle.includes('applicationIdSuffix'),'No generar un paquete distinto en cada Beta');
assert.ok(gradle.includes('versionCode 9011'),'Versión de Android incrementada');
const script=desktop.script,start=script.indexOf('function applyFilmMembershipLists('),end=script.indexOf('function toggleSearchedFilm(',start);
assert.ok(start>=0&&end>start);
function parse(t){const lines=t.split('\n'),entries=[];
 lines.forEach((l,idx)=>{const m=l.match(/^([☐☑])\s+(.*?)(?:\s+\((\d{4})\))?$/);
  if(m)entries.push({idx,checked:m[1]==='☑',title:m[2],name:m[2],year:m[3]||''});
 });return {lines,entries};
}
const key=f=>f.title+'|'+f.year;
const apply=new Function('parse','groupKey',script.slice(start,end)+';return applyFilmMembershipLists;')(parse,key);
const lists=[
 {id:'1',name:'Primera',text:'☑ A (2024)\n☐ B (2025)'},
 {id:'2',name:'Segunda',text:'☐ A (2024)'},
 {id:'3',name:'Tercera',text:'☐ C (2018)'}
];
const ids=lists.map(x=>x.id);
assert.equal(apply(lists,{title:'A',year:'2024',checked:true},new Set(['2','3']),new Map()),2);
assert.equal(lists[0].text,'☐ B (2025)','No modifica otras películas');
assert.equal(lists[1].text,'☐ A (2024)','No duplica');
assert.equal(lists[2].text,'☐ C (2018)\n☑ A (2024)','Conserva estado visto');
assert.deepEqual(lists.map(x=>x.id),ids,'No modifica IDs');
const a=script.indexOf('function filmMembershipIndex('),b=script.indexOf('function numberedListName(',a);
const index=new Function('libraryState','parse','buildIdentityGroups','groupKey','let membershipCache=null;'+script.slice(a,b)+';return filmMembershipIndex;')(
 {lists:[{id:'a',name:'A',text:'☐ Viaje (2020)'},{id:'b',name:'B',text:'☑ Journey (2020)'}]},parse,()=>new Map(),()=>314
)();
assert.equal(index.rows.length,1,'Una ficha por identidad TMDB');
assert.deepEqual(index.rows[0]._lists.map(x=>x.n),[1,2],'Listas numeradas');
assert.equal(index.rows[0].checked,true,'Estado global visto');
assert.deepEqual(index.rows[0]._searchNames,['Viaje','Journey'],'Alias en búsqueda global');
console.log('OK: sintaxis Windows+Android, sincronización preservada, búsqueda y listas y estadísticas');
