package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your/envdiff/internal/baseline"
	"github.com/your/envdiff/internal/diff"
)

func makeResult(key string, kind diff.Kind) diff.Result {
	return diff.Result{Key: key, Kind: kind}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	err := baseline.Save(path, "test", []diff.Result{makeResult("FOO", diff.MissingInRight)})
	if err != nil {
		t.Fatalf("Save: unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	input := []diff.Result{
		makeResult("DB_HOST", diff.MissingInRight),
		makeResult("API_KEY", diff.ValueMismatch),
	}
	if err := baseline.Save(path, "ci", input); err != nil {
		t.Fatalf("Save: %v", err)
	}

	b, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if b.Label != "ci" {
		t.Errorf("label: got %q, want %q", b.Label, "ci")
	}
	if len(b.Results) != 2 {
		t.Fatalf("results len: got %d, want 2", len(b.Results))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := baseline.Load("/nonexistent/baseline.json")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCompare_NewResults(t *testing.T) {
	base := &baseline.Baseline{
		Results: []diff.Result{
			makeResult("FOO", diff.MissingInRight),
		},
	}
	current := []diff.Result{
		makeResult("FOO", diff.MissingInRight), // already in baseline
		makeResult("BAR", diff.ValueMismatch),   // new
	}
	novel := baseline.Compare(base, current)
	if len(novel) != 1 {
		t.Fatalf("novel len: got %d, want 1", len(novel))
	}
	if novel[0].Key != "BAR" {
		t.Errorf("novel key: got %q, want %q", novel[0].Key, "BAR")
	}
}

func TestCompare_EmptyBaseline(t *testing.T) {
	base := &baseline.Baseline{Results: []diff.Result{}}
	current := []diff.Result{makeResult("X", diff.MissingInLeft)}
	novel := baseline.Compare(base, current)
	if len(novel) != 1 {
		t.Fatalf("expected 1 novel result, got %d", len(novel))
	}
}
