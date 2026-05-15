package redact_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/redact"
)

func TestValue_RedactionEnabled_MasksAll(t *testing.T) {
	opts := redact.Options{Enabled: true, Mask: "***"}
	got := redact.Value("SOME_VAR", "plaintext", opts)
	if got != "***" {
		t.Errorf("expected *** got %q", got)
	}
}

func TestValue_RedactionDisabled_ReturnsOriginal(t *testing.T) {
	opts := redact.Options{Enabled: false, Mask: "***"}
	got := redact.Value("SOME_VAR", "plaintext", opts)
	if got != "plaintext" {
		t.Errorf("expected plaintext got %q", got)
	}
}

func TestValue_SensitiveSuffix_AlwaysMasked(t *testing.T) {
	opts := redact.Options{
		Enabled:           false,
		Mask:              "[hidden]",
		SensitiveSuffixes: []string{"_SECRET", "_TOKEN"},
	}

	cases := []struct {
		key string
	}{
		{"DB_SECRET"},
		{"AUTH_TOKEN"},
		{"db_secret"},  // case-insensitive
		{"AUTH_token"}, // mixed case
	}

	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			got := redact.Value(tc.key, "s3cr3t", opts)
			if got != "[hidden]" {
				t.Errorf("key %q: expected [hidden] got %q", tc.key, got)
			}
		})
	}
}

func TestValue_NonSensitiveKey_NotMaskedWhenDisabled(t *testing.T) {
	opts := redact.Options{
		Enabled:           false,
		Mask:              "***",
		SensitiveSuffixes: []string{"_SECRET"},
	}
	got := redact.Value("APP_ENV", "production", opts)
	if got != "production" {
		t.Errorf("expected production got %q", got)
	}
}

func TestValue_EmptyMask_UsesDefault(t *testing.T) {
	opts := redact.Options{Enabled: true, Mask: ""}
	got := redact.Value("KEY", "value", opts)
	if got != "***" {
		t.Errorf("expected *** got %q", got)
	}
}

func TestDefaultOptions_SensitiveSuffixesNonEmpty(t *testing.T) {
	opts := redact.DefaultOptions()
	if len(opts.SensitiveSuffixes) == 0 {
		t.Error("expected at least one sensitive suffix in DefaultOptions")
	}
	if !opts.Enabled {
		t.Error("expected Enabled to be true in DefaultOptions")
	}
}
