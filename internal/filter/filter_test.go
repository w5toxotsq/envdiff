package filter_test

import (
	"testing"

	"github.com/yourusername/envdiff/internal/diff"
	"github.com/yourusername/envdiff/internal/filter"
)

func makeResult(entries ...diff.Entry) diff.Result {
	return diff.Result{Entries: entries}
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	input := makeResult(
		diff.Entry{Key: "FOO", Kind: diff.KindMissingRight},
		diff.Entry{Key: "BAR", Kind: diff.KindMismatch},
	)
	out := filter.Apply(input, filter.Options{})
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	input := makeResult(
		diff.Entry{Key: "DB_HOST", Kind: diff.KindMismatch},
		diff.Entry{Key: "APP_PORT", Kind: diff.KindMissingRight},
		diff.Entry{Key: "DB_PORT", Kind: diff.KindMissingLeft},
	)
	out := filter.Apply(input, filter.Options{Prefix: "DB_"})
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
	for _, e := range out.Entries {
		if e.Key != "DB_HOST" && e.Key != "DB_PORT" {
			t.Errorf("unexpected key %q", e.Key)
		}
	}
}

func TestApply_OnlyMissing(t *testing.T) {
	input := makeResult(
		diff.Entry{Key: "FOO", Kind: diff.KindMissingRight},
		diff.Entry{Key: "BAR", Kind: diff.KindMismatch},
		diff.Entry{Key: "BAZ", Kind: diff.KindMissingLeft},
	)
	out := filter.Apply(input, filter.Options{OnlyMissing: true})
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestApply_OnlyMismatched(t *testing.T) {
	input := makeResult(
		diff.Entry{Key: "FOO", Kind: diff.KindMissingRight},
		diff.Entry{Key: "BAR", Kind: diff.KindMismatch},
		diff.Entry{Key: "BAZ", Kind: diff.KindMismatch},
	)
	out := filter.Apply(input, filter.Options{OnlyMismatched: true})
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestApply_EmptyResult(t *testing.T) {
	out := filter.Apply(makeResult(), filter.Options{Prefix: "DB_"})
	if len(out.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(out.Entries))
	}
}

func TestApply_PrefixCaseInsensitive(t *testing.T) {
	input := makeResult(
		diff.Entry{Key: "DB_HOST", Kind: diff.KindMismatch},
	)
	out := filter.Apply(input, filter.Options{Prefix: "db_"})
	if len(out.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out.Entries))
	}
}
