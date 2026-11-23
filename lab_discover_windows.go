//go:build windows

package main

import (
 "context"
 "embed"
 "encoding/json"
 "io"
 "net/http"
 "net/url"
 "strconv"
 "strings"
 "time"
 "unicode/utf8"
)

//go:embed lab_discover_windows.html
var labDiscoveryFiles embed.FS

const labDiscoverySeparator = "<!-- LAB_DISCOVERY_SCRIPT -->"

func injectLabDiscovery(page string) string {
 if !appIsLab { return page }
 raw, err := labDiscoveryFiles.ReadFile("lab_discover_windows.html")
 if err != nil { return page }
 fragments := strings.SplitN(string(raw), labDiscoverySeparator, 2)
 if len(fragments)!=2 { return page }
 nav := `<button class="nav-btn" type="button" id="labDiscoverNav">✦ Descubrir <span class="lab-nav-tag">LAB</span></button>`
 page = strings.Replace(page, `<button class="nav-btn" type="button" data-view="statistics">`, nav+`<button class="nav-btn" type="button" data-view="statistics">`, 1)
 page = strings.Replace(page, "🧪 SELECTOR CINE 365 LAB · v9 Beta 14 Preview", "🧪 SELECTOR CINE 365 · LAB-02.3 WINDOWS", 1)
 // Los botones específicos se inyectan solo en Windows LAB; plantilla.html permanece intacta.
 // En Windows LAB, el enlace externo de cada fila sobra: ya está dentro de Detalles.
 page = strings.Replace(page, "b.append(togg,pill,fix,link);frag.appendChild(b);", "b.append(togg,pill,labCreateDetailButton(f),fix);frag.appendChild(b);", 1)
 page = strings.Replace(page, `<input type="search" id="search"`, `<div class="lab-prefetch-banner" id="labPrefetchBanner" role="status" aria-live="polite"><span id="labPrefetchText">Fichas TMDB: comprobando caché…</span><button type="button" id="labPrefetchPause" hidden>Pausar</button></div><input type="search" id="search"`, 1)
 // La precarga se inicia tras cargar biblioteca e identificaciones TMDB, y se actualiza al editar listas.
 page = strings.Replace(page, "void cloudAutoReceive(true);});", "void cloudAutoReceive(true);labQueuePrefetchSoon();});", 1)
 page = strings.Replace(page, "function scheduleSave(){membershipCache=null;", "function scheduleSave(){labQueuePrefetchSoon();membershipCache=null;", 1)
 page = strings.Replace(page, "function showView(name){if(", "function showView(name){labLeave();if(name==='library')labQueuePrefetchSoon();if(", 1)
 page = strings.Replace(page, "platformConfigured=true;identityConfigured=true;closePlatformModal();", "platformConfigured=true;identityConfigured=true;closePlatformModal();labQueuePrefetchSoon();", 1)
 page = strings.Replace(page, "</main>", fragments[0]+"\n</main>", 1)
 page = strings.Replace(page, "\n})();\n</script>", "\n"+fragments[1]+"\n})();\n</script>", 1)
 return page
}

type labMovie struct {
 ID int `json:"id"`
 Title string `json:"title"`
 OriginalTitle string `json:"original_title,omitempty"`
 ReleaseDate string `json:"release_date,omitempty"`
 Overview string `json:"overview,omitempty"`
 PosterPath string `json:"poster_path,omitempty"`
 VoteAverage float64 `json:"vote_average"`
 VoteCount int `json:"vote_count"`
}
func safeLabMovie(m labMovie) labMovie {
 if m.ID<=0||m.ID>999999999 {m.ID=0}
 if !validPosterPath(m.PosterPath) {m.PosterPath=""}
 if utf8.RuneCountInString(m.Title)>200 {m.Title=string([]rune(m.Title)[:200])}
 if utf8.RuneCountInString(m.Overview)>900 {m.Overview=string([]rune(m.Overview)[:900])}
 return m
}
type labMoviePage struct {
 Page int `json:"page"`
 TotalPages int `json:"total_pages"`
 TotalResults int `json:"total_results"`
 Results []labMovie `json:"results"`
}
type labGenre struct { Name string `json:"name"` }
type labPerson struct { Name string `json:"name"`; Character string `json:"character"`; Job string `json:"job"` }
type labVideo struct { Key string `json:"key"`; Site string `json:"site"`; Type string `json:"type"`; Official bool `json:"official"` }
type labProvider struct { ProviderName string `json:"provider_name"` }
type labProviders struct {
 Link string `json:"link,omitempty"`
 Flatrate []labProvider `json:"flatrate,omitempty"`
 Free []labProvider `json:"free,omitempty"`
 Ads []labProvider `json:"ads,omitempty"`
 Rent []labProvider `json:"rent,omitempty"`
 Buy []labProvider `json:"buy,omitempty"`
}
func labProviderNames(items []labProvider) []string {
 seen:=map[string]bool{}
 names:=make([]string,0,len(items))
 for _,p:=range items {s:=strings.TrimSpace(p.ProviderName);if s!=""&&!seen[s]{seen[s]=true;names=append(names,s)}}
 return names
}
func labTrailer(vs []labVideo) string {
 for pass:=0;pass<2;pass++ {
  for _,v:=range vs {
   if v.Site!="YouTube"||v.Type!="Trailer"||pass==0&&!v.Official||len(v.Key)!=11 {continue}
   valid:=true
   for _,c:=range v.Key {if !(c>='a'&&c<='z'||c>='A'&&c<='Z'||c>='0'&&c<='9'||c=='_'||c=='-'){valid=false;break}}
   if valid {return "https://www.youtube.com/watch?v="+v.Key}
  }
 }
 return ""
}

