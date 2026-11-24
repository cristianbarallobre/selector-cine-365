'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');

const stable=fs.readFileSync('plantilla.html','utf8');
const injector=fs.readFileSync('lab_reactions_windows.go','utf8');
const reactions=fs.readFileSync('lab_reactions_windows.html','utf8');
const filters=fs.readFileSync('lab_filters_inject_windows.go','utf8');
const rx=reactions.split('<!-- LAB_REACTIONS_SCRIPT -->')[1];
assert.ok(rx,'El JS de Favoritos debe existir');
new vm.Script('(()=>{\n'+rx+'\n})();');
for(const needle of [
 "const labStats=labStatsDataset();const n=labStats.entries.length;",
 "'watched').textContent=labStats.seen;",
 "'remaining').textContent=labStats.pending.length;",
 "labStatsDataset().pending.filter(f=>{",
 "labStatsDataset().pending.length;",
 "const seen=labStatsDataset().entries.filter(f=>f.checked)",
 "labFavoritesMode?'♥ Favoritos'",
 "$('markWatched').disabled=virtualMode||labFavoritesMode;",
 "function toggle(idx){\\n if(virtualMode||labFavoritesMode)return;",
 "labStatsDataset().entries;",
 "labRefreshStatsPanels();"
])assert.ok(injector.includes(needle),'Falta integración en HTML de LAB: '+needle);
assert.ok(filters.includes("labStatsHeroValues(seen,known,minutes)"),'Las estadísticas extendidas siguen activas');
assert.match(injector,/function renderList\(\)\{labFiltersRefreshChoices\(\);if\(labFavoritesMode&&!labFavoriteCount\(\)\)\{labFavoritesMode=false;labRefreshStatsPanels\(\);\}/);
assert.doesNotMatch(stable,/labStatsDataset\(/,'No alterar la versión estable');

function slice(src,first,after){
 const a=src.indexOf(first);assert.ok(a>=0,'No se encuentra '+first);
 const b=src.indexOf(after,a);assert.ok(b>a,'No se encuentra cierre de '+first);
 return src.slice(a,b);
}
function replaceOnce(src,old,next){
 assert.equal(src.split(old).length,2,'Esperaba un único fragmento: '+old);
 return src.replace(old,next);
}
const intro=slice(rx,"const labFavoriteID=","function labWriteReaction(f,changes){");
let update=slice(stable,"function update(){","// Dragging is enabled");
update=replaceOnce(update,
 "function update(){membershipCache=null;parsed=parse(rawText);const n=parsed.entries.length;",
 "function update(){membershipCache=null;parsed=parse(rawText);const labStats=labStatsDataset();const n=labStats.entries.length;");
update=replaceOnce(update,
 "$('total').textContent=n;$('watched').textContent=parsed.seen;$('remaining').textContent=parsed.pending.length;",
 "$('total').textContent=n;$('watched').textContent=labStats.seen;$('remaining').textContent=labStats.pending.length;");
update=replaceOnce(update,
 "const pct=n?Math.round(100*parsed.seen/n):0;",
 "const pct=n?Math.round(100*labStats.seen/n):0;");
let lottery=slice(stable,"function eligibleFilms(){","function updatePlatformChoices(){");
lottery=replaceOnce(lottery,"if(selectedSortPlatform==='all')return parsed.pending;","if(selectedSortPlatform==='all')return labStatsDataset().pending;");
lottery=replaceOnce(lottery,"return parsed.pending.filter(f=>{","return labStatsDataset().pending.filter(f=>{");
lottery=replaceOnce(lottery,"const count=eligibleFilms().length;const total=parsed.pending.length;","const count=eligibleFilms().length;const total=labStatsDataset().pending.length;");
let statistics=slice(stable,"function renderStats(){","\n\n$('filmPrev')");
statistics=replaceOnce(statistics,
 "if(!parsed)return;const seen=parsed.entries.filter(f=>f.checked),hero=$('statsHero')",
 "if(!parsed)return;const seen=labStatsDataset().entries.filter(f=>f.checked),hero=$('statsHero')");
statistics=replaceOnce(statistics,
 "$('statsTitle').textContent=virtualMode?'Películas coincidentes · solo lectura':activeCollection()?.name||'';",
 "$('statsTitle').textContent=labFavoritesMode?'♥ Favoritos':virtualMode?'Películas coincidentes · solo lectura':activeCollection()?.name||'';");
statistics=replaceOnce(statistics,
 "for(const [label,val] of [['Películas vistas',seen.length],['Pendientes',parsed.pending.length],['Horas registradas',minutes?Math.floor(minutes/60)+' h '+(minutes%60)+' min':'—'],['Fichas TMDB',known.filter(x=>x.metadata_version>=2).length]]){",
 "for(const [label,val] of labStatsHeroValues(seen,known,minutes)){");

class Element{
 constructor(){this.textContent='';this.children=[];this.style={};this.attrs={};this.disabled=false;}
 replaceChildren(){this.children=[];}
 append(...children){this.children.push(...children);}
 appendChild(node){this.children.push(node);}
 setAttribute(k,v){this.attrs[k]=v;}
}
const ids=['total','watched','remaining','progressFill','progress','libraryCount','choose',
 'sortCount','statsTitle','statsHero','statsSections'];
const ui=Object.fromEntries(ids.map(id=>[id,new Element()]));
const parse=text=>{
 const entries=[];
 for(const row of String(text||'').split('\n')){
  const m=/^([☐☑]) (.+?)(?: \((\d{4})\))?$/.exec(row);
  if(m)entries.push({idx:entries.length,title:m[2],name:m[2],year:m[3]||'',checked:m[1]==='☑'});
 }
 const pending=entries.filter(f=>!f.checked);
 return {entries,pending,seen:entries.length-pending.length};
};
const lists=[
 {id:'main',name:'Quiero ver',text:'☑ Visto original (2020)\n☐ Pendiente (2021)\n☐ Solo normal (2019)'},
 {id:'other',name:'Otra lista',text:'☑ Visto original (2020)\n☐ Otro pendiente (2022)'}
];
const state={lists,active_id:'main',reactions:{
 'visto original|2020':{favorite:true,title:'Visto original',year:'2020',changed_at:11,rating:5},
 'pendiente|2021':{favorite:true,title:'Pendiente',year:'2021',changed_at:12},
 'fuera de listas|2025':{favorite:true,title:'Fuera de listas',year:'2025',changed_at:13}
}};
const platformData=new Map([
 ['visto original|2020',{runtime:120,genres:['Ciencia ficción'],countries:['EE. UU.'],director:'Directora A',metadata_version:2}],
 ['pendiente|2021',{runtime:100,providers:['Filmin'],state:'ok'}],
 ['solo normal|2019',{runtime:105,providers:['Filmin'],state:'ok'}]
]);
const ctx={
 $:id=>{assert.ok(ui[id],id);return ui[id];},
 document:{createElement:()=>new Element()},
 Map,Set,Array,String,Number,Object,Math,Date,parse,
 libraryState:state,parsed:parse(lists[0].text),rawText:lists[0].text,
 membershipCache:null,virtualMode:false,selectedSortPlatform:'all',
 platformData, chosen:null, statPages:new Map(),
 activeCollection:()=>lists.find(x=>x.id===state.active_id),
 providerKey:f=>String(f.title).toLocaleLowerCase('es')+'|'+f.year,
 filmKey:f=>String(f.title).toLocaleLowerCase('es')+'|'+f.year,
 subscriptions:r=>r.providers||[],
 renderList:()=>{},updateSortAvailability:undefined,update:undefined,
 updateChoicePlatforms:()=>{},resetChoice:()=>{ctx.chosen=null;},
 statRows:(heading,entries)=>{const x=new Element();x.title=heading;x.entries=entries;return x;},
 labStatsHeroValues:(seen,known,minutes)=>[
  ['Películas vistas',seen.length],['Horas registradas',minutes?Math.floor(minutes/60)+' h':'—'],
  ['Nota media TMDB','—']
 ],
};
vm.createContext(ctx);
vm.runInContext(intro+'\n'+lottery+'\n'+statistics+'\n'+update,ctx);
let x=code=>vm.runInContext(code,ctx);
ctx.update();
ctx.renderStats();
assert.equal(ui.total.textContent,3);
assert.equal(ui.watched.textContent,1);
assert.equal(ui.remaining.textContent,2);
assert.equal(ui.sortCount.textContent,'2 pendientes');
assert.equal(ui.progressFill.style.width,'33%');
assert.equal(ui.statsTitle.textContent,'Quiero ver');
assert.equal(ui.statsHero.children[0].children[0].textContent,1);
assert.equal(ctx.eligibleFilms().length,2);

// Favoritos no debe heredar 3 títulos ni el porcentaje de «Quiero ver».
x('labFavoritesMode=true');
ctx.update();
ctx.renderStats();
assert.equal(ui.total.textContent,3);
assert.equal(ui.watched.textContent,1);
assert.equal(ui.remaining.textContent,2);
assert.equal(ui.sortCount.textContent,'2 pendientes');
assert.equal(ui.progressFill.style.width,'33%');
assert.equal(ui.statsTitle.textContent,'♥ Favoritos');
assert.equal(ctx.eligibleFilms().length,2);
assert.equal(ui.statsHero.children[0].children[0].textContent,1);
assert.equal(ui.statsSections.children.length,4,'Actualizar gráficas de década, género, país y dirección');
ctx.selectedSortPlatform='Filmin';
ctx.updateSortAvailability();
assert.equal(ctx.eligibleFilms().length,1,'El sorteo usa solo las pendientes favoritas disponibles');
assert.equal(ui.sortCount.textContent,'1 de 2 pendientes disponibles');

state.reactions['pendiente|2021'].favorite=false;
ctx.labRefreshStatsPanels();
assert.equal(ui.total.textContent,2);
assert.equal(ui.watched.textContent,1);
assert.equal(ui.remaining.textContent,1);
assert.equal(ui.sortCount.textContent,'0 de 1 pendientes disponibles');
assert.equal(ui.choose.disabled,true);
assert.equal(ui.progressFill.style.width,'50%');
assert.equal(ui.statsTitle.textContent,'♥ Favoritos');

state.reactions['fuera de listas|2025'].favorite=false;
ctx.selectedSortPlatform='all';
ctx.labRefreshStatsPanels();
assert.equal(ui.total.textContent,1);
assert.equal(ui.watched.textContent,1);
assert.equal(ui.remaining.textContent,0);
assert.equal(ui.sortCount.textContent,'0 pendientes');
assert.equal(ui.progressFill.style.width,'100%');
assert.equal(ui.choose.disabled,true);
assert.equal(ctx.labStatsDataset().entries[0].title,'Visto original');
assert.equal(ui.statsHero.children[0].children[0].textContent,1);

// Volver a la última lista normal restituye exclusivamente sus cifras.
x('labFavoritesMode=false');
ctx.update();
ctx.renderStats();
assert.equal(ui.total.textContent,3);
assert.equal(ui.watched.textContent,1);
assert.equal(ui.remaining.textContent,2);
assert.equal(ui.statsTitle.textContent,'Quiero ver');
assert.equal(ui.sortCount.textContent,'2 pendientes');
assert.equal(lists[0].text,'☑ Visto original (2020)\n☐ Pendiente (2021)\n☐ Solo normal (2019)');
assert.equal(lists[1].text,'☑ Visto original (2020)\n☐ Otro pendiente (2022)');
console.log('OK: Favoritos tiene estadísticas, progreso y sorteo propios; al salir se restaura la lista real sin modificar TXT');
