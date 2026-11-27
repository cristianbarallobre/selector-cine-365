//go:build windows

package main

import (
 "bytes"
 "os"
 "path/filepath"
 "testing"
)

func TestWindowsDPAPIUserBoundSessionRoundTrip(t *testing.T){
 dir:=t.TempDir()
 store:=protectedSessionStore(dir)
 secret:="windows-only-test-refresh-token-THIS-MUST-BE-ENCRYPTED"
 original:=graphStoredSession{
  Version:graphSessionVersion,ClientID:desktopMicrosoftAppID,
  RootID:"test-folder-id",Refresh:secret,
 }
 if err:=store.save(original);err!=nil{t.Fatalf("DPAPI could not protect session: %v",err)}
 path:=filepath.Join(dir,"onedrive_sesion_dpapi.bin")
 onDisk,err:=os.ReadFile(path)
 if err!=nil{t.Fatal(err)}
 if len(onDisk)==0||bytes.Contains(onDisk,[]byte(secret))||bytes.Contains(onDisk,[]byte(original.RootID)){
  t.Fatal("DPAPI stored a token or drive id in plaintext")
 }
 got,err:=store.load()
 if err!=nil{t.Fatalf("Windows user could not recover own session: %v",err)}
 if got!=original{t.Fatal("Windows DPAPI session changed values during restart")}
 if err=store.clear();err!=nil{t.Fatal(err)}
 if _,err=os.Stat(path);!os.IsNotExist(err){t.Fatalf("disconnect must delete protected session file: %v",err)}
 if _,err=store.load();!os.IsNotExist(err){t.Fatalf("cleared session must not be recovered: %v",err)}
}

func TestWindowsDPAPIRejectsTampering(t *testing.T){
 dir:=t.TempDir()
 store:=protectedSessionStore(dir)
 if err:=store.save(graphStoredSession{Version:graphSessionVersion,ClientID:desktopMicrosoftAppID,RootID:"folder",Refresh:"private-refresh"});err!=nil{t.Fatal(err)}
 path:=filepath.Join(dir,"onedrive_sesion_dpapi.bin")
 data,err:=os.ReadFile(path);if err!=nil{t.Fatal(err)}
 if len(data)<16{t.Fatal("encrypted DPAPI record unexpectedly short")}
 // A malformed encrypted record must not be treated as another account.
 data=data[:8]
 if err=os.WriteFile(path,data,0600);err!=nil{t.Fatal(err)}
 if _,err=store.load();err==nil{t.Fatal("corrupted protected token was accepted")}
}
