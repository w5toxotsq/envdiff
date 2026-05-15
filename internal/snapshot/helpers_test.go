package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your/envdiff/internal/diff"
	"github.com/your/envdiff/internal/snapshot"
)

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")

	if err := os.WriteFile(path, []byte(`{not valid json`), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := snapshot.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestSave_EmptyResults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")

	if err := snapshot.Save(path, ".env", ".env.staging", []diff.Result{}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(snap.Results) != 0 {
		t.Errorf("Results len = %d, want 0", len(snap.Results))
	}
}

func TestSave_NilResults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nil.json")

	if err := snapshot.Save(path, "x", "y", nil); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if snap.LeftFile != "x" || snap.RightFile != "y" {
		t.Errorf("files = %q / %q, want x / y", snap.LeftFile, snap.RightFile)
	}
}
