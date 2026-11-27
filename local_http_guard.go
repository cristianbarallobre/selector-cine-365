package main

import (
 "net/http"
 "strings"
)

// Trust only the exact loopback listener address. This blocks DNS rebinding,
// cross-origin browser requests and framing of the token-bearing local page.
func localHTTPGuard(host string,next http.Handler) http.Handler {
 origin:="http://"+host
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Referrer-Policy","no-referrer")
  w.Header().Set("X-Frame-Options","DENY")
  w.Header().Set("Content-Security-Policy","frame-ancestors 'none'; base-uri 'none'; object-src 'none'; form-action 'self'")
  if r.Host!=host {http.Error(w,"Dirección local no autorizada",http.StatusForbidden);return}
  if o:=r.Header.Get("Origin");o!=""&&o!=origin{http.Error(w,"Origen no autorizado",http.StatusForbidden);return}
  if site:=strings.ToLower(r.Header.Get("Sec-Fetch-Site"));site!=""&&site!="same-origin"&&site!="none"{http.Error(w,"Origen no autorizado",http.StatusForbidden);return}
  next.ServeHTTP(w,r)
 })
}
