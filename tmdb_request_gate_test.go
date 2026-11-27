package main
import("testing";"net/http";"context";"time")
func TestSharedTMDBCooldown(t *testing.T){
 c:=&http.Client{}
 deferTMDBRequests(c,"2")
 if _,e:=acquireTMDBRequest(context.Background(),c);e==nil{t.Fatal("cooldown ignored")}
 other:=&http.Client{}
 release,e:=acquireTMDBRequest(context.Background(),other);if e!=nil{t.Fatal(e)};release()
 for i:=0;i<cap(tmdbSlots);i++{tmdbSlots<-struct{}{}}
 defer func(){for i:=0;i<cap(tmdbSlots);i++{<-tmdbSlots}}()
 ctx,cancel:=context.WithTimeout(context.Background(),time.Millisecond);defer cancel()
 if _,e:=acquireTMDBRequest(ctx,other);e==nil{t.Fatal("request ignored deadline")}
}
