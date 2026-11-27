package main

import (
 "encoding/json"
 "errors"
 "net/http"
 "strings"
 "unicode/utf8"
)

// Portable history is deliberately independent of device clocks and secrets.
// Activity is newest first, the same order used by the activity viewer.
type completeLibraryBackup struct{
 libraryState
 CompleteVersion int `json:"complete_backup_version"`
 Activity []activityEntry `json:"activity"`
}
func validateCompleteBackup(b completeLibraryBackup)error{
 if b.CompleteVersion!=1||len(b.Activity)>activityLimit{return errors.New("Formato de copia completa no válido")}
 if err:=validateLibrary(b.libraryState);err!=nil{return err}
 for _,e:=range b.Activity{
  switch e.Kind{case "seen","unseen","film_add","film_remove","film_edit","list_create","list_delete","list_rename","favorite_add","favorite_remove","rating","identity","order","comparison","restore","bulk":default:return errors.New("Movimiento de historial no válido")}
  if e.At<0||e.At>4102444800000{return errors.New("Fecha de historial no válida")}
  for _,s:=range []string{e.Film,e.List,e.Detail}{if !utf8.ValidString(s)||utf8.RuneCountInString(s)>1000{return errors.New("Texto de historial no válido")}}
 }
 return nil
}
// Accept old library-only files, complete portable files and local checkpoints.
// Never import device clocks, tokens or cloud configuration from a backup.
func decodeRecoveryBackup(data []byte)(libraryState,*[]activityEntry,error){
 var b completeLibraryBackup
 if err:=json.Unmarshal(data,&b);err!=nil{return libraryState{},nil,err}
 if err:=validateLibrary(b.libraryState);err!=nil{return libraryState{},nil,err}
 if b.CompleteVersion!=0{
  if err:=validateCompleteBackup(b);err!=nil{return libraryState{},nil,err}
  return b.libraryState,&b.Activity,nil
 }
 var disk libraryDisk
 if err:=json.Unmarshal(data,&disk);err!=nil{return libraryState{},nil,err}
 if disk.LocalSync!=nil{
  s:=&desktopSyncStore{state:*disk.LocalSync}
  b.CompleteVersion=1;b.Activity=s.localActivity()
  if err:=validateCompleteBackup(b);err!=nil{return libraryState{},nil,err}
  return b.libraryState,&b.Activity,nil
 }
 return b.libraryState,nil,nil
}
func(a *runningApp) completeBackupLocked()completeLibraryBackup{
 return completeLibraryBackup{libraryState:a.library,CompleteVersion:1,Activity:a.syncStore.localActivity()}
}
func(a *runningApp) serveCompleteBackup(w http.ResponseWriter,r *http.Request){
 if !a.auth(r)||a.extraBackup==nil{http.Error(w,"No autorizado",403);return}
 if a.syncStore==nil{http.Error(w,"Historial local no disponible",503);return}
 if r.URL.Path=="/api/extra-backup/export"{
  if r.Method!=http.MethodGet{http.Error(w,"Método no permitido",405);return}
  a.mu.Lock();data,err:=json.Marshal(a.completeBackupLocked());a.mu.Unlock()
  if err!=nil||len(data)>libraryFileLimit{http.Error(w,"No se pudo preparar la copia",500);return}
  w.Header().Set("Content-Type","application/json; charset=utf-8");_,_=w.Write(data);return
 }
 if r.Method!=http.MethodPost{http.Error(w,"Método no permitido",405);return}
 var b completeLibraryBackup
 if json.NewDecoder(http.MaxBytesReader(w,r.Body,libraryFileLimit)).Decode(&b)!=nil||validateCompleteBackup(b)!=nil{http.Error(w,"Copia completa no válida",400);return}
 a.mu.Lock();defer a.mu.Unlock()
 if r.Header.Get("If-Match")!=a.libraryETag(){http.Error(w,"La biblioteca cambió; no se ha restaurado",412);return}
 // Keep a full rescue point before any mutation, including the current history.
 if err:=a.saveSyncRecovery(a.library,"_antes_restaurar");err!=nil{http.Error(w,"No se pudo guardar la copia previa; no se ha restaurado",500);return}
 next:=b.libraryState;recordMembershipEdits(a.library,&next)
 mirror,metadata,err:=a.commitProtectedLibraryWithHistory(next,nil,nil,"",&b.Activity,"restore")
 if err!=nil{http.Error(w,err.Error(),500);return}
 a.library=next;a.libraryRevision++;a.extraBackup.request();a.resolvePendingCopies(next,false)
 w.Header().Set("ETag",a.libraryETag())
 if mirror{w.Header().Set("X-Cine-Legacy","pending")}
 if metadata!=nil{diagnosticError("guardado","espejo de historial pendiente",metadata)}
 w.WriteHeader(204)
}
func completeBackupHTML(page string)string{
 old:=`$('downloadBackup').addEventListener('click',()=>{
 syncActiveText();const now=new Date().toISOString().slice(0,10);
 const backup={...libraryState,backup_created_at:new Date().toISOString()};
 saveDownload('Cine365_copia_'+now+'.json',JSON.stringify(backup,null,2),'application/json;charset=utf-8');
});`
 next:=`$('downloadBackup').addEventListener('click',async()=>{
 const button=$('downloadBackup');button.disabled=true;
 try{
  if(!await flushLibrary())throw new Error('Hay cambios pendientes de guardar');
  const backup=await cloudRequest('/api/extra-backup/export');
  saveDownload('Cine365_copia_completa_'+new Date().toISOString().slice(0,10)+'.json',JSON.stringify(backup,null,2),'application/json;charset=utf-8');
 }catch(e){alert('No se pudo exportar la copia completa: '+e.message);}finally{button.disabled=false;}
});`
 return strings.NewReplacer(
  old,next,
  "Se sustituirán todas las listas actuales.","Se sustituirán todas las listas actuales y, si la copia contiene historial, se restaurará también.",
  "async function writeLibrary(payload){","async function writeLibrary(payload,complete=false){",
  "('/api/library',{method:'POST',headers,body:payload,cache:'no-store'},12000)","(complete?'/api/extra-backup/restore':'/api/library',{method:'POST',headers,body:payload,cache:'no-store'},12000)",
  "libraryBase=payload;","if(complete){const clean=JSON.parse(payload);delete clean.complete_backup_version;delete clean.activity;libraryBase=JSON.stringify(clean);}else libraryBase=payload;",
  "const saved=await writeLibrary(payload);","const saved=await writeLibrary(payload,state.complete_backup_version===1);delete state.complete_backup_version;delete state.activity;",
 ).Replace(page)
}
