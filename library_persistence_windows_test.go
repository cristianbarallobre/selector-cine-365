package main

import (
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"
)

func TestLibraryLockedJSONKeepsPreviousCommit(t *testing.T) {
	dir := t.TempDir()
	old := initialLibrary("Before")
	a := &runningApp{library: old, libraryPath: filepath.Join(dir, "colecciones.json"), savePath: filepath.Join(dir, "legacy.txt"), token: "test"}
	if err := saveLibrary(a.libraryPath, a.savePath, old); err != nil {
		t.Fatal(err)
	}
	name, err := syscall.UTF16PtrFromString(a.libraryPath)
	if err != nil {
		t.Fatal(err)
	}
	// Allow reads and writes but deny replacement. An in-place fallback
	// would overwrite the locked original despite the failed rename.
	h, err := syscall.CreateFile(name, syscall.GENERIC_READ, syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE, nil, syscall.OPEN_EXISTING, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.CloseHandle(h)
	w := postTestLibrary(t, a, initialLibrary("After"))
	txt, err := os.ReadFile(a.savePath)
	if w.Code != 500 || err != nil || string(txt) != "Before" || !reflect.DeepEqual(a.library, old) || !reflect.DeepEqual(loadLibrary(a.libraryPath, ""), old) {
		t.Fatalf("locked file: HTTP %d TXT=%q err=%v", w.Code, txt, err)
	}
	files, err := filepath.Glob(filepath.Join(dir, ".cine-guardar-*"))
	if err != nil || len(files) != 0 {
		t.Fatalf("temporary files: %v err=%v", files, err)
	}
}
