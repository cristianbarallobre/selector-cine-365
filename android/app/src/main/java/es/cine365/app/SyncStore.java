package es.cine365.app;

import org.json.JSONArray;
import org.json.JSONObject;
import java.io.File;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.text.Normalizer;
import java.util.HashSet;
import java.util.Iterator;
import java.util.Locale;
import java.util.Set;
import java.util.UUID;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Selector Cine 365 v9: causal per-film revisions. Independent of platform and clock skew.
 * Files from older installations are read as legacy snapshots, never edited in place.
 */
final class SyncStore {
    static final int PROTOCOL=2;
    private static final Pattern CHECK=Pattern.compile("^\\s*[☐□☑☒✓✔✅]\\s*");
    private static final Pattern NUMBER=Pattern.compile("^(?:\\d{1,4}[.)]|[-•])\\s*");
    private static final Pattern YEAR=Pattern.compile("\\s+\\((\\d{4})\\)$");
    private final File stateFile;
    private JSONObject state;

    SyncStore(File localDir) {
        stateFile=new File(localDir,"sync_dispositivo_v2.json");
        try {state=new JSONObject(new String(Files.readAllBytes(stateFile.toPath()),StandardCharsets.UTF_8));}
        catch(Exception ignored) {state=new JSONObject();}
        try {
            state.put("protocol",PROTOCOL);
            if(state.optString("device_id").isEmpty())state.put("device_id","android-"+UUID.randomUUID());
            if(state.optJSONObject("revisions")==null)state.put("revisions",new JSONObject());
            if(state.optJSONObject("known") ==null)state.put("known",new JSONObject());
            if(state.optJSONArray("history")==null)state.put("history",new JSONArray());
            if(state.optInt("seq")<0)state.put("seq",0);
            save();
        }catch(Exception e){throw new IllegalStateException(e);}
    }
    synchronized String deviceId(){return state.optString("device_id");}
    synchronized JSONObject revisions()throws Exception {return new JSONObject(state.getJSONObject("revisions").toString());}
    synchronized JSONObject known()throws Exception {return new JSONObject(state.getJSONObject("known").toString());}
    private synchronized void appendHistory(String key,JSONObject rev,String origin)throws Exception{
        JSONArray history=state.optJSONArray("history");
        if(history==null)history=new JSONArray();
        history.put(new JSONObject()
            .put("at",System.currentTimeMillis())
            .put("changed_at",rev.optLong("changed_at"))
            .put("key",key)
            .put("checked",rev.optBoolean("checked"))
            .put("device",rev.optString("device"))
            .put("origin",origin));
        if(history.length()>80){
            JSONArray recent=new JSONArray();
            for(int i=history.length()-80;i<history.length();i++)recent.put(history.get(i));
            history=recent;
        }
        state.put("history",history);
    }
    synchronized JSONObject history()throws Exception{
        JSONArray saved=state.optJSONArray("history"),events=new JSONArray();
        if(saved!=null)for(int i=saved.length()-1;i>=0;i--)events.put(saved.getJSONObject(i));
        return new JSONObject().put("device_id",deviceId()).put("events",events);
    }
    synchronized void accept(JSONObject revisions, JSONObject known)throws Exception {
        JSONObject previous=state.getJSONObject("revisions");
        for(Iterator<String> keys=revisions.keys();keys.hasNext();){
            String key=keys.next();JSONObject next=revisions.optJSONObject(key),old=previous.optJSONObject(key);
            if(next!=null&&(old==null||old.optBoolean("checked")!=next.optBoolean("checked")
                ||!old.optString("device").equals(next.optString("device"))
                ||old.optLong("seq")!=next.optLong("seq")))appendHistory(key,next,"recibido");
        }
        state.put("revisions",new JSONObject(revisions.toString()));
        state.put("known",new JSONObject(known.toString()));
        save();
    }
    private synchronized void save()throws Exception {Backend.atomicWrite(stateFile,state.toString(2));}

