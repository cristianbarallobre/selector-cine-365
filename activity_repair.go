package main
import("encoding/json";"errors";"os";"path/filepath";"slices")

// Restore only the exact committed history, never an approximate older history.
// All candidates are read-only; the current library is not replaced.
func restoreExactActivitySnapshot(root,ref string) error {
 if !activityRefPattern.MatchString(ref){return errors.New("Referencia de historial no válida")}
 candidates:=[]string{filepath.Join(root,"biblioteca_guardada.json")}
 for _,group:=range []struct{dir string; names func(string)([]string,error)}{
  {"copias_rapidas_lab",func(dir string)([]string,error){return labBackupCandidates(dir,labQuickBackupName)}},
  {"copias_automaticas",func(dir string)([]string,error){return labBackupCandidates(dir,labDailyBackupName)}},
  {"sync_recuperacion",func(dir string)([]string,error){return labBackupCandidates(dir,syncRecoveryName)}},
 }{
  dir:=filepath.Join(root,group.dir);names,_:=group.names(dir)
  for i:=len(names)-1;i>=0;i--{if !stringsHasPendingSuffix(names[i]){candidates=append(candidates,filepath.Join(dir,names[i]))}}
 }
 b:=newExtraBackup(root)
 if target:=b.target(b.config);target!=""{
  candidates=append([]string{target,filepath.Join(filepath.Dir(target),"colecciones_anterior.json")},candidates...)
 }
 for _,file:=range candidates{
  data,e:=readRecoveryFile(file);if e!=nil{continue}
  _,history,e:=decodeRecoveryBackup(data);if e!=nil||history==nil{continue}
  events:=append([]activityEntry(nil),(*history)...);slices.Reverse(events)
  canonical,e:=json.Marshal(events);if e!=nil{continue}
  if len(events)==0&&libraryBytesFingerprint(canonical)!=ref{canonical=[]byte("[]")}
  if libraryBytesFingerprint(canonical)!=ref{continue}
  // Preserve a damaged regular snapshot before replacing it.
  path,e:=activitySnapshotPath(root,ref);if e!=nil{return e}
  if old,e:=readExtraBackupFile(path);e==nil{
   if libraryBytesFingerprint(old)!=ref{if e=preserveDamagedCopy(filepath.Dir(path),ref,old);e!=nil{return e}}
  }else if !os.IsNotExist(e){return e}
  return writeActivitySnapshot(root,ref,canonical)
 }
 return errors.New("No se encontró una copia con la huella exacta del historial")
}
