package main

import (
 "encoding/json"
 "os"
 "path/filepath"
 "strings"
 "testing"
 "time"
)

func labTestLibrary(name string) libraryState{
 return libraryState{Version:1,ActiveID:"main",Lists:[]filmCollection{
  {ID:"main",Name:name,Text:"☑ Origen (2000)\n☐ Pendiente (2001)"},
 }}
}
func TestLabRecoveryBadLibraryDoesNotGetOverwritten(t *testing.T){
 if !protectedLibraryEnabled{t.Skip("Solo LAB")}
 dir:=t.TempDir()
 path:=filepath.Join(dir,"colecciones.json")
 bad:=[]byte("{datos de biblioteca dañados, no borrar}")
 if err:=os.WriteFile(path,bad,0600);err!=nil{t.Fatal(err)}
 _,_,err:=labLoadLibraryProtected(path,"☑ Inicial (2000)")
 if err==nil||!strings.Contains(err.Error(),"dañado"){t.Fatalf("Esperaba rechazo al iniciar: %v",err)}
 raw,e:=os.ReadFile(path)
 if e!=nil||string(raw)!=string(bad){t.Fatal("No se debe reemplazar el JSON dañado por datos iniciales")}
}
func TestLabRecoveryFromValidSnapshotQuarantinesDamagedOriginal(t *testing.T){
 if !protectedLibraryEnabled{t.Skip("Solo LAB")}
 dir:=t.TempDir()
 path:=filepath.Join(dir,"colecciones.json")
 original:=[]byte("{contenido roto}")
 if err:=os.WriteFile(path,original,0600);err!=nil{t.Fatal(err)}
 good:=labTestLibrary("Recuperada")
 data,_:=json.Marshal(good)
 folder:=filepath.Join(dir,"copias_automaticas")
 if err:=os.MkdirAll(folder,0700);err!=nil{t.Fatal(err)}
 if err:=os.WriteFile(filepath.Join(folder,"2026-11-22.json"),data,0600);err!=nil{t.Fatal(err)}
 state,recovered,err:=labLoadLibraryProtected(path,"☐ Inicial (2000)")
 if err!=nil||!recovered{t.Fatalf("Copia válida no restaurada: %v",err)}
 if state.Lists[0].Name!="Recuperada"{t.Fatal("Se recuperó una lista incorrecta")}
 raw,err:=os.ReadFile(path)
 if err!=nil||!labValidBackup(raw){t.Fatal("El JSON activo no fue reparado")}
 items,err:=os.ReadDir(dir);if err!=nil{t.Fatal(err)}
 found:=false
 for _,item:=range items{
  if strings.HasPrefix(item.Name(),"colecciones.json.corrupt_"){
   old,_:=os.ReadFile(filepath.Join(dir,item.Name()))
   found=string(old)==string(original)
  }
 }
 if !found{t.Fatal("El JSON dañado se ha perdido en lugar de guardarse aparte")}
}
func TestLabRecoverySkipsInvalidMostRecentCopy(t *testing.T){
 if !protectedLibraryEnabled{t.Skip("Solo LAB")}
 dir:=t.TempDir()
 path:=filepath.Join(dir,"colecciones.json")
 os.WriteFile(path,[]byte("{roto}"),0600)
 folder:=filepath.Join(dir,"copias_automaticas")
 os.MkdirAll(folder,0700)
 good,_:=json.Marshal(labTestLibrary("Válida"))
 a:=filepath.Join(folder,"2026-11-01.json")
 b:=filepath.Join(folder,"2026-11-02.json")
 os.WriteFile(a,good,0600)
 os.WriteFile(b,[]byte("{roto}"),0600)
 old:=time.Now().Add(-time.Hour)
 _=os.Chtimes(a,old,old)
 state,recovered,err:=labLoadLibraryProtected(path,"☐ No usar")
 if err!=nil||!recovered||state.Lists[0].Name!="Válida"{t.Fatalf("No omitió la copia reciente inválida: %v",err)}
}
func TestLabBackupsHaveBoundedRetention(t *testing.T){
 if !protectedLibraryEnabled{t.Skip("Solo LAB")}
 dir:=t.TempDir()
 path:=filepath.Join(dir,"colecciones.json")
 good,_:=json.Marshal(labTestLibrary("Una lista"))
 labQuickBackup(path,good)
 fastDir:=filepath.Join(dir,"copias_rapidas_lab")
 first,err:=labBackupCandidates(fastDir,labQuickBackupName)
 if err!=nil||len(first)!=1{t.Fatalf("Primera copia rápida no creada: %v %v",first,err)}
 labQuickBackup(path,good)
 second,_:=labBackupCandidates(fastDir,labQuickBackupName)
 if len(second)!=1{t.Fatal("No crear una copia con cada pulsación o cambio de estado")}
 for i:=0;i<labQuickBackupCount+10;i++{
  name:=time.Date(2026,1,1,0,i,0,0,time.UTC).Format("20060102T150405.000000000Z")+".json"
  _=os.WriteFile(filepath.Join(fastDir,name),good,0600)
 }
 labPruneNamedBackups(fastDir,labQuickBackupName,labQuickBackupCount)
 after,_:=labBackupCandidates(fastDir,labQuickBackupName)
 if len(after)!=labQuickBackupCount{t.Fatalf("Sin límite para las copias rápidas: %d",len(after))}
 daily:=filepath.Join(dir,"copias_automaticas")
 os.MkdirAll(daily,0700)
 for i:=0;i<labDailyBackupCount+12;i++{
  name:=time.Date(2024,1,1,0,0,0,0,time.UTC).AddDate(0,0,i).Format("2006-01-02")+".json"
  _=os.WriteFile(filepath.Join(daily,name),good,0600)
 }
 // Los archivos ajenos, aunque estén en el mismo directorio, no se eliminan.
 keep:=filepath.Join(daily,"mi_copia_manual.json")
 os.WriteFile(keep,good,0600)
 labPruneLibraryBackups(path)
 remaining,_:=labBackupCandidates(daily,labDailyBackupName)
 if len(remaining)!=labDailyBackupCount{t.Fatalf("El historial diario no está acotado: %d",len(remaining))}
 if _,err:=os.Stat(keep);err!=nil{t.Fatal("Se eliminó un archivo manual")}
}
