package main

import (
 "encoding/json"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestBeta10RestoreKeepsBeforeBackupAndRecordsUnwatch(t *testing.T){
 folder:=t.TempDir()
 active:=initialLibrary("☑ 001. Casablanca (1942)")
 old:=initialLibrary("☐ 001. Casablanca (1942)")
 state:=filepath.Join(folder,"colecciones.json")
 txt:=filepath.Join(folder,"mis_365_peliculas.txt")
 if err:=saveLibrary(state,txt,active);err!=nil{t.Fatal(err)}
 src:=filepath.Join(folder,"sync_recuperacion","20261120T100000.000000000.json")
 b,_:=json.Marshal(old)
 if err:=cloudWriteFile(src,b);err!=nil{t.Fatal(err)}
 store,err:=newDesktopSyncStore(folder);if err!=nil{t.Fatal(err)}
 a:=&runningApp{libraryPath:state,savePath:txt,library:active,syncStore:store}
 before,err:=a.localRecoveryList();if err!=nil{t.Fatal(err)}
 if len(before)==0||before[0].ID==""{t.Fatalf("no recovery choices: %+v",before)}
 result,err:=a.restoreLocalRecovery("rec:20261120T100000.000000000.json")
 if err!=nil{t.Fatal(err)}
 if result["restored"]!=true||syncWatched(a.library)!=0{t.Fatalf("not restored: %+v",result)}
 if syncStatuses(a.library)["casablanca|1942"]{t.Fatal("restore did not clear mark")}
 _,revs,_:=store.parts()
 revision,ok:=revs["casablanca|1942"];if !ok||revision.Checked{t.Fatalf("restore not recorded as local unwatch: %+v",revs)}
 history:=store.history()
 if len(history)==0||history[0].Origin!="local"||history[0].Checked{t.Fatalf("restore not in change history: %+v",history)}
 entries,err:=os.ReadDir(filepath.Join(folder,"sync_recuperacion"));if err!=nil{t.Fatal(err)}
 if len(entries)<2{t.Fatal("restore failed to create new rescue snapshot")}
 foundPrevious:=false
 for _,entry:=range entries{
  if strings.Contains(entry.Name(),"antes_restaurar"){
   snap,e:=os.ReadFile(filepath.Join(folder,"sync_recuperacion",entry.Name()))
   if e!=nil{t.Fatal(e)}
   var previous libraryState
   if e=json.Unmarshal(snap,&previous);e!=nil{t.Fatal(e)}
   foundPrevious=syncWatched(previous)==1
  }
 }
 if !foundPrevious{t.Fatal("missing rescue with original watched state")}
}

func TestBeta10RecoveryRejectsArbitraryPathsAndCorruptSnapshots(t *testing.T){
 dir:=t.TempDir();a:=&runningApp{libraryPath:filepath.Join(dir,"colecciones.json")}
 for _,id:=range []string{"../colecciones.json","rec:../outside.json","day:../../secrets.json","oops:20261120.json","rec:a.json"}{
  if _,err:=a.restoreLocalRecovery(id);err==nil{t.Fatalf("unsafe recovery id accepted: %q",id)}
 }
 bad:=filepath.Join(dir,"sync_recuperacion","20261120_bad.json")
 if err:=cloudWriteFile(bad,[]byte("{invalid"));err!=nil{t.Fatal(err)}
 if _,err:=a.restoreLocalRecovery("rec:20261120_bad.json");err==nil{t.Fatal("accepted corrupt JSON")}
}
