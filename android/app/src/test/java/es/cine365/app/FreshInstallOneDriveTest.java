package es.cine365.app;

import org.json.JSONArray;
import org.json.JSONObject;
import org.junit.Test;
import static org.junit.Assert.*;

/**
 * An Android reinstall starts with the bundled local starter list and a NEW
 * device identity. OneDrive may already contain a Windows snapshot and an
 * older Android snapshot. The first merge must preserve the newest known
 * Windows watched flags and import extra lists without replacing cloud data.
 */
public final class FreshInstallOneDriveTest {
    private static JSONObject list(String id,String name,String rows)throws Exception{
        return new JSONObject().put("id",id).put("name",name).put("text",rows);
    }
    private static JSONObject library(JSONObject... lists)throws Exception{
        JSONArray all=new JSONArray();
        for(JSONObject item:lists)all.put(item);
        return new JSONObject().put("version",1).put("active_id",lists[0].getString("id")).put("lists",all);
    }
    private static JSONObject rev(boolean checked,String device,int seq,JSONObject clock)throws Exception{
        return new JSONObject().put("checked",checked).put("device",device)
            .put("seq",seq).put("clock",clock).put("changed_at",1000L+seq);
    }
    private static JSONObject peer(String id,JSONObject lib,JSONObject revisions,JSONObject known)throws Exception{
        return new JSONObject().put("protocol",2).put("device_id",id).put("library",lib)
            .put("revisions",revisions).put("known",known);
    }
    @Test public void cleanAndroidImportsExistingWindowsAndIgnoresObsoletePhoneFlags()throws Exception{
        JSONObject fresh=library(list("atlas","Selector Cine 365 · 365",
            "☐ 001. Casablanca (1942)\n☐ 002. Metropolis (1927)"));
        JSONObject windowsClock=new JSONObject().put("android-old",1).put("desktop-windows",3);
        JSONObject windowsRev=rev(true,"desktop-windows",3,windowsClock);
        JSONObject windowsLibrary=library(
            list("atlas","Mi lista de Windows","☑ 001. Casablanca (1942)\n☐ 002. Metropolis (1927)"),
            list("lista_extra","Lista adicional de Windows","☑ The Matrix (1999)"));
        JSONObject fromWindows=peer("desktop-windows",windowsLibrary,
            new JSONObject().put("casablanca|1942",windowsRev),windowsClock);
        SyncStore.Merge imported=SyncStore.merge(fresh,new JSONObject(),new JSONObject(),fromWindows,false);
        assertEquals(0,imported.conflicts.length());
        assertEquals(2,imported.library.getJSONArray("lists").length());
        assertEquals("Mi lista de Windows",imported.library.getJSONArray("lists").getJSONObject(0).getString("name"));
        assertTrue(SyncStore.statuses(imported.library).getBoolean("casablanca|1942"));
        assertTrue(SyncStore.statuses(imported.library).getBoolean("the matrix|1999"));

        JSONObject oldPhoneLibrary=library(list("atlas","Selector Cine 365 · 365",
            "☐ 001. Casablanca (1942)\n☐ 002. Metropolis (1927)"));
        JSONObject oldRevision=rev(false,"android-old",1,new JSONObject().put("android-old",1));
        JSONObject oldPhone=peer("android-old",oldPhoneLibrary,
            new JSONObject().put("casablanca|1942",oldRevision),
            new JSONObject().put("android-old",1));
        SyncStore.Merge merged=SyncStore.merge(imported.library,imported.revisions,imported.known,oldPhone,false);
        assertEquals("The older phone should be causally superseded, not create a conflict",0,merged.conflicts.length());
        assertTrue("Old Android must not unwatch newer Windows revision",
            SyncStore.statuses(merged.library).getBoolean("casablanca|1942"));
        assertEquals(2,merged.library.getJSONArray("lists").length());
    }
    @Test public void freshAndroidDoesNotSilentlyResolveConcurrentOppositeFlags()throws Exception{
        JSONObject local=library(list("atlas","Móvil nuevo","☑ Casablanca (1942)"));
        JSONObject localRev=rev(true,"new-android",1,new JSONObject().put("new-android",1));
        JSONObject windows=library(list("atlas","Windows","☐ Casablanca (1942)"));
        JSONObject remoteRev=rev(false,"desktop-windows",1,new JSONObject().put("desktop-windows",1));
        SyncStore.Merge result=SyncStore.merge(local,
            new JSONObject().put("casablanca|1942",localRev),new JSONObject().put("new-android",1),
            peer("desktop-windows",windows,new JSONObject().put("casablanca|1942",remoteRev),
                new JSONObject().put("desktop-windows",1)),true);
        assertEquals("Opposite independent edits require explicit user decision",1,result.conflicts.length());
    }
    @Test public void deletedWindowsListDoesNotReturnFromOldAndroidCopy()throws Exception{
        JSONObject fresh=library(list("atlas","Principal","☐ Casablanca (1942)"));
        fresh.put("deleted_lists",new JSONObject().put("old_extra",true));
        JSONObject oldPhone=library(
            list("atlas","Principal","☐ Casablanca (1942)"),
            list("old_extra","Lista eliminada","☑ Another Film (2010)"));
        SyncStore.Merge result=SyncStore.merge(fresh,new JSONObject(),new JSONObject(),
            peer("android-old",oldPhone,new JSONObject(),new JSONObject()),true);
        assertEquals(1,result.library.getJSONArray("lists").length());
        assertFalse(SyncStore.statuses(result.library).has("another film|2010"));
    }
}
