package main

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

type syncMembershipRevision struct {
	Present   bool  `json:"present"`
	ChangedAt int64 `json:"changed_at"`
}

type syncMembershipChanges map[string]map[string]syncMembershipRevision

func validateMembershipChanges(changes syncMembershipChanges) error {
	for id, entries := range changes {
		if !listIDPattern.MatchString(id) {
			return fmt.Errorf("identificador inválido en la pertenencia a listas")
		}
		for key, revision := range entries {
			if len(key) > 1024 || !utf8.ValidString(key) || !strings.Contains(key, "|") ||
				revision.ChangedAt <= 0 || revision.ChangedAt > 4102444800000 {
				return fmt.Errorf("revisión inválida en la pertenencia a listas")
			}
		}
	}
	return nil
}

func mergeMembershipChanges(local, remote syncMembershipChanges) syncMembershipChanges {
	if len(local) == 0 && len(remote) == 0 {
		return nil
	}
	result := make(syncMembershipChanges, len(local)+len(remote))
	for _, changes := range []syncMembershipChanges{local, remote} {
		for id, entries := range changes {
			if result[id] == nil {
				result[id] = make(map[string]syncMembershipRevision, len(entries))
			}
			for key, next := range entries {
				prior, exists := result[id][key]
				if !exists || next.ChangedAt > prior.ChangedAt ||
					next.ChangedAt == prior.ChangedAt && prior.Present && !next.Present {
					result[id][key] = next
				}
			}
		}
	}
	return result
}

func recordMembershipEdits(before libraryState, after *libraryState) {
	changes := mergeMembershipChanges(before.MembershipChanges, after.MembershipChanges)
	previous := make(map[string]filmCollection, len(before.Lists))
	for _, list := range before.Lists {
		previous[list.ID] = list
	}
	now := time.Now().UnixMilli()
	for _, list := range after.Lists {
		old, exists := previous[list.ID]
		if !exists || old.Text == list.Text {
			continue
		}
		prior, current := syncListFilmKeys(old), syncListFilmKeys(list)
		entries := changes[list.ID]
		for key := range prior {
			if !current[key] && len(key) <= 1024 {
				if entries == nil {
					entries = map[string]syncMembershipRevision{}
				}
				entries[key] = syncMembershipRevision{Present: false, ChangedAt: max(now, entries[key].ChangedAt+1)}
			}
		}
		for key := range current {
			if revision, tracked := entries[key]; tracked && !prior[key] {
				entries[key] = syncMembershipRevision{Present: true, ChangedAt: max(now, revision.ChangedAt+1)}
			}
		}
		if len(entries) != 0 {
			if changes == nil {
				changes = syncMembershipChanges{}
			}
			changes[list.ID] = entries
		}
	}
	for id := range after.DeletedLists {
		delete(changes, id)
	}
	after.MembershipChanges = changes
}

func filterRemovedMemberships(list filmCollection, changes syncMembershipChanges) filmCollection {
	revisions := changes[list.ID]
	if len(revisions) == 0 {
		return list
	}
	lines := strings.Split(list.Text, "\n")
	kept := lines[:0]
	for _, line := range lines {
		film, valid := parseSyncFilm(line)
		if revision, tracked := revisions[film.key]; valid && tracked && !revision.Present {
			continue
		}
		kept = append(kept, line)
	}
	list.Text = strings.Join(kept, "\n")
	return list
}
