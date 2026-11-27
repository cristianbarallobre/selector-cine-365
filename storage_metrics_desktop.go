//go:build windows || darwin

package main
import("bytes";"io";"os";"path/filepath";"io/fs")
func cineSumUsage(items ...labFileUsage)labFileUsage{
 var total labFileUsage;for _,u:=range items{total.Files+=u.Files;total.Bytes+=u.Bytes};return total
}
func cineDirectoryUsage(root string)labFileUsage{
 var total labFileUsage
 _=filepath.WalkDir(root,func(path string,d fs.DirEntry,e error)error{
  if e!=nil{return nil}
  if d.Type().IsRegular(){if info,e:=d.Info();e==nil{total.Files++;total.Bytes+=info.Size()}}
  return nil
 });return total
}
func cineLegacyCacheNames(root,cache string)[]string{
 if root==cache{return nil}
 names:=[]string{"plataformas_cache.json","tmdb_fichas_lab.json"}
 for _,dir:=range []string{"carteles","imagenes_descubrir"}{
  entries,_:=os.ReadDir(filepath.Join(root,dir));for _,e:=range entries{if e.Type().IsRegular()&&validPosterPath("/"+e.Name()){names=append(names,filepath.Join(dir,e.Name()))}}
 }
 return names
}
func cineLegacyCacheUsage(root,cache string)labFileUsage{
 var total labFileUsage
 for _,name:=range cineLegacyCacheNames(root,cache){if info,e:=os.Lstat(filepath.Join(root,name));e==nil&&info.Mode().IsRegular(){total.Files++;total.Bytes+=info.Size()}}
 return total
}
func cineIdenticalFiles(old,new string)bool{
 for _,path:=range []string{old,new}{info,e:=os.Lstat(path);if e!=nil||!info.Mode().IsRegular()||info.Size()>64<<20{return false};real,e:=filepath.EvalSymlinks(path);if e!=nil||real!=path{return false}}
 a,e:=os.Open(old);if e!=nil{return false};defer a.Close()
 b,e:=os.Open(new);if e!=nil{return false};defer b.Close()
 ai,e:=a.Stat();if e!=nil{return false};bi,e:=b.Stat();if e!=nil||os.SameFile(ai,bi)||ai.Size()!=bi.Size(){return false}
 x,y:=make([]byte,64<<10),make([]byte,64<<10)
 for{n,e1:=io.ReadFull(a,x);m,e2:=io.ReadFull(b,y);if n!=m||!bytes.Equal(x[:n],y[:m]){return false};if e1==io.EOF||e1==io.ErrUnexpectedEOF{return e1==e2};if e1!=nil||e2!=nil{return false}}
}
func cineCleanVerifiedLegacyCache(root,cache string)(int,error){
 n:=0
 source,e:=filepath.EvalSymlinks(root);if e!=nil{return 0,e}
 target,e:=filepath.EvalSymlinks(cache);if os.IsNotExist(e){return 0,nil};if e!=nil{return 0,e}
 if source==target{return 0,nil}
 root,cache=source,target
 for _,name:=range cineLegacyCacheNames(root,cache){
  old,new:=filepath.Join(root,name),filepath.Join(cache,name)
  if cineIdenticalFiles(old,new){if e:=os.Remove(old);e!=nil{return n,e};n++}
 }
 return n,nil
}
