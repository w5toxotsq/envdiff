package mask_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/mask"
)

func TestValue_MaskingDisabled_ReturnsOriginal(t *testing.T) {
	opts := mask.DefaultOptions()
	opts.Enabled = false

	got := mask.Value("DB_PASSWORD", "supersecret", opts)
	if got != "supersecret" {
		t.Fatalf("expected original value, got %q", got)
	}
}

func TestValue_SecretKey_IsMasked(t *testing.T) {
	opts := mask.DefaultOptions()

	got := mask.Value("API_SECRET", "abc123", opts)
	if got != "***" {
		t.Fatalf("expected masked value, got %q", got)
	}
}

func TestValue_NonSensitiveKey_NotMasked(t *testing.T) {
	opts := mask.DefaultOptions()

	got := mask.Value("APP_ENV", "production", opts)
	if got != "production" {
		t.Fatalf("expected original value, got %q", got)
	}
}

func TestValue_AllowlistedKey_NotMasked(t *testing.T) {
	opts := mask.DefaultOptions()
	opts.Allowlist = []string{"PUBLIC_KEY"}

	got := mask.Value("PUBLIC_KEY", "open", opts)
	if got != "open" {
		t.Fatalf("expected allowlisted value to be unmasked, got %q", got)
	}
}

func TestValue_AllowlistCaseInsensitive(t *testing.T) {
	opts := mask.DefaultOptions()
	opts.Allowlist = []string{"public_key"}

	got := mask.Value("PUBLIC_KEY", "open", opts)
	if got != "open" {
		t.Fatalf("expected case-insensitive allowlist match, got %q", got)
	}
}

func TestValue_CustomMask(t *testing.T) {
	opts := mask.DefaultOptions()
	opts.Mask = "[REDACTED]"

	got := mask.Value("DB_PASSWORD", "hunter2", opts)
	if got != "[REDACTED]" {
		t.Fatalf("expected custom mask string, got %q", got)
	}
}

func TestValue_EmptyMask_UsesDefault(t *testing.T) {
	opts := mask.DefaultOptions()
	opts.Mask = ""

	got := mask.Value("AUTH_TOKEN", "tok", opts)
	if got != "***" {
		t.Fatalf("expected default mask ***, got %q", got)
	}
}

func TestDefaultOptions_PatternsNonEmpty(t *testing.T) {
	opts := mask.DefaultOptions()
	if len(opts.Patterns) == 0 {
		t.Fatal("expected default patterns to be non-empty")
	}
	if !opts.Enabled {
		t.Fatal("expected masking to be enabled by default")
	}
}
