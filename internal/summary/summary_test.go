package summary_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/summary"
)

func makeResults(kinds ...diff.Kind) []diff.Result {
	results := make([]diff.Result, 0, len(kinds))
	for i, k := range kinds {
		results = append(results, diff.Result{
			Key:  fmt.Sprintf("KEY_%d", i),
			Kind: k,
		})
	}
	return results
}

func TestCompute_Empty(t *testing.T) {
	s := summary.Compute(nil)
	if s.Total != 0 || s.Missing != 0 || s.Extra != 0 || s.Mismatched != 0 {
		t.Errorf("expected zero stats, got %+v", s)
	}
	if s.HasDifferences() {
		t.Error("expected HasDifferences to be false for empty results")
	}
}

func TestCompute_AllKinds(t *testing.T) {
	results := []diff.Result{
		{Key: "A", Kind: diff.KindMissingInRight},
		{Key: "B", Kind: diff.KindMissingInLeft},
		{Key: "C", Kind: diff.KindValueMismatch},
		{Key: "D", Kind: diff.KindMissingInRight},
	}
	s := summary.Compute(results)
	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if s.Missing != 2 {
		t.Errorf("expected Missing=2, got %d", s.Missing)
	}
	if s.Extra != 1 {
		t.Errorf("expected Extra=1, got %d", s.Extra)
	}
	if s.Mismatched != 1 {
		t.Errorf("expected Mismatched=1, got %d", s.Mismatched)
	}
	if !s.HasDifferences() {
		t.Error("expected HasDifferences to be true")
	}
}

func TestBreakdown_OmitsZeros(t *testing.T) {
	results := []diff.Result{
		{Key: "X", Kind: diff.KindValueMismatch},
	}
	s := summary.Compute(results)
	b := summary.Breakdown(s)
	if _, ok := b["missing"]; ok {
		t.Error("expected missing key to be absent from breakdown")
	}
	if _, ok := b["extra"]; ok {
		t.Error("expected extra key to be absent from breakdown")
	}
	if b["mismatched"] != 1 {
		t.Errorf("expected mismatched=1, got %d", b["mismatched"])
	}
}

func TestBreakdown_Empty(t *testing.T) {
	s := summary.Compute(nil)
	b := summary.Breakdown(s)
	if len(b) != 0 {
		t.Errorf("expected empty breakdown, got %v", b)
	}
}
