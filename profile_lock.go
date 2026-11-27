package main

import (
	"errors"
	"os"
	"path/filepath"
)

var errProfileInUse = errors.New("Selector Cine 365 ya está abierto con este perfil")

func lockProfile(dir string) (func(), error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return lockProfileFile(filepath.Join(dir, ".selector-cine.lock"))
}
