package main

import (
 "bytes"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "os"
 "os/exec"
 "path/filepath"
 "strings"
 "testing"
)

func TestExtraBackupUsesIndependentLocalCopy(t *testing.T){
 a:=durableTestApp(t)
 a.extraBackup=newExtraBackup(filepath.Dir(a.libraryPath))
 b:=a.extraBackup
 folder:=t.TempDir()
 if err:=b.configure(folder);err!=nil{t.Fatal(err)}
 b.copy(a)
 target:=b.target(b.config)
 first,err:=os.ReadFile(target);if err!=nil{t.Fatal(err)}
 var library libraryState
 if json.Unmarshal(first,&library)!=nil||validateLibrary(library)!=nil{t.Fatal("invalid backup")}
 if bytes.Contains(first,[]byte("_local_sync")){t.Fatal("private sync metadata copied")}
 a.library.Lists[0].Text="☑ Origen (2010)"
 b.request();b.copy(a)
 second,err:=os.ReadFile(target);if err!=nil||!bytes.Contains(second,[]byte("☑")){t.Fatal("copy not updated",err)}
 previous,err:=os.ReadFile(filepath.Join(filepath.Dir(target),"colecciones_anterior.json"))
 if err!=nil||!bytes.Equal(first,previous){t.Fatal("previous copy lost",err)}
 if err:=b.configure("");err!=nil{t.Fatal(err)}
 if _,err:=os.Stat(target);err!=nil{t.Fatal("disable removed backup",err)}
 if err:=b.configure(filepath.Dir(a.libraryPath));err==nil{t.Fatal("accepted primary data folder")}
 if err:=b.configure(folder);err!=nil{t.Fatal(err)}
 b.copy(a)
 other:=newExtraBackup(t.TempDir())
 if err:=other.configure(folder);err!=nil{t.Fatal(err)}
 if b.target(b.config)==other.target(other.config){t.Fatal("two installations share a destination")}
 if err:=os.WriteFile(target,[]byte("damaged backup"),0600);err!=nil{t.Fatal(err)}
 b.request();b.copy(a)
 if b.lastError!=""||b.pending{t.Fatal("backup not repaired",b.lastError)}
 repaired,_:=os.ReadFile(target);if !bytes.Equal(repaired,second){t.Fatal("backup not rebuilt")}
 damaged,_:=os.ReadFile(filepath.Join(filepath.Dir(target),"colecciones_danada_"+libraryBytesFingerprint([]byte("damaged backup"))+".json"))
 if string(damaged)!="damaged backup"{t.Fatal("damaged copy not preserved")}
 if a.library.Lists[0].Text!="☑ Origen (2010)"{t.Fatal("backup failure modified local library")}
}
func TestLocalOnlyDesktopSurface(t *testing.T){
 page,err:=makeHTML("☐ Origen (2010)","test")
 if err!=nil{t.Fatal(err)}
 page,err=localOnlyHTML(page);if err!=nil{t.Fatal(err)}
 for _,required:=range []string{"id=\"extraBackupFolder\"","/api/extra-backup/status","/api/extra-backup/export","/api/extra-backup/restore","async function writeLibrary(payload,complete=false)","const saved=await writeLibrary(payload,state.complete_backup_version===1)","id=\"cloudRecoveryList\"","id=\"cloudRecoveryRestore\"","/api/sync/recovery/restore"}{
  if !strings.Contains(page,required){t.Fatalf("missing %s",required)}
 }
 for _,forbidden:=range []string{"id=\"cloudHeaderSync\"","id=\"cloudGraphConnect\"","id=\"cloudSyncPanel\"","/api/onedrive-desktop/","/api/cloud/auto","/api/cloud/config","/api/cloud/browse","Sincronizar con la nube"}{
  if strings.Contains(page,forbidden){t.Fatalf("old cloud feature still exposed: %s",forbidden)}
 }
 // The actual served script, after all injections, must parse successfully.
 if node,e:=exec.LookPath("node");e==nil{
  start:=strings.Index(page,"<script>\n(() => {")
  end:=strings.LastIndex(page,"</script>")
  if start<0||end<=start{t.Fatal("main script absent")}
  file:=filepath.Join(t.TempDir(),"runtime.js")
  if err:=os.WriteFile(file,[]byte(page[start+len("<script>\n"):end]),0600);err!=nil{t.Fatal(err)}
  if out,e:=exec.Command(node,"--check",file).CombinedOutput();e!=nil{t.Fatalf("runtime JS syntax: %v: %s",e,out)}
 }
 a:=durableTestApp(t);a.extraBackup=newExtraBackup(filepath.Dir(a.libraryPath))
 for _,path:=range []string{"/api/cloud/status","/api/cloud/config","/api/cloud/auto","/api/onedrive-desktop/status","/api/onedrive-desktop/start","/?code=old"}{
  req:=httptest.NewRequest(http.MethodGet,path,nil);req.Header.Set("X-Atlas-Token","test")
  w:=httptest.NewRecorder();a.serveHTTP(w,req)
  if w.Code!=404{t.Fatalf("%s still active: %d",path,w.Code)}
 }
 req:=httptest.NewRequest(http.MethodGet,"/api/extra-backup/status",nil)
 w:=httptest.NewRecorder();a.serveHTTP(w,req)
 if w.Code!=403{t.Fatalf("unauthorized backup endpoint: %d",w.Code)}
 req.Header.Set("X-Atlas-Token","test")
 w=httptest.NewRecorder();a.serveHTTP(w,req)
 if w.Code!=200{t.Fatalf("backup endpoint: %d",w.Code)}
}
func TestExtraBackupScheduledOnlyAfterLocalCommit(t *testing.T){
 if !desktopFeatures{t.Skip("desktop")}
 a:=durableTestApp(t)
 a.extraBackup=newExtraBackup(filepath.Dir(a.libraryPath))
 if err:=a.extraBackup.configure(t.TempDir());err!=nil{t.Fatal(err)}
 a.extraBackup.copy(a)
 a.extraBackup.mu.Lock();pending:=a.extraBackup.pending;a.extraBackup.mu.Unlock()
 if pending{t.Fatal("initial copy not confirmed")}
 next:=a.library;next.Lists=append([]filmCollection(nil),next.Lists...);next.Lists[0].Text="☑ Origen (2010)"
 w:=postTestLibrary(t,a,next)
 if w.Code!=204{t.Fatal(w.Code,w.Body.String())}
 a.extraBackup.mu.Lock();pending=a.extraBackup.pending;a.extraBackup.mu.Unlock()
 if !pending{t.Fatal("successful local save did not schedule copy")}
 a.extraBackup.copy(a)
 data,err:=os.ReadFile(a.extraBackup.target(a.extraBackup.config))
 if err!=nil||!bytes.Contains(data,[]byte("☑")){t.Fatal("copy did not follow committed library",err)}
}

func TestExtraBackupIgnoresNavigationOnly(t *testing.T){
 a:=durableTestApp(t);a.extraBackup=newExtraBackup(filepath.Dir(a.libraryPath));b:=a.extraBackup
 a.library.Lists=append(a.library.Lists,filmCollection{ID:"other",Name:"Otra",Text:"☐ Prueba (2011)"})
 if e:=b.configure(t.TempDir());e!=nil{t.Fatal(e)};b.copy(a)
 path:=b.target(b.config);before,e:=os.ReadFile(path);if e!=nil{t.Fatal(e)}
 a.library.ActiveID="other";b.request();b.copy(a)
 after,e:=os.ReadFile(path);if e!=nil||!bytes.Equal(before,after){t.Fatal("navigation rewrote backup",e)}
 if _,e=os.Stat(filepath.Join(filepath.Dir(path),"colecciones_anterior.json"));!os.IsNotExist(e){t.Fatal("navigation generated previous copy")}
 a.library.Lists[1].Text="☑ Prueba (2011)";b.request();b.copy(a)
 after,e=os.ReadFile(path);if e!=nil||bytes.Equal(before,after){t.Fatal("real edit not copied",e)}
}
