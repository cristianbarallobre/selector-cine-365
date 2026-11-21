package main

import (
 "errors"
 "net/url"
 "os/exec"
 "runtime"
)

// Open the official Microsoft authorization page in the system browser,
// not in the embedded WebView. User passwords never pass through the app.
func openDesktopMicrosoftBrowser(link string)error{
 address,err:=url.Parse(link)
 if err!=nil||address.Scheme!="https"||address.Hostname()!="login.microsoftonline.com"{
  return errors.New("Dirección de Microsoft no reconocida")
 }
 switch runtime.GOOS{
 case "windows":return exec.Command("rundll32.exe","url.dll,FileProtocolHandler",link).Start()
 case "darwin":return exec.Command("open",link).Start()
 default:return errors.New("Abre la dirección oficial de Microsoft desde el enlace de la aplicación")
 }
}
