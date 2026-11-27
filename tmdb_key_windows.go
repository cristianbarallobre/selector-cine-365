//go:build windows
package main
func readTMDBKey(path string)(string,error){
 s,e:=readTMDBSettings(path);if e!=nil{return "",e}
 if len(s.Protected)>0{raw,e:=dpapiOpen(s.Protected);if e!=nil{return "",e};defer zeroBytes(raw);return string(raw),nil}
 if s.Key!=""{if e=writeTMDBKey(path,s.Key);e!=nil{diagnosticError("tmdb","proteger credencial antigua",e)}}
 return s.Key,nil
}
func writeTMDBKey(path,key string)error{
 if key==""{return writeTMDBSettings(path,platformSettings{})}
 raw:=[]byte(key);defer zeroBytes(raw)
 sealed,e:=dpapiSeal(raw);if e!=nil{return e}
 return writeTMDBSettings(path,platformSettings{Protected:sealed})
}
