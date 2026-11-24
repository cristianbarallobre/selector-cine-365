//go:build windows

package main

import (
 "embed"
 "strings"
)

//go:embed lab_smart_sort_windows.html
var labSmartSortAssets embed.FS

func injectLabSmartSort(page string) string {
 raw,err:=labSmartSortAssets.ReadFile("lab_smart_sort_windows.html")
 if err!=nil{return page}
 parts:=strings.SplitN(string(raw),"<!-- LAB_SMART_HTML -->",2)
 if len(parts)!=2{return page}
 fragments:=strings.SplitN(parts[1],"<!-- LAB_SMART_SCRIPT -->",2)
 if len(fragments)!=2{return page}
 // Solo Windows LAB: la plantilla y los selectores de la versión estable no se modifican.
 required:=[]string{
  "</head>",
  `<div id="sortCount" class="sort-count" aria-live="polite"></div>`,
  `<div class="platform-detail" id="resultPlatforms"></div>`,
  `<div class="pick" aria-live="polite"><img id="resultPoster"`,
  "function updateSortAvailability(){",
  "function updateChoicePlatforms(){",
  "function resetChoice(){chosen=null;",
 }
 for _,marker:=range required {if !strings.Contains(page,marker){return page}}
 // Los tres textos son el valor inicial, el botón visible y el menú reconstruido con JS.
 // El valor interno "all" sigue incluyendo todas las películas pendientes.
 page=strings.ReplaceAll(page,"Todas las películas pendientes","Todas las plataformas")
 page=strings.Replace(page,"</head>",parts[0]+"\n</head>",1)
 page=strings.Replace(page,
  `<div id="sortCount" class="sort-count" aria-live="polite"></div>`,
  fragments[0]+`<div id="sortCount" class="sort-count" aria-live="polite"></div>`,1)
 page=strings.Replace(page,
  `<div class="platform-detail" id="resultPlatforms"></div>`,
  `<div class="platform-detail" id="resultPlatforms"></div><div class="lab-smart-reason" id="labSmartReason" hidden></div>`,1)
 // Botón dentro de la tarjeta, esquina superior derecha, tanto en sorteo normal como inteligente.
 page=strings.Replace(page,
  `<div class="pick" aria-live="polite"><img id="resultPoster"`,
  `<div class="pick" aria-live="polite"><button id="labSmartResultDetails" class="lab-smart-details-btn" type="button" hidden disabled aria-label="Ver detalles de la película sorteada">Detalles</button><img id="resultPoster"`,1)
 page=strings.Replace(page,"function updateSortAvailability(){",
  "function updateSortAvailability(){queueMicrotask(labSmartUpdateCount);",1)
 page=strings.Replace(page,"function updateChoicePlatforms(){",
  "function updateChoicePlatforms(){labSmartSyncDetailsButton();labSmartSyncSeenAction();",1)
 page=strings.Replace(page,"function resetChoice(){chosen=null;",
  "function resetChoice(){chosen=null;labSmartSyncDetailsButton();labSmartSyncSeenAction();$('labSmartReason').hidden=true;",1)
 // Descubrir ya insertó JS en la IIFE, así que se busca su último cierre real.
 closeAt:=strings.LastIndex(page,"})();")
 scriptEnd:=strings.LastIndex(page,"</script>")
 if closeAt<0||scriptEnd<0||closeAt>=scriptEnd{return page}
 page=page[:closeAt]+"\n"+fragments[1]+"\n"+page[closeAt:]
 return page
}
