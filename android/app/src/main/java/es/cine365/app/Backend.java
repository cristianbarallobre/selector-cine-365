package es.cine365.app;

import android.content.Context;
import android.net.Uri;
import org.json.JSONArray;
import org.json.JSONObject;
import java.io.File;
import java.io.FileOutputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.time.OffsetDateTime;
import java.time.LocalDate;

final class Backend {
    static final class ApiResponse {
        int status; String body,cloudHeader="";
        ApiResponse(int s,String b){status=s;body=b;}
        ApiResponse cloud(String s){cloudHeader=s;return this;}
    }
    private final MainActivity activity;
    final File dir;
    final Tmdb tmdb;
    final CloudBackup cloud;
    final OneDriveCloud onedrive;
    final SyncStore sync;
    private final Object filesLock=new Object();
    Backend(MainActivity activity){
        this.activity=activity;
        dir=activity.getFilesDir();
        sync=new SyncStore(dir);
        tmdb=new Tmdb(dir);
        onedrive=new OneDriveCloud(activity);
        cloud=new CloudBackup(activity,dir,sync,onedrive);
        if(!new File(dir,"colecciones.json").exists()){
            try{
                byte[] data;
                try(java.io.InputStream in=activity.getAssets().open("lista_inicial.txt"); java.io.ByteArrayOutputStream out=new java.io.ByteArrayOutputStream()) {
                    byte[] buffer=new byte[8192];int n;while((n=in.read(buffer))!=-1)out.write(buffer,0,n);
                    data=out.toByteArray();
                }
                JSONObject state=new JSONObject().put("version",1).put("active_id","atlas")
                    .put("lists",new JSONArray().put(new JSONObject().put("id","atlas").put("name","Selector Cine 365 · 365")
                    .put("text",new String(data,StandardCharsets.UTF_8))));
                storeState(state);
            }catch(Exception ignored){}
        }
    }
    ApiResponse handle(String path,String method,String body)throws Exception{
        switch(path){
            case "/api/ping":return new ApiResponse(204,"");
            case "/api/quit":return new ApiResponse(204,"");
            case "/api/library":
                if("GET".equals(method))return ok(readOrDefault("colecciones.json","{}"));
                if("POST".equals(method)){
                    if(body.length()>3_000_000)return fail(413,"Archivo demasiado grande");
                    JSONObject state=new JSONObject(body);
                    if(!validLibrary(state))return fail(400,"Colecciones no válidas");
                    synchronized(filesLock){
                        JSONObject previous=new JSONObject(readOrDefault("colecciones.json","{}"));
                        sync.observe(validLibrary(previous)?previous:null,state);
                        storeState(state);
                    }
                    JSONObject s=cloud.backup();
                    return new ApiResponse(204,"").cloud(!s.optBoolean("enabled",false)?"":s.optString("last_error","").isEmpty()?"ok":"error");
                }
                return fail(405,"Método incorrecto");
            case "/api/identities":return ok(tmdb.identities().toString());
            case "/api/identities/resolve":{
                JSONObject req=new JSONObject(body);
                return ok(tmdb.identify(req.optString("title"),req.optInt("year"),req.optBoolean("force")).toString());
            }
            case "/api/platforms/status":return ok(new JSONObject().put("configured",tmdb.hasKey()).toString());
            case "/api/platforms/cache":return ok(tmdb.cachedEntries().toString());
            case "/api/platforms/config":{
                String key=new JSONObject(body).optString("key","").trim();
                if(key.length()<12)return fail(400,"Clave TMDB incorrecta");
                tmdb.configure(key);return ok(new JSONObject().put("configured",true).toString());
            }
            case "/api/platforms/lookup":{
                JSONObject p=new JSONObject(body);
                return ok(tmdb.lookup(p.optString("title"),p.optInt("year"),p.optBoolean("force")).toString());
            }
            case "/api/platforms/match":{
                JSONObject p=new JSONObject(body);
                tmdb.match(p.optString("title"),p.optInt("year"),p.optInt("id"));
                cloud.backup();return ok(new JSONObject().put("saved",true).toString());
            }
            case "/api/onedrive/status":return ok(onedrive.status().toString());
            case "/api/onedrive/connect":return ok(onedrive.begin().toString());
            case "/api/onedrive/poll":return ok(onedrive.poll().toString());
            case "/api/onedrive/folders":return ok(onedrive.folders().toString());
            case "/api/onedrive/create":{
                String name=new JSONObject(body).optString("name","Pruebas_Selector_Cine_365");
                return ok(onedrive.createRoot(name).toString());
            }
            case "/api/onedrive/disconnect":{
                if("onedrive".equals(cloud.status().optString("provider")))cloud.configure("");
                onedrive.logout();
                return ok(onedrive.status().toString());
            }
            case "/api/cloud/status":return ok(cloud.status().toString());
            case "/api/cloud/browse":return ok(new JSONObject().put("folder",activity.selectFolder()).toString());
            case "/api/cloud/config":{
                String uri=new JSONObject(body).optString("folder","");
                cloud.configure(uri);
                return ok(cloud.status().toString());
            }
            case "/api/cloud/backup":return ok(cloud.backup().toString());
            case "/api/cloud/preview":{
                synchronized(filesLock){return ok(cloud.preview(this).toString());}
            }
            case "/api/cloud/sync":{
                JSONObject options=new JSONObject(body.isEmpty()?"{}":body);
                synchronized(filesLock){return ok(cloud.synchronize(this,options.optString("token"),
                    options.optJSONObject("decisions")==null?new JSONObject():options.getJSONObject("decisions")).toString());}
            }
            case "/api/save":return new ApiResponse(204,"");
            default:return fail(404,"Ruta inexistente");
        }
    }
    static ApiResponse ok(String json){return new ApiResponse(200,json);}
    static ApiResponse fail(int code,String msg)throws Exception{return new ApiResponse(code,new JSONObject().put("error",msg).toString());}
    static boolean validLibrary(JSONObject s){
        try {
            if(s.optInt("version")!=1)return false;
            JSONArray lists=s.getJSONArray("lists");
            if(lists.length()<1||lists.length()>40)return false;
            boolean active=false;java.util.HashSet<String> ids=new java.util.HashSet<>();
            for(int i=0;i<lists.length();i++){
                JSONObject x=lists.getJSONObject(i);
                String id=x.getString("id"),title=x.getString("name"),text=x.getString("text");
                if(!id.matches("[A-Za-z0-9_-]{1,60}")||!ids.add(id)||title.trim().isEmpty()||title.length()>100||text.length()>900000)return false;
                if(id.equals(s.optString("active_id")))active=true;
            }
            return active;
        }catch(Exception e){return false;}
    }
    void storeStateImported(JSONObject state)throws Exception{
        if(!validLibrary(state))throw new Exception("Biblioteca remota no válida");
        storeState(state);
    }
    void storeState(JSONObject state)throws Exception{
        String txt=state.toString(2);
        File old=new File(dir,"colecciones.json"), backupDir=new File(dir,"copias_automaticas");
        if(old.exists()){
            backupDir.mkdirs();File snap=new File(backupDir,LocalDate.now()+".json");
            if(!snap.exists())Files.copy(old.toPath(),snap.toPath());
        }
        atomicWrite(old,txt);
        JSONArray lists=state.getJSONArray("lists");
        for(int i=0;i<lists.length();i++){
            JSONObject l=lists.getJSONObject(i);
            if("atlas".equals(l.optString("id")))atomicWrite(new File(dir,"mis_365_peliculas.txt"),l.optString("text"));
        }
    }
    static void atomicWrite(File file,String contents)throws Exception{
        File temp=new File(file.getParentFile(),file.getName()+".tmp");
        try(FileOutputStream f=new FileOutputStream(temp)){f.write(contents.getBytes(StandardCharsets.UTF_8)); f.getFD().sync();}
        if(!temp.renameTo(file)){
            try(FileOutputStream f=new FileOutputStream(file)){f.write(contents.getBytes(StandardCharsets.UTF_8));}
            temp.delete();
        }
    }
    String readOrDefault(String name,String def){
        try{return new String(Files.readAllBytes(new File(dir,name).toPath()),StandardCharsets.UTF_8);}catch(Exception e){return def;}
    }
}
