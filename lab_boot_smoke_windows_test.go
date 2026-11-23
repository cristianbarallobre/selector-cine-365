//go:build windows && lab

package main

import (
 "context"
 "fmt"
 "net/http"
 "net/http/httptest"
 "os"
 "os/exec"
 "path/filepath"
 "strings"
 "testing"
 "time"
)

// La regresión LAB-02.10: antes se llamaba update() (que invoca renderList)
// cuando todavía no se había inicializado const labFilterDecade.
// El HTML compilado debe instalar TODOS los módulos antes del arranque.
func TestLabDiscoveryScriptInjection(t *testing.T){
 raw,err:=resources.ReadFile("plantilla.html");if err!=nil{t.Fatal(err)}
 labRaw,err:=labDiscoveryFiles.ReadFile("lab_discover_windows.html");if err!=nil{t.Fatal(err)}
 parts:=strings.SplitN(string(labRaw),labDiscoverySeparator,2)
 if len(parts)!=2{t.Fatal("el separador de Descubrir ha desaparecido")}
 prior:=string(raw)
 after:=injectLabDiscovery(prior)
 t.Logf("HTML original: cierre IIFE=%d; fragmento Discover=%d bytes; JS inyectado=%t; cierre resultante=%d",strings.Count(prior,"\n})();\n</script>"),len(parts[1]),strings.Contains(after,"function labLeave(){"),strings.Count(after,"\n})();\n</script>"))
 if !strings.Contains(after,"function labLeave(){"){t.Fatal("Descubrir no ha insertado labLeave en el script original")}
 if !strings.Contains(after,"function labCreateDetailButton(f)"){t.Fatal("Descubrir no ha insertado el botón Detalles en JS")}
}

func TestLabStartupAfterAllModules(t *testing.T){
 page,err:=makeHTML("☐ Ejemplo (2024)\n☑ Otra (1988)","smoke-token")
 if err!=nil{t.Fatal(err)}
 bootPrefix:="updateListPlatformOptions();update();renderCollectionControls();updatePlatformProgress();void initLibrary().then("
 if strings.Count(page,bootPrefix)!=1{t.Fatalf("el arranque debe existir exactamente una vez; hay %d",strings.Count(page,bootPrefix))}
 bootAt:=strings.Index(page,bootPrefix)
 required:=[]string{
  "const labSmartEnable=$('labSmartEnabled');",
  "const labFilterDecade=$('labFilterDecade')",
  "document.body.appendChild(labFiltersOverlay);",
  "window.fetch=labResilientFetch;",
  "function labFiltersRefreshChoices(){",
  "function labFiltersMatches(f){",
  "function labFiltersSort(matches){",
  "function labLeave(){",
  "const labNav=$('labDiscoverNav')",
  "labFilterScopeSlot.appendChild(",
  "labFilterStatusSlot.appendChild(",
  "labFilterPlatformSlot.appendChild(",
 }
 for _,snippet:=range required {
  at:=strings.Index(page,snippet)
  if at<0{t.Errorf("falta inicialización: %s",snippet);continue}
  if at>=bootAt{t.Errorf("componente sin inicializar al arrancar: %s",snippet)}
 }
 if strings.Count(page,`id="labFilterOpen"`)!=1{t.Fatal("botón de filtros ausente/duplicado")}
 if strings.Count(page,`id="labFiltersOverlay"`)!=1{t.Fatal("ventana de filtros ausente/duplicada")}
 if strings.Contains(page,"function renderList(){labFiltersRefreshChoices();")&&
  strings.Index(page,"const labFilterDecade=")>bootAt {
  t.Fatal("ReferenceError por acceso a const labFilterDecade antes de inicializar")
 }
}

