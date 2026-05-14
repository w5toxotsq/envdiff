package config_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/config"
)

func TestLoadFile_OnlyMissingFlag(t *testing.T) {
	path := writeYAML(t, "only_missing: true\n")
	cfg := config.DefaultConfig()
	if err := config.LoadFile(path, &cfg); err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	if !cfg.OnlyMissing {
		t.Error("expected OnlyMissing true after load")
	}
}

func TestLoadFile_OnlyMismatchedFlag(t *testing.T) {
	path := writeYAML(t, "only_mismatched: true\n")
	cfg := config.DefaultConfig()
	if err := config.LoadFile(path, &cfg); err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	if !cfg.OnlyMismatched {
		t.Error("expected OnlyMismatched true after load")
	}
}

func TestLoadFile_EmptyFile(t *testing.T) {
	path := writeYAML(t, "")
	cfg := config.DefaultConfig()
	origFormat := cfg.Format
	if err := config.LoadFile(path, &cfg); err != nil {
		t.Fatalf("LoadFile on empty file: %v", err)
	}
	if cfg.Format != origFormat {
		t.Errorf("format changed on empty file load: got %q", cfg.Format)
	}
}

func TestDefaultConfig_ZeroFiles(t *testing.T) {
	cfg := config.DefaultConfig()
	if len(cfg.Files) != 0 {
		t.Errorf("expected no files in default config, got %v", cfg.Files)
	}
}
