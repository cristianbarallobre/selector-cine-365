package main

import (
	"errors"
	"testing"
)

func TestProfileLockBlocksSecondOpenerAndReleases(t *testing.T) {
	dir := t.TempDir()
	release, err := lockProfile(dir)
	if err != nil {
		t.Fatal(err)
	}
	second, err := lockProfile(dir)
	release()
	if err == nil {
		second()
		t.Fatal("the profile had two owners")
	}
	if !errors.Is(err, errProfileInUse) {
		t.Fatal(err)
	}
	release, err = lockProfile(dir)
	if err != nil {
		t.Fatal(err)
	}
	release()
}
