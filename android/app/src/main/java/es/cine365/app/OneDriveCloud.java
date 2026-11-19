package es.cine365.app;

import android.content.Context;
import android.net.Uri;
import android.security.keystore.KeyGenParameterSpec;
import android.security.keystore.KeyProperties;
import android.util.Base64;
import org.json.JSONArray;
import org.json.JSONObject;
import java.io.ByteArrayOutputStream;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.security.KeyStore;
import java.util.Arrays;
import java.util.Locale;
import javax.crypto.Cipher;
import javax.crypto.KeyGenerator;
import javax.crypto.SecretKey;
import javax.crypto.spec.GCMParameterSpec;

/**
 * Microsoft Graph connector for personal OneDrive accounts.
 * OAuth device-code flow needs only a PUBLIC client ID registered by the publisher.
 * No app secret, passwords or tokens pass through the embedded web UI.
 * Persistent refresh tokens are encrypted with Android Keystore AES-GCM.
 */
final class OneDriveCloud {
    private static final String AUTH="https://login.microsoftonline.com/consumers/oauth2/v2.0/";
    private static final String GRAPH="https://graph.microsoft.com/v1.0/";
    private static final String SCOPE="offline_access Files.ReadWrite User.Read";
    private static final String SECRET_KEY="selector_cine_365_microsoft_tokens_v1";
    private final Context context;
    private final String clientId;
    private String accessToken="",refreshToken="";
    private long expiresAt=0;
    private String deviceCode="",userCode="",verificationUrl="";
    private long deviceExpiresAt=0,nextPollAt=0;
    private int pollInterval=5;

