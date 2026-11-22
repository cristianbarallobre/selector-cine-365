//go:build windows

package main

import (
 "context"
 "encoding/json"
 "errors"
 "fmt"
 "net/http"
 "net/url"
 "os"
 "path/filepath"
 "strconv"
 "strings"
 "sync"
 "time"
 "unicode/utf8"
)

// Descargas de fichas aisladas del JSON de listas y de cualquier carpeta de OneDrive.
const (
 labMetaTTL=30*24*time.Hour
 labExtraTTL=7*24*time.Hour
 labProvidersTTL=24*time.Hour
 labPrefetchDelay=900*time.Millisecond
)
type labFilmJob struct {
 Title string `json:"title"`
 Year int `json:"year"`
 ID int `json:"id"`
}
type labDetailRecord struct {
 Data map[string]any `json:"data"`
 MetaAt int64 `json:"meta_at"`
 ExtraAt int64 `json:"extra_at"`
 ProvidersAt int64 `json:"providers_at"`
}
type labCacheDisk struct {
 Version int `json:"version"`
 Items map[string]labDetailRecord `json:"items"`
 Resolved map[string]int `json:"resolved"`
}
type labCatalogCache struct {
 mu sync.Mutex
 file string
 app *runningApp
 disk labCacheDisk
 targets map[string]labFilmJob
 queue []labFilmJob
 queued map[string]bool
 running bool
 paused bool
 busy string
 errors int
 lastError string
}
var labCatalogMu sync.Mutex
var labCatalogByPath=map[string]*labCatalogCache{}
func (a *runningApp) labCatalog() *labCatalogCache {
 path:=filepath.Join(filepath.Dir(a.savePath),"tmdb_fichas_lab.json")
 labCatalogMu.Lock();defer labCatalogMu.Unlock()
 if c:=labCatalogByPath[path];c!=nil{return c}
 c:=&labCatalogCache{file:path,app:a,disk:labCacheDisk{Version:1,Items:map[string]labDetailRecord{},Resolved:map[string]int{}},targets:map[string]labFilmJob{},queued:map[string]bool{}}
 if data,err:=os.ReadFile(path);err==nil&&len(data)<25<<20 {
  var saved labCacheDisk
  if json.Unmarshal(data,&saved)==nil&&saved.Version==1 {
   if saved.Items!=nil{c.disk.Items=saved.Items}
   if saved.Resolved!=nil{c.disk.Resolved=saved.Resolved}
  }
 }
 labCatalogByPath[path]=c
 return c
}
func labID(id int) string{return strconv.Itoa(id)}
func labJobKey(j labFilmJob) string {
 // Stable for the entire queue, even after resolving a title to an ID.
 if strings.TrimSpace(j.Title)!=""{return normaliseTitle(j.Title)+"|"+strconv.Itoa(j.Year)}
 return "tmdb:"+strconv.Itoa(j.ID)
}
func (c *labCatalogCache) fileSaveLocked() error {
 if err:=os.MkdirAll(filepath.Dir(c.file),0700);err!=nil{return err}
 data,err:=json.Marshal(c.disk);if err!=nil{return err}
 if len(data)>25<<20{return errors.New("Caché de fichas demasiado grande")}
 f,err:=os.CreateTemp(filepath.Dir(c.file),".fichas-lab-*");if err!=nil{return err}
 defer os.Remove(f.Name())
 if err=f.Chmod(0600);err!=nil{_ =f.Close();return err}
 if _,err=f.Write(data);err!=nil{_ =f.Close();return err}
 if err=f.Close();err!=nil{return err}
 return os.Rename(f.Name(),c.file)
}
func (c *labCatalogCache) record(id int) (labDetailRecord,bool) {
 c.mu.Lock();defer c.mu.Unlock()
 v,ok:=c.disk.Items[labID(id)]
 return v,ok&&len(v.Data)>0
}
func (c *labCatalogCache) remember(id int,rec labDetailRecord,identity string) {
 if id<=0||len(rec.Data)==0{return}
 c.mu.Lock();defer c.mu.Unlock()
 if len(c.disk.Items)>8000{return}
 c.disk.Items[labID(id)]=rec
 if identity!=""&&len(c.disk.Resolved)<12000{c.disk.Resolved[identity]=id}
 if err:=c.fileSaveLocked();err!=nil{c.lastError="No se pudo guardar la caché: "+err.Error()}
}
func labExpired(at int64,ttl time.Duration) bool{return at<=0||time.Since(time.Unix(at,0))>=ttl}
func labNeedsRefresh(r labDetailRecord) bool {
 return len(r.Data)==0||labExpired(r.MetaAt,labMetaTTL)||labExpired(r.ExtraAt,labExtraTTL)||labExpired(r.ProvidersAt,labProvidersTTL)
}
func labCopyMap(src map[string]any) map[string]any {
 dst:=make(map[string]any,len(src)+4)
 for k,v:=range src{dst[k]=v}
 return dst
}
func labFreshDetail(ctx context.Context,id int,key string,prior labDetailRecord) (labDetailRecord,error) {
 if id<=0||id>999999999{return prior,errors.New("ID TMDB no válido")}
 now:=time.Now().Unix()
 rec:=prior
 rec.Data=labCopyMap(prior.Data)
 path:="/movie/"+labID(id)
 if len(rec.Data)==0||labExpired(prior.MetaAt,labMetaTTL) {
  var d struct {
   labMovie
   Runtime int `json:"runtime"`
   Genres []labGenre `json:"genres"`
   Credits struct {Cast []labPerson `json:"cast"`;Crew []labPerson `json:"crew"`} `json:"credits"`
   Videos struct {Results []labVideo `json:"results"`} `json:"videos"`
  }
  if err:=requestTMDB(ctx,path,url.Values{"language":{"es-ES"},"append_to_response":{"credits,videos"}},key,&d);err!=nil{return prior,err}
  if d.ID!=id||d.Title==""{return prior,errors.New("TMDB no devolvió una película válida")}
  d.labMovie=safeLabMovie(d.labMovie)
  director:=""
  for _,p:=range d.Credits.Crew {if p.Job=="Director"&&p.Name!="" {if director!=""{director+=" / "};director+=p.Name;if len(director)>220{break}}}
  cast:=make([]map[string]string,0,10)
  for _,p:=range d.Credits.Cast {if len(cast)>=10{break};if p.Name!=""{cast=append(cast,map[string]string{"name":p.Name,"character":p.Character})}}
  genres:=make([]string,0,len(d.Genres))
  for _,g:=range d.Genres{if g.Name!=""{genres=append(genres,g.Name)}}
  trailer:=labTrailer(d.Videos.Results)
  if trailer=="" {
   var en struct{Results []labVideo `json:"results"`}
   if requestTMDB(ctx,path+"/videos",url.Values{"language":{"en-US"}},key,&en)==nil{trailer=labTrailer(en.Results)}
  }
  rec.Data["id"]=d.ID;rec.Data["title"]=d.Title;rec.Data["original_title"]=d.OriginalTitle
  rec.Data["release_date"]=d.ReleaseDate;rec.Data["overview"]=d.Overview
  rec.Data["poster_path"]=d.PosterPath;rec.Data["vote_average"]=d.VoteAverage
  rec.Data["vote_count"]=d.VoteCount;rec.Data["runtime"]=d.Runtime
  rec.Data["genres"]=genres;rec.Data["director"]=director;rec.Data["cast"]=cast
  rec.Data["trailer"]=trailer;rec.Data["tmdb_url"]=fmt.Sprintf("https://www.themoviedb.org/movie/%d",id)
  rec.MetaAt=now;rec.ExtraAt=now
 } else if labExpired(prior.ExtraAt,labExtraTTL) {
  var d struct {VoteAverage float64 `json:"vote_average"`;VoteCount int `json:"vote_count"`;Videos struct {Results []labVideo `json:"results"`} `json:"videos"`}
  if requestTMDB(ctx,path,url.Values{"language":{"es-ES"},"append_to_response":{"videos"}},key,&d)==nil {
   rec.Data["vote_average"]=d.VoteAverage;rec.Data["vote_count"]=d.VoteCount
   if trailer:=labTrailer(d.Videos.Results);trailer!=""{rec.Data["trailer"]=trailer}
   rec.ExtraAt=now
  }
 }
 if len(rec.Data)==0{return prior,errors.New("Ficha vacía")}
 if labExpired(prior.ProvidersAt,labProvidersTTL) {
  var p struct{Results map[string]labProviders `json:"results"`}
  if err:=requestTMDB(ctx,path+"/watch/providers",url.Values{},key,&p);err==nil{
   es:=p.Results["ES"]
   rec.Data["subscriptions"]=labProviderNames(es.Flatrate)
   rec.Data["free"]=labProviderNames(es.Free);rec.Data["ads"]=labProviderNames(es.Ads)
   rec.Data["rent"]=labProviderNames(es.Rent);rec.Data["buy"]=labProviderNames(es.Buy)
   rec.Data["watch_link"]=es.Link;rec.Data["availability_note"]=""
   rec.ProvidersAt=now
  }else if prior.ProvidersAt==0 {
   rec.Data["availability_note"]="Plataformas temporalmente no disponibles."
  }
 }
 return rec,nil
}
func labResolveMovie(ctx context.Context,j labFilmJob,key string) (int,error) {
 if j.ID>0{return j.ID,nil}
 var result struct{Results []struct{
  ID int `json:"id"`
  Title string `json:"title"`
  Original string `json:"original_title"`
  Released string `json:"release_date"`
 } `json:"results"`}
 params:=url.Values{"language":{"es-ES"},"query":{j.Title},"include_adult":{"false"}}
 if err:=requestTMDB(ctx,"/search/movie",params,key,&result);err!=nil{return 0,err}
 id,_:=chooseMovie(j.Title,j.Year,result.Results)
 if id<=0{return 0,errors.New("sin identificación TMDB fiable")}
 return id,nil
}
func (c *labCatalogCache) enqueue(jobs []labFilmJob) {
 c.mu.Lock()
 for _,j:=range jobs {
  j.Title=strings.TrimSpace(j.Title)
  if j.ID<=0&&(utf8.RuneCountInString(j.Title)<2||utf8.RuneCountInString(j.Title)>150){continue}
  if j.Year<0||j.Year>2200||j.ID<0||j.ID>999999999{continue}
  k:=labJobKey(j)
  if k==""{continue}
  if len(c.targets)>=1800{if _,ok:=c.targets[k];!ok{continue}}
  c.targets[k]=j
  if j.ID==0{j.ID=c.disk.Resolved[k]}
  existing:=c.disk.Items[labID(j.ID)]
  if j.ID>0&&!labNeedsRefresh(existing){continue}
  if c.queued[k]||c.busy==k{continue}
  c.queue=append(c.queue,j);c.queued[k]=true
 }
 start:=!c.running&&!c.paused&&len(c.queue)>0
 if start{c.running=true}
 c.mu.Unlock()
 if start{go c.runWorker()}
}
func (c *labCatalogCache) queueRefresh(id int) {
 if id>0{c.enqueue([]labFilmJob{{ID:id}})}
}
func (c *labCatalogCache) status() map[string]any {
 c.mu.Lock();defer c.mu.Unlock()
 cached:=0
 for k,j:=range c.targets{
  id:=j.ID
  if id==0{id=c.disk.Resolved[k]}
  if id>0&&len(c.disk.Items[labID(id)].Data)>0{cached++}
 }
 return map[string]any{
  "total":len(c.targets),"cached":cached,"pending":len(c.queue),
  "running":c.running,"paused":c.paused,"current":c.busy,
  "errors":c.errors,"last_error":c.lastError,
 }
}
func (c *labCatalogCache) runWorker() {
 for {
  c.mu.Lock()
  if c.paused||len(c.queue)==0 {
   c.running=false;c.busy="";c.mu.Unlock();return
  }
  j:=c.queue[0];c.queue=c.queue[1:]
  k:=labJobKey(j);delete(c.queued,k);c.busy=k
  if j.ID==0{j.ID=c.disk.Resolved[k]}
  c.mu.Unlock()
  select{case <-c.app.stopCh: c.mu.Lock();c.running=false;c.busy="";c.mu.Unlock();return;default:}
  key:=c.app.platforms.getKey()
  var err error
  if key=="" {err=errors.New("TMDB sin conectar")}
  if err==nil {
   ctx,cancel:=context.WithTimeout(context.Background(),35*time.Second)
   if j.ID==0 {j.ID,err=labResolveMovie(ctx,j,key)}
   if err==nil {
    prior,_:=c.record(j.ID)
    if labNeedsRefresh(prior){
     var fresh labDetailRecord
     fresh,err=labFreshDetail(ctx,j.ID,key,prior)
     if err==nil{c.remember(j.ID,fresh,k)}
    }else{c.mu.Lock();c.disk.Resolved[k]=j.ID;c.mu.Unlock()}
   }
   cancel()
  }
  c.mu.Lock()
  if err!=nil {c.errors++;c.lastError=err.Error()}
  c.busy=""
  c.mu.Unlock()
  delay:=labPrefetchDelay
  if err!=nil&&strings.Contains(err.Error(),"Límite temporal"){delay=60*time.Second;c.enqueue([]labFilmJob{j})}
  select{case <-c.app.stopCh: c.mu.Lock();c.running=false;c.mu.Unlock();return;case <-time.After(delay):}
 }
}
func (a *runningApp) serveLabPrefetch(w http.ResponseWriter,r *http.Request) {
 if !appIsLab||!a.auth(r){http.Error(w,"No autorizado",403);return}
 w.Header().Set("Content-Type","application/json; charset=utf-8")
 c:=a.labCatalog()
 switch r.URL.Path {
 case "/api/discover/prefetch":
  if r.Method==http.MethodPost {
   var b struct{Films []labFilmJob `json:"films"`}
   if err:=json.NewDecoder(http.MaxBytesReader(w,r.Body,250000)).Decode(&b);err!=nil||len(b.Films)>1800{http.Error(w,"Lista de precarga no válida",400);return}
   if a.platforms.getKey()!=""{c.enqueue(b.Films)}
  }else if r.Method!=http.MethodGet {http.Error(w,"Método no permitido",405);return}
 case "/api/discover/prefetch/pause":
  if r.Method!=http.MethodPost{http.Error(w,"Método no permitido",405);return}
  var b struct{Paused bool `json:"paused"`}
  if json.NewDecoder(http.MaxBytesReader(w,r.Body,1024)).Decode(&b)!=nil {http.Error(w,"JSON inválido",400);return}
  c.mu.Lock();c.paused=b.Paused
  restart:=!b.Paused&&!c.running&&len(c.queue)>0
  if restart{c.running=true}
  c.mu.Unlock()
  if restart{go c.runWorker()}
 default:http.NotFound(w,r);return
 }
 _=json.NewEncoder(w).Encode(c.status())
}
// Un título preparado puede abrir su ficha sin buscar de nuevo en TMDB.
func (a *runningApp) labCachedIdentity(w http.ResponseWriter,r *http.Request) {
 if r.Method!=http.MethodGet {http.Error(w,"Método no permitido",405);return}
 title:=strings.TrimSpace(r.URL.Query().Get("title"))
 year,err:=strconv.Atoi(r.URL.Query().Get("year"))
 if err!=nil||year<0||year>2200||utf8.RuneCountInString(title)<2||utf8.RuneCountInString(title)>150 {
  http.Error(w,"Título o año no válido",400);return
 }
 c:=a.labCatalog()
 k:=labJobKey(labFilmJob{Title:title,Year:year})
 c.mu.Lock()
 id:=c.disk.Resolved[k]
 if len(c.disk.Items[labID(id)].Data)==0 {id=0}
 c.mu.Unlock()
 w.Header().Set("Content-Type","application/json; charset=utf-8")
 _=json.NewEncoder(w).Encode(map[string]int{"id":id})
}

