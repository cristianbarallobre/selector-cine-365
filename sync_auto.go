package main

import (
 "errors"
 "reflect"
 "time"
)

// cloudAutoSync imports only conflict-free changes after the user has explicitly
// paired a folder in the manual preview/confirmation flow. No silent initial pairing.
func (a *runningApp) cloudAutoSync() (map[string]any,error) {
 a.syncPlanMu.Lock();defer a.syncPlanMu.Unlock()
 if a.cloud==nil || a.syncStore==nil {
  return map[string]any{"changed":false,"skipped":"disabled"},nil
 }
 folder:=a.cloud.status().Folder
 if folder=="" || !a.syncStore.isPaired(folder) {
  return map[string]any{"changed":false,"skipped":"not_paired","needs_review":folder!=""},nil
 }
 a.mu.Lock()
 local,revision:=a.library,a.libraryRevision
 a.mu.Unlock()
 plan,err:=a.cloudPlanFor(local,folder)
 if err!=nil{return nil,err}
 a.mu.Lock()
 if revision!=a.libraryRevision||folder!=a.cloud.status().Folder{
  a.mu.Unlock()
  return map[string]any{"changed":false,"skipped":"local_changed"},nil
 }
 a.syncLastCheck.Store(time.Now().UnixMilli())
 if len(plan.Conflicts)!=0{
  a.syncPendingConflict.Store(true)
  a.mu.Unlock()
  return map[string]any{"changed":false,"needs_review":true,"conflicts":len(plan.Conflicts)},nil
 }
 a.syncPendingConflict.Store(false)
 if plan.Sources==0{
  a.mu.Unlock()
  return map[string]any{"changed":false,"skipped":"no_remote_sources"},nil
 }
 _,beforeRevisions,beforeKnown:=a.syncStore.parts()
 libChanged:=!reflect.DeepEqual(a.library,plan.Library)
 metadataChanged:=!reflect.DeepEqual(beforeRevisions,plan.Revisions)||!reflect.DeepEqual(beforeKnown,plan.Known)
 if !libChanged&&!metadataChanged {
  a.mu.Unlock()
  return map[string]any{"changed":false},nil
 }
 if libChanged {
  if err=a.saveSyncRecovery(a.library,"");err!=nil{a.mu.Unlock();return nil,err}
  if err=saveLibrary(a.libraryPath,a.savePath,plan.Library);err!=nil{a.mu.Unlock();return nil,err}
  a.library=plan.Library
  a.libraryRevision++
  a.syncLastReceived.Store(time.Now().UnixMilli())
 }
 if err=a.syncStore.accept(plan,nil,folder);err!=nil{a.mu.Unlock();return nil,err}
 a.syncPreviewToken=""
 a.mu.Unlock()
 // Remote recovery points remain separate from the local rescue snapshot.
 cloudResult:=a.cloud.backup()
 if cloudResult.LastError!=""{return nil,errors.New(cloudResult.LastError)}
 if err=a.publishDeviceSnapshot();err!=nil{return nil,err}
 return map[string]any{
  "changed":libChanged,"metadata_updated":metadataChanged,
  "lists":len(plan.Library.Lists),"seen":syncWatched(plan.Library),
  "sources":plan.Sources,
 },nil
}
