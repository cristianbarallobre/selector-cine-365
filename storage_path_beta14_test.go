package main
import (
 "os"
 "path/filepath"
 "testing"
)
func TestFreshDataFolderHasCurrentProductName(t *testing.T) {
 path:=getSavePath()
 if got:=filepath.Base(filepath.Dir(path));got!="Selector_Cine_365" {t.Fatalf("unexpected storage folder: %q",got)}
 if filepath.Base(path)!="mis_365_peliculas.txt"{t.Fatalf("wrong TXT filename: %q",path)}
 if got:=filepath.Base(filepath.Dir(filepath.Join(filepath.Dir(path),"colecciones.json")));got!="Selector_Cine_365"{t.Fatal("library folder differs from main data folder")}
 if _,err:=os.Stat(filepath.Join(filepath.Dir(path),"../AtlasDelCine"));err==nil{
  // No migration: preexisting legacy folder must be left alone and ignored.
  t.Log("Legacy files remain untouched; fresh installations use only Selector_Cine_365")
 }
}
