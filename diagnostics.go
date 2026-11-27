package main

import (
 "context"
 "errors"
 "fmt"
 "io"
 "net"
 "net/http"
 "os"
 "path/filepath"
 "regexp"
 "runtime"
 "strings"
 "sync"
 "sync/atomic"
 "time"
)

const diagnosticVersion = "9.2.20"
const diagnosticMaxBytes = 2 * 1024 * 1024
var diagnosticEnabled atomic.Bool
var diagnosticMu sync.Mutex
var diagnosticRole = "motor"
var diagnosticLast = map[string]time.Time{}
var diagnosticRequestID atomic.Uint64

func diagnosticDir() string {
 return filepath.Join(filepath.Dir(getSavePath()), "logs")
}
func diagnosticPath() string {
 return filepath.Join(diagnosticDir(), "diagnostico-"+diagnosticRole+".log")
}

var diagnosticURL = regexp.MustCompile("(?i)https?://[^\\s]+")
var diagnosticWindowsPath = regexp.MustCompile("(?i)[a-z]:\\\\[^\\r\\n]+")
var diagnosticUnixPath = regexp.MustCompile("/(?:Users|home|var|private|tmp)/[^\\r\\n]+")
var diagnosticSecret = regexp.MustCompile("(?i)(?:bearer\\s+|(?:api_key|access_token|refresh_token|client_secret|authorization|code|token)\\s*[:=]\\s*)[^\\s,;]+")
var diagnosticHTTPCode = regexp.MustCompile(`(?i)(?:HTTP |OneDrive respondió |Microsoft temporalmente no disponible \()([1-5][0-9]{2})`)
var diagnosticEmail = regexp.MustCompile("[A-Za-z0-9._%+\\-]+@[A-Za-z0-9.\\-]+\\.[A-Za-z]{2,}")
func sanitizeDiagnostic(s string) string {
 s=diagnosticURL.ReplaceAllString(s,"[URL omitida]")
 s=diagnosticWindowsPath.ReplaceAllString(s,"[ruta omitida]")
 s=diagnosticUnixPath.ReplaceAllString(s,"[ruta omitida]")
 s=diagnosticSecret.ReplaceAllString(s,"[credencial omitida]")
 s=diagnosticEmail.ReplaceAllString(s,"[correo omitido]")
 s=strings.Map(func(r rune) rune {if r<32 || r==127{return ' '};return r},s)
 if len(s)>1500{s=s[:1500]}
 return s
}

func rotateDiagnostic(path string) {
 if st,err:=os.Stat(path);err==nil && st.Size()>=diagnosticMaxBytes {
  _=os.Remove(path+".2")
  _=os.Rename(path+".1",path+".2")
  _=os.Rename(path,path+".1")
 }
}
func diagnosticEvent(level, component, message string) {
 if !diagnosticEnabled.Load(){return}
 diagnosticMu.Lock()
 defer diagnosticMu.Unlock()
 path:=diagnosticPath()
 if err:=os.MkdirAll(filepath.Dir(path),0700);err!=nil{return}
 rotateDiagnostic(path)
 f,err:=os.OpenFile(path,os.O_CREATE|os.O_WRONLY|os.O_APPEND,0600)
 if err!=nil{return}
 defer f.Close()
 _,_=fmt.Fprintf(f,"%s version=%s pid=%d %s %s %s\n",time.Now().Format(time.RFC3339Nano),diagnosticVersion,os.Getpid(),level,component,sanitizeDiagnostic(message))
}
func logDiagnostic(message string) { diagnosticEvent("INFO","nativo",message) }

