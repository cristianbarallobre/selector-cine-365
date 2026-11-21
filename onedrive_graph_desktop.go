package main

import (
 "bytes"
 "crypto/rand"
 "crypto/sha256"
 "crypto/subtle"
 "encoding/base64"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "net/http"
 "net/url"
 "os"
 "strings"
 "sync"
 "time"
)

const desktopMicrosoftAppID = "696458fe-0672-47e5-9fdd-aa0332e9884b"
const microsoftOAuthEndpoint = "https://login.microsoftonline.com/consumers/oauth2/v2.0/"
const desktopGraphEndpoint = "https://graph.microsoft.com/v1.0/"
const desktopGraphScope = "offline_access Files.ReadWrite User.Read"

// Windows protects the refresh token using current-user DPAPI. Access tokens
// remain in RAM; no token is exposed to browser JS, logs or GitHub artifacts.
// Other platforms keep the memory-only behavior until native secure storage is added.
type desktopGraph struct {
 mu sync.Mutex
 tokenMu sync.Mutex
 publishMu sync.Mutex
 client *http.Client
 authBase,graphBase string
 token,refresh,rootID,lastError,loginURL,pendingState,pendingVerifier,redirectURI string
 expiry,timeLimit time.Time
 pending,callbackUsed bool
 generation uint64
 lastCheck time.Time
 sessionStore graphSessionStore
 sessionSaved,sessionRestored bool
 persistError string
}

func newDesktopGraph(configDirectory ...string)*desktopGraph {
 g:=&desktopGraph{
  client:&http.Client{Timeout:35*time.Second,CheckRedirect:func(req *http.Request,via []*http.Request)error{
   if req.URL.Scheme!="https"{return errors.New("Redirección de OneDrive no segura")}
   if len(via)>5{return errors.New("Demasiadas redirecciones de OneDrive")}
   if req.URL.Hostname()!="graph.microsoft.com"{req.Header.Del("Authorization")}
   return nil
  }},
  authBase:microsoftOAuthEndpoint,graphBase:desktopGraphEndpoint,
 }
 if len(configDirectory)>0&&configDirectory[0]!=""{
  g.sessionStore=windowsProtectedSessionStore(configDirectory[0])
  g.restoreSession()
 }
 return g
}

// Restore the encrypted refresh token without requesting the Internet at startup.
// The first Graph request refreshes it, so opening the app also works offline.
func (g *desktopGraph) restoreSession(){
 if g.sessionStore==nil{return}
 saved,err:=g.sessionStore.load()
 if errors.Is(err,os.ErrNotExist){return}
 if err!=nil{
  g.persistError="No se pudo recuperar la autorización Microsoft guardada. Conecta de nuevo para restaurarla."
  return
 }
 if !validGraphStoredSession(saved){
  g.persistError="La autorización guardada no es válida. Conecta Microsoft de nuevo."
  return
 }
 g.rootID=saved.RootID
 g.refresh=saved.Refresh
 g.token=""
 g.sessionSaved=true
 g.sessionRestored=true
}

// Called under g.mu. Never store the short-lived access token.
func (g *desktopGraph) rememberSessionLocked(){
 if g.sessionStore==nil{return}
 saved:=graphStoredSession{
  Version:graphSessionVersion,ClientID:desktopMicrosoftAppID,
  RootID:g.rootID,Refresh:g.refresh,
 }
 if err:=g.sessionStore.save(saved);err!=nil{
  g.sessionSaved=false
  g.persistError="No se pudo guardar la sesión Microsoft protegida en Windows. La conexión puede requerir iniciar sesión al volver a abrir."
  return
 }
 g.sessionSaved=true
 g.persistError=""
}

