package main

import (
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "fmt"
 "net/http"
 "net/url"
 "path/filepath"
 "reflect"
 "strings"
 "time"
)

func (a *runningApp) serveGraphAuthHTTP(w http.ResponseWriter,r *http.Request){
 w.Header().Set("Content-Type","application/json; charset=utf-8")
 if !a.auth(r){http.Error(w,"No autorizado",403);return}
 respond:=func(code int,v any){w.WriteHeader(code);_ = json.NewEncoder(w).Encode(v)}
 if a.graph==nil{respond(503,map[string]string{"error":"Microsoft Graph no está disponible"});return}
 switch r.URL.Path {
 case "/api/onedrive-desktop/status":
  if r.Method!=http.MethodGet{http.Error(w,"Método no permitido",405);return}
  respond(200,a.graph.status())
 case "/api/onedrive-desktop/start":
  if r.Method!=http.MethodPost{http.Error(w,"Método no permitido",405);return}
  if a.cloud.status().Enabled {
   respond(409,map[string]string{"error":"Ya hay una carpeta configurada. Para cambiar a Graph, exporta primero una copia JSON y desactiva explícitamente la configuración anterior."});return
  }
  data,err:=a.graph.begin()
  if err!=nil{respond(409,map[string]string{"error":err.Error()});return}
  respond(200,data)
 case "/api/onedrive-desktop/disconnect":
  if r.Method!=http.MethodPost{http.Error(w,"Método no permitido",405);return}
  a.graph.logout()
  a.mu.Lock();a.syncPreviewToken="";a.mu.Unlock()
  respond(200,a.graph.status())
 default:http.NotFound(w,r)
 }
}

