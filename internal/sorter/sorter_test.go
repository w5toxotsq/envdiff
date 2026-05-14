package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/sorter"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "ZEBRA", Kind: diff.KindMismatch, LeftVal: "a", RightVal: "b"},
		{Key: "ALPHA", Kind: diff.KindMissingInRight, LeftVal: "x"},
		{Key: "MANGO", Kind: diff.KindMissingInLeft, RightVal: "y"},
		{Key: "BETA", Kind: diff.KindMismatch, LeftVal: "c", RightVal: "d"},
	}
}

func TestSort_ByKey(t *testing.T) {
	results := makeResults()
	opts := sorter.Options{Order: sorter.ByKey}
	sorted := sorter.Sort(results, opts)

	expectedKeys := []string{"ALPHA", "BETA", "MANGO", "ZEBRA"}
	for i, key := range expectedKeys {
		if sorted[i].Key != key {
			t.Errorf("index %d: expected key %q, got %q", i, key, sorted[i].Key)
		}
	}
}

func TestSort_ByKind(t *testing.T) {
	results := makeResults()
	opts := sorter.Options{Order: sorter.ByKind}
	sorted := sorter.Sort(results, opts)

	for i := 1; i < len(sorted); i++ {
		if sorted[i].Kind < sorted[i-1].Kind {
			t.Errorf("results not sorted by kind at index %d: %v before %v", i, sorted[i-1].Kind, sorted[i].Kind)
		}
		if sorted[i].Kind == sorted[i-1].Kind && sorted[i].Key < sorted[i-1].Key {
			t.Errorf("results with same kind not sorted by key at index %d", i)
		}
	}
}

func TestSort_DefaultOptions(t *testing.T) {
	opts := sorter.DefaultOptions()
	if opts.Order != sorter.ByKey {
		t.Errorf("expected default order ByKey, got %v", opts.Order)
	}
}

func TestSort_EmptySlice(t *testing.T) {
	results := []diff.Result{}
	sorted := sorter.Sort(results, sorter.DefaultOptions())
	if len(sorted) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(sorted))
	}
}

func TestSort_OriginalUnmodified(t *testing.T) {
	results := makeResults()
	originalFirst := results[0].Key
	sorter.Sort(results, sorter.DefaultOptions())
	if results[0].Key != originalFirst {
		t.Errorf("original slice was modified: expected first key %q, got %q", originalFirst, results[0].Key)
	}
}
