package rename_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/rename"
)

func TestFromDiff_DetectsRename(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", Kind: diff.MissingInRight, LeftValue: "localhost", RightValue: ""},
		{Key: "DATABASE_HOST", Kind: diff.MissingInLeft, LeftValue: "", RightValue: "localhost"},
	}

	mappings := rename.FromDiff(results)

	if len(mappings) != 1 {
		t.Fatalf("expected 1 mapping, got %d", len(mappings))
	}
	if mappings[0].OldKey != "DB_HOST" {
		t.Errorf("expected OldKey=DB_HOST, got %q", mappings[0].OldKey)
	}
	if mappings[0].NewKey != "DATABASE_HOST" {
		t.Errorf("expected NewKey=DATABASE_HOST, got %q", mappings[0].NewKey)
	}
}

func TestFromDiff_NoMatchReturnsEmpty(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", Kind: diff.MissingInRight, LeftValue: "localhost", RightValue: ""},
		{Key: "OTHER_KEY", Kind: diff.MissingInLeft, LeftValue: "", RightValue: "differentvalue"},
	}

	mappings := rename.FromDiff(results)
	if len(mappings) != 0 {
		t.Errorf("expected 0 mappings, got %d", len(mappings))
	}
}

func TestFromDiff_IgnoresMismatchKind(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", Kind: diff.Mismatch, LeftValue: "localhost", RightValue: "remotehost"},
	}

	mappings := rename.FromDiff(results)
	if len(mappings) != 0 {
		t.Errorf("expected 0 mappings for Mismatch kind, got %d", len(mappings))
	}
}

func TestFromDiff_EmptyResults(t *testing.T) {
	mappings := rename.FromDiff(nil)
	if mappings != nil && len(mappings) != 0 {
		t.Errorf("expected empty mappings for nil input, got %v", mappings)
	}
}

func TestDefaultOptions_Values(t *testing.T) {
	opts := rename.DefaultOptions()
	if opts.DryRun {
		t.Error("DryRun should default to false")
	}
	if opts.FailOnMissing {
		t.Error("FailOnMissing should default to false")
	}
}
