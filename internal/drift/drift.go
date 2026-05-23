// Package drift detects when a .env file has changed relative to a
// previously recorded snapshot, reporting which keys were added, removed,
// or modified between two points in time.
package drift

import (
	"fmt"
	"sort"

	"github.com/user/envdiff/internal/snapshot"
)

// Change describes a single key-level change detected between a snapshot
// and the current environment map.
type Change struct {
	Key    string
	Kind   ChangeKind
	OldVal string
	NewVal string
}

// ChangeKind classifies the type of drift observed for a key.
type ChangeKind int

const (
	Added    ChangeKind = iota // key exists now but not in snapshot
	Removed                    // key existed in snapshot but not now
	Modified                   // key exists in both but value differs
)

// String returns a human-readable label for the ChangeKind.
func (k ChangeKind) String() string {
	switch k {
	case Added:
		return "added"
	case Removed:
		return "removed"
	case Modified:
		return "modified"
	default:
		return fmt.Sprintf("unknown(%d)", int(k))
	}
}

// Detect compares a saved snapshot against the current env map and returns
// the list of changes. An empty slice means no drift was detected.
func Detect(snap *snapshot.Snapshot, current map[string]string) []Change {
	if snap == nil {
		return nil
	}

	var changes []Change

	// Build a lookup from the snapshot results.
	snapshotVals := make(map[string]string, len(snap.Results))
	for _, r := range snap.Results {
		snapshotVals[r.Key] = r.Left
	}

	// Keys removed or modified.
	for key, oldVal := range snapshotVals {
		newVal, exists := current[key]
		if !exists {
			changes = append(changes, Change{Key: key, Kind: Removed, OldVal: oldVal})
		} else if newVal != oldVal {
			changes = append(changes, Change{Key: key, Kind: Modified, OldVal: oldVal, NewVal: newVal})
		}
	}

	// Keys added since snapshot.
	for key, newVal := range current {
		if _, existed := snapshotVals[key]; !existed {
			changes = append(changes, Change{Key: key, Kind: Added, NewVal: newVal})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		if changes[i].Kind != changes[j].Kind {
			return changes[i].Kind < changes[j].Kind
		}
		return changes[i].Key < changes[j].Key
	})

	return changes
}
