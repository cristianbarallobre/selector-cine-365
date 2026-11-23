//go:build windows

package main

import (
 "embed"
 "strings"
)

//go:embed lab_smart_sort_windows.html
var labSmartSortAssets embed.FS

func injectLabSmartSort(page string) string {
 if !appIsLab{return page}
 raw,err:=labSmartSortAssets.ReadFile("lab_smart_sort_windows.html")
 if err!=nil{return page}
 parts:=strings.SplitN(string(raw),"<!-- LAB_SMART_HTML -->",2)
 if len(parts)!=2{return page}
 fragments:=strings.SplitN(parts[1],"<!-- LAB_SMART_SCRIPT -->",2)
 if len(fragments)!=2{return page}
 // No se modifica plantilla.html: todos los cambios de UI son exclusivos de Windows LAB.
 var required=[]string{
  "</head>",
  `<div id="sortCount" class="sort-count" aria-live="polite"></div>`,
  `<div class="platform-detail" id="resultPlatforms"></div>`,
  "function updateSortAvailability(){",
  "function resetChoice(){chosen=null;",
  "\n})();\n</script>",
 }
 for _,marker:=range required {if !strings.Contains(page,marker){return page}}
 page=strings.Replace(page,"</head>",parts[0]+"\n</head>",1)
 page=strings.Replace(page,
  `<div id="sortCount" class="sort-count" aria-live="polite"></div>`,
  fragments[0]+`<div id="sortCount" class="sort-count" aria-live="polite"></div>`,1)
 page=strings.Replace(page,
  `<div class="platform-detail" id="resultPlatforms"></div>`,
  `<div class="platform-detail" id="resultPlatforms"></div><div class="lab-smart-reason" id="labSmartReason" hidden></div>`,1)
 page=strings.Replace(page,"function updateSortAvailability(){",
  "function updateSortAvailability(){queueMicrotask(labSmartUpdateCount);",1)
 page=strings.Replace(page,"function resetChoice(){chosen=null;",
  "function resetChoice(){chosen=null;$('labSmartReason').hidden=true;",1)
 page=strings.Replace(page,"\n})();\n</script>","\n"+fragments[1]+"\n})();\n</script>",1)
 return page
}
