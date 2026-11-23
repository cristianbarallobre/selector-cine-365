//go:build windows && lab

package main

import (
 "bytes"
 "encoding/json"
 "image"
 "image/png"
 "net/http"
 "net/http/httptest"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func labAppearanceRequest(t *testing.T,a *runningApp,method,path string,data []byte,authorized bool) *httptest.ResponseRecorder {
 t.Helper()
 req:=httptest.NewRequest(method,path,bytes.NewReader(data))
 if authorized {req.Header.Set("X-Atlas-Token","test-token")}
 recorder:=httptest.NewRecorder()
 a.serveLabAppearance(recorder,req)
 return recorder
}

func TestLabAppearanceDefaultsAndValidation(t *testing.T) {
 defaultTheme:=labAppearanceDefault()
 if defaultTheme.Preset!="original"||!labAppearanceValid(defaultTheme){t.Fatal("la apariencia inicial debe ser la original")}
 for _,preset:=range []string{"cinema","space","classic","custom"} {
  s:=defaultTheme;s.Preset=preset
  if !labAppearanceValid(s){t.Fatalf("fondo predeterminado no válido: %s",preset)}
 }
 for _,s:=range []labAppearanceSettings{
  {Preset:"http://somewhere.example",Opacity:86,Blur:4},
  {Preset:"cinema",Opacity:50,Blur:4},
  {Preset:"space",Opacity:101,Blur:4},
  {Preset:"classic",Opacity:86,Blur:40},
 } {if labAppearanceValid(s){t.Fatalf("opciones incorrectas aceptadas: %+v",s)}}
}

func TestLabAppearanceOnlyLocalAndAuthenticated(t *testing.T) {
 root:=t.TempDir()
 path:=filepath.Join(root,"Selector_Cine_365_LAB","mis_365_peliculas.txt")
 a:=&runningApp{savePath:path,token:"test-token"}
 if response:=labAppearanceRequest(t,a,http.MethodGet,"/api/lab/appearance",nil,false);response.Code!=403 {
  t.Fatalf("la configuración debe requerir token: %d",response.Code)
 }
 saved:=labAppearanceSettings{Preset:"cinema",Opacity:80,Blur:7}
 data,_:=json.Marshal(saved)
 put:=labAppearanceRequest(t,a,http.MethodPut,"/api/lab/appearance",data,true)
 if put.Code!=200{t.Fatalf("no se pudo guardar una configuración válida: %d %s",put.Code,put.Body.String())}
 get:=labAppearanceRequest(t,a,http.MethodGet,"/api/lab/appearance",nil,true)
 if get.Code!=200{t.Fatalf("no se pudo leer la configuración: %d",get.Code)}
 var loaded labAppearanceSettings
 if err:=json.Unmarshal(get.Body.Bytes(),&loaded);err!=nil{t.Fatal(err)}
 if loaded!=saved{t.Fatalf("la configuración no sobrevivió al guardado: %+v",loaded)}
 if _,err:=os.Stat(filepath.Join(root,"Selector_Cine_365_LAB","apariencia_lab.json"));err!=nil{t.Fatal("falta archivo de apariencia LAB:",err)}
 if _,err:=os.Stat(filepath.Join(root,"Selector_Cine_365","apariencia_lab.json"));err==nil{t.Fatal("se alteraron carpetas estables")}
 if _,err:=os.Stat(filepath.Join(root,"Selector_Cine_365_LAB","colecciones.json"));err==nil{t.Fatal("no debe tocar las listas")}
 if result:=labAppearanceRequest(t,a,http.MethodPut,"/api/lab/appearance",[]byte(`{"preset":"bad","opacity":80,"blur":5}`),true);result.Code!=400{
  t.Fatalf("debe rechazar fondos desconocidos: %d",result.Code)
 }
}

func TestLabAppearanceUploadAndRestore(t *testing.T) {
 a:=&runningApp{savePath:filepath.Join(t.TempDir(),"Selector_Cine_365_LAB","mis_365_peliculas.txt"),token:"test-token"}
 img:=image.NewRGBA(image.Rect(0,0,150,150))
 var buf bytes.Buffer
 if err:=png.Encode(&buf,img);err!=nil{t.Fatal(err)}
 upload:=labAppearanceRequest(t,a,http.MethodPost,"/api/lab/appearance/wallpaper",buf.Bytes(),true)
 if upload.Code!=200{t.Fatalf("no se pudo guardar PNG: %d, %s",upload.Code,upload.Body.String())}
 get:=labAppearanceRequest(t,a,http.MethodGet,"/api/lab/appearance/wallpaper",nil,true)
 if get.Code!=200||!strings.HasPrefix(get.Header().Get("Content-Type"),"image/png"){
  t.Fatalf("fondo no recuperable: %d, %s",get.Code,get.Header().Get("Content-Type"))
 }
 if a.labAppearanceLoad().Preset!="custom"{t.Fatal("al subir una imagen debería activarse su fondo")}
 invalid:=labAppearanceRequest(t,a,http.MethodPost,"/api/lab/appearance/wallpaper",[]byte("<html>no es una imagen</html>"),true)
 if invalid.Code!=400{t.Fatalf("debe rechazar imágenes no válidas: %d",invalid.Code)}
 after:=a.labAppearanceLoad()
 if after.Preset!="custom"{t.Fatal("una carga incorrecta no debe reemplazar el fondo actual")}
 if err:=a.labAppearanceSave(labAppearanceDefault());err!=nil{t.Fatal(err)}
 if a.labAppearanceLoad().Preset!="original"{t.Fatal("Restablecer debe devolver el fondo original")}
 if len(get.Body.Bytes())!=len(buf.Bytes()){t.Fatal("los bytes de imagen recuperados son diferentes")}
}

func TestLabAppearanceRenderedOnlyInLab(t *testing.T) {
 page,err:=makeHTML("☐ El padrino (1972)","test-token")
 if err!=nil{t.Fatal(err)}
 for _,fragment:=range []string{
  "labAppearanceNav","labAppearanceView","labAppearanceOpacity",
  "labAppearanceBlur","labAppearanceChoose","labAppearanceReset",
  "body.lab-appearance-on::before","labAppearanceLeave();",
 } {
  if !strings.Contains(page,fragment){t.Errorf("falta elemento Apariencia: %s",fragment)}
 }
 if !strings.Contains(page,"labDiscoverView"){t.Fatal("Apariencia no debe eliminar Descubrir")}
 if strings.Count(page,"id=\"labAppearanceNav\"")!=1{t.Fatal("debe existir una sola pestaña Apariencia")}
}
