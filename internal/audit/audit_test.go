package audit_test

import (
	"strings"
	"testing"
	"time"

	"github.com/your/envdiff/internal/audit"
	"github.com/your/envdiff/internal/diff"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Kind: diff.MissingInRight},
		{Key: "API_KEY", Kind: diff.ValueMismatch, LeftVal: "abc", RightVal: "xyz"},
		{Key: "PORT", Kind: diff.MissingInLeft},
	}
}

func TestBuild_Counts(t *testing.T) {
	results := makeResults()
	entry := audit.Build(".env.a", ".env.b", results, audit.DefaultOptions())

	if entry.Total != 3 {
		t.Errorf("Total: got %d, want 3", entry.Total)
	}
	if entry.Missing != 2 {
		t.Errorf("Missing: got %d, want 2", entry.Missing)
	}
	if entry.Mismatch != 1 {
		t.Errorf("Mismatch: got %d, want 1", entry.Mismatch)
	}
}

func TestBuild_FilesRecorded(t *testing.T) {
	entry := audit.Build(".env.staging", ".env.prod", nil, audit.DefaultOptions())

	if entry.LeftFile != ".env.staging" {
		t.Errorf("LeftFile: got %q", entry.LeftFile)
	}
	if entry.RightFile != ".env.prod" {
		t.Errorf("RightFile: got %q", entry.RightFile)
	}
}

func TestBuild_TimestampSet(t *testing.T) {
	before := time.Now().UTC()
	entry := audit.Build("a", "b", nil, audit.DefaultOptions())
	after := time.Now().UTC()

	if entry.Timestamp.Before(before) || entry.Timestamp.After(after) {
		t.Errorf("Timestamp %v not in expected range", entry.Timestamp)
	}
}

func TestBuild_IncludeResults_False(t *testing.T) {
	entry := audit.Build("a", "b", makeResults(), audit.DefaultOptions())
	if entry.Results != nil {
		t.Error("expected Results to be nil when IncludeResults is false")
	}
}

func TestBuild_IncludeResults_True(t *testing.T) {
	opts := audit.Options{IncludeResults: true}
	entry := audit.Build("a", "b", makeResults(), opts)
	if len(entry.Results) != 3 {
		t.Errorf("Results len: got %d, want 3", len(entry.Results))
	}
}

func TestSummary_ContainsFiles(t *testing.T) {
	entry := audit.Build(".env.dev", ".env.prod", makeResults(), audit.DefaultOptions())
	s := entry.Summary()

	if !strings.Contains(s, ".env.dev") {
		t.Errorf("Summary missing left file: %s", s)
	}
	if !strings.Contains(s, ".env.prod") {
		t.Errorf("Summary missing right file: %s", s)
	}
}
