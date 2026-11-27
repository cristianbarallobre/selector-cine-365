//go:build windows

package main

import (
 "errors"
 "fmt"
 "runtime"
 "syscall"
 "unsafe"
)

const dpapiNoUI = 0x1

type dpapiBlob struct {
 Length uint32
 Data *byte
}

var (
 dpapiCrypt32 = syscall.NewLazyDLL("crypt32.dll")
 dpapiKernel = syscall.NewLazyDLL("kernel32.dll")
 dpapiProtect = dpapiCrypt32.NewProc("CryptProtectData")
 dpapiUnprotect = dpapiCrypt32.NewProc("CryptUnprotectData")
 dpapiLocalFree = dpapiKernel.NewProc("LocalFree")
)

func dpapiInput(data []byte) dpapiBlob {
 if len(data)==0{return dpapiBlob{}}
 return dpapiBlob{Length:uint32(len(data)),Data:&data[0]}
}

// CryptProtectData uses the current Windows user profile, not the computer-wide
// key. The app never writes an unencrypted local refresh token to disk.
func dpapiSeal(plaintext []byte)([]byte,error){
 if len(plaintext)==0||len(plaintext)>32*1024{return nil,errors.New("Credencial local de tamaño no permitido")}
 in:=dpapiInput(plaintext)
 var out dpapiBlob
 description:=syscall.StringToUTF16Ptr("Selector Cine 365 - Credenciales")
 ok,_,callErr:=dpapiProtect.Call(
  uintptr(unsafe.Pointer(&in)),uintptr(unsafe.Pointer(description)),
  0,0,0,dpapiNoUI,uintptr(unsafe.Pointer(&out)))
 runtime.KeepAlive(plaintext)
 if ok==0{return nil,dpapiError("No se pudo proteger la sesión local con Windows",callErr)}
 if out.Data==nil||out.Length==0{return nil,errors.New("Windows no devolvió ninguna credencial cifrada")}
 defer dpapiLocalFree.Call(uintptr(unsafe.Pointer(out.Data)))
 if out.Length>64*1024{return nil,errors.New("Windows devolvió una credencial cifrada demasiado grande")}
 return append([]byte(nil),unsafe.Slice(out.Data,int(out.Length))...),nil
}

func dpapiOpen(ciphertext []byte)([]byte,error){
 if len(ciphertext)==0||len(ciphertext)>64*1024{return nil,errors.New("Credencial local cifrada inválida")}
 in:=dpapiInput(ciphertext)
 var out dpapiBlob
 ok,_,callErr:=dpapiUnprotect.Call(
  uintptr(unsafe.Pointer(&in)),0,0,0,0,dpapiNoUI,uintptr(unsafe.Pointer(&out)))
 runtime.KeepAlive(ciphertext)
 if ok==0{return nil,dpapiError("Windows no pudo recuperar la sesión local de este usuario",callErr)}
 if out.Data==nil||out.Length==0{return nil,errors.New("Windows devolvió una sesión local vacía")}
 defer dpapiLocalFree.Call(uintptr(unsafe.Pointer(out.Data)))
 if out.Length>32*1024{return nil,errors.New("La credencial local descifrada es demasiado grande")}
 return append([]byte(nil),unsafe.Slice(out.Data,int(out.Length))...),nil
}
func dpapiError(prefix string,err error)error{
 if err==nil||errors.Is(err,syscall.Errno(0)){return errors.New(prefix)}
 return fmt.Errorf("%s: %w",prefix,err)
}

