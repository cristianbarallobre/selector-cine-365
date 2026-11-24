//go:build windows && lab

package main

import (
 "bytes"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestLabMaintenanceIntegratedInHTML(t *testing.T){
 page,err:=makeHTML("☐ Película de ejemplo (2025)","lab-maintenance-test")
 if err!=nil{t.Fatal(err)}
 for _,expected:=range []string{
  "Mantenimiento de Windows LAB",
  `id="labMaintenanceDetails"`,`id="labMaintenanceMetrics"`,
  "function labPreviewRecovery(){",
  "async function labPreviewRecovery(){",
  `'/api/sync/recovery/preview?id='+encodeURIComponent(id)`,
  `'/api/lab/maintenance'`,
  `{action:'clear_posters'}`,
  "Versiones estables en GitHub",
 }{
  if !strings.Contains(page,expected){t.Errorf("Falta integración de mantenimiento: %s",expected)}
 }
 if strings.Count(page,"function labPreviewRecovery(){")!=1{t.Fatal("Se duplicó el visor de copias")}
}
func TestLabMaintenanceAuthenticatedDoesNotDeleteUserData(t *testing.T){
 tmp:=t.TempDir()
 a:=&runningApp{token:"solo-local",libraryPath:filepath.Join(tmp,"colecciones.json"),platforms:newAvailabilityServer(tmp,tmp)}
 normal:=httptest.NewRecorder()
 a.serveLabMaintenance(normal,httptest.NewRequest(http.MethodGet,"/api/lab/maintenance",nil))
 if normal.Code!=403{t.Fatalf("El mantenimiento no exige autenticación: %d",normal.Code)}
 request:=httptest.NewRequest(http.MethodGet,"/api/lab/maintenance",nil)
 request.Header.Set("X-Atlas-Token","solo-local")
 w:=httptest.NewRecorder();a.serveLabMaintenance(w,request)
 if w.Code!=200{t.Fatalf("Estado no disponible: %d %s",w.Code,w.Body.String())}
 var status map[string]any
 if json.Unmarshal(w.Body.Bytes(),&status)!=nil{t.Fatal("Estado no es JSON")}
 if status["webview_managed_separately"]!=true{t.Fatal("Debe indicar límite separado de WebView2")}
 bad:=httptest.NewRequest(http.MethodPost,"/api/lab/maintenance",strings.NewReader(`{"action":"delete_library"}`))
 bad.Header.Set("X-Atlas-Token","solo-local")
 w=httptest.NewRecorder();a.serveLabMaintenance(w,bad)
 if w.Code!=400{t.Fatalf("Acción peligrosa aceptada: %d",w.Code)}
 posterDir:=filepath.Join(appCacheRoot(tmp),"carteles")
 if err:=os.MkdirAll(posterDir,0700);err!=nil{t.Fatal(err)}
 picture:=filepath.Join(posterDir,"abcdef.jpg")
 os.WriteFile(picture,[]byte("imagen"),0600)
 user:=filepath.Join(tmp,"colecciones.json")
 os.WriteFile(user,[]byte("biblioteca"),0600)
 request=httptest.NewRequest(http.MethodPost,"/api/lab/maintenance",bytes.NewBufferString(`{"action":"clear_posters"}`))
 request.Header.Set("X-Atlas-Token","solo-local")
 w=httptest.NewRecorder();a.serveLabMaintenance(w,request)
 if w.Code!=200{t.Fatalf("Limpieza fallida: %d %s",w.Code,w.Body.String())}
 if _,err:=os.Stat(picture);!os.IsNotExist(err){t.Fatal("No retiró el cartel")}
 if data,err:=os.ReadFile(user);err!=nil||string(data)!="biblioteca"{t.Fatal("Limpieza tocó la biblioteca")}
}
func TestLabRecoveryPreviewRequiresValidLibrary(t *testing.T){
 dir:=t.TempDir()
 a:=&runningApp{token:"123",libraryPath:filepath.Join(dir,"colecciones.json")}
 folder:=filepath.Join(dir,"copias_rapidas_lab")
 os.MkdirAll(folder,0700)
 name:="20261124T140000.000000000Z.json"
 good,_:=json.Marshal(labTestLibrary("Test"))
 os.WriteFile(filepath.Join(folder,name),good,0600)
 report,err:=a.previewLocalRecovery("quick:"+name)
 if err!=nil{t.Fatal(err)}
 if report["valid"]!=true||report["lists"]!=1{t.Fatalf("No presenta datos correctos: %v",report)}
 if _,err:=a.previewLocalRecovery("quick:../colecciones.json");err==nil{t.Fatal("Se aceptó ruta arbitraria")}
 os.WriteFile(filepath.Join(folder,name),[]byte("{mal}"),0600)
 if _,err:=a.previewLocalRecovery("quick:"+name);err==nil{t.Fatal("Se aceptó una copia corrupta")}
}
