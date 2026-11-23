//go:build !windows

package main

import "net/http"

// La sección Apariencia y sus rutas existen únicamente en Windows LAB.
func injectLabAppearance(page string) string { return page }
func (a *runningApp) serveLabAppearance(w http.ResponseWriter,r *http.Request) {http.NotFound(w,r)}