// Log a safe category rather than an arbitrary error string, which can include
// request URLs, API keys, account information or movie titles.
func diagnosticErrorKind(err error) string {
 if err==nil{return "none"}
 if errors.Is(err,context.DeadlineExceeded){return "timeout"}
 if errors.Is(err,context.Canceled){return "cancelled"}
 if errors.Is(err,os.ErrPermission){return "permission_denied"}
 if errors.Is(err,os.ErrNotExist){return "file_not_found"}
 var n net.Error
 if errors.As(err,&n) && n.Timeout(){return "timeout"}
 s:=strings.ToLower(err.Error())
 if match:=diagnosticHTTPCode.FindStringSubmatch(s);len(match)==2{return "http_"+match[1]}
 switch {
 case strings.Contains(s,"no space") || strings.Contains(s,"disk full") || strings.Contains(s,"not enough space"):return "disk_full"
 case strings.Contains(s,"access is denied") || strings.Contains(s,"permission denied"):return "permission_denied"
 case strings.Contains(s,"sharing violation") || strings.Contains(s,"being used by another process"):return "file_in_use"
 case strings.Contains(s,"connection refused"):return "connection_refused"
 case strings.Contains(s,"no such host") || strings.Contains(s,"dns"):return "dns"
 case strings.Contains(s,"certificate") || strings.Contains(s,"tls"):return "tls"
 case strings.Contains(s,"429") || strings.Contains(s,"límite temporal") || strings.Contains(s,"limitado las solicitudes"):return "rate_limit"
 case strings.Contains(s,"401") || strings.Contains(s,"403") || strings.Contains(s,"sin autorización") || strings.Contains(s,"sesión microsoft necesita"):return "authorization"
 case strings.Contains(s,"412") || strings.Contains(s,"conflict") || strings.Contains(s,"biblioteca cambió"):return "revision_conflict"
 }
 return fmt.Sprintf("%T",err)
}
func diagnosticError(component, operation string, err error) {
 if err!=nil{diagnosticEvent("ERROR",component,operation+" cause="+diagnosticErrorKind(err))}
}
func diagnosticStart() {
 if len(os.Args)>1 && os.Args[1]=="--webview2-host"{diagnosticRole="ventana"}
 diagnosticEnabled.Store(true)
 diagnosticEvent("INFO","sesion",fmt.Sprintf("inicio os=%s arch=%s go=%s lab=%t",runtime.GOOS,runtime.GOARCH,runtime.Version(),appIsLab))
}
func diagnosticPanic() {
 if value:=recover();value!=nil{
  diagnosticEvent("ERROR","sesion",fmt.Sprintf("panic tipo=%T",value))
  pcs:=make([]uintptr,32)
  frames:=runtime.CallersFrames(pcs[:runtime.Callers(3,pcs)])
  for { frame,more:=frames.Next(); diagnosticEvent("ERROR","stack",fmt.Sprintf("%s:%d",frame.Function,frame.Line));if !more{break} }
  panic(value)
 }
}

// Avoid logging user-controlled paths and query strings.
var diagnosticRoutes = map[string]bool{
 "/":true,"/api/library":true,"/api/library/draft":true,"/api/save":true,"/api/beacon-save":true,"/api/quit":true,
 "/api/ping":true,"/api/sync/history":true,"/api/sync/revisions":true,"/api/sync/recovery/list":true,
 "/api/sync/recovery/pending":true,"/api/sync/recovery/preview":true,"/api/sync/recovery/restore":true,
 "/api/cloud/status":true,"/api/cloud/pending":true,"/api/cloud/config":true,"/api/cloud/backup":true,
 "/api/cloud/browse":true,"/api/cloud/preview":true,"/api/cloud/sync":true,"/api/cloud/auto":true,
 "/api/cloud/auto-folder":true,"/api/cloud/mode":true,"/api/onedrive-desktop/status":true,
 "/api/onedrive-desktop/start":true,"/api/onedrive-desktop/disconnect":true,"/api/lab/maintenance":true,
 "/api/platforms/status":true,"/api/platforms/cache":true,"/api/platforms/config":true,
 "/api/platforms/lookup":true,"/api/platforms/match":true,"/api/platforms/overrides":true,
 "/api/platforms/overrides/import":true,"/api/identities":true,"/api/identities/resolve":true,
 "/api/discover/list":true,"/api/discover/search":true,"/api/discover/detail":true,"/api/discover/person":true,
 "/api/discover/resolve":true,"/api/discover/library-cache":true,"/api/discover/poster":true,"/api/poster":true,
}
type diagnosticResponse struct{http.ResponseWriter; status int; bytes int}
func (w *diagnosticResponse) Unwrap() http.ResponseWriter{return w.ResponseWriter}
func (w *diagnosticResponse) WriteHeader(code int){if w.status!=0{return};w.status=code;w.ResponseWriter.WriteHeader(code)}
func (w *diagnosticResponse) Write(p []byte)(int,error){if w.status==0{w.WriteHeader(200)};n,e:=w.ResponseWriter.Write(p);w.bytes+=n;return n,e}
func diagnosticHeader(h http.Header,name string)string{
 switch v:=h.Get(name);v{case "saved","pending","ok","error","review":return v}
 return "-"
}
func diagnosticHTTP(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if !diagnosticEnabled.Load(){next.ServeHTTP(w,r);return}
  started:=time.Now(); rw:=&diagnosticResponse{ResponseWriter:w}
  route:=r.URL.Path;if !diagnosticRoutes[route]{route="other"}
  method:=r.Method;if method!="GET"&&method!="POST"{method="other"}
  important:=method=="POST" && (route=="/api/library"||strings.HasPrefix(route,"/api/cloud/")||strings.HasPrefix(route,"/api/sync/recovery/")||route=="/api/quit"||route=="/api/save")
  id:=diagnosticRequestID.Add(1)
  if important {diagnosticEvent("INFO","http",fmt.Sprintf("inicio id=%d %s %s",id,method,route))}
  defer diagnosticPanic()
  defer func(){
   status:=rw.status;if status==0{status=200}
   elapsed:=time.Since(started)
   if status>=400||important||elapsed>3*time.Second {
    level:="INFO";if status>=400{level="WARN"};if status>=500{level="ERROR"}
    diagnosticEvent(level,"http",fmt.Sprintf("fin id=%d %s %s status=%d ms=%d bytes=%d recovery=%s txt=%s cloud=%s",id,method,route,status,elapsed.Milliseconds(),rw.bytes,diagnosticHeader(rw.Header(),"X-Cine-Recovery"),diagnosticHeader(rw.Header(),"X-Cine-Legacy"),diagnosticHeader(rw.Header(),"X-Cine-Cloud")))
   }
  }()
  next.ServeHTTP(rw,r)
 })
}

