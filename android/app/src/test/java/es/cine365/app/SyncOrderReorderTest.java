package es.cine365.app;

import org.json.JSONArray;
import org.json.JSONObject;
import org.junit.Test;
import static org.junit.Assert.*;

public class SyncOrderReorderTest {
 @Test public void newerOrderIsAppliedAndOlderOrderCannotUndoIt() throws Exception {
  JSONObject local=library("☐ Uno (2001)\n☐ Dos (2002)\n☐ Tres (2003)",100);
  JSONObject remote=library("☐ Tres (2003)\n☐ Uno (2001)\n☐ Dos (2002)",200);
  JSONObject snapshot=new JSONObject().put("protocol",2).put("device_id","desktop").put("library",remote)
     .put("revisions",new JSONObject()).put("known",new JSONObject());
  SyncStore.Merge result=SyncStore.merge(local,new JSONObject(),new JSONObject(),snapshot);
  assertEquals("☐ Tres (2003)\n☐ Uno (2001)\n☐ Dos (2002)",result.library.getJSONArray("lists").getJSONObject(0).getString("text"));
  assertEquals(200,result.library.getJSONArray("lists").getJSONObject(0).getLong("order_changed_at"));
  JSONObject stale=library("☐ Dos (2002)\n☐ Uno (2001)\n☐ Cuatro (2004)",150);
  snapshot.put("library",stale);
  SyncStore.Merge after=SyncStore.merge(result.library,new JSONObject(),new JSONObject(),snapshot);
  String text=after.library.getJSONArray("lists").getJSONObject(0).getString("text");
  assertTrue(text.startsWith("☐ Tres (2003)\n☐ Uno (2001)\n☐ Dos (2002)"));
  assertTrue(text.contains("Cuatro (2004)"));
 }
 @Test public void unnumberedHeadingsKeepTheirLocations() throws Exception {
  String text="Cabecera\n☐ 001. Uno (2001)\n\n☑ 002. Dos (2002)";
  String other="☑ 002. Dos (2002)\n☐ 001. Uno (2001)";
  String moved=SyncStore.orderedFilmRows(text,other);
  assertEquals("Cabecera\n☑ 002. Dos (2002)\n\n☐ 001. Uno (2001)",moved);
 }
 private static JSONObject library(String text,long timestamp) throws Exception {
  return new JSONObject().put("version",1).put("active_id","atlas")
    .put("lists",new JSONArray().put(new JSONObject().put("id","atlas").put("name","Prueba")
    .put("text",text).put("order_changed_at",timestamp)));
 }
}
