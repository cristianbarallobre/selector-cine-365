package main
import("testing";"path/filepath";"os";"strings";"runtime")
func TestTMDBNativeCredentialMigration(t *testing.T){
 if !desktopFeatures{t.Skip("official desktop")}
 path:=filepath.Join(t.TempDir(),"plataformas_ajustes.json")
 key:=strings.Repeat("a",32)
 if e:=writeTMDBSettings(path,platformSettings{Key:key});e!=nil{t.Fatal(e)}
 got,e:=readTMDBKey(path);if e!=nil||got!=key{t.Fatal("migration failed",e)}
 raw,e:=os.ReadFile(path);if e!=nil{t.Fatal(e)}
 if strings.Contains(string(raw),key){t.Fatal("plaintext remained in settings on "+runtime.GOOS)}
 got,e=readTMDBKey(path);if e!=nil||got!=key{t.Fatal("native key did not reopen",e)}
 if e=writeTMDBKey(path,"");e!=nil{t.Fatal("disconnect",e)}
 got,e=readTMDBKey(path);if e!=nil||got!=""{t.Fatal("credential remained configured",e)}
 fresh:=strings.Repeat("b",32);if e=writeTMDBKey(path,fresh);e!=nil{t.Fatal(e)}
 got,e=readTMDBKey(path);if e!=nil||got!=fresh{t.Fatal("reconnect failed",e)}
 if e=writeTMDBKey(path,"");e!=nil{t.Fatal(e)}
}
