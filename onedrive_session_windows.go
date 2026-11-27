//go:build windows

package main

import("os";"path/filepath")

func protectedSessionStore(dir string) graphSessionStore {
 return &windowsGraphSessionStore{path:filepath.Join(dir,"onedrive_sesion_dpapi.bin")}
}

type windowsGraphSessionStore struct{path string}
func (s *windowsGraphSessionStore) load()(graphStoredSession,error){
 bytes,err:=os.ReadFile(s.path)
 if err!=nil{return graphStoredSession{},err}
 raw,err:=dpapiOpen(bytes)
 if err!=nil{return graphStoredSession{},err}
 defer zeroBytes(raw)
 return decodeGraphSession(raw)
}
func (s *windowsGraphSessionStore) save(session graphStoredSession)error{
 raw,err:=encodeGraphSession(session)
 if err!=nil{return err}
 defer zeroBytes(raw)
 protected,err:=dpapiSeal(raw)
 if err!=nil{return err}
 defer zeroBytes(protected)
 return cloudWriteFile(s.path,protected)
}
func (s *windowsGraphSessionStore) clear()error{
 err:=os.Remove(s.path)
 if os.IsNotExist(err){return nil}
 return err
}
