//go:build !windows && !darwin

package main

// En otras plataformas la interfaz estable permanece sin cambios.
func injectLabReactions(page string) string{return page}
