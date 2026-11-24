// Executed inside the built page, in WebKit and WebView2.
(() => {
  const el=id=>document.getElementById(id);
  const check=(value,message)=>{if(!value)throw Error(message)};
  check(document.title.includes('9.2.0'),'Versión incorrecta');
  for(const id of ['labDiscoverNav','labFilterOpen','labSmartEnabled','labMaintenanceRefresh','labFavoritesCollection']){
    check(el(id),'Falta '+id);
  }
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
