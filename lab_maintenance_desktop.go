//go:build windows || darwin

package main

import (
 "embed"
 "encoding/json"
 "net/http"
 "os"
 "path/filepath"
 "runtime"
 "strings"
)

//go:embed lab_maintenance_windows.html
var labMaintenanceAssets embed.FS

// Añadir el panel dentro de Conectar TMDB / Sincronización, solo Windows LAB.
func injectLabMaintenance(page string) string {
 raw,err:=labMaintenanceAssets.ReadFile("lab_maintenance_windows.html")
 if err!=nil{return page}
 pieces:=strings.SplitN(string(raw),"<!-- LAB_MAINTENANCE_HTML -->",2)
 if len(pieces)!=2{return page}
 tail:=strings.SplitN(pieces[1],"<!-- LAB_MAINTENANCE_SCRIPT -->",2)
 if len(tail)!=2{return page}
 anchors:=[]string{
  `<p class="cloud-status" id="cloudRecoveryStatus" aria-live="polite"></p>`,
  `async function cloudRequest(path,method='GET',payload){`,
  `updateListPlatformOptions();update();renderCollectionControls();updatePlatformProgress();void initLibrary().then(`,
  `$('cloudRecoveryList').addEventListener('change',()=>{$('cloudRecoveryRestore').disabled=!$('cloudRecoveryList').value;});`,
 }
 for _,anchor:=range anchors{if !strings.Contains(page,anchor){return page}}
 page=strings.Replace(page,"</head>",pieces[0]+"\n</head>",1)
 page=strings.Replace(page,`<p class="cloud-status" id="cloudRecoveryStatus" aria-live="polite"></p>`,
  `<p class="cloud-status" id="cloudRecoveryStatus" aria-live="polite"></p>`+tail[0],1)
 page=strings.Replace(page,`(item.id.startsWith('rec:')?'Antes de sincronizar':'Copia diaria')`,
  `(item.id.startsWith('rec:')?'Antes de sincronizar':item.id.startsWith('quick:')?'Copia rápida':'Copia diaria')`,1)
 page=strings.Replace(page,
  `$('cloudRecoveryList').addEventListener('change',()=>{$('cloudRecoveryRestore').disabled=!$('cloudRecoveryList').value;});`,
  `$('cloudRecoveryList').addEventListener('change',()=>{void labPreviewRecovery();});`,1)
 page=strings.Replace(page,
  "updateListPlatformOptions();update();renderCollectionControls();updatePlatformProgress();void initLibrary().then(",
  tail[1]+"\nupdateListPlatformOptions();update();renderCollectionControls();updatePlatformProgress();void initLibrary().then(",1)
 return page
}

type labFileUsage struct {
 Files int `json:"files"`
 Bytes int64 `json:"bytes"`
}
func labLocalFileUsage(path string) labFileUsage {
 var usage labFileUsage
 info,err:=os.Stat(path)
 if err==nil&&info.Mode().IsRegular(){usage.Files=1;usage.Bytes=info.Size()}
 return usage
}
func labLocalPosterUsage(path string) labFileUsage {
 files,bytes:=posterDiskUsage(path)
 return labFileUsage{Files:files,Bytes:bytes}
}
func (a *runningApp) serveLabMaintenance(w http.ResponseWriter,r *http.Request) {
 w.Header().Set("Content-Type","application/json; charset=utf-8")
 if !a.auth(r){http.Error(w,"No autorizado",403);return}
 root:=filepath.Dir(a.libraryPath)
 cacheRoot:=appCacheRoot(root)
 posterDir:=filepath.Join(cacheRoot,"carteles")
 discoveryDir:=discoveryImageDir(a.libraryPath)
 respond:=func(code int,data any){w.WriteHeader(code);_ =json.NewEncoder(w).Encode(data)}
 switch r.Method{
 case http.MethodGet:
  daily,_:=labBackupCandidates(filepath.Join(root,"copias_automaticas"),labDailyBackupName)
  quick,_:=labBackupCandidates(filepath.Join(root,"copias_rapidas_lab"),labQuickBackupName)
  var mem runtime.MemStats;runtime.ReadMemStats(&mem)
  info:=map[string]any{
   "go_heap_bytes":mem.HeapAlloc,"go_reserved_bytes":mem.Sys,"go_workers":runtime.NumGoroutine(),
   "library_path":root,"cache_path":cacheRoot,
   "library":labLocalFileUsage(a.libraryPath),
   "tmdb_details":labLocalFileUsage(filepath.Join(cacheRoot,"tmdb_fichas_lab.json")),
   "platform_cache":labLocalFileUsage(filepath.Join(cacheRoot,"plataformas_cache.json")),
   "posters":labLocalPosterUsage(posterDir),
   "discovery_images":labLocalPosterUsage(discoveryDir),
   "backup_daily":len(daily),"backup_quick":len(quick),
   "backup_warning":libraryBackupIssue.Load(),
   "library_recovered":a.libraryRecovered,
   "tmdb_connected":a.platforms!=nil&&a.platforms.configured(),
   "cloud_connected":a.cloud!=nil&&a.cloud.status().Enabled,
   "onedrive_connected":a.graph!=nil&&a.graph.active(),
   "sync_pending_review":a.syncPendingConflict.Load(),
   "sync_last_check_ms":a.syncLastCheck.Load(),
   "sync_last_received_ms":a.syncLastReceived.Load(),
   "webview_managed_separately":true,
  }
  respond(200,info)
 case http.MethodPost:
  var req struct{Action string `json:"action"`}
  if err:=json.NewDecoder(http.MaxBytesReader(w,r.Body,4096)).Decode(&req);err!=nil{
   respond(400,map[string]string{"error":"Solicitud no válida"});return
  }
  switch req.Action{
  case "prune_posters":
   n,err:=prunePosterCache(posterDir,appPosterMaxFiles,appPosterMaxBytes,appPosterMaxAge)
   if err!=nil{respond(500,map[string]string{"error":"No se han podido revisar los carteles"});return}
   extra,err:=a.discoveryImages.maintain(discoveryDir,false)
   if err!=nil{respond(500,map[string]string{"error":"No se han podido revisar las imágenes de Descubrir"});return}
   n+=extra
   respond(200,map[string]any{"removed":n})
  case "clear_posters":
   n,err:=prunePosterCache(posterDir,0,0,0)
   if err!=nil{respond(500,map[string]string{"error":"No se han podido eliminar las miniaturas"});return}
   extra,err:=a.discoveryImages.maintain(discoveryDir,true)
   if err!=nil{respond(500,map[string]string{"error":"No se han podido eliminar las imágenes de Descubrir"});return}
   n+=extra
   respond(200,map[string]any{"removed":n})
  default:respond(400,map[string]string{"error":"Acción de mantenimiento no permitida"})
  }
 default:respond(405,map[string]string{"error":"Método no permitido"})
 }
}
