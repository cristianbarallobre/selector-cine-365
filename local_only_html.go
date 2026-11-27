package main

import (
 "fmt"
 "strings"
)

// Produce the local-only desktop surface after the existing desktop panels
// have been injected. Retain the local recovery and history tools.
func localOnlyHTML(page string)(string,error){
 var err error
 cut:=func(start,end,replacement string){
  if err!=nil{return}
  i:=strings.Index(page,start)
  if i<0{err=fmt.Errorf("no se pudo preparar la interfaz local: %s",start);return}
  j:=strings.Index(page[i+len(start):],end)
  if j<0{err=fmt.Errorf("no se pudo preparar la interfaz local: %s",end);return}
  page=page[:i]+replacement+page[i+len(start)+j:]
 }
 cut("    <div id=\"cloudHeaderControls\" hidden>","  </header>","")
 cut("      <div class=\"cloud-box\" aria-label=\"Configuración de sincronización y copias de seguridad\">","        <details class=\"cloud-guide-extra\">\n          <summary>Actividad reciente, historial y copias de recuperación</summary>",`      <div class="cloud-box" aria-label="Guardado y copia adicional">
        <div class="cloud-title">📁 Guardado y copia adicional</div>
        <p>La biblioteca se guarda siempre en este equipo. Si quieres, conserva una segunda copia automática en una carpeta de tu elección.</p>
        <p class="cloud-status">La copia incluye listas, estados de visionado, favoritos, valoraciones e identificaciones de películas. No incluye claves de acceso, sesiones ni cachés.</p>
        <label for="extraBackupFolder">Ubicación adicional (opcional)</label>
        <input id="extraBackupFolder" class="cloud-path" type="text" autocomplete="off" placeholder="Elige una carpeta del equipo">
        <div class="cloud-buttons">
          <button type="button" class="tool" id="extraBackupBrowse">📁 Elegir carpeta</button>
          <button type="button" class="tool" id="extraBackupActivate">Guardar ubicación</button>
          <button type="button" class="tool" id="extraBackupNow">Copiar ahora</button>
          <button type="button" class="tool" id="extraBackupDisable">Desactivar copia</button>
        </div>
        <p id="extraBackupStatus" class="cloud-status" role="status" aria-live="polite">Comprobando ubicación…</p>
        <p id="extraBackupTarget" class="cloud-status"></p>
        <p>Si eliges una carpeta gestionada por otro programa, ese programa se encarga de copiarla a sus servidores. Selector Cine 365 no inicia sesión, no recibe cambios y no mezcla listas de varios dispositivos. Para recuperar una copia, utiliza Restaurar copia y selecciona su JSON.</p>
`)
 cut("        <div class=\"sync-desktop-panel\" id=\"cloudSyncPanel\" hidden>","      </div>\n    </div>\n  </section>","")
 cut("<dialog id=\"cloudSyncDialog\"","<script>\n(() => {","")
 cut("let cineGraphConnected=false,","const cineSyncDeviceLabel=",`let cloudManualBusy=false,extraBackupDirty=false,extraBackupBusy=false;
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
`)
 cut("let syncPreviewToken='',syncPreviewEdit=0,cloudManualBusy=false;","// Subscription/free/ad-supported services only;",`void cloudRefresh();
void cineRefreshHistory();
void cineRefreshRecovery();

`)
 cut("<article><h3>5. OneDrive y sincronización</h3>","<article><h3>6. Copias y seguridad</h3>",`<article><h3>5. Copia adicional</h3><p>En <strong>Mis listas → Guardado y copia adicional</strong> puedes elegir una segunda ubicación para tus datos. La biblioteca principal permanece en el equipo. La copia se actualiza automáticamente y no se importa ni se combina con otras bibliotecas.</p></article>
        `)
 if err!=nil{return "",err}
 page=strings.NewReplacer(
  "La nube y la clave TMDB son opcionales.","La copia adicional y la clave TMDB son opcionales.",
  "Mis listas → OneDrive y copias de seguridad → Actividad reciente, historial y copias de recuperación","Mis listas → Guardado y copia adicional → Actividad reciente, historial y copias de recuperación",
  "procedimientos paso a paso, sincronización, solución de problemas y recuperación","procedimientos paso a paso, guardado, solución de problemas y recuperación",
  "Consultar cambios locales y recibidos de otros dispositivos","Consultar cambios guardados",
  "Los cambios que haces aquí se guardan primero en el PC. Los cambios recibidos desde otros dispositivos aparecen indicados en este historial.","Los cambios que haces aquí se guardan en este equipo. El historial conserva también los movimientos antiguos.",
  "Las horas son orientativas. Para decidir qué cambio prevalece se comprueban las revisiones y el orden causal, no solo el reloj.","Las horas corresponden a los movimientos registrados.",
  "Los cambios restaurados se compartirán cuando pulses «Sincronizar con la nube».","La copia adicional, si está activada, se actualizará después.",
  "Antes de sincronizar","Copia de recuperación",
  "guardado, restauración o sincronización","guardado o restauración",
  "void checkCloudIncoming();","",
  "const longRequest=/^\\/api\\/cloud\\/(?:auto|sync|preview|backup)(?:\\?|$)/.test(path);","const longRequest=false;",
  "path!=='/api/cloud/browse'","path!=='/api/extra-backup/browse'",
  "'OneDrive: '","'Guardado: '",
  "if(path.startsWith('/api/cloud/')||path.startsWith('/api/onedrive-desktop/'))return 'OneDrive';","if(path.startsWith('/api/extra-backup/'))return 'Copia adicional';",
  "✓ Guardado local + carpeta de nube","✓ Guardado",
  "✓ Local · ⚠ Copia en nube","✓ Local · ⚠ Historial local",
  "✓ Guardado local · Revisar sincronización","✓ Guardado",
  "const sync=data.onedrive_connected?'OneDrive conectado':data.cloud_connected?'Carpeta de nube configurada':'Sin nube conectada';","const sync='Guardado local';",
  "(data.sync_pending_review?'⚠ Sincronización pendiente de revisión. ':'')+","",
  "No se tocarán las listas, fichas personales ni OneDrive.","No se tocarán las listas ni fichas personales.",
 ).Replace(page)
 return page,nil
}
