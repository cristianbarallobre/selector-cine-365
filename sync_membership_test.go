package main

import (
	"strings"
	"testing"
)

func TestSyncMembershipRemovalSurvivesOldSnapshotAndReadd(t *testing.T) {
	before := initialLibrary("☐ Maridos (1970)\n☐ Insomnio (2002)")
	before.Lists = append(before.Lists, filmCollection{ID: "other", Name: "Otra lista", Text: "☐ Maridos (1970)"})
	removed := before
	removed.Lists = append([]filmCollection(nil), before.Lists...)
	removed.Lists[0].Text = "☐ Insomnio (2002)"
	recordMembershipEdits(before, &removed)

	// Another device has an old copy of Maridos and a legitimate new film.
	peer := before
	peer.Lists = append([]filmCollection(nil), before.Lists...)
	peer.Lists[0].Text += "\n☐ Origen (2010)"
	merged, _, _, _, err := syncMergeVersioned(removed, nil, nil, syncSnapshot{Protocol: 2, Library: peer}, true)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(merged.Lists[0].Text, "Maridos") || !strings.Contains(merged.Lists[0].Text, "Origen") ||
		!strings.Contains(merged.Lists[1].Text, "Maridos") {
		t.Fatal("La fusión resucitó la pertenencia retirada o perdió otros cambios", merged.Lists)
	}

	// A device on the updated version can receive the removal in the opposite direction.
	received, _, _, _, err := syncMergeVersioned(peer, nil, nil, syncSnapshot{Protocol: 2, Library: merged}, true)
	if err != nil || strings.Contains(received.Lists[0].Text, "Maridos") {
		t.Fatal("La retirada no llegó al otro equipo", err)
	}

	readded := merged
	readded.Lists = append([]filmCollection(nil), merged.Lists...)
	readded.Lists[0].Text += "\n☐ Maridos (1970)"
	recordMembershipEdits(merged, &readded)
	received, _, _, _, err = syncMergeVersioned(readded, nil, nil, syncSnapshot{Protocol: 2, Library: merged}, true)
	if err != nil || !strings.Contains(received.Lists[0].Text, "Maridos") {
		t.Fatal("Una retirada vieja sobrescribió la nueva incorporación", err)
	}
}

