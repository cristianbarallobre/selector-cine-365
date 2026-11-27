//go:build !windows && !darwin

package main

// Los nuevos filtros y estadísticas solo se activan en Windows LAB.
func injectLabLibraryFilters(page string) string {return page}
