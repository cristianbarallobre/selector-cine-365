//go:build darwin && cgo

package main

/*
#cgo LDFLAGS: -framework Security -framework Foundation
#include <stdlib.h>
#include <Security/SecBase.h>
int cineKeychainRead(const char *, void **, int *);
int cineKeychainSave(const char *, const void *, int);
int cineKeychainDelete(const char *);
*/
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

type macSessionStore string

func protectedSessionStore(dir string) graphSessionStore { return macSessionStore(dir) }

func (s macSessionStore) load() (graphStoredSession, error) {
	account := C.CString(string(s))
	defer C.free(unsafe.Pointer(account))
	var data unsafe.Pointer
	var length C.int
	status := C.cineKeychainRead(account, &data, &length)
	if status == C.errSecItemNotFound {
		return graphStoredSession{}, os.ErrNotExist
	}
	if status != 0 {
		return graphStoredSession{}, fmt.Errorf("Llavero de macOS: %d", status)
	}
	defer C.free(data)
	raw := unsafe.Slice((*byte)(data), int(length))
	defer zeroBytes(raw)
	return decodeGraphSession(raw)
}

func (s macSessionStore) save(session graphStoredSession) error {
	raw, err := encodeGraphSession(session)
	if err != nil {
		return err
	}
	defer zeroBytes(raw)
	account := C.CString(string(s))
	defer C.free(unsafe.Pointer(account))
	if status := C.cineKeychainSave(account, unsafe.Pointer(&raw[0]), C.int(len(raw))); status != 0 {
		return fmt.Errorf("Llavero de macOS: %d", status)
	}
	return nil
}

func (s macSessionStore) clear() error {
	account := C.CString(string(s))
	defer C.free(unsafe.Pointer(account))
	if status := C.cineKeychainDelete(account); status != 0 && status != C.errSecItemNotFound {
		return fmt.Errorf("Llavero de macOS: %d", status)
	}
	return nil
}
