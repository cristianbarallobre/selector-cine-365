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
 // On disk, the 40 permitted lists can each contain 900,000 bytes.
 // JSON escaping can expand each byte sixfold; allow room for metadata too.
 libraryFileLimit=256<<20
)
var labQuickBackupName=regexp.MustCompile(`^\d{8}T\d{6}\.\d{9}Z\.json$`)
var labDailyBackupName=regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\.json$`)

func labValidBackup(data []byte) bool {
 if len(data)==0||len(data)>libraryFileLimit {return false}
 _,_,err:=decodeRecoveryBackup(data)
 return err==nil
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
 if err!=nil{return}
 const budget=int64(512<<20)
 var bytes int64
 kept:=0
 for i:=len(names)-1;i>=0;i--{
  file:=filepath.Join(dir,names[i]);info,e:=os.Lstat(file)
  if e!=nil||!info.Mode().IsRegular(){continue}
  if kept>0&&(kept>=keep||bytes+info.Size()>budget){_ = os.Remove(file);continue}
  kept++;bytes+=info.Size()
 }
}
func labQuickBackup(path string,previous []byte) {
 if !protectedLibraryEnabled||!labValidBackup(previous){return}
 dir:=filepath.Join(filepath.Dir(path),"copias_rapidas_lab")
 names,err:=labBackupCandidates(dir,labQuickBackupName)
 if err!=nil{recordBackupIssue("listar copias rapidas",err);return}
 now:=time.Now().UTC()
 if len(names)>0 {
  // Los nombres usan UTC y precisión nanosegundo, ordenables alfabéticamente.
  last:=strings.TrimSuffix(names[len(names)-1],".json")
  stamp,err:=time.Parse("20060102T150405.000000000Z",last)
  if err==nil&&now.Sub(stamp)<labBackupMinInterval{return}
 }
 if err:=os.MkdirAll(dir,0700);err!=nil{recordBackupIssue("crear carpeta de copias rapidas",err);return}
 filename:=filepath.Join(dir,now.Format("20060102T150405.000000000Z")+".json")
 // No se copian API keys, tokens ni caché: solo el JSON válido de biblioteca.
 if err:=replaceLibraryFile(filename,previous);err==nil{
  labPruneNamedBackups(dir,labQuickBackupName,labQuickBackupCount)
 }else{recordBackupIssue("guardar copia rapida",err)}
}
func labPruneLibraryBackups(path string) {
 if !protectedLibraryEnabled{return}
 labPruneNamedBackups(filepath.Join(filepath.Dir(path),"copias_automaticas"),labDailyBackupName,labDailyBackupCount)
}

// Ante colecciones.json dañado, NUNCA iniciar con una biblioteca nueva y
// sobrescribir el archivo. Recuperar la última copia válida dejando el archivo
// original intacto en un nombre de cuarentena.
func labLoadLibraryProtected(path,legacy string) (libraryState,bool,error) {
 data,err:=os.ReadFile(path)
 missing:=errors.Is(err,os.ErrNotExist)
 if err!=nil&&!missing{return libraryState{},false,fmt.Errorf("No se puede leer la biblioteca: %w",err)}
 if labValidBackup(data) {
  var state libraryState
  if err:=json.Unmarshal(data,&state);err==nil{return state,false,nil}
 }
 root:=filepath.Dir(path)
 candidates:=[]struct{dir string;pattern *regexp.Regexp}{
  {filepath.Join(root,"copias_rapidas_lab"),labQuickBackupName},
  {filepath.Join(root,"copias_automaticas"),labDailyBackupName},
  {filepath.Join(root,"sync_recuperacion"),syncRecoveryName},
 }
 type backup struct{filename string;when time.Time}
 found:=[]backup{}
 evidence:=false
 for _,c:=range candidates{
  names,e:=labBackupCandidates(c.dir,c.pattern)
  if e!=nil{evidence=true;recordBackupIssue("listar copias para recuperación",e);continue}
  for _,name:=range names{
   if stringsHasPendingSuffix(name){continue}
   evidence=true
   file:=filepath.Join(c.dir,name)
   info,e:=os.Stat(file);if e==nil&&info.Size()<=libraryFileLimit{
    found=append(found,backup{file,info.ModTime()})
   }
  }
 }
 first:=filepath.Join(root,"biblioteca_guardada.json")
 if info,e:=os.Stat(first);e==nil{evidence=true;if info.Mode().IsRegular()&&info.Size()<=libraryFileLimit{found=append(found,backup{first,info.ModTime()})}}else if !os.IsNotExist(e){evidence=true}
 sort.Slice(found,func(i,j int)bool{return found[i].when.After(found[j].when)})
 for _,candidate:=range found{
  restored,e:=os.ReadFile(candidate.filename)
  if e!=nil||!labValidBackup(restored){continue}
  var state libraryState
  if e=json.Unmarshal(restored,&state);e!=nil{continue}
  quarantine:=path+".corrupt_"+time.Now().UTC().Format("20060102T150405.000000000")+".json"
  if !missing {if e=os.Rename(path,quarantine);e!=nil{
   return libraryState{},false,fmt.Errorf("La biblioteca está dañada y no se pudo poner el original a salvo: %w",e)
  }}
  if e=replaceLibraryFile(path,restored);e!=nil{
   // En caso de error, devolver el original sin tocar la copia buena.
   if !missing {_=os.Rename(quarantine,path)}
   return libraryState{},false,fmt.Errorf("No se pudo restaurar la copia local: %w",e)
  }
  return state,true,nil
 }
 if missing&&!evidence{return initialLibrary(legacy),false,nil}
 return libraryState{},false,fmt.Errorf("El archivo %s está dañado. No se ha sobrescrito porque no hay ninguna copia válida. Conserva el archivo y restaura una copia manualmente",path)
}
