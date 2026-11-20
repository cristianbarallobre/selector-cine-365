package es.cine365.app;

import org.json.JSONArray;
import org.json.JSONObject;
import java.io.File;
import java.net.URL;
import java.net.URLEncoder;
import java.net.HttpURLConnection;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.text.Normalizer;
import java.time.Instant;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashSet;
import java.util.LinkedHashSet;
import java.util.Set;

final class Tmdb {
    private final File dir;
    private String key="";
    private JSONObject cache=new JSONObject(),overrides=new JSONObject();
    Tmdb(File dir){
        this.dir=dir;
        try{key=new JSONObject(readFile("plataformas_ajustes.json","{}" )).optString("api_key","");}catch(Exception ignored){}
        try{cache=new JSONObject(readFile("plataformas_cache.json","{}" )).optJSONObject("entries");if(cache==null)cache=new JSONObject();}catch(Exception ignored){}
        try{overrides=new JSONObject(readFile("identificaciones_tmdb.json","{}"));}catch(Exception ignored){}
    }
    private JSONObject identitiesCache(){
        try{return new JSONObject(readFile("identidades_cine.json","{}")).optJSONObject("entries");}
        catch(Exception ignored){return null;}
    }
    synchronized JSONObject identities()throws Exception{
        JSONObject entries=new JSONObject();
        JSONObject records=identitiesCache();
        if(records!=null)for(java.util.Iterator<String>it=records.keys();it.hasNext();){
            String k=it.next();int id=records.optJSONObject(k)==null?0:records.optJSONObject(k).optInt("tmdb_id");
            if(id>0)entries.put(k,id);
        }
        for(java.util.Iterator<String>it=cache.keys();it.hasNext();){
            String k=it.next();int id=cache.optJSONObject(k)==null?0:cache.optJSONObject(k).optInt("tmdb_id");
            if(id>0)entries.put(k,id);
        }
        for(java.util.Iterator<String>it=overrides.keys();it.hasNext();){
            String k=it.next();if(overrides.optInt(k)>0)entries.put(k,overrides.optInt(k));
        }
        return new JSONObject().put("configured",hasKey()).put("entries",entries);
    }
    synchronized JSONObject identify(String title,int year,boolean force)throws Exception{
        if(!hasKey())return new JSONObject().put("state","unconfigured");
        String k=cacheKey(title,year);
        if(!force){
            int id=identities().getJSONObject("entries").optInt(k);
            if(id>0)return new JSONObject().put("state","ok").put("tmdb_id",id);
            // También recordar resultados negativos para no repetir búsquedas al abrir.
            JSONObject records=identitiesCache(),prior=records==null?null:records.optJSONObject(k);
            if(prior!=null && !"ok".equals(prior.optString("state")))
                return new JSONObject(prior.toString());
        }
        JSONObject data=lookup(title,year,force);int id=data.optInt("tmdb_id");
        JSONObject records=identitiesCache();if(records==null)records=new JSONObject();
        records.put(k,new JSONObject().put("tmdb_id",id).put("state",id>0?"ok":data.optString("state","not-found"))
            .put("updated",java.time.OffsetDateTime.now().toString()));
        Backend.atomicWrite(new File(dir,"identidades_cine.json"),new JSONObject().put("entries",records).toString(2));
        return new JSONObject().put("state",id>0?"ok":data.optString("state","not-found"))
            .put("tmdb_id",id).put("note",data.optString("note",""));
    }
    synchronized JSONObject manualOverrides()throws Exception {
        return new JSONObject(overrides.toString());
    }
    // Importa solo IDs de TMDB. Las claves API personales nunca salen del móvil.
    synchronized void importManualOverrides(JSONObject changes)throws Exception {
        if(changes==null||changes.length()>10000)throw new Exception("Identificaciones TMDB inválidas");
        JSONObject next=new JSONObject(overrides.toString()),records=identitiesCache();
        if(records==null)records=new JSONObject();
        boolean modified=false;
        for(java.util.Iterator<String>it=changes.keys();it.hasNext();){
            String k=it.next();Object v=changes.get(k);
            if(k.length()<3||k.length()>275||!k.contains("|")||!(v instanceof Number))
                throw new Exception("Identificación TMDB inválida");
            int id=((Number)v).intValue();
            if(id<0||id>999999999)throw new Exception("ID de TMDB incorrecto");
            if(id==0){
                if(next.has(k)){next.remove(k);modified=true;}
                if(records.has(k)){records.remove(k);modified=true;}
                cache.remove(k);
            }else if(next.optInt(k)!=id){
                next.put(k,id);modified=true;
                records.put(k,new JSONObject().put("tmdb_id",id).put("state","ok")
                    .put("updated",java.time.OffsetDateTime.now().toString()));
                cache.remove(k);
            }
        }
        if(!modified)return;
        Backend.atomicWrite(new File(dir,"identificaciones_tmdb.json"),next.toString(2));
        Backend.atomicWrite(new File(dir,"identidades_cine.json"),new JSONObject().put("entries",records).toString(2));
        Backend.atomicWrite(new File(dir,"plataformas_cache.json"),new JSONObject().put("entries",cache).toString());
        overrides=next;
    }
    synchronized boolean hasKey(){return !key.isEmpty();}
    synchronized String getKey(){return key;}
    synchronized void configure(String k)throws Exception{
        if(k.length()<12||k.length()>512)throw new Exception("Clave no válida");
        request("/configuration","",k);
        Backend.atomicWrite(new File(dir,"plataformas_ajustes.json"),new JSONObject().put("api_key",k).toString(2));
        key=k;
    }
    private static String enc(String s){try{return URLEncoder.encode(s,"UTF-8");}catch(Exception e){return s;}}
    static String clean(String s){
        String normal=Normalizer.normalize(s.toLowerCase(java.util.Locale.ROOT),Normalizer.Form.NFD).replaceAll("\\p{M}","");
        return normal.replaceAll("[^a-z0-9]+"," ").trim().replaceAll(" +"," ");
    }
    private static Set<String> words(String x){return new HashSet<>(Arrays.asList(clean(x).split(" ")));}
    static double similarity(String a,String b){
        String aa=clean(a),bb=clean(b);
        if(aa.equals(bb) && !aa.isEmpty())return 1;
        if(aa.isEmpty()||bb.isEmpty())return 0;
        if(aa.length()>12 && bb.contains(aa) || bb.length()>12 && aa.contains(bb))return .83;
        Set<String> x=words(aa),y=words(bb);int count=0;
        for(String s:x)if(y.contains(s))count++;
        return (double)count/(x.size()+y.size()-count);
    }
    static ArrayList<String> variants(String title){
        LinkedHashSet<String> vals=new LinkedHashSet<>();
        vals.add(title);
        for(String divider:new String[]{". ",": "," — "," - "}){
            int pos=title.indexOf(divider);if(pos>=4)vals.add(title.substring(0,pos));
        }
        int idx=title.indexOf(" (");if(idx>3)vals.add(title.substring(0,idx));
        for(String x:new ArrayList<>(vals))vals.add(Normalizer.normalize(x,Normalizer.Form.NFD).replaceAll("\\p{M}",""));
        ArrayList<String> out=new ArrayList<>();for(String v:vals)if(v.trim().length()>=4 && out.size()<7)out.add(v.trim());
        return out;
    }
    private JSONObject request(String path,String query,String key)throws Exception {
        String url="https://api.themoviedb.org/3"+path;
        if(query!=null&&!query.isEmpty())url+="?"+query;
        if(!key.startsWith("eyJ")&&!key.contains("."))url+=(url.contains("?")?"&":"?")+"api_key="+enc(key);
        HttpURLConnection c=(HttpURLConnection)new URL(url).openConnection();
        c.setConnectTimeout(10000);c.setReadTimeout(15000);c.setRequestProperty("Accept","application/json");
        if(key.startsWith("eyJ")||key.contains("."))c.setRequestProperty("Authorization","Bearer "+key);
        try{
            int status=c.getResponseCode();
            if(status==401||status==403)throw new Exception("Clave de TMDB no válida");
            if(status==429)throw new Exception("TMDB ha limitado temporalmente las consultas");
            if(status!=200)throw new Exception("TMDB HTTP "+status);
            try(InputStream in=c.getInputStream()){
                java.io.ByteArrayOutputStream bytes=new java.io.ByteArrayOutputStream();
                byte[] buf=new byte[8192];int n;while((n=in.read(buf))!=-1){bytes.write(buf,0,n);if(bytes.size()>3_000_000)break;}
                return new JSONObject(bytes.toString("UTF-8"));
            }
        }finally{c.disconnect();}
    }
    private static int year(String s){try{return Integer.parseInt(s.substring(0,4));}catch(Exception e){return 0;}}
    private static String cacheKey(String title,int yr){return clean(title)+"|"+yr;}
    synchronized void match(String title,int year,int id)throws Exception{
        if(title.isEmpty()||id<0||id>999_999_999)throw new Exception("ID de película no válido");
        if(id>0){JSONObject movie=request("/movie/"+id,"language=es-ES",getKey());
            int y=year(movie.optString("release_date"));if(y>0&&year>0&&Math.abs(y-year)>2)throw new Exception("El año no coincide con la película");
        }
        String cacheKey=cacheKey(title,year);
        if(id>0)overrides.put(cacheKey,id);else overrides.remove(cacheKey);
        Backend.atomicWrite(new File(dir,"identificaciones_tmdb.json"),overrides.toString(2));
        cache.remove(cacheKey);Backend.atomicWrite(new File(dir,"plataformas_cache.json"),new JSONObject().put("entries",cache).toString());
    }
    private static JSONArray providers(JSONObject es,String name){
        JSONArray items=es.optJSONArray(name),result=new JSONArray();
        if(items==null)return result;
        Set<String> seen=new HashSet<>();
        for(int i=0;i<items.length();i++){
            String title=items.optJSONObject(i)==null?"":items.optJSONObject(i).optString("provider_name");
            if(!title.isEmpty()&&seen.add(title))result.put(title);
        }
        return result;
    }
    private static JSONObject empty(String state,String note)throws Exception{
        return new JSONObject().put("state",state).put("note",note).put("subscriptions",new JSONArray())
            .put("free",new JSONArray()).put("ads",new JSONArray()).put("rent",new JSONArray()).put("buy",new JSONArray());
    }
    // Datos guardados en disco: lectura local al abrir, sin conexión a TMDB.
    synchronized JSONObject cachedEntries()throws Exception{
        return new JSONObject().put("entries",new JSONObject(cache.toString()));
    }
    synchronized JSONObject lookup(String title,int releaseYear,boolean force)throws Exception{
        if(title==null||title.trim().isEmpty()||releaseYear<0||releaseYear>2200)return empty("error","Película no válida");
        String ck=cacheKey(title,releaseYear);
        JSONObject saved=cache.optJSONObject(ck);
        // Sin caducidad por tiempo: refrescar solo cuando se pide explícitamente,
        // se incorpora una película nueva o se corrige su identificación.
        if(!force&&saved!=null&&saved.optInt("metadata_version")>=3)return saved;
        String api=getKey();if(api.isEmpty())return empty("error","Configura una clave de TMDB");
        JSONObject result;
        try{
            int id=overrides.optInt(ck),fromYear=releaseYear;
            if(id==0){
                double best=0;
                outer:for(String query:variants(title)){
                    JSONObject found=request("/search/movie","language=es-ES&include_adult=false&query="+enc(query),api);
                    JSONArray movies=found.optJSONArray("results");if(movies==null)continue;
                    for(int i=0;i<movies.length();i++){
                        JSONObject f=movies.optJSONObject(i);if(f==null)continue;
                        int y=year(f.optString("release_date")),diff=Math.abs(y-fromYear);
                        if(y>0&&fromYear>0&&diff>2)continue;
                        double score=Math.max(similarity(query,f.optString("title")),similarity(query,f.optString("original_title")));
                        if(y>0&&fromYear>0){score+=diff==0?.1:diff==1?.03:-.12;}
                        if(score>best && score>=.8){best=score;id=f.optInt("id");}
                    }
                    if(id>0)break outer;
                }
            }
            if(id==0){result=empty("not-found","");result.put("date",Instant.now().toString()).put("metadata_version",3);}
            else{
                JSONObject info=request("/movie/"+id,"language=es-ES&append_to_response=credits",api);
                JSONObject prov=request("/movie/"+id+"/watch/providers","",api);
                JSONObject es=prov.optJSONObject("results");es=es==null?null:es.optJSONObject("ES");
                result=empty(es==null?"no-data":"ok","");
                result.put("tmdb_id",id).put("title",info.optString("title"))
                    .put("vote_average",info.optDouble("vote_average")).put("vote_count",info.optInt("vote_count"))
                    .put("poster_path",info.optString("poster_path")).put("runtime",info.optInt("runtime"));
                JSONArray g=new JSONArray(),c=new JSONArray();
                JSONArray genres=info.optJSONArray("genres");
                if(genres!=null)for(int j=0;j<genres.length();j++)g.put(genres.optJSONObject(j).optString("name"));
                JSONArray nations=info.optJSONArray("production_countries");
                if(nations!=null)for(int j=0;j<nations.length();j++)c.put(nations.optJSONObject(j).optString("name"));
                result.put("genres",g).put("countries",c);
                StringBuilder director=new StringBuilder();
                JSONObject credits=info.optJSONObject("credits");JSONArray crew=credits==null?null:credits.optJSONArray("crew");
                if(crew!=null)for(int j=0;j<crew.length();j++){
                    JSONObject person=crew.optJSONObject(j);
                    if(person!=null&&"Director".equals(person.optString("job"))){if(director.length()>0)director.append(" / ");director.append(person.optString("name"));}
                }
                result.put("director",director.toString());
                result.put("film_url",es==null?"":es.optString("link"));
                if(es!=null){result.put("subscriptions",providers(es,"flatrate")).put("free",providers(es,"free"))
                   .put("ads",providers(es,"ads")).put("rent",providers(es,"rent")).put("buy",providers(es,"buy"));}
                result.put("date",Instant.now().toString()).put("metadata_version",3);
            }
            cache.put(ck,result);Backend.atomicWrite(new File(dir,"plataformas_cache.json"),new JSONObject().put("entries",cache).toString());
            return result;
        }catch(Exception e){
            // Si no hay conexión no insistimos en cada inicio ni tras marcar una vista.
            JSONObject failure=empty("error",e.getMessage()==null?"No se puede consultar TMDB":e.getMessage());
            failure.put("date",Instant.now().toString()).put("metadata_version",3);
            cache.put(ck,failure);
            Backend.atomicWrite(new File(dir,"plataformas_cache.json"),new JSONObject().put("entries",cache).toString());
            return failure;
        }
    }

