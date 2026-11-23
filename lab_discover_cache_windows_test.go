//go:build windows && lab

package main

import (
 "encoding/json"
 "os"
 "path/filepath"
 "testing"
 "time"
 "strings"
)

func TestLabCacheTTL(t *testing.T) {
 now:=time.Now().Unix()
 fresh:=labDetailRecord{Data:map[string]any{"title":"El padrino"},MetaAt:now,ExtraAt:now,ProvidersAt:now}
 if labNeedsRefresh(fresh){t.Fatal("una ficha recién guardada no debe caducar")}
 oldProviders:=fresh;oldProviders.ProvidersAt=now-25*3600
 if !labNeedsRefresh(oldProviders){t.Fatal("plataformas de más de 24 horas deberían actualizarse")}
 oldTrailer:=fresh;oldTrailer.ExtraAt=now-8*86400
 if !labNeedsRefresh(oldTrailer){t.Fatal("tráiler de más de siete días debería actualizarse")}
 oldDetails:=fresh;oldDetails.MetaAt=now-31*86400
 if !labNeedsRefresh(oldDetails){t.Fatal("datos de más de 30 días deberían actualizarse")}
 if labJobKey(labFilmJob{Title:"El padrino",Year:1972})!=labJobKey(labFilmJob{Title:"El padrino",Year:1972,ID:238}){
  t.Fatal("el trabajo debe conservar su identidad tras resolver un ID")
 }
}
func TestLabCachePersistsSeparately(t *testing.T) {
 file:=filepath.Join(t.TempDir(),"tmdb_fichas_lab.json")
 c:=&labCatalogCache{file:file,disk:labCacheDisk{Version:1,Items:map[string]labDetailRecord{},Resolved:map[string]int{}}}
 saved:=labDetailRecord{Data:map[string]any{"id":238,"title":"El padrino"},MetaAt:time.Now().Unix()}
 c.remember(238,saved,"el padrino|1972")
 raw,err:=os.ReadFile(file)
 if err!=nil {t.Fatal("no se ha creado la caché persistente:",err)}
 var disk labCacheDisk
 if err=json.Unmarshal(raw,&disk);err!=nil{t.Fatal(err)}
 if disk.Items["238"].Data["title"]!="El padrino"||disk.Resolved["el padrino|1972"]!=238{
  t.Fatal("datos o identificaciones no recuperables")
 }
 if _,err:=os.Stat(filepath.Join(filepath.Dir(file),"colecciones.json"));err==nil{t.Fatal("la caché jamás debe crear o modificar colecciones")}
}
func TestLabPrefetchDeduplicatesAndPauses(t *testing.T){
 c:=&labCatalogCache{
  paused:true,
  disk:labCacheDisk{Version:1,Items:map[string]labDetailRecord{},Resolved:map[string]int{}},
  targets:map[string]labFilmJob{},queued:map[string]bool{},
 }
 movie:=labFilmJob{Title:"El padrino",Year:1972}
 c.enqueue([]labFilmJob{movie,movie,movie})
 if len(c.queue)!=1||len(c.targets)!=1||c.running {t.Fatal("la precarga en pausa debe crear un solo trabajo")}
 c.enqueue([]labFilmJob{{Title:"El padrino",Year:1972,ID:238}})
 if len(c.queue)!=1{t.Fatal("resolver un ID no debe duplicar el trabajo")}
}

func TestLabWindowsHTMLInjection(t *testing.T) {
 page,err:=makeHTML("☐ El padrino (1972)","test-token")
 if err!=nil{t.Fatal(err)}
 for _,required:=range []string{
  "labCreateDetailButton(f)",
  "b.append(togg,pill,labCreateDetailButton(f),fix);frag.appendChild(b);",
  "labPrefetchBanner",
  "labQueuePrefetchSoon()",
  "labDetailOverlay",
 } {
  if !strings.Contains(page,required){t.Errorf("falta elemento LAB en Windows: %s",required)}
 }
 if strings.Contains(page,"b.append(togg,pill,fix,link);frag.appendChild(b);"){
  t.Fatal("el acceso ↗ externo no debe mostrarse en la fila de Windows LAB")
 }
}