func (g *desktopGraph) existingBase()(string,error){
 root:=g.pairingID()
 if root==""{return "",errors.New("Conecta primero OneDrive")}
 id:=strings.TrimPrefix(root,"onedrive://item/")
 for _,name:=range []string{"actual","sincronizacion"}{
  item,exists,err:=g.child(id,name)
  if err!=nil{return "",err}
  if exists&&item.Folder!=nil{return id,nil}
 }
 item,exists,err:=g.child(id,cloudDirectoryName)
 if err!=nil{return "",err}
 if exists{
  if item.Folder==nil{return "",errors.New("La carpeta de copias de OneDrive está ocupada por un archivo")}
  return item.ID,nil
 }
 return "",nil
}
func (g *desktopGraph) ensureBase()(string,error){
 base,err:=g.existingBase();if err!=nil||base!=""{return base,err}
 root:=strings.TrimPrefix(g.pairingID(),"onedrive://item/")
 return g.ensureDir(root,cloudDirectoryName)
}
func (g *desktopGraph) childFolder(parent,name string)(string,error){
 if parent==""{return "",nil}
 item,exists,err:=g.child(parent,name)
 if err!=nil{return "",err}
 if !exists{return "",nil}
 if item.Folder==nil{return "",fmt.Errorf("%s no es una carpeta de OneDrive",name)}
 return item.ID,nil
}
func (a *runningApp) graphPlanFor(local libraryState)(syncPlan,error){
 if a.graph==nil||!a.graph.active(){return syncPlan{},errors.New("Conecta primero tu cuenta de OneDrive")}
 if a.syncStore==nil{return syncPlan{},errors.New("Sincronización no inicializada")}
 if err:=validateLibrary(local);err!=nil{return syncPlan{},err}
 pairing:=a.graph.pairingID()
 id,revision,known:=a.syncStore.parts()
 strict:=a.syncStore.isPaired(pairing)
 p:=syncPlan{Library:local,Revisions:revision,Known:known,Conflicts:[]syncConflict{}}
 h:=sha256.New()
 localBytes,_:=json.Marshal(local)
 h.Write(localBytes)
 process:=func(name string,data []byte,legacy bool)error {
  sum:=sha256.Sum256(data);fmt.Fprintf(h,"%s:%x|",name,sum)
  var peer syncSnapshot
  if legacy {
   peer.Protocol=1;peer.DeviceID="windows-legacy"
   if err:=json.Unmarshal(data,&peer.Library);err!=nil{return fmt.Errorf("biblioteca antigua corrupta: %w",err)}
   p.Legacy=true
  }else{
   if err:=json.Unmarshal(data,&peer);err!=nil{return fmt.Errorf("snapshot remoto corrupto (%s): %w",name,err)}
   if peer.DeviceID==id{return nil}
  }
  merged,rev,k,issues,err:=syncMergeVersioned(p.Library,p.Revisions,p.Known,peer,strict)
  if err!=nil{return err}
  p.Library=merged;p.Revisions=rev;p.Known=k;p.Sources++
  seen:=map[string]bool{}
  for _,issue:=range p.Conflicts{seen[issue.Key]=true}
  for _,issue:=range issues{if !seen[issue.Key]{p.Conflicts=append(p.Conflicts,issue);seen[issue.Key]=true}}
  return nil
 }
 base,err:=a.graph.existingBase();if err!=nil{return p,err}
 if base!="" {
  if !strict {
   actual,err:=a.graph.childFolder(base,"actual");if err!=nil{return p,err}
   if actual!=""{
    item,exists,err:=a.graph.child(actual,"colecciones.json");if err!=nil{return p,err}
    if exists{
     raw,err:=a.graph.readItem(item.ID);if err!=nil{return p,err}
     if err=process("actual/colecciones.json",raw,true);err!=nil{return p,err}
    }
   }
  }
  syncDir,err:=a.graph.childFolder(base,"sincronizacion");if err!=nil{return p,err}
  devices,err:=a.graph.childFolder(syncDir,"dispositivos");if err!=nil{return p,err}
  if devices!=""{
   list,err:=a.graph.list("me/drive/items/"+url.PathEscape(devices)+"/children?$select=id,name,folder,file&$top=200")
   if err!=nil{return p,err}
   if len(list)>100{return p,errors.New("Demasiadas copias de dispositivo; no se sincronizarán datos parciales")}
   for _,item:=range list{
    if item.File==nil||!strings.HasSuffix(item.Name,".json"){continue}
    raw,err:=a.graph.readItem(item.ID);if err!=nil{return p,err}
    if err=process("sincronizacion/dispositivos/"+item.Name,raw,false);err!=nil{return p,err}
   }
  }
 }
 p.Token=hex.EncodeToString(h.Sum(nil))
 return p,nil
}
func (a *runningApp) graphSyncPreview()(map[string]any,error){
 a.mu.Lock();defer a.mu.Unlock()
 p,err:=a.graphPlanFor(a.library);if err!=nil{return nil,err}
 a.syncPreviewToken=p.Token
 return map[string]any{
  "token":p.Token,"sources":p.Sources,"legacy":p.Legacy,
  "local_lists":len(a.library.Lists),"local_seen":syncWatched(a.library),
  "result_lists":len(p.Library.Lists),"result_seen":syncWatched(p.Library),
  "conflicts":p.Conflicts,"requires_confirmation":true,
 },nil
}
func (a *runningApp) graphSyncConfirm(token string,decisions map[string]bool)(map[string]any,error){
 a.mu.Lock()
 if token==""||token!=a.syncPreviewToken{a.mu.Unlock();return nil,errors.New("Revisa primero los datos de OneDrive")}
 p,err:=a.graphPlanFor(a.library)
 if err!=nil{a.mu.Unlock();return nil,err}
 if token!=p.Token{a.mu.Unlock();return nil,errors.New("OneDrive ha cambiado desde la vista previa. Revisa otra vez")}
 for _,issue:=range p.Conflicts{
  if _,ok:=decisions[issue.Key];!ok{
   a.mu.Unlock();return nil,fmt.Errorf("Falta resolver el conflicto %s",issue.Key)
  }
 }
 before,err:=json.MarshalIndent(a.library,"","  ")
 if err!=nil{a.mu.Unlock();return nil,err}
 recovery:=filepath.Join(filepath.Dir(a.libraryPath),"sync_recuperacion",time.Now().UTC().Format("20060102T150405.000000000")+".json")
 if err=cloudWriteFile(recovery,before);err!=nil{a.mu.Unlock();return nil,err}
 p.Library=syncApplyDecisions(p.Library,decisions)
 if err=saveLibrary(a.libraryPath,a.savePath,p.Library);err!=nil{a.mu.Unlock();return nil,err}
 a.library=p.Library
 if err=a.syncStore.accept(p,decisions,a.graph.pairingID());err!=nil{a.mu.Unlock();return nil,err}
 a.syncPreviewToken=""
 a.syncPendingConflict.Store(false)
 a.syncLastCheck.Store(time.Now().UnixMilli())
 a.syncLastReceived.Store(time.Now().UnixMilli())
 a.mu.Unlock()
 if err=a.graphPublishSnapshot();err!=nil{return nil,err}
 return map[string]any{"changed":true,"sources":p.Sources,"lists":len(p.Library.Lists),"seen":syncWatched(p.Library)},nil
}
func (a *runningApp) graphPublishSnapshot()error{
 if a.graph!=nil{a.graph.publishMu.Lock();defer a.graph.publishMu.Unlock()}
 if a.graph==nil||!a.graph.active(){return errors.New("OneDrive no está conectado")}
 if a.syncStore==nil||!a.syncStore.isPaired(a.graph.pairingID()){
  return errors.New("Revisa y confirma primero la sincronización")
 }
 a.mu.Lock();lib:=a.library;a.mu.Unlock()
 bytes,err:=a.syncStore.snapshot(lib);if err!=nil{return err}
 base,err:=a.graph.ensureBase();if err!=nil{return err}
 syncDir,err:=a.graph.ensureDir(base,"sincronizacion");if err!=nil{return err}
 devices,err:=a.graph.ensureDir(syncDir,"dispositivos");if err!=nil{return err}
 id,_,_:=a.syncStore.parts()
 name:=id+".json"
 old,exists,err:=a.graph.child(devices,name);if err!=nil{return err}
 // The snapshot contains a timestamp, but an identical payload still needs no PUT.
 if exists {
  previous,err:=a.graph.readItem(old.ID);if err!=nil{return err}
  var previousSnapshot,currentSnapshot syncSnapshot
  if json.Unmarshal(previous,&previousSnapshot)==nil&&json.Unmarshal(bytes,&currentSnapshot)==nil {
   previousSnapshot.UpdatedAt=0;currentSnapshot.UpdatedAt=0
   if reflect.DeepEqual(previousSnapshot,currentSnapshot){return nil}
  }
 }
 if err=a.graph.writeItem(devices,name,bytes);err!=nil{return err}
 actual,exists,err:=a.graph.child(devices,name)
 if err!=nil{return err}
 if !exists{return errors.New("OneDrive no confirmó la copia del dispositivo")}
 confirm,err:=a.graph.readItem(actual.ID);if err!=nil{return err}
 if sha256.Sum256(confirm)!=sha256.Sum256(bytes){return errors.New("OneDrive no confirmó el contenido de la copia")}
 return nil
}
func (a *runningApp) graphAutoSync(force bool)(map[string]any,error){
 if !a.graph.active()||a.syncStore==nil{
  return map[string]any{"changed":false,"skipped":"disconnected"},nil
 }
 pairing:=a.graph.pairingID()
 if !a.syncStore.isPaired(pairing){
  return map[string]any{"changed":false,"skipped":"not_paired","needs_review":true},nil
 }
 a.graph.mu.Lock()
 if !force&&time.Since(a.graph.lastCheck)<90*time.Second{
  a.graph.mu.Unlock()
  return map[string]any{"changed":false,"skipped":"cooldown"},nil
 }
 a.graph.lastCheck=time.Now()
 a.graph.mu.Unlock()
 a.mu.Lock()
 p,err:=a.graphPlanFor(a.library)
 if err!=nil{a.mu.Unlock();return nil,err}
 a.syncLastCheck.Store(time.Now().UnixMilli())
 if len(p.Conflicts)>0{
  a.syncPendingConflict.Store(true);a.mu.Unlock()
  return map[string]any{"changed":false,"needs_review":true,"conflicts":len(p.Conflicts)},nil
 }
 a.syncPendingConflict.Store(false)
 if p.Sources==0{
  a.mu.Unlock()
  if force {if err:=a.graphPublishSnapshot();err!=nil{return nil,err}}
  return map[string]any{"changed":false,"skipped":"no_remote_sources"},nil
 }
 _,oldRev,oldKnown:=a.syncStore.parts()
 libChanged:=!reflect.DeepEqual(a.library,p.Library)
 metadataChanged:=!reflect.DeepEqual(oldRev,p.Revisions)||!reflect.DeepEqual(oldKnown,p.Known)
 if !libChanged&&!metadataChanged{
  a.mu.Unlock()
  if force {if err:=a.graphPublishSnapshot();err!=nil{return nil,err}}
  return map[string]any{"changed":false},nil
 }
 if libChanged{
  before,_:=json.MarshalIndent(a.library,"","  ")
  recovery:=filepath.Join(filepath.Dir(a.libraryPath),"sync_recuperacion",time.Now().UTC().Format("20060102T150405.000000000")+".json")
  if err=cloudWriteFile(recovery,before);err!=nil{a.mu.Unlock();return nil,err}
  if err=saveLibrary(a.libraryPath,a.savePath,p.Library);err!=nil{a.mu.Unlock();return nil,err}
  a.library=p.Library
  a.syncLastReceived.Store(time.Now().UnixMilli())
 }
 if err=a.syncStore.accept(p,nil,pairing);err!=nil{a.mu.Unlock();return nil,err}
 a.syncPreviewToken=""
 a.mu.Unlock()
 if err=a.graphPublishSnapshot();err!=nil{return nil,err}
 return map[string]any{"changed":libChanged,"metadata_updated":metadataChanged,
  "lists":len(p.Library.Lists),"seen":syncWatched(p.Library),"sources":p.Sources},nil
}

func (a *runningApp) queueGraphUpload(){
 if a.graphUploadCh==nil||a.cloud==nil||!a.cloud.autoMode(){return}
 select{case a.graphUploadCh<-struct{}{}:default:}
}

// One small batch after edits: serialize uploads, no per-keystroke network,
// never publish before the user has reviewed the initial pairing.
func (a *runningApp) graphPublisher(){
 for{
  select {
  case <-a.stopCh:return
  case <-a.graphUploadCh:
   timer:=time.NewTimer(3*time.Second)
   select{case <-timer.C:case <-a.stopCh:timer.Stop();return}
   // Drain duplicate notifications without dropping the newest library state.
   for {
    select{case <-a.graphUploadCh:continue
    default:goto publish}
   }
  publish:
   if a.graph==nil||!a.graph.active()||a.syncStore==nil||!a.syncStore.isPaired(a.graph.pairingID()){continue}
   if err:=a.graphPublishSnapshot();err!=nil{
    a.graph.mu.Lock();a.graph.lastError=err.Error();a.graph.mu.Unlock()
   }else{
    a.graph.mu.Lock();a.graph.lastError="";a.graph.mu.Unlock()
   }
  }
 }
}
