//go:build !windows && !darwin

package main

import "errors"

func chooseBackupFolder() (string, error) { return "", errors.New("Solo disponible en Windows") }
