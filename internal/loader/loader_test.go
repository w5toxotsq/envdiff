package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/loader"
)

func writeTempEnv(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return path
}

func TestLoad_ValidEnvFile(t *testing.T) {
	path := writeTempEnv(t, ".env", "KEY=value\nFOO=bar\n")
	env, err := loader.Load(path, loader.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", env["KEY"])
	}
	if env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", env["FOO"])
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := loader.Load("/nonexistent/.env", loader.Options{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var notFound *loader.ErrFileNotFound
	if _, ok := err.(*loader.ErrFileNotFound); !ok {
		_ = notFound
		t.Errorf("expected ErrFileNotFound, got %T: %v", err, err)
	}
}

func TestLoad_InvalidExtension(t *testing.T) {
	path := writeTempEnv(t, "config.txt", "KEY=value\n")
	_, err := loader.Load(path, loader.Options{})
	if err == nil {
		t.Fatal("expected error for invalid extension, got nil")
	}
	if _, ok := err.(*loader.ErrInvalidExtension); !ok {
		t.Errorf("expected ErrInvalidExtension, got %T: %v", err, err)
	}
}

func TestLoad_SkipExtensionCheck(t *testing.T) {
	path := writeTempEnv(t, "config.txt", "KEY=value\n")
	env, err := loader.Load(path, loader.Options{SkipExtensionCheck: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", env["KEY"])
	}
}

func TestLoad_DotEnvPrefixedFile(t *testing.T) {
	path := writeTempEnv(t, ".env.production", "DB=prod\n")
	env, err := loader.Load(path, loader.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["DB"] != "prod" {
		t.Errorf("expected DB=prod, got %q", env["DB"])
	}
}
