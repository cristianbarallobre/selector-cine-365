package main
import("fmt";"strings";"testing")
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
