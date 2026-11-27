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
 // The base template now contains only local features. Expand explicit
 // module slots instead of deleting historical markup at runtime.
 if !strings.Contains(page,"__CINE_LOCAL_BACKUP_MARKUP__")&&strings.Contains(page,"id=\"extraBackupFolder\""){return page,nil}
 for _,slot:=range []struct{marker,body string}{
  {"__CINE_LOCAL_BACKUP_MARKUP__",localBackupMarkup},
  {"__CINE_LOCAL_BACKUP_SCRIPT__",localBackupScript},
  {"__CINE_LOCAL_ACTIVITY_SCRIPT__",localActivityScript},
 }{
  if strings.Count(page,slot.marker)!=1{return "",fmt.Errorf("Módulo de interfaz ausente o duplicado: %s",slot.marker)}
  page=strings.Replace(page,slot.marker,slot.body,1)
 }
 page=strings.NewReplacer(
  "Actividad reciente de películas","Historial de actividad",
  "<div id=\"cloudHistory\"","<div style=\"display:flex;gap:8px;flex-wrap:wrap;margin:8px 0\"><input class=\"cloud-path\" style=\"flex:1;min-width:160px\" id=\"activitySearch\" placeholder=\"Buscar película, lista o acción\" aria-label=\"Buscar actividad\"><select id=\"activityFilter\" aria-label=\"Tipo de movimiento\"><option value=\"\">Todos los movimientos</option></select></div><div style=\"display:flex;align-items:center;gap:8px;margin-bottom:8px\"><button class=\"tool\" id=\"activityPrev\">Anterior</button><span id=\"activityPage\"></span><button class=\"tool\" id=\"activityNext\">Siguiente</button></div><div id=\"cloudHistory\"",
  "if(name==='collections')renderCollectionControls();","if(name==='collections'){renderCollectionControls();void cineRefreshHistory();}",
  "await cineRefreshRecovery();\n  $('cloudRecoveryStatus').textContent=data.warning","await cineRefreshRecovery();void cineRefreshHistory();\n  $('cloudRecoveryStatus').textContent=data.warning",
  ".cine-sync-history{max-height:220px;",".cine-sync-history{max-height:560px;",
  "La nube y la clave TMDB son opcionales.","La copia adicional y la clave TMDB son opcionales.",
  "procedimientos paso a paso, sincronización, solución de problemas y recuperación","procedimientos paso a paso, guardado, solución de problemas y recuperación",
  "Consultar cambios locales y recibidos de otros dispositivos","Consultar cambios guardados",
  "Los cambios que haces aquí se guardan primero en el PC. Los cambios recibidos desde otros dispositivos aparecen indicados en este historial.","Los cambios que haces aquí se guardan en este equipo. El historial conserva también los movimientos antiguos.",
  "Los cambios restaurados se compartirán cuando pulses «Sincronizar con la nube».","La copia adicional, si está activada, se actualizará después.",
  "Antes de sincronizar","Copia de recuperación",
  "guardado, restauración o sincronización","guardado o restauración",
  "void checkCloudIncoming();","",
  "✓ Guardado local + carpeta de nube","✓ Guardado",
  "✓ Local · ⚠ Copia en nube","✓ Local · ⚠ Historial local",
  "✓ Guardado local · Revisar sincronización","✓ Guardado",
  "return msg.join(' | ')||'Ver alternativas en JustWatch';}","return (msg.join(' | ')||'Ver alternativas en JustWatch')+' · '+cinePlatformAge(r);}",
  "async function loadPlatforms(force){","function cinePlatformExpired(entry){const at=Date.parse(entry?.date||'');return (!entry?.cine_retry_at||Date.now()>=entry.cine_retry_at)&&(!Number.isFinite(at)||Date.now()-at>=(entry.state==='error'?900000:86400000));}\nasync function loadPlatforms(force,staleOnly=false){",
  "const entries=force?all:all.filter(f=>{const cached=platformData.get(providerKey(f));return !cached||cached.state==='partial';});","const entries=force?all:all.filter(f=>{const cached=platformData.get(providerKey(f));return staleOnly?!!cached&&cinePlatformExpired(cached):!cached||cached.state==='partial'||cinePlatformExpired(cached);}).slice(0,staleOnly?24:all.length);",
  "storePlatformLookup(providerKey(f),r);if(r.tmdb_id>0)","if(cinePlatformExpired(r))r.cine_retry_at=Date.now()+900000;storePlatformLookup(providerKey(f),r);if(r.tmdb_id>0)",
  "if(current===platformRun)storePlatformLookup(providerKey(f),{state:'error',note:err.message});","if(current===platformRun){const old=platformData.get(providerKey(f));if(old)old.cine_retry_at=Date.now()+900000;else storePlatformLookup(providerKey(f),{state:'error',note:err.message,date:new Date().toISOString()});}",
 ).Replace(page)
 return completeBackupHTML(page),nil
}
