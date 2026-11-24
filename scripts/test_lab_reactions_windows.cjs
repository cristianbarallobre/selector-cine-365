'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');

const template=fs.readFileSync('lab_reactions_windows.html','utf8');
const source=template.split('<!-- LAB_REACTIONS_SCRIPT -->')[1];
const page=fs.readFileSync('lab_discover_windows.html','utf8');
const injections=fs.readFileSync('lab_reactions_windows.go','utf8');
assert.ok(source);
new vm.Script('(()=>{\n'+source+'\n})();',{filename:'lab_reactions_windows.html'});
assert.ok(template.includes('😫'),'First emotion');
for(const emoji of ['😫','☹️','😐','🙂','😄'])assert.ok(template.includes(emoji),'Missing '+emoji);
assert.match(template,/id="labRatingOverlay"/);
assert.match(injections,/labFavoritesMode\|\|virtualMode&&!global/);
assert.match(injections,/labFilmRowActions\(f,pill,labCreateDetailButton\(f\),fix\)/);
assert.match(template,/#view-library \.lab-film-action-layout\{[\s\S]*?grid-template-columns:72px 122px/);
assert.match(template,/#view-library \.lab-film-action-grid\{[\s\S]*?grid-template-rows:31px 31px/);
assert.match(template,/#view-library \.lab-film-action-grid \.pill\{[\s\S]*?display:flex!important/);
assert.match(template,/@media\(max-width:650px\)/);
assert.match(injections,/labFavoritesFilms\(\)/);
assert.match(injections,/labFavoritesOpen/);
assert.match(injections,/const favoriteBadge=labFavoriteMembershipButton\(f\)/,'El indicador se inserta junto a los números');
assert.match(injections,/if\(favoriteBadge\)membership\.appendChild\(favoriteBadge\)/,'No mostrar corazón en películas que no son favoritas');
assert.match(template,/#view-library \.film-memberships \.lab-membership-heart/,'Corazón discreto en el bloque de listas');
assert.match(page,/labDiscoverFavoriteButton\(film\)/);
assert.match(page,/id="labDetailFavorite"/);

class FakeNode{
 constructor(tag='div'){
  this.tagName=tag;this.children=[];this.events={};this.attrs={};this.className='';this.style={};
  this.hidden=true;this.textContent='';this.value='';this.isConnected=true;this.dataset={};this.classSet=new Set();
  this.classList={
   add:k=>this.classSet.add(k),
   remove:k=>this.classSet.delete(k),
   toggle:(k,v)=>v?this.classSet.add(k):this.classSet.delete(k),
   contains:k=>this.classSet.has(k)
  };
 }
 append(...args){this.children.push(...args);}
 appendChild(node){this.children.push(node);}
 replaceChildren(...args){this.children=[...args];}
 addEventListener(name,cb){this.events[name]=cb;}
 setAttribute(key,value){this.attrs[key]=value;}
 querySelector(s){
  if(s==='button')return this.children.find(x=>x.tagName==='button')||null;
  if(s==='button.lab-rating-selected')return this.children.find(x=>x.className==='lab-rating-selected')||null;
  if(s==='.lab-heart-button'||s==='.lab-rate-button')return this.children.find(x=>x.className.split(' ').includes(s.slice(1)))||null;
  return null;
 }
 focus(){doc.activeElement=this;}
 click(){this.events.click?.({stopPropagation(){}});}
}
const ids=['labRatingOverlay','labRatingTitle','labRatingOptions','labRatingClear','labRatingClose','labFavoritesCollection','labFavoritesCount','labFavoritesOpen','labDetailFavorite','filmList'];
const dom=Object.fromEntries(ids.map(x=>[x,new FakeNode()]));
let renderCalls=0,collections=0,saves=0,shown='';
const doc={activeElement:new FakeNode('button'),createElement:tag=>new FakeNode(tag),
 querySelectorAll:()=>[],documentElement:{clientWidth:950,clientHeight:720}};
const parse=text=>{
 const entries=text.split('\n').filter(Boolean).map((row,idx)=>{
  const m=/^[☐☑] (.+?)(?: \((\d{4})\))?$/.exec(row);
  return m?{idx,checked:row[0]==='☑',title:m[1],year:m[2]||'',name:m[1]}:null;
 }).filter(Boolean);
 return {entries,seen:entries.filter(x=>x.checked).length};
};
const libraryState={
 active_id:'lista1',
 lists:[{id:'lista1',text:'☑ Una película (2024)\n☐ Otra (2025)',name:'Primera'},
        {id:'lista2',text:'☑ Una película (2024)',name:'Segunda'}]
};
const ctx={
 $:id=>{assert.ok(dom[id],id);return dom[id];},
 window:{addEventListener(){}},document:doc,libraryState,
 filmKey:f=>f.title.toLowerCase()+'|'+(Number(f.year)||0),
 knownTMDB:()=>0,parse,Date,Math,Number,String,Array,Object,
 labCurrentMovie:null,
 scheduleSave:()=>{saves++},
 renderList:()=>{renderCalls++},
 renderCollectionControls:()=>{collections++},
 notify:()=>{},
 syncActiveText:()=>{},
 activeCollection:()=>libraryState.lists[0],
 restartList:()=>{},
 showView:x=>{shown=x},
 virtualMode:false,rawText:'',
};
vm.createContext(ctx);
vm.runInContext(source,ctx);
const evaluate=js=>vm.runInContext(js,ctx);
const watched={title:'Una película',year:'2024',checked:true};
const pending={title:'Otra',year:'2025',checked:false};
assert.equal(evaluate('labFavoriteCount()'),0,'No mostrar lista sin favoritos');
assert.equal(dom.labFavoritesCollection.hidden,true);
const watchedControls=ctx.labReactionControls(watched);
assert.equal(watchedControls.children.length,2,'Vistas: corazón y valoración');
assert.equal(watchedControls.children[0].textContent,'♡');
assert.equal(watchedControls.children[1].children[1].textContent,'Valorar');
const pendingControls=ctx.labReactionControls(pending);
assert.equal(pendingControls.children.length,1,'Pendientes: solo corazón');
assert.equal(ctx.labFavoriteMembershipButton(watched),null,'Sin favorito no se muestra corazón junto a Listas');
assert.equal(ctx.labFavoriteMembershipButton(pending),null,'Sin favorito no aparece indicador de lista automática');

// La valoración ocupa una columna a la izquierda. La derecha presenta
// estado/corazón arriba y Detalles/edición debajo, con el mismo ancho en todas.
const seenStatus=new FakeNode('span');seenStatus.textContent='Vista';
const seenDetails=new FakeNode('button');seenDetails.textContent='Detalles';
const seenEdit=new FakeNode('button');seenEdit.textContent='✎';
const seenRow=ctx.labFilmRowActions(watched,seenStatus,seenDetails,seenEdit);
assert.equal(seenRow.className,'lab-film-action-layout');
assert.equal(seenRow.children.length,2);
const [seenScore,seenRight]=seenRow.children;
assert.equal(seenScore.className,'lab-film-score-slot');
assert.equal(seenScore.children.length,1,'Las vistas tienen valoración a la izquierda');
assert.equal(seenScore.children[0].className,'lab-rate-button');
assert.equal(seenRight.className,'lab-film-action-grid');
assert.equal(seenRight.children.length,4,'Dos filas y dos columnas de acciones');
assert.equal(seenRight.children[0],seenStatus,'Primera fila, columna izquierda: estado');
assert.equal(seenRight.children[1].className,'lab-heart-button','Primera fila, columna derecha: corazón');
assert.equal(seenRight.children[2],seenDetails,'Segunda fila, columna izquierda: Detalles');
assert.equal(seenRight.children[3],seenEdit,'Segunda fila, columna derecha: edición');

const pendingStatus=new FakeNode('span');pendingStatus.textContent='Pendiente';
const pendingDetails=new FakeNode('button');pendingDetails.textContent='Detalles';
const pendingEdit=new FakeNode('button');pendingEdit.textContent='✎';
const pendingRow=ctx.labFilmRowActions(pending,pendingStatus,pendingDetails,pendingEdit);
assert.equal(pendingRow.className,'lab-film-action-layout');
assert.equal(pendingRow.children[0].className,seenScore.className);
assert.equal(pendingRow.children[0].children.length,0,'Pendientes reservan espacio, sin permitir valorar');
assert.equal(pendingRow.children[1].className,seenRight.className);
assert.equal(pendingRow.children[1].children[0],pendingStatus);
assert.equal(pendingRow.children[1].children[1].className,'lab-heart-button');
assert.equal(pendingRow.children[1].children[2],pendingDetails);
assert.equal(pendingRow.children[1].children[3],pendingEdit);
watchedControls.children[0].click();
assert.equal(libraryState.reactions['una película|2024'].favorite,true);
assert.equal(evaluate('labFavoriteCount()'),1);
assert.equal(dom.labFavoritesCollection.hidden,false);
assert.equal(dom.labFavoritesCount.textContent,'1 película · lista automática');
assert.equal(saves,1,'Se guarda inmediatamente la elección');
const smallFavoriteHeart=ctx.labFavoriteMembershipButton(watched);
assert.ok(smallFavoriteHeart,'Al marcar una película como favorita aparece el corazón junto a Listas');
assert.equal(smallFavoriteHeart.tagName,'button');
assert.equal(smallFavoriteHeart.className,'lab-membership-heart');
assert.equal(smallFavoriteHeart.textContent,'♥');
assert.equal(smallFavoriteHeart.title,'En Favoritos');
assert.equal(smallFavoriteHeart.attrs['aria-pressed'],'true');
assert.equal(smallFavoriteHeart.attrs['aria-label'],'Quitar «Una película» de Favoritos');
ctx.labOpenFavorites();
assert.equal(evaluate('labFavoritesMode'),true);
assert.equal(shown,'library');
const favorites=ctx.labFavoritesFilms();
assert.equal(favorites.length,1,'Sin duplicar película de dos listas');
assert.equal(favorites[0].checked,true,'Mantiene estado de vista de la lista original');

ctx.labOpenRating(pending);
assert.equal(dom.labRatingOverlay.hidden,true,'No puntuar pendientes');
ctx.labOpenRating(watched);
assert.equal(dom.labRatingOverlay.hidden,false);
assert.equal(dom.labRatingOptions.children.length,5);
assert.equal(dom.labRatingOptions.children[0].attrs['aria-label'],'1 de 5: Muy mala');
dom.labRatingOptions.children[3].click();
assert.equal(libraryState.reactions['una película|2024'].rating,4);
assert.equal(libraryState.reactions['una película|2024'].favorite,true,'Puntuar no quita corazón');
assert.equal(dom.labRatingOverlay.hidden,true);
assert.equal(saves,2);
ctx.labOpenRating(watched);
assert.equal(dom.labRatingClear.hidden,false);
dom.labRatingClear.click();
assert.equal(libraryState.reactions['una película|2024'].rating,0,'La nota se puede quitar');
assert.equal(libraryState.reactions['una película|2024'].favorite,true,'Quitar nota no quita favorito');

ctx.labToggleFavorite(watched);
assert.equal(libraryState.reactions['una película|2024'].favorite,false);
assert.equal(evaluate('labFavoriteCount()'),0,'La lista desaparece cuando ya no hay corazones');
assert.equal(evaluate('labFavoritesMode'),false);
assert.equal(dom.labFavoritesCollection.hidden,true);
assert.equal(libraryState.lists.length,2,'No se duplica ninguna lista persistente');
assert.equal(libraryState.lists[0].text,'☑ Una película (2024)\n☐ Otra (2025)','El TXT original no se altera');
assert.ok(renderCalls>=3&&collections>=3);
const discovered={id:123,title:'Película descubierta',release_date:'2026-10-01'};
const button=ctx.labDiscoverFavoriteButton(discovered);
button.click();
const rec=libraryState.reactions['película descubierta|2026'];
assert.equal(rec.tmdb_id,123,'Conserva ID real de TMDB al añadir desde Descubrir');
assert.equal(rec.favorite,true);
assert.equal(ctx.labFavoritesFilms().length,1,'También incluye favoritos descubiertos fuera de la Biblioteca');
const discoveredBadge=ctx.labFavoriteMembershipButton({title:'Película descubierta',year:'2026'});
assert.equal(discoveredBadge.textContent,'♥','Comparte favorito con botón grande y Descubrir');
let stopped=false;
const savedBefore=saves;
discoveredBadge.events.click({stopPropagation(){stopped=true;}});
assert.equal(stopped,true,'El clic en el corazón no debe activar otro control de la fila');
assert.equal(libraryState.reactions['película descubierta|2026'].favorite,false,'El pequeño corazón permite quitar de Favoritos');
assert.equal(saves,savedBefore+1,'Quitar desde Listas guarda una sola vez');
assert.equal(ctx.labFavoriteMembershipButton({title:'Película descubierta',year:'2026'}),null,'El corazón desaparece al quitar Favoritos');
assert.equal(evaluate('labFavoriteCount()'),0,'La lista virtual desaparece sin favoritos');
assert.equal(libraryState.lists.length,2,'No crea ni borra listas normales');
console.log('OK: corazón junto a Listas solo si es favorito, clic para quitar, botones alineados y sin duplicar datos');
