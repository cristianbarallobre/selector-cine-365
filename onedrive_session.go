package main

import (
 "encoding/json"
 "errors"
 "strings"
)

const graphSessionVersion = 1

// graphSessionStore is implemented by Windows user-bound DPAPI storage.
// It is injectable to verify restart and token rotation without secrets.
type graphSessionStore interface {
 load()(graphStoredSession,error)
 save(graphStoredSession)error
 clear()error
}

type graphStoredSession struct {
 Version int `json:"version"`
 ClientID string `json:"client_id"`
 RootID string `json:"onedrive_folder_id"`
 Refresh string `json:"refresh_token"`
}

func validGraphStoredSession(s graphStoredSession)bool{
 return s.Version==graphSessionVersion&&s.ClientID==desktopMicrosoftAppID&&
  strings.TrimSpace(s.RootID)!=""&&len(s.RootID)<=1024&&
  strings.TrimSpace(s.Refresh)!=""&&len(s.Refresh)<=16384
}

func encodeGraphSession(s graphStoredSession)([]byte,error){
 if !validGraphStoredSession(s){return nil,errors.New("Sesión Microsoft incompleta: no se puede conservar")}
 return json.Marshal(s)
}

func decodeGraphSession(data []byte)(graphStoredSession,error){
 var s graphStoredSession
 if len(data)==0||len(data)>32*1024{return s,errors.New("Credencial Microsoft guardada de tamaño inválido")}
 if err:=json.Unmarshal(data,&s);err!=nil{return graphStoredSession{},errors.New("Credencial Microsoft guardada dañada")}
 if !validGraphStoredSession(s){return graphStoredSession{},errors.New("La credencial guardada no pertenece a esta versión de Selector Cine 365")}
 return s,nil
}

func zeroBytes(data []byte){for i:=range data{data[i]=0}}
