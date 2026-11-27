package main
import("encoding/json";"os";"path/filepath";"testing")
func TestExactHistoryRepairFromPortableCopy(t *testing.T){
 if !desktopFeatures{t.Skip("desktop")}
 a:=durableTestApp(t);root:=filepath.Dir(a.libraryPath)
 a.extraBackup=newExtraBackup(root)
 if _,_,e:=a.commitProtectedLibrary(a.library,nil,nil,"");e!=nil{t.Fatal(e)}
 next:=initialLibrary("☑ Origen (2010)")
 if _,_,e:=a.commitProtectedLibrary(next,nil,nil,"");e!=nil{t.Fatal(e)}
 a.library=next
 if e:=a.extraBackup.configure(t.TempDir());e!=nil{t.Fatal(e)}
 a.extraBackup.copy(a)
 raw,e:=os.ReadFile(a.libraryPath);if e!=nil{t.Fatal(e)}
 var disk libraryDisk;if e=json.Unmarshal(raw,&disk);e!=nil{t.Fatal(e)}
 path,e:=activitySnapshotPath(root,disk.ActivityRef);if e!=nil{t.Fatal(e)}
 if e=os.WriteFile(path,[]byte("corrupt"),0600);e!=nil{t.Fatal(e)}
 if e=os.Remove(a.syncStore.path);e!=nil{t.Fatal(e)}
 state,_,e:=labLoadLibraryProtected(a.libraryPath,"")
 if e!=nil||state.Lists[0].Text!=next.Lists[0].Text{t.Fatal(state,e)}
 cp:=readLibraryCheckpoint(a.libraryPath);if cp==nil||len(cp.Activity)!=1||cp.Activity[0].Kind!="seen"{t.Fatal(cp)}
 current,_:=os.ReadFile(a.libraryPath);if string(current)!=string(raw){t.Fatal("library replaced")}
 damaged,e:=os.ReadFile(filepath.Join(filepath.Dir(path),disk.ActivityRef+"_danada_"+libraryBytesFingerprint([]byte("corrupt"))+".json"))
 if e!=nil||string(damaged)!="corrupt"{t.Fatal("damaged original not preserved",e)}
}
