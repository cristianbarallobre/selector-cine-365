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
 data,e:=os.ReadFile(filepath.Join(filepath.Dir(b.target(b.config)),"historial_actividad.json"))
 if e!=nil{t.Fatal(e)}
 var v struct{Version int;Events []activityEntry}
 if json.Unmarshal(data,&v)!=nil||v.Version!=1{t.Fatal(string(data))}
}
