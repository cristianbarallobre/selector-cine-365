package es.cine365.app;

import android.content.Context;
import android.content.UriPermission;
import android.database.Cursor;
import android.net.Uri;
import android.provider.DocumentsContract;
import org.json.JSONObject;
import org.json.JSONArray;
import java.security.MessageDigest;
import java.util.Collections;
import java.util.HashSet;
import java.util.Set;
import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.InputStream;
import java.io.OutputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.time.LocalDate;
import java.time.OffsetDateTime;
import java.util.ArrayList;
import java.util.List;

/** Android SAF mirror. Does not overwrite Windows's actual/colecciones.json. */
final class CloudBackup {
    private static final String ROOT="Selector_Cine_365_Copias";
    private final MainActivity activity;
    private final File appDir;
    private final SyncStore sync;
    private final OneDriveCloud oneDrive;
    private String folder="",lastBackup="",lastSync="",lastError="";
    private boolean activated=false,autoNeedsReview=false;
    private String previewToken="";

    CloudBackup(MainActivity a,File files,SyncStore sync,OneDriveCloud oneDrive){
        this.activity=a;this.appDir=files;this.sync=sync;this.oneDrive=oneDrive;
        folder=a.getPreferences(Context.MODE_PRIVATE).getString("backup_uri","");
        lastBackup=a.getPreferences(Context.MODE_PRIVATE).getString("last_backup","");
        lastSync=a.getPreferences(Context.MODE_PRIVATE).getString("last_sync","");
        activated=a.getPreferences(Context.MODE_PRIVATE).getBoolean("sync_activated_v2",false);
    }
    synchronized void configure(String uriValue)throws Exception{
        uriValue=uriValue.trim();
        if(!uriValue.isEmpty()){
            Uri uri=Uri.parse(uriValue);
            if("onedrive".equals(uri.getScheme())){
                if(!"item".equals(uri.getAuthority())||uri.getLastPathSegment()==null||uri.getLastPathSegment().isEmpty())
                    throw new Exception("Selecciona una carpeta de OneDrive desde la lista.");
                if(!oneDrive.connected())throw new Exception("Primero conecta tu cuenta Microsoft.");
            }else if(!"content".equals(uri.getScheme()))throw new Exception("Selecciona una carpeta válida.");
            boolean read=false,write=false;
            if("content".equals(uri.getScheme())){
                for(UriPermission permission:activity.getContentResolver().getPersistedUriPermissions()){
                    if(uri.equals(permission.getUri())){read=permission.isReadPermission();write=permission.isWritePermission();}
                }
                if(!read||!write)throw new Exception("La carpeta elegida necesita permiso de lectura y escritura");
            }
        }
        if(!folder.equals(uriValue)){
            activated=false;autoNeedsReview=false;previewToken="";
            activity.getPreferences(Context.MODE_PRIVATE).edit().putBoolean("sync_activated_v2",false).apply();
        }
        folder=uriValue;lastError="";
        activity.getPreferences(Context.MODE_PRIVATE).edit().putString("backup_uri",folder).apply();
    }
    synchronized JSONObject status()throws Exception{
        return new JSONObject().put("enabled",!folder.isEmpty()).put("activated",activated).put("needs_review",!folder.isEmpty()&&(!activated||autoNeedsReview)).put("folder",folder).put("target",folder)
            .put("last_backup",lastBackup).put("last_sync",lastSync).put("last_error",lastError)
            .put("provider",folder.startsWith("onedrive://")?"onedrive":folder.isEmpty()?"none":"folder")
            .put("auto_sync",!folder.isEmpty()&&activated);
    }
    private Uri tree(){return Uri.parse(folder);}
    private boolean graphMode(){return folder.startsWith("onedrive://");}
    private static Uri driveUri(String id){return Uri.parse("onedrive://item/"+Uri.encode(id));}
    private static String driveId(Uri uri){return uri.getLastPathSegment();}
    private Uri document(Uri child){
        if(graphMode())return child;
        Uri tree=tree();if(child.equals(tree))return DocumentsContract.buildDocumentUriUsingTree(tree,DocumentsContract.getTreeDocumentId(tree));
        return child;
    }
    private Uri findChild(Uri parent,String name)throws Exception{
        if(graphMode()){
            JSONObject item=oneDrive.find(driveId(parent),name);
            return item==null?null:driveUri(item.getString("id"));
        }
        Uri root=tree();Uri item=document(parent);
        String id=DocumentsContract.getDocumentId(item);
        Uri children=DocumentsContract.buildChildDocumentsUriUsingTree(root,id);
        try(Cursor cursor=activity.getContentResolver().query(children,
            new String[]{DocumentsContract.Document.COLUMN_DOCUMENT_ID,DocumentsContract.Document.COLUMN_DISPLAY_NAME},null,null,null)){
            if(cursor==null)throw new Exception("El proveedor de nube no permite consultar esta carpeta");
            while(cursor.moveToNext())if(name.equals(cursor.getString(1)))
                return DocumentsContract.buildDocumentUriUsingTree(root,cursor.getString(0));
        }
        return null;
    }
    private List<Uri> children(Uri parent)throws Exception{
        if(graphMode()){
            List<Uri> files=new ArrayList<>();
            JSONArray items=oneDrive.allChildren(driveId(parent));
            for(int i=0;i<items.length();i++){
                JSONObject item=items.getJSONObject(i);
                if(item.has("file")&&item.optString("name").toLowerCase(java.util.Locale.ROOT).endsWith(".json"))
                    files.add(driveUri(item.getString("id")));
            }
            return files;
        }
        Uri root=tree();Uri item=document(parent);List<Uri> list=new ArrayList<>();
        Uri query=DocumentsContract.buildChildDocumentsUriUsingTree(root,DocumentsContract.getDocumentId(item));
        try(Cursor cursor=activity.getContentResolver().query(query,
                new String[]{DocumentsContract.Document.COLUMN_DOCUMENT_ID,DocumentsContract.Document.COLUMN_DISPLAY_NAME},null,null,null)){
            if(cursor==null)throw new Exception("No se puede explorar la carpeta seleccionada; comprueba permisos y conexión");
            while(cursor.moveToNext()){
                String name=cursor.getString(1);
                if(name!=null&&name.toLowerCase(java.util.Locale.ROOT).endsWith(".json"))
                    list.add(DocumentsContract.buildDocumentUriUsingTree(root,cursor.getString(0)));
            }
        }
        return list;
    }
    private Uri ensureDirectory(Uri parent,String name)throws Exception{
        if(graphMode())return driveUri(oneDrive.newDirectory(driveId(parent),name).getString("id"));
        Uri item=findChild(parent,name);if(item!=null)return item;
        Uri created=DocumentsContract.createDocument(activity.getContentResolver(),document(parent),DocumentsContract.Document.MIME_TYPE_DIR,name);
        if(created==null)throw new Exception("El proveedor no permite crear la carpeta «"+name+"»");
        return created;
    }
    private Uri base()throws Exception{
        Uri root=tree();
        // Compatible con cualquier carpeta de proveedor SAF compatible, tanto raíz como subcarpeta.
        // No se necesita una cuenta de un proveedor específico.
        if(findChild(root,"actual")!=null||findChild(root,"sincronizacion")!=null)return root;
        return ensureDirectory(root,ROOT);
    }
    private Uri existingBase()throws Exception{
        Uri root=tree();if(findChild(root,"actual")!=null||findChild(root,"sincronizacion")!=null)return root;
        return findChild(root,ROOT);
    }
    private byte[] read(Uri uri)throws Exception{
        if(graphMode())return oneDrive.read(driveId(uri));
        try(InputStream in=activity.getContentResolver().openInputStream(uri);
            ByteArrayOutputStream out=new ByteArrayOutputStream()){
            if(in==null)throw new Exception("No se puede leer documento de nube");
            byte[] buf=new byte[8192];int n;while((n=in.read(buf))>=0){out.write(buf,0,n);if(out.size()>8*1024*1024)throw new Exception("El archivo de nube es demasiado grande");}
            return out.toByteArray();
        }
    }
    private void putFile(Uri folder,String name,byte[] bytes,boolean overwrite)throws Exception{
        if(graphMode()){oneDrive.write(driveId(folder),name,bytes,overwrite);return;}
        Uri file=findChild(folder,name);
        if(file!=null&&!overwrite)return;
        if(file==null)file=DocumentsContract.createDocument(activity.getContentResolver(),document(folder),
            name.endsWith(".json")?"application/json":"text/plain",name);
        if(file==null)throw new Exception("No se puede crear el archivo «"+name+"»");
        try(OutputStream out=activity.getContentResolver().openOutputStream(file,"wt")){
            if(out==null)throw new Exception("No se puede escribir en la carpeta seleccionada");
            out.write(bytes);out.flush();
        }
    }
    synchronized JSONObject backup()throws Exception{
        // Never publish the out-of-box library merely because the user selected a folder.
        if(folder.isEmpty()||!activated)return status();
        try {
            JSONObject library=new JSONObject(new String(Files.readAllBytes(new File(appDir,"colecciones.json").toPath()),StandardCharsets.UTF_8));
            if(!Backend.validLibrary(library))throw new Exception("No hay biblioteca local válida");
            Uri base=base(),devices=ensureDirectory(ensureDirectory(base,"sincronizacion"),"dispositivos");
            byte[] data=sync.snapshot(library).toString(2).getBytes(StandardCharsets.UTF_8);
            putFile(devices,sync.deviceId()+".json",data,true);
            Uri own=findChild(devices,sync.deviceId()+".json");
            if(own==null||!java.util.Arrays.equals(read(own),data))throw new Exception("El proveedor no confirmó la escritura de la sincronización");
            // Daily rescue point for mobile, separate from the Windows backup.
            Uri mobile=ensureDirectory(ensureDirectory(base,"sincronizacion"),"respaldo_movil");
            putFile(mobile,LocalDate.now()+"_"+sync.deviceId()+".json",data,false);
            lastBackup=OffsetDateTime.now().toString();lastError="";
            activity.getPreferences(Context.MODE_PRIVATE).edit().putString("last_backup",lastBackup).apply();
        }catch(Exception e){lastError=e.getMessage()==null?"Error al acceder a la carpeta compartida":e.getMessage();}
        return status();
    }
    private static final class Plan {
        JSONObject library,revisions,known;
        JSONArray conflicts=new JSONArray();
        int sources=0;
        boolean hasLegacy=false;
        String fingerprint="";
    }
    private static String hexHash(byte[] input)throws Exception{
        byte[] hash=MessageDigest.getInstance("SHA-256").digest(input);
        StringBuilder out=new StringBuilder();for(byte b:hash)out.append(String.format(java.util.Locale.ROOT,"%02x",b&255));
        return out.toString();
    }
    private Plan buildPlan(Backend backend)throws Exception{
        if(folder.isEmpty())throw new Exception("Selecciona antes una carpeta compartida");
        Plan p=new Plan();p.library=new JSONObject(backend.readOrDefault("colecciones.json","{}"));
        if(!Backend.validLibrary(p.library))throw new Exception("Biblioteca local inválida");
        p.revisions=sync.revisions();p.known=sync.known();
        // Incluir también el estado local en la huella: editar después de la vista previa
        // obliga a revisar otra vez antes de importar datos.
        StringBuilder digest=new StringBuilder("Cine365-sync-v2|local:"+hexHash(
            backend.readOrDefault("colecciones.json","").getBytes(StandardCharsets.UTF_8))+"|");
        Uri root=existingBase();
        if(root!=null){
            Uri actual=findChild(root,"actual");
            if(actual!=null){Uri file=findChild(actual,"colecciones.json");if(file!=null){
                byte[] bytes=read(file);digest.append("legacy:").append(hexHash(bytes)).append('|');
                JSONObject remote=new JSONObject(new String(bytes,StandardCharsets.UTF_8));
                if(!Backend.validLibrary(remote))throw new Exception("La copia Windows no es válida; no se han cambiado tus datos");
                JSONObject snapshot=new JSONObject().put("protocol",1).put("device_id","windows-legacy")
                    .put("library",remote).put("revisions",new JSONObject());
                SyncStore.Merge m=SyncStore.merge(p.library,p.revisions,p.known,snapshot);
                p.library=m.library;p.revisions=m.revisions;p.known=m.known;
                p.sources++;p.hasLegacy=true;
            }}
            Uri syncDir=findChild(root,"sincronizacion"),deviceDir=syncDir==null?null:findChild(syncDir,"dispositivos");
            if(deviceDir!=null){
                List<Uri> snaps=children(deviceDir);
                Collections.sort(snaps,(a,b)->a.toString().compareTo(b.toString()));
                for(Uri uri:snaps){
                    byte[] bytes=read(uri);digest.append("peer:").append(uri.toString()).append(':').append(hexHash(bytes)).append('|');
                    JSONObject snapshot=new JSONObject(new String(bytes,StandardCharsets.UTF_8));
                    if(sync.deviceId().equals(snapshot.optString("device_id")))continue;
                    SyncStore.Merge m=SyncStore.merge(p.library,p.revisions,p.known,snapshot);
                    p.library=m.library;p.revisions=m.revisions;p.known=m.known;p.sources++;
                    for(int i=0;i<m.conflicts.length();i++){
                        JSONObject item=m.conflicts.getJSONObject(i);
                        boolean exists=false;
                        for(int j=0;j<p.conflicts.length();j++)if(item.optString("key").equals(p.conflicts.getJSONObject(j).optString("key")))exists=true;
                        if(!exists)p.conflicts.put(item);
                    }
                }
            }
        }
        // La huella incluye los archivos remotos y la biblioteca local original.
        p.fingerprint=hexHash(digest.toString().getBytes(StandardCharsets.UTF_8));
        return p;
    }
    synchronized JSONObject preview(Backend backend)throws Exception{
        try {
            Plan plan=buildPlan(backend);
            previewToken=plan.fingerprint;
            JSONObject local=new JSONObject(backend.readOrDefault("colecciones.json","{}"));
            return status().put("token",previewToken).put("sources",plan.sources)
                .put("legacy",plan.hasLegacy).put("local_lists",local.getJSONArray("lists").length())
                .put("local_seen",SyncStore.watchedCount(local))
                .put("result_lists",plan.library.getJSONArray("lists").length())
                .put("result_seen",SyncStore.watchedCount(plan.library))
                .put("conflicts",plan.conflicts).put("requires_confirmation",true);
        }catch(Exception e){lastError=e.getMessage();throw e;}
    }
    synchronized JSONObject synchronize(Backend backend,String expected,JSONObject decisions)throws Exception{
        if(folder.isEmpty())throw new Exception("Configura primero una carpeta compartida");
        if(previewToken.isEmpty()||!previewToken.equals(expected))throw new Exception("Revisa primero los datos para sincronizar");
        try {
            Plan plan=buildPlan(backend);
            if(!plan.fingerprint.equals(expected))throw new Exception("La carpeta compartida ha cambiado desde la vista previa; vuelve a revisar");
            JSONObject resolved=new JSONObject();
            for(int i=0;i<plan.conflicts.length();i++){
                JSONObject conflict=plan.conflicts.getJSONObject(i);
                String key=conflict.optString("key");
                if(!decisions.has(key))throw new Exception("Quedan decisiones pendientes para «"+key+"». No se ha sobrescrito nada.");
                resolved.put(key,decisions.optBoolean(key));
            }
            if(resolved.length()>0){
                plan.library=SyncStore.applyStatuses(plan.library,resolved);
            }
            // Guardar la biblioteca con copia de recuperación antes de modificar
            // el historial causal o publicar cambios en otra nube.
            backend.storeStateImported(plan.library);
            if(resolved.length()>0){
                sync.resolve(plan.revisions,plan.known,resolved);
                plan.revisions=sync.revisions();plan.known=sync.known();
            }else sync.accept(plan.revisions,plan.known);
            activated=true;autoNeedsReview=false;
            activity.getPreferences(Context.MODE_PRIVATE).edit().putBoolean("sync_activated_v2",true).apply();
            JSONObject backup=backup();
            if(!backup.optString("last_error").isEmpty())throw new Exception(backup.optString("last_error"));
            lastSync=OffsetDateTime.now().toString();
            activity.getPreferences(Context.MODE_PRIVATE).edit().putString("last_sync",lastSync).apply();
            previewToken="";
            return status().put("changed",true).put("sources",plan.sources)
                .put("lists",plan.library.getJSONArray("lists").length())
                .put("seen",SyncStore.watchedCount(plan.library));
        }catch(Exception e){lastError=e.getMessage();throw e;}
    }
    /**
     * Automatic inbound merge is deliberately opt-in by first manual pairing.
     * New list additions and causal seen/pending changes are imported; opposite
     * concurrent decisions are left untouched until the user reviews them.
     */
    synchronized JSONObject autoSynchronize(Backend backend)throws Exception{
        if(folder.isEmpty()||!activated)return status().put("changed",false).put("skipped","not_paired");
        if(!previewToken.isEmpty())return status().put("changed",false).put("skipped","manual_review_open");
        try{
            Plan plan=buildPlan(backend);
            if(plan.conflicts.length()>0){
                autoNeedsReview=true;
                return status().put("changed",false).put("conflicts",plan.conflicts.length());
            }
            if(plan.sources==0)return status().put("changed",false).put("skipped","no_remote_sources");
            JSONObject before=new JSONObject(backend.readOrDefault("colecciones.json","{}"));
            if(!Backend.validLibrary(before))throw new Exception("Biblioteca local no válida");
            boolean changed=!before.toString().equals(plan.library.toString());
            boolean metadataChanged=!sync.revisions().toString().equals(plan.revisions.toString())
                ||!sync.known().toString().equals(plan.known.toString());
            if(!changed&&!metadataChanged){
                autoNeedsReview=false;
                return status().put("changed",false);
            }
            // Creates the same daily local recovery point as a manual import.
            if(changed)backend.storeStateImported(plan.library);
            sync.accept(plan.revisions,plan.known);
            autoNeedsReview=false;
            JSONObject result=backup();
            if(!result.optString("last_error").isEmpty())throw new Exception(result.optString("last_error"));
            lastSync=OffsetDateTime.now().toString();
            activity.getPreferences(Context.MODE_PRIVATE).edit().putString("last_sync",lastSync).apply();
            return status().put("changed",changed).put("metadata_updated",metadataChanged)
                .put("lists",plan.library.getJSONArray("lists").length())
                .put("seen",SyncStore.watchedCount(plan.library))
                .put("sources",plan.sources);
        }catch(Exception e){
            lastError=e.getMessage()==null?"No se pudo revisar la nube":e.getMessage();
            throw e;
        }
    }

}
