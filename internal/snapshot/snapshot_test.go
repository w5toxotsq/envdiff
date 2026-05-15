package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your/envdiff/internal/diff"
	"github.com/your/envdiff/internal/snapshot"
)

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	results := []diff.Result{
		{Key: "FOO", Kind: diff.MissingInRight},
	}

	if err := snapshot.Save(path, ".env", ".env.prod", results); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	results := []diff.Result{
		{Key: "DB_HOST", Kind: diff.ValueMismatch, LeftValue: "localhost", RightValue: "prod-db"},
		{Key: "SECRET", Kind: diff.MissingInLeft},
	}

	if err := snapshot.Save(path, "a.env", "b.env", results); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if snap.LeftFile != "a.env" {
		t.Errorf("LeftFile = %q, want %q", snap.LeftFile, "a.env")
	}
	if snap.RightFile != "b.env" {
		t.Errorf("RightFile = %q, want %q", snap.RightFile, "b.env")
	}
	if len(snap.Results) != 2 {
		t.Fatalf("Results len = %d, want 2", len(snap.Results))
	}
	if snap.Results[0].Key != "DB_HOST" {
		t.Errorf("Results[0].Key = %q, want %q", snap.Results[0].Key, "DB_HOST")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSave_SetsCreatedAt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	before := time.Now().UTC()

	if err := snapshot.Save(path, ".env", ".env.prod", nil); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if snap.CreatedAt.Before(before) {
		t.Errorf("CreatedAt %v is before test start %v", snap.CreatedAt, before)
	}
}
