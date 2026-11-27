let cloudManualBusy=false,extraBackupDirty=false,extraBackupBusy=false;
function extraBackupShow(s){
 const input=$('extraBackupFolder'),status=$('extraBackupStatus');
 if(!extraBackupDirty&&document.activeElement!==input)input.value=s.folder||'';
 $('extraBackupNow').disabled=!s.enabled||s.busy||extraBackupBusy;
 $('extraBackupDisable').disabled=!s.enabled||extraBackupBusy;
 $('extraBackupTarget').textContent=s.target?'Archivo: '+s.target:'';
 status.classList.toggle('warn',Boolean(s.error));
 status.textContent=!s.enabled?'Solo guardado principal en este equipo.':
   s.error?'⚠ '+s.error:
   s.busy?'⟳ Creando copia adicional…':
   s.pending?'◷ Copia adicional pendiente…':
   s.last_backup?'✓ Copia adicional guardada · '+new Date(s.last_backup).toLocaleString('es-ES'):'Copia adicional configurada.';
}
async function cloudRefresh(){
 try{extraBackupShow(await cloudRequest('/api/extra-backup/status'));}
 catch(e){$('extraBackupStatus').textContent='⚠ No se pudo consultar la copia adicional.';}
}
async function extraBackupAction(action,payload){
 if(extraBackupBusy)return;
 extraBackupBusy=true;
 const buttons=['extraBackupBrowse','extraBackupActivate','extraBackupNow','extraBackupDisable'].map($);
 buttons.forEach(b=>b.disabled=true);
 try{
  const result=await cloudRequest('/api/extra-backup/'+action,'POST',payload);
  if(action==='browse'){
   if(result.folder){$('extraBackupFolder').value=result.folder;extraBackupDirty=true;$('extraBackupStatus').textContent='Carpeta elegida. Pulsa Guardar ubicación.';}
  }else{
   if(action==='config'){extraBackupDirty=false;$('extraBackupFolder').value=result.folder||'';}
   extraBackupShow(result);
  }
 }catch(e){$('extraBackupStatus').textContent='⚠ '+e.message;}
 finally{
  extraBackupBusy=false;buttons[0].disabled=false;buttons[1].disabled=false;
  if(action!=='browse')void cloudRefresh();else{buttons[2].disabled=false;buttons[3].disabled=false;}
 }
}
$('extraBackupFolder').addEventListener('input',()=>{extraBackupDirty=true;});
$('extraBackupBrowse').addEventListener('click',()=>void extraBackupAction('browse'));
$('extraBackupActivate').addEventListener('click',()=>{
 const folder=$('extraBackupFolder').value.trim();
 if(!folder){$('extraBackupStatus').textContent='Selecciona una carpeta primero.';return;}
 void extraBackupAction('config',{folder});
});
$('extraBackupNow').addEventListener('click',()=>void extraBackupAction('now'));
$('extraBackupDisable').addEventListener('click',()=>{
 if(confirm('¿Desactivar la copia adicional? Los archivos ya guardados se conservarán.'))void extraBackupAction('config',{folder:''});
});
setInterval(()=>{if(!document.hidden&&activeView==='collections'&&!extraBackupBusy)void cloudRefresh();},2500);
 const cineRefreshStalePlatforms=()=>{if(!document.hidden&&libraryReady&&platformConfigured&&!platformStatusUnavailable&&!platformActive&&activeView==='library'&&parsed.entries.some(f=>{const e=platformData.get(providerKey(f));return !!e&&cinePlatformExpired(e);}))void loadPlatforms(false,true);};
 setInterval(()=>{cineUpdatePlatformStamps();cineRefreshStalePlatforms();},60000);
 setTimeout(cineRefreshStalePlatforms,5000);
