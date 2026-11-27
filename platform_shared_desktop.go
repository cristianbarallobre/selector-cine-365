//go:build windows || darwin

package main
import("encoding/json";"time")

func cineValidAvailability(r providerResult) (time.Time,bool){
 at,e:=time.Parse(time.RFC3339,r.Date)
 return at,e==nil&&r.TMDBID>0&&(r.State=="ok"||r.State=="no-data")&&!at.After(time.Now().Add(5*time.Minute))
}
func cineProviderToDetail(rec labDetailRecord,r providerResult) labDetailRecord{
 at,ok:=cineValidAvailability(r);if !ok||at.Unix()<=rec.ProvidersAt{return rec}
 rec.Data=labCopyMap(rec.Data)
 rec.Data["subscriptions"]=r.Subscriptions;rec.Data["free"]=r.Free;rec.Data["ads"]=r.Ads
 rec.Data["rent"]=r.Rent;rec.Data["buy"]=r.Buy;rec.Data["watch_link"]=r.FilmURL;rec.Data["availability_note"]=""
 rec.ProvidersAt=at.Unix()
 return rec
}
func(c *labCatalogCache) mergeLibraryAvailability(id int,rec labDetailRecord)labDetailRecord{
 if c.app==nil||c.app.platforms==nil||id<=0{return rec}
 p:=c.app.platforms
 p.mu.Lock();defer p.mu.Unlock()
 var best providerResult;var newest time.Time
 for key,r:=range p.cache{
  if r.TMDBID!=id||(p.overrides[key]>0&&p.overrides[key]!=id){continue}
  at,ok:=cineValidAvailability(r);if ok&&at.After(newest){best=r;newest=at}
 }
 return cineProviderToDetail(rec,best)
}
func cineDetailAvailability(id int,rec labDetailRecord)(providerResult,bool){
 var r providerResult
 if id<=0||rec.ProvidersAt<=0||time.Unix(rec.ProvidersAt,0).After(time.Now().Add(5*time.Minute)){return r,false}
 data,e:=json.Marshal(rec.Data);if e!=nil||json.Unmarshal(data,&r)!=nil{return r,false}
 r.TMDBID=id;r.State="ok";r.Date=time.Unix(rec.ProvidersAt,0).UTC().Format(time.RFC3339)
 r.FilmURL,_=rec.Data["watch_link"].(string)
 return r,true
}
func(c *labCatalogCache) publishLibraryAvailability(id int,rec labDetailRecord){
 if c.app==nil||c.app.platforms==nil{return}
 incoming,ok:=cineDetailAvailability(id,rec);if !ok{return}
 p:=c.app.platforms
 p.mu.Lock();defer p.mu.Unlock()
 for key,prior:=range p.cache{
  if prior.TMDBID!=id||(p.overrides[key]>0&&p.overrides[key]!=id){continue}
  at,valid:=cineValidAvailability(prior);if valid&&at.Unix()>=rec.ProvidersAt{continue}
  prior.State="ok";prior.Date=incoming.Date;prior.FilmURL=incoming.FilmURL
  prior.Subscriptions=incoming.Subscriptions;prior.Free=incoming.Free;prior.Ads=incoming.Ads
  prior.Rent=incoming.Rent;prior.Buy=incoming.Buy;prior.Note=""
  // Do not falsely certify metadata that was only partial.
  p.putCachedLocked(key,prior)
 }
}
func cineKeepNewerDetailAvailability(id int,rec,prior labDetailRecord)labDetailRecord{
 if prior.ProvidersAt<=rec.ProvidersAt{return rec}
 if incoming,ok:=cineDetailAvailability(id,prior);ok{return cineProviderToDetail(rec,incoming)}
 return rec
}

func cineKeepLibraryProvider(prior,next providerResult)bool{
 old,ok:=cineValidAvailability(prior);if !ok{return false}
 at,valid:=cineValidAvailability(next)
 return !valid||!at.After(old)
}
