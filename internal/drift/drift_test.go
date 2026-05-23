package drift_test

import (
	"testing"
	"time"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/drift"
	"github.com/user/envdiff/internal/snapshot"
)

func makeSnapshot(pairs map[string]string) *snapshot.Snapshot {
	var results []diff.Result
	for k, v := range pairs {
		results = append(results, diff.Result{Key: k, Left: v, Right: v})
	}
	return &snapshot.Snapshot{CreatedAt: time.Now(), Results: results}
}

func TestDetect_NoDrift(t *testing.T) {
	snap := makeSnapshot(map[string]string{"FOO": "bar", "BAZ": "qux"})
	current := map[string]string{"FOO": "bar", "BAZ": "qux"}

	changes := drift.Detect(snap, current)
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %d", len(changes))
	}
}

func TestDetect_AddedKey(t *testing.T) {
	snap := makeSnapshot(map[string]string{"FOO": "bar"})
	current := map[string]string{"FOO": "bar", "NEW_KEY": "value"}

	changes := drift.Detect(snap, current)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != drift.Added || changes[0].Key != "NEW_KEY" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDetect_RemovedKey(t *testing.T) {
	snap := makeSnapshot(map[string]string{"FOO": "bar", "OLD_KEY": "gone"})
	current := map[string]string{"FOO": "bar"}

	changes := drift.Detect(snap, current)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != drift.Removed || changes[0].Key != "OLD_KEY" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
	if changes[0].OldVal != "gone" {
		t.Errorf("expected OldVal=gone, got %q", changes[0].OldVal)
	}
}

func TestDetect_ModifiedKey(t *testing.T) {
	snap := makeSnapshot(map[string]string{"FOO": "old"})
	current := map[string]string{"FOO": "new"}

	changes := drift.Detect(snap, current)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	c := changes[0]
	if c.Kind != drift.Modified {
		t.Errorf("expected Modified, got %s", c.Kind)
	}
	if c.OldVal != "old" || c.NewVal != "new" {
		t.Errorf("unexpected values: old=%q new=%q", c.OldVal, c.NewVal)
	}
}

func TestDetect_NilSnapshot(t *testing.T) {
	changes := drift.Detect(nil, map[string]string{"FOO": "bar"})
	if changes != nil {
		t.Errorf("expected nil for nil snapshot, got %v", changes)
	}
}

func TestChangeKind_String(t *testing.T) {
	cases := []struct {
		kind drift.ChangeKind
		want string
	}{
		{drift.Added, "added"},
		{drift.Removed, "removed"},
		{drift.Modified, "modified"},
		{drift.ChangeKind(99), "unknown(99)"},
	}
	for _, tc := range cases {
		if got := tc.kind.String(); got != tc.want {
			t.Errorf("ChangeKind(%d).String() = %q, want %q", tc.kind, got, tc.want)
		}
	}
}
