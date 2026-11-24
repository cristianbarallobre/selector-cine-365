//go:build !windows && !darwin

package main

import "net/http"

// La función Descubrir LAB solo existe en Windows. macOS y Linux no reciben interfaz ni API.
func injectLabDiscovery(page string) string {return page}
func (a *runningApp) serveLabDiscovery(w http.ResponseWriter,r *http.Request) {http.NotFound(w,r)}
