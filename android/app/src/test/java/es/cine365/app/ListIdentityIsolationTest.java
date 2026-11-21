package es.cine365.app;

import org.json.JSONArray;
import org.json.JSONObject;
import org.junit.Test;
import static org.junit.Assert.*;

public final class ListIdentityIsolationTest {
    private JSONObject list(String id,String name,String text)throws Exception{
        return new JSONObject().put("id",id).put("name",name).put("text",text);
    }
    private JSONObject library(JSONObject... lists)throws Exception{
        JSONArray arr=new JSONArray();for(JSONObject l:lists)arr.put(l);
        return new JSONObject().put("version",1).put("active_id",lists[0].getString("id")).put("lists",arr);
    }
    private JSONObject snapshot(String device,JSONObject library)throws Exception{
        return new JSONObject().put("protocol",2).put("device_id",device).put("library",library);
    }
    private SyncStore.Merge merge(JSONObject local,JSONObject remote)throws Exception{
        return SyncStore.merge(local,new JSONObject(),new JSONObject(),snapshot("desktop-test",remote),false);
    }
    @Test public void defaultStarterCannotSwallowDifferentWindowsList()throws Exception{
        JSONObject android=library(list("atlas","Selector Cine 365 · 365",
            "☐ Casablanca (1942)\n☐ Metropolis (1927)\n☐ Citizen Kane (1941)"));
        JSONObject windows=library(list("atlas","Mi colección personal",
            "☐ Parasite (2019)\n☐ Roma (2018)\n☐ The Matrix (1999)"),
            list("otra","Otra lista","☐ Mulholland Drive (2001)"));
        try{
            merge(android,windows);
            fail("An unrelated Windows list must never be appended to the Android starter");
        }catch(Exception e){
            assertTrue(e.getMessage(),e.getMessage().contains("Sincronización detenida"));
            assertTrue(e.getMessage().contains("atlas"));
        }
        assertFalse("Preview must not mutate local source",android.toString().contains("Parasite"));
    }
    @Test public void sameIdWithSameGenericNameIsStillNotEnough()throws Exception{
        JSONObject a=library(list("mis_pelis","Mis películas","☐ Parasite (2019)\n☐ Roma (2018)"));
        JSONObject b=library(list("mis_pelis","Mis películas","☐ Alien (1979)\n☐ Arrival (2016)"));
        try{merge(a,b);fail("Unrelated lists with same ID/name must be blocked");}
        catch(Exception e){assertTrue(e.getMessage().contains("Sincronización detenida"));}
    }
    @Test public void legitimateSharedListMayAddFilms()throws Exception{
        JSONObject a=library(list("mis_pelis","Favoritas","☐ Parasite (2019)\n☐ Roma (2018)"));
        JSONObject b=library(list("mis_pelis","Favoritas","☐ Parasite (2019)\n☐ Roma (2018)\n☐ Alien (1979)"),
            list("lista_b","Otra lista","☐ Akira (1988)"));
        SyncStore.Merge merged=merge(a,b);
        assertEquals(2,merged.library.getJSONArray("lists").length());
        assertEquals(3,SyncStore.statuses(library(merged.library.getJSONArray("lists").getJSONObject(0))).length());
        assertFalse(merged.library.getJSONArray("lists").getJSONObject(0).getString("text").contains("Akira"));
    }
    @Test public void exampleListAndRenamedListWithSameFilmsCanMerge()throws Exception{
        JSONObject a=library(list("atlas","Selector Cine 365 · 365","☐ Parasite (2019)\n☐ Roma (2018)"));
        JSONObject b=library(list("atlas","Colección de cine","☐ Parasite (2019)\n☐ Roma (2018)"));
        SyncStore.Merge merged=merge(a,b);
        assertEquals("Colección de cine",merged.library.getJSONArray("lists").getJSONObject(0).getString("name"));
        assertEquals(2,SyncStore.statuses(merged.library).length());
    }
    @Test public void pollutedSupersetCannotExpandCorrectList()throws Exception{
        StringBuilder clean=new StringBuilder();
        StringBuilder contaminated=new StringBuilder();
        for(int i=0;i<40;i++){
            String line="☐ Película "+i+" (2000)";
            clean.append(line).append("\n");contaminated.append(line).append("\n");
        }
        for(int i=0;i<40;i++)contaminated.append("☐ Película ajena "+i+" (2010)\n");
        JSONObject windows=library(list("lista","Favoritas",clean.toString()));
        JSONObject android=library(list("lista","Favoritas",contaminated.toString()));
        try{merge(windows,android);fail("A polluted superset must not be accepted as a larger legitimate list");}
        catch(Exception e){assertTrue(e.getMessage(),e.getMessage().contains("Sincronización detenida"));}
        assertEquals(40,SyncStore.statuses(windows).length());
    }

    @Test public void onlyAnUnchangedSampleMayBeReplacedByWindows()throws Exception{
        String sample="☐ Casablanca (1942)\n☐ Metropolis (1927)";
        JSONObject clean=library(list("atlas","Selector Cine 365 · 365",sample));
        assertTrue(SyncStore.isUntouchedStarter(clean,new JSONObject(),new JSONObject(),sample));
        assertFalse(SyncStore.isUntouchedStarter(clean,new JSONObject().put("key",true),new JSONObject(),sample));
        assertFalse(SyncStore.isUntouchedStarter(clean,new JSONObject(),new JSONObject().put("desktop-x",2),sample));
        clean.getJSONArray("lists").getJSONObject(0).put("text",sample+"\n☐ Persona (1966)");
        assertFalse("Real user additions must be preserved",SyncStore.isUntouchedStarter(clean,new JSONObject(),new JSONObject(),sample));
        clean.getJSONArray("lists").getJSONObject(0).put("text",sample);
        clean.getJSONArray("lists").getJSONObject(0).put("name","Mi lista personal");
        assertFalse("Renamed samples are user data",SyncStore.isUntouchedStarter(clean,new JSONObject(),new JSONObject(),sample));
    }
}
