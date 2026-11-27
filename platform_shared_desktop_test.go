//go:build windows || darwin

package main
import("testing";"time";"path/filepath")
func TestSharedProviderAvailabilityUsesNewestSuccessfulResult(t *testing.T){
 now:=time.Now().UTC().Truncate(time.Second)
 p:=newAvailabilityServer(t.TempDir())
 p.cache["film|2000"]=providerResult{TMDBID:42,State:"ok",Date:now.Format(time.RFC3339),Subscriptions:[]string{"Netflix"},Runtime:123}
 p.cache["alias|2000"]=providerResult{TMDBID:42,State:"error",Date:now.Add(time.Minute).Format(time.RFC3339)}
 p.cache["different|2000"]=providerResult{TMDBID:99,State:"ok",Date:now.Format(time.RFC3339),Subscriptions:[]string{"Other"}}
 c:=&labCatalogCache{app:&runningApp{platforms:p},file:filepath.Join(t.TempDir(),"details.json"),disk:labCacheDisk{Items:map[string]labDetailRecord{},Resolved:map[string]int{}}}
 prior:=labDetailRecord{Data:map[string]any{"title":"Test","subscriptions":[]string{"Old"}},ProvidersAt:now.Add(-time.Hour).Unix()}
 fresh:=c.mergeLibraryAvailability(42,prior)
 if fresh.ProvidersAt!=now.Unix()||fresh.Data["subscriptions"].([]string)[0]!="Netflix"||prior.Data["subscriptions"].([]string)[0]!="Old"{t.Fatal(fresh,prior)}
 if c.mergeLibraryAvailability(99,fresh).ProvidersAt!=fresh.ProvidersAt{t.Fatal("same-date availability replaced")}
 newer:=fresh;newer.Data=labCopyMap(fresh.Data);newer.Data["subscriptions"]=[]string{"New"};newer.ProvidersAt=now.Add(time.Second).Unix()
 c.publishLibraryAvailability(42,newer)
 if got:=p.cache["film|2000"];got.Subscriptions[0]!="New"||got.Runtime!=123{t.Fatal(got)}
 if p.cache["different|2000"].Subscriptions[0]!="Other"{t.Fatal("different movie modified")}
 c.publishLibraryAvailability(42,prior)
 if p.cache["film|2000"].Subscriptions[0]!="New"{t.Fatal("stale response overwrote newer data")}
 p.flushProviderCache()
}