// Verificar SINTAXIS DEL JS COMPLETO INYECTADO, no solo cada módulo por separado.
// Son las mismas cadenas que ejecuta WebView2.
func TestLabMergedJavaScriptSyntax(t *testing.T){
 node,err:=exec.LookPath("node")
 if err!=nil{t.Skip("Node no disponible para verificación sintáctica")}
 page,err:=makeHTML("☐ Ejemplo (2024)","smoke-token")
 if err!=nil{t.Fatal(err)}
 lastStart:=strings.LastIndex(page,"<script>")
 lastEnd:=strings.LastIndex(page,"</script>")
 if lastStart<0||lastEnd<=lastStart{t.Fatal("no se encontró el script principal")}
 src:=page[lastStart+len("<script>"):lastEnd]
 path:=filepath.Join(t.TempDir(),"app-lab.js")
 if err:=os.WriteFile(path,[]byte(src),0600);err!=nil{t.Fatal(err)}
 cmd:=exec.Command(node,"--check",path)
 out,err:=cmd.CombinedOutput()
 if err!=nil{t.Fatalf("JavaScript completo no compila: %v\n%s",err,out)}
}

// Prueba de integración con un navegador real en Windows; detecta excepciones
// iniciales, no solo compilaciones de Go y presencia de elementos HTML.
func TestLabStartupHeadlessBrowser(t *testing.T){
 if os.Getenv("CI")==""{t.Skip("Prueba de navegador habilitada en CI de Windows")}
 var browser string
 candidates:=[]string{
  filepath.Join(os.Getenv("ProgramFiles(x86)"),"Microsoft","Edge","Application","msedge.exe"),
  filepath.Join(os.Getenv("ProgramFiles"),"Google","Chrome","Application","chrome.exe"),
  filepath.Join(os.Getenv("ProgramFiles(x86)"),"Google","Chrome","Application","chrome.exe"),
 }
 for _,p:=range candidates{if info,err:=os.Stat(p);err==nil&&!info.IsDir(){browser=p;break}}
 if browser==""{t.Skip("No hay Edge/Chrome instalados para prueba de interfaz")}
 page,err:=makeHTML("☐ Ejemplo de prueba (2024)\n☑ Otra película (1988)","smoke-token")
 if err!=nil{t.Fatal(err)}
 watchScript:=`<script>
 window.__labSmokeErrors=[];
 window.addEventListener('error',e=>window.__labSmokeErrors.push(e.message||String(e.error)));
 window.addEventListener('unhandledrejection',e=>window.__labSmokeErrors.push(String(e.reason)));
 window.addEventListener('DOMContentLoaded',()=>{
   setTimeout(()=>{
     try{
       const list=document.getElementById('filmList');
       const films=list?list.querySelectorAll('.film-item').length:0;
       const outside=document.querySelector('#view-library .search-scope')!==null||
         document.querySelector('#view-library .filters')!==null||
         document.querySelector('#view-library .library-platform-row')!==null;
       const filters=document.getElementById('labFilterOpen');
       filters?.click();
       const filtersOpened=!document.getElementById('labFiltersOverlay')?.hidden;
       const platformList=document.getElementById('listPlatformOptions');
       const platformAll=document.getElementById('labFilterAllPlatforms');
       const allSelectedByDefault=platformAll?.checked===true;
       const noPlatformDropdown=platformList?.parentElement?.id==='labFilterVisiblePlatforms'&&
         document.getElementById('listPlatformMenu')?.hidden===true;
       const platformsVisible=(platformList?.getBoundingClientRect().height||0)>0;
       const platformOptionsPresent=(platformList?.querySelectorAll('input[type="checkbox"]').length||0)>=5;
       const platformFirst=platformList?.querySelector('input[type="checkbox"]');
       if(platformFirst){platformFirst.checked=true;platformFirst.dispatchEvent(new Event('change',{bubbles:true}));}
       const allDeselectedAfterFilter=platformAll?.checked===false;
       const resetOutside=document.getElementById('labFilterResetAll')?.parentElement?.id==='labLibraryToolbar';
       document.getElementById('labFilterDone')?.click();
       const filtersClosed=document.getElementById('labFiltersOverlay')?.hidden===true;
       document.getElementById('labFilterResetAll')?.click();
       const allRestoredAfterReset=platformAll?.checked===true;
       document.querySelector('[data-view="statistics"]')?.click();
       const statsVisible=document.getElementById('view-statistics')?.hidden===false;
       const statsCards=document.getElementById('statsHero')?.children.length||0;
       document.getElementById('labDiscoverNav')?.click();
       const discoverVisible=document.getElementById('labDiscoverView')?.hidden===false;
       const genreControl=document.getElementById('labDiscoverGenre');
       const durationControl=document.getElementById('labDiscoverDuration');
       const discoverFiltersVisible=!!genreControl&&!genreControl.disabled&&
         !!durationControl&&!durationControl.disabled;
       if(genreControl){genreControl.value='878';genreControl.dispatchEvent(new Event('change',{bubbles:true}));}
       if(durationControl){durationControl.value='short';durationControl.dispatchEvent(new Event('change',{bubbles:true}));}
       const discoverFiltersSelected=genreControl?.value==='878'&&durationControl?.value==='short';
       const baseNav=document.querySelector('.top-nav [data-view="library"]');
       const helpNav=document.querySelector('.top-nav [data-view="instructions"]');
       const helpSameHeight=Math.abs((baseNav?.getBoundingClientRect().height||0)-
         (helpNav?.getBoundingClientRect().height||-100))<=1;
       document.querySelector('[data-view="collections"]')?.click();
       const collectionsVisible=document.getElementById('view-collections')?.hidden===false;
       document.querySelector('[data-view="instructions"]')?.click();
       const instructionsVisible=document.getElementById('view-instructions')?.hidden===false;
       document.querySelector('[data-view="library"]')?.click();
       const libraryVisible=document.getElementById('view-library')?.hidden===false;
       // Validar el sorteo clásico, cuya acción se había roto al no
       // insertar labCreateDetailButton en algunos entornos.
       document.getElementById('choose')?.click();
       const resultText=document.getElementById('resultTitle')?.textContent?.trim()||'';
       const sortWorked=resultText.length>0&&!resultText.includes('aparecerá aquí');
       const detailsVisible=document.getElementById('labSmartResultDetails')?.hidden===false;
       const result={films,outside,filtersOpened,filtersClosed,allSelectedByDefault,
         noPlatformDropdown,platformsVisible,platformOptionsPresent,
         allDeselectedAfterFilter,resetOutside,allRestoredAfterReset,statsVisible,statsCards,
         discoverVisible,discoverFiltersVisible,discoverFiltersSelected,helpSameHeight,
         collectionsVisible,instructionsVisible,libraryVisible,
         sortWorked,detailsVisible,errors:window.__labSmokeErrors.slice(0,5)};
       document.body.setAttribute('data-lab-smoke',encodeURIComponent(JSON.stringify(result)));
     }catch(err){
       document.body.setAttribute('data-lab-smoke',encodeURIComponent(JSON.stringify({error:String(err)})));
     }
   },1000);
 });
 </script>`
 // El detector debe instalarse antes del script principal y su evento DCL
 // se dispara incluso cuando el arranque de la app no llega a completarse.
 page=strings.Replace(page,"</head>",watchScript+"\n</head>",1)
 srv:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path=="/"{
   w.Header().Set("Content-Type","text/html; charset=utf-8")
   _,_=w.Write([]byte(page));return
  }
  w.Header().Set("Content-Type","application/json")
  w.WriteHeader(http.StatusNotFound)
  _,_=w.Write([]byte(`{"error":"sin backend en prueba UI"}`))
 }))
 defer srv.Close()
 ctx,cancel:=context.WithTimeout(context.Background(),35*time.Second);defer cancel()
 args:=[]string{
  "--headless=new","--disable-gpu","--no-first-run","--no-default-browser-check",
  "--disable-background-networking","--disable-extensions","--no-sandbox",
  "--user-data-dir="+filepath.Join(t.TempDir(),"chrome-profile"),
  "--virtual-time-budget=4500","--dump-dom",srv.URL+"/",
 }
 cmd:=exec.CommandContext(ctx,browser,args...)
 out,err:=cmd.CombinedOutput()
 if err!=nil{t.Fatalf("el navegador de prueba falló: %v\n%s",err,string(out[max(0,len(out)-1500):]))}
 marker:=`data-lab-smoke="`
 index:=strings.Index(string(out),marker)
 if index<0{t.Fatalf("el navegador no completó la prueba de arranque: %s",string(out[max(0,len(out)-1200):]))}
 end:=strings.Index(string(out)[index+len(marker):],`"`)
 if end<0{t.Fatal("marcador de humo truncado")}
 value:=string(out)[index+len(marker):index+len(marker)+end]
 // Los valores JSON están percent-encoded para evitar escapes HTML.
 if strings.Contains(value,"films%22%3A0")||strings.Contains(value,"outside%22%3Atrue")||
  strings.Contains(value,"filtersOpened%22%3Afalse")||strings.Contains(value,"filtersClosed%22%3Afalse")||
  strings.Contains(value,"allSelectedByDefault%22%3Afalse")||strings.Contains(value,"noPlatformDropdown%22%3Afalse")||
  strings.Contains(value,"platformsVisible%22%3Afalse")||strings.Contains(value,"platformOptionsPresent%22%3Afalse")||
  strings.Contains(value,"allDeselectedAfterFilter%22%3Afalse")||strings.Contains(value,"resetOutside%22%3Afalse")||
  strings.Contains(value,"allRestoredAfterReset%22%3Afalse")||
  strings.Contains(value,"statsVisible%22%3Afalse")||strings.Contains(value,"statsCards%22%3A0")||
  strings.Contains(value,"discoverVisible%22%3Afalse")||
  strings.Contains(value,"discoverFiltersVisible%22%3Afalse")||
  strings.Contains(value,"discoverFiltersSelected%22%3Afalse")||
  strings.Contains(value,"helpSameHeight%22%3Afalse")||
  strings.Contains(value,"collectionsVisible%22%3Afalse")||
  strings.Contains(value,"instructionsVisible%22%3Afalse")||strings.Contains(value,"libraryVisible%22%3Afalse")||
  strings.Contains(value,"sortWorked%22%3Afalse")||strings.Contains(value,"detailsVisible%22%3Afalse")||
  strings.Contains(value,"error%22%3A"){
   t.Fatalf("arranque no funcional en navegador: %s",value)
 }
 for _,expected:=range []string{"films%22%3A","outside%22%3Afalse","filtersOpened%22%3Atrue","filtersClosed%22%3Atrue",
 "allSelectedByDefault%22%3Atrue","noPlatformDropdown%22%3Atrue","platformsVisible%22%3Atrue",
 "platformOptionsPresent%22%3Atrue","allDeselectedAfterFilter%22%3Atrue",
 "resetOutside%22%3Atrue","allRestoredAfterReset%22%3Atrue","statsVisible%22%3Atrue","statsCards%22%3A5","discoverVisible%22%3Atrue","discoverFiltersVisible%22%3Atrue",
 "discoverFiltersSelected%22%3Atrue","helpSameHeight%22%3Atrue","collectionsVisible%22%3Atrue","instructionsVisible%22%3Atrue","libraryVisible%22%3Atrue","sortWorked%22%3Atrue","detailsVisible%22%3Atrue"}{
  if !strings.Contains(value,expected){t.Errorf("prueba de navegador sin verificar %s: %s",expected,value)}
 }
 if t.Failed(){t.Log(fmt.Sprintf("Edge/Chrome devolvió %d bytes",len(out)))}
}
