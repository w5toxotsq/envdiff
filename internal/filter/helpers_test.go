package filter_test

import (
	"testing"

	"github.com/yourusername/envdiff/internal/diff"
	"github.com/yourusername/envdiff/internal/filter"
)

// TestApply_OriginalUnmodified ensures Apply does not mutate the source Result.
func TestApply_OriginalUnmodified(t *testing.T) {
	original := makeResult(
		diff.Entry{Key: "FOO", Kind: diff.KindMissingRight},
		diff.Entry{Key: "BAR", Kind: diff.KindMismatch},
	)

	_ = filter.Apply(original, filter.Options{OnlyMissing: true})

	if len(original.Entries) != 2 {
		t.Fatalf("original result was mutated: expected 2 entries, got %d", len(original.Entries))
	}
}

// TestApply_CombinedPrefixAndOnlyMismatched verifies both filters apply together.
func TestApply_CombinedPrefixAndOnlyMismatched(t *testing.T) {
	input := makeResult(
		diff.Entry{Key: "DB_HOST", Kind: diff.KindMismatch},
		diff.Entry{Key: "DB_USER", Kind: diff.KindMissingRight},
		diff.Entry{Key: "APP_KEY", Kind: diff.KindMismatch},
	)
	out := filter.Apply(input, filter.Options{
		Prefix:         "DB_",
		OnlyMismatched: true,
	})
	if len(out.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out.Entries))
	}
	if out.Entries[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", out.Entries[0].Key)
	}
}
