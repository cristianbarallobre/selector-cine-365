//go:build !windows

package main

import (
 "os"
 "path/filepath"
)

// Flush the directory entry after the atomic rename where supported. Once
// renamed, the commit is already visible; a directory flush failure cannot
// be reported as if the old library were still present.
func replaceLibraryDestination(source,destination string)error{
 if err:=os.Rename(source,destination);err!=nil{return err}
 dir,err:=os.Open(filepath.Dir(destination))
 if err==nil{err=dir.Sync();_ = dir.Close()}
 if err!=nil{diagnosticError("guardado","confirmar directorio tras sustitucion",err)}
 return nil
}
