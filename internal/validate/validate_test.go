package validate_test

import (
	"testing"

	"github.com/envdiff/envdiff/internal/diff"
	"github.com/envdiff/envdiff/internal/validate"
)

func makeResult(key string, kind diff.Kind, left, right string) diff.Result {
	return diff.Result{Key: key, Kind: kind, Left: left, Right: right}
}

func TestRun_NoIssues(t *testing.T) {
	results := []diff.Result{
		makeResult("DATABASE_URL", diff.KindMismatch, "postgres://a", "postgres://b"),
	}
	opts := validate.DefaultOptions()
	issues := validate.Run(results, opts)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestRun_InvalidKeyFormat(t *testing.T) {
	results := []diff.Result{
		makeResult("bad-key", diff.KindMissingRight, "", ""),
		makeResult("GOOD_KEY", diff.KindMissingRight, "", ""),
	}
	opts := validate.DefaultOptions()
	issues := validate.Run(results, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d: %v", len(issues), issues)
	}
	if issues[0].Key != "bad-key" {
		t.Errorf("expected issue for 'bad-key', got %q", issues[0].Key)
	}
}

func TestRun_EmptyValue(t *testing.T) {
	results := []diff.Result{
		makeResult("API_KEY", diff.KindMismatch, "", "secret"),
	}
	opts := validate.DefaultOptions()
	issues := validate.Run(results, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "API_KEY" {
		t.Errorf("unexpected key %q", issues[0].Key)
	}
}

func TestRun_URLSchemeMismatch(t *testing.T) {
	results := []diff.Result{
		makeResult("REDIS_URL", diff.KindMismatch, "redis://localhost", "rediss://prod"),
	}
	opts := validate.DefaultOptions()
	opts.CheckURLValues = true
	issues := validate.Run(results, opts)
	// may also include empty-value issues; find the scheme one
	found := false
	for _, iss := range issues {
		if iss.Key == "REDIS_URL" && contains(iss.Message, "scheme") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected URL scheme mismatch issue, got %v", issues)
	}
}

func TestRun_DisabledChecks(t *testing.T) {
	results := []diff.Result{
		makeResult("bad-key", diff.KindMismatch, "", ""),
	}
	opts := validate.Options{} // all disabled
	issues := validate.Run(results, opts)
	if len(issues) != 0 {
		t.Fatalf("expected no issues with all checks disabled, got %v", issues)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
