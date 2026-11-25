package main

import (
 "crypto/rand"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "fmt"
 "os"
 "path/filepath"
 "regexp"
 "sort"
 "strconv"
 "strings"
 "sync"
 "time"
)

const syncProtocol = 2

type syncRevision struct {
 Checked bool `json:"checked"`
 Device string `json:"device"`
 Seq int64 `json:"seq"`
 Clock map[string]int64 `json:"clock"`
 ChangedAt int64 `json:"changed_at"`
}
type syncHistoryEntry struct {
 At int64 `json:"at"`
 ChangedAt int64 `json:"changed_at,omitempty"`
 Key string `json:"key"`
 Checked bool `json:"checked"`
 Device string `json:"device"`
 Origin string `json:"origin"`
}
type syncDeviceState struct {
 History []syncHistoryEntry `json:"history,omitempty"`
 Protocol int `json:"protocol"`
 DeviceID string `json:"device_id"`
 Seq int64 `json:"seq"`
 Known map[string]int64 `json:"known"`
 Revisions map[string]syncRevision `json:"revisions"`
 PairedFolder string `json:"paired_folder,omitempty"`
}
type syncSnapshot struct {
 Protocol int `json:"protocol"`
 DeviceID string `json:"device_id"`
 UpdatedAt int64 `json:"updated_at,omitempty"`
 Library libraryState `json:"library"`
 Known map[string]int64 `json:"known,omitempty"`
 Revisions map[string]syncRevision `json:"revisions,omitempty"`
}
type syncConflict struct {
 Key string `json:"key"`
 Local bool `json:"local"`
 Remote bool `json:"remote"`
 LocalDevice string `json:"local_device"`
 RemoteDevice string `json:"remote_device"`
}
type syncPlan struct {
 Library libraryState
 Revisions map[string]syncRevision
 Known map[string]int64
 Conflicts []syncConflict
 Sources int
 Legacy bool
 Token string
}
type desktopSyncStore struct {
 mu sync.Mutex
 path string
 state syncDeviceState
}

func newDesktopSyncStore(dir string) (*desktopSyncStore,error) {
 s:=&desktopSyncStore{path:filepath.Join(dir,"sync_dispositivo_v2.json")}
 data,err:=os.ReadFile(s.path)
 if err==nil {
  if e:=json.Unmarshal(data,&s.state);e!=nil{return nil,fmt.Errorf("estado de sincronización dañado (no se ha sobrescrito): %w",e)}
 } else if !os.IsNotExist(err){return nil,err}
 if s.state.DeviceID=="" {
  seed:=make([]byte,16);if _,err:=rand.Read(seed);err!=nil{return nil,err}
  s.state.DeviceID="desktop-"+hex.EncodeToString(seed)
 }
 s.state.Protocol=syncProtocol
 if s.state.Known==nil{s.state.Known=map[string]int64{}}
 if s.state.Revisions==nil{s.state.Revisions=map[string]syncRevision{}}
 if err:=s.saveLocked();err!=nil{return nil,err}
 return s,nil
}
func (s *desktopSyncStore) saveLocked() error {
 b,e:=json.MarshalIndent(s.state,"","  ");if e!=nil{return e}
 return cloudWriteFile(s.path,b)
}
func (s *desktopSyncStore) appendHistoryLocked(key string,rev syncRevision,origin string) {
 entry:=syncHistoryEntry{
  At:time.Now().UnixMilli(),ChangedAt:rev.ChangedAt,
  Key:key,Checked:rev.Checked,Device:rev.Device,Origin:origin,
 }
 s.state.History=append(s.state.History,entry)
 if len(s.state.History)>80{s.state.History=s.state.History[len(s.state.History)-80:]}
}
func (s *desktopSyncStore) history() []syncHistoryEntry {
 s.mu.Lock();defer s.mu.Unlock()
 events:=make([]syncHistoryEntry,0,len(s.state.History))
 for i:=len(s.state.History)-1;i>=0;i--{events=append(events,s.state.History[i])}
 return events
}
func (s *desktopSyncStore) isPaired(folder string) bool {
 if s==nil{return false}
 s.mu.Lock();defer s.mu.Unlock()
 return folder!=""&&s.state.PairedFolder==filepath.Clean(folder)
}
func (s *desktopSyncStore) parts() (string,map[string]syncRevision,map[string]int64) {
 s.mu.Lock();defer s.mu.Unlock()
 rev:=make(map[string]syncRevision,len(s.state.Revisions))
 for k,v:=range s.state.Revisions {rev[k]=v}
 known:=make(map[string]int64,len(s.state.Known))
 for k,v:=range s.state.Known {known[k]=v}
 return s.state.DeviceID,rev,known
}

