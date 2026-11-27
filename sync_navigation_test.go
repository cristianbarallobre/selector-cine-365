package main
import("path/filepath";"testing")
func TestDesktopNavigationDoesNotRecordActivity(t *testing.T){
 if !desktopFeatures{t.Skip("desktop")}
 a:=durableTestApp(t);a.extraBackup=newExtraBackup(filepath.Dir(a.libraryPath))
 a.library.Lists=append(a.library.Lists,filmCollection{ID:"second",Name:"Otra",Text:"☐ Otra (2011)"})
 if w:=postTestLibrary(t,a,a.library);w.Code!=204{t.Fatal(w.Code,w.Body.String())}
 navigation:=a.library;navigation.ActiveID="second"
 if !sameShareableLibrary(a.library,navigation){t.Fatal("navigation treated as content")}
 if w:=postTestLibrary(t,a,navigation);w.Code!=204{t.Fatal(w.Code,w.Body.String())}
 if len(a.syncStore.localActivity())!=0{t.Fatal("navigation recorded")}
 if got:=loadLibrary(a.libraryPath,"");got.ActiveID!="second"{t.Fatal("navigation not saved")}
 edited:=navigation;edited.Lists=append([]filmCollection(nil),navigation.Lists...);edited.Lists[1].Text="☑ Otra (2011)"
 if w:=postTestLibrary(t,a,edited);w.Code!=204{t.Fatal(w.Code,w.Body.String())}
 if len(a.syncStore.localActivity())!=1{t.Fatal("edit missing")}
 again:=edited;again.ActiveID="atlas"
 if w:=postTestLibrary(t,a,again);w.Code!=204{t.Fatal(w.Code,w.Body.String())}
 if len(a.syncStore.localActivity())!=1{t.Fatal("navigation changed history")}
}
