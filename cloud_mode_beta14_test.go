package main

import "testing"

func TestDesktopSyncModeDefaultsAutomaticAndPersists(t *testing.T) {
 dir:=t.TempDir()
 first:=newCloudBackups(dir)
 if !first.autoMode(){t.Fatal("existing desktop users must keep automatic mode by default")}
 if err:=first.setAutoMode(false);err!=nil{t.Fatal(err)}
 second:=newCloudBackups(dir)
 if second.autoMode(){t.Fatal("manual mode must survive restart")}
 if err:=second.setAutoMode(true);err!=nil{t.Fatal(err)}
 third:=newCloudBackups(dir)
 if !third.autoMode(){t.Fatal("automatic mode must survive restart")}
}