    static final class Film {
        final String key;final boolean seen;
        Film(String k,boolean s){key=k;seen=s;}
    }
    static Film film(String line) {
        String original=line.trim();
        if(original.isEmpty()||original.startsWith("#")||
            original.toUpperCase(Locale.ROOT).startsWith("ATLAS DEL CINE")||
            original.toUpperCase(Locale.ROOT).startsWith("SELECTOR CINE 365")||
            original.matches("[-_= ]+"))return null;
        boolean watched=original.matches("^[☑☒✓✔✅].*");
        String text=NUMBER.matcher(CHECK.matcher(original).replaceFirst("")).replaceFirst("");
        Matcher m=YEAR.matcher(text);String year="";
        if(m.find()){year=m.group(1);text=text.substring(0,m.start());}
        String normalized=Normalizer.normalize(text.toLowerCase(Locale.ROOT),Normalizer.Form.NFD).replaceAll("\\p{M}","")
            .replace("æ","ae").replace("œ","oe")
            .replaceAll("[^\\p{L}\\p{N}]+"," ").trim().replaceAll(" +"," ");
        if(normalized.isEmpty())return null;
        return new Film(normalized+"|"+year,watched);
    }
    static String checkedLine(String line,boolean checked){
        if(film(line)==null)return line;
        return (checked?"☑ ":"☐ ")+CHECK.matcher(line.trim()).replaceFirst("");
    }
    static JSONObject statuses(JSONObject library)throws Exception {
        JSONObject result=new JSONObject();JSONArray lists=library.getJSONArray("lists");
        for(int i=0;i<lists.length();i++)for(String line:lists.getJSONObject(i).getString("text").split("\\n",-1)) {
            Film film=film(line);if(film!=null)result.put(film.key,result.optBoolean(film.key)||film.seen);
        }
        return result;
    }
    static int watchedCount(JSONObject library)throws Exception {
        JSONObject marks=statuses(library);int n=0;
        for(Iterator<String> keys=marks.keys();keys.hasNext();)if(marks.optBoolean(keys.next()))n++;
        return n;
    }
    static JSONObject joinClocks(JSONObject a, JSONObject b)throws Exception {
        JSONObject out=new JSONObject(a==null?"{}":a.toString());
        if(b!=null)for(Iterator<String>keys=b.keys();keys.hasNext();){String k=keys.next();out.put(k,Math.max(out.optLong(k),b.optLong(k)));}
        return out;
    }
    // 1: a causally newer, -1: b newer, 0: same, 2: concurrent.
    static int compareClocks(JSONObject a,JSONObject b){
        if(a==null||b==null)return 2;
        Set<String>keys=new HashSet<>();for(Iterator<String>it=a.keys();it.hasNext();)keys.add(it.next());
        for(Iterator<String>it=b.keys();it.hasNext();)keys.add(it.next());
        boolean newer=false,older=false;
        for(String k:keys){if(a.optLong(k)>b.optLong(k))newer=true;else if(a.optLong(k)<b.optLong(k))older=true;}
        return newer&&older?2:newer?1:older?-1:0;
    }
    private synchronized JSONObject newRevision(boolean watched)throws Exception {
        long seq=state.optLong("seq")+1;state.put("seq",seq);
        JSONObject vector=known();vector.put(deviceId(),seq);
        state.put("known",vector);
        return new JSONObject().put("checked",watched).put("device",deviceId())
            .put("seq",seq).put("clock",vector).put("changed_at",System.currentTimeMillis());
    }
    synchronized void observe(JSONObject before,JSONObject after)throws Exception {
        JSONObject old=before==null?new JSONObject():statuses(before),now=statuses(after),revisions=state.getJSONObject("revisions");
        boolean changed=false;
        for(Iterator<String>it=now.keys();it.hasNext();){String key=it.next(),ignored="";boolean next=now.optBoolean(key);
            if((old.has(key)&&old.optBoolean(key)!=next)||(!old.has(key)&&next&&before!=null)) {
                JSONObject revision=newRevision(next);
                revisions.put(key,revision);
                appendHistory(key,revision,"local");
                changed=true;
            }
        }
        if(changed)save();
    }
    synchronized JSONObject snapshot(JSONObject library)throws Exception {
        return new JSONObject().put("protocol",PROTOCOL).put("device_id",deviceId())
            .put("updated_at",System.currentTimeMillis()).put("library",library)
            .put("known",known()).put("revisions",revisions());
    }

