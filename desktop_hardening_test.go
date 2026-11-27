package main
import("testing";"time";"path/filepath";"os";"encoding/json")
func TestProviderCacheExpiryHardening(t *testing.T){
 now:=time.Now()
 for _,tc:=range []struct{age time.Duration;state string;expired bool}{{time.Hour,"ok",false},{25*time.Hour,"ok",true},{16*time.Minute,"error",true},{5*time.Minute,"error",false}}{
  r:=providerResult{Date:now.Add(-tc.age).Format(time.RFC3339),State:tc.state}
  if providerCacheExpired(r,now)!=tc.expired{t.Fatal(tc)}
 }
 p:=&availabilityServer{cache:map[string]providerResult{"x":{State:"ok",TMDBID:123,MetadataVersion:3,Subscriptions:[]string{"Test"}}}}
 if r:=p.finishAvailability("x",0,providerResult{State:"error",TMDBID:123});r.State!="ok"||len(r.Subscriptions)!=1{t.Fatal(r)}
}
func TestExtraActivityArchiveHardening(t *testing.T){
 a:=durableTestApp(t);b:=newExtraBackup(filepath.Dir(a.libraryPath))
 if e:=b.configure(t.TempDir());e!=nil{t.Fatal(e)}
 b.copy(a)
 data,e:=os.ReadFile(b.target(b.config))
 if e!=nil{t.Fatal(e)}
 var v completeLibraryBackup
 if json.Unmarshal(data,&v)!=nil||v.CompleteVersion!=1{t.Fatal(string(data))}
 if _,e=os.Stat(filepath.Join(filepath.Dir(b.target(b.config)),"historial_actividad.json"));!os.IsNotExist(e){t.Fatal("duplicate activity archive")}
}

func TestRestoreRecoveryCompleteHistory(t *testing.T){
 if !desktopFeatures{t.Skip("desktop")}
 a:=durableTestApp(t);a.extraBackup=newExtraBackup(filepath.Dir(a.libraryPath))
 dir:=filepath.Join(filepath.Dir(a.libraryPath),"sync_recuperacion");if e:=os.MkdirAll(dir,0700);e!=nil{t.Fatal(e)}
 name:="20261127T150000.000000000Z_antes_restaurar.json"
 b:=completeLibraryBackup{libraryState:a.library,CompleteVersion:1,Activity:[]activityEntry{{At:1,Kind:"film_add",Film:"Recuperada"}}}
 raw,_:=json.Marshal(b);if e:=os.WriteFile(filepath.Join(dir,name),raw,0600);e!=nil{t.Fatal(e)}
 preview,e:=a.previewLocalRecovery("rec:"+name);if e!=nil||preview["has_history"]!=true{t.Fatal(preview,e)}
 if _,e=a.restoreLocalRecovery("rec:"+name,a.libraryETag());e!=nil{t.Fatal(e)}
 events:=a.syncStore.localActivity();if len(events)!=2||events[1].Film!="Recuperada"{t.Fatal(events)}
}
