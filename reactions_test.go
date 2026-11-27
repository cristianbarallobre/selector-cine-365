package main

import "testing"

func reactionLibrary() libraryState {
 return libraryState{Version:1,ActiveID:"atlas",Lists:[]filmCollection{{ID:"atlas",Name:"Lista",Text:"☑ Ejemplo (2024)\n☐ Otra (2025)"}}}
}
func TestRatingsAndFavoritesAreValidated(t *testing.T) {
 l:=reactionLibrary()
 l.Reactions=map[string]filmReaction{"ejemplo|2024":{Favorite:true,Rating:4,Title:"Ejemplo",Year:"2024",ChangedAt:1000,TMDBID:100}}
 if err:=validateLibrary(l);err!=nil{t.Fatalf("Valoración válida rechazada: %v",err)}
 r:=l.Reactions["ejemplo|2024"]
 r.Rating=6;l.Reactions["ejemplo|2024"]=r
 if err:=validateLibrary(l);err==nil{t.Fatal("Se aceptó una nota superior a 5")}
 r.Rating=0;r.Favorite=false;l.Reactions["ejemplo|2024"]=r
 if err:=validateLibrary(l);err!=nil{t.Fatalf("Una valoración borrada debe conservar tombstone: %v",err)}
 r.TMDBID=1000000000;l.Reactions["ejemplo|2024"]=r
 if err:=validateLibrary(l);err==nil{t.Fatal("Se aceptó ID TMDB inválido")}
}
func TestReactionSyncNewerChangedAtWins(t *testing.T) {
 local:=reactionLibrary()
 remote:=reactionLibrary()
 local.Reactions=map[string]filmReaction{"ejemplo|2024":{Favorite:false,Rating:0,Title:"Ejemplo",Year:"2024",ChangedAt:300}}
 remote.Reactions=map[string]filmReaction{"ejemplo|2024":{Favorite:true,Rating:5,Title:"Ejemplo",Year:"2024",ChangedAt:200},"otra|2025":{Favorite:true,Rating:0,Title:"Otra",Year:"2025",ChangedAt:400}}
 result,_,_,_,err:=syncMergeVersioned(local,nil,nil,syncSnapshot{Protocol:syncProtocol,Library:remote},false)
 if err!=nil{t.Fatal(err)}
 if result.Reactions["ejemplo|2024"].Favorite||result.Reactions["ejemplo|2024"].Rating!=0{t.Fatal("La copia antigua resucitó un corazón o una nota eliminada")}
 if !result.Reactions["otra|2025"].Favorite{t.Fatal("No se incorporó un favorito nuevo del otro equipo")}
 // El registro más reciente de otro dispositivo sí debe sustituir al local.
 remote.Reactions["ejemplo|2024"]=filmReaction{Favorite:true,Rating:3,Title:"Ejemplo",Year:"2024",ChangedAt:450}
 result,_,_,_,err=syncMergeVersioned(local,nil,nil,syncSnapshot{Protocol:syncProtocol,Library:remote},false)
 if err!=nil{t.Fatal(err)}
 if !result.Reactions["ejemplo|2024"].Favorite||result.Reactions["ejemplo|2024"].Rating!=3{t.Fatal("El cambio cronológicamente posterior no llegó")}
 // Una copia antigua que no conoce reacciones no puede vaciar el estado.
 remote.Reactions=nil
 result,_,_,_,err=syncMergeVersioned(local,nil,nil,syncSnapshot{Protocol:syncProtocol,Library:remote},false)
 if err!=nil{t.Fatal(err)}
 if _,ok:=result.Reactions["ejemplo|2024"];!ok{t.Fatal("Un cliente antiguo borró la valoración")}
}

func TestFavoritesOrderMergeKeepsNewerOrder(t *testing.T) {
 local:=initialLibrary("☐ Ejemplo (2024)")
 remote:=initialLibrary("☐ Ejemplo (2024)")
 local.FavoritesOrder=&favoritesOrder{Keys:[]string{"ejemplo|2024","otra|2025"},ChangedAt:200}
 remote.FavoritesOrder=&favoritesOrder{Keys:[]string{"otra|2025","ejemplo|2024"},ChangedAt:100}
 peer:=syncSnapshot{Protocol:syncProtocol,Library:remote}
 older,_,_,_,err:=syncMerge(local,nil,nil,peer)
 if err!=nil||older.FavoritesOrder==nil||older.FavoritesOrder.Keys[0]!="ejemplo|2024"{t.Fatalf("El orden antiguo sustituyó al reciente: %v",err)}
 peer.Library.FavoritesOrder.ChangedAt=300
 newer,_,_,_,err:=syncMerge(local,nil,nil,peer)
 if err!=nil||newer.FavoritesOrder==nil||newer.FavoritesOrder.Keys[0]!="otra|2025"{t.Fatalf("No se aplicó el orden más reciente: %v",err)}
}
