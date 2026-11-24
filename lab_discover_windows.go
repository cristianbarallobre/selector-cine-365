//go:build windows

package main

import (
 "context"
 "embed"
 "encoding/json"
 "io"
 "net/http"
 "net/url"
 "sort"
 "strconv"
 "strings"
 "time"
 "unicode/utf8"
)

//go:embed lab_discover_windows.html
var labDiscoveryFiles embed.FS

const labDiscoverySeparator = "<!-- LAB_DISCOVERY_SCRIPT -->"

func injectLabDiscovery(page string) string {
 // Git convierte los HTML a CRLF en runners Windows; unificar LF ANTES
 // de insertar JS, evitando dejar botones sin sus manejadores.
 page=strings.ReplaceAll(page,"\r\n","\n")
 raw, err := labDiscoveryFiles.ReadFile("lab_discover_windows.html")
 if err != nil { return page }
 fragments := strings.SplitN(strings.ReplaceAll(string(raw),"\r\n","\n"), labDiscoverySeparator, 2)
 if len(fragments)!=2 { return page }
 nav := `<button class="nav-btn" type="button" id="labDiscoverNav">✦ Descubrir <span class="lab-nav-tag">LAB</span></button>`
 page = strings.Replace(page, `<button class="nav-btn" type="button" data-view="statistics">`, nav+`<button class="nav-btn" type="button" data-view="statistics">`, 1)
 page = strings.Replace(page, "🧪 SELECTOR CINE 365 LAB · v9 Beta 14 Preview", "🧪 SELECTOR CINE 365 · LAB-02.15 WINDOWS", 1)
 // Los botones específicos se inyectan solo en Windows LAB; plantilla.html permanece intacta.
 // En Windows LAB, el enlace externo de cada fila sobra: ya está dentro de Detalles.
 page = strings.Replace(page, "b.append(togg,pill,fix,link);frag.appendChild(b);", "b.append(togg,pill,labCreateDetailButton(f),fix);frag.appendChild(b);", 1)
 page = strings.Replace(page, `<input type="search" id="search"`, `<div class="lab-prefetch-banner" id="labPrefetchBanner" role="status" aria-live="polite"><span id="labPrefetchText">Fichas TMDB: comprobando caché…</span><button type="button" id="labPrefetchPause" hidden>Pausar</button></div><input type="search" id="search"`, 1)
 // La precarga se inicia tras cargar biblioteca e identificaciones TMDB, y se actualiza al editar listas.
 page = strings.Replace(page, "void cloudAutoReceive(true);});", "void cloudAutoReceive(true);labQueuePrefetchSoon();});", 1)
 page = strings.Replace(page, "function scheduleSave(){membershipCache=null;", "function scheduleSave(){labQueuePrefetchSoon();membershipCache=null;", 1)
 page = strings.Replace(page, "function showView(name){if(", "function showView(name){labLeave();if(name==='library')labQueuePrefetchSoon();if(", 1)
 // La confirmación nativa confirm() revela 127.0.0.1 y el puerto local.
 // Solo Windows LAB usa el diálogo visual; el código de borrado y sincronización permanece intacto.
 page = strings.Replace(page, "function deleteList(id){", "async function deleteList(id){", 1)
 page = strings.Replace(page,
  "if(!confirm('¿Eliminar «'+l.name+'»? Esta lista dejará de aparecer en la aplicación.'))return;",
  "if(!await labConfirmDeleteList(l))return; if(libraryState.lists.length<2||!libraryState.lists.some(item=>item.id===id))return;", 1)

 // El gesto de arrastre nativo se interrumpe al reconstruir las filas durante
 // el cambio de página. En Windows LAB usamos Pointer Events con captura en
 // document.body, conservando el gesto incluso tras renderList().
 page = strings.Replace(page,
  " handle.draggable=true;",
  " handle.draggable=false;",
  1)
 // Los números de «Listas» deben abrir la lista íntegra y señalar la película,
 // nunca introducir su título en el buscador de Biblioteca.
 page = strings.Replace(page,
  `function openFilmInList(m){
 if(!m)return;
 if(m.id!==libraryState.active_id||virtualMode)switchCollection(m.id);
 selectSearchScope('current');search=m.title;$('search').value=search;
 filmPage=0;showView('library');renderList();
}`,
  `function openFilmInList(m){labJumpToFilmList(m);}`,1)
 // El resaltado también se aplica al volver a dibujar las filas por cambios de
 // caché, plataformas o tamaño durante los segundos de localización.
 page = strings.Replace(page,
  "b.dataset.idx=f.idx;b.setAttribute('role','listitem');",
  "b.dataset.idx=f.idx;b.setAttribute('role','listitem');if(labListJumpHighlight&&(labListJumpHighlight.listID===libraryState.active_id||labListJumpHighlight.listID===labFavoriteID&&labFavoritesMode)&&labListJumpHighlight.rawIndex===f.idx)b.classList.add('lab-list-jump-target');",
  1)
 // Aprovecha el hueco sin reducir el número de filas: conserva la capacidad de 96px+7px.
 // Cuando sobra espacio, estira proporcionalmente hasta 110px; en ventanas ajustadas mantiene 96px.
 page = strings.Replace(page,
  "const pageSize=Math.max(1,Math.floor((areaHeight+7)/103));",
  "const labUsableHeight=Math.max(0,areaHeight-2);const pageSize=Math.max(1,Math.floor((labUsableHeight+7)/103));const labRowHeight=Math.min(110,Math.max(96,Math.floor((labUsableHeight-7*(pageSize-1))/pageSize)));parent.style.setProperty('--lab-film-row-height',labRowHeight+'px');",
  1)
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
// Géneros de película TMDB; validar la lista de IDs antes de construir
// consultas externas (nunca reenviar parámetros arbitrarios desde la UI).
var labDiscoverGenreIDs=map[int]bool{
 28:true,12:true,16:true,35:true,80:true,99:true,18:true,10751:true,
 14:true,36:true,27:true,10402:true,9648:true,10749:true,878:true,
 10770:true,53:true,10752:true,37:true,
}
type labGenre struct { Name string `json:"name"` }
type labPerson struct { ID int `json:"id"`; Name string `json:"name"`; Character string `json:"character"`; Job string `json:"job"`; ProfilePath string `json:"profile_path"` }
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
 if !a.auth(r) {http.Error(w,"No autorizado",403);return}
 if r.URL.Path=="/api/discover/prefetch"||r.URL.Path=="/api/discover/prefetch/pause" {a.serveLabPrefetch(w,r);return}
 if r.Method!=http.MethodGet {http.Error(w,"Método no permitido",405);return}
 if r.URL.Path=="/api/discover/detail" {a.labDetailFromCache(w,r);return}
 if r.URL.Path=="/api/discover/person" {a.labPersonDetail(w,r);return}
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
   // Géneros filtrados directamente en TMDB, sin consultas por película.
   genreText:=r.URL.Query().Get("genre")
   genreID:=0
   if genreText!=""&&genreText!="any" {
    n,err:=strconv.Atoi(genreText)
    if err!=nil||!labDiscoverGenreIDs[n] {fail(400,"Género no válido");return}
    genreID=n
   }
   if genreID>0 {
    endpoint="/discover/movie"
    params.Set("region","ES")
    params.Set("sort_by","popularity.desc")
    if genreID>0 {params.Set("with_genres",strconv.Itoa(genreID))}
    switch kind {
    case "top":params.Set("sort_by","vote_average.desc");params.Set("vote_count.gte","200")
    case "upcoming":
     params.Set("sort_by","primary_release_date.asc")
     params.Set("primary_release_date.gte",time.Now().Format("2006-01-02"))
    }
   }
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

  
// Fichas de intérpretes de TMDB: se sirven desde el backend LAB, nunca con la API key en el navegador.
type labPersonFilm struct {
 ID int `json:"id"`
 Title string `json:"title"`
 ReleaseDate string `json:"release_date"`
 Character string `json:"character"`
 Popularity float64 `json:"popularity"`
 VoteAverage float64 `json:"vote_average"`
 VoteCount int `json:"vote_count"`
 CastOrder int `json:"order"`
}
func (a *runningApp) labPersonDetail(w http.ResponseWriter,r *http.Request) {
 w.Header().Set("Content-Type","application/json; charset=utf-8")
 fail:=func(code int,msg string){w.WriteHeader(code);_ =json.NewEncoder(w).Encode(map[string]string{"error":msg})}
 id,err:=strconv.Atoi(r.URL.Query().Get("id"))
 if err!=nil||id<=0||id>999999999{fail(400,"ID de persona no válido");return}
 key:=a.platforms.getKey()
 if key==""{fail(412,"Conecta TMDB para consultar intérpretes.");return}
 ctx,cancel:=context.WithTimeout(r.Context(),22*time.Second);defer cancel()
 var p struct {
  ID int `json:"id"`
  Name string `json:"name"`
  Biography string `json:"biography"`
  Birthday string `json:"birthday"`
  Deathday string `json:"deathday"`
  PlaceOfBirth string `json:"place_of_birth"`
  KnownFor string `json:"known_for_department"`
  ProfilePath string `json:"profile_path"`
  Credits struct{Cast []labPersonFilm `json:"cast"`} `json:"movie_credits"`
 }
 params:=url.Values{"language":{"es-ES"},"append_to_response":{"movie_credits"}}
 if err:=requestTMDB(ctx,"/person/"+strconv.Itoa(id),params,key,&p);err!=nil{fail(502,err.Error());return}
 if p.ID!=id||strings.TrimSpace(p.Name)==""{fail(502,"TMDB no devolvió una persona válida");return}
 // TMDB devuelve biografías localizadas por separado. Si falta es-ES,
 // intentar otra variante en español antes de ofrecer el original inglés
 // claramente identificado. No mostrar automáticamente inglés como español.
 biographyLanguage:="es-ES"
 originalBiography:=""
 if strings.TrimSpace(p.Biography)==""{
  var spanish struct{Biography string `json:"biography"`}
  if requestTMDB(ctx,"/person/"+strconv.Itoa(id),url.Values{"language":{"es-MX"}},key,&spanish)==nil&&strings.TrimSpace(spanish.Biography)!=""{
   p.Biography=spanish.Biography
   biographyLanguage="es-MX"
  }else{
   biographyLanguage="unavailable"
   var english struct{Biography string `json:"biography"`}
   if requestTMDB(ctx,"/person/"+strconv.Itoa(id),url.Values{"language":{"en-US"}},key,&english)==nil{
    originalBiography=english.Biography
   }
  }
 }
 limit:=func(s string,n int)string{v:=[]rune(strings.TrimSpace(s));if len(v)>n{v=v[:n]};return string(v)}
 films:=make([]labPersonFilm,0,len(p.Credits.Cast))
 seen:=map[int]bool{}
 for _,film:=range p.Credits.Cast {
  if film.ID<=0||film.ID>999999999||seen[film.ID]||strings.TrimSpace(film.Title)==""{continue}
  seen[film.ID]=true
  film.Title=limit(film.Title,180)
  film.Character=limit(film.Character,100)
  films=append(films,film)
 }
 // La clasificación se hace en la interfaz sobre toda la filmografía.
 // No recortar a 16 aquí: una obra importante podría quedar fuera de los
 // primeros puestos por popularidad, bloqueando otros criterios.
 sort.SliceStable(films,func(i,j int)bool{
  if films[i].ID!=films[j].ID{return films[i].ID<films[j].ID}
  return films[i].Title<films[j].Title
 })
 if !validPosterPath(p.ProfilePath){p.ProfilePath=""}
 _=json.NewEncoder(w).Encode(map[string]any{
  "id":p.ID,"name":limit(p.Name,180),"biography":limit(p.Biography,6500),
  "biography_language":biographyLanguage,"biography_original":limit(originalBiography,6500),
  "birthday":limit(p.Birthday,10),"deathday":limit(p.Deathday,10),
  "place_of_birth":limit(p.PlaceOfBirth,200),"known_for_department":limit(p.KnownFor,80),
  "profile_path":p.ProfilePath,"movies":films,
  "tmdb_url":"https://www.themoviedb.org/person/"+strconv.Itoa(id),
 })
}