    OneDriveCloud(Context context){
        this.context=context;
        clientId=BuildConfig.MICROSOFT_CLIENT_ID.trim();
        try {JSONObject saved=new JSONObject(decrypt(context.getSharedPreferences("onedrive",Context.MODE_PRIVATE).getString("secrets","")));
            accessToken=saved.optString("access_token");refreshToken=saved.optString("refresh_token");
            expiresAt=saved.optLong("expires_at");
        }catch(Exception ignored){}
    }
    synchronized boolean available(){return clientId.matches("[0-9a-fA-F-]{36}");}
    synchronized boolean connected(){return available() && (!refreshToken.isEmpty()||(!accessToken.isEmpty()&&System.currentTimeMillis()<expiresAt));}
    synchronized JSONObject status()throws Exception{
        return new JSONObject().put("available",available()).put("connected",connected())
            .put("authentication","device_code").put("provider","onedrive");
    }
    synchronized JSONObject begin()throws Exception{
        if(!available())throw new Exception("OneDrive aún no está activado: falta registrar Selector Cine 365 en Microsoft.");
        JSONObject answer=new JSONObject(new String(http(AUTH+"devicecode","POST",
            form("client_id",clientId,"scope",SCOPE),null,"application/x-www-form-urlencoded"),StandardCharsets.UTF_8));
        deviceCode=answer.getString("device_code");userCode=answer.getString("user_code");
        verificationUrl=answer.getString("verification_uri");
        if(!verificationUrl.startsWith("https://microsoft.com/")&&!verificationUrl.startsWith("https://www.microsoft.com/"))
            throw new Exception("Dirección de inicio de sesión no reconocida");
        deviceExpiresAt=System.currentTimeMillis()+answer.optInt("expires_in",900)*1000L;
        pollInterval=Math.max(5,answer.optInt("interval",5));
        nextPollAt=0;
        return new JSONObject().put("code",userCode).put("url",verificationUrl)
            .put("expires_in",answer.optInt("expires_in",900)).put("pending",true);
    }
    synchronized JSONObject poll()throws Exception{
        if(deviceCode.isEmpty())return new JSONObject().put("pending",false).put("connected",connected());
        long now=System.currentTimeMillis();
        if(now>deviceExpiresAt){deviceCode="";throw new Exception("El código de Microsoft ha caducado; vuelve a conectar.");}
        if(now<nextPollAt)return new JSONObject().put("pending",true).put("retry_in",(nextPollAt-now+999)/1000);
        nextPollAt=now+pollInterval*1000L;
        try {
            JSONObject tokens=tokenRequest(form("grant_type","urn:ietf:params:oauth:grant-type:device_code",
                "client_id",clientId,"device_code",deviceCode));
            store(tokens);
            deviceCode="";userCode="";
            return new JSONObject().put("pending",false).put("connected",true);
        }catch(OAuthPending pending){
            if("slow_down".equals(pending.code)){pollInterval+=5;nextPollAt=now+pollInterval*1000L;}
            if("authorization_declined".equals(pending.code)||"expired_token".equals(pending.code)){
                deviceCode="";throw new Exception("Inicio de sesión cancelado o caducado.");
            }
            return new JSONObject().put("pending",true).put("retry_in",pollInterval);
        }
    }
    synchronized void logout(){
        accessToken="";refreshToken="";expiresAt=0;deviceCode="";
        context.getSharedPreferences("onedrive",Context.MODE_PRIVATE).edit().remove("secrets").apply();
    }
    private static final class OAuthPending extends Exception {
        final String code;
        OAuthPending(String code){this.code=code;}
    }
    private JSONObject tokenRequest(byte[] body)throws Exception{
        HttpURLConnection c=null;
        try{
            c=connection(AUTH+"token","POST");
            c.setRequestProperty("Content-Type","application/x-www-form-urlencoded");
            c.setDoOutput(true);
            try(OutputStream out=c.getOutputStream()){out.write(body);}
            int status=c.getResponseCode();
            byte[] bytes=readStream(status>=400?c.getErrorStream():c.getInputStream(),512*1024);
            JSONObject result=new JSONObject(new String(bytes,StandardCharsets.UTF_8));
            if(status>=400){
                String error=result.optString("error");
                if(Arrays.asList("authorization_pending","slow_down","authorization_declined","expired_token").contains(error))throw new OAuthPending(error);
                if("invalid_grant".equals(error))throw new Exception("Tu sesión de OneDrive ha caducado. Conecta de nuevo tu cuenta Microsoft.");
                throw new Exception("Microsoft no autorizó la conexión: "+result.optString("error_description",error));
            }
            return result;
        }finally{if(c!=null)c.disconnect();}
    }
    private void store(JSONObject token)throws Exception{
        accessToken=token.getString("access_token");
        refreshToken=token.optString("refresh_token",refreshToken);
        expiresAt=System.currentTimeMillis()+Math.max(60,token.optLong("expires_in",3600)-120)*1000L;
        JSONObject packed=new JSONObject().put("access_token",accessToken).put("refresh_token",refreshToken).put("expires_at",expiresAt);
        context.getSharedPreferences("onedrive",Context.MODE_PRIVATE).edit().putString("secrets",encrypt(packed.toString())).apply();
    }
    private synchronized String token()throws Exception{
        if(!connected())throw new Exception("Primero conecta tu cuenta de OneDrive.");
        if(System.currentTimeMillis()<expiresAt&&!accessToken.isEmpty())return accessToken;
        if(refreshToken.isEmpty())throw new Exception("Conecta de nuevo OneDrive; no hay sesión válida.");
        store(tokenRequest(form("grant_type","refresh_token","client_id",clientId,"refresh_token",refreshToken,"scope",SCOPE)));
        return accessToken;
    }
    private byte[] graph(String path,String method,byte[] body,String type)throws Exception{
        if(!path.startsWith("me/drive/"))throw new Exception("Ruta de OneDrive no autorizada");
        String bearer=token();
        String url=GRAPH+path;
        HttpURLConnection c=null;
        try{
            c=connection(url,method);c.setInstanceFollowRedirects(false);
            c.setRequestProperty("Authorization","Bearer "+bearer);
            if(body!=null){c.setDoOutput(true);c.setRequestProperty("Content-Type",type);
                try(OutputStream out=c.getOutputStream()){out.write(body);}}
            int status=c.getResponseCode();
            if((status==302||status==301)&&"GET".equals(method)&&path.endsWith("/content")){
                String redirected=c.getHeaderField("Location");
                if(redirected==null||!redirected.startsWith("https://"))
                    throw new Exception("OneDrive devolvió una descarga no segura");
                // Never send the Graph bearer token to a download CDN.
                return http(redirected,"GET",null,null,null);
            }
            byte[] result=readStream(status>=400?c.getErrorStream():c.getInputStream(),9*1024*1024);
            if(status>=400){
                String msg=new String(result,StandardCharsets.UTF_8);
                try {msg=new JSONObject(msg).optJSONObject("error").optString("message",msg);}catch(Exception ignored){}
                if(status==401)throw new Exception("OneDrive necesita volver a iniciar sesión.");
                if(status==429)throw new Exception("OneDrive limita temporalmente las solicitudes. Vuelve a intentarlo más tarde.");
                throw new Exception("Error de OneDrive ("+status+"): "+msg);
            }
            return result;
        }finally{if(c!=null)c.disconnect();}
    }
    synchronized JSONObject folders()throws Exception{
        JSONArray array=new JSONArray();
        String path="me/drive/root/children?$select=id,name,folder&$top=200";
        for(int page=0;page<5;page++){
            JSONObject result=new JSONObject(new String(graph(path,"GET",null,null),StandardCharsets.UTF_8));
            JSONArray list=result.optJSONArray("value");if(list==null)break;
            for(int i=0;i<list.length();i++){
                JSONObject item=list.getJSONObject(i);
                if(item.has("folder"))array.put(new JSONObject().put("id",item.getString("id")).put("name",item.getString("name")));
            }
            String next=result.optString("@odata.nextLink","");
            if(next.isEmpty())break;
            if(!next.startsWith(GRAPH+"me/drive/"))throw new Exception("Paginación de OneDrive no reconocida");
            path=next.substring(GRAPH.length());
        }
        return new JSONObject().put("folders",array);
    }
    synchronized JSONObject createRoot(String name)throws Exception{
        if(!name.matches("[^<>:\"/\\\\|?*]{1,70}")||name.endsWith(".")||name.trim().isEmpty())
            throw new Exception("Nombre de carpeta inválido");
        JSONArray list=folders().getJSONArray("folders");
        for(int i=0;i<list.length();i++){
            JSONObject folder=list.getJSONObject(i);
            if(name.equalsIgnoreCase(folder.getString("name")))return folder;
        }
        JSONObject body=new JSONObject().put("name",name).put("folder",new JSONObject())
            .put("@microsoft.graph.conflictBehavior","fail");
        return toJson(graph("me/drive/root/children","POST",body.toString().getBytes(StandardCharsets.UTF_8),"application/json"));
    }
    synchronized JSONObject find(String parentId,String name)throws Exception{
        JSONArray items=allChildren(parentId);
        for(int i=0;i<items.length();i++){JSONObject item=items.getJSONObject(i);
            if(name.equals(item.optString("name")))return item;}
        return null;
    }
    synchronized JSONArray allChildren(String parentId)throws Exception{
        JSONArray all=new JSONArray();
        String path="me/drive/items/"+Uri.encode(parentId)+"/children?$select=id,name,folder,file&$top=200";
        for(int page=0;page<25;page++){
            JSONObject result=toJson(graph(path,"GET",null,null));
            JSONArray list=result.optJSONArray("value");if(list==null)break;
            for(int i=0;i<list.length();i++)all.put(list.getJSONObject(i));
            String next=result.optString("@odata.nextLink","");
            if(next.isEmpty())break;
            if(!next.startsWith(GRAPH+"me/drive/"))throw new Exception("Paginación no reconocida de Microsoft Graph");
            path=next.substring(GRAPH.length());
        }
        return all;
    }
    synchronized JSONObject newDirectory(String parentId,String name)throws Exception{
        JSONObject existing=find(parentId,name);
        if(existing!=null){if(!existing.has("folder"))throw new Exception("Existe un archivo con el nombre «"+name+"»");return existing;}
        JSONObject body=new JSONObject().put("name",name).put("folder",new JSONObject())
            .put("@microsoft.graph.conflictBehavior","fail");
        return toJson(graph("me/drive/items/"+Uri.encode(parentId)+"/children","POST",body.toString().getBytes(StandardCharsets.UTF_8),"application/json"));
    }
    synchronized byte[] read(String itemId)throws Exception {
        return graph("me/drive/items/"+Uri.encode(itemId)+"/content","GET",null,null);
    }
    synchronized void write(String parentId,String name,byte[] data,boolean overwrite)throws Exception{
        JSONObject existing=find(parentId,name);
        if(existing!=null&&!overwrite)return;
        if(existing!=null&&existing.has("folder"))throw new Exception("La ruta es una carpeta y no un archivo: "+name);
        // Only the caller's device snapshot (or unique daily rescue file) may be written.
        graph("me/drive/items/"+Uri.encode(parentId)+":/"+Uri.encode(name)+":/content","PUT",data,"application/json");
    }
    private static JSONObject toJson(byte[] b)throws Exception{return new JSONObject(new String(b,StandardCharsets.UTF_8));}
    private static byte[] form(String...values)throws Exception{
        StringBuilder body=new StringBuilder();
        for(int i=0;i<values.length;i+=2){if(i>0)body.append('&');body.append(URLEncoder.encode(values[i],"UTF-8")).append('=').append(URLEncoder.encode(values[i+1],"UTF-8"));}
        return body.toString().getBytes(StandardCharsets.UTF_8);
    }
    private static HttpURLConnection connection(String url,String method)throws Exception{
        HttpURLConnection c=(HttpURLConnection)new URL(url).openConnection();
        c.setRequestMethod(method);c.setConnectTimeout(15000);c.setReadTimeout(30000);
        c.setRequestProperty("Accept","application/json");c.setUseCaches(false);
        return c;
    }
    private static byte[] http(String url,String method,byte[] body,String token,String type)throws Exception{
        HttpURLConnection c=null;
        try{
            c=connection(url,method);
            if(token!=null)c.setRequestProperty("Authorization","Bearer "+token);
            if(body!=null){c.setDoOutput(true);if(type!=null)c.setRequestProperty("Content-Type",type);try(OutputStream out=c.getOutputStream()){out.write(body);}}
            int status=c.getResponseCode();
            byte[] b=readStream(status>=400?c.getErrorStream():c.getInputStream(),9*1024*1024);
            if(status>=400)throw new Exception("Microsoft respondió "+status+": "+new String(b,StandardCharsets.UTF_8));
            return b;
        }finally{if(c!=null)c.disconnect();}
    }
    private static byte[] readStream(InputStream in,int max)throws Exception{
        if(in==null)throw new Exception("Respuesta vacía del servidor");
        try(InputStream stream=in;ByteArrayOutputStream out=new ByteArrayOutputStream()){
            byte[] b=new byte[8192];int n;
            while((n=stream.read(b))!=-1){out.write(b,0,n);if(out.size()>max)throw new Exception("El archivo supera el tamaño permitido");}
            return out.toByteArray();
        }
    }
    private SecretKey key()throws Exception {
        KeyStore store=KeyStore.getInstance("AndroidKeyStore");store.load(null);
        if(!store.containsAlias(SECRET_KEY)){
            KeyGenerator generator=KeyGenerator.getInstance(KeyProperties.KEY_ALGORITHM_AES,"AndroidKeyStore");
            generator.init(new KeyGenParameterSpec.Builder(SECRET_KEY,KeyProperties.PURPOSE_ENCRYPT|KeyProperties.PURPOSE_DECRYPT)
                .setBlockModes(KeyProperties.BLOCK_MODE_GCM).setEncryptionPaddings(KeyProperties.ENCRYPTION_PADDING_NONE).build());
            generator.generateKey();
        }
        return ((KeyStore.SecretKeyEntry)store.getEntry(SECRET_KEY,null)).getSecretKey();
    }
    private String encrypt(String plain)throws Exception{
        Cipher c=Cipher.getInstance("AES/GCM/NoPadding");c.init(Cipher.ENCRYPT_MODE,key());
        byte[] iv=c.getIV(),encrypted=c.doFinal(plain.getBytes(StandardCharsets.UTF_8));
        byte[] both=new byte[iv.length+encrypted.length];System.arraycopy(iv,0,both,0,iv.length);
        System.arraycopy(encrypted,0,both,iv.length,encrypted.length);
        return Base64.encodeToString(both,Base64.NO_WRAP);
    }
    private String decrypt(String encoded)throws Exception{
        if(encoded.isEmpty())return "{}";
        byte[] all=Base64.decode(encoded,Base64.DEFAULT);
        if(all.length<29)throw new Exception("Credenciales dañadas");
        Cipher c=Cipher.getInstance("AES/GCM/NoPadding");c.init(Cipher.DECRYPT_MODE,key(),new GCMParameterSpec(128,all,0,12));
        return new String(c.doFinal(all,12,all.length-12),StandardCharsets.UTF_8);
    }
}
