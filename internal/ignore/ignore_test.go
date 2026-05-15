package ignore_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/diff"
	"github.com/your-org/envdiff/internal/ignore"
)

func makeResult(key string, kind diff.Kind) diff.Result {
	return diff.Result{Key: key, Kind: kind}
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	results := []diff.Result{
		makeResult("FOO", diff.KindMismatch),
		makeResult("BAR", diff.KindMissingLeft),
	}
	out := ignore.Apply(results, ignore.Options{})
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestApply_ExactKeyIgnored(t *testing.T) {
	results := []diff.Result{
		makeResult("SECRET", diff.KindMismatch),
		makeResult("HOST", diff.KindMismatch),
	}
	out := ignore.Apply(results, ignore.Options{Keys: []string{"SECRET"}})
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Key != "HOST" {
		t.Errorf("expected HOST, got %s", out[0].Key)
	}
}

func TestApply_PatternIgnored(t *testing.T) {
	results := []diff.Result{
		makeResult("CI_BUILD_ID", diff.KindMissingRight),
		makeResult("CI_BRANCH", diff.KindMismatch),
		makeResult("APP_ENV", diff.KindMismatch),
	}
	out := ignore.Apply(results, ignore.Options{Patterns: []string{"CI_"}})
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV, got %s", out[0].Key)
	}
}

func TestApply_OriginalUnmodified(t *testing.T) {
	original := []diff.Result{
		makeResult("IGNORED", diff.KindMismatch),
		makeResult("KEPT", diff.KindMismatch),
	}
	_ = ignore.Apply(original, ignore.Options{Keys: []string{"IGNORED"}})
	if len(original) != 2 {
		t.Error("original slice was modified")
	}
}

func TestApply_CombinedKeysAndPatterns(t *testing.T) {
	results := []diff.Result{
		makeResult("SECRET", diff.KindMismatch),
		makeResult("INTERNAL_FLAG", diff.KindMissingLeft),
		makeResult("APP_PORT", diff.KindMismatch),
	}
	out := ignore.Apply(results, ignore.Options{
		Keys:     []string{"SECRET"},
		Patterns: []string{"INTERNAL_"},
	})
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Key != "APP_PORT" {
		t.Errorf("expected APP_PORT, got %s", out[0].Key)
	}
}
