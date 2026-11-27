package main

import (
 _ "embed"
 "fmt"
 "strings"
)

//go:embed ui/local_backup.html
var localBackupMarkup string

//go:embed ui/local_backup.js
var localBackupScript string

//go:embed ui/activity.js
var localActivityScript string

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
 cut("      <div class=\"cloud-box\" aria-label=\"Configuración de sincronización y copias de seguridad\">","        <details class=\"cloud-guide-extra\">\n          <summary>Actividad reciente, historial y copias de recuperación</summary>",localBackupMarkup)
 cut("        <div class=\"sync-desktop-panel\" id=\"cloudSyncPanel\" hidden>","      </div>\n    </div>\n  </section>","")
 cut("<dialog id=\"cloudSyncDialog\"","<script>\n(() => {","")
 cut("let cineGraphConnected=false,","const cineSyncDeviceLabel=",localBackupScript)
 cut("async function cineRefreshHistory(){","async function cineRefreshRecovery(){",localActivityScript)
 cut("let syncPreviewToken='',syncPreviewEdit=0,cloudManualBusy=false;","// Subscription/free/ad-supported services only;",`void cloudRefresh();
void cineRefreshHistory();
void cineRefreshRecovery();

`)
 cut("<article><h3>5. OneDrive y sincronización</h3>","<article><h3>6. Copias y seguridad</h3>",`<article><h3>5. Copia adicional</h3><p>En <strong>Mis listas → Guardado y copia adicional</strong> puedes elegir una segunda ubicación para tus datos. La biblioteca principal permanece en el equipo. La copia se actualiza automáticamente y no se importa ni se combina con otras bibliotecas.</p></article>
        `)
 if err!=nil{return "",err}
 page=strings.NewReplacer(
  "Actividad reciente de películas","Historial de actividad",
  "<div id=\"cloudHistory\"","<div style=\"display:flex;gap:8px;flex-wrap:wrap;margin:8px 0\"><input class=\"cloud-path\" style=\"flex:1;min-width:160px\" id=\"activitySearch\" placeholder=\"Buscar película, lista o acción\" aria-label=\"Buscar actividad\"><select id=\"activityFilter\" aria-label=\"Tipo de movimiento\"><option value=\"\">Todos los movimientos</option></select></div><div style=\"display:flex;align-items:center;gap:8px;margin-bottom:8px\"><button class=\"tool\" id=\"activityPrev\">Anterior</button><span id=\"activityPage\"></span><button class=\"tool\" id=\"activityNext\">Siguiente</button></div><div id=\"cloudHistory\"",
  "if(name==='collections')renderCollectionControls();","if(name==='collections'){renderCollectionControls();void cineRefreshHistory();}",
  "await cineRefreshRecovery();\n  $('cloudRecoveryStatus').textContent=data.warning","await cineRefreshRecovery();void cineRefreshHistory();\n  $('cloudRecoveryStatus').textContent=data.warning",
  ".cine-sync-history{max-height:220px;",".cine-sync-history{max-height:560px;",
  "Las horas son orientativas.","Las horas corresponden al guardado confirmado.",
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
  "return msg.join(' | ')||'Ver alternativas en JustWatch';}","return (msg.join(' | ')||'Ver alternativas en JustWatch')+' · '+cinePlatformAge(r);}",
  "badge.textContent=name;","badge.textContent=name+(r?.date&&Date.now()-Date.parse(r.date)>=86400000?' · datos guardados':'');",
  "label.textContent=x[0]+'  ';","label.textContent=x[0]+' · '+cinePlatformAge(r)+'  ';",
  "async function loadPlatforms(force){","function cinePlatformAge(entry){const at=Date.parse(entry?.date||'');if(!Number.isFinite(at))return 'Sin fecha de comprobación';return (Date.now()-at>=86400000?'Datos guardados · ':'Comprobado · ')+new Date(at).toLocaleString('es-ES');}\nfunction cinePlatformExpired(entry){const at=Date.parse(entry?.date||'');return (!entry?.cine_retry_at||Date.now()>=entry.cine_retry_at)&&(!Number.isFinite(at)||Date.now()-at>=(entry.state==='error'?900000:86400000));}\nasync function loadPlatforms(force,staleOnly=false){",
  "const entries=force?all:all.filter(f=>{const cached=platformData.get(providerKey(f));return !cached||cached.state==='partial';});","const entries=force?all:all.filter(f=>{const cached=platformData.get(providerKey(f));return staleOnly?!!cached&&cinePlatformExpired(cached):!cached||cached.state==='partial'||cinePlatformExpired(cached);}).slice(0,staleOnly?24:all.length);",
  "storePlatformLookup(providerKey(f),r);if(r.tmdb_id>0)","if(cinePlatformExpired(r))r.cine_retry_at=Date.now()+900000;storePlatformLookup(providerKey(f),r);if(r.tmdb_id>0)",
  "if(current===platformRun)storePlatformLookup(providerKey(f),{state:'error',note:err.message});","if(current===platformRun){const old=platformData.get(providerKey(f));if(old)old.cine_retry_at=Date.now()+900000;else storePlatformLookup(providerKey(f),{state:'error',note:err.message,date:new Date().toISOString()});}",
  "No se tocarán las listas, fichas personales ni OneDrive.","No se tocarán las listas ni fichas personales.",
 ).Replace(page)
 return completeBackupHTML(page),nil
}
