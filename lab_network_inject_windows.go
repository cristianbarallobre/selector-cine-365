//go:build windows

package main

import (
 "embed"
 "strings"
)

//go:embed lab_network_windows.js
var labNetworkFrontend embed.FS

func injectLabNetwork(page string) string {
 if !appIsLab{return page}
 raw,err:=labNetworkFrontend.ReadFile("lab_network_windows.js")
 if err!=nil{return page}
 // El interceptor se instala al principio de la IIFE: así protege también
 // las primeras consultas realizadas al abrir la aplicación.
 anchor:=`const $=id=>document.getElementById(id);`
 if !strings.Contains(page,anchor){return page}
 page=strings.Replace(page,anchor,anchor+"\n"+string(raw),1)

 // Sin repetir escrituras, el estado del guardado es incierto si se pierde
 // la respuesta. Evita transmitir que no se guardó con certeza.
 page=strings.Replace(page,
  `$('saveState').textContent='⚠ Error al guardar'`,
  `$('saveState').textContent='⚠ Guardado no confirmado · comprueba la lista antes de cerrar'`,1)

 // Con HTTP 5xx OneDrive puede responder texto en vez de JSON. Mostrar
 // el error del servidor, no un SyntaxError/Failed to fetch genérico.
 page=strings.Replace(page,
  `const r=await fetch(path,opts);const data=await r.json();
 if(!r.ok)throw new Error(data.error||'No se pudo completar la operación');return data;`,
  `const r=await fetch(path,opts);const raw=await r.text();
 let data={};
 try{data=raw?JSON.parse(raw):{}}catch(e){data={error:raw.slice(0,160)||'Respuesta no válida del servidor local'};}
 if(!r.ok)throw new Error('OneDrive: '+(data.error||'HTTP '+r.status));return data;`,1)
 return page
}
