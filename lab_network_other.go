//go:build !windows && !darwin

package main

import "net/http"

// Los diagnósticos de red y el refuerzo de fetch son exclusivos de Windows LAB.
func injectLabNetwork(page string) string{return page}
func (a *runningApp) serveLabNetworkLog(w http.ResponseWriter,r *http.Request){http.NotFound(w,r)}
