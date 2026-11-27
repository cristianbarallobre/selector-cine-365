package main
import("encoding/json";"os";"path/filepath";"testing")
func TestActivitySnapshotCommitAndRecovery(t *testing.T){
 if !desktopFeatures{t.Skip("desktop")}
 a:=durableTestApp(t);a.extraBackup=newExtraBackup(filepath.Dir(a.libraryPath))
 if _,_,e:=a.commitProtectedLibrary(a.library,nil,nil,"");e!=nil{t.Fatal(e)}
 next:=initialLibrary("☑ Origen (2010)")
 if _,_,e:=a.commitProtectedLibrary(next,nil,nil,"");e!=nil{t.Fatal(e)}
 raw,e:=os.ReadFile(a.libraryPath);if e!=nil{t.Fatal(e)}
 var disk libraryDisk;if e=json.Unmarshal(raw,&disk);e!=nil||disk.ActivityRef==""||len(disk.LocalSync.Activity)!=0{t.Fatal("history is not separated",e)}
 cp:=readLibraryCheckpoint(a.libraryPath);if cp==nil||len(cp.Activity)!=1||cp.Activity[0].Kind!="seen"{t.Fatal(cp)}
 root:=filepath.Dir(a.libraryPath)
 path,e:=activitySnapshotPath(root,disk.ActivityRef);if e!=nil{t.Fatal(e)}
 if e=os.Remove(path);e!=nil{t.Fatal(e)}
 // Mirror can recover the exact committed activity, even if snapshot is gone.
 if cp=readLibraryCheckpoint(a.libraryPath);cp==nil||len(cp.Activity)!=1{t.Fatal("missing fallback")}
 if e=os.Remove(a.syncStore.path);e!=nil{t.Fatal(e)}
 if _,_,e=labLoadLibraryProtected(a.libraryPath,"");e==nil{t.Fatal("must not silently discard history")}
 after,_:=os.ReadFile(a.libraryPath);if string(after)!=string(raw){t.Fatal("changed original")}
}
func TestActivitySnapshotBackupSelfContained(t *testing.T){
 if !desktopFeatures{t.Skip("desktop")}
 a:=durableTestApp(t);a.extraBackup=newExtraBackup(filepath.Dir(a.libraryPath))
 if _,_,e:=a.commitProtectedLibrary(a.library,nil,nil,"");e!=nil{t.Fatal(e)}
 raw,e:=os.ReadFile(filepath.Join(filepath.Dir(a.libraryPath),"biblioteca_guardada.json"));if e!=nil{t.Fatal(e)}
 var disk libraryDisk;if json.Unmarshal(raw,&disk)!=nil||disk.ActivityRef!=""||disk.LocalSync==nil{t.Fatal("backup references an external snapshot")}
 _,history,e:=decodeRecoveryBackup(raw);if e!=nil||history==nil{t.Fatal("backup not portable",e)}
}
