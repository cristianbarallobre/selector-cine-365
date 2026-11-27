package main
import("net/http";"net/http/httptest";"testing")
func TestLocalHTTPGuard(t *testing.T){
 const host="127.0.0.1:49152"
 h:=localHTTPGuard(host,http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){w.WriteHeader(204)}))
 for _,c:=range []struct{host,origin,site string;code int}{
  {host,"","",204},{host,"http://"+host,"same-origin",204},{host,"","none",204},
  {"external.invalid:49152","","",403},{host,"https://external.invalid","",403},
  {host,"","cross-site",403},{host,"","same-site",403},{host,"null","",403},
 }{
  r:=httptest.NewRequest("GET","http://"+c.host+"/",nil)
  if c.origin!=""{r.Header.Set("Origin",c.origin)}
  if c.site!=""{r.Header.Set("Sec-Fetch-Site",c.site)}
  w:=httptest.NewRecorder();h.ServeHTTP(w,r)
  if w.Code!=c.code||w.Header().Get("Referrer-Policy")!="no-referrer"||w.Header().Get("X-Frame-Options")!="DENY"{t.Fatalf("%+v => %d",c,w.Code)}
 }
 a:=&runningApp{token:"test"}
 r:=httptest.NewRequest("POST","http://"+host+"/api/library?token=test",nil)
 if a.auth(r){t.Fatal("query token accepted for mutation")}
 r.Header.Set("X-Atlas-Token","test");if !a.auth(r){t.Fatal("header token rejected")}
}
