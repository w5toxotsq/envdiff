package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envdiff/internal/config"
)

func writeYAML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "envdiff.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeYAML: %v", err)
	}
	return p
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.Format != "text" {
		t.Errorf("expected format 'text', got %q", cfg.Format)
	}
	if cfg.ShowValues {
		t.Error("expected ShowValues false by default")
	}
}

func TestLoadFile_Basic(t *testing.T) {
	path := writeYAML(t, "format: json\nshow_values: true\nprefix: APP_\n")
	cfg := config.DefaultConfig()
	if err := config.LoadFile(path, &cfg); err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	if cfg.Format != "json" {
		t.Errorf("expected format 'json', got %q", cfg.Format)
	}
	if !cfg.ShowValues {
		t.Error("expected ShowValues true after load")
	}
	if cfg.Prefix != "APP_" {
		t.Errorf("expected prefix 'APP_', got %q", cfg.Prefix)
	}
}

func TestLoadFile_DoesNotOverwriteExisting(t *testing.T) {
	path := writeYAML(t, "format: json\nprefix: DB_\n")
	cfg := config.DefaultConfig()
	cfg.Format = "text" // already set — should not be overwritten
	cfg.Prefix = "MY_"
	if err := config.LoadFile(path, &cfg); err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	if cfg.Format != "text" {
		t.Errorf("format should not be overwritten, got %q", cfg.Format)
	}
	if cfg.Prefix != "MY_" {
		t.Errorf("prefix should not be overwritten, got %q", cfg.Prefix)
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	cfg := config.DefaultConfig()
	err := config.LoadFile("/nonexistent/path/envdiff.yaml", &cfg)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadFile_Files(t *testing.T) {
	path := writeYAML(t, "files:\n  - .env.dev\n  - .env.prod\n")
	cfg := config.DefaultConfig()
	if err := config.LoadFile(path, &cfg); err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	if len(cfg.Files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(cfg.Files))
	}
	if cfg.Files[0] != ".env.dev" || cfg.Files[1] != ".env.prod" {
		t.Errorf("unexpected files: %v", cfg.Files)
	}
}
