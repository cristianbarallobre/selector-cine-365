package es.cine365.app;

import org.json.JSONArray;
import org.json.JSONObject;
import org.junit.Test;
import static org.junit.Assert.*;

/** Regression: a later Android unwatch must defeat an older Windows watch. */
public final class SyncStoreTest {
    private static JSONObject library(boolean seen)throws Exception{
        return new JSONObject().put("version",1).put("active_id","atlas")
            .put("lists",new JSONArray().put(new JSONObject().put("id","atlas")
                .put("name","Selector Cine 365 · 365")
                .put("text",(seen?"☑ ":"☐ ")+"001. Casablanca (1942)")));
    }
    private static JSONObject revision(boolean checked,String device,int seq,JSONObject clock)throws Exception{
        return new JSONObject().put("checked",checked).put("device",device)
            .put("seq",seq).put("clock",clock).put("changed_at",1000L+seq);
    }
    private static JSONObject peer(JSONObject lib,String device,JSONObject revs,JSONObject known)throws Exception{
        return new JSONObject().put("protocol",2).put("device_id",device)
            .put("library",lib).put("revisions",revs).put("known",known);
    }
    @Test public void laterAndroidUnwatchBeatsWindowsSeen()throws Exception{
        JSONObject windowsSeen=revision(true,"windows",1,new JSONObject().put("windows",1));
        JSONObject localRevs=new JSONObject().put("casablanca|1942",windowsSeen);
        JSONObject known=new JSONObject().put("windows",1);
        JSONObject newer=revision(false,"android",1,new JSONObject().put("windows",1).put("android",1));
        JSONObject android=peer(library(false),"android",
            new JSONObject().put("casablanca|1942",newer),
            new JSONObject().put("windows",1).put("android",1));
        SyncStore.Merge merged=SyncStore.merge(library(true),localRevs,known,android,true);
        assertEquals(0,merged.conflicts.length());
        assertFalse(SyncStore.statuses(merged.library).getBoolean("casablanca|1942"));
        assertFalse(merged.revisions.getJSONObject("casablanca|1942").getBoolean("checked"));
    }
    @Test public void pairedLegacySeenWithoutRevisionMustNotResurrect()throws Exception{
        JSONObject noLongerSeen=revision(false,"android",2,new JSONObject().put("windows",1).put("android",2));
        JSONObject localRevs=new JSONObject().put("casablanca|1942",noLongerSeen);
        SyncStore.Merge merged=SyncStore.merge(library(false),localRevs,
            new JSONObject().put("windows",1).put("android",2),
            peer(library(true),"old-beta",new JSONObject(),new JSONObject()),true);
        assertEquals(0,merged.conflicts.length());
        assertFalse(SyncStore.statuses(merged.library).getBoolean("casablanca|1942"));
    }
    @Test public void concurrentOppositeDecisionsRequireManualReview()throws Exception{
        JSONObject local=revision(true,"windows",1,new JSONObject().put("windows",1));
        JSONObject remote=revision(false,"android",1,new JSONObject().put("android",1));
        SyncStore.Merge merged=SyncStore.merge(library(true),
            new JSONObject().put("casablanca|1942",local),new JSONObject().put("windows",1),
            peer(library(false),"android",new JSONObject().put("casablanca|1942",remote),
                new JSONObject().put("android",1)),true);
        assertEquals(1,merged.conflicts.length());
    }
}
