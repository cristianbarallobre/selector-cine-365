//go:build windows || darwin

package main

// Diagnóstico local y acotado de las peticiones de la interfaz Windows LAB.
// No registra URL completas, consultas, claves, tokens, títulos ni cuerpos.

import (
 "encoding/json"
 "fmt"
 "net/http"
 "os"
 "path/filepath"
 "regexp"
 "strings"
 "sync"
 "time"
)

const labNetMaxEventBytes=1024
const labNetMaxLogBytes=256*1024
var labNetPathRE=regexp.MustCompile(`^/api/[a-z0-9/_-]{1,72}$`)
var labNetLogMu sync.Mutex

type labNetworkEvent struct {
 Route string `json:"route"`
 Method string `json:"method"`
 Kind string `json:"kind"`
 Status int `json:"status"`
 Attempts int `json:"attempts"`
}

func labNetSource(route string) string {
 switch {
 case strings.HasPrefix(route,"/api/discover/"),
      strings.HasPrefix(route,"/api/platforms/"),
      strings.HasPrefix(route,"/api/identities"):
  return "tmdb"
 case strings.HasPrefix(route,"/api/extra-backup/"):
  return "backup"
 default:
  return "local"
 }
}
func labNetEventValid(e labNetworkEvent) bool {
 if !labNetPathRE.MatchString(e.Route)||e.Route=="/api/lab/network-log"||e.Route=="/api/ping"{return false}
 switch e.Method {case "GET","HEAD","POST","PUT","PATCH","DELETE":default:return false}
 switch e.Kind {case "network","http":default:return false}
 if e.Status<0||e.Status>599||e.Attempts<1||e.Attempts>3{return false}
 return (e.Kind=="network"&&e.Status==0)||(e.Kind=="http"&&e.Status>=400)
}
func (a *runningApp) labNetworkLogPath() string {
 return filepath.Join(appCacheRoot(filepath.Dir(a.savePath)),"diagnostico_red_lab.txt")
}
func (a *runningApp) serveLabNetworkLog(w http.ResponseWriter,r *http.Request) {
 if !a.auth(r){http.Error(w,"No autorizado",http.StatusForbidden);return}
 if r.Method!=http.MethodPost{http.Error(w,"Método no permitido",http.StatusMethodNotAllowed);return}
 w.Header().Set("Cache-Control","no-store")
 var event labNetworkEvent
 dec:=json.NewDecoder(http.MaxBytesReader(w,r.Body,labNetMaxEventBytes))
 dec.DisallowUnknownFields()
 if err:=dec.Decode(&event);err!=nil||!labNetEventValid(event){
  http.Error(w,"Registro no válido",http.StatusBadRequest);return
 }
 path:=a.labNetworkLogPath()
 labNetLogMu.Lock()
 defer labNetLogMu.Unlock()
 if err:=os.MkdirAll(filepath.Dir(path),0700);err!=nil{
  http.Error(w,"No se pudo guardar el diagnóstico",500);return
 }
 if info,err:=os.Stat(path);err==nil&&info.Size()>=labNetMaxLogBytes {
  _=os.Remove(path+".anterior")
  _=os.Rename(path,path+".anterior")
 }
 f,err:=os.OpenFile(path,os.O_CREATE|os.O_WRONLY|os.O_APPEND,0600)
 if err!=nil{http.Error(w,"No se pudo guardar el diagnóstico",500);return}
 _,err=fmt.Fprintf(f,"%s | %s | %s | %s | %s | HTTP %d | intentos %d\n",
  time.Now().UTC().Format(time.RFC3339),labNetSource(event.Route),
  event.Route,event.Method,event.Kind,event.Status,event.Attempts)
 closeErr:=f.Close()
 if err!=nil||closeErr!=nil{http.Error(w,"No se pudo guardar el diagnóstico",500);return}
 w.WriteHeader(http.StatusNoContent)
}