var syncMark = regexp.MustCompile("^[☐□☑☒✓✔✅]\\s*")
var syncNumber = regexp.MustCompile("^(?:[0-9]{1,4}[.)]|[-•])\\s*")
var syncYear = regexp.MustCompile("\\s+\\(([0-9]{4})\\)$")
var syncPunctuation = regexp.MustCompile("[^a-z0-9]+")
var syncAccents = strings.NewReplacer("á","a","à","a","ä","a","â","a","é","e","è","e","ë","e","ê","e","í","i","ì","i","ï","i","ó","o","ò","o","ö","o","ô","o","ú","u","ù","u","ü","u","û","u","ñ","n","ç","c","æ","ae","œ","oe")

type syncFilm struct{ key string; checked bool }
func parseSyncFilm(line string) (syncFilm,bool) {
 text:=strings.TrimSpace(line)
 if text==""||strings.HasPrefix(text,"#")||strings.HasPrefix(strings.ToUpper(text),"SELECTOR CINE 365")||
   strings.HasPrefix(strings.ToUpper(text),"ATLAS DEL CINE"){return syncFilm{},false}
 trimmed:=syncMark.ReplaceAllString(text,"")
 trimmed=syncNumber.ReplaceAllString(trimmed,"")
 year:=""
 if idx:=syncYear.FindStringSubmatchIndex(trimmed);idx!=nil{
  year=trimmed[idx[2]:idx[3]];trimmed=trimmed[:idx[0]]
 }
 normalized:=syncPunctuation.ReplaceAllString(syncAccents.Replace(strings.ToLower(strings.TrimSpace(trimmed)))," ")
 normalized=strings.Join(strings.Fields(normalized)," ")
 if normalized==""{return syncFilm{},false}
 return syncFilm{normalized+"|"+year,strings.ContainsAny(string([]rune(text)[0]),"☑☒✓✔✅")},true
}
func syncStatuses(l libraryState) map[string]bool {
 out:=map[string]bool{}
 for _,list:=range l.Lists {
  for _,line:=range strings.Split(list.Text,"\n") {
   if f,ok:=parseSyncFilm(line);ok {out[f.key]=out[f.key]||f.checked}
  }
 }
 return out
}
func syncWatched(l libraryState) int {
 n:=0;for _,v:=range syncStatuses(l){if v{n++}}
 return n
}
func syncCheckedLine(line string,seen bool) string {
 // Mantener título, número y año tal como los escribió el usuario.
 if _,ok:=parseSyncFilm(line);!ok{return line}
 cleaned:=syncMark.ReplaceAllString(strings.TrimSpace(line),"")
 if seen{return "☑ "+cleaned};return "☐ "+cleaned
}
func syncCopyMap(src map[string]int64) map[string]int64 {
 dst:=map[string]int64{};for k,v:=range src{dst[k]=v};return dst
}
func syncJoin(a,b map[string]int64) map[string]int64 {
 out:=syncCopyMap(a);for k,v:=range b{if v>out[k]{out[k]=v}};return out
}
func syncCompare(a,b map[string]int64) int {
 newer,older:=false,false
 keys:=map[string]bool{}
 for k:=range a{keys[k]=true};for k:=range b{keys[k]=true}
 for k:=range keys{if a[k]>b[k]{newer=true}else if a[k]<b[k]{older=true}}
 if newer&&older{return 2};if newer{return 1};if older{return -1};return 0
}
func (s *desktopSyncStore) observe(before,after libraryState) error {
 if s==nil{return nil}
 s.mu.Lock();defer s.mu.Unlock()
 old,current:=syncStatuses(before),syncStatuses(after)
 dirty:=false
 for key,next:=range current {
  previous,exists:=old[key]
  if !exists { if !next {continue} } else if previous==next {continue}
  s.state.Seq++
  vector:=syncCopyMap(s.state.Known)
  vector[s.state.DeviceID]=s.state.Seq
  s.state.Known=syncJoin(s.state.Known,vector)
  revision:=syncRevision{Checked:next,Device:s.state.DeviceID,Seq:s.state.Seq,Clock:vector,ChangedAt:time.Now().UnixMilli()}
  s.state.Revisions[key]=revision
  s.appendHistoryLocked(key,revision,"local")
  dirty=true
 }
 if dirty{return s.saveLocked()};return nil
}
func (s *desktopSyncStore) snapshot(l libraryState) ([]byte,error) {
 if s==nil{return nil,errors.New("sincronización no inicializada")}
 id,rev,known:=s.parts()
 return json.MarshalIndent(syncSnapshot{Protocol:syncProtocol,DeviceID:id,UpdatedAt:time.Now().UnixMilli(),Library:l,Known:known,Revisions:rev},"","  ")
}
func (s *desktopSyncStore) accept(p syncPlan,decisions map[string]bool,folder string) error {
 if s==nil{return errors.New("sincronización no inicializada")}
 s.mu.Lock();defer s.mu.Unlock()
 known:=syncJoin(p.Known,s.state.Known)
 revisions:=make(map[string]syncRevision,len(p.Revisions))
 for k,v:=range p.Revisions{revisions[k]=v}
 for key,seen:=range decisions{
  s.state.Seq++
  vector:=syncCopyMap(known);vector[s.state.DeviceID]=s.state.Seq
  known=syncJoin(known,vector)
  revisions[key]=syncRevision{Checked:seen,Device:s.state.DeviceID,Seq:s.state.Seq,Clock:vector,ChangedAt:time.Now().UnixMilli()}
 }
 for key,rev:=range revisions{
  old,had:=s.state.Revisions[key]
  if !had||old.Checked!=rev.Checked||old.Device!=rev.Device||old.Seq!=rev.Seq {
   s.appendHistoryLocked(key,rev,"recibido")
  }
 }
 s.state.Known=known
 s.state.Revisions=revisions
 s.state.PairedFolder=filepath.Clean(folder)
 return s.saveLocked()
}
// syncOrderedFilmRows changes *only* the positions of film rows. It keeps
// unnumbered headings, blank lines, exact film text, existing marks, and films
// added locally. A peer's older order cannot roll back a newer local movement.
func syncOrderedFilmRows(current,remote string) string {
 lines:=strings.Split(current,"\n")
 ranks:=map[string]int{}
 for _,line:=range strings.Split(remote,"\n"){
  f,ok:=parseSyncFilm(line)
  if ok {if _,seen:=ranks[f.key];!seen {ranks[f.key]=len(ranks)}}
 }
 type row struct { line string; original int; rank int }
 indices:=[]int{}
 rows:=[]row{}
 for i,line:=range lines {
  f,ok:=parseSyncFilm(line)
  if !ok{continue}
  rank,found:=ranks[f.key];if !found{rank=len(ranks)+i}
  indices=append(indices,i)
  rows=append(rows,row{line:line,original:i,rank:rank})
 }
 sort.SliceStable(rows,func(i,j int)bool{
  if rows[i].rank==rows[j].rank{return rows[i].original<rows[j].original}
  return rows[i].rank<rows[j].rank
 })
 for i,idx:=range indices{lines[idx]=rows[i].line}
 return strings.Join(lines,"\n")
}

