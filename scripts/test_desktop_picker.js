// Executed inside the built page, in WebKit and WebView2.
(() => {
  const el=id=>document.getElementById(id);
  const check=(value,message)=>{if(!value)throw Error(message)};
  check(document.title.endsWith('9.2.11'),'Versión incorrecta');
  check(typeof window.cinePrepareClose==='function','Falta el guardado previo al cierre');
  for(const id of ['labDiscoverNav','labFilterOpen','labSmartEnabled','labMaintenanceRefresh','labFavoritesCollection']){
    check(el(id),'Falta '+id);
  }
  const picker=el('listSwitcherButton').getBoundingClientRect();
  const sortPanel=document.querySelector('#view-library section[aria-label="Sorteo de películas"]').getBoundingClientRect();
  check(Math.abs(picker.left-sortPanel.left)<1 && Math.abs(picker.width-sortPanel.width)<1,'Selector de listas desalineado');
  check(el('membershipConfirm')?.hidden,'Falta la confirmación integrada de listas');
  check(el('matchStatus'),'Falta el aviso de identificación TMDB');
  const menu=el('sortPlatformMenu'),button=el('sortPlatformButton');
  if(menu.hidden)button.click();
  const option=name=>el('sortPlatformOptions').querySelector('[data-value="'+name+'"]');
  const chosen=()=>Array.from(el('platformSelect').selectedOptions,o=>o.value).sort().join(',');
  option('Netflix').click();option('Filmin').click();
  check(!menu.hidden && chosen()==='Filmin,Netflix','La segunda selección sustituyó a la primera');
  check(option('Netflix').getAttribute('aria-selected')==='true','Netflix sin marcar');
  check(option('Filmin').getAttribute('aria-selected')==='true','Filmin sin marcar');
  option('Netflix').click();check(chosen()==='Filmin','Desmarcar afectó a la otra plataforma');
  option('any').click();check(chosen()==='any','Streaming debe ser exclusivo');
  el('sortPlatformReset').click();check(chosen()==='all','Restablecer todas');
  return true;
})()
