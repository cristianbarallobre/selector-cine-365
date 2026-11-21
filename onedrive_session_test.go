package main

import (
 "encoding/json"
 "errors"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
 "time"
)

// In-memory fake tests the complete session lifecycle on every OS. Actual
// Windows encryption and disk isolation are tested in a windows-only file.
type testGraphSessionStore struct{
 saved graphStoredSession
 present bool
 writes,reads,deletes int
 failSave,failClear bool
}
func (s *testGraphSessionStore) load()(graphStoredSession,error){
 s.reads++
 if !s.present{return graphStoredSession{},errors.New("unavailable")}
 return s.saved,nil
}
func (s *testGraphSessionStore) save(sess graphStoredSession)error{
 if s.failSave{return errors.New("test disk failure")}
 s.saved=sess;s.present=true;s.writes++
 return nil
}
func (s *testGraphSessionStore) clear()error{
 if s.failClear{return errors.New("test delete failure")}
 s.present=false;s.saved=graphStoredSession{};s.deletes++
 return nil
}

func TestMicrosoftSessionRestoreRefreshRotationAndDisconnect(t *testing.T){
 cache:=&testGraphSessionStore{}
 g:=newDesktopGraph()
 g.sessionStore=cache
 g.rootID="one-drive-folder-id"
 g.refresh="initial-test-refresh"
 g.token="short-lived-test-access"
 g.expiry=time.Now().Add(time.Hour)
 g.mu.Lock();g.rememberSessionLocked();g.mu.Unlock()
 if !cache.present||cache.writes!=1||cache.saved.RootID!="one-drive-folder-id" {t.Fatal("session not saved")}
 if cache.saved.Refresh!="initial-test-refresh"{t.Fatal("wrong refresh token stored")}
 if stringValue(g.status(),"unused")!=""{t.Fatal("unexpected test helper")}

 var seenRefresh string
 oauth:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path!="/token"{t.Errorf("unexpected OAuth path %s",r.URL.Path);w.WriteHeader(404);return}
  if err:=r.ParseForm();err!=nil{t.Error(err)}
  if r.Form.Get("grant_type")!="refresh_token"{t.Error("must refresh without opening browser")}
  seenRefresh=r.Form.Get("refresh_token")
  _=json.NewEncoder(w).Encode(map[string]any{
   "access_token":"rotated-test-access","refresh_token":"rotated-test-refresh","expires_in":3600,
  })
 }))
 defer oauth.Close()
 restored:=newDesktopGraph()
 restored.sessionStore=cache
 restored.restoreSession()
 restored.authBase=oauth.URL+"/"
 restored.client=oauth.Client()
 if !restored.active(){t.Fatal("saved refresh token should restore connected state")}
 if restored.token!=""{t.Fatal("access tokens must never be persisted")}
 if restored.rootID!="one-drive-folder-id" {t.Fatal("drive binding lost")}
 if restored.status()["session_restored"]!=true {t.Fatal("restored state not advertised")}
 if restored.status()["session_remembered"]!=true {t.Fatal("session must be remembered")}
 token,err:=restored.accessToken()
 if err!=nil||token!="rotated-test-access"{t.Fatalf("automatic refresh failed: %v",err)}
 if seenRefresh!="initial-test-refresh"{t.Fatal("refresh token lost on restart")}
 if cache.saved.Refresh!="rotated-test-refresh"||cache.writes!=2 {t.Fatal("rotated token not securely persisted")}
 again:=newDesktopGraph()
 again.sessionStore=cache;again.restoreSession()
 if !again.active()||again.refresh!="rotated-test-refresh" {t.Fatal("second app update should reuse updated session")}
 if err=again.logout();err!=nil{t.Fatal(err)}
 if cache.present||cache.deletes!=1||again.active() {t.Fatal("disconnect must delete saved session")}
 third:=newDesktopGraph();third.sessionStore=cache;third.restoreSession()
 if third.active(){t.Fatal("account must not return after logout")}
}

func TestMicrosoftSessionTransientOutagePreservesAccount(t *testing.T){
 cache:=&testGraphSessionStore{}
 cache.saved=graphStoredSession{Version:graphSessionVersion,ClientID:desktopMicrosoftAppID,RootID:"folder",Refresh:"refresh"}
 cache.present=true
 g:=newDesktopGraph();g.sessionStore=cache;g.restoreSession()
 oauth:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.WriteHeader(503);_,_=w.Write([]byte(`{"error":"temporarily_unavailable"}`))
 }))
 defer oauth.Close()
 g.authBase=oauth.URL+"/";g.client=oauth.Client()
 if _,err:=g.accessToken();err==nil{t.Fatal("expected temporary provider error")}
 if !g.active()||!cache.present||cache.deletes!=0{t.Fatal("temporary outage must not revoke locally stored session")}
}

func TestMicrosoftSessionRevokedTokenRequiresReconnect(t *testing.T){
 cache:=&testGraphSessionStore{
  saved:graphStoredSession{Version:graphSessionVersion,ClientID:desktopMicrosoftAppID,RootID:"folder",Refresh:"revoked"},
  present:true,
 }
 g:=newDesktopGraph();g.sessionStore=cache;g.restoreSession()
 oauth:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  _=json.NewEncoder(w).Encode(map[string]any{"error":"invalid_grant"})
 }))
 defer oauth.Close()
 g.authBase=oauth.URL+"/";g.client=oauth.Client()
 if _,err:=g.accessToken();err==nil||!strings.Contains(err.Error(),"Vuelve a conectar"){
  t.Fatalf("expected new login request, got %v",err)
 }
 if g.active()||cache.present{t.Fatal("revoked token should not keep a fake connection")}
}

func TestMicrosoftSessionDeletionErrorKeepsState(t *testing.T){
 cache:=&testGraphSessionStore{present:true,failClear:true}
 g:=newDesktopGraph();g.sessionStore=cache;g.refresh="active-refresh";g.rootID="root"
 if err:=g.logout();err==nil{t.Fatal("failed credential deletion cannot be reported as successful logout")}
 if !g.active(){t.Fatal("failure to delete saved token must not hide existing connection")}
}

func TestMicrosoftSessionCannotStoreInvalidOrForeignApplicationTokens(t *testing.T){
 for _,sess:=range []graphStoredSession{
  {},
  {Version:graphSessionVersion,ClientID:"unrelated-app",RootID:"folder",Refresh:"token"},
  {Version:graphSessionVersion,ClientID:desktopMicrosoftAppID,RootID:"",Refresh:"token"},
  {Version:graphSessionVersion,ClientID:desktopMicrosoftAppID,RootID:"folder",Refresh:""},
 }{
  if _,err:=encodeGraphSession(sess);err==nil{t.Fatalf("invalid session accepted: %+v",sess.Version)}
 }
}
