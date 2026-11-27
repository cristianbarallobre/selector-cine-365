package main

import (
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "net/http"
 "os"
 "path/filepath"
 "regexp"
 "strconv"
 "strings"
 "sync"
)

const uiDraftLimit = 64 << 20
var uiDraftSession = regexp.MustCompile(`^[a-f0-9]{32}$`)
var uiDraftName = regexp.MustCompile(`^ui-(?:motor-)?[a-f0-9]{32}_cambios_pendientes\.json$`)
var uiDraftMu sync.Mutex

type uiDraftDisk struct{
 libraryState
 Revision uint64 `json:"_pending_revision"`
}

// One draft per UI session. An acknowledgement for an older revision cannot
// delete a newer draft. This file is independent of the HTTP motor on Windows.
func updateUIDraft(dir,session string,revision uint64,action,payload string) error{
 return updateUIDraftAt(dir,"ui-",session,revision,action,payload)
}
func updateUIDraftAt(dir,prefix,session string,revision uint64,action,payload string) error{
 if !uiDraftSession.MatchString(session)||revision==0||len(payload)>libraryFileLimit{return errors.New("borrador no válido")}
 uiDraftMu.Lock()
 defer uiDraftMu.Unlock()
 root:=filepath.Join(dir,"sync_recuperacion")
 path:=filepath.Join(root,prefix+session+pendingRecoverySuffix+".json")
 data,err:=os.ReadFile(path)
 var previous uiDraftDisk
 if err==nil{
  if json.Unmarshal(data,&previous)!=nil||validateLibrary(previous.libraryState)!=nil{return errors.New("borrador existente no válido")}
 }else if !errors.Is(err,os.ErrNotExist){return err}
 switch action{
 case "clear":
  if err==nil&&previous.Revision<=revision{return os.Remove(path)}
  return nil
 case "put":
  if err==nil&&previous.Revision>revision{return nil}
  var state libraryState
  if json.Unmarshal([]byte(payload),&state)!=nil||validateLibrary(state)!=nil{return errors.New("borrador de biblioteca no válido")}
  data,err=json.Marshal(uiDraftDisk{libraryState:state,Revision:revision})
  if err!=nil{return err}
  if len(data)>libraryFileLimit{return errors.New("borrador demasiado grande")}
  if err=os.MkdirAll(root,0700);err!=nil{return err}
  return replaceLibraryFile(path,data)
 default:return errors.New("acción de borrador no válida")
 }
}

// Only validated fixed paths and protocol fields are accepted from the native UI.
func nativeDraftMessage(dir,message string) string{
 p:=strings.SplitN(message,":",5)
 if len(p)!=5||p[0]!="cine-draft"||!uiDraftSession.MatchString(p[1]){return ""}
 revision,err:=strconv.ParseUint(p[2],10,53)
 if err!=nil||revision==0{return ""}
 status:="ok"
 if err=updateUIDraft(dir,p[1],revision,p[3],p[4]);err!=nil{
  status="error"
  diagnosticError("recuperacion","persistir borrador de interfaz",err)
 }
 return fmt.Sprintf("cine-draft-ack:%s:%d:%s",p[1],revision,status)
}
func (a *runningApp) serveUIDraft(w http.ResponseWriter,r *http.Request){
 if !desktopFeatures||r.Method!=http.MethodPost||!a.auth(r){http.Error(w,"No autorizado",403);return}
 body,err:=io.ReadAll(http.MaxBytesReader(w,r.Body,libraryFileLimit+256))
 if err!=nil{http.Error(w,"Borrador demasiado grande",413);return}
 ack:=nativeDraftMessage(filepath.Dir(a.libraryPath),string(body))
 if ack==""||strings.HasSuffix(ack,":error"){http.Error(w,"No se pudo guardar el borrador",500);return}
 w.WriteHeader(204)
}

// Preserve the same session draft, rather than producing a new duplicate
// recovery file for every failed save. Legacy clients keep their old behavior.
func (a *runningApp) preservePendingRequest(w http.ResponseWriter,r *http.Request,state libraryState){
 p:=strings.Split(r.Header.Get("X-Cine-Draft"),":")
 if len(p)==2&&uiDraftSession.MatchString(p[0]){
  revision,err:=strconv.ParseUint(p[1],10,53)
  if err==nil&&revision>0{
   path:=filepath.Join(filepath.Dir(a.libraryPath),"sync_recuperacion","ui-"+p[0]+pendingRecoverySuffix+".json")
   data,e:=readRecoveryFile(path)
   var draft uiDraftDisk
   if e==nil&&json.Unmarshal(data,&draft)==nil&&draft.Revision>=revision&&validateLibrary(draft.libraryState)==nil{
    w.Header().Set("X-Cine-Recovery","saved")
    return
   }
   payload,marshalErr:=json.Marshal(state)
   if marshalErr==nil&&updateUIDraftAt(filepath.Dir(a.libraryPath),"ui-motor-",p[0],revision,"put",string(payload))==nil{
    w.Header().Set("X-Cine-Recovery","saved")
    return
   }
  }
 }
 a.preservePendingSave(w,state)
}
func (a *runningApp) clearServerDraft(r *http.Request){
 p:=strings.Split(r.Header.Get("X-Cine-Draft"),":")
 if len(p)!=2||!uiDraftSession.MatchString(p[0]){return}
 revision,err:=strconv.ParseUint(p[1],10,53)
 if err==nil&&revision>0{
  diagnosticError("recuperacion","retirar borrador confirmado",updateUIDraftAt(filepath.Dir(a.libraryPath),"ui-motor-",p[0],revision,"clear",""))
 }
}
