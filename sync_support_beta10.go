package main

import (
 "encoding/json"
 "errors"
 "fmt"
 "net/http"
 "os"
 "path/filepath"
 "regexp"
 "sort"
 "strings"
 "time"
)

// Recovery backups stay on this computer. No cloud file or arbitrary path is
// accepted as a restore source, and restoration always creates a rescue point.
var safeRecoveryName=regexp.MustCompile(`^[0-9A-Za-z._-]{5,95}\.json$`)
type recoveryInfo struct {
 ID string `json:"id"`
 Name string `json:"name"`
 WhenMS int64 `json:"when_ms"`
 Bytes int64 `json:"bytes"`
}
func (a *runningApp) localRecoveryList() ([]recoveryInfo,error){
 root:=filepath.Dir(a.libraryPath)
 all:=[]recoveryInfo{}
 directories:=[]struct{prefix,path string}{
  {"rec",filepath.Join(root,"sync_recuperacion")},
  {"day",filepath.Join(root,"copias_automaticas")},
 }
 if protectedLibraryEnabled{directories=append(directories,struct{prefix,path string}{"quick",filepath.Join(root,"copias_rapidas_lab")})}
 for _,d:=range directories{
  entries,err:=os.ReadDir(d.path)
  if os.IsNotExist(err){continue}
  if err!=nil{return nil,err}
  for _,entry:=range entries{
   if entry.IsDir()||!safeRecoveryName.MatchString(entry.Name()){continue}
   info,err:=entry.Info();if err!=nil||info.Size()>3<<20{continue}
   all=append(all,recoveryInfo{
    ID:d.prefix+":"+entry.Name(),Name:entry.Name(),
    WhenMS:info.ModTime().UnixMilli(),Bytes:info.Size(),
   })
  }
 }
 sort.Slice(all,func(i,j int)bool{return all[i].WhenMS>all[j].WhenMS})
 if len(all)>50{all=all[:50]}
 return all,nil
}
// Solo información resumida de la copia elegida, sin exponer los TXT.
func (a *runningApp) previewLocalRecovery(id string) (map[string]any,error){
 prefix,name,ok:=strings.Cut(id,":")
 if !ok||!safeRecoveryName.MatchString(name){return nil,errors.New("copia no válida")}
 root:=filepath.Dir(a.libraryPath)
 dir:=""
 switch prefix{
 case "rec":dir=filepath.Join(root,"sync_recuperacion")
 case "day":dir=filepath.Join(root,"copias_automaticas")
 case "quick":if protectedLibraryEnabled{dir=filepath.Join(root,"copias_rapidas_lab")}
 }
 if dir==""{return nil,errors.New("origen no permitido")}
 data,err:=os.ReadFile(filepath.Join(dir,name))
 if err!=nil{return nil,err}
 if len(data)==0||len(data)>3<<20{return nil,errors.New("Copia demasiado grande o vacía")}
 var state libraryState
 if json.Unmarshal(data,&state)!=nil||validateLibrary(state)!=nil{
  return nil,errors.New("La copia no es válida. No es seguro restaurarla")
 }
 favorites:=0
 for _,rec:=range state.Reactions{if rec.Favorite{favorites++}}
 return map[string]any{
  "lists":len(state.Lists),"seen":syncWatched(state),"favorites":favorites,
  "valid":true,"message":"La copia ha superado la validación de la biblioteca",
 },nil
}
func (a *runningApp) restoreLocalRecovery(id string) (map[string]any,error) {
 prefix,name,ok:=strings.Cut(id,":")
 if !ok||!safeRecoveryName.MatchString(name) {return nil,errors.New("copia de recuperación no válida")}
 root:=filepath.Dir(a.libraryPath)
 dir:=""
 switch prefix {
 case "rec":dir=filepath.Join(root,"sync_recuperacion")
 case "day":dir=filepath.Join(root,"copias_automaticas")
 case "quick":if protectedLibraryEnabled{dir=filepath.Join(root,"copias_rapidas_lab")}
 default:return nil,errors.New("origen de recuperación no permitido")
 }
 filename:=filepath.Join(dir,name)
 data,err:=os.ReadFile(filename);if err!=nil{return nil,err}
 if len(data)==0||len(data)>3<<20{return nil,errors.New("tamaño de copia de recuperación no válido")}
 var restored libraryState
 if err=json.Unmarshal(data,&restored);err!=nil{return nil,errors.New("la copia JSON no es válida")}
 if err=validateLibrary(restored);err!=nil{return nil,fmt.Errorf("no se ha restaurado: %w",err)}
 a.mu.Lock()
 before:=a.library
 oldBytes,err:=json.MarshalIndent(before,"","  ")
 if err!=nil{a.mu.Unlock();return nil,err}
 rescue:=filepath.Join(root,"sync_recuperacion",time.Now().UTC().Format("20060102T150405.000000000")+"_antes_restaurar.json")
 if err=cloudWriteFile(rescue,oldBytes);err!=nil{a.mu.Unlock();return nil,err}
 if err=saveLibrary(a.libraryPath,a.savePath,restored);err!=nil{a.mu.Unlock();return nil,err}
 a.library=restored
 if a.syncStore!=nil{
  if err=a.syncStore.observe(before,restored);err!=nil{a.mu.Unlock();return nil,err}
 }
 a.syncPreviewToken=""
 a.mu.Unlock()
 warning:=""
 if a.cloud!=nil&&a.syncStore!=nil{
  folder:=a.cloud.status().Folder
  if a.syncStore.isPaired(folder) {
   saved:=a.cloud.backup()
   if saved.LastError!=""{warning=saved.LastError}else if err=a.publishDeviceSnapshot();err!=nil{warning=err.Error()}
  }
 }
 return map[string]any{
  "restored":true,"lists":len(restored.Lists),"seen":syncWatched(restored),"warning":warning,
 },nil
}
func (a *runningApp) serveSyncSupport(w http.ResponseWriter,r *http.Request){
 w.Header().Set("Content-Type","application/json; charset=utf-8")
 if !a.auth(r){http.Error(w,"No autorizado",403);return}
 respond:=func(code int,obj any){w.WriteHeader(code);_ =json.NewEncoder(w).Encode(obj)}
 switch r.URL.Path {
 case "/api/sync/revisions":
  if r.Method!=http.MethodGet{respond(405,map[string]string{"error":"Método incorrecto"});return}
  if a.syncStore==nil{respond(200,map[string]any{"revisions":map[string]syncRevision{}});return}
  _,revisions,_:=a.syncStore.parts()
  respond(200,map[string]any{"revisions":revisions})
 case "/api/sync/history":
  if r.Method!=http.MethodGet{respond(405,map[string]string{"error":"Método incorrecto"});return}
  if a.syncStore==nil{respond(200,map[string]any{"events":[]syncHistoryEntry{}});return}
  respond(200,map[string]any{"events":a.syncStore.history()})
 case "/api/sync/recovery/list":
  if r.Method!=http.MethodGet{respond(405,map[string]string{"error":"Método incorrecto"});return}
  copies,err:=a.localRecoveryList()
  if err!=nil{respond(500,map[string]string{"error":err.Error()});return}
  respond(200,map[string]any{"copies":copies})
 case "/api/sync/recovery/preview":
  if !protectedLibraryEnabled{respond(404,map[string]string{"error":"No disponible"});return}
  if r.Method!=http.MethodGet{respond(405,map[string]string{"error":"Método incorrecto"});return}
  info,err:=a.previewLocalRecovery(r.URL.Query().Get("id"))
  if err!=nil{respond(409,map[string]string{"error":err.Error()});return}
  respond(200,info)
 case "/api/sync/recovery/restore":
  if r.Method!=http.MethodPost{respond(405,map[string]string{"error":"Método incorrecto"});return}
  var request struct{ID string `json:"id"`}
  if err:=json.NewDecoder(http.MaxBytesReader(w,r.Body,4096)).Decode(&request);err!=nil{
   respond(400,map[string]string{"error":"Solicitud incorrecta"});return
  }
  restored,err:=a.restoreLocalRecovery(request.ID)
  if err!=nil{respond(409,map[string]string{"error":err.Error()});return}
  respond(200,restored)
 default:http.NotFound(w,r)
 }
}
