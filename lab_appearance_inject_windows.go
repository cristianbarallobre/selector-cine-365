//go:build windows

package main

import (
 "embed"
 "strings"
)

//go:embed lab_appearance_windows.html
var labAppearanceFiles embed.FS

func injectLabAppearance(page string) string {
 if !appIsLab{return page}
 raw,err:=labAppearanceFiles.ReadFile("lab_appearance_windows.html")
 if err!=nil{return page}
 parts:=strings.SplitN(string(raw),"<!-- LAB_APPEARANCE_SCRIPT -->",2)
 if len(parts)!=2{return page}
 nav:=`<button class="nav-btn" type="button" id="labAppearanceNav">🎨 Apariencia <span class="lab-nav-tag">LAB</span></button>`
 constMarker:=`<button class="nav-btn" type="button" data-view="instructions">❔ Instrucciones</button>`
 page=strings.Replace(page,constMarker,constMarker+"\n    "+nav,1)
 // Incluye navegación programática (p. ej. al abrir una lista), no solo pulsaciones en las pestañas.
 page=strings.Replace(page,"function showView(name){labLeave();","function showView(name){labLeave();labAppearanceLeave();",1)
 page=strings.Replace(page,"</main>",parts[0]+"\n</main>",1)
 page=strings.Replace(page,"\n})();\n</script>","\n"+parts[1]+"\n})();\n</script>",1)
 return page
}
