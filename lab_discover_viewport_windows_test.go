//go:build windows

package main

import (
 "context"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "net/url"
 "os"
 "os/exec"
 "path/filepath"
 "strconv"
 "strings"
 "sync"
 "testing"
 "time"
)

// Reproduce exactamente el fallo de LAB-02.12 en Edge: seis columnas,
// ocho títulos en TMDB página 1 y espacio para otras dos filas.
func TestLabDiscoverFillsVisibleRowsInEdge(t *testing.T){
 if os.Getenv("CI")==""{t.Skip("La prueba de Edge se ejecuta en CI Windows")}
 browser:=""
 for _,p:=range []string{
  filepath.Join(os.Getenv("ProgramFiles(x86)"),"Microsoft","Edge","Application","msedge.exe"),
  filepath.Join(os.Getenv("ProgramFiles"),"Google","Chrome","Application","chrome.exe"),
  filepath.Join(os.Getenv("ProgramFiles(x86)"),"Google","Chrome","Application","chrome.exe"),
 }{
  if info,err:=os.Stat(p);err==nil&&!info.IsDir(){browser=p;break}
 }
 if browser==""{t.Skip("Edge o Chrome no disponibles")}
 initial:="☐ Prueba visual (2026)"
 page,err:=makeHTML(initial,"smoke-token")
 if err!=nil{t.Fatal(err)}
 script:=`<script>
window.addEventListener('DOMContentLoaded',()=>{
 function waitForLibrary(){
   if(document.getElementById("libraryApp")?.hidden){setTimeout(waitForLibrary,50);return;}
   document.getElementById('labDiscoverNav')?.click();
   setTimeout(()=>{
     const grid=document.getElementById('labMoviesGrid');
     const cards=[...(grid?.querySelectorAll('.lab-movie-card')||[])];
     const rows=new Set(cards.map(card=>Math.round(card.getBoundingClientRect().top)));
     const genre=document.getElementById('labDiscoverGenre');
     const initial={cards:cards.length,rows:rows.size,width:grid?.clientWidth||0,
       pageLabel:document.getElementById('labPageLabel')?.textContent||''};
     if(genre){genre.value='878';genre.dispatchEvent(new Event('change',{bubbles:true}));}
     setTimeout(()=>{
       const latest=[...(grid?.querySelectorAll('.lab-movie-card')||[])];
       const filtered={
         cards:latest.length,rows:new Set(latest.map(c=>Math.round(c.getBoundingClientRect().top))).size,
         genre:genre?.value,
         pageLabel:document.getElementById('labPageLabel')?.textContent||'',
         discoverVisible:document.getElementById('labDiscoverView')?.hidden===false
       };
       grid?.querySelector('.lab-card-details')?.click();
       setTimeout(()=>{
         const castNode=document.getElementById('labDetailCast');
         const cast={
           actor:castNode?.querySelector('.lab-cast-actor')?.textContent||'',
           role:castNode?.querySelector('.lab-cast-role')?.textContent||'',
           roleLabel:castNode?.querySelector('.lab-cast-role')?getComputedStyle(castNode.querySelector('.lab-cast-role'),'::before').content:'',
           director:castNode?.querySelector('.lab-cast-director strong')?.textContent||'',
           visible:document.getElementById('labDetailOverlay')?.hidden===false
         };
         document.body.setAttribute('data-lab-grid',encodeURIComponent(JSON.stringify({initial,filtered,cast})));
       },450);
     },1000);
   },1000);
 }
 waitForLibrary();
});
</script>`
 page=strings.Replace(page,"</head>",script+"\n</head>",1)
 var mu sync.Mutex
 requests:=map[string]int{}
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path=="/api/library"{
   serveSmokeLibrary(w,r,initial);return
  }
  if r.URL.Path=="/"{
   w.Header().Set("Content-Type","text/html; charset=utf-8")
   _,_=w.Write([]byte(page));return
  }
  if r.URL.Path=="/api/discover/detail"{
   w.Header().Set("Content-Type","application/json")
   _=json.NewEncoder(w).Encode(map[string]any{
     "id":101,"title":"Película de prueba","release_date":"2026-10-01","poster_path":"",
     "runtime":110,"vote_count":77,"vote_average":7.7,"director":"Kitao Sakurai",
     "cast":[]map[string]string{
       {"name":"Andrew Koji","character":"Ryu"},
       {"name":"Noah Centineo","character":"Ken Masters"},
     },
     "genres":[]string{"Acción"},"tmdb_url":"https://www.themoviedb.org/movie/101",
   })
   return
  }
  if r.URL.Path=="/api/discover/list"{
   pg,_:=strconv.Atoi(r.URL.Query().Get("page"))
   if pg<1 {pg=1}
   key:=r.URL.Query().Get("genre")
   mu.Lock();requests[key]++;mu.Unlock()
   movies:=make([]map[string]any,8)
   for i:=0;i<8;i++{
    movies[i]=map[string]any{"id":pg*100+i+1,"title":"Película "+strconv.Itoa(pg)+"-"+strconv.Itoa(i+1),
     "release_date":"2026-10-01","vote_count":80,"vote_average":7.4,"poster_path":""}
   }
   w.Header().Set("Content-Type","application/json")
   _=json.NewEncoder(w).Encode(map[string]any{"page":pg,"total_pages":500,"total_results":4000,"results":movies})
   return
  }
  w.Header().Set("Content-Type","application/json")
  w.WriteHeader(http.StatusNotFound)
  _,_=w.Write([]byte(`{"error":"sin datos en la prueba"}`))
 }))
 defer server.Close()
 ctx,cancel:=context.WithTimeout(context.Background(),40*time.Second);defer cancel()
 args:=[]string{
  "--headless=new","--disable-gpu","--no-first-run","--no-default-browser-check",
  "--disable-background-networking","--disable-extensions","--no-sandbox",
  "--window-size=1480,960",
  "--user-data-dir="+filepath.Join(t.TempDir(),"profile"),
  "--virtual-time-budget=7600","--dump-dom",server.URL+"/",
 }
 output,err:=exec.CommandContext(ctx,browser,args...).CombinedOutput()
 if err!=nil{t.Fatalf("Edge no completó prueba de rejilla: %v (%s)",err,string(output[max(0,len(output)-900):]))}
 text:=string(output)
 prefix:=`data-lab-grid="`
 at:=strings.Index(text,prefix)
 if at<0{t.Fatalf("El navegador no generó datos de rejilla: %s",text[max(0,len(text)-950):])}
 text=text[at+len(prefix):]
 end:=strings.IndexByte(text,'"')
 if end<0{t.Fatal("Marcador de rejilla incompleto")}
 decoded,err:=url.QueryUnescape(text[:end]);if err!=nil{t.Fatal(err)}
 var result struct {
  Initial struct{ Cards int; Rows int; Width int; PageLabel string } `json:"initial"`
  Filtered struct{ Cards int; Rows int; Genre string; PageLabel string; DiscoverVisible bool } `json:"filtered"`
  Cast struct{ Actor string; Role string; RoleLabel string; Director string; Visible bool } `json:"cast"`
 }
 if err:=json.Unmarshal([]byte(decoded),&result);err!=nil{t.Fatalf("JSON de navegador inválido: %v (%s)",err,decoded)}
 if result.Initial.Width<850{t.Errorf("No se ha reproducido el tamaño de la captura (ancho %d)",result.Initial.Width)}
 if result.Initial.Cards<18||result.Initial.Rows<3{
  t.Errorf("La página 1 sigue dejando hueco con ocho títulos por llamada: %+v",result.Initial)
 }
 if result.Filtered.Cards<18||result.Filtered.Rows<3||result.Filtered.Genre!="878"||!result.Filtered.DiscoverVisible {
  t.Errorf("La rejilla con el filtro de género no se completa: %+v",result.Filtered)
 }
 if result.Cast.Actor!="Andrew Koji"||result.Cast.Role!="Ryu"||
   result.Cast.Director!="Kitao Sakurai"||!result.Cast.Visible{
  t.Errorf("Los créditos de la ficha no se distinguen en Edge: %+v",result.Cast)
 }
 if !strings.Contains(result.Cast.RoleLabel,"Personaje:"){
  t.Errorf("El personaje no está etiquetado en la ficha: %+v",result.Cast)
 }
 mu.Lock()
 general:=requests[""]
 filtered:=requests["878"]
 mu.Unlock()
 if general<3{t.Errorf("Para llenar seis columnas hacen falta >= 3 páginas TMDB; realizadas %d",general)}
 if filtered<3{t.Errorf("Los filtros no han provocado suficientes peticiones TMDB; realizadas %d",filtered)}
 t.Logf("Edge confirma rejilla completa: primera=%+v, filtrada=%+v, peticiones general=%d filtradas=%d",
  result.Initial,result.Filtered,general,filtered)
}