// Reject unrelated libraries which happen to reuse a legacy list ID. A new
// Android installation always starts with "atlas" and can otherwise append
// hundreds of unrelated film lines into the selected Windows list.
func syncListFilmKeys(list filmCollection) map[string]bool {
 keys:=map[string]bool{}
 for _,line:=range strings.Split(list.Text,"\n"){
  if f,ok:=parseSyncFilm(line);ok{keys[f.key]=true}
 }
 return keys
}
func syncListIdentityError(local,remote filmCollection) error {
 ours,theirs:=syncListFilmKeys(local),syncListFilmKeys(remote)
 if len(ours)==0||len(theirs)==0{
  if local.Name!=remote.Name&&len(ours)!=len(theirs){
   return fmt.Errorf("Sincronización detenida: dos listas distintas comparten el ID %s. Revisa las copias antes de continuar",local.ID)
  }
  return nil
 }
 matches:=0
 for key:=range theirs{if ours[key]{matches++}}
 starterA:=strings.HasPrefix(local.Name,"Selector Cine 365 ·")
 starterB:=strings.HasPrefix(remote.Name,"Selector Cine 365 ·")
 differentSets:=len(ours)!=len(theirs)||matches!=len(ours)
 starterCollision:=starterA!=starterB&&differentSets
 differentNamedLists:=local.Name!=remote.Name&&differentSets
 smaller:=len(ours);if len(theirs)<smaller{smaller=len(theirs)}
 unrelated:=matches*2<smaller
 // An already contaminated peer can contain the whole true list plus
 // hundreds of films from other lists. Overlap alone does not prove identity.
 extraLeft,extraRight:=len(ours)-matches,len(theirs)-matches
 excess:=extraLeft;if extraRight>excess{excess=extraRight}
 unexpectedExpansion:=excess>=24&&excess*5>=smaller
 if starterCollision||differentNamedLists||unrelated||unexpectedExpansion{
  return fmt.Errorf("Sincronización detenida: «%s» y «%s» comparten el ID %s pero parecen listas diferentes. No se han mezclado ni enviado datos. Conserva las copias antes de continuar",local.Name,remote.Name,local.ID)
 }
 return nil
}
func syncMerge(local libraryState,revisions map[string]syncRevision,known map[string]int64,peer syncSnapshot) (libraryState,map[string]syncRevision,map[string]int64,[]syncConflict,error) {
 return syncMergeVersioned(local,revisions,known,peer,false)
}
// With an already paired device, ONLY explicit causal revisions may change a mark.
// Raw legacy/snapshot ✓ flags are not operations and must never resurrect old marks.
func syncMergeVersioned(local libraryState,revisions map[string]syncRevision,known map[string]int64,peer syncSnapshot,strict bool) (libraryState,map[string]syncRevision,map[string]int64,[]syncConflict,error) {
 if peer.Protocol!=1&&peer.Protocol!=syncProtocol{return local,nil,nil,nil,errors.New("versión de sincronización desconocida")}
 if err:=validateLibrary(peer.Library);err!=nil{return local,nil,nil,nil,fmt.Errorf("biblioteca remota no válida: %w",err)}
 // Copia profunda: un error de la fuente no modifica las listas en memoria.
 payload,_:=json.Marshal(local)
 var result libraryState
 if err:=json.Unmarshal(payload,&result);err!=nil{return local,nil,nil,nil,err}
 // Los borrados se conservan como tombstones. Una copia vieja sin el campo
 // deleted_lists no puede resucitar una lista que el usuario ya eliminó.
 if result.DeletedLists==nil{result.DeletedLists=map[string]bool{}}
 for id,deleted:=range peer.Library.DeletedLists{if deleted{result.DeletedLists[id]=true}}
 // Solo los IDs de TMDB confirmados; jamás incluir o sincronizar la clave API.
 if result.TMDBOverrides==nil{result.TMDBOverrides=map[string]int{}}
 for key,id:=range peer.Library.TMDBOverrides{
  if _,exists:=result.TMDBOverrides[key];!exists{result.TMDBOverrides[key]=id}
 }
 // La puntuación y el corazón comparten un registro versionado por película.
 // Un cambio posterior prevalece; los registros vacíos no se eliminan para que
 // una copia antigua no resucite favoritos descartados.
 if result.Reactions==nil{result.Reactions=map[string]filmReaction{}}
 for key,remote:=range peer.Library.Reactions{
  localRec,exists:=result.Reactions[key]
  if !exists||remote.ChangedAt>localRec.ChangedAt||
   (remote.ChangedAt==localRec.ChangedAt&&
    (remote.Title+"|"+remote.Year+"|"+strconv.Itoa(remote.Rating)+"|"+strconv.FormatBool(remote.Favorite))>
    (localRec.Title+"|"+localRec.Year+"|"+strconv.Itoa(localRec.Rating)+"|"+strconv.FormatBool(localRec.Favorite))) {
   result.Reactions[key]=remote
  }
 }
 // El orden viaja por separado para que reordenar no cambie corazones ni notas.
 if remote:=peer.Library.FavoritesOrder;remote!=nil{
  current:=result.FavoritesOrder
  if current==nil||remote.ChangedAt>current.ChangedAt||
   (remote.ChangedAt==current.ChangedAt&&strings.Join(remote.Keys,"\x00")>strings.Join(current.Keys,"\x00")){
   result.FavoritesOrder=remote
  }
 }
 remoteMarks:=syncStatuses(peer.Library)
 for _,other:=range peer.Library.Lists {
  if result.DeletedLists[other.ID]{continue}

  idx:=-1
  for j:=range result.Lists{if result.Lists[j].ID==other.ID{idx=j;break}}
  if idx<0 {
   if len(result.Lists)>=40{return local,nil,nil,nil,errors.New("se alcanzó el máximo de listas; no se ha importado nada")}
   result.Lists=append(result.Lists,other);continue
  }
  own:=&result.Lists[idx]
  if err:=syncListIdentityError(*own,other);err!=nil{return local,nil,nil,nil,err}
  if strings.HasPrefix(own.Name,"Selector Cine 365 ·")&&!strings.HasPrefix(other.Name,"Selector Cine 365 ·"){own.Name=other.Name}
  keys:=map[string]bool{}
  for _,line:=range strings.Split(own.Text,"\n"){if f,ok:=parseSyncFilm(line);ok{keys[f.key]=true}}
  for _,line:=range strings.Split(other.Text,"\n"){
   if f,ok:=parseSyncFilm(line);ok&&!keys[f.key]{
    if own.Text!=""&&!strings.HasSuffix(own.Text,"\n"){own.Text+="\n"}
    own.Text+=line;keys[f.key]=true
   }
  }
  if other.OrderChangedAt>own.OrderChangedAt {
   own.Text=syncOrderedFilmRows(own.Text,other.Text)
   own.OrderChangedAt=other.OrderChangedAt
  }
 }
 // Eliminar del resultado las listas marcadas como borradas, aunque
 // ya estuvieran presentes localmente; resolver la selección sin perder listas.
 if len(result.DeletedLists)>0{
  kept:=make([]filmCollection,0,len(result.Lists))
  for _,list:=range result.Lists{if !result.DeletedLists[list.ID]{kept=append(kept,list)}}
  if len(kept)==0{return local,nil,nil,nil,errors.New("No se puede fusionar: todas las listas aparecen borradas")}
  result.Lists=kept
  active:=false
  for _,list:=range kept{if list.ID==result.ActiveID{active=true}}
  if !active{result.ActiveID=kept[0].ID}
  ids:=map[string]bool{}
  for _,list:=range kept{ids[list.ID]=true}
  filtered:=make([]string,0,len(result.ComparisonIDs))
  for _,id:=range result.ComparisonIDs{if ids[id]{filtered=append(filtered,id)}}
  if len(filtered)<2{filtered=nil}
  result.ComparisonIDs=filtered
 }
 if len(result.ComparisonIDs)==0&&len(peer.Library.ComparisonIDs)>0{
  valid:=map[string]bool{};for _,l:=range result.Lists{valid[l.ID]=true}
  keep:=true;for _,id:=range peer.Library.ComparisonIDs{if !valid[id]{keep=false}}
  if keep{result.ComparisonIDs=append([]string(nil),peer.Library.ComparisonIDs...)}
 }
 joined:=make(map[string]syncRevision,len(revisions))
 for k,v:=range revisions{joined[k]=v}
 allKnown:=syncJoin(known,peer.Known)
 conflicts:=[]syncConflict{}
 for key,theirs:=range peer.Revisions{
  ours,exists:=joined[key]
  if !exists{joined[key]=theirs;continue}
  if ours.Device==theirs.Device{if theirs.Seq>ours.Seq{joined[key]=theirs};continue}
  direction:=syncCompare(ours.Clock,theirs.Clock)
  if direction==-1{joined[key]=theirs;continue}
  if direction==1{continue}
  if ours.Checked==theirs.Checked{
   ours.Clock=syncJoin(ours.Clock,theirs.Clock);joined[key]=ours;continue
  }
  conflicts=append(conflicts,syncConflict{Key:key,Local:ours.Checked,Remote:theirs.Checked,LocalDevice:ours.Device,RemoteDevice:theirs.Device})
 }
 allSeen:=syncStatuses(result)
 if strict {
  // Respect the user's current local status for every film; no OR from an
  // unversioned peer library, including films added in old snapshots.
  allSeen=syncStatuses(local)
 } else {
  for key,seen:=range remoteMarks{if seen{allSeen[key]=true}}
 }
 for i:=range result.Lists{
  lines:=strings.Split(result.Lists[i].Text,"\n")
  for j,line:=range lines{
   f,ok:=parseSyncFilm(line);if !ok{continue}
   seen:=allSeen[f.key]
   if rev,ok:=joined[f.key];ok{seen=rev.Checked}
   lines[j]=syncCheckedLine(line,seen)
  }
  result.Lists[i].Text=strings.Join(lines,"\n")
 }
 if err:=validateLibrary(result);err!=nil{return local,nil,nil,nil,fmt.Errorf("fusión cancelada: %w",err)}
 sort.Slice(conflicts,func(i,j int)bool{return conflicts[i].Key<conflicts[j].Key})
 return result,joined,allKnown,conflicts,nil
}
func syncApplyDecisions(l libraryState,decisions map[string]bool) libraryState {
 for i:=range l.Lists{
  lines:=strings.Split(l.Lists[i].Text,"\n")
  for j,line:=range lines{if f,ok:=parseSyncFilm(line);ok{if v,yes:=decisions[f.key];yes{lines[j]=syncCheckedLine(line,v)}}}
  l.Lists[i].Text=strings.Join(lines,"\n")
 }
 return l
}
func cloudSyncRoot(folder string) string {
 if appIsLab { return filepath.Join(folder,cloudDirectoryName) }
 if _,err:=os.Stat(filepath.Join(folder,"actual"));err==nil{return folder}
 if _,err:=os.Stat(filepath.Join(folder,"sincronizacion"));err==nil{return folder}
 return filepath.Join(folder,cloudDirectoryName)
}
func readSyncFile(path string) ([]byte,error) {
 file,err:=os.Open(path);if err!=nil{return nil,err};defer file.Close()
 st,err:=file.Stat();if err!=nil{return nil,err}
 if st.IsDir()||st.Size()>8*1024*1024{return nil,fmt.Errorf("archivo compartido demasiado grande: %s",filepath.Base(path))}
 return os.ReadFile(path)
}
func (a *runningApp) cloudPlanFor(local libraryState,folder string) (syncPlan,error) {
 if a.syncStore==nil{return syncPlan{},errors.New("falta inicializar la sincronización")}
 if folder==""{return syncPlan{},errors.New("selecciona una carpeta compartida")}
 if err:=validateLibrary(local);err!=nil{return syncPlan{},err}
 id,revisions,known:=a.syncStore.parts()
 strict:=a.syncStore.isPaired(folder)
 p:=syncPlan{Library:local,Revisions:revisions,Known:known,Conflicts:[]syncConflict{}}
 hash:=sha256.New()
 localBytes,_:=json.Marshal(local)
 hash.Write(localBytes)
 root:=cloudSyncRoot(folder)
 process:=func(name string,b []byte,legacy bool)error{
  sum:=sha256.Sum256(b)
  fmt.Fprintf(hash,"%s:%x|",name,sum)
  var peer syncSnapshot
  if legacy {
   peer.Protocol=1;peer.DeviceID="windows-legacy"
   if err:=json.Unmarshal(b,&peer.Library);err!=nil{return fmt.Errorf("copia antigua corrupta: %w",err)}
   p.Legacy=true
  }else{
   if err:=json.Unmarshal(b,&peer);err!=nil{return fmt.Errorf("archivo de dispositivo corrupto (%s): %w",name,err)}
   if peer.DeviceID==id{return nil}
  }
  merged,r,k,issues,err:=syncMergeVersioned(p.Library,p.Revisions,p.Known,peer,strict)
  if err!=nil{return err}
  p.Library=merged;p.Revisions=r;p.Known=k;p.Sources++
  existing:=map[string]bool{}
  for _,c:=range p.Conflicts{existing[c.Key]=true}
  for _,c:=range issues{if !existing[c.Key]{p.Conflicts=append(p.Conflicts,c);existing[c.Key]=true}}
  return nil
 }
 legacy:=filepath.Join(root,"actual","colecciones.json")
 if !strict {
  if b,e:=readSyncFile(legacy);e==nil{
   if err:=process("actual/colecciones.json",b,true);err!=nil{return p,err}
  }else if !os.IsNotExist(e){return p,e}
 }
 peersDir:=filepath.Join(root,"sincronizacion","dispositivos")
 entries,err:=os.ReadDir(peersDir)
 if err!=nil&&!os.IsNotExist(err){return p,err}
 if len(entries)>100{return p,errors.New("demasiados archivos de sincronización para validar")}
 for _,it:=range entries{
  if it.IsDir()||!strings.HasSuffix(it.Name(),".json"){continue}
  b,e:=readSyncFile(filepath.Join(peersDir,it.Name()));if e!=nil{return p,e}
  if e=process("sincronizacion/dispositivos/"+it.Name(),b,false);e!=nil{return p,e}
 }
 p.Token=hex.EncodeToString(hash.Sum(nil))
 return p,nil
}
func (a *runningApp) cloudSyncPreview() (map[string]any,error) {
 if a.syncStore==nil{return nil,errors.New("sincronización no inicializada")}
 a.syncPlanMu.Lock();defer a.syncPlanMu.Unlock()
 folder:=a.cloud.status().Folder
 a.mu.Lock()
 local,revision:=a.library,a.libraryRevision
 a.mu.Unlock()
 p,err:=a.cloudPlanFor(local,folder);if err!=nil{return nil,err}
 a.mu.Lock();defer a.mu.Unlock()
 if revision!=a.libraryRevision||folder!=a.cloud.status().Folder{return nil,errors.New("La biblioteca o la carpeta ha cambiado. Vuelve a revisar la sincronización")}
 a.syncPreviewToken=p.Token
 return map[string]any{"token":p.Token,"sources":p.Sources,"legacy":p.Legacy,
  "local_lists":len(a.library.Lists),"local_seen":syncWatched(a.library),
  "result_lists":len(p.Library.Lists),"result_seen":syncWatched(p.Library),
  "conflicts":p.Conflicts,"requires_confirmation":true},nil
}
func (a *runningApp) cloudSyncConfirm(token string,decisions map[string]bool) (map[string]any,error) {
 if a.syncStore==nil{return nil,errors.New("sincronización no inicializada")}
 a.syncPlanMu.Lock();defer a.syncPlanMu.Unlock()
 a.mu.Lock()
 folder:=a.cloud.status().Folder
 if token==""||token!=a.syncPreviewToken{a.mu.Unlock();return nil,errors.New("revisa los datos antes de sincronizar")}
 local,revision:=a.library,a.libraryRevision
 a.mu.Unlock()
 p,err:=a.cloudPlanFor(local,folder)
 if err!=nil{return nil,err}
 a.mu.Lock()
 if p.Token!=token||token!=a.syncPreviewToken||revision!=a.libraryRevision||folder!=a.cloud.status().Folder{a.mu.Unlock();return nil,errors.New("la carpeta o la biblioteca ha cambiado: vuelve a revisar")}
 required:=map[string]bool{}
 for _,c:=range p.Conflicts{required[c.Key]=true}
 for key:=range required{if _,ok:=decisions[key];!ok{a.mu.Unlock();return nil,fmt.Errorf("falta resolver el conflicto: %s",key)}}
 // Copia adicional inmediata, independiente de la copia diaria habitual.
 if err:=a.saveSyncRecovery(a.library,"");err!=nil{a.mu.Unlock();return nil,err}
 p.Library=syncApplyDecisions(p.Library,decisions)
 if err:=saveLibrary(a.libraryPath,a.savePath,p.Library);err!=nil{a.mu.Unlock();return nil,err}
 a.library=p.Library
 a.libraryRevision++
 if err:=a.syncStore.accept(p,decisions,folder);err!=nil{a.mu.Unlock();return nil,err}
 a.syncPreviewToken=""
 a.syncPendingConflict.Store(false)
 a.syncLastCheck.Store(time.Now().UnixMilli())
 a.syncLastReceived.Store(time.Now().UnixMilli())
 a.mu.Unlock()
 saved:=a.cloud.backup()
 if saved.LastError!=""{return nil,errors.New(saved.LastError)}
 if err:=a.publishDeviceSnapshot();err!=nil{return nil,err}
 return map[string]any{"changed":true,"sources":p.Sources,"lists":len(p.Library.Lists),"seen":syncWatched(p.Library)},nil
}
func (a *runningApp) publishDeviceSnapshot() error {
 if a.syncStore==nil||a.cloud==nil{return nil}
 folder:=a.cloud.status().Folder
 if !a.syncStore.isPaired(folder){return nil}
 a.mu.Lock();l:=a.library;a.mu.Unlock()
 data,err:=a.syncStore.snapshot(l);if err!=nil{return err}
 root:=cloudSyncRoot(folder)
 id,_,_:=a.syncStore.parts()
 target:=filepath.Join(root,"sincronizacion","dispositivos",id+".json")
 if err:=cloudWriteFile(target,data);err!=nil{return err}
 readBack,err:=os.ReadFile(target);if err!=nil{return err}
 if sha256.Sum256(readBack)!=sha256.Sum256(data){return errors.New("el proveedor no confirmó el archivo de este dispositivo")}
 return nil
}
