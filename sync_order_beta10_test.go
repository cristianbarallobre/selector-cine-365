package main

import (
 "testing"
)

func TestBeta10PairedIgnoresUnversionedSeenSnapshot(t *testing.T){
 key:="casablanca|1942"
 pending:=initialLibrary("☐ 001. Casablanca (1942)")
 watched:=initialLibrary("☑ 001. Casablanca (1942)")
 ours:=syncRevision{Checked:false,Device:"remote",Seq:3,Clock:map[string]int64{"windows":1,"remote":3}}
 peer:=syncSnapshot{Protocol:2,DeviceID:"old-beta",Library:watched,Known:map[string]int64{},Revisions:map[string]syncRevision{}}
 got,revs,_,conflicts,err:=syncMergeVersioned(pending,map[string]syncRevision{key:ours},map[string]int64{"windows":1,"remote":3},peer,true)
 if err!=nil{t.Fatal(err)}
 if syncWatched(got)!=0||revs[key].Checked||len(conflicts)!=0{t.Fatalf("historic seen resurrected: %+v %+v",revs,conflicts)}
}
func TestBeta10PairedConcurrentStillConflicts(t *testing.T){
 key:="casablanca|1942"
 watched:=initialLibrary("☑ 001. Casablanca (1942)")
 pending:=initialLibrary("☐ 001. Casablanca (1942)")
 ours:=syncRevision{Checked:true,Device:"windows",Seq:2,Clock:map[string]int64{"windows":2}}
 theirs:=syncRevision{Checked:false,Device:"remote",Seq:1,Clock:map[string]int64{"remote":1}}
 peer:=syncSnapshot{Protocol:2,DeviceID:"remote",Library:pending,Revisions:map[string]syncRevision{key:theirs}}
 _,_,_,conflicts,err:=syncMergeVersioned(watched,map[string]syncRevision{key:ours},map[string]int64{"windows":2},peer,true)
 if err!=nil{t.Fatal(err)}
 if len(conflicts)!=1{t.Fatalf("concurrent opposite edits must be reviewed: %+v",conflicts)}
}
