//go:build windows

package main

import (
 "context"
 "embed"
 "encoding/json"
 "fmt"
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
 page = strings.Replace(page, "🧪 SELECTOR CINE 365 LAB ·", "🧪 SELECTOR CINE 365 · LAB-02 ·", 1)
 // El botón Detalles se añade solo al renderizado Windows LAB, sin cambiar plantilla.html.
 page = strings.Replace(page, "titleLine.appendChild(title);", "titleLine.append(title,labCreateDetailButton(f));", 1)
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
 if r.Method!=http.MethodGet {http.Error(w,"Método no permitido",405);return}
 if r.URL.Path=="/api/discover/poster" {a.labDiscoveryPoster(w,r);return}
 w.Header().Set("Content-Type","application/json; charset=utf-8")
 fail:=func(code int,msg string){w.WriteHeader(code);_ =json.NewEncoder(w).Encode(map[string]string{"error":msg})}
 key:=a.platforms.getKey()
 if key=="" {fail(412,"Conecta TMDB en Biblioteca para utilizar Descubrir.");return}
 ctx,cancel:=context.WithTimeout(r.Context(),22*time.Second);defer cancel()
 switch r.URL.Path {
 case "/api/discover/list", "/api/discover/search":
  page:=1
  if raw:=r.URL.Query().Get("page");raw!="" {n,err:=strconv.Atoi(raw);if err!=nil||n<1||n>5 {fail(400,"Página no válida");return};page=n}
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
  if data.TotalPages>5 {data.TotalPages=5}
  filtered:=make([]labMovie,0,len(data.Results))
  for _,m:=range data.Results {m=safeLabMovie(m);if m.ID>0&&m.Title!="" {filtered=append(filtered,m)}}
  data.Results=filtered
  _=json.NewEncoder(w).Encode(data)
 case "/api/discover/detail":
  id,err:=strconv.Atoi(r.URL.Query().Get("id"))
  if err!=nil||id<1||id>999999999 {fail(400,"Identificador de película no válido");return}
  path:="/movie/"+strconv.Itoa(id)
  var detail struct {
   labMovie
   Runtime int `json:"runtime"`
   Genres []labGenre `json:"genres"`
   Credits struct { Cast []labPerson `json:"cast"`; Crew []labPerson `json:"crew"` } `json:"credits"`
   Videos struct { Results []labVideo `json:"results"` } `json:"videos"`
  }
  if err:=requestTMDB(ctx,path,url.Values{"language":{"es-ES"},"append_to_response":{"credits,videos"}},key,&detail);err!=nil {fail(502,err.Error());return}
  if detail.ID!=id||detail.Title=="" {fail(502,"TMDB devolvió una ficha inválida");return}
  detail.labMovie=safeLabMovie(detail.labMovie)
  director:=""
  for _,p:=range detail.Credits.Crew {if p.Job=="Director"&&p.Name!="" {if director!="" {director+=" / "};director+=p.Name;if len(director)>220{break}}}
  cast:=make([]map[string]string,0,10)
  for _,p:=range detail.Credits.Cast {if len(cast)>=10 {break};if p.Name!="" {cast=append(cast,map[string]string{"name":p.Name,"character":p.Character})}}
  genres:=make([]string,0,len(detail.Genres));for _,g:=range detail.Genres {if g.Name!="" {genres=append(genres,g.Name)}}
  trailer:=labTrailer(detail.Videos.Results)
  if trailer=="" {var videos struct{Results []labVideo `json:"results"`};if requestTMDB(ctx,path+"/videos",url.Values{"language":{"en-US"}},key,&videos)==nil {trailer=labTrailer(videos.Results)}}
  var providers struct{ Results map[string]labProviders `json:"results"` }
  note:=""
  if err:=requestTMDB(ctx,path+"/watch/providers",url.Values{},key,&providers);err!=nil {note="La disponibilidad no se ha podido consultar"}
  es:=providers.Results["ES"]
  _=json.NewEncoder(w).Encode(map[string]any{
   "id":detail.ID,"title":detail.Title,"original_title":detail.OriginalTitle,"release_date":detail.ReleaseDate,
   "overview":detail.Overview,"poster_path":detail.PosterPath,"vote_average":detail.VoteAverage,
   "vote_count":detail.VoteCount,"runtime":detail.Runtime,"genres":genres,"director":director,"cast":cast,"trailer":trailer,
   "subscriptions":labProviderNames(es.Flatrate),"free":labProviderNames(es.Free),"ads":labProviderNames(es.Ads),
   "rent":labProviderNames(es.Rent),"buy":labProviderNames(es.Buy),"watch_link":es.Link,
   "availability_note":note,"tmdb_url":fmt.Sprintf("https://www.themoviedb.org/movie/%d",id),
  })
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
