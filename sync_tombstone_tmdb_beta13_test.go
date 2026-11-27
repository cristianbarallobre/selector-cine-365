package main

import (
 "encoding/json"
 "os"
 "path/filepath"
 "testing"
)

func syncTestLists() libraryState {
 return libraryState{Version:1,ActiveID:"a",Lists:[]filmCollection{
  {ID:"a",Name:"Lista A",Text:"☐ Casablanca (1942)"},
  {ID:"b",Name:"Lista B",Text:"☑ Ciudadano Kane (1941)"},
 }}
}
func TestDeletedListNeverResurrectsFromOutdatedPeer(t *testing.T){
 old:=syncTestLists()
 deleted:=syncTestLists()
 deleted.Lists=deleted.Lists[:1]
 deleted.DeletedLists=map[string]bool{"b":true}
 peer:=syncSnapshot{Protocol:2,DeviceID:"peer-old",Library:old}
 got,_,_,conflicts,err:=syncMergeVersioned(deleted,nil,nil,peer,true)
 if err!=nil{t.Fatal(err)}
 if len(got.Lists)!=1||got.Lists[0].ID!="a"||!got.DeletedLists["b"]||len(conflicts)!=0{
  t.Fatalf("stale list resurrected: %+v, conflicts %+v",got,conflicts)
 }
 // When the other device has not yet seen the deletion, it must ALSO remove B.
 peer=syncSnapshot{Protocol:2,DeviceID:"windows-updated",Library:deleted}
 received,_,_,_,err:=syncMergeVersioned(old,nil,nil,peer,true)
 if err!=nil{t.Fatal(err)}
 if len(received.Lists)!=1||received.Lists[0].ID!="a"||!received.DeletedLists["b"]{
  t.Fatalf("peer delete not applied: %+v",received)
 }
 // A third snapshot from an obsolete device cannot recreate B after import.
 returned,_,_,_,err:=syncMergeVersioned(received,nil,nil,syncSnapshot{Protocol:2,Library:old},true)
 if err!=nil{t.Fatal(err)}
 if len(returned.Lists)!=1||returned.Lists[0].ID!="a"{t.Fatal("deleted list restored by old snapshot")}
}
func TestSharedManualTMDBOverrideAndReset(t *testing.T){
 a:=syncTestLists();b:=syncTestLists()
 a.TMDBOverrides=map[string]int{"ciudadano kane|1941":111,"casablanca|1942":0}
 b.TMDBOverrides=map[string]int{"vertigo|1958":222}
 got,_,_,_,err:=syncMergeVersioned(b,nil,nil,syncSnapshot{Protocol:2,Library:a},true)
 if err!=nil{t.Fatal(err)}
 if got.TMDBOverrides["ciudadano kane|1941"]!=111||got.TMDBOverrides["vertigo|1958"]!=222{t.Fatalf("IDs lost: %v",got.TMDBOverrides)}
 v,exists:=got.TMDBOverrides["casablanca|1942"]
 if !exists||v!=0{t.Fatal("explicit reset must persist")}
 payload,_:=json.Marshal(got)
 if string(payload)==""{t.Fatal("snapshot empty")}
 if err:=validateLibrary(got);err!=nil{t.Fatal(err)}
}
func TestTMDBKeySurvivesRestartWithoutSharingSecret(t *testing.T){
 dir:=t.TempDir()
 a:=newAvailabilityServer(dir)
 key:="test-key-12345678901234567890"
 if err:=a.saveKey(key);err!=nil{t.Fatal(err)}
 b:=newAvailabilityServer(dir)
 if !b.configured()||b.getKey()!=key{t.Fatal("saved API key lost after restart")}
 // The key must NOT occur in the cloud library JSON.
 lib:=syncTestLists();lib.TMDBOverrides=map[string]int{"casablanca|1942":4321}
 snapshot,_:=json.Marshal(lib)
 if len(snapshot)==0||string(snapshot)==key{t.Fatal("bad state")}
 if err:=b.importManualOverrides(lib.TMDBOverrides);err!=nil{t.Fatal(err)}
 again:=newAvailabilityServer(dir)
 if again.manualOverrides()["casablanca|1942"]!=4321{t.Fatal("manual mapping not persisted")}
 if _,err:=os.Stat(filepath.Join(dir,"plataformas_ajustes.json"));err!=nil{t.Fatal(err)}
}
func TestDeletedListKeepsValidSelectionAndComparison(t *testing.T){
 l:=syncTestLists();l.ActiveID="b";l.ComparisonIDs=[]string{"a","b"}
 remote:=syncTestLists();remote.Lists=remote.Lists[:1];remote.DeletedLists=map[string]bool{"b":true}
 got,_,_,_,err:=syncMergeVersioned(l,nil,nil,syncSnapshot{Protocol:2,Library:remote},true)
 if err!=nil{t.Fatal(err)}
 if got.ActiveID!="a"||len(got.ComparisonIDs)!=0||len(got.Lists)!=1{t.Fatalf("bad selection after deletion %+v",got)}
}