var diagnosticUIEvents=map[string]bool{
 "ui_started":true,"local_connection_lost":true,"local_connection_restored":true,"js_error":true,
 "unhandled_rejection":true,"save_failed":true,"save_confirmed":true,"retry_save":true,
 "emergency_export":true,"close_failed":true,"close_ready":true,"sync_started":true,"sync_failed":true,
 "sync_finished":true,"cloud_check_failed":true,"close_forced":true,"close_cancelled":true,"renderer_terminated":true,
}
var diagnosticJSLocation=regexp.MustCompile("^js_error:[0-9]{1,6}:[0-9]{1,6}$")
func validDiagnosticUIEvent(event string) bool {return diagnosticUIEvents[event]||diagnosticJSLocation.MatchString(event)}
func diagnosticUIEvent(event string) {
 if !validDiagnosticUIEvent(event){return}
 key:=event
 if strings.HasPrefix(event,"js_error:"){key="js_error"}
 // Suppress repeated identical notifications in a short interval.
 diagnosticMu.Lock()
 now:=time.Now();last:=diagnosticLast[key]
 if now.Sub(last)<time.Second{diagnosticMu.Unlock();return}
 diagnosticLast[key]=now
 diagnosticMu.Unlock()
 level:="INFO";if strings.HasSuffix(key,"failed")||key=="local_connection_lost"||key=="js_error"||key=="unhandled_rejection"||key=="close_forced"||key=="renderer_terminated"{level="WARN"}
 diagnosticEvent(level,"interfaz",event)
}

func (a *runningApp) serveDiagnostics(w http.ResponseWriter,r *http.Request){
 if !desktopFeatures || !a.auth(r){http.Error(w,"No autorizado",403);return}
 w.Header().Set("Cache-Control","no-store")
 if r.Method==http.MethodPost{
  body,err:=io.ReadAll(http.MaxBytesReader(w,r.Body,128));if err!=nil||!validDiagnosticUIEvent(string(body)){http.Error(w,"Evento no válido",400);return}
  diagnosticUIEvent(string(body));w.WriteHeader(204);return
 }
 if r.Method!=http.MethodGet{http.Error(w,"Método no permitido",405);return}
 w.Header().Set("Content-Type","text/plain; charset=utf-8")
 w.Header().Set("Content-Disposition","attachment; filename=Selector_Cine_365_diagnostico.log")
 diagnosticMu.Lock();defer diagnosticMu.Unlock()
 _,_=fmt.Fprintln(w,"Selector Cine 365 - diagnóstico local; fechas con zona horaria. Sin contenidos de biblioteca ni credenciales.")
 for _,role:=range []string{"motor","ventana"}{
  for _,suffix:=range []string{".2",".1",""}{
   name:="diagnostico-"+role+".log"+suffix
   f,err:=os.Open(filepath.Join(diagnosticDir(),name));if err!=nil{continue}
   _,_=fmt.Fprintln(w,"\n--- "+name+" ---")
   _,_=io.Copy(w,io.LimitReader(f,diagnosticMaxBytes+4096))
   _=f.Close()
  }
 }
}