    // Aceptamos únicamente rutas de imágenes de TMDB, sin dominios ni rutas arbitrarias.
    byte[] poster(String path)throws Exception{
        if(path==null||!path.matches("/[A-Za-z0-9_-]{1,160}\\.(?:jpg|jpeg|png|webp)"))
            throw new Exception("Ruta de cartel inválida");
        File folder=new File(dir,"carteles_tmdb");
        if(!folder.isDirectory()&&!folder.mkdirs())throw new Exception("No se pudo crear caché de carteles");
        File local=new File(folder,path.substring(1));
        if(local.isFile()&&local.length()>128&&local.length()<=2_500_000)return Files.readAllBytes(local.toPath());
        HttpURLConnection conn=(HttpURLConnection)new URL("https://image.tmdb.org/t/p/w342"+path).openConnection();
        conn.setConnectTimeout(8000);conn.setReadTimeout(12000);conn.setRequestProperty("Accept","image/*");
        byte[] data;
        try{
            if(conn.getResponseCode()!=200||conn.getContentType()==null||!conn.getContentType().startsWith("image/"))
                throw new Exception("Cartel no disponible");
            try(InputStream in=conn.getInputStream();java.io.ByteArrayOutputStream out=new java.io.ByteArrayOutputStream()){
                byte[] buffer=new byte[8192];int n;
                while((n=in.read(buffer))!=-1){
                    if(out.size()+n>2_500_000)throw new Exception("Cartel demasiado grande");
                    out.write(buffer,0,n);
                }
                data=out.toByteArray();
            }
        }finally{conn.disconnect();}
        if(data.length<=128)throw new Exception("Cartel vacío");
        java.nio.file.Path temp=Files.createTempFile(folder.toPath(),"tmdb-", ".tmp");
        try{
            Files.write(temp,data);
            try{Files.move(temp,local.toPath(),java.nio.file.StandardCopyOption.ATOMIC_MOVE,java.nio.file.StandardCopyOption.REPLACE_EXISTING);}
            catch(java.nio.file.AtomicMoveNotSupportedException e){Files.move(temp,local.toPath(),java.nio.file.StandardCopyOption.REPLACE_EXISTING);}
        }finally{Files.deleteIfExists(temp);}
        return data;
    }

    private String readFile(String name,String defaultValue){
        try{return new String(Files.readAllBytes(new File(dir,name).toPath()),StandardCharsets.UTF_8);}catch(Exception e){return defaultValue;}
    }
}
