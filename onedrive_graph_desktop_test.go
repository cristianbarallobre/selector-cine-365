package main

import (
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "strings"
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