func (g *desktopGraph) status() map[string]any {
 g.mu.Lock();defer g.mu.Unlock()
 return map[string]any{
  "connected":(g.token!=""||g.refresh!="")&&g.rootID!="",
  "pending":g.pending,
  "session_remembered":g.sessionSaved,
  "session_restored":g.sessionRestored,
  "secure_storage":g.sessionStore!=nil,
  "persistence_warning":g.persistError,
  "login_url":g.loginURL,
  "last_error":g.lastError,
  "provider":"onedrive-direct",
 }
}
func (g *desktopGraph) active()bool{
 g.mu.Lock();defer g.mu.Unlock()
 return (g.token!=""||g.refresh!="")&&g.rootID!=""
}
func (g *desktopGraph) pairingID()string{
 g.mu.Lock();defer g.mu.Unlock()
 if g.rootID==""{return ""}
 return "onedrive://item/"+g.rootID
}
func (g *desktopGraph) logout()error{
 // Serialize logout with token refresh. A failed deletion must not be reported
 // as a successful disconnect, otherwise the old account could reappear.
 g.tokenMu.Lock();defer g.tokenMu.Unlock()
 g.mu.Lock();defer g.mu.Unlock()
 if g.sessionStore!=nil{
  if err:=g.sessionStore.clear();err!=nil{
   return errors.New("No se pudo eliminar la sesión cifrada. Cierra otras instancias y vuelve a pulsar Desconectar.")
  }
 }
 g.generation++;g.token="";g.refresh="";g.rootID="";g.pending=false;g.callbackUsed=false;g.pendingState="";g.pendingVerifier="";g.loginURL="";g.lastError=""
 g.sessionSaved=false;g.sessionRestored=false;g.persistError=""
 return nil
}
func (g *desktopGraph) oauthForm(endpoint string,values url.Values)(map[string]any,error){
 req,err:=http.NewRequest("POST",g.authBase+endpoint,strings.NewReader(values.Encode()))
 if err!=nil{return nil,err}
 req.Header.Set("Content-Type","application/x-www-form-urlencoded")
 resp,err:=g.client.Do(req);if err!=nil{return nil,err}
 defer resp.Body.Close()
 b,err:=io.ReadAll(io.LimitReader(resp.Body,1<<20));if err!=nil{return nil,err}
 var data map[string]any
 if err=json.Unmarshal(b,&data);err!=nil{return nil,fmt.Errorf("Microsoft devolvió una respuesta inválida: %w",err)}
 if resp.StatusCode>=500{return nil,fmt.Errorf("Microsoft temporalmente no disponible (%d)",resp.StatusCode)}
 return data,nil
}
func stringValue(m map[string]any,k string)string{s,_:=m[k].(string);return s}
func numericValue(m map[string]any,k string)int{v,_:=m[k].(float64);return int(v)}
func (g *desktopGraph) setRedirect(uri string){
 g.mu.Lock();g.redirectURI=uri;g.mu.Unlock()
}
func (g *desktopGraph) callbackHost()string{
 g.mu.Lock();uri:=g.redirectURI;g.mu.Unlock()
 u,err:=url.Parse(uri);if err!=nil{return ""}
 return u.Host
}
func randomPKCE()(string,error){
 bytes:=make([]byte,32)
 if _,err:=rand.Read(bytes);err!=nil{return "",err}
 return base64.RawURLEncoding.EncodeToString(bytes),nil
}
func (g *desktopGraph) begin()(map[string]any,error){
 g.mu.Lock()
 if g.pending{g.mu.Unlock();return g.status(),nil}
 redirect:=g.redirectURI
 g.mu.Unlock()
 if redirect==""{return nil,errors.New("La dirección local de retorno de Microsoft no está preparada")}
 verifier,err:=randomPKCE();if err!=nil{return nil,err}
 state,err:=randomPKCE();if err!=nil{return nil,err}
 challengeBytes:=sha256.Sum256([]byte(verifier))
 challenge:=base64.RawURLEncoding.EncodeToString(challengeBytes[:])
 query:=url.Values{
  "client_id":{desktopMicrosoftAppID},
  "response_type":{"code"},
  "redirect_uri":{redirect},
  "response_mode":{"query"},
  "scope":{desktopGraphScope},
  "code_challenge":{challenge},
  "code_challenge_method":{"S256"},
  "state":{state},
  "prompt":{"select_account"},
 }
 address:=g.authBase+"authorize?"+query.Encode()
 g.mu.Lock()
 g.generation++
 g.token="";g.refresh="";g.rootID="";g.lastError=""
 g.sessionSaved=false;g.sessionRestored=false;g.persistError=""
 g.loginURL=address;g.pendingState=state;g.pendingVerifier=verifier;g.pending=true;g.callbackUsed=false
 g.timeLimit=time.Now().Add(10*time.Minute)
 g.mu.Unlock()
 return g.status(),nil
}
func (g *desktopGraph) loginFailure(gen uint64,err error){
 g.mu.Lock();defer g.mu.Unlock()
 if g.generation!=gen{return}
 g.token="";g.refresh="";g.rootID="";g.pending=false;g.callbackUsed=false
 g.pendingState="";g.pendingVerifier="";g.loginURL=""
 g.lastError=err.Error()
}
// finish is called only on a local loopback callback with an unguessable state.
// The authorization code and verifier are never passed to browser JavaScript.
func (g *desktopGraph) finish(state,code,oauthError string)error{
 g.mu.Lock()
 if !g.pending||g.callbackUsed||time.Now().After(g.timeLimit){
  g.mu.Unlock();return errors.New("La sesión de Microsoft ha caducado o ya fue utilizada")
 }
 if subtle.ConstantTimeCompare([]byte(state),[]byte(g.pendingState))!=1{
  g.mu.Unlock();return errors.New("La respuesta de Microsoft no coincide con esta conexión")
 }
 gen:=g.generation
 verifier,redirect:=g.pendingVerifier,g.redirectURI
 g.callbackUsed=true
 g.pendingState="";g.pendingVerifier="";g.loginURL=""
 g.mu.Unlock()
 if oauthError!=""{
  err:=fmt.Errorf("Microsoft canceló el acceso: %s",oauthError)
  g.loginFailure(gen,err);return err
 }
 if code==""||len(code)>8192{
  err:=errors.New("Microsoft no devolvió un código de autorización válido")
  g.loginFailure(gen,err);return err
 }
 data,err:=g.oauthForm("token",url.Values{
  "grant_type":{"authorization_code"},"client_id":{desktopMicrosoftAppID},
  "code":{code},"redirect_uri":{redirect},"code_verifier":{verifier},
  "scope":{desktopGraphScope},
 })
 if err!=nil{g.loginFailure(gen,err);return err}
 access:=stringValue(data,"access_token")
 if access==""{
  err=fmt.Errorf("Microsoft rechazó el inicio de sesión: %s",stringValue(data,"error_description"))
  g.loginFailure(gen,err);return err
 }
 g.mu.Lock()
 if g.generation!=gen{g.mu.Unlock();return errors.New("La conexión fue cancelada")}
 g.token=access;g.refresh=stringValue(data,"refresh_token")
 lifetime:=numericValue(data,"expires_in");if lifetime<120{lifetime=3600}
 g.expiry=time.Now().Add(time.Duration(lifetime-90)*time.Second)
 g.mu.Unlock()
 root,err:=g.ensureRoot()
 if err!=nil{g.loginFailure(gen,err);return err}
 g.mu.Lock();defer g.mu.Unlock()
 if g.generation!=gen{return errors.New("La conexión fue cancelada")}
 g.rootID=root;g.pending=false;g.callbackUsed=false
 g.pendingState="";g.pendingVerifier="";g.loginURL="";g.lastError=""
 g.sessionRestored=false
 g.rememberSessionLocked()
 return nil
}
func (g *desktopGraph) accessToken()(string,error){
 g.tokenMu.Lock();defer g.tokenMu.Unlock()
 g.mu.Lock()
 if g.token==""&&g.refresh==""{g.mu.Unlock();return "",errors.New("Conecta primero tu cuenta Microsoft")}
 if g.token!=""&&time.Now().Before(g.expiry){t:=g.token;g.mu.Unlock();return t,nil}
 refresh:=g.refresh
 g.mu.Unlock()
 if refresh==""{return "",errors.New("La sesión de Microsoft ha caducado. Vuelve a conectar.")}
 data,err:=g.oauthForm("token",url.Values{
  "grant_type":{"refresh_token"},"client_id":{desktopMicrosoftAppID},
  "refresh_token":{refresh},"scope":{desktopGraphScope},
 })
 if err!=nil{return "",err}
 token:=stringValue(data,"access_token")
 if token==""{
  // Only an explicit revocation should forget the saved session. Transient
  // network failures are already returned by oauthForm and remain retryable.
  if code:=stringValue(data,"error");code=="invalid_grant"||code=="interaction_required"{
   g.mu.Lock()
   if g.refresh==refresh{
    if g.sessionStore!=nil{
     if err:=g.sessionStore.clear();err!=nil{
      g.persistError="No se pudo borrar la sesión Microsoft caducada. Desconecta y vuelve a conectar."
     }
    }
    g.token="";g.refresh="";g.rootID="";g.sessionSaved=false;g.sessionRestored=false
    g.lastError="La autorización Microsoft ha caducado o ha sido revocada. Vuelve a conectar."
   }
   g.mu.Unlock()
  }
  return "",errors.New("La autorización Microsoft ha caducado o ha sido revocada. Vuelve a conectar")
 }
 g.mu.Lock();defer g.mu.Unlock()
 if g.refresh!=refresh{return "",errors.New("La cuenta Microsoft ha cambiado; vuelve a conectar")}
 g.token=token
 if next:=stringValue(data,"refresh_token");next!=""{
  g.refresh=next
  g.rememberSessionLocked()
 }
 ttl:=numericValue(data,"expires_in");if ttl<120{ttl=3600}
 g.expiry=time.Now().Add(time.Duration(ttl-90)*time.Second)
 return g.token,nil
}
func (g *desktopGraph) request(path,method string,body []byte,contentType string)([]byte,error){
 if !strings.HasPrefix(path,"me/drive/"){return nil,errors.New("Ruta Microsoft Graph no autorizada")}
 token,err:=g.accessToken();if err!=nil{return nil,err}
 req,err:=http.NewRequest(method,g.graphBase+path,bytes.NewReader(body));if err!=nil{return nil,err}
 req.Header.Set("Authorization","Bearer "+token)
 req.Header.Set("Accept","application/json")
 if contentType!=""{req.Header.Set("Content-Type",contentType)}
 resp,err:=g.client.Do(req);if err!=nil{return nil,err}
 defer resp.Body.Close()
 if resp.StatusCode==401{return nil,errors.New("La sesión Microsoft necesita volver a conectarse")}
 if resp.StatusCode==429{return nil,errors.New("OneDrive ha limitado las solicitudes; inténtalo más tarde")}
 raw,err:=io.ReadAll(io.LimitReader(resp.Body,8*1024*1024+1));if err!=nil{return nil,err}
 if len(raw)>8*1024*1024{return nil,errors.New("Archivo de OneDrive demasiado grande")}
 if resp.StatusCode<200||resp.StatusCode>=300{return nil,fmt.Errorf("OneDrive respondió %d: %.180s",resp.StatusCode,raw)}
 return raw,nil
}
type graphDriveItem struct {
 ID string `json:"id"`
 Name string `json:"name"`
 Folder *json.RawMessage `json:"folder,omitempty"`
 File *json.RawMessage `json:"file,omitempty"`
}
func (g *desktopGraph) list(path string)([]graphDriveItem,error){
 list:=[]graphDriveItem{}
 for n:=0;n<20;n++ {
  b,err:=g.request(path,"GET",nil,"");if err!=nil{return nil,err}
  var response struct {
   Value []graphDriveItem `json:"value"`
   Next string `json:"@odata.nextLink"`
  }
  if err=json.Unmarshal(b,&response);err!=nil{return nil,err}
  if response.Value==nil{return nil,errors.New("OneDrive devolvió un listado incompleto")}
  list=append(list,response.Value...)
  if response.Next==""{return list,nil}
  if !strings.HasPrefix(response.Next,g.graphBase+"me/drive/"){return nil,errors.New("Paginación Microsoft Graph no reconocida")}
  path=strings.TrimPrefix(response.Next,g.graphBase)
 }
 return nil,errors.New("OneDrive devolvió demasiadas páginas; no se sincronizarán datos parciales")
}
func (g *desktopGraph) ensureRoot()(string,error){
 items,err:=g.list("me/drive/root/children?$select=id,name,folder&$top=200");if err!=nil{return "",err}
 for _,item:=range items{
  if strings.EqualFold(item.Name,"Selector_Cine_365"){
   if item.Folder==nil{return "",errors.New("Selector_Cine_365 existe en OneDrive, pero es un archivo")}
   return item.ID,nil
  }
 }
 b,_:=json.Marshal(map[string]any{"name":"Selector_Cine_365","folder":map[string]any{},"@microsoft.graph.conflictBehavior":"fail"})
 result,err:=g.request("me/drive/root/children","POST",b,"application/json")
 if err!=nil{return "",err}
 var item graphDriveItem
 if err=json.Unmarshal(result,&item);err!=nil{return "",err}
 if item.ID==""||item.Folder==nil{return "",errors.New("OneDrive no confirmó la creación de carpeta")}
 return item.ID,nil
}
func (g *desktopGraph) child(parent,name string)(graphDriveItem,bool,error){
 children,err:=g.list("me/drive/items/"+url.PathEscape(parent)+"/children?$select=id,name,folder,file&$top=200")
 if err!=nil{return graphDriveItem{},false,err}
 for _,item:=range children{if item.Name==name{return item,true,nil}}
 return graphDriveItem{},false,nil
}
func (g *desktopGraph) ensureDir(parent,name string)(string,error){
 item,exists,err:=g.child(parent,name);if err!=nil{return "",err}
 if exists {
  if item.Folder==nil{return "",errors.New("Existe un archivo con el nombre de una carpeta de sincronización")}
  return item.ID,nil
 }
 body,_:=json.Marshal(map[string]any{"name":name,"folder":map[string]any{},"@microsoft.graph.conflictBehavior":"fail"})
 data,err:=g.request("me/drive/items/"+url.PathEscape(parent)+"/children","POST",body,"application/json")
 if err!=nil{return "",err}
 if err=json.Unmarshal(data,&item);err!=nil{return "",err}
 if item.ID==""||item.Folder==nil{return "",errors.New("OneDrive no creó el directorio")}
 return item.ID,nil
}
func (g *desktopGraph) readItem(id string)([]byte,error){
 return g.request("me/drive/items/"+url.PathEscape(id)+"/content","GET",nil,"")
}
func (g *desktopGraph) writeItem(parent,name string,bytes []byte)error{
 if len(bytes)>8*1024*1024{return errors.New("La copia supera el tamaño permitido")}
 _,err:=g.request("me/drive/items/"+url.PathEscape(parent)+":/"+url.PathEscape(name)+":/content","PUT",bytes,"application/json")
 return err
}
