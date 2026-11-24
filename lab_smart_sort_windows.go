//go:build windows

package main

import (
 "encoding/json"
 "net/http"
 "strconv"
)

// Metadatos de la caché de Descubrir, solo lectura y sin conexiones externas.
// Permite combinar géneros/duración/puntuaciones sin cargar fichas individuales.
type labSmartFilmMeta struct {
 Runtime int `json:"runtime"`
 Genres []string `json:"genres"`
 VoteAverage float64 `json:"vote_average"`
 VoteCount int `json:"vote_count"`
}
type labSmartCacheReply struct {
 Items map[string]labSmartFilmMeta `json:"items"`
 Resolved map[string]int `json:"resolved"`
}
func labSmartMeta(rec labDetailRecord) labSmartFilmMeta {
 data:=rec.Data
 m:=labSmartFilmMeta{}
 if v,ok:=data["runtime"].(float64);ok&&v>=0&&v<=600{m.Runtime=int(v)}
 if v,ok:=data["runtime"].(int);ok&&v>=0&&v<=600{m.Runtime=v}
 if v,ok:=data["vote_average"].(float64);ok&&v>=0&&v<=10{m.VoteAverage=v}
 if v,ok:=data["vote_count"].(float64);ok&&v>=0&&v<=100000000{m.VoteCount=int(v)}
 if v,ok:=data["vote_count"].(int);ok&&v>=0&&v<=100000000{m.VoteCount=v}
 if arr,ok:=data["genres"].([]string);ok{
  for _,s:=range arr {if len(s)>0&&len(s)<80&&len(m.Genres)<12{m.Genres=append(m.Genres,s)}}
 }else if arr,ok:=data["genres"].([]any);ok{
  for _,v:=range arr{if s,ok:=v.(string);ok&&len(s)>0&&len(s)<80&&len(m.Genres)<12{m.Genres=append(m.Genres,s)}}
 }
 return m
}
func (a *runningApp) serveLabSmartCache(w http.ResponseWriter,r *http.Request) {
 if r.Method!=http.MethodGet {http.Error(w,"Método no permitido",405);return}
 if !a.auth(r){http.Error(w,"No autorizado",403);return}
 c:=a.labCatalog()
 reply:=labSmartCacheReply{Items:make(map[string]labSmartFilmMeta),Resolved:make(map[string]int)}
 c.mu.Lock()
 for key,rec:=range c.disk.Items {
  id,err:=strconv.Atoi(key)
  if err!=nil||id<1||id>999999999||len(rec.Data)==0{continue}
  reply.Items[key]=labSmartMeta(rec)
  if len(reply.Items)>=2500{break}
 }
 for key,id:=range c.disk.Resolved{
  if id>0&&id<=999999999&&len(reply.Resolved)<4000{
   if _,exists:=reply.Items[strconv.Itoa(id)];exists{reply.Resolved[key]=id}
  }
 }
 c.mu.Unlock()
 w.Header().Set("Cache-Control","no-store")
 w.Header().Set("Content-Type","application/json; charset=utf-8")
 _=json.NewEncoder(w).Encode(reply)
}
