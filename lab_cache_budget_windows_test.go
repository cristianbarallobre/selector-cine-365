//go:build windows && lab

package main

import (
 "encoding/json"
 "os"
 "path/filepath"
 "strconv"
 "strings"
 "testing"
 "time"
)

func TestLabCacheDirectorySeparateFromRoaming(t *testing.T){
 local:=t.TempDir()
 t.Setenv("LOCALAPPDATA",local)
 old:=filepath.Join(t.TempDir(),"Roaming",appDataDirectoryName)
 got:=appCacheRoot(old)
 if got!=filepath.Join(local,appDataDirectoryName,"cache_lab"){t.Fatalf("La caché debe estar bajo LocalAppData: %s",got)}
 if got==old{t.Fatal("Caché comparte Roaming")}
}
func TestLabCatalogEvictsOldestInsteadOfStoppingSaving(t *testing.T){
 dir:=t.TempDir()
 file:=filepath.Join(dir,"tmdb_fichas_lab.json")
 c:=&labCatalogCache{file:file,disk:labCacheDisk{Version:1,Items:map[string]labDetailRecord{},Resolved:map[string]int{}}}
 for i:=1;i<=labCacheMaxEntries+5;i++{
  id:=strconv.Itoa(i)
  c.disk.Items[id]=labDetailRecord{Data:map[string]any{"id":i,"title":"Película"},LastUsedAt:int64(i)}
 }
 if !c.pruneLocked(){t.Fatal("La caché no eliminó fichas antiguas")}
 if len(c.disk.Items)>labCacheMaxEntries{t.Fatalf("Superó el límite: %d",len(c.disk.Items))}
 if len(c.disk.Items["1"].Data)>0{t.Fatal("No expulsó la ficha menos usada")}
 if len(c.disk.Items[strconv.Itoa(labCacheMaxEntries+5)].Data)==0{t.Fatal("Perdió la ficha más reciente")}
 if err:=c.fileSaveLocked();err!=nil{t.Fatal(err)}
 raw,err:=os.ReadFile(file);if err!=nil{t.Fatal(err)}
 if len(raw)>labCacheMaxBytes{t.Fatal("El archivo siguió siendo demasiado grande")}
 var disk labCacheDisk
 if json.Unmarshal(raw,&disk)!=nil||len(disk.Items)!=labCacheMaxEntries{t.Fatal("El archivo compacto no se puede recargar")}
 // Al recordar una ficha nueva tampoco debe interrumpirse la persistencia.
 c.remember(999999,labDetailRecord{Data:map[string]any{"id":999999,"title":"Nueva"}}, "nueva|2026")
 raw,err=os.ReadFile(file);if err!=nil{t.Fatal(err)}
 if !strings.Contains(string(raw),"Nueva"){t.Fatal("Las nuevas fichas no se guardan tras limpiar")}
 if len(c.disk.Items)>labCacheMaxEntries{t.Fatal("Caché sin límite tras remember")}
}
func TestLabCacheAccessRefreshesLRU(t *testing.T){
 now:=time.Now().Unix()
 c:=&labCatalogCache{disk:labCacheDisk{Items:map[string]labDetailRecord{"9":{Data:map[string]any{"title":"Favorita"},MetaAt:now-15*86400}}}}
 rec,ok:=c.record(9)
 if !ok||rec.LastUsedAt<now-5{t.Fatal("Consultar una película debe marcarla como recientemente usada")}
}
func TestLabCacheResolvedBudgetWithoutChangingManualMappings(t *testing.T){
 c:=&labCatalogCache{disk:labCacheDisk{Version:1,Items:map[string]labDetailRecord{},Resolved:map[string]int{}}}
 for i:=0;i<labCacheMaxResolved+50;i++{c.disk.Resolved["titulo-"+strconv.Itoa(i)+"|2025"]=i+1}
 if !c.pruneLocked(){t.Fatal("Las coincidencias descartables superaron el máximo")}
 if len(c.disk.Resolved)>labCacheMaxResolved{t.Fatal("No se redujo el índice desechable")}
}
