package main

import (
 "errors"
 "net/http"
 "net/http/httptest"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestDiagnosticPrivacy(t *testing.T) {
 raw:="https://api.example.test/?api_key=SECRET_VALUE user@example.test Bearer ANOTHER_SECRET C:\\Users\\Private\\file.txt"
 safe:=sanitizeDiagnostic(raw)
 for _,secret:=range []string{"SECRET_VALUE","user@example.test","ANOTHER_SECRET","Private"}{
  if strings.Contains(safe,secret){t.Fatalf("private value survived: %s",secret)}
 }
 if validDiagnosticUIEvent("save_failed:private title"){t.Fatal("arbitrary UI messages must be rejected")}
 if !validDiagnosticUIEvent("js_error:123:4"){t.Fatal("safe numeric location rejected")}
 if got:=diagnosticErrorKind(errors.New("OneDrive respondió 503: confidential"));got!="http_503"{t.Fatal(got)}
}

func TestDiagnosticHTTPAndRotation(t *testing.T){
 dir:=t.TempDir()
 t.Setenv("APPDATA",dir)
 t.Setenv("HOME",dir)
 t.Setenv("XDG_CONFIG_HOME",dir)
 old:=diagnosticEnabled.Load()
 diagnosticEnabled.Store(true)
 t.Cleanup(func(){diagnosticEnabled.Store(old)})
 if !strings.HasPrefix(diagnosticDir(),dir){t.Fatalf("not isolated: %s",diagnosticDir())}
 h:=diagnosticHTTP(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("X-Cine-Recovery","saved")
  w.WriteHeader(http.StatusPreconditionFailed)
  _,_=w.Write([]byte("private library content"))
 }))
 req:=httptest.NewRequest(http.MethodPost,"/api/library?token=PRIVATE_TOKEN",strings.NewReader("private payload"))
 response:=httptest.NewRecorder()
 h.ServeHTTP(response,req)
 if response.Code!=412{t.Fatal(response.Code)}
 diagnosticUIEvent("local_connection_lost")
 diagnosticUIEvent("bad PRIVATE_EVENT")
 body,err:=os.ReadFile(diagnosticPath());if err!=nil{t.Fatal(err)}
 log:=string(body)
 for _,want:=range []string{"inicio id=","fin id=","POST /api/library","status=412","recovery=saved","local_connection_lost","version=9.2.19"}{
  if !strings.Contains(log,want){t.Fatalf("missing %s",want)}
 }
 for _,secret:=range []string{"PRIVATE_TOKEN","private payload","private library content","PRIVATE_EVENT"}{
  if strings.Contains(log,secret){t.Fatalf("logged private data: %s",secret)}
 }
 if err:=os.WriteFile(diagnosticPath(),make([]byte,diagnosticMaxBytes),0600);err!=nil{t.Fatal(err)}
 diagnosticEvent("INFO","test","rotation marker")
 if _,err:=os.Stat(diagnosticPath()+".1");err!=nil{t.Fatal(err)}
 if data,err:=os.ReadFile(diagnosticPath());err!=nil||!strings.Contains(string(data),"rotation marker"){t.Fatal("rotation failed")}
 if !desktopFeatures{return}
 a:=&runningApp{token:"test-token"}
 denied:=httptest.NewRecorder()
 a.serveDiagnostics(denied,httptest.NewRequest("GET","/api/diagnostics",nil))
 if denied.Code!=403{t.Fatalf("export should be protected: %d",denied.Code)}
 export:=httptest.NewRequest("GET","/api/diagnostics",nil)
 export.Header.Set("X-Atlas-Token","test-token")
 ok:=httptest.NewRecorder()
 a.serveDiagnostics(ok,export)
 if ok.Code!=200||!strings.Contains(ok.Body.String(),"rotation marker"){t.Fatal("export failed")}
 if _,err:=os.Stat(filepath.Join(diagnosticDir(),"diagnostico-"+diagnosticRole+".log"));err!=nil{t.Fatal(err)}
}
