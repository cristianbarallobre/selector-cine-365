package main
import("testing";"encoding/json";"net/http";"net/http/httptest";"strings";"path/filepath")
func TestCompleteBackupRestore(t *testing.T){
 if !desktopFeatures{t.Skip("desktop")}
 a:=durableTestApp(t);a.extraBackup=newExtraBackup(filepath.Dir(a.libraryPath))
 b:=completeLibraryBackup{libraryState:a.library,CompleteVersion:1,Activity:[]activityEntry{{At:1,Kind:"film_add",Film:"Test"}}}
 raw,_:=json.Marshal(b)
 req:=httptest.NewRequest(http.MethodPost,"/api/extra-backup/restore",strings.NewReader(string(raw)));req.Header.Set("X-Atlas-Token","test");req.Header.Set("If-Match",a.libraryETag())
 w:=httptest.NewRecorder();a.serveHTTP(w,req);if w.Code!=204{t.Fatal(w.Code,w.Body.String())}
 cp:=readLibraryCheckpoint(a.libraryPath);if cp==nil||len(cp.Activity)<2||cp.Activity[0].Film!="Test"{t.Fatal("history not committed")}
 page,e:=makeHTML("☐ Origen (2010)","test");if e!=nil{t.Fatal(e)};page,e=localOnlyHTML(page);if e!=nil{t.Fatal(e)}
 for _,marker:=range []string{"/api/extra-backup/export","/api/extra-backup/restore","state.complete_backup_version===1"}{if !strings.Contains(page,marker){t.Fatal(marker)}}
}
