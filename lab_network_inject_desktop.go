//go:build windows || darwin

package main

import (
	_ "embed"
	"strings"
)

//go:embed lab_network_windows.js
var labNetworkFrontend string

func injectLabNetwork(page string) string {
	return strings.Replace(page, "/* CINE_NETWORK_INTERCEPTOR */", labNetworkFrontend, 1)
}
