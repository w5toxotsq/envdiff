package patch_test

import (
	"testing"

	"github.com/your/envdiff/internal/diff"
	"github.com/your/envdiff/internal/patch"
)

func TestApply_EmptyResults_NoChanges(t *testing.T) {
	path := writeTempEnv(t, "A=1\n")
	dst := map[string]string{"A": "1"}

	out, err := patch.Apply(path, dst, nil, patch.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(out.Changes))
	}
}

func TestApply_MissingInLeft_IsIgnored(t *testing.T) {
	path := writeTempEnv(t, "A=1\n")
	dst := map[string]string{"A": "1"}
	results := []diff.Result{
		{Key: "X", LeftValue: "", RightValue: "extra", Kind: diff.MissingInLeft},
	}

	out, err := patch.Apply(path, dst, results, patch.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Changes) != 0 {
		t.Errorf("expected 0 changes for MissingInLeft, got %d", len(out.Changes))
	}
}

func TestDefaultOptions_Values(t *testing.T) {
	opts := patch.DefaultOptions()
	if opts.DryRun {
		t.Error("expected DryRun=false by default")
	}
	if opts.OnlyMissing {
		t.Error("expected OnlyMissing=false by default")
	}
}

func TestApply_Change_KindPreserved(t *testing.T) {
	path := writeTempEnv(t, "A=old\n")
	dst := map[string]string{"A": "old"}
	results := []diff.Result{
		{Key: "A", LeftValue: "new", RightValue: "old", Kind: diff.Mismatch},
	}

	out, err := patch.Apply(path, dst, results, patch.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(out.Changes))
	}
	if out.Changes[0].Kind != diff.Mismatch {
		t.Errorf("expected Kind=Mismatch, got %v", out.Changes[0].Kind)
	}
	if out.Changes[0].OldValue != "old" {
		t.Errorf("expected OldValue=old, got %s", out.Changes[0].OldValue)
	}
}
