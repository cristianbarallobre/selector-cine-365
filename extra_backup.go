package main

import (
 "bytes"
 "encoding/json"
 "errors"
 "net/http"
 "io"
 "slices"
 "os"
 "path/filepath"
 "strings"
 "sync"
 "time"
)

type extraBackupConfig struct {
 Folder string `json:"folder"`
 ID string `json:"id"`
}
type extraBackup struct {
 mu sync.Mutex
 configPath, source string
 config extraBackupConfig
 ready bool
 pending,busy bool
 last string
 lastError string
 generation uint64
 wake chan struct{}
}
func newExtraBackup(dir string)*extraBackup {
 b:=&extraBackup{source:dir,configPath:filepath.Join(dir,"copia_adicional_ajustes.json"),wake:make(chan struct{},1)}
 data,err:=os.ReadFile(b.configPath)
 if err==nil {
  if json.Unmarshal(data,&b.config)!=nil||!uiDraftSession.MatchString(b.config.ID){
   b.lastError="La configuración de copia adicional no es válida. Elige otra vez una carpeta.";b.config=extraBackupConfig{}
  }
 }else if !os.IsNotExist(err){b.lastError="No se pudo leer la configuración de copia adicional."}
 if b.config.ID=="" {b.config.ID,_=randomToken()}
 b.ready=uiDraftSession.MatchString(b.config.ID)
 b.pending=b.config.Folder!=""
 return b
}
func (b *extraBackup) target(c extraBackupConfig)string{
 if c.Folder==""{return ""}
 return filepath.Join(c.Folder,"Selector_Cine_365_Copias",c.ID,"colecciones.json")
}
func (b *extraBackup) status()map[string]any{
 b.mu.Lock();defer b.mu.Unlock()
 return map[string]any{"folder":b.config.Folder,"target":b.target(b.config),"enabled":b.config.Folder!="","pending":b.pending,"busy":b.busy,"last_backup":b.last,"error":b.lastError}
}
func withinFolder(base,path string)bool{
 rel,err:=filepath.Rel(base,path)
 return err==nil&&(rel=="."||rel!=".."&&!strings.HasPrefix(rel,".."+string(os.PathSeparator)))
}
func(b *extraBackup) configure(folder string)error{
 folder=strings.TrimSpace(folder)
 if folder!=""{
  if !filepath.IsAbs(folder){return errors.New("Selecciona una carpeta existente con una ruta completa")}
  resolved,err:=filepath.EvalSymlinks(folder)
  if err!=nil{return errors.New("La carpeta elegida no está disponible")}
  info,err:=os.Stat(resolved);if err!=nil||!info.IsDir(){return errors.New("La ubicación debe ser una carpeta")}
  source,err:=filepath.EvalSymlinks(b.source);if err!=nil{return errors.New("No se pudo verificar la carpeta de datos")}
  if withinFolder(source,resolved){return errors.New("Elige una ubicación distinta de la carpeta de datos de la aplicación")}
  folder=resolved
 }
 b.mu.Lock();defer b.mu.Unlock()
 if !b.ready{return errors.New("No se pudo crear el identificador de copia de este equipo")}
 next:=b.config;next.Folder=folder
 data,err:=json.Marshal(next);if err!=nil{return err}
 if err=os.MkdirAll(b.source,0700);err!=nil{return err}
 if err=replaceLibraryFile(b.configPath,data);err!=nil{return err}
 b.config=next;b.generation++;b.last="";b.lastError="";b.pending=folder!=""
 b.signalLocked()
 return nil
}
func(b *extraBackup) signalLocked(){select{case b.wake<-struct{}{}:default:}}
func(b *extraBackup) request(){
 if b==nil{return}
 b.mu.Lock();defer b.mu.Unlock()
 if b.config.Folder!=""{b.pending=true;b.generation++;b.signalLocked()}
}
func(b *extraBackup) run(a *runningApp){
 ticker:=time.NewTicker(30*time.Second);defer ticker.Stop()
 b.request()
 for{select{case <-a.stopCh:return;case <-b.wake:b.copy(a);case <-ticker.C:b.copy(a)}}
}
func(b *extraBackup) copy(a *runningApp){
 b.mu.Lock()
 if !b.pending||b.busy||b.config.Folder==""{b.mu.Unlock();return}
 c,generation:=b.config,b.generation;b.busy=true;b.mu.Unlock()
 a.mu.Lock();data,err:=json.MarshalIndent(a.library,"","  ")
 var history []activityEntry
 hasHistory:=a.syncStore!=nil
 if hasHistory{history=a.syncStore.localActivity();data,err=json.MarshalIndent(completeLibraryBackup{libraryState:a.library,CompleteVersion:1,Activity:history},"","  ")}
 a.mu.Unlock()
 if err==nil{if len(data)>libraryFileLimit{err=errors.New("la biblioteca excede el límite de copia")}else{err=b.write(c,data);}}
 // The reported time is the file's actual modification time, not the
 // time of an unchanged verification or navigation-only save.
 savedAt:=""
 if err==nil{if st,e:=os.Lstat(b.target(c));e==nil&&st.Mode().IsRegular(){savedAt=st.ModTime().Format(time.RFC3339)}else{err=errors.New("No se pudo comprobar la fecha de la copia adicional")}}
 b.mu.Lock();defer b.mu.Unlock()
 b.busy=false
 if b.config!=c{return}
 if err!=nil{b.lastError="No se pudo actualizar la copia adicional. Se conserva el guardado principal.";diagnosticError("copia","escribir copia adicional",err);return}
 b.lastError="";b.last=savedAt
 if generation==b.generation{b.pending=false}else{b.signalLocked()}
}
// Unlike ordinary recovery reads, a zero-byte damaged copy is eligible for
// quarantine. Do not follow symlinks or read unbounded files.
func readExtraBackupFile(path string)([]byte,error){
 expected,e:=os.Lstat(path);if e!=nil{return nil,e}
 if !expected.Mode().IsRegular()||expected.Size()>libraryFileLimit{return nil,errors.New("Archivo de copia no seguro o demasiado grande")}
 f,e:=os.Open(path);if e!=nil{return nil,e};defer f.Close()
 st,e:=f.Stat();if e!=nil{return nil,e}
 if !st.Mode().IsRegular()||!os.SameFile(expected,st)||st.Size()>libraryFileLimit{return nil,errors.New("Archivo de copia cambiado")}
 data,e:=io.ReadAll(io.LimitReader(f,libraryFileLimit+1));if e!=nil{return nil,e}
 if len(data)>libraryFileLimit{return nil,errors.New("Archivo de copia demasiado grande")}
 return data,nil
}
// Name by content so repeated failures cannot create duplicate quarantines.
func preserveDamagedCopy(dir,prefix string,data []byte)error{
 path:=filepath.Join(dir,prefix+"_danada_"+libraryBytesFingerprint(data)+".json")
 prior,e:=readExtraBackupFile(path)
 if e==nil{if bytes.Equal(prior,data){return nil};return errors.New("La cuarentena existente no coincide")}
 if !os.IsNotExist(e){return e}
 if e=replaceLibraryFile(path,data);e!=nil{return e}
 saved,e:=readExtraBackupFile(path);if e!=nil||!bytes.Equal(saved,data){return errors.New("No se pudo verificar el archivo dañado conservado")}
 return nil
}
func sameExtraBackupContent(a,b libraryState,ah,bh *[]activityEntry)bool{
 if (ah==nil)!=(bh==nil){return false}
 if ah!=nil&&!slices.Equal(*ah,*bh){return false}
 return sameShareableLibrary(a,b)
}
func(b *extraBackup) write(c extraBackupConfig,data []byte)error{
 if !uiDraftSession.MatchString(c.ID){return errors.New("identificador de copia inválido")}
 root,err:=filepath.EvalSymlinks(c.Folder);if err!=nil{return err}
 info,err:=os.Stat(root);if err!=nil||!info.IsDir(){return errors.New("carpeta de copia no disponible")}
 source,err:=filepath.EvalSymlinks(b.source);if err!=nil{return err}
 if withinFolder(source,root){return errors.New("destino dentro de los datos principales")}
 // A unique per-installation directory prevents two PCs overwriting each other.
 dir:=root
 for _,part:=range []string{"Selector_Cine_365_Copias",c.ID}{
  dir=filepath.Join(dir,part)
  if err=os.Mkdir(dir,0700);err!=nil&&!os.IsExist(err){return err}
  info,err=os.Lstat(dir);if err!=nil||!info.IsDir()||info.Mode()&os.ModeSymlink!=0{return errors.New("carpeta de copia no segura")}
  real,e:=filepath.EvalSymlinks(dir);if e!=nil||!withinFolder(root,real)||withinFolder(source,real){return errors.New("destino de copia no seguro")}
 }
 next,nextHistory,err:=decodeRecoveryBackup(data);if err!=nil{return errors.New("biblioteca no válida")}
 target:=filepath.Join(dir,"colecciones.json")
 previous:=filepath.Join(dir,"colecciones_anterior.json")
 for _,name:=range []string{target,previous}{
  if info,e:=os.Lstat(name);e==nil{if !info.Mode().IsRegular()||info.Mode()&os.ModeSymlink!=0{return errors.New("archivo de copia no seguro")}}else if !os.IsNotExist(e){return e}
 }
 old,err:=readExtraBackupFile(target)
 if err==nil{
  if bytes.Equal(old,data){return nil}
  prior,priorHistory,decodeErr:=decodeRecoveryBackup(old)
  if decodeErr!=nil{
   if err=preserveDamagedCopy(dir,"colecciones",old);err!=nil{return err}
   diagnosticEvent("WARN","copia","copia adicional dañada conservada; reconstruyendo desde el guardado principal")
  }else{
   if sameExtraBackupContent(prior,next,priorHistory,nextHistory){return nil}
   if err=replaceLibraryFile(previous,old);err!=nil{return err}
   daily:=filepath.Join(dir,"copia_"+time.Now().UTC().Format("2006-01-02")+".json")
   if _,e:=os.Lstat(daily);os.IsNotExist(e){if e=replaceLibraryFile(daily,old);e!=nil{return e}}else if e!=nil{return e}
   entries,e:=os.ReadDir(dir);if e!=nil{return e}
   var names []string
   for _,entry:=range entries{if entry.Type().IsRegular()&&strings.HasPrefix(entry.Name(),"copia_")&&labDailyBackupName.MatchString(strings.TrimPrefix(entry.Name(),"copia_")){names=append(names,entry.Name())}}
   var bytes int64;kept:=0
   for i:=len(names)-1;i>=0;i--{file:=filepath.Join(dir,names[i]);info,e:=os.Lstat(file);if e!=nil||!info.Mode().IsRegular(){continue};if kept>0&&(kept>=14||bytes+info.Size()>512<<20){if e=os.Remove(file);e!=nil{return e};continue};kept++;bytes+=info.Size()}
  }
  latest,e:=readExtraBackupFile(target);if e!=nil||!bytes.Equal(latest,old){return errors.New("La copia adicional cambió durante la operación; se reintentará")}
 }else if !os.IsNotExist(err){return err}
 if err=replaceLibraryFile(target,data);err!=nil{return err}
 actual,err:=readRecoveryFile(target)
 if err!=nil||!bytes.Equal(actual,data){return errors.New("no se pudo verificar la copia adicional")}
 return nil
}
func(b *extraBackup) finish(a *runningApp,limit time.Duration){
 if b==nil{return}
 done:=make(chan struct{});deadline:=time.Now().Add(limit)
 go func(){defer close(done);for{
  b.mu.Lock();pending,busy:=b.pending,b.busy;b.mu.Unlock()
  if !pending&&!busy||time.Now().After(deadline){return}
  if !busy{b.copy(a);return}
  time.Sleep(50*time.Millisecond)
 }}()
 select{case <-done:case <-time.After(limit):diagnosticEvent("WARN","copia","copia adicional pendiente al cerrar")}
}
func(a *runningApp) serveExtraBackup(w http.ResponseWriter,r *http.Request){
 if !a.auth(r)||a.extraBackup==nil{http.Error(w,"No autorizado",403);return}
 b:=a.extraBackup
 w.Header().Set("Content-Type","application/json; charset=utf-8")
 fail:=func(code int,e error){w.WriteHeader(code);_=json.NewEncoder(w).Encode(map[string]string{"error":e.Error()})}
 if r.URL.Path=="/api/extra-backup/status"&&r.Method==http.MethodGet{_=json.NewEncoder(w).Encode(b.status());return}
 if r.Method!=http.MethodPost{http.Error(w,"Método no permitido",405);return}
 switch r.URL.Path{
 case "/api/extra-backup/browse":
  folder,err:=chooseBackupFolder();if err!=nil{fail(400,err);return}
  _=json.NewEncoder(w).Encode(map[string]string{"folder":folder});return
 case "/api/extra-backup/config":
  var req struct{Folder *string `json:"folder"`}
  if json.NewDecoder(http.MaxBytesReader(w,r.Body,32<<10)).Decode(&req)!=nil||req.Folder==nil{fail(400,errors.New("Carpeta no válida"));return}
  if err:=b.configure(*req.Folder);err!=nil{fail(400,err);return}
 case "/api/extra-backup/now":b.request()
 default:http.NotFound(w,r);return
 }
 _=json.NewEncoder(w).Encode(b.status())
}
