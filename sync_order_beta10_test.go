package main

import (
 "encoding/json"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestBeta10PairedUnwatchOverridesOldWindowsWatch(t *testing.T){
 root:=t.TempDir();localDir:=filepath.Join(root,"desktop");shared:=filepath.Join(root,"OneDrive")
 if err:=os.MkdirAll(shared,0700);err!=nil{t.Fatal(err)}
 key:="casablanca|1942"
 pending:=initialLibrary("☐ 001. Casablanca (1942)")
 watched:=initialLibrary("☑ 001. Casablanca (1942)")
 libraryPath:=filepath.Join(localDir,"colecciones.json")
 txt:=filepath.Join(localDir,"mis_365_peliculas.txt")
 if err:=saveLibrary(libraryPath,txt,pending);err!=nil{t.Fatal(err)}
 store,err:=newDesktopSyncStore(localDir);if err!=nil{t.Fatal(err)}
 if err=store.observe(pending,watched);err!=nil{t.Fatal(err)}
 _,winRevs,winKnown:=store.parts()
 if !winRevs[key].Checked{t.Fatal("Windows did not record first viewed action")}
 cloud:=newCloudBackups(localDir);if err=cloud.configure(shared);err!=nil{t.Fatal(err)}
 app:=&runningApp{libraryPath:libraryPath,savePath:txt,library:watched,cloud:cloud,syncStore:store}
 if err=saveLibrary(libraryPath,txt,watched);err!=nil{t.Fatal(err)}
 if err=store.accept(syncPlan{Revisions:winRevs,Known:winKnown},nil,shared);err!=nil{t.Fatal(err)}
 // An old "actual" copy in OneDrive still has the Windows checked mark.
 old,_:=json.Marshal(watched)
 if err=cloudWriteFile(filepath.Join(cloudSyncRoot(shared),"actual","colecciones.json"),old);err!=nil{t.Fatal(err)}
 // Android already received Windows' revision, so its following unwatch is causal.
 androidPending:=syncSnapshot{Protocol:2,DeviceID:"android-sim",Library:pending,
  Known:map[string]int64{store.state.DeviceID:1,"android-sim":1},
  Revisions:map[string]syncRevision{key:{
   Checked:false,Device:"android-sim",Seq:1,ChangedAt:2000,
   Clock:map[string]int64{store.state.DeviceID:1,"android-sim":1},
  }},
 }
 b,_:=json.Marshal(androidPending)
 target:=filepath.Join(cloudSyncRoot(shared),"sincronizacion","dispositivos","android-sim.json")
 if err=cloudWriteFile(target,b);err!=nil{t.Fatal(err)}
 result,err:=app.cloudAutoSync();if err!=nil{t.Fatal(err)}
 if result["needs_review"]==true{t.Fatal("causal unwatch falsely treated as conflict")}
 if result["changed"]!=true{t.Fatalf("newer unwatch was not imported: %+v",result)}
 if syncWatched(app.library)!=0{t.Fatal("stale Windows mark resurrected after Android unwatch")}
 bytes,_:=os.ReadFile(libraryPath)
 if !strings.Contains(string(bytes),"☐ 001."){t.Fatal("unwatch was not persisted")}
 result,err=app.cloudAutoSync();if err!=nil{t.Fatal(err)}
 if result["changed"]==true{t.Fatal("stale legacy watch returned on second sync")}
}
func TestBeta10PairedIgnoresUnversionedSeenSnapshot(t *testing.T){
 key:="casablanca|1942"
 pending:=initialLibrary("☐ 001. Casablanca (1942)")
 watched:=initialLibrary("☑ 001. Casablanca (1942)")
 ours:=syncRevision{Checked:false,Device:"android",Seq:3,Clock:map[string]int64{"windows":1,"android":3}}
 peer:=syncSnapshot{Protocol:2,DeviceID:"old-beta",Library:watched,Known:map[string]int64{},Revisions:map[string]syncRevision{}}
 got,revs,_,conflicts,err:=syncMergeVersioned(pending,map[string]syncRevision{key:ours},map[string]int64{"windows":1,"android":3},peer,true)
 if err!=nil{t.Fatal(err)}
 if syncWatched(got)!=0||revs[key].Checked||len(conflicts)!=0{t.Fatalf("historic seen resurrected: %+v %+v",revs,conflicts)}
}
func TestBeta10PairedConcurrentStillConflicts(t *testing.T){
 key:="casablanca|1942"
 watched:=initialLibrary("☑ 001. Casablanca (1942)")
 pending:=initialLibrary("☐ 001. Casablanca (1942)")
 ours:=syncRevision{Checked:true,Device:"windows",Seq:2,Clock:map[string]int64{"windows":2}}
 theirs:=syncRevision{Checked:false,Device:"android",Seq:1,Clock:map[string]int64{"android":1}}
 peer:=syncSnapshot{Protocol:2,DeviceID:"android",Library:pending,Revisions:map[string]syncRevision{key:theirs}}
 _,_,_,conflicts,err:=syncMergeVersioned(watched,map[string]syncRevision{key:ours},map[string]int64{"windows":2},peer,true)
 if err!=nil{t.Fatal(err)}
 if len(conflicts)!=1{t.Fatalf("concurrent opposite edits must be reviewed: %+v",conflicts)}
}
