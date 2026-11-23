//go:build windows

package main

import (
 "bytes"
 "encoding/json"
 "errors"
 "image"
 _ "image/jpeg"
 _ "image/png"
 "io"
 "net/http"
 "os"
 "path/filepath"
 "strconv"
 "sync"
)

const labWallpaperLimit=12<<20

type labAppearanceSettings struct {
 Preset string `json:"preset"`
 Opacity int `json:"opacity"`
 Blur int `json:"blur"`
}
var labAppearanceMutex sync.Mutex

func labAppearanceDefault() labAppearanceSettings {
 return labAppearanceSettings{Preset:"original",Opacity:86,Blur:4}
}
func labAppearanceValid(v labAppearanceSettings) bool {
 switch v.Preset {
 case "original","cinema","space","classic","custom":
 default:return false
 }
 return v.Opacity>=60&&v.Opacity<=100&&v.Blur>=0&&v.Blur<=16
}
func (a *runningApp) labAppearancePath(name string) string {
 return filepath.Join(filepath.Dir(a.savePath),name)
}
func (a *runningApp) labAppearanceLoad() labAppearanceSettings {
 d:=labAppearanceDefault()
 raw,err:=os.ReadFile(a.labAppearancePath("apariencia_lab.json"))
 if err!=nil||len(raw)>4096{return d}
 if json.Unmarshal(raw,&d)!=nil||!labAppearanceValid(d){return labAppearanceDefault()}
 if d.Preset=="custom" {
  if _,err:=os.Stat(a.labAppearancePath("fondo_personalizado_lab.bin"));err!=nil {
   d.Preset="original"
  }
 }
 return d
}
func labAppearanceAtomicWrite(path string,data []byte) error {
 dir:=filepath.Dir(path)
 if err:=os.MkdirAll(dir,0700);err!=nil{return err}
 f,err:=os.CreateTemp(dir,".apariencia-lab-*")
 if err!=nil{return err}
 defer os.Remove(f.Name())
 if err=f.Chmod(0600);err!=nil{_ =f.Close();return err}
 if _,err=f.Write(data);err!=nil{_ =f.Close();return err}
 if err=f.Close();err!=nil{return err}
 return os.Rename(f.Name(),path)
}
func (a *runningApp) labAppearanceSave(s labAppearanceSettings) error {
 if !labAppearanceValid(s){return errors.New("opciones de apariencia incorrectas")}
 if s.Preset=="custom" {
  if _,err:=os.Stat(a.labAppearancePath("fondo_personalizado_lab.bin"));err!=nil {
   return errors.New("selecciona primero una imagen personalizada")
  }
 }
 b,err:=json.MarshalIndent(s,"","  ")
 if err!=nil{return err}
 return labAppearanceAtomicWrite(a.labAppearancePath("apariencia_lab.json"),b)
}
func labImageCheck(data []byte) (string,error) {
 if len(data)==0||len(data)>labWallpaperLimit{return "",errors.New("la imagen supera el límite de 12 MB")}
 mime:=http.DetectContentType(data)
 if mime!="image/jpeg"&&mime!="image/png" {
  return "",errors.New("solo se permiten fotografías JPG o PNG")
 }
 config,format,err:=image.DecodeConfig(bytes.NewReader(data))
 if err!=nil||config.Width<100||config.Height<100||config.Width>8192||config.Height>8192 {
  return "",errors.New("la imagen está dañada o tiene dimensiones no admitidas")
 }
 if format!="jpeg"&&format!="png"{return "",errors.New("formato de imagen no admitido")}
 return mime,nil
}
func (a *runningApp) serveLabAppearance(w http.ResponseWriter,r *http.Request) {
 if !appIsLab{http.NotFound(w,r);return}
 if !a.auth(r){http.Error(w,"No autorizado",http.StatusForbidden);return}
 w.Header().Set("Cache-Control","no-store")
 switch r.URL.Path {
 case "/api/lab/appearance":
  if r.Method==http.MethodGet {
   labAppearanceMutex.Lock();s:=a.labAppearanceLoad();labAppearanceMutex.Unlock()
   w.Header().Set("Content-Type","application/json; charset=utf-8")
   _=json.NewEncoder(w).Encode(s);return
  }
  if r.Method!=http.MethodPut {http.Error(w,"Método no permitido",405);return}
  var s labAppearanceSettings
  decoder:=json.NewDecoder(http.MaxBytesReader(w,r.Body,4096))
  decoder.DisallowUnknownFields()
  if decoder.Decode(&s)!=nil||!labAppearanceValid(s){http.Error(w,"Opciones no válidas",400);return}
  labAppearanceMutex.Lock()
  err:=a.labAppearanceSave(s)
  labAppearanceMutex.Unlock()
  if err!=nil{http.Error(w,err.Error(),400);return}
  w.Header().Set("Content-Type","application/json; charset=utf-8")
  _=json.NewEncoder(w).Encode(s)
 case "/api/lab/appearance/wallpaper":
  if r.Method==http.MethodGet {
   labAppearanceMutex.Lock()
   file,err:=os.Open(a.labAppearancePath("fondo_personalizado_lab.bin"))
   if err!=nil{labAppearanceMutex.Unlock();http.NotFound(w,r);return}
   data,err:=io.ReadAll(io.LimitReader(file,labWallpaperLimit+1))
   _=file.Close();labAppearanceMutex.Unlock()
   if err!=nil{http.Error(w,"No se pudo leer la imagen",500);return}
   mime,err:=labImageCheck(data)
   if err!=nil{http.Error(w,"Imagen no válida",404);return}
   w.Header().Set("Content-Type",mime)
   w.Header().Set("Content-Length",strconv.Itoa(len(data)))
   _,_=w.Write(data);return
  }
  if r.Method!=http.MethodPost{http.Error(w,"Método no permitido",405);return}
  data,err:=io.ReadAll(http.MaxBytesReader(w,r.Body,labWallpaperLimit))
  if err!=nil{http.Error(w,"La imagen supera los 12 MB",413);return}
  if _,err:=labImageCheck(data);err!=nil{http.Error(w,err.Error(),400);return}
  labAppearanceMutex.Lock()
  err=labAppearanceAtomicWrite(a.labAppearancePath("fondo_personalizado_lab.bin"),data)
  if err==nil{
   s:=a.labAppearanceLoad()
   s.Preset="custom"
   err=a.labAppearanceSave(s)
  }
  labAppearanceMutex.Unlock()
  if err!=nil{http.Error(w,"No se pudo guardar la imagen",500);return}
  w.Header().Set("Content-Type","application/json; charset=utf-8")
  _=json.NewEncoder(w).Encode(map[string]string{"status":"ok"})
 default:http.NotFound(w,r)
 }
}
