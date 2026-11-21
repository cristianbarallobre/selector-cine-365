package main

import (
 "bytes"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "net/http"
 "net/url"
 "strings"
 "sync"
 "time"
)

const desktopMicrosoftAppID = "696458fe-0672-47e5-9fdd-aa0332e9884b"
const microsoftDeviceEndpoint = "https://login.microsoftonline.com/consumers/oauth2/v2.0/"
const desktopGraphEndpoint = "https://graph.microsoft.com/v1.0/"
const desktopGraphScope = "offline_access Files.ReadWrite User.Read"

// desktopGraph stores tokens only in memory. Never write refresh tokens to
// unencrypted files, browser JS, GitHub artifacts, or logging. Reauthentication
// after app restart is intentional until native secure token storage ships.
type desktopGraph struct {
 mu sync.Mutex
 client *http.Client
 authBase,graphBase string
 token,refresh,rootID,lastError,deviceCode,userCode,verifyURL string
 expiry,timeLimit time.Time
 pending bool
 lastCheck time.Time
}

func newDesktopGraph()*desktopGraph {
 return &desktopGraph{
  client:&http.Client{Timeout:35*time.Second},
  authBase:microsoftDeviceEndpoint,graphBase:desktopGraphEndpoint,
 }
}
func (g *desktopGraph) status() map[string]any {
 g.mu.Lock();defer g.mu.Unlock()
 return map[string]any{
  "connected":g.token!=""&&g.rootID!="",
  "pending":g.pending,
  "verification_uri":g.verifyURL,
  "user_code":g.userCode,
  "last_error":g.lastError,
  "provider":"onedrive-direct",
 }
}
func (g *desktopGraph) active()bool{
 g.mu.Lock();defer g.mu.Unlock()
 return g.token!=""&&g.rootID!=""
}
func (g *desktopGraph) pairingID()string{
 g.mu.Lock();defer g.mu.Unlock()
 if g.rootID==""{return ""}
 return "onedrive://item/"+g.rootID
}
func (g *desktopGraph) logout(){
 g.mu.Lock();defer g.mu.Unlock()
 g.token="";g.refresh="";g.rootID="";g.pending=false;g.deviceCode="";g.userCode="";g.verifyURL="";g.lastError=""
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
func (g *desktopGraph) begin() (map[string]any,error) {
 g.mu.Lock()
 if g.pending {g.mu.Unlock();return g.status(),nil}
 g.token="";g.refresh="";g.rootID="";g.lastError=""
 g.mu.Unlock()
 data,err:=g.oauthForm("devicecode",url.Values{"client_id":{desktopMicrosoftAppID},"scope":{desktopGraphScope}})
 if err!=nil{return nil,err}
 device:=stringValue(data,"device_code")
 if device==""{return nil,fmt.Errorf("Microsoft no autorizó el inicio de sesión: %s",stringValue(data,"error_description"))}
 verify:=stringValue(data,"verification_uri")
 code:=stringValue(data,"user_code")
 if !strings.HasPrefix(verify,"https://"){return nil,errors.New("Dirección de validación Microsoft no segura")}
 expires:=numericValue(data,"expires_in")
 if expires<30||expires>3600{return nil,errors.New("Microsoft devolvió una expiración de inicio de sesión inválida")}
 interval:=numericValue(data,"interval")
 if interval<5{interval=5}
 g.mu.Lock()
 g.deviceCode=device;g.userCode=code;g.verifyURL=verify;g.pending=true
 g.timeLimit=time.Now().Add(time.Duration(expires)*time.Second)
 g.mu.Unlock()
 go g.poll(device,interval)
 return g.status(),nil
}
func (g *desktopGraph) poll(device string,interval int){
 for {
  time.Sleep(time.Duration(interval)*time.Second)
  g.mu.Lock()
  active:=g.pending&&g.deviceCode==device&&time.Now().Before(g.timeLimit)
  g.mu.Unlock()
  if !active{return}
  data,err:=g.oauthForm("token",url.Values{
   "grant_type":{"urn:ietf:params:oauth:grant-type:device_code"},
   "client_id":{desktopMicrosoftAppID},"device_code":{device},
  })
  if err!=nil{g.failLogin(device,err);return}
  code:=stringValue(data,"error")
  if code=="authorization_pending"{continue}
  if code=="slow_down"{interval+=5;continue}
  if code!=""{g.failLogin(device,fmt.Errorf("Microsoft no autorizó la conexión: %s",code));return}
  access:=stringValue(data,"access_token")
  if access==""{g.failLogin(device,errors.New("Microsoft no devolvió acceso a OneDrive"));return}
  g.mu.Lock()
  if !g.pending||g.deviceCode!=device{g.mu.Unlock();return}
  g.token=access;g.refresh=stringValue(data,"refresh_token")
  lifetime:=numericValue(data,"expires_in")
  if lifetime<120{lifetime=3600}
  g.expiry=time.Now().Add(time.Duration(lifetime-90)*time.Second)
  g.mu.Unlock()
  root,err:=g.ensureRoot()
  if err!=nil{g.failLogin(device,err);return}
  g.mu.Lock()
  if g.pending&&g.deviceCode==device{
   g.rootID=root
   g.pending=false;g.deviceCode="";g.userCode="";g.verifyURL="";g.lastError=""
  }
  g.mu.Unlock()
  return
 }
}
func (g *desktopGraph) failLogin(device string,err error){
 g.mu.Lock();defer g.mu.Unlock()
 if g.deviceCode!=device{return}
 g.token="";g.refresh="";g.rootID="";g.pending=false
 g.deviceCode="";g.userCode="";g.verifyURL=""
 g.lastError=err.Error()
}
func (g *desktopGraph) accessToken()(string,error){
 g.mu.Lock()
 if g.token==""{g.mu.Unlock();return "",errors.New("Conecta primero tu cuenta Microsoft")}
 if time.Now().Before(g.expiry){t:=g.token;g.mu.Unlock();return t,nil}
 refresh:=g.refresh
 g.mu.Unlock()
 if refresh==""{return "",errors.New("La sesión de Microsoft ha caducado. Vuelve a conectar.")}
 data,err:=g.oauthForm("token",url.Values{
  "grant_type":{"refresh_token"},"client_id":{desktopMicrosoftAppID},
  "refresh_token":{refresh},"scope":{desktopGraphScope},
 })
 if err!=nil{return "",err}
 token:=stringValue(data,"access_token")
 if token==""{return "",errors.New("Microsoft ha caducado la sesión; vuelve a conectar")}
 g.mu.Lock();defer g.mu.Unlock()
 if g.refresh!=refresh{return "",errors.New("La cuenta Microsoft ha cambiado; vuelve a conectar")}
 g.token=token
 if next:=stringValue(data,"refresh_token");next!=""{g.refresh=next}
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
