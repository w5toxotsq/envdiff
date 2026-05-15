package suggest_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/suggest"
)

func makeResult(key string, kind diff.Kind) diff.Result {
	return diff.Result{Key: key, Kind: kind}
}

func TestRun_NoResults_ReturnsEmpty(t *testing.T) {
	out := suggest.Run(nil, []string{"DB_PASSWORD", "DB_HOST"}, suggest.DefaultOptions())
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d suggestions", len(out))
	}
}

func TestRun_MismatchedKindSkipped(t *testing.T) {
	results := []diff.Result{makeResult("DB_PASSWORD", diff.KindMismatch)}
	out := suggest.Run(results, []string{"DB_PASWORD"}, suggest.DefaultOptions())
	if len(out) != 0 {
		t.Fatalf("expected no suggestions for mismatch kind, got %d", len(out))
	}
}

func TestRun_MissingInRight_FindsTypo(t *testing.T) {
	results := []diff.Result{makeResult("DB_PASSWORD", diff.KindMissingInRight)}
	known := []string{"DB_PASWORD", "DB_HOST", "UNRELATED_KEY_XYZ"}
	out := suggest.Run(results, known, suggest.DefaultOptions())
	if len(out) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(out))
	}
	if len(out[0].Candidates) == 0 {
		t.Fatal("expected at least one candidate")
	}
	if out[0].Candidates[0] != "DB_PASWORD" {
		t.Errorf("expected DB_PASWORD as top candidate, got %q", out[0].Candidates[0])
	}
}

func TestRun_MissingInLeft_FindsTypo(t *testing.T) {
	results := []diff.Result{makeResult("APP_SECRETT", diff.KindMissingInLeft)}
	known := []string{"APP_SECRET", "APP_KEY"}
	out := suggest.Run(results, known, suggest.DefaultOptions())
	if len(out) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(out))
	}
	if out[0].Candidates[0] != "APP_SECRET" {
		t.Errorf("expected APP_SECRET, got %q", out[0].Candidates[0])
	}
}

func TestRun_NoCloseMatch_OmitsSuggestion(t *testing.T) {
	results := []diff.Result{makeResult("COMPLETELY_DIFFERENT", diff.KindMissingInRight)}
	known := []string{"XYZ", "ABC"}
	out := suggest.Run(results, known, suggest.DefaultOptions())
	if len(out) != 0 {
		t.Fatalf("expected no suggestions, got %d", len(out))
	}
}

func TestRun_MaxCandidatesRespected(t *testing.T) {
	results := []diff.Result{makeResult("KEY", diff.KindMissingInRight)}
	// All within distance 1 of "KEY"
	known := []string{"KEX", "KAY", "KFY", "KGY"}
	opts := suggest.DefaultOptions()
	opts.MaxCandidates = 2
	out := suggest.Run(results, known, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1 suggestion entry, got %d", len(out))
	}
	if len(out[0].Candidates) > 2 {
		t.Errorf("expected at most 2 candidates, got %d", len(out[0].Candidates))
	}
}

func TestRun_ExactMatchExcluded(t *testing.T) {
	results := []diff.Result{makeResult("DB_HOST", diff.KindMissingInRight)}
	known := []string{"DB_HOST", "DB_HOSE"}
	out := suggest.Run(results, known, suggest.DefaultOptions())
	if len(out) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(out))
	}
	for _, c := range out[0].Candidates {
		if c == "DB_HOST" {
			t.Error("exact match should not appear as candidate")
		}
	}
}
