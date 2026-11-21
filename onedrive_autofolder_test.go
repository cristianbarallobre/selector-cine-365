package main

import (
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestDiscoverWindowsOneDriveWithoutManualFolder(t *testing.T) {
 home:=t.TempDir()
 root:=filepath.Join(home,"OneDrive")
 if err:=os.Mkdir(root,0700);err!=nil{t.Fatal(err)}
 folder,err:=discoverOneDriveFolder(home,"windows",func(k string)string{
  if k=="OneDriveConsumer"{return root}
  return ""
 })
 if err!=nil{t.Fatal(err)}
 want:=filepath.Join(root,"Selector_Cine_365")
 if folder!=want{t.Fatalf("folder=%q want %q",folder,want)}
 if _,err:=os.Stat(want);!os.IsNotExist(err){t.Fatalf("discovery must never create or upload: %v",err)}
}

func TestDiscoverPreservesExistingFolder(t *testing.T) {
 home:=t.TempDir()
 root:=filepath.Join(home,"OneDrive")
 want:=filepath.Join(root,"Selector_Cine_365")
 if err:=os.MkdirAll(want,0700);err!=nil{t.Fatal(err)}
 folder,err:=discoverOneDriveFolder(home,"windows",func(string)string{return root})
 if err!=nil||folder!=want{t.Fatalf("folder=%q err=%v",folder,err)}
}

func TestDiscoverRejectsAmbiguousAccounts(t *testing.T) {
 home:=t.TempDir()
 one,two:=filepath.Join(home,"First"),filepath.Join(home,"Second")
 for _,path:=range []string{one,two}{if err:=os.MkdirAll(path,0700);err!=nil{t.Fatal(err)}}
 _,err:=discoverOneDriveFolder(home,"windows",func(k string)string{
  if k=="OneDriveConsumer"{return one}
  if k=="OneDriveCommercial"{return two}
  return ""
 })
 if err==nil||!strings.Contains(err.Error(),"cuenta"){t.Fatalf("expected ambiguous-account refusal: %v",err)}
}

func TestDiscoverMacCloudStorage(t *testing.T) {
 home:=t.TempDir()
 root:=filepath.Join(home,"Library","CloudStorage","OneDrive-Personal")
 if err:=os.MkdirAll(root,0700);err!=nil{t.Fatal(err)}
 folder,err:=discoverOneDriveFolder(home,"darwin",func(string)string{return ""})
 if err!=nil||folder!=filepath.Join(root,"Selector_Cine_365"){t.Fatalf("folder=%q err=%v",folder,err)}
}

func TestDiscoverDoesNotChooseUnknownPlatform(t *testing.T) {
 _,err:=discoverOneDriveFolder(t.TempDir(),"linux",func(string)string{return ""})
 if err==nil{t.Fatal("must not pick an unsupported OS")}
}
