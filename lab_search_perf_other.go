//go:build !windows && !darwin

package main

func injectLabSearchPerformance(page string) string{return page}
