package main

// Salvaguardas adicionales de Windows LAB; no sustituyen los flujos de
// sincronización, confirmación y recuperación que ya tiene la aplicación.
import (
 "encoding/json"
 "errors"
 "fmt"
 "os"
 "path/filepath"
 "regexp"
 "sort"
 "strings"
 "time"
)

const (
 labQuickBackupCount=24
 labDailyBackupCount=60
 labBackupMinInterval=15*time.Minute
 labBackupLimit=3<<20
)
var labQuickBackupName=regexp.MustCompile(`^\d{8}T\d{6}\.\d{9}Z\.json$`)
var labDailyBackupName=regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\.json$`)

func labValidBackup(data []byte) bool {
 if len(data)==0||len(data)>labBackupLimit {return false}
 var state libraryState
 return json.Unmarshal(data,&state)==nil&&validateLibrary(state)==nil
}
func labBackupCandidates(dir string, pattern *regexp.Regexp) ([]string,error) {
 entries,err:=os.ReadDir(dir)
 if errors.Is(err,os.ErrNotExist){return nil,nil}
 if err!=nil{return nil,err}
 names:=make([]string,0,len(entries))
 for _,entry:=range entries{
  if !entry.Type().IsRegular()||!pattern.MatchString(entry.Name()){continue}
  names=append(names,entry.Name())
 }
 sort.Strings(names)
 return names,nil
}
func labPruneNamedBackups(dir string,pattern *regexp.Regexp,keep int) {
 names,err:=labBackupCandidates(dir,pattern)
 if err!=nil||len(names)<=keep{return}
 for _,name:=range names[:len(names)-keep]{
  _=os.Remove(filepath.Join(dir,name))
 }
}
func labQuickBackup(path string,previous []byte) {
 if !appIsLab||!labValidBackup(previous){return}
 dir:=filepath.Join(filepath.Dir(path),"copias_rapidas_lab")
 names,err:=labBackupCandidates(dir,labQuickBackupName)
 if err!=nil{return}
 now:=time.Now().UTC()
 if len(names)>0 {
  // Los nombres usan UTC y precisión nanosegundo, ordenables alfabéticamente.
  last:=strings.TrimSuffix(names[len(names)-1],".json")
  stamp,err:=time.Parse("20060102T150405.000000000Z",last)
  if err==nil&&now.Sub(stamp)<labBackupMinInterval{return}
 }
 if err:=os.MkdirAll(dir,0700);err!=nil{return}
 filename:=filepath.Join(dir,now.Format("20060102T150405.000000000Z")+".json")
 // No se copian API keys, tokens ni caché: solo el JSON válido de biblioteca.
 if err:=os.WriteFile(filename,previous,0600);err==nil{
  labPruneNamedBackups(dir,labQuickBackupName,labQuickBackupCount)
 }
}
func labPruneLibraryBackups(path string) {
 if !appIsLab{return}
 labPruneNamedBackups(filepath.Join(filepath.Dir(path),"copias_automaticas"),labDailyBackupName,labDailyBackupCount)
}

// Ante colecciones.json dañado, NUNCA iniciar con una biblioteca nueva y
// sobrescribir el archivo. Recuperar la última copia válida dejando el archivo
// original intacto en un nombre de cuarentena.
func labLoadLibraryProtected(path,legacy string) (libraryState,bool,error) {
 data,err:=os.ReadFile(path)
 if errors.Is(err,os.ErrNotExist){return initialLibrary(legacy),false,nil}
 if err!=nil{return libraryState{},false,fmt.Errorf("No se puede leer la biblioteca: %w",err)}
 if labValidBackup(data) {
  var state libraryState
  if err:=json.Unmarshal(data,&state);err==nil{return state,false,nil}
 }
 root:=filepath.Dir(path)
 candidates:=[]struct{dir string;pattern *regexp.Regexp}{
  {filepath.Join(root,"copias_rapidas_lab"),labQuickBackupName},
  {filepath.Join(root,"copias_automaticas"),labDailyBackupName},
 }
 type backup struct{filename string;when time.Time}
 found:=[]backup{}
 for _,c:=range candidates{
  names,e:=labBackupCandidates(c.dir,c.pattern)
  if e!=nil{continue}
  for _,name:=range names{
   file:=filepath.Join(c.dir,name)
   info,e:=os.Stat(file);if e==nil&&info.Size()<=labBackupLimit{
    found=append(found,backup{file,info.ModTime()})
   }
  }
 }
 sort.Slice(found,func(i,j int)bool{return found[i].when.After(found[j].when)})
 for _,candidate:=range found{
  restored,e:=os.ReadFile(candidate.filename)
  if e!=nil||!labValidBackup(restored){continue}
  var state libraryState
  if e=json.Unmarshal(restored,&state);e!=nil{continue}
  quarantine:=path+".corrupt_"+time.Now().UTC().Format("20060102T150405.000000000")+".json"
  if e=os.Rename(path,quarantine);e!=nil{
   return libraryState{},false,fmt.Errorf("La biblioteca está dañada y no se pudo poner el original a salvo: %w",e)
  }
  if e=os.WriteFile(path,restored,0600);e!=nil{
   // En caso de error, devolver el original sin tocar la copia buena.
   _=os.Rename(quarantine,path)
   return libraryState{},false,fmt.Errorf("No se pudo restaurar la copia local: %w",e)
  }
  return state,true,nil
 }
 return libraryState{},false,fmt.Errorf("El archivo %s está dañado. No se ha sobrescrito porque no hay ninguna copia válida. Conserva el archivo y restaura una copia manualmente",path)
}
