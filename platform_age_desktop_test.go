//go:build windows || darwin

package main

import (
 "testing"
 "time"
)
func TestPlatformAgeDetailResponse(t *testing.T){
 rec:=labDetailRecord{Data:map[string]any{"title":"Test","providers_checked_at":"no fiable"},MetaAt:1700000000,ProvidersAt:1600000000}
 data:=labDetailWithAvailabilityDate(rec)
 if data["providers_checked_at"]!=time.Unix(rec.ProvidersAt,0).UTC().Format(time.RFC3339){t.Fatal(data)}
 if rec.Data["providers_checked_at"]!="no fiable"{t.Fatal("se modificó la caché")}
 rec.ProvidersAt=0
 if _,ok:=labDetailWithAvailabilityDate(rec)["providers_checked_at"];ok{t.Fatal("fecha inventada")}
}
