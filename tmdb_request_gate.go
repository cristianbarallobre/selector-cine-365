package main
import("context";"errors";"net/http";"strconv";"sync";"time")
// Shared by discovery, identity, library and IMDb ID resolution.
var tmdbSlots=make(chan struct{},4)
var tmdbRateState struct{sync.Mutex;client *http.Client;until time.Time}
func tmdbCoolingDown(client *http.Client)bool{
 tmdbRateState.Lock();defer tmdbRateState.Unlock()
 if tmdbRateState.client!=client{tmdbRateState.client=client;tmdbRateState.until=time.Time{}}
 return time.Now().Before(tmdbRateState.until)
}
func acquireTMDBRequest(ctx context.Context,client *http.Client)(func(),error){
 limited:=errors.New("Límite temporal de TMDB. Prueba a actualizar más tarde")
 if tmdbCoolingDown(client){return nil,limited}
 select{case <-ctx.Done():return nil,ctx.Err();case tmdbSlots<-struct{}{}:}
 release:=func(){<-tmdbSlots}
 if tmdbCoolingDown(client){release();return nil,limited}
 return release,nil
}
func deferTMDBRequests(client *http.Client,header string){
 delay:=time.Minute
 if n,e:=strconv.Atoi(header);e==nil&&n>0&&n<=300{delay=time.Duration(n)*time.Second}else if date,e:=http.ParseTime(header);e==nil{delay=time.Until(date)}
 if delay<time.Second{delay=time.Second};if delay>5*time.Minute{delay=5*time.Minute}
 tmdbRateState.Lock();defer tmdbRateState.Unlock()
 if tmdbRateState.client!=client{tmdbRateState.client=client;tmdbRateState.until=time.Time{}}
 until:=time.Now().Add(delay);if until.After(tmdbRateState.until){tmdbRateState.until=until}
}
