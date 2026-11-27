package main

import (
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
  "os"
 "path/filepath"
 "bytes"
 "sync/atomic"
 "time"
)

// The local checkpoint is never part of API responses or cloud snapshots.
type libraryDisk struct {
 libraryState
 LocalSync *syncDeviceState `json:"_local_sync,omitempty"`
}

func libraryFingerprint(path string) (string,error) {
 data,err:=readRecoveryFile(path)
 if errors.Is(err,os.ErrNotExist){return "missing",nil}
 if err!=nil{return "",err}
 return libraryBytesFingerprint(data),nil
}

func libraryBytesFingerprint(data []byte) string{
 sum:=sha256.Sum256(data)
 return hex.EncodeToString(sum[:])
}
func readLibraryCheckpoint(path string) *syncDeviceState {
 data,err:=readRecoveryFile(path);if err!=nil{return nil}
 var disk libraryDisk
 if json.Unmarshal(data,&disk)!=nil||disk.LocalSync==nil{return nil}
 s:=disk.LocalSync
 if s.Checkpoint==""||s.DeviceID==""||s.Protocol!=syncProtocol||s.Known==nil||s.Revisions==nil{return nil}
 return s
}

// Recover metadata from the same atomic library commit, even if the separate
// metadata file could not be written. A matching checkpoint may have newer
// publication acknowledgements, so preserve it.
func recoverSyncCheckpoint(dir string, checkpoint *syncDeviceState, recovered ...libraryState) (*desktopSyncStore,error) {
 path:=filepath.Join(dir,"sync_dispositivo_v2.json")
 var existing syncDeviceState
 data,err:=os.ReadFile(path)
 existingOK:=err==nil&&json.Unmarshal(data,&existing)==nil&&existing.DeviceID!=""&&existing.Known!=nil&&existing.Revisions!=nil
 var store *desktopSyncStore
 if checkpoint==nil{
  store,err=newDesktopSyncStore(dir)
  if err!=nil{return nil,err}
 }else if existingOK&&existing.Checkpoint==checkpoint.Checkpoint&&existing.DeviceID==checkpoint.DeviceID{
  store=&desktopSyncStore{path:path,state:existing}
 }else{
  store=&desktopSyncStore{path:path,state:*checkpoint}
  store.state.GraphUploadPending=true
  // A rollback must never reuse causal sequence numbers from a newer save.
  if existingOK&&existing.DeviceID==store.state.DeviceID&&existing.Seq>store.state.Seq{
   store.state.Seq=existing.Seq
   store.state.Known=syncJoin(store.state.Known,existing.Known)
   store.state.PairedFolder=""
  }
 }
 if len(recovered)>0{
  // Restoring an older local backup is not consent to publish it. Force
  // a fresh review, and do not let mismatching sidecar marks change it at startup.
  statuses:=syncStatuses(recovered[0])
  for key,revision:=range store.state.Revisions{
   if seen,ok:=statuses[key];!ok||seen!=revision.Checked{delete(store.state.Revisions,key)}
  }
  store.state.PairedFolder=""
  store.state.GraphUploadPending=true
 }
 if checkpoint!=nil||len(recovered)>0{
  if err:=store.saveLocked();err!=nil{
   diagnosticError("sincronizacion","reparar estado desde biblioteca",err)
   // Keep the checkpoint in memory; the committed library still owns it.
  }
 }
 return store,nil
}

