package main
import("encoding/json";"os")
func readTMDBSettings(path string)(platformSettings,error){
 var s platformSettings
 data,e:=os.ReadFile(path);if e!=nil{return s,e}
 e=json.Unmarshal(data,&s);return s,e
}
func writeTMDBSettings(path string,s platformSettings)error{
 data,e:=json.Marshal(s);if e!=nil{return e}
 return replaceLibraryFile(path,data)
}
