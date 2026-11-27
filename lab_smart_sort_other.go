//go:build !windows && !darwin

package main

import "net/http"

func injectLabSmartSort(page string) string { return page }
func (a *runningApp) serveLabSmartCache(w http.ResponseWriter,r *http.Request){http.NotFound(w,r)}
