package audit_test

import (
	"strings"
	"testing"

	"github.com/your/envdiff/internal/audit"
	"github.com/your/envdiff/internal/diff"
)

func TestBuild_EmptyResults(t *testing.T) {
	entry := audit.Build("a", "b", []diff.Result{}, audit.DefaultOptions())

	if entry.Total != 0 {
		t.Errorf("Total: got %d, want 0", entry.Total)
	}
	if entry.Missing != 0 {
		t.Errorf("Missing: got %d, want 0", entry.Missing)
	}
	if entry.Mismatch != 0 {
		t.Errorf("Mismatch: got %d, want 0", entry.Mismatch)
	}
}

func TestBuild_NilResults(t *testing.T) {
	entry := audit.Build("a", "b", nil, audit.DefaultOptions())

	if entry.Total != 0 {
		t.Errorf("Total: got %d, want 0", entry.Total)
	}
}

func TestDefaultOptions_IncludeResultsFalse(t *testing.T) {
	opts := audit.DefaultOptions()
	if opts.IncludeResults {
		t.Error("DefaultOptions.IncludeResults should be false")
	}
}

func TestSummary_ContainsCounts(t *testing.T) {
	results := makeResults()
	entry := audit.Build("a", "b", results, audit.DefaultOptions())
	s := entry.Summary()

	for _, want := range []string{"total:3", "missing:2", "mismatch:1"} {
		if !strings.Contains(s, want) {
			t.Errorf("Summary %q missing %q", s, want)
		}
	}
}
