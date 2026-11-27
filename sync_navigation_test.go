package main

import (
 "path/filepath"
 "testing"
)

func TestDesktopNavigationDoesNotMarkUploadPending(t *testing.T){
 if !desktopFeatures{t.Skip("desktop")}
 a:=durableTestApp(t)
 a.library.Lists=append(a.library.Lists,filmCollection{ID:"second",Name:"Otra",Text:"☐ Otra (2011)"})
 a.cloud=newCloudBackups(filepath.Dir(a.libraryPath))
 if err:=a.cloud.configure(t.TempDir());err!=nil{t.Fatal(err)}
 // Establish a saved state, then simulate acknowledgement of real changes.
 if w:=postTestLibrary(t,a,a.library);w.Code!=204{t.Fatal(w.Code,w.Body.String())}
 if err:=a.syncStore.acknowledgeGraphUpload(a.syncStore.graphPublicationVersion());err!=nil{t.Fatal(err)}
 navigation:=a.library;navigation.ActiveID="second"
 if !sameShareableLibrary(a.library,navigation){t.Fatal("navigation treated as content")}
 if w:=postTestLibrary(t,a,navigation);w.Code!=204{t.Fatal(w.Code,w.Body.String())}
 if a.syncStore.graphUploadPending(){t.Fatal("navigation scheduled an upload")}
 if got:=loadLibrary(a.libraryPath,"");got.ActiveID!="second"{t.Fatal("navigation not saved locally")}
 edited:=navigation;edited.Lists=append([]filmCollection(nil),navigation.Lists...);edited.Lists[1].Text="☑ Otra (2011)"
 if w:=postTestLibrary(t,a,edited);w.Code!=204{t.Fatal(w.Code,w.Body.String())}
 if !a.syncStore.graphUploadPending(){t.Fatal("real edit not scheduled")}
 again:=edited;again.ActiveID="atlas"
 if w:=postTestLibrary(t,a,again);w.Code!=204{t.Fatal(w.Code,w.Body.String())}
 if !a.syncStore.graphUploadPending(){t.Fatal("navigation erased pending real edit")}
}
