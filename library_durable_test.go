package main

import (
 "encoding/json"
 "os"
 "path/filepath"
 "strings"
 "strconv"
 "testing"
)

func durableTestApp(t *testing.T)*runningApp{
 t.Helper()
 dir:=t.TempDir()
 store,err:=newDesktopSyncStore(dir);if err!=nil{t.Fatal(err)}
 return &runningApp{libraryPath:filepath.Join(dir,"colecciones.json"),savePath:filepath.Join(dir,"mis_365_peliculas.txt"),library:initialLibrary("☐ Origen (2010)"),libraryDigest:"missing",syncStore:store,token:"test"}
}

func TestDurableCheckpointSurvivesFailedMetadataWrite(t *testing.T){
 if !desktopFeatures{t.Skip("desktop")}
 a:=durableTestApp(t)
 if w:=postTestLibrary(t,a,a.library);w.Code!=204{t.Fatal(w.Code,w.Body.String())}
 originalPath:=a.syncStore.path
 blocked:=filepath.Join(filepath.Dir(originalPath),"blocked")
 if err:=os.Mkdir(blocked,0700);err!=nil{t.Fatal(err)}
 a.syncStore.path=blocked
 next:=initialLibrary("☑ Origen (2010)")
 w:=postTestLibrary(t,a,next)
 if w.Code!=204||w.Header().Get("X-Cine-Cloud")!="error"{t.Fatal(w.Code,w.Body.String(),w.Header())}
 cp:=readLibraryCheckpoint(a.libraryPath)
 if cp==nil||!cp.Revisions["origen|2010"].Checked{t.Fatal("checkpoint does not preserve revision")}
 reopened,err:=recoverSyncCheckpoint(filepath.Dir(originalPath),cp)
 if err!=nil{t.Fatal(err)}
 _,revisions,_:=reopened.parts()
 if !revisions["origen|2010"].Checked{t.Fatal("lost revision after restart")}
 if !strings.Contains(loadLibrary(a.libraryPath,"").Lists[0].Text,"☑"){t.Fatal("lost library")}
 // The checkpoint is local-only, including in the legacy folder backup.
 folder:=t.TempDir()
 cloud:=newCloudBackups(filepath.Dir(a.libraryPath))
 if err=cloud.configure(folder);err!=nil{t.Fatal(err)}
 if status:=cloud.backup();status.LastError!=""{t.Fatal(status.LastError)}
 shared,err:=os.ReadFile(filepath.Join(cloudSyncRoot(folder),"actual","colecciones.json"))
 if err!=nil||strings.Contains(string(shared),"_local_sync"){t.Fatal("private checkpoint leaked",err)}
}
func TestDurableExternalChangeAndMissingLibrary(t *testing.T){
 if !protectedLibraryEnabled{t.Skip("desktop")}
 a:=durableTestApp(t)
 if w:=postTestLibrary(t,a,a.library);w.Code!=204{t.Fatal(w.Code,w.Body.String())}
 external:=initialLibrary("☑ Externa (2001)")
 raw,_:=json.Marshal(external)
 if err:=os.WriteFile(a.libraryPath,raw,0600);err!=nil{t.Fatal(err)}
 if w:=postTestLibrary(t,a,initialLibrary("☑ Otra (2010)"));w.Code!=500{t.Fatal(w.Code,w.Body.String())}
 if got,_:=os.ReadFile(a.libraryPath);string(got)!=string(raw){t.Fatal("external change overwritten")}
 if err:=os.Remove(a.libraryPath);err!=nil{t.Fatal(err)}
 state,recovered,err:=labLoadLibraryProtected(a.libraryPath,"Muestra")
 if err!=nil||!recovered||state.Lists[0].Text!="☐ Origen (2010)"{t.Fatal(state,recovered,err)}
}
func TestDurableBackupFailureWarnsWithoutLosingCommit(t *testing.T){
 if !desktopFeatures{t.Skip("desktop")}
 old:=libraryBackupIssue.Load();libraryBackupIssue.Store(false);t.Cleanup(func(){libraryBackupIssue.Store(old)})
 a:=durableTestApp(t)
 if w:=postTestLibrary(t,a,a.library);w.Code!=204{t.Fatal(w.Code)}
 if err:=os.WriteFile(filepath.Join(filepath.Dir(a.libraryPath),"copias_automaticas"),[]byte("blocked"),0600);err!=nil{t.Fatal(err)}
 next:=initialLibrary("☑ Origen (2010)")
 w:=postTestLibrary(t,a,next)
 if w.Code!=204||w.Header().Get("X-Cine-Backup")!="error"||loadLibrary(a.libraryPath,"").Lists[0].Text!=next.Lists[0].Text{t.Fatal(w.Code,w.Header(),w.Body.String())}
}
func TestDurableDraftRevisionAndRecovery(t *testing.T){
 dir:=t.TempDir()
 session:=strings.Repeat("a",32)
 before,_:=json.Marshal(initialLibrary("☐ Origen (2010)"))
 after,_:=json.Marshal(initialLibrary("☑ Origen (2010)"))
 put:=func(rev,payload string)string{return nativeDraftMessage(dir,"cine-draft:"+session+":"+rev+":put:"+payload)}
 if !strings.HasSuffix(put("1",string(before)),":ok")||!strings.HasSuffix(put("2",string(after)),":ok"){t.Fatal("draft not acknowledged")}
 if !strings.HasSuffix(put("1",string(before)),":ok"){t.Fatal("late old draft")}
 if !strings.HasSuffix(nativeDraftMessage(dir,"cine-draft:"+session+":1:clear:"),":ok"){t.Fatal("clear")}
 file:=filepath.Join(dir,"sync_recuperacion","ui-"+session+pendingRecoverySuffix+".json")
 data,err:=os.ReadFile(file);if err!=nil{t.Fatal(err)}
 var draft uiDraftDisk
 if json.Unmarshal(data,&draft)!=nil||draft.Revision!=2||draft.Lists[0].Text!="☑ Origen (2010)"{t.Fatal("newer draft destroyed")}
 a:=&runningApp{libraryPath:filepath.Join(dir,"colecciones.json")}
 copies,err:=a.localRecoveryList();if err!=nil||len(copies)!=1{t.Fatal(copies,err)}
 if nativeDraftMessage(dir,"cine-draft:../bad:3:put:{}")!=""{t.Fatal("invalid session accepted")}
 if !strings.HasSuffix(put("3","{}"),":error"){t.Fatal("invalid payload accepted")}
 if !strings.HasSuffix(nativeDraftMessage(dir,"cine-draft:"+session+":2:clear:"),":ok"){t.Fatal("clear latest")}
 if _,err=os.Stat(file);!os.IsNotExist(err){t.Fatal("resolved draft still exists",err)}
}
func TestDurableTombstonesDoNotBlockLifetimeSaves(t *testing.T){
 state:=initialLibrary("Cine")
 state.DeletedLists=map[string]bool{}
 for i:=0;i<1200;i++{state.DeletedLists["deleted_"+strconv.Itoa(i)]=true}
 if err:=validateLibrary(state);err!=nil{t.Fatal(err)}
}
