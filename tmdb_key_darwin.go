//go:build darwin && cgo
package main
/*
#cgo LDFLAGS: -framework Security -framework Foundation
#include <stdlib.h>
#include <Security/SecBase.h>
int cineTMDBKeyRead(const char *, void **, int *);
int cineTMDBKeySave(const char *, const void *, int);
*/
import "C"
import("fmt";"unsafe";"crypto/sha256";"path/filepath")
func tmdbKeyAccount(path string)string{return fmt.Sprintf("%x",sha256.Sum256([]byte(filepath.Clean(path))))}
func readTMDBKeychain(path string)(string,error){
 account:=C.CString(tmdbKeyAccount(path));defer C.free(unsafe.Pointer(account))
 var data unsafe.Pointer;var length C.int
 status:=C.cineTMDBKeyRead(account,&data,&length)
 if status!=0{return "",fmt.Errorf("Llavero de macOS: %d",status)}
 defer C.free(data)
 raw:=unsafe.Slice((*byte)(data),int(length));defer zeroBytes(raw)
 return string(raw),nil
}
func readTMDBKey(path string)(string,error){
 s,e:=readTMDBSettings(path);if e!=nil{return "",e}
 if s.Keychain{return readTMDBKeychain(path)}
 if s.Key!=""{if e=writeTMDBKey(path,s.Key);e!=nil{diagnosticError("tmdb","proteger credencial antigua",e)}}
 return s.Key,nil
}
func writeTMDBKey(path,key string)error{
 if key==""{return writeTMDBSettings(path,platformSettings{})}
 raw:=[]byte(key);defer zeroBytes(raw)
 account:=C.CString(tmdbKeyAccount(path));defer C.free(unsafe.Pointer(account))
 status:=C.cineTMDBKeySave(account,unsafe.Pointer(&raw[0]),C.int(len(raw)))
 if status!=0{return fmt.Errorf("Llavero de macOS: %d",status)}
 if prior,e:=readTMDBSettings(path);e==nil&&prior.Keychain&&prior.Key==""{return nil}
 return writeTMDBSettings(path,platformSettings{Keychain:true})
}
