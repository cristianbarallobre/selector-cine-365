package main

import (
	"strings"
	"testing"
)

func TestSyncMergeRemoteWatchedAndLocalPending(t *testing.T) {
	local := initialLibrary("☐ 001. Ciudadano Kane (1941)\n☐ 002. Casablanca (1942)")
	remote := initialLibrary("☑ 001. Ciudadano Kane (1941)\n☐ 002. Casablanca (1942)")
	peer := syncSnapshot{Protocol: 2, DeviceID: "remote-test", Library: remote, Known: map[string]int64{}, Revisions: map[string]syncRevision{}}
	merged, _, _, conflicts, err := syncMerge(local, map[string]syncRevision{}, map[string]int64{}, peer)
	if err != nil {
		t.Fatal(err)
	}
	if len(conflicts) != 0 || syncWatched(merged) != 1 {
		t.Fatalf("merge: %v, watched %d", conflicts, syncWatched(merged))
	}
	if !strings.Contains(merged.Lists[0].Text, "☑ 001.") {
		t.Fatal("no se importó la marca de Remote")
	}
}
func TestSyncMergeCausalUnwatch(t *testing.T) {
	l := initialLibrary("☑ 001. Casablanca (1942)")
	remote := initialLibrary("☐ 001. Casablanca (1942)")
	key := "casablanca|1942"
	localRev := syncRevision{Checked: true, Device: "remote", Seq: 1, Clock: map[string]int64{"remote": 1}}
	remoteRev := syncRevision{Checked: false, Device: "desktop", Seq: 2, Clock: map[string]int64{"remote": 1, "desktop": 2}}
	peer := syncSnapshot{Protocol: 2, DeviceID: "desktop", Library: remote, Known: map[string]int64{"remote": 1, "desktop": 2}, Revisions: map[string]syncRevision{key: remoteRev}}
	result, rev, _, conflicts, err := syncMerge(l, map[string]syncRevision{key: localRev}, map[string]int64{"remote": 1}, peer)
	if err != nil {
		t.Fatal(err)
	}
	if len(conflicts) != 0 || syncWatched(result) != 0 || rev[key].Checked {
		t.Fatalf("causal newer unwatch lost: %v, %+v", conflicts, rev)
	}
}
func TestSyncMergeConcurrentOpposite(t *testing.T) {
	l := initialLibrary("☑ 001. Casablanca (1942)")
	remote := initialLibrary("☐ 001. Casablanca (1942)")
	key := "casablanca|1942"
	left := syncRevision{Checked: true, Device: "desktop", Seq: 1, Clock: map[string]int64{"desktop": 1}}
	right := syncRevision{Checked: false, Device: "remote", Seq: 1, Clock: map[string]int64{"remote": 1}}
	p := syncSnapshot{Protocol: 2, DeviceID: "remote", Library: remote, Known: map[string]int64{"remote": 1}, Revisions: map[string]syncRevision{key: right}}
	_, _, _, issues, err := syncMerge(l, map[string]syncRevision{key: left}, map[string]int64{"desktop": 1}, p)
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 || issues[0].Key != key {
		t.Fatalf("missing conflict: %+v", issues)
	}
}
