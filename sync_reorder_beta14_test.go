package main

import (
 "strings"
 "testing"
)

func TestFilmOrderDoesNotAlterRowsOrMarks(t *testing.T){
 current:="Mi lista\n☐ 001. Ciudadano Kane (1941)\n\n☑ 002. Vértigo (1958)\n☐ Otra (2002)"
 remote:="☑ 002. Vértigo (1958)\n☐ 001. Ciudadano Kane (1941)"
 got:=syncOrderedFilmRows(current,remote)
 want:="Mi lista\n☑ 002. Vértigo (1958)\n\n☐ 001. Ciudadano Kane (1941)\n☐ Otra (2002)"
 if got!=want{t.Fatalf("wrong order:\n%s\nwanted:\n%s",got,want)}
}
func TestFilmOrderSyncMergesNewerRemoteOrder(t *testing.T){
 local:=initialLibrary("☐ Uno (2001)\n☐ Dos (2002)\n☐ Tres (2003)")
 local.Lists[0].OrderChangedAt=100
 remote:=initialLibrary("☐ Tres (2003)\n☐ Uno (2001)\n☐ Dos (2002)")
 remote.Lists[0].OrderChangedAt=200
 merged,_,_,_,err:=syncMerge(local,nil,nil,syncSnapshot{Protocol:syncProtocol,DeviceID:"android",Library:remote})
 if err!=nil{t.Fatal(err)}
 if merged.Lists[0].Text!=remote.Lists[0].Text{t.Fatalf("newer order not accepted: %q",merged.Lists[0].Text)}
 if merged.Lists[0].OrderChangedAt!=200{t.Fatalf("lost order timestamp")}
 // An old snapshot may add new films, but cannot revert the order.
 stale:=initialLibrary("☐ Dos (2002)\n☐ Uno (2001)\n☐ Tres (2003)\n☐ Cuatro (2004)")
 stale.Lists[0].OrderChangedAt=150
 same,_,_,_,err:=syncMerge(merged,nil,nil,syncSnapshot{Protocol:syncProtocol,DeviceID:"old",Library:stale})
 if err!=nil{t.Fatal(err)}
 if !strings.HasPrefix(same.Lists[0].Text,"☐ Tres (2003)\n☐ Uno (2001)\n☐ Dos (2002)"){t.Fatalf("stale reorder changed local order: %q",same.Lists[0].Text)}
 if !strings.Contains(same.Lists[0].Text,"Cuatro"){t.Fatal("lost remote addition")}
}
