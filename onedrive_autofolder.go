package main

import (
 "errors"
 "fmt"
 "os"
 "path/filepath"
 "runtime"
 "strings"
)

// discoverOneDriveFolder locates the same user-owned folder used by Android.
// This is the local OneDrive-client compatibility path, not a Graph login.
// Never guess between several accounts: choosing the wrong root could fork sync.
func discoverOneDriveFolder(home, platform string, getenv func(string) string) (string, error) {
 roots:=[]string{}
 add:=func(p string) {
  if strings.TrimSpace(p)==""{return}
  p=filepath.Clean(p)
  st,err:=os.Stat(p);if err!=nil||!st.IsDir(){return}
  for _,old:=range roots{if strings.EqualFold(old,p){return}}
  roots=append(roots,p)
 }
 if platform=="windows" {
  for _,key:=range []string{"OneDriveConsumer","OneDrive","OneDriveCommercial"} {add(getenv(key))}
  add(filepath.Join(home,"OneDrive"))
 } else if platform=="darwin" {
  dirs,err:=os.ReadDir(filepath.Join(home,"Library","CloudStorage"))
  if err==nil {
   for _,dir:=range dirs{
    if dir.IsDir() && strings.HasPrefix(strings.ToLower(dir.Name()),"onedrive-"){
     add(filepath.Join(home,"Library","CloudStorage",dir.Name()))
    }
   }
  }
  add(filepath.Join(home,"OneDrive"))
 } else {return "", errors.New("la detección de OneDrive solo está disponible en Windows y macOS")}
 if len(roots)==0{return "",errors.New("No se ha detectado una carpeta OneDrive local. La conexión directa por Microsoft Graph todavía no está activada en escritorio.")}
 existing:=[]string{}
 for _,root:=range roots {
  target:=filepath.Join(root,oneDriveFolderName)
  if info,err:=os.Stat(target);err==nil&&info.IsDir(){existing=append(existing,target)}
 }
 if len(existing)==1{return existing[0],nil}
 if len(existing)>1{return "",errors.New("Se han encontrado varias bibliotecas Selector_Cine_365 en cuentas OneDrive diferentes; selecciona la correcta manualmente")}
 if len(roots)!=1{return "",fmt.Errorf("Se han encontrado %d ubicaciones OneDrive posibles; selecciona la cuenta correcta manualmente",len(roots))}
 return filepath.Join(roots[0],oneDriveFolderName),nil
}

func discoverCurrentOneDriveFolder() (string,error) {
 home,err:=os.UserHomeDir();if err!=nil{return "",err}
 return discoverOneDriveFolder(home,runtime.GOOS,os.Getenv)
}
