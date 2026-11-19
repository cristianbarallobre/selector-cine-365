//go:build !windows && !darwin

package main

import "errors"

func chooseCloudFolder() (string, error) { return "", errors.New("Solo disponible en Windows") }