// Caller holds a.mu. Prepare all sync metadata in memory, then commit it
// together with the library. A failed sidecar write cannot lose revisions.
func (a *runningApp) commitProtectedLibrary(next libraryState, plan *syncPlan, decisions map[string]bool, folder string, activityAction ...string) (mirrorPending bool, metadataErr error, err error) {
 if a.libraryDigest!=""{
  current,e:=libraryFingerprint(a.libraryPath)
  if e!=nil{return false,nil,e}
  if current!=a.libraryDigest{return false,nil,errors.New("La biblioteca cambió en disco. No se sobrescribió; exporta los cambios y reinicia la aplicación")}
 }
 if a.syncStore==nil||!desktopFeatures{
  mirrorPending,err=commitLibraryCheckpoint(a.libraryPath,a.savePath,next,nil,a.libraryDigest)
 }else{
  s:=a.syncStore
  s.mu.Lock()
  defer s.mu.Unlock()
  data,e:=json.Marshal(s.state);if e!=nil{return false,nil,e}
  clone:=&desktopSyncStore{graphUploadVersion:s.graphUploadVersion}
  if e=json.Unmarshal(data,&clone.state);e!=nil{return false,nil,e}
  if a.extraBackup!=nil{appendLocalActivity(&clone.state,a.library,next);if len(activityAction)>0&&activityAction[0]=="restore"{appendActivityEvent(&clone.state,activityEntry{At:time.Now().UnixMilli(),Kind:"restore",Detail:"Copia local restaurada"})}}
  if plan!=nil{e=clone.accept(*plan,decisions,folder)}else{
   e=clone.observe(a.library,next)
   if e==nil&&!sameShareableLibrary(a.library,next)&&a.cloud!=nil&&(a.cloud.status().Enabled||a.graph!=nil&&a.graph.active()){e=clone.markGraphUploadPending()}
  }
  if e!=nil{return false,nil,e}
  checkpoint,e:=randomToken();if e!=nil{return false,nil,e}
  clone.state.Checkpoint=checkpoint
  mirrorPending,err=commitLibraryCheckpoint(a.libraryPath,a.savePath,next,&clone.state,a.libraryDigest)
  if err==nil{
   s.state=clone.state
   s.graphUploadVersion=clone.graphUploadVersion
   metadataErr=s.saveLocked()
  }
 }
 if err==nil&&a.libraryDigest!=""{
  var cp *syncDeviceState
  if desktopFeatures&&a.syncStore!=nil{cp=&a.syncStore.state}
  data,e:=json.MarshalIndent(libraryDisk{libraryState:next,LocalSync:cp},"","  ")
  if e!=nil{a.libraryDigest="unavailable";diagnosticError("guardado","calcular huella tras escritura",e)}else{a.libraryDigest=libraryBytesFingerprint(data)}
 }
 return
}

// Remove only failed snapshots proven equivalent to a committed library.
// Never delete an unrelated recovery or a newer draft.
func (a *runningApp) resolvePendingCopies(state libraryState, includeNative bool) {
 dir:=filepath.Join(filepath.Dir(a.libraryPath),"sync_recuperacion")
 entries,err:=os.ReadDir(dir);if err!=nil{return}
 state.MembershipChanges=nil
 expected,_:=json.Marshal(state)
 for _,entry:=range entries{
  if !entry.Type().IsRegular()||!(syncRecoveryName.MatchString(entry.Name())||includeNative&&uiDraftName.MatchString(entry.Name()))||!stringsHasPendingSuffix(entry.Name()){continue}
  path:=filepath.Join(dir,entry.Name())
  data,e:=readRecoveryFile(path);if e!=nil{continue}
  var prior libraryState
  if json.Unmarshal(data,&prior)!=nil{continue}
  prior.MembershipChanges=nil
  // Navigation alone does not distinguish an unresolved edit.
  prior.ActiveID=state.ActiveID
  encoded,e:=json.Marshal(prior)
  if e==nil&&bytes.Equal(encoded,expected){if e=os.Remove(path);e!=nil{diagnosticError("recuperacion","retirar copia resuelta",e)}}
 }
}
func stringsHasPendingSuffix(name string) bool{
 suffix:=pendingRecoverySuffix+".json"
 return len(name)>=len(suffix)&&name[len(name)-len(suffix):]==suffix
}

var libraryBackupIssue atomic.Bool
func recordBackupIssue(operation string,err error){
 if err!=nil{libraryBackupIssue.Store(true);diagnosticError("recuperacion",operation,err)}
}

// ActiveID is purely local navigation. Serialize with the existing JSON
// contract so empty optional collections and nil are equivalent as well.
func sameShareableLibrary(a,b libraryState) bool {
 a.ActiveID=""; b.ActiveID=""
 x,e:=json.Marshal(a); y,f:=json.Marshal(b)
 return e==nil&&f==nil&&bytes.Equal(x,y)
}
