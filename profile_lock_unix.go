//go:build !windows

package main

import (
	"errors"
	"os"
	"syscall"
)

func lockProfileFile(path string) (func(), error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = file.Close()
		if errors.Is(err, syscall.EWOULDBLOCK) {
			return nil, errProfileInUse
		}
		return nil, err
	}
	return func() { _ = file.Close() }, nil
}
