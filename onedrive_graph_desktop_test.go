package main

import (
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "strings"
 "sync"
 "path/filepath"
 "io"
 "testing"
 "time"
)

func TestDesktopGraphRejectsForeignRequest(t *testing.T){
 g:=newDesktopGraph()
 _,err:=g.request("https://evil.example/","GET",nil,"")
 if err==nil||!strings.Contains(err.Error(),"no autorizada"){t.Fatalf("expected unsafe URL rejection, got %v",err)}
}
func TestDesktopGraphFolderDetectionAndReadOnlyPlan(t *testing.T){
 counts:=map[string]int{}
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.Header.Get("Authorization")!="Bearer test-access"{t.Errorf("missing scoped authorization");w.WriteHeader(401);return}
  counts[r.Method+" "+r.URL.Path]++
  w.Header().Set("Content-Type","application/json")
  var payload any=map[string]any{"value":[]any{}}
  switch r.URL.Path{
  case "/me/drive/root/children":
   payload=map[string]any{"value":[]any{map[string]any{"id":"root","name":"Selector_Cine_365","folder":map[string]any{}}}}
  case "/me/drive/items/root/children":
   payload=map[string]any{"value":[]any{map[string]any{"id":"sync","name":"sincronizacion","folder":map[string]any{}}}}
  case "/me/drive/items/sync/children":
   payload=map[string]any{"value":[]any{map[string]any{"id":"dev","name":"dispositivos","folder":map[string]any{}}}}
  case "/me/drive/items/dev/children":
   payload=map[string]any{"value":[]any{}}
  default:t.Errorf("unexpected path: %s",r.URL.Path);w.WriteHeader(404);return
  }
  _=json.NewEncoder(w).Encode(payload)
 }))
 defer server.Close()
 g:=newDesktopGraph()
 g.client=server.Client();g.graphBase=server.URL+"/"
 g.token="test-access";g.rootID="root";g.expiry=time.Now().Add(time.Hour)
 folder,err:=g.ensureRoot()
 if err!=nil||folder!="root"{t.Fatalf("existing root: %q %v",folder,err)}
 store,err:=newDesktopSyncStore(t.TempDir())
 if err!=nil{t.Fatal(err)}
 app:=&runningApp{graph:g,syncStore:store}
 lib:=initialLibrary("☐ Película de prueba (2024)")
 plan,err:=app.graphPlanFor(lib)
 if err!=nil{t.Fatal(err)}
 if plan.Sources!=0||plan.Token==""{t.Fatalf("unexpected preview: %+v",plan)}
 for k:=range counts{if strings.HasPrefix(k,"POST ")||strings.HasPrefix(k,"PUT "){t.Errorf("preview wrote cloud data: %s",k)}}
}
func TestDesktopGraphRejectsUntrustedPagination(t *testing.T){
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Content-Type","application/json")
  _,_=w.Write([]byte("{\"value\":[],\"@odata.nextLink\":\"https://evil.example/steal\"}"))
 }))
 defer server.Close()
 g:=newDesktopGraph();g.client=server.Client();g.graphBase=server.URL+"/"
 g.token="test-access";g.rootID="root";g.expiry=time.Now().Add(time.Hour)
 _,err:=g.list("me/drive/root/children")
 if err==nil||!strings.Contains(err.Error(),"no reconocida"){t.Fatalf("unsafe pagination should fail: %v",err)}
}

func TestGraphPreviewConfirmAndroidSnapshotWithoutOverwritingPeer(t *testing.T){
 remote:=initialLibrary("☑ Película compartida (2024)")
 remoteJSON,err:=json.Marshal(syncSnapshot{Protocol:2,DeviceID:"android-peer",Library:remote,Known:map[string]int64{},Revisions:map[string]syncRevision{}})
 if err!=nil{t.Fatal(err)}
 var guard sync.Mutex
 var ownBytes []byte
 ownName:=""
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  guard.Lock();defer guard.Unlock()
  w.Header().Set("Content-Type","application/json")
  switch {
  case r.Method=="GET"&&r.URL.Path=="/me/drive/items/root/children":
   _,_=w.Write([]byte("{\"value\":[{\"id\":\"sync\",\"name\":\"sincronizacion\",\"folder\":{}}]}"))
  case r.Method=="GET"&&r.URL.Path=="/me/drive/items/sync/children":
   _,_=w.Write([]byte("{\"value\":[{\"id\":\"dev\",\"name\":\"dispositivos\",\"folder\":{}}]}"))
  case r.Method=="GET"&&r.URL.Path=="/me/drive/items/dev/children":
   items:=[]map[string]any{{"id":"remote","name":"android-peer.json","file":map[string]any{}}}
   if ownBytes!=nil{items=append(items,map[string]any{"id":"own","name":ownName,"file":map[string]any{}})}
   _=json.NewEncoder(w).Encode(map[string]any{"value":items})
  case r.Method=="GET"&&r.URL.Path=="/me/drive/items/remote/content":
   _,_=w.Write(remoteJSON)
  case r.Method=="PUT"&&strings.HasPrefix(r.URL.Path,"/me/drive/items/dev:/")&&strings.HasSuffix(r.URL.Path,":/content"):
   // This test refuses any write to the Android peer snapshot.
   ownName=strings.TrimSuffix(strings.TrimPrefix(r.URL.Path,"/me/drive/items/dev:/"),":/content")
   ownBytes,_=io.ReadAll(r.Body)
   _,_=w.Write([]byte("{\"id\":\"own\"}"))
  case r.Method=="GET"&&r.URL.Path=="/me/drive/items/own/content":
   _,_=w.Write(ownBytes)
  default:
   t.Errorf("unexpected Graph operation %s %s",r.Method,r.URL.Path)
   w.WriteHeader(404)
  }
 }))
 defer server.Close()
 g:=newDesktopGraph()
 g.client=server.Client();g.graphBase=server.URL+"/"
 g.token="test-access";g.rootID="root";g.expiry=time.Now().Add(time.Hour)
 dir:=t.TempDir()
 store,err:=newDesktopSyncStore(dir);if err!=nil{t.Fatal(err)}
 lib:=initialLibrary("☐ Película compartida (2024)")
 app:=&runningApp{graph:g,syncStore:store,library:lib,libraryPath:filepath.Join(dir,"colecciones.json"),savePath:filepath.Join(dir,"peliculas.txt")}
 if err=saveLibrary(app.libraryPath,app.savePath,lib);err!=nil{t.Fatal(err)}
 preview,err:=app.graphSyncPreview();if err!=nil{t.Fatal(err)}
 if n,ok:=preview["sources"].(int);!ok||n!=1{t.Fatalf("expected one Android snapshot, got %v",preview["sources"])}
 if ownBytes!=nil{t.Fatal("preview must not upload")}
 _,err=app.graphSyncConfirm(preview["token"].(string),map[string]bool{})
 if err!=nil{t.Fatal(err)}
 if ownBytes==nil{t.Fatal("confirmation did not publish own-device snapshot")}
 if ownName=="android-peer.json"{t.Fatal("must never replace another device snapshot")}
 if syncWatched(app.library)!=1{t.Fatalf("Android watched flag not merged: %+v",app.library)}
}
