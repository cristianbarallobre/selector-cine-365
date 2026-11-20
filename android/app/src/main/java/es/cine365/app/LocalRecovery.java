package es.cine365.app;

import org.json.JSONArray;
import org.json.JSONObject;
import java.io.File;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;

/** Local-only recovery. Never accepts arbitrary paths or cloud files for restoration. */
final class LocalRecovery {
    private final Backend backend;
    private final File root;
    LocalRecovery(Backend b){backend=b;root=b.dir;}

    private File sourceFile(String id)throws Exception{
        String[] parts=id.split(":",2);
        if(parts.length!=2||!parts[1].matches("[A-Za-z0-9._-]{5,95}\\.json"))
            throw new Exception("Identificador de copia no válido");
        String folder;
        if("rec".equals(parts[0]))folder="sync_recuperacion";
        else if("day".equals(parts[0]))folder="copias_automaticas";
        else throw new Exception("Esta carpeta no se puede restaurar");
        return new File(new File(root,folder),parts[1]);
    }
    JSONObject list()throws Exception{
        JSONArray output=new JSONArray();
        List<JSONObject> copies=new ArrayList<>();
        for(String type:new String[]{"rec","day"}){
            File dir=new File(root,"rec".equals(type)?"sync_recuperacion":"copias_automaticas");
            File[] files=dir.listFiles();
            if(files==null)continue;
            for(File f:files){
                if(!f.isFile()||f.length()>3_000_000||!f.getName().matches("[A-Za-z0-9._-]{5,95}\\.json"))continue;
                copies.add(new JSONObject().put("id",type+":"+f.getName()).put("name",f.getName())
                    .put("when_ms",f.lastModified()).put("bytes",f.length()));
            }
        }
        copies.sort(Comparator.comparingLong((JSONObject j)->j.optLong("when_ms")).reversed());
        for(int i=0;i<Math.min(25,copies.size());i++)output.put(copies.get(i));
        return new JSONObject().put("copies",output);
    }
    JSONObject restore(String id)throws Exception{
        File file=sourceFile(id);
        if(!file.isFile()||file.length()>3_000_000||file.length()<2)
            throw new Exception("No se encontró una copia local válida");
        JSONObject state=new JSONObject(new String(Files.readAllBytes(file.toPath()),StandardCharsets.UTF_8));
        if(!Backend.validLibrary(state))throw new Exception("La copia no contiene una biblioteca válida");
        JSONObject before=new JSONObject(backend.readOrDefault("colecciones.json","{}"));
        if(!Backend.validLibrary(before))throw new Exception("La biblioteca actual no es válida; no se ha restaurado");
        File recoveryDir=new File(root,"sync_recuperacion");
        if(!recoveryDir.isDirectory()&&!recoveryDir.mkdirs())throw new Exception("No se pudo crear la carpeta de recuperación");
        File rescue=new File(recoveryDir,System.currentTimeMillis()+"_"+System.nanoTime()+"_antes_restaurar.json");
        Backend.atomicWrite(rescue,before.toString(2));
        backend.storeStateImported(state);
        backend.sync.observe(before,state);
        JSONObject result=new JSONObject().put("restored",true)
            .put("lists",state.getJSONArray("lists").length())
            .put("seen",SyncStore.watchedCount(state));
        if(backend.cloud.status().optBoolean("activated")){
            JSONObject published=backend.cloud.backup();
            if(!published.optString("last_error").isEmpty())result.put("warning",published.optString("last_error"));
        }
        return result;
    }
}