func (a *runningApp) serveLabDiscovery(w http.ResponseWriter,r *http.Request) {
 if !appIsLab {http.NotFound(w,r);return}
 if !a.auth(r) {http.Error(w,"No autorizado",403);return}
 if r.URL.Path=="/api/discover/prefetch"||r.URL.Path=="/api/discover/prefetch/pause" {a.serveLabPrefetch(w,r);return}
 if r.Method!=http.MethodGet {http.Error(w,"Método no permitido",405);return}
 if r.URL.Path=="/api/discover/detail" {a.labDetailFromCache(w,r);return}
 if r.URL.Path=="/api/discover/resolve" {a.labCachedIdentity(w,r);return}
 if r.URL.Path=="/api/discover/poster" {a.labDiscoveryPoster(w,r);return}
 w.Header().Set("Content-Type","application/json; charset=utf-8")
 fail:=func(code int,msg string){w.WriteHeader(code);_ =json.NewEncoder(w).Encode(map[string]string{"error":msg})}
 key:=a.platforms.getKey()
 if key=="" {fail(412,"Conecta TMDB en Biblioteca para utilizar Descubrir.");return}
 ctx,cancel:=context.WithTimeout(r.Context(),22*time.Second);defer cancel()
 switch r.URL.Path {
 case "/api/discover/list", "/api/discover/search":
  page:=1
  if raw:=r.URL.Query().Get("page");raw!="" {n,err:=strconv.Atoi(raw);if err!=nil||n<1||n>500 {fail(400,"Página no válida");return};page=n}
  endpoint:="/search/movie"
  params:=url.Values{"language":{"es-ES"},"include_adult":{"false"},"page":{strconv.Itoa(page)}}
  if r.URL.Path=="/api/discover/search" {
   q:=strings.TrimSpace(r.URL.Query().Get("q"))
   if n:=utf8.RuneCountInString(q);n<2||n>100 {fail(400,"Introduce entre 2 y 100 caracteres.");return}
   params.Set("query",q)
  } else {
   kind:=r.URL.Query().Get("kind")
   switch kind {
   case "trending":endpoint="/trending/movie/week"
   case "popular":endpoint="/movie/popular"
   case "top":endpoint="/movie/top_rated"
   case "upcoming":endpoint="/movie/upcoming"
   default:fail(400,"Categoría no válida");return
   }
   if kind!="trending" {params.Set("region","ES")}
  }
  var data labMoviePage
  if err:=requestTMDB(ctx,endpoint,params,key,&data);err!=nil {fail(502,err.Error());return}
  if data.TotalPages>500 {data.TotalPages=500}
  filtered:=make([]labMovie,0,len(data.Results))
  for _,m:=range data.Results {m=safeLabMovie(m);if m.ID>0&&m.Title!="" {filtered=append(filtered,m)}}
  data.Results=filtered
  _=json.NewEncoder(w).Encode(data)
 default:http.NotFound(w,r)
 }
}
func (a *runningApp) labDiscoveryPoster(w http.ResponseWriter,r *http.Request) {
 path:=r.URL.Query().Get("path")
 if !validPosterPath(path){http.Error(w,"Imagen no válida",400);return}
 ctx,cancel:=context.WithTimeout(r.Context(),17*time.Second);defer cancel()
 req,err:=http.NewRequestWithContext(ctx,"GET","https://image.tmdb.org/t/p/w342"+path,nil)
 if err!=nil {http.Error(w,"Imagen no disponible",502);return}
 resp,err:=tmdbClient.Do(req);if err!=nil {http.Error(w,"Imagen no disponible",502);return}
 defer resp.Body.Close()
 if resp.StatusCode!=200 {http.Error(w,"Imagen no disponible",502);return}
 mime:="image/jpeg";if strings.HasSuffix(path,".png"){mime="image/png"}else if strings.HasSuffix(path,".webp"){mime="image/webp"}
 w.Header().Set("Content-Type",mime);w.Header().Set("Cache-Control","private, max-age=86400")
 _,_=io.Copy(w,io.LimitReader(resp.Body,2<<20))
}
