package diff

import (
	"testing"
)

func TestKind_String(t *testing.T) {
	cases := []struct {
		kind Kind
		want string
	}{
		{MissingInRight, "missing_in_right"},
		{MissingInLeft, "missing_in_left"},
		{ValueMismatch, "value_mismatch"},
		{Kind(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.kind.String(); got != tc.want {
			t.Errorf("Kind(%d).String() = %q, want %q", tc.kind, got, tc.want)
		}
	}
}

func TestResult_HasDifferences(t *testing.T) {
	empty := Result{}
	if empty.HasDifferences() {
		t.Error("expected HasDifferences() = false for empty result")
	}

	withEntries := Result{
		Entries: []Entry{
			{Key: "FOO", Kind: MissingInRight},
		},
	}
	if !withEntries.HasDifferences() {
		t.Error("expected HasDifferences() = true for result with entries")
	}
}

func TestResult_ByKind(t *testing.T) {
	r := Result{
		Entries: []Entry{
			{Key: "A", Kind: MissingInRight},
			{Key: "B", Kind: MissingInLeft},
			{Key: "C", Kind: ValueMismatch},
			{Key: "D", Kind: MissingInRight},
		},
	}

	missRight := r.ByKind(MissingInRight)
	if len(missRight) != 2 {
		t.Errorf("ByKind(MissingInRight) returned %d entries, want 2", len(missRight))
	}

	missLeft := r.ByKind(MissingInLeft)
	if len(missLeft) != 1 {
		t.Errorf("ByKind(MissingInLeft) returned %d entries, want 1", len(missLeft))
	}

	mismatch := r.ByKind(ValueMismatch)
	if len(mismatch) != 1 {
		t.Errorf("ByKind(ValueMismatch) returned %d entries, want 1", len(mismatch))
	}

	none := r.ByKind(Kind(99))
	if len(none) != 0 {
		t.Errorf("ByKind(unknown) returned %d entries, want 0", len(none))
	}
}