    static final class Merge {
        final JSONObject library,revisions,known;
        final JSONArray conflicts;
        final boolean changed;
        Merge(JSONObject l,JSONObject r,JSONObject k,JSONArray issues,boolean c){library=l;revisions=r;known=k;conflicts=issues;changed=c;}
    }
    static Merge merge(JSONObject local,JSONObject revisions,JSONObject known,JSONObject peer)throws Exception {
        return merge(local,revisions,known,peer,false);
    }
    // Once paired, only explicit per-film revisions may change the viewed flag.
    // Raw old snapshot flags are NOT changes, and may not resurrect an old ✓.
    static Merge merge(JSONObject local,JSONObject revisions,JSONObject known,JSONObject peer,boolean strict)throws Exception {
        if(peer.optInt("protocol")!=1&&peer.optInt("protocol")!=PROTOCOL)
            throw new Exception("Formato de sincronización no reconocido");
        JSONObject remote=peer.getJSONObject("library");
        if(!Backend.validLibrary(remote))throw new Exception("Biblioteca compartida dañada; no se han cambiado tus datos");
        JSONObject result=new JSONObject(local.toString()), remoteMarks=statuses(remote);
        // Borrados que sobreviven a copias desactualizadas y a reinicios.
        JSONObject deleted=result.optJSONObject("deleted_lists");
        deleted=deleted==null?new JSONObject():new JSONObject(deleted.toString());
        JSONObject remoteDeleted=remote.optJSONObject("deleted_lists");
        if(remoteDeleted!=null)for(Iterator<String>it=remoteDeleted.keys();it.hasNext();){
            String id=it.next();if(remoteDeleted.optBoolean(id))deleted.put(id,true);
        }
        result.put("deleted_lists",deleted);
        // Se comparten solo IDs de película. Nunca claves API de TMDB.
        JSONObject tmdb=result.optJSONObject("tmdb_overrides");
        tmdb=tmdb==null?new JSONObject():new JSONObject(tmdb.toString());
        JSONObject remoteTMDB=remote.optJSONObject("tmdb_overrides");
        if(remoteTMDB!=null)for(Iterator<String>it=remoteTMDB.keys();it.hasNext();){
            String key=it.next(),localID=tmdb.optString(key,"");
            int id=remoteTMDB.optInt(key);
            if(localID.isEmpty()&&id>0)tmdb.put(key,id);
        }
        result.put("tmdb_overrides",tmdb);
        JSONArray left=result.getJSONArray("lists"),right=remote.getJSONArray("lists");
        for(int i=0;i<right.length();i++){
            JSONObject other=right.getJSONObject(i),own=null;
            if(deleted.optBoolean(other.optString("id")))continue;
            for(int j=0;j<left.length();j++)if(left.getJSONObject(j).getString("id").equals(other.getString("id"))){own=left.getJSONObject(j);break;}
            if(own==null){if(left.length()<40)left.put(new JSONObject(other.toString()));continue;}
            if(own.optString("name").startsWith("Selector Cine 365 ·")&&!other.optString("name").startsWith("Selector Cine 365 ·"))
                own.put("name",other.getString("name"));
            Set<String>keys=new HashSet<>();String before=own.getString("text");
            for(String line:before.split("\\n",-1)){Film f=film(line);if(f!=null)keys.add(f.key);}
            StringBuilder combined=new StringBuilder(before);
            for(String line:other.getString("text").split("\\n",-1)){
                Film f=film(line);if(f==null||keys.contains(f.key))continue;
                if(combined.length()>0&&combined.charAt(combined.length()-1)!='\n')combined.append('\n');
                combined.append(line);keys.add(f.key);
            }
            own.put("text",combined.toString());
        }
        // Una tombstone local o remota se impone incluso si el otro dispositivo
        // mantiene una copia anterior con el mismo ID de lista.
        if(deleted.length()>0){
            JSONArray keep=new JSONArray();
            for(int i=0;i<left.length();i++){
                JSONObject list=left.getJSONObject(i);
                if(!deleted.optBoolean(list.optString("id")))keep.put(list);
            }
            if(keep.length()==0)throw new Exception("La fusión borraría todas las listas. Operación cancelada.");
            left=keep;result.put("lists",left);
            boolean active=false;
            for(int i=0;i<left.length();i++)if(left.getJSONObject(i).optString("id").equals(result.optString("active_id")))active=true;
            if(!active)result.put("active_id",left.getJSONObject(0).getString("id"));
            JSONArray chosen=result.optJSONArray("comparison_ids"),cleaned=new JSONArray();
            if(chosen!=null)for(int i=0;i<chosen.length();i++){
                String id=chosen.optString(i);
                for(int j=0;j<left.length();j++)if(id.equals(left.getJSONObject(j).optString("id"))){cleaned.put(id);break;}
            }
            if(cleaned.length()>=2)result.put("comparison_ids",cleaned);
            else result.remove("comparison_ids");
        }
        if(result.optJSONArray("comparison_ids")==null&&remote.optJSONArray("comparison_ids")!=null)
            { JSONArray choices=remote.getJSONArray("comparison_ids"), valid=new JSONArray();
              for(int i=0;i<choices.length();i++){
                  String id=choices.optString(i);
                  for(int j=0;j<left.length();j++)if(id.equals(left.getJSONObject(j).optString("id"))){valid.put(id);break;}
              }
              if(valid.length()>=2)result.put("comparison_ids",valid);
            }
        JSONObject combinedKnown=joinClocks(known,peer.optJSONObject("known"));
        JSONObject joined=new JSONObject(revisions.toString()),incoming=peer.optJSONObject("revisions");
        JSONArray conflicts=new JSONArray();
        if(incoming!=null)for(Iterator<String>it=incoming.keys();it.hasNext();){
            String key=it.next();JSONObject theirs=incoming.optJSONObject(key),ours=joined.optJSONObject(key);
            if(theirs==null)continue;
            if(ours==null){joined.put(key,theirs);continue;}
            if(ours.optString("device").equals(theirs.optString("device"))){
                if(theirs.optLong("seq")>ours.optLong("seq"))joined.put(key,theirs);
                continue;
            }
            int direction=compareClocks(ours.optJSONObject("clock"),theirs.optJSONObject("clock"));
            if(direction==-1){joined.put(key,theirs);continue;}
            if(direction==1)continue;
            if(ours.optBoolean("checked")==theirs.optBoolean("checked")){
                JSONObject merged=new JSONObject(ours.toString());
                merged.put("clock",joinClocks(ours.optJSONObject("clock"),theirs.optJSONObject("clock")));
                joined.put(key,merged);continue;
            }
            conflicts.put(new JSONObject().put("key",key).put("local",ours.optBoolean("checked"))
                .put("remote",theirs.optBoolean("checked"))
                .put("local_device",ours.optString("device")).put("remote_device",theirs.optString("device")));
        }
        JSONObject allSeen=strict?statuses(local):statuses(result);
        if(!strict)for(Iterator<String>it=remoteMarks.keys();it.hasNext();){String key=it.next();
            if(remoteMarks.optBoolean(key))allSeen.put(key,true);
        }
        for(int i=0;i<left.length();i++){
            JSONObject list=left.getJSONObject(i);String[]lines=list.getString("text").split("\\n",-1);
            for(int j=0;j<lines.length;j++){
                Film film=film(lines[j]);if(film==null)continue;
                JSONObject rev=joined.optJSONObject(film.key);
                boolean marked=rev!=null?rev.optBoolean("checked"):allSeen.optBoolean(film.key);
                lines[j]=checkedLine(lines[j],marked);
            }
            list.put("text",String.join("\n",lines));
        }
        if(!Backend.validLibrary(result))throw new Exception("La fusión no produjo una biblioteca válida");
        return new Merge(result,joined,combinedKnown,conflicts,!result.toString().equals(local.toString()));
    }
    synchronized void resolve(JSONObject revisions,JSONObject known,JSONObject selections)throws Exception {
        JSONObject mergedKnown=joinClocks(known,this.known());
        state.put("known",mergedKnown);
        for(Iterator<String>it=selections.keys();it.hasNext();){
            String key=it.next();
            revisions.put(key,newRevision(selections.optBoolean(key)));
        }
        accept(revisions,joinClocks(mergedKnown,this.known()));
    }
    static JSONObject applyStatuses(JSONObject library,JSONObject changes)throws Exception {
        JSONObject result=new JSONObject(library.toString());
        JSONArray lists=result.getJSONArray("lists");
        for(int i=0;i<lists.length();i++){
            JSONObject list=lists.getJSONObject(i);String[]lines=list.getString("text").split("\\n",-1);
            for(int j=0;j<lines.length;j++){
                Film f=film(lines[j]);if(f!=null&&changes.has(f.key))lines[j]=checkedLine(lines[j],changes.optBoolean(f.key));
            }
            list.put("text",String.join("\n",lines));
        }
        return result;
    }
}
