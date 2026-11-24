//go:build !windows

package main

import "net/http"

func injectLabMaintenance(page string) string{return page}
func (a *runningApp) serveLabMaintenance(w http.ResponseWriter,r *http.Request){http.NotFound(w,r)}
