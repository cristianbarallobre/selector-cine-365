package main
import("fmt";"strings";"testing";"path/filepath";"time")
func BenchmarkLibraryActivity(b *testing.B){
 for _,n:=range []int{1000,5000,10000}{
  b.Run(fmt.Sprint(n),func(b *testing.B){
   var rows strings.Builder
   for i:=0;i<n;i++{fmt.Fprintf(&rows,"☐ Película %d (2020)\n",i)}
   before:=initialLibrary(rows.String());after:=before
   after.Lists=append([]filmCollection(nil),before.Lists...)
   after.Lists[0].Text+="☑ Nueva (2025)"
   b.ReportAllocs();b.ResetTimer()
   for i:=0;i<b.N;i++{_ = activityDiff(before,after)}
  })
 }
}

func BenchmarkDesktopCommit(b *testing.B){
 if !desktopFeatures{b.Skip("desktop")}
 for _,n:=range []int{1000,5000,10000}{
  b.Run(fmt.Sprint(n),func(b *testing.B){
   var rows strings.Builder;for i:=0;i<n;i++{fmt.Fprintf(&rows,"☐ Película %d (2020)\n",i)}
   dir:=b.TempDir();s,e:=newDesktopSyncStore(dir);if e!=nil{b.Fatal(e)}
   a:=&runningApp{library:initialLibrary(rows.String()),libraryPath:filepath.Join(dir,"colecciones.json"),savePath:filepath.Join(dir,"mis_365_peliculas.txt"),libraryDigest:"missing",syncStore:s,extraBackup:newExtraBackup(dir)}
   if _,_,e=a.commitProtectedLibrary(a.library,nil,nil,"");e!=nil{b.Fatal(e)}
   b.ReportAllocs();b.ResetTimer()
   for i:=0;i<b.N;i++{
    next:=a.library;next.Lists=append([]filmCollection(nil),next.Lists...)
    if i%2==0{next.Lists[0].Text="☑"+strings.TrimPrefix(next.Lists[0].Text,"☐")}else{next.Lists[0].Text="☐"+strings.TrimPrefix(next.Lists[0].Text,"☑")}
    if _,_,e=a.commitProtectedLibrary(next,nil,nil,"");e!=nil{b.Fatal(e)}
    a.library=next
   }
  })
 }
}
func BenchmarkProviderCacheHit(b *testing.B){
 p:=&availabilityServer{cache:map[string]providerResult{"test":{State:"ok",TMDBID:1,MetadataVersion:3,Date:time.Now().Format(time.RFC3339)}}}
 b.ReportAllocs();b.ResetTimer()
 for i:=0;i<b.N;i++{if _,ok:=p.getCached("test",false);!ok{b.Fatal("cache miss")}}
}
