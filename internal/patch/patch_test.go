package patch_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your/envdiff/internal/diff"
	"github.com/your/envdiff/internal/patch"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	return f.Name()
}

func TestApply_DryRun_DoesNotWriteFile(t *testing.T) {
	path := writeTempEnv(t, "A=1\n")
	dst := map[string]string{"A": "1"}
	results := []diff.Result{
		{Key: "B", LeftValue: "2", RightValue: "", Kind: diff.MissingInRight},
	}

	original, _ := os.ReadFile(path)
	opts := patch.DefaultOptions()
	opts.DryRun = true

	out, err := patch.Apply(path, dst, results, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.DryRun {
		t.Error("expected DryRun=true in result")
	}
	if len(out.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(out.Changes))
	}

	after, _ := os.ReadFile(path)
	if string(after) != string(original) {
		t.Error("file was modified during dry run")
	}
}

func TestApply_WritesMissingKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.env")
	_ = os.WriteFile(path, []byte("A=1\n"), 0o644)

	dst := map[string]string{"A": "1"}
	results := []diff.Result{
		{Key: "B", LeftValue: "hello", RightValue: "", Kind: diff.MissingInRight},
	}

	out, err := patch.Apply(path, dst, results, patch.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(out.Changes))
	}
	if out.Changes[0].Key != "B" || out.Changes[0].NewValue != "hello" {
		t.Errorf("unexpected change: %+v", out.Changes[0])
	}
}

func TestApply_OnlyMissing_SkipsMismatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.env")
	_ = os.WriteFile(path, []byte("A=old\n"), 0o644)

	dst := map[string]string{"A": "old"}
	results := []diff.Result{
		{Key: "A", LeftValue: "new", RightValue: "old", Kind: diff.Mismatch},
		{Key: "B", LeftValue: "added", RightValue: "", Kind: diff.MissingInRight},
	}

	opts := patch.DefaultOptions()
	opts.OnlyMissing = true

	out, err := patch.Apply(path, dst, results, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Changes) != 1 {
		t.Errorf("expected 1 change (only missing), got %d", len(out.Changes))
	}
	if out.Changes[0].Key != "B" {
		t.Errorf("expected change for B, got %s", out.Changes[0].Key)
	}
}
