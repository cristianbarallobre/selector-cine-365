package main

import (
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "path/filepath"
 "strings"
 "testing"
)
func TestActivityChangesAndMigration(t *testing.T){
 a:=initialLibrary("☐ Uno (2001)\n☑ Dos (2002)")
 b:=a;b.Lists=append([]filmCollection(nil),a.Lists...)
 b.Lists[0].Name="Nueva";b.Lists[0].Text="☑ Uno (2001)\n☐ Tres (2003)"
 b.Lists=append(b.Lists,filmCollection{ID:"otra",Name:"Otra",Text:"☐ Uno (2001)"})
 b.Reactions=map[string]filmReaction{"uno|2001":{Title:"Uno",Year:"2001",Favorite:true,Rating:4}}
 b.TMDBOverrides=map[string]int{"uno|2001":123}
 state:=syncDeviceState{History:[]syncHistoryEntry{{At:1,Key:"dos|2002",Checked:true,Origin:"local"}}}
 appendLocalActivity(&state,a,b)
 kinds:=map[string]int{};for _,e:=range state.Activity{kinds[e.Kind]++}
 for _,k:=range []string{"list_rename","list_create","seen","film_add","film_remove","favorite_add","rating","identity"}{if kinds[k]==0{t.Fatal("missing",k,state.Activity)}}
 if kinds["seen"]!=2{t.Fatal("old activity lost or duplicated")}
 n:=len(state.Activity);c:=b;c.ActiveID="otra";appendLocalActivity(&state,b,c)
 if len(state.Activity)!=n{t.Fatal("navigation recorded")}
 appendLocalActivity(&state,c,c);if len(state.Activity)!=n{t.Fatal("identical save recorded")}
 for i:=0;i<activityLimit+10;i++{appendActivityEvent(&state,activityEntry{At:int64(i),Kind:"seen"})}
 if len(state.Activity)!=activityLimit{t.Fatal("retention",len(state.Activity))}
}
func TestActivityCommitAndRestart(t *testing.T){
 if !desktopFeatures{t.Skip("desktop")}
 a:=durableTestApp(t);a.extraBackup=newExtraBackup(filepath.Dir(a.libraryPath))
 a.syncStore.state.History=[]syncHistoryEntry{{At:1,Key:"anterior|2000",Checked:true,Origin:"local"}}
 next:=a.library;next.Lists=append([]filmCollection(nil),next.Lists...);next.Lists[0].Text="☐ Origen (2010)\n☐ Nueva (2020)"
 if w:=postTestLibrary(t,a,next);w.Code!=204{t.Fatal(w.Code,w.Body.String())}
 events:=a.syncStore.localActivity()
 if len(events)!=2||events[0].Kind!="film_add"||events[1].Film!="anterior · 2000"{t.Fatal(events)}
 cp:=readLibraryCheckpoint(a.libraryPath);if cp==nil||len(cp.Activity)!=2{t.Fatal("history not committed")}
 restarted,err:=recoverSyncCheckpoint(filepath.Dir(a.libraryPath),cp);if err!=nil||len(restarted.localActivity())!=2{t.Fatal(err)}
 req:=httptest.NewRequest(http.MethodGet,"/api/sync/history",nil);req.Header.Set("X-Atlas-Token","test")
 w:=httptest.NewRecorder();a.serveHTTP(w,req)
 var response struct{Events []activityEntry}
 if w.Code!=200||json.Unmarshal(w.Body.Bytes(),&response)!=nil||len(response.Events)!=2{t.Fatal(w.Code,w.Body.String())}
 // Force a failed save before any history can be committed.
 a.libraryDigest="not-current"
 changed:=next;changed.Lists=append([]filmCollection(nil),next.Lists...);changed.Lists[0].Text+="\n☐ Fallo (2021)"
 if w:=postTestLibrary(t,a,changed);w.Code==204{t.Fatal("unexpected success")}
 if len(a.syncStore.localActivity())!=2{t.Fatal("failed save recorded")}
 page,err:=makeHTML("☐ Origen (2010)","test");if err!=nil{t.Fatal(err)}
 page,err=localOnlyHTML(page);if err!=nil{t.Fatal(err)}
 for _,marker:=range []string{"id=\"activitySearch\"","id=\"activityFilter\"","id=\"activityPrev\"","id=\"activityNext\"","máximo 10.000"}{if !strings.Contains(page,marker){t.Fatal("missing",marker)}}
}
