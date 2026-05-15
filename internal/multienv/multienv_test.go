package multienv_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/multienv"
)

func writeTempEnv(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestLoad_Basic(t *testing.T) {
	p1 := writeTempEnv(t, ".env.dev", "FOO=bar\nBAZ=qux\n")
	p2 := writeTempEnv(t, ".env.prod", "FOO=barprod\nSECRET=yes\n")

	set, err := multienv.Load(
		[]string{"dev", "prod"},
		[]string{p1, p2},
		multienv.LoadOptions{SkipExtensionCheck: true},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(set) != 2 {
		t.Errorf("expected 2 envs, got %d", len(set))
	}
	if set["dev"]["FOO"] != "bar" {
		t.Errorf("dev FOO: got %q, want %q", set["dev"]["FOO"], "bar")
	}
	if set["prod"]["FOO"] != "barprod" {
		t.Errorf("prod FOO: got %q, want %q", set["prod"]["FOO"], "barprod")
	}
}

func TestLoad_MismatchedLengths(t *testing.T) {
	_, err := multienv.Load(
		[]string{"dev"},
		[]string{},
		multienv.LoadOptions{},
	)
	if err == nil {
		t.Fatal("expected error for mismatched lengths, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := multienv.Load(
		[]string{"dev"},
		[]string{"/nonexistent/.env"},
		multienv.LoadOptions{},
	)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
