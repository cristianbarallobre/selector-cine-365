package main

// Immutable, content-addressed activity snapshots. The library's atomic
// replacement is the commit point. Every recovery copy remains self-contained;
// a snapshot written before a failed library commit cannot become visible.
import (
 "encoding/json"
 "errors"
 "fmt"
 "os"
 "path/filepath"
 "regexp"
 "sort"
 "strings"
)
var activityRefPattern=regexp.MustCompile("^[0-9a-f]{64}$")
const activityStorageDir="historial_local"
func encodeLocalLibrary(state libraryState, checkpoint *syncDeviceState)([]byte,string,[]byte,error){
 disk:=libraryDisk{libraryState:state,LocalSync:checkpoint}
 if checkpoint==nil||!checkpoint.ActivityMigrated{
  data,e:=json.MarshalIndent(disk,"","  ");return data,"",nil,e
 }
 events,e:=json.Marshal(checkpoint.Activity);if e!=nil{return nil,"",nil,e}
 if len(events)>libraryFileLimit{return nil,"",nil,errors.New("Historial demasiado grande")}
 ref:=libraryBytesFingerprint(events)
 metadata:=*checkpoint;metadata.Activity=nil
 disk.LocalSync=&metadata;disk.ActivityRef=ref
 data,e:=json.MarshalIndent(disk,"","  ");return data,ref,events,e
}
func activitySnapshotPath(root,ref string)(string,error){
 if !activityRefPattern.MatchString(ref){return "",errors.New("Referencia de historial no válida")}
 dir:=filepath.Join(root,activityStorageDir)
 if info,e:=os.Lstat(dir);e==nil{if !info.IsDir()||info.Mode()&os.ModeSymlink!=0{return "",errors.New("Carpeta de historial no segura")}}else if !os.IsNotExist(e){return "",e}
 return filepath.Join(dir,ref+".json"),nil
}
func readActivitySnapshot(root,ref string)([]activityEntry,error){
 path,e:=activitySnapshotPath(root,ref);if e!=nil{return nil,e}
 data,e:=readRecoveryFile(path)
 var events []activityEntry
 if e==nil&&libraryBytesFingerprint(data)==ref&&json.Unmarshal(data,&events)==nil&&len(events)<=activityLimit{return events,nil}
 // The existing metadata mirror is retained as a fallback and for old
 // releases. Never accept a mirror whose history does not match the hash.
 var mirror syncDeviceState
 data,e=readRecoveryFile(filepath.Join(root,"sync_dispositivo_v2.json"))
 if e==nil&&json.Unmarshal(data,&mirror)==nil{
  canonical,e:=json.Marshal(mirror.Activity)
  if e==nil&&libraryBytesFingerprint(canonical)==ref&&len(mirror.Activity)<=activityLimit{return mirror.Activity,nil}
 }
 return nil,errors.New("No se encontró una copia íntegra del historial confirmado")
}
func writeActivitySnapshot(root,ref string,data []byte)error{
 path,e:=activitySnapshotPath(root,ref);if e!=nil{return e}
 if libraryBytesFingerprint(data)!=ref{return errors.New("Huella de historial incorrecta")}
 if old,e:=readRecoveryFile(path);e==nil&&libraryBytesFingerprint(old)==ref{return nil}
 if info,e:=os.Lstat(path);e==nil&&!info.Mode().IsRegular(){return errors.New("Archivo de historial no seguro")}else if e!=nil&&!os.IsNotExist(e){return e}
 if e=os.MkdirAll(filepath.Dir(path),0700);e!=nil{return e}
 if e=replaceLibraryFile(path,data);e!=nil{return fmt.Errorf("No se pudo guardar el historial: %w",e)}
 return nil
}
func expandActivityReference(root string,data []byte)([]byte,error){
 var disk libraryDisk
 if e:=json.Unmarshal(data,&disk);e!=nil{return nil,e}
 if disk.ActivityRef==""{return data,nil}
 if disk.LocalSync==nil{return nil,errors.New("Faltan metadatos de historial")}
 events,e:=readActivitySnapshot(root,disk.ActivityRef);if e!=nil{return nil,e}
 disk.LocalSync.Activity=events;disk.LocalSync.ActivityMigrated=true;disk.ActivityRef=""
 expanded,e:=json.MarshalIndent(disk,"","  ");if e!=nil{return nil,e}
 if len(expanded)>libraryFileLimit{return nil,errors.New("La copia completa supera el límite")}
 return expanded,nil
}
func pruneActivitySnapshots(root,active string){
 dir:=filepath.Join(root,activityStorageDir)
 entries,e:=os.ReadDir(dir);if e!=nil{return}
 protected:=map[string]bool{active:true}
 // Keep the activity of quarantined libraries: the raw damaged originals
 // remain available for manual diagnosis and recovery.
 rootEntries,e:=os.ReadDir(root);if e!=nil{return}
 for _,item:=range rootEntries{
  if !item.Type().IsRegular()||!strings.HasPrefix(item.Name(),"colecciones.json.corrupt_"){continue}
  raw,e:=readRecoveryFile(filepath.Join(root,item.Name()));if e!=nil{continue}
  var disk libraryDisk;if json.Unmarshal(raw,&disk)==nil&&activityRefPattern.MatchString(disk.ActivityRef){protected[disk.ActivityRef]=true}
 }
 type entry struct{name string;when int64}
 candidates:=[]entry{}
 for _,item:=range entries{
  name:=strings.TrimSuffix(item.Name(),".json")
  if !item.Type().IsRegular()||!strings.HasSuffix(item.Name(),".json")||!activityRefPattern.MatchString(name)||protected[name]{continue}
  if info,e:=item.Info();e==nil{candidates=append(candidates,entry{item.Name(),info.ModTime().UnixNano()})}
 }
 sort.Slice(candidates,func(i,j int)bool{return candidates[i].when>candidates[j].when})
 // Keep up to three recent inactive snapshots, in addition to the active
 // and any quarantine references. All standard recovery backups embed history.
 for i,item:=range candidates{if i>=3{if e=os.Remove(filepath.Join(dir,item.name));e!=nil{diagnosticError("historial","retirar instantánea antigua",e)}}}
}
