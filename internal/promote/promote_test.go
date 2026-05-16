package promote_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/promote"
)

func makeResults(entries ...diff.Result) []diff.Result { return entries }

func TestPromote_AddsMissingKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"A": "1"}
	results := makeResults(diff.Result{Key: "B", Kind: diff.KindMissingInRight})

	r, err := promote.Promote(src, dst, results, promote.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Changes) != 1 || r.Changes[0].Key != "B" {
		t.Fatalf("expected one change for B, got %v", r.Changes)
	}
	if dst["B"] != "2" {
		t.Errorf("expected dst[B]=2, got %q", dst["B"])
	}
}

func TestPromote_SkipsMismatchWithoutOverwrite(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	results := makeResults(diff.Result{Key: "A", Kind: diff.KindMismatch})

	r, err := promote.Promote(src, dst, results, promote.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Changes) != 0 {
		t.Errorf("expected no changes, got %v", r.Changes)
	}
	if dst["A"] != "old" {
		t.Errorf("expected dst[A] unchanged, got %q", dst["A"])
	}
}

func TestPromote_OverwritesMismatchWhenEnabled(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	results := makeResults(diff.Result{Key: "A", Kind: diff.KindMismatch})

	opts := promote.DefaultOptions()
	opts.Overwrite = true
	r, err := promote.Promote(src, dst, results, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Changes) != 1 || !r.Changes[0].Overwritten {
		t.Errorf("expected one overwrite change, got %v", r.Changes)
	}
	if dst["A"] != "new" {
		t.Errorf("expected dst[A]=new, got %q", dst["A"])
	}
}

func TestPromote_DryRun_DoesNotModifyDst(t *testing.T) {
	src := map[string]string{"B": "99"}
	dst := map[string]string{}
	results := makeResults(diff.Result{Key: "B", Kind: diff.KindMissingInRight})

	opts := promote.DefaultOptions()
	opts.DryRun = true
	r, err := promote.Promote(src, dst, results, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Changes) != 1 {
		t.Errorf("expected one reported change, got %d", len(r.Changes))
	}
	if _, ok := dst["B"]; ok {
		t.Error("dry-run must not modify dst")
	}
}

func TestPromote_NilSrc_ReturnsError(t *testing.T) {
	_, err := promote.Promote(nil, map[string]string{}, nil, promote.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil src")
	}
}

func TestPromote_NilDst_ReturnsError(t *testing.T) {
	_, err := promote.Promote(map[string]string{}, nil, nil, promote.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil dst")
	}
}
