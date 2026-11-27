//go:build windows || darwin

package main
import("testing";"os";"path/filepath")
func TestVerifiedLegacyCleanupDoesNotTouchPersonalData(t *testing.T){
 root,cache:=t.TempDir(),t.TempDir()
 for _,pair:=range []struct{name,old,new string}{{"plataformas_cache.json","identical","identical"},{"tmdb_fichas_lab.json","old","new"},{"colecciones.json","library","library"}}{
  if e:=os.WriteFile(filepath.Join(root,pair.name),[]byte(pair.old),0600);e!=nil{t.Fatal(e)}
  if e:=os.WriteFile(filepath.Join(cache,pair.name),[]byte(pair.new),0600);e!=nil{t.Fatal(e)}
 }
 n,e:=cineCleanVerifiedLegacyCache(root,cache)
 if e!=nil||n!=1{t.Fatal(n,e)}
 for _,name:=range []string{"colecciones.json","tmdb_fichas_lab.json"}{if _,e:=os.Stat(filepath.Join(root,name));e!=nil{t.Fatal(name,e)}}
 if _,e:=os.Stat(filepath.Join(root,"plataformas_cache.json"));!os.IsNotExist(e){t.Fatal(e)}
}
