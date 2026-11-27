package main

// Las cachés son archivos descartables. No tocar los JSON de listas,
// identificaciones manuales, configuración ni tokens de sincronización.
import (
 "errors"
 "io"
 "os"
 "path/filepath"
 "sort"
 "time"
)

const (
 appPosterMaxFiles=2500
 appPosterMaxBytes=int64(150<<20)
 appPosterMaxAge=120*24*time.Hour
)

// Copia una caché antigua si aún no existe en LocalAppData. Mantiene intacto
// el origen Roaming, que puede seguir usando otra versión de la aplicación.
func copyLegacyCacheIfMissing(oldPath,newPath string,maxBytes int64) error {
 if oldPath==newPath{return nil}
 if _,err:=os.Stat(newPath);err==nil{return nil}else if !errors.Is(err,os.ErrNotExist){return err}
 src,err:=os.Open(oldPath);if errors.Is(err,os.ErrNotExist){return nil};if err!=nil{return err}
 defer src.Close()
 stat,err:=src.Stat();if err!=nil{return err}
 if !stat.Mode().IsRegular()||stat.Size()>maxBytes||stat.Size()<=0{return nil}
 if err:=os.MkdirAll(filepath.Dir(newPath),0700);err!=nil{return err}
 tmp,err:=os.CreateTemp(filepath.Dir(newPath),".cache-copia-*");if err!=nil{return err}
 defer os.Remove(tmp.Name())
 if err:=tmp.Chmod(0600);err!=nil{tmp.Close();return err}
 n,err:=io.Copy(tmp,io.LimitReader(src,maxBytes+1))
 if err!=nil||n!=stat.Size()||n>maxBytes{tmp.Close();return errors.New("copia de caché incompleta")}
 if err:=tmp.Close();err!=nil{return err}
 if _,err:=os.Stat(newPath);err==nil{return nil}else if !errors.Is(err,os.ErrNotExist){return err}
 return os.Rename(tmp.Name(),newPath)
}

// Solo se inspeccionan y eliminan imágenes con nombres TMDB válidos.
// No seguir enlaces simbólicos ni recorrer directorios del usuario.
func prunePosterCache(dir string,maxFiles int,maxBytes int64,expiry time.Duration) (int,error) {
 type fileRecord struct{name string;size int64;modified time.Time}
 entries,err:=os.ReadDir(dir)
 if errors.Is(err,os.ErrNotExist){return 0,nil};if err!=nil{return 0,err}
 items:=make([]fileRecord,0,len(entries))
 now:=time.Now()
 removed:=0
 var total int64
 for _,entry:=range entries{
  if !entry.Type().IsRegular()||!validPosterPath("/"+entry.Name()){continue}
  info,err:=entry.Info();if err!=nil{continue}
  file:=filepath.Join(dir,entry.Name())
  if expiry>0&&now.Sub(info.ModTime())>expiry{
   if os.Remove(file)==nil{removed++}
   continue
  }
  items=append(items,fileRecord{entry.Name(),info.Size(),info.ModTime()})
  total+=info.Size()
 }
 sort.Slice(items,func(i,j int)bool{return items[i].modified.Before(items[j].modified)})
 remaining:=len(items)
 for _,item:=range items{
  if remaining<=maxFiles&&total<=maxBytes{break}
  if os.Remove(filepath.Join(dir,item.name))==nil {removed++;remaining--;total-=item.size}
 }
 return removed,nil
}
