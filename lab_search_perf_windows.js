// Windows LAB: índice de Biblioteca, metadatos de filtros y búsqueda amortiguada.
// Nunca persiste datos: se reconstruye a partir de las listas oficiales.
const labSearchDelayMs=150;
let labSearchTimer=null;
function labScheduleLibrarySearch(value){
 search=String(value||'');
 filmPage=0;
 if(labSearchTimer!==null)clearTimeout(labSearchTimer);
 labSearchTimer=setTimeout(()=>{
  labSearchTimer=null;
  renderList();
 },labSearchDelayMs);
}
function labCancelLibrarySearch(){
 if(labSearchTimer!==null)clearTimeout(labSearchTimer);
 labSearchTimer=null;
}

// Únicamente se vuelve a analizar el TXT de la lista cuyo contenido cambió.
// Las claves de caché se basan en ID y contenido, nunca en el número visible.
const labListParseCache=new Map();
let labGlobalIndexCache=null;
let labGlobalIndexOrder='';
let labGlobalIndexIdentityRevision=-1;
let labGlobalIndexGeneration=0;
let labIdentityVersion=0;
let labMetadataVersion=0;
let labFilterChoicesCache=null;
function labIndexedEntries(list){
 const cached=labListParseCache.get(list.id);
 if(cached&&cached.text===list.text)return cached.entries;
 const entries=parse(list.text||'').entries;
 labListParseCache.set(list.id,{text:list.text,entries});
 return entries;
}
function labTrackMapChanges(map,kind){
 const originalSet=map.set;
 const originalDelete=map.delete;
 const originalClear=map.clear;
 map.set=function(key,value){
  const before=this.get(key);
  if(before!==value){
   if(kind==='identity')labIdentityVersion++;
   else{
    labMetadataVersion++;
    if(Number(before?.tmdb_id||0)!==Number(value?.tmdb_id||0))labIdentityVersion++;
   }
  }
  return originalSet.call(this,key,value);
 };
 map.delete=function(key){
  if(this.has(key)){
   if(kind==='identity')labIdentityVersion++;
   else{
    labMetadataVersion++;
    if(Number(this.get(key)?.tmdb_id||0)>0)labIdentityVersion++;
   }
  }
  return originalDelete.call(this,key);
 };
 map.clear=function(){
  if(this.size){
   if(kind==='identity')labIdentityVersion++;
   else{labIdentityVersion++;labMetadataVersion++;}
  }
  return originalClear.call(this);
 };
}
labTrackMapChanges(identityMap,'identity');
labTrackMapChanges(platformData,'metadata');

function labGlobalMembershipIndex(){
 const ids=[];
 let changed=!labGlobalIndexCache;
 for(const list of libraryState.lists){
  ids.push(list.id);
  const cached=labListParseCache.get(list.id);
  if(!cached||cached.text!==list.text){
   labIndexedEntries(list);
   changed=true;
  }
 }
 const order=ids.join('\u001f');
 if(order!==labGlobalIndexOrder)changed=true;
 const currentIDs=new Set(ids);
 for(const key of labListParseCache.keys()){
  if(!currentIDs.has(key)){labListParseCache.delete(key);changed=true;}
 }
 if(labGlobalIndexIdentityRevision!==labIdentityVersion)changed=true;
 // Renombrar una lista no altera su TXT, pero sí las etiquetas de pertenencia.
 if(!changed&&labGlobalIndexCache){
  const names=labGlobalIndexCache.listNames;
  changed=libraryState.lists.some(l=>names.get(l.id)!==l.name);
 }
 if(!changed){
  membershipCache=labGlobalIndexCache;
  return labGlobalIndexCache;
 }
 const groups=buildIdentityGroups(),byKey=new Map();
 const membershipIds=new Map();
 for(const [index,list] of libraryState.lists.entries()){
  for(const film of labIndexedEntries(list)){
   const key=groupKey(film,groups);
   let record=byKey.get(key);
   if(!record){
    record={...film,_global:true,_group:key,_lists:[],_searchNames:[]};
    byKey.set(key,record);
    membershipIds.set(key,new Set());
   }
   record.checked=record.checked||film.checked;
   record._searchNames.push(film.name);
   const present=membershipIds.get(key);
   if(!present.has(list.id)){
    record._lists.push({id:list.id,n:index+1,name:list.name,title:film.title,year:film.year});
    present.add(list.id);
   }
  }
 }
 labGlobalIndexCache={
  rows:[...byKey.values()],byKey,groups,
  listNames:new Map(libraryState.lists.map(l=>[l.id,l.name]))
 };
 labGlobalIndexOrder=order;
 labGlobalIndexIdentityRevision=labIdentityVersion;
 labGlobalIndexGeneration++;
 membershipCache=labGlobalIndexCache;
 return labGlobalIndexCache;
}

// Crear las opciones de género y década solo si cambia el conjunto de
// películas, los metadatos, la lista activa o el tipo de búsqueda.
function labShouldRefreshFilterChoices(){
 const global=searchScope==='all';
 const scope=global?'all':(typeof labFavoritesMode!=='undefined'&&labFavoritesMode?'favorites':virtualMode?'comparison':'current');
 let dataset=parsed;
 if(global){
  dataset=filmMembershipIndex();
 }
 const prior=labFilterChoicesCache;
 if(prior&&prior.scope===scope&&prior.dataset===dataset&&
   prior.activeID===libraryState.active_id&&prior.metadataVersion===labMetadataVersion){
  return false;
 }
 labFilterChoicesCache={scope,dataset,activeID:libraryState.active_id,metadataVersion:labMetadataVersion};
 return true;
}
