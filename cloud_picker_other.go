//go:build !windows

package main

import "errors"

func chooseCloudFolder() (string, error) { return "", errors.New("Solo disponible en Windows") }
