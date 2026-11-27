//go:build !windows && (!darwin || !cgo)
package main
// Compatibility for source builds on unsupported systems. Official desktop
// releases use native protection.
func readTMDBKey(path string)(string,error){s,e:=readTMDBSettings(path);return s.Key,e}
func writeTMDBKey(path,key string)error{return writeTMDBSettings(path,platformSettings{Key:key})}
