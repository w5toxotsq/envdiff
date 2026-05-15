package baseline_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/your/envdiff/internal/baseline"
	"github.com/your/envdiff/internal/diff"
)

func TestSave_SetsCreatedAt(t *testing.T) {
	before := time.Now().UTC()
	dir := t.TempDir()
	path := filepath.Join(dir, "b.json")

	if err := baseline.Save(path, "", nil); err != nil {
		t.Fatalf("Save: %v", err)
	}

	b, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if b.CreatedAt.Before(before) {
		t.Errorf("CreatedAt %v is before test start %v", b.CreatedAt, before)
	}
}

func TestSave_NilResults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "b.json")

	if err := baseline.Save(path, "", nil); err != nil {
		t.Fatalf("Save with nil results: %v", err)
	}
	b, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if b.Results == nil {
		t.Error("expected non-nil Results slice after round-trip")
	}
}

func TestCompare_AllAlreadyInBaseline(t *testing.T) {
	results := []diff.Result{
		makeResult("A", diff.MissingInLeft),
		makeResult("B", diff.ValueMismatch),
	}
	base := &baseline.Baseline{Results: results}
	novel := baseline.Compare(base, results)
	if len(novel) != 0 {
		t.Errorf("expected 0 novel results, got %d", len(novel))
	}
}