func (a *runningApp) labDetailFromCache(w http.ResponseWriter,r *http.Request) {
 w.Header().Set("Content-Type","application/json; charset=utf-8")
 id,err:=strconv.Atoi(r.URL.Query().Get("id"))
 if err!=nil||id<=0||id>999999999{http.Error(w,"ID de película no válido",400);return}
 c:=a.labCatalog()
 if prior,ok:=c.record(id);ok{
  data:=labCopyMap(prior.Data)
  data["cache_state"]="guardada"
  if labNeedsRefresh(prior){
   data["cache_state"]="actualizando"
   if a.platforms.getKey()!=""{c.queueRefresh(id)}
  }
  _=json.NewEncoder(w).Encode(data)
  return
 }
 key:=a.platforms.getKey()
 if key==""{http.Error(w,"Conecta TMDB; esta ficha aún no está en caché.",412);return}
 ctx,cancel:=context.WithTimeout(r.Context(),25*time.Second);defer cancel()
 rec,err:=labFreshDetail(ctx,id,key,labDetailRecord{})
 if err!=nil{http.Error(w,err.Error(),502);return}
 c.remember(id,rec,"")
 data:=labCopyMap(rec.Data);data["cache_state"]="descargada"
 _=json.NewEncoder(w).Encode(data)
}
