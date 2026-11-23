//go:build windows && lab

package main

import (
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "path/filepath"
 "strings"
 "testing"
)

func TestLabSmartMetaReadOnlyEndpoint(t *testing.T) {
 a:=&runningApp{savePath:filepath.Join(t.TempDir(),"Selector_Cine_365_LAB","lista.txt"),token:"secret"}
 c:=a.labCatalog()
 c.mu.Lock()
 c.disk.Items["123"]=labDetailRecord{Data:map[string]any{
  "runtime":float64(87),
  "genres":[]any{"Comedia","Drama"},
  "vote_average":float64(8.1),
  "vote_count":float64(550),
 }}
 c.disk.Resolved["ejemplo|2024"]=123
 c.mu.Unlock()
 for _,tt:=range []struct{method string;token string;want int}{
  {"POST","secret",http.StatusMethodNotAllowed},
  {"GET","bad",http.StatusForbidden},
  {"GET","secret",http.StatusOK},
 } {
  req:=httptest.NewRequest(tt.method,"/api/discover/smart-cache",nil)
  req.Header.Set("X-Atlas-Token",tt.token)
  recorder:=httptest.NewRecorder()
  a.serveLabSmartCache(recorder,req)
  if recorder.Code!=tt.want{t.Fatalf("%s %s: %d (esperado %d)",tt.method,tt.token,recorder.Code,tt.want)}
  if tt.want!=http.StatusOK{continue}
  var body labSmartCacheReply
  if err:=json.Unmarshal(recorder.Body.Bytes(),&body);err!=nil{t.Fatal(err)}
  film,exists:=body.Items["123"]
  if !exists||film.Runtime!=87||film.VoteCount!=550||film.VoteAverage!=8.1||len(film.Genres)!=2 {
   t.Fatalf("metadatos incompletos: %#v",film)
  }
  if body.Resolved["ejemplo|2024"]!=123{t.Fatal("la identificación debe poder usarse sin consultar TMDB")}
  if recorder.Header().Get("Cache-Control")!="no-store"{t.Fatal("la respuesta debe evitar cachés intermediarias")}
 }
}
func TestLabSmartInjectedOnlyWindowsLab(t *testing.T){
 page,err:=makeHTML("☐ Ejemplo (2024)","token")
 if err!=nil{t.Fatal(err)}
 for _,snippet:=range []string{
  "labSmartEnabled","labSmartSeen","labSmartDuration","labSmartGenre",
  "labSmartReason","labSmartChoose","labSmartFetchCache","labSmartResultDetails","labSmartSyncDetailsButton",
  "LAB-02.14 WINDOWS","labDiscoverView",
  "function updateSortAvailability(){queueMicrotask(labSmartUpdateCount);",
  "function resetChoice(){chosen=null;labSmartSyncDetailsButton();labSmartSyncSeenAction();$('labSmartReason').hidden=true;",
 } {
  if !strings.Contains(page,snippet){t.Errorf("elemento LAB-02.14 ausente: %s",snippet)}
 }
 if strings.Count(page,`id="labSmartEnabled"`)!=1 {t.Fatal("el interruptor debe existir solo una vez")}
 if !strings.Contains(page,"Elegir una película al azar"){t.Fatal("el sorteo clásico debe conservarse")}
 if strings.Contains(page,"Todas las películas pendientes"){t.Fatal("el selector superior debe decir Todas las plataformas")}
 if strings.Count(page,"Todas las plataformas")<2 {t.Fatal("el selector y su reconstrucción deben usar Todas las plataformas")}
 if strings.Count(page,`id="labSmartResultDetails"`)!=1 {t.Fatal("debe haber un solo botón Detalles en la tarjeta sorteada")}
 if !strings.Contains(page,"labOpenFilmDetails(chosen)") {t.Fatal("la tarjeta sorteada debe abrir la ficha de detalle existente")}
 if strings.Contains(page,"labSmartMood")||strings.Contains(page,"labSmartInfo"){t.Fatal("el selector ¿Qué te apetece? y las aclaraciones deben desaparecer")}
 for _,label:=range []string{"No vistas (pendientes)","Vistas","Vistas y no vistas"}{
  if !strings.Contains(page,label){t.Errorf("opción de estado ausente: %s",label)}
 }
 if !strings.Contains(page,"function labSmartPool()"){t.Fatal("el modo inteligente debe incluir estados vistos y pendientes")}
 if !strings.Contains(page,"labSmartSyncSeenAction"){t.Fatal("la acción vista/pendiente debe ajustarse a la película")}
 if strings.Contains(page,"swShowMaximized"){t.Fatal("el tamaño de ventana no debe inyectarse mediante JavaScript")}
}
