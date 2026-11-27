package main

import (
 "encoding/json"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestCloudAutoSyncImportsPairedChangesWithoutOverwritingPeer(t *testing.T){
 dir:=t.TempDir()
 localDir:=filepath.Join(dir,"pc")
 cloudDir:=filepath.Join(dir,"onedrive")
 if err:=os.MkdirAll(cloudDir,0700);err!=nil{t.Fatal(err)}
 local:=initialLibrary("☐ 001. Ciudadano Kane (1941)")
 path:=filepath.Join(localDir,"colecciones.json")
 savePath:=filepath.Join(localDir,"mis_365_peliculas.txt")
 if err:=saveLibrary(path,savePath,local);err!=nil{t.Fatal(err)}
 store,err:=newDesktopSyncStore(localDir);if err!=nil{t.Fatal(err)}
 cloud:=newCloudBackups(localDir);if err=cloud.configure(cloudDir);err!=nil{t.Fatal(err)}
 app:=&runningApp{syncStore:store,cloud:cloud,library:local,libraryPath:path,savePath:savePath}
 first,err:=app.cloudAutoSync();if err!=nil{t.Fatal(err)}
 if first["skipped"]!="not_paired"{t.Fatalf("auto synced before pairing: %+v",first)}
 preview,err:=app.cloudSyncPreview();if err!=nil{t.Fatal(err)}
 if _,err=app.cloudSyncConfirm(preview["token"].(string),map[string]bool{});err!=nil{t.Fatal(err)}
 peer:=syncSnapshot{Protocol:2,DeviceID:"peer-example",Library:initialLibrary("☑ 001. Ciudadano Kane (1941)"),
  Known:map[string]int64{"peer-example":1},Revisions:map[string]syncRevision{"ciudadano kane|1941":{
   Checked:true,Device:"peer-example",Seq:1,Clock:map[string]int64{"peer-example":1},
  }},
 }
 b,_:=json.Marshal(peer)
 remote:=filepath.Join(cloudSyncRoot(cloudDir),"sincronizacion","dispositivos","peer-example.json")
 if err=cloudWriteFile(remote,b);err!=nil{t.Fatal(err)}
 result,err:=app.cloudAutoSync();if err!=nil{t.Fatal(err)}
 if result["changed"]!=true{t.Fatalf("remote change not imported: %+v",result)}
 data,err:=os.ReadFile(path);if err!=nil{t.Fatal(err)}
 if !strings.Contains(string(data),"☑ 001."){t.Fatal("remote mark lost")}
 peerBytes,err:=os.ReadFile(remote);if err!=nil{t.Fatal(err)}
 if string(peerBytes)!=string(b){t.Fatal("overwrote peer snapshot")}
 if _,err=os.Stat(filepath.Join(localDir,"sync_recuperacion"));err!=nil{t.Fatalf("no recovery copy: %v",err)}
 result,err=app.cloudAutoSync();if err!=nil{t.Fatal(err)}
 if result["changed"]==true{t.Fatalf("reimported without changes: %+v",result)}
}

func TestCloudAutoSyncStopsAtConcurrentConflict(t *testing.T){
 dir:=t.TempDir();localDir:=filepath.Join(dir,"pc");shared:=filepath.Join(dir,"nube")
 if err:=os.MkdirAll(shared,0700);err!=nil{t.Fatal(err)}
 pending:=initialLibrary("☐ 001. Casablanca (1942)")
 viewed:=initialLibrary("☑ 001. Casablanca (1942)")
 path:=filepath.Join(localDir,"colecciones.json");savePath:=filepath.Join(localDir,"mis_365_peliculas.txt")
 if err:=saveLibrary(path,savePath,viewed);err!=nil{t.Fatal(err)}
 store,err:=newDesktopSyncStore(localDir);if err!=nil{t.Fatal(err)}
 if err=store.observe(pending,viewed);err!=nil{t.Fatal(err)}
 id,revs,known:=store.parts()
 if err=store.accept(syncPlan{Revisions:revs,Known:known},nil,shared);err!=nil{t.Fatal(err)}
 cloud:=newCloudBackups(localDir);if err=cloud.configure(shared);err!=nil{t.Fatal(err)}
 remote:=syncSnapshot{Protocol:2,DeviceID:"peer-second",Library:pending,
  Known:map[string]int64{"peer-second":1},Revisions:map[string]syncRevision{
   "casablanca|1942":{Checked:false,Device:"peer-second",Seq:1,Clock:map[string]int64{"peer-second":1}},
  }}
 b,_:=json.Marshal(remote)
 if err=cloudWriteFile(filepath.Join(cloudSyncRoot(shared),"sincronizacion","dispositivos","peer-second.json"),b);err!=nil{t.Fatal(err)}
 app:=&runningApp{syncStore:store,cloud:cloud,library:viewed,libraryPath:path,savePath:savePath}
 result,err:=app.cloudAutoSync();if err!=nil{t.Fatal(err)}
 if result["needs_review"]!=true{t.Fatalf("expected manual review: %+v",result)}
 data,_:=os.ReadFile(path)
 if !strings.Contains(string(data),"☑ 001."){t.Fatal("conflict silently overwritten")}
 if id==""{t.Fatal("missing device")}
}
