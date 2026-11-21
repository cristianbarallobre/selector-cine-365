package main

import (
 "fmt"
 "strings"
 "testing"
)

func TestWindowsRejectsUnrelatedListWithSameLegacyID(t *testing.T){
 win:=initialLibrary("☐ Casablanca (1942)\n☐ Metropolis (1927)\n☐ Ciudadano Kane (1941)")
 win.Lists[0].Name="Mi lista de Windows"
 android:=initialLibrary("☐ Parasite (2019)\n☐ Roma (2018)\n☐ The Matrix (1999)")
 peer:=syncSnapshot{Protocol:2,DeviceID:"android-old",Library:android}
 _,_,_,_,err:=syncMerge(win,map[string]syncRevision{},map[string]int64{},peer)
 if err==nil||!strings.Contains(err.Error(),"Sincronización detenida"){t.Fatalf("unrelated same-ID lists must be blocked: %v",err)}
 if strings.Contains(win.Lists[0].Text,"Parasite"){t.Fatal("unsafe preview changed Windows source")}
}

func TestWindowsLegitimateListExtensionStillWorks(t *testing.T){
 win:=initialLibrary("☐ Casablanca (1942)\n☐ Metropolis (1927)")
 win.Lists[0].Name="Mi lista de Windows"
 android:=initialLibrary("☐ Casablanca (1942)\n☐ Metropolis (1927)\n☐ Parasite (2019)")
 android.Lists[0].Name=win.Lists[0].Name
 peer:=syncSnapshot{Protocol:2,DeviceID:"android-same",Library:android}
 out,_,_,_,err:=syncMerge(win,map[string]syncRevision{},map[string]int64{},peer)
 if err!=nil{t.Fatal(err)}
 if !strings.Contains(out.Lists[0].Text,"Parasite"){t.Fatal("same list should accept newly added film")}
}

func TestWindowsStarterAndRenamedListCannotSwallowSubset(t *testing.T){
 // Even full overlap of a small subset is unsafe when a pristine 365 starter
 // shares the reserved "atlas" ID with a differently named Windows collection.
 win:=initialLibrary("☐ Casablanca (1942)\n☐ Metropolis (1927)\n☐ Parasite (2019)")
 android:=initialLibrary("☐ Casablanca (1942)\n☐ Parasite (2019)")
 android.Lists[0].Name="Mi selección de películas"
 peer:=syncSnapshot{Protocol:2,DeviceID:"android-other",Library:android}
 _,_,_,_,err:=syncMerge(win,map[string]syncRevision{},map[string]int64{},peer)
 if err==nil||!strings.Contains(err.Error(),"Sincronización detenida"){t.Fatalf("starter/subset collision not blocked: %v",err)}
}

func TestWindowsRejectsPreviouslyPollutedSuperset(t *testing.T){
 var expected,polluted strings.Builder
 for i:=0;i<40;i++{
  line:=fmt.Sprintf("☐ Película %d (2000)\\n",i)
  expected.WriteString(line);polluted.WriteString(line)
 }
 for i:=0;i<40;i++{polluted.WriteString(fmt.Sprintf("☐ Otra película %d (2010)\\n",i))}
 correct:=initialLibrary(expected.String())
 correct.Lists[0].Name="Mi lista"
 bad:=initialLibrary(polluted.String())
 bad.Lists[0].Name="Mi lista"
 peer:=syncSnapshot{Protocol:2,DeviceID:"android-contaminated",Library:bad}
 _,_,_,_,err:=syncMerge(correct,map[string]syncRevision{},map[string]int64{},peer)
 if err==nil||!strings.Contains(err.Error(),"Sincronización detenida"){t.Fatalf("polluted supersets must be rejected: %v",err)}
}
