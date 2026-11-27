//go:build windows && lab

package main

import (
 "encoding/json"
 "net/http/httptest"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func labNetworkTestRequest(a *runningApp,method,token string,data []byte) *httptest.ResponseRecorder{
 r:=httptest.NewRequest(method,"/api/lab/network-log",strings.NewReader(string(data)))
 if token!=""{r.Header.Set("X-Atlas-Token",token)}
 w:=httptest.NewRecorder()
 a.serveLabNetworkLog(w,r)
 return w
}
func TestLabNetworkEventValidation(t *testing.T){
 good:=labNetworkEvent{Route:"/api/extra-backup/now",Method:"POST",Kind:"network",Status:0,Attempts:1}
 if !labNetEventValid(good){t.Fatal("un fallo de red legítimo debe registrarse")}
 cases:=[]labNetworkEvent{
  {Route:"/api/extra-backup/now?token=secret",Method:"POST",Kind:"network",Attempts:1},
  {Route:"https://tmdb.org/api",Method:"GET",Kind:"network",Attempts:1},
  {Route:"/api/lab/network-log",Method:"POST",Kind:"network",Attempts:1},
  {Route:"/api/ping",Method:"GET",Kind:"network",Attempts:1},
  {Route:"/api/library",Method:"POST",Kind:"network",Attempts:4},
  {Route:"/api/library",Method:"POST",Kind:"network",Attempts:1,Status:503},
  {Route:"/api/library",Method:"POST",Kind:"http",Attempts:1,Status:200},
  {Route:"/api/library",Method:"EXPLOIT",Kind:"network",Attempts:1},
 }
 for _,event:=range cases {if labNetEventValid(event){t.Errorf("evento peligroso aceptado: %+v",event)}}
 for _,test:=range []struct{path,want string}{
  {"/api/discover/detail","tmdb"},
  {"/api/platforms/cache","tmdb"},
  {"/api/extra-backup/status","backup"},
  {"/api/extra-backup/now","backup"},
  {"/api/library","local"},
 }{if got:=labNetSource(test.path);got!=test.want{t.Errorf("%s: %s, esperado %s",test.path,got,test.want)}}
}
func TestLabNetworkLogIsPrivateAndBounded(t *testing.T){
 root:=t.TempDir()
 local:=t.TempDir()
 t.Setenv("LOCALAPPDATA",local)
 a:=&runningApp{savePath:filepath.Join(root,appDataDirectoryName,"lista.txt"),token:"session-secret"}
 e:=labNetworkEvent{Route:"/api/extra-backup/now",Method:"POST",Kind:"network",Status:0,Attempts:1}
 b,err:=json.Marshal(e);if err!=nil{t.Fatal(err)}
 if got:=labNetworkTestRequest(a,"POST","",b).Code;got!=403{t.Fatalf("token obligatorio: %d",got)}
 if got:=labNetworkTestRequest(a,"GET","session-secret",b).Code;got!=405{t.Fatalf("solo POST: %d",got)}
 if got:=labNetworkTestRequest(a,"POST","session-secret",b).Code;got!=204{t.Fatalf("el evento válido falló: %d",got)}
 logPath:=a.labNetworkLogPath()
 if filepath.Dir(logPath)!=filepath.Join(local,appDataDirectoryName,"cache_lab"){t.Fatalf("archivo fuera de caché local LAB: %s",logPath)}
 data,err:=os.ReadFile(logPath);if err!=nil{t.Fatal(err)}
 if !strings.Contains(string(data),"backup | /api/extra-backup/now | POST | network | HTTP 0 | intentos 1"){t.Fatal("registro de diagnóstico incompleto")}
 if strings.Contains(string(data),"session-secret")||strings.Contains(string(data),"token="){t.Fatal("el diagnóstico reveló credenciales")}
 // En Windows los permisos se expresan mediante ACL, no mediante bits POSIX 0600.
 // La privacidad se verifica por carpeta LAB separada, token de API y ausencia de secretos.
 if info,err:=os.Stat(logPath);err!=nil||info.IsDir(){t.Fatalf("el registro debe ser un archivo normal en LAB: %v, %v",info,err)}
 // No se admiten URLs con queries ni cuerpos adicionales.
 invalid:=[]byte(`{"route":"/api/discover/detail?title=secreto","method":"GET","kind":"network","status":0,"attempts":1}`)
 if got:=labNetworkTestRequest(a,"POST","session-secret",invalid).Code;got!=400{t.Fatalf("se aceptó ruta con parámetros privados: %d",got)}
 body:=append([]byte(`{"route":"/api/library","method":"GET","kind":"network","status":0,"attempts":1,`),[]byte(strings.Repeat("x",2000))...)
 if got:=labNetworkTestRequest(a,"POST","session-secret",body).Code;got!=400{t.Fatalf("debe rechazar cuerpos de más de 1024 bytes: %d",got)}
 if err:=os.WriteFile(logPath,[]byte(strings.Repeat("A",labNetMaxLogBytes+10)),0600);err!=nil{t.Fatal(err)}
 e=labNetworkEvent{Route:"/api/discover/detail",Method:"GET",Kind:"http",Status:502,Attempts:3}
 b,_=json.Marshal(e)
 if got:=labNetworkTestRequest(a,"POST","session-secret",b).Code;got!=204{t.Fatalf("la rotación falló: %d",got)}
 if _,err:=os.Stat(logPath+".anterior");err!=nil{t.Fatal("falta rotación del registro: ",err)}
 newLog,err:=os.ReadFile(logPath);if err!=nil{t.Fatal(err)}
 if !strings.Contains(string(newLog),"tmdb | /api/discover/detail | GET | http | HTTP 502 | intentos 3"){t.Fatal("la rotación perdió el evento nuevo")}
 if len(newLog)>2048{t.Fatal("el archivo de registro nuevo no está acotado")}
}
func TestLabNetworkIsInjectedInFinalHTML(t *testing.T){
 page,err:=makeHTML("☐ Ejemplo (2024)","secret")
 if err!=nil{t.Fatal(err)}
 first:=strings.Index(page,"const labOriginalFetch=window.fetch.bind(window);")
 boot:=strings.LastIndex(page,"void initLibrary()")
 if first<0||boot<0||first>=boot{t.Fatal("el interceptador debe instalarse antes de iniciar las primeras peticiones")}
 for _,s:=range []string{
  "window.fetch=labResilientFetch;",
  "const tries=safeRead?3:1;",
  "labNetworkReport(info,'network',0,attempt);",
  "Conexión interrumpida al guardar o cambiar datos",
  "Guardado no confirmado · comprueba la lista",
  "Guardado: ",
 } {
  if !strings.Contains(page,s){t.Errorf("no se inyectó la protección: %s",s)}
 }
 if strings.Count(page,"window.fetch=labResilientFetch;")!=1{t.Fatal("el interceptador solo debe instalarse una vez")}
}
