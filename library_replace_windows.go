//go:build windows

package main

import (
 "path/filepath"
 "runtime"
 "strings"
 "syscall"
 "unsafe"
)

var libraryMoveFileEx=syscall.NewLazyDLL("kernel32.dll").NewProc("MoveFileExW")

func libraryWindowsPath(path string)(*uint16,error){
 absolute,err:=filepath.Abs(path);if err!=nil{return nil,err}
 extended:="\\\\?\\"
 unc:="\\\\"
 if !strings.HasPrefix(absolute,extended)&&len(absolute)>=240{
  if strings.HasPrefix(absolute,unc){absolute=extended+"UNC\\"+strings.TrimPrefix(absolute,unc)}else{absolute=extended+absolute}
 }
 return syscall.UTF16PtrFromString(absolute)
}

// No fallback that truncates the existing destination. Request write-through
// for the same-volume rename; the temporary file was already flushed.
func replaceLibraryDestination(source,destination string)error{
 from,err:=libraryWindowsPath(source);if err!=nil{return err}
 to,err:=libraryWindowsPath(destination);if err!=nil{return err}
 result,_,callErr:=libraryMoveFileEx.Call(uintptr(unsafe.Pointer(from)),uintptr(unsafe.Pointer(to)),0x1|0x8)
 runtime.KeepAlive(from);runtime.KeepAlive(to)
 if result==0{return callErr}
 return nil
}
