package main

import (
 "fmt"
 "slices"
 "sort"
 "strings"
 "time"
)

const activityLimit=10000
type activityEntry struct {
 At int64 `json:"at"`
 Kind string `json:"kind"`
 Film string `json:"film,omitempty"`
 List string `json:"list,omitempty"`
 Detail string `json:"detail,omitempty"`
}
type activityFilm struct{title string; seen bool}
func activityFilms(text string)(map[string]activityFilm,[]string){
 m:=map[string]activityFilm{};var order []string
 for _,line:=range strings.Split(text,"\n"){
  f,ok:=parseSyncFilm(line);if !ok{continue}
  title:=syncNumber.ReplaceAllString(syncMark.ReplaceAllString(strings.TrimSpace(line),""),"")
  if _,ok=m[f.key];!ok{order=append(order,f.key)}
  m[f.key]=activityFilm{title,f.checked}
 }
 return m,order
}
func activityKeys[T any](m map[string]T)[]string{
 out:=make([]string,0,len(m));for k:=range m{out=append(out,k)};sort.Strings(out);return out
}
func activityDiff(before,after libraryState)[]activityEntry{
 now:=time.Now().UnixMilli()
 var events []activityEntry
 omitted:=0
 add:=func(kind,film,list,detail string){
  if len(events)>=activityLimit-1{omitted++;return}
  limit:=func(s string)string{r:=[]rune(s);if len(r)>300{return string(r[:300])+"…"};return s}
  events=append(events,activityEntry{now,kind,limit(film),limit(list),limit(detail)})
 }
 old:=map[string]filmCollection{};next:=map[string]filmCollection{}
 for _,l:=range before.Lists{old[l.ID]=l}
 for _,l:=range after.Lists{next[l.ID]=l}
 for _,l:=range before.Lists{if _,ok:=next[l.ID];!ok{add("list_delete","",l.Name,"Lista eliminada")}}
 for _,l:=range after.Lists{
  prev,exists:=old[l.ID]
  if !exists{add("list_create","",l.Name,"Lista creada")}else if prev.Name!=l.Name{add("list_rename","",l.Name,prev.Name+" → "+l.Name)}
  if exists&&prev.Text==l.Text{continue}
  a,ao:=activityFilms(prev.Text);b,bo:=activityFilms(l.Text)
  for _,k:=range ao{if _,ok:=b[k];!ok{add("film_remove",a[k].title,l.Name,"Película retirada")}}
  for _,k:=range bo{
   f:=b[k];p,ok:=a[k]
   if !ok{detail:="Pendiente";if f.seen{detail="Vista"};add("film_add",f.title,l.Name,detail);continue}
   if p.seen!=f.seen{kind:="unseen";if f.seen{kind="seen"};add(kind,f.title,l.Name,"")}
   if p.title!=f.title{add("film_edit",f.title,l.Name,p.title+" → "+f.title)}
  }
  var x,y []string
  for _,k:=range ao{if _,ok:=b[k];ok{x=append(x,k)}}
  for _,k:=range bo{if _,ok:=a[k];ok{y=append(y,k)}}
  if !slices.Equal(x,y){add("order","",l.Name,"Orden de películas modificado")}
 }
 var x,y []string
 for _,l:=range before.Lists{if _,ok:=next[l.ID];ok{x=append(x,l.ID)}}
 for _,l:=range after.Lists{if _,ok:=old[l.ID];ok{y=append(y,l.ID)}}
 if !slices.Equal(x,y){add("order","","","Orden de listas modificado")}
 reactions:=map[string]bool{}
 for k:=range before.Reactions{reactions[k]=true};for k:=range after.Reactions{reactions[k]=true}
 for _,k:=range activityKeys(reactions){
  a,b:=before.Reactions[k],after.Reactions[k];title:=b.Title;if title==""{title=a.Title};year:=b.Year;if year==""{year=a.Year};if year!=""{title+=" ("+year+")"}
  if a.Favorite!=b.Favorite{kind:="favorite_remove";if b.Favorite{kind="favorite_add"};add(kind,title,"Favoritos","")}
  if a.Rating!=b.Rating{add("rating",title,"",ratingActivity(a.Rating)+" → "+ratingActivity(b.Rating))}
 }
 var af,bf []string
 if before.FavoritesOrder!=nil{af=before.FavoritesOrder.Keys};if after.FavoritesOrder!=nil{bf=after.FavoritesOrder.Keys}
 if !slices.Equal(af,bf){add("order","","Favoritos","Orden de Favoritos modificado")}
 ids:=map[string]bool{};for k:=range before.TMDBOverrides{ids[k]=true};for k:=range after.TMDBOverrides{ids[k]=true}
 for _,k:=range activityKeys(ids){a,b:=before.TMDBOverrides[k],after.TMDBOverrides[k];if a!=b{add("identity",strings.ReplaceAll(k,"|"," · "),"","TMDB: "+identityActivity(a)+" → "+identityActivity(b))}}
 if !slices.Equal(before.ComparisonIDs,after.ComparisonIDs)&& (len(before.ComparisonIDs)>0||len(after.ComparisonIDs)>0){add("comparison","","","Selección de listas comparadas modificada")}
 if omitted>0{events=append(events,activityEntry{At:now,Kind:"bulk",Detail:fmt.Sprintf("%d modificaciones adicionales en el mismo guardado",omitted)})}
 return events
}
func ratingActivity(n int)string{if n==0{return "Sin valoración"};return fmt.Sprintf("%d/5",n)}
func identityActivity(n int)string{if n==0{return "Automática"};return fmt.Sprint(n)}
func appendLocalActivity(s *syncDeviceState,before,after libraryState){
 if !s.ActivityMigrated{
  for _,e:=range s.History{kind:="unseen";if e.Checked{kind="seen"};s.Activity=append(s.Activity,activityEntry{At:e.At,Kind:kind,Film:strings.ReplaceAll(e.Key,"|"," · "),Detail:"Historial anterior · "+e.Origin})}
  s.ActivityMigrated=true
 }
 s.Activity=append(s.Activity,activityDiff(before,after)...)
 if len(s.Activity)>activityLimit{s.Activity=append([]activityEntry(nil),s.Activity[len(s.Activity)-activityLimit:]...)}
}
func(s *desktopSyncStore) localActivity()[]activityEntry{
 s.mu.Lock();defer s.mu.Unlock()
 state:=syncDeviceState{Activity:append([]activityEntry(nil),s.state.Activity...),ActivityMigrated:s.state.ActivityMigrated,History:s.state.History}
 appendLocalActivity(&state,libraryState{},libraryState{})
 out:=state.Activity
 for i,j:=0,len(out)-1;i<j;i,j=i+1,j-1{out[i],out[j]=out[j],out[i]}
 return out
}

func appendActivityEvent(s *syncDeviceState,e activityEntry){s.Activity=append(s.Activity,e);if len(s.Activity)>activityLimit{s.Activity=append([]activityEntry(nil),s.Activity[len(s.Activity)-activityLimit:]...)}}
