let activityItems=[],activityPage=0,activityLoading=false;
const activityLabels={seen:'Vista',unseen:'Pendiente',film_add:'Película añadida',film_remove:'Película retirada',film_edit:'Título modificado',list_create:'Lista creada',list_delete:'Lista eliminada',list_rename:'Lista renombrada',favorite_add:'Favorito añadido',favorite_remove:'Favorito retirado',rating:'Valoración personal',identity:'Identificación TMDB',order:'Orden modificado',comparison:'Comparación modificada',restore:'Copia restaurada',bulk:'Modificación masiva'};
function renderActivity(){
 const q=$('activitySearch').value.normalize('NFD').replace(/[\u0300-\u036f]/g,'').toLowerCase();
 const type=$('activityFilter').value;
 const matches=activityItems.filter(e=>(!type||e.kind===type)&&[e.film,e.list,e.detail,activityLabels[e.kind]].join(' ').normalize('NFD').replace(/[\u0300-\u036f]/g,'').toLowerCase().includes(q));
 const pages=Math.max(1,Math.ceil(matches.length/20));activityPage=Math.min(activityPage,pages-1);
 const target=$('cloudHistory');target.replaceChildren();target.classList.toggle('is-empty',!matches.length);
 if(!matches.length)target.textContent='No hay movimientos que coincidan.';
 for(const e of matches.slice(activityPage*20,activityPage*20+20)){
  const row=document.createElement('div');row.className='cine-sync-history-row';
  const title=document.createElement('strong');title.textContent=[activityLabels[e.kind]||e.kind,e.film].filter(Boolean).join(' · ');
  const info=document.createElement('span');info.textContent=[e.list,e.detail,cineSyncAt(e.at)].filter(Boolean).join(' · ');
  row.append(title,info);target.append(row);
 }
 $('activityPage').textContent=(activityPage+1)+' / '+pages+' · '+matches.length+' movimientos';
 $('activityPrev').disabled=activityPage===0;$('activityNext').disabled=activityPage+1>=pages;
}
async function cineRefreshHistory(){
 if(activityLoading)return;
 activityLoading=true;const button=$('cloudHistoryRefresh');button.disabled=true;
 try{const data=await cloudRequest('/api/sync/history');activityItems=data.events||[];renderActivity();$('cloudHistoryChecked').textContent='Historial actualizado · '+activityItems.length+' movimientos conservados (máximo 10.000).';}
 catch(e){$('cloudHistoryChecked').textContent='⚠ '+e.message;}finally{button.disabled=false;activityLoading=false;}
}
$('activitySearch').addEventListener('input',()=>{activityPage=0;renderActivity();});
$('activityFilter').addEventListener('change',()=>{activityPage=0;renderActivity();});
$('activityPrev').addEventListener('click',()=>{activityPage--;renderActivity();});
$('activityNext').addEventListener('click',()=>{activityPage++;renderActivity();});
for(const [value,label] of Object.entries(activityLabels))$('activityFilter').append(new Option(label,value));
