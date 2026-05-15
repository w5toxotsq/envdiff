package lint_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/lint"
)

func TestRun_NoIssues(t *testing.T) {
	env := map[string]string{"APP_NAME": "myapp", "PORT": "8080"}
	keys := []string{"APP_NAME", "PORT"}
	issues := lint.Run(env, keys, lint.DefaultOptions())
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestRun_InvalidNaming(t *testing.T) {
	env := map[string]string{"app_name": "myapp", "1INVALID": "x"}
	keys := []string{"app_name", "1INVALID"}
	issues := lint.Run(env, keys, lint.DefaultOptions())
	if len(issues) != 2 {
		t.Fatalf("expected 2 naming issues, got %d", len(issues))
	}
	for _, iss := range issues {
		if iss.Severity != lint.SeverityWarning {
			t.Errorf("expected warning severity, got %s", iss.Severity)
		}
	}
}

func TestRun_EmptyValue(t *testing.T) {
	env := map[string]string{"SECRET": "", "TOKEN": "   "}
	keys := []string{"SECRET", "TOKEN"}
	opts := lint.Options{CheckEmptyValues: true}
	issues := lint.Run(env, keys, opts)
	if len(issues) != 2 {
		t.Fatalf("expected 2 empty-value issues, got %d", len(issues))
	}
	for _, iss := range issues {
		if iss.Severity != lint.SeverityWarning {
			t.Errorf("expected warning severity, got %s", iss.Severity)
		}
	}
}

func TestRun_DuplicateKeys(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	keys := []string{"PORT", "PORT", "PORT"}
	opts := lint.Options{CheckDuplicates: true}
	issues := lint.Run(env, keys, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 duplicate issue, got %d", len(issues))
	}
	if issues[0].Severity != lint.SeverityError {
		t.Errorf("expected error severity for duplicate, got %s", issues[0].Severity)
	}
}

func TestRun_DisabledChecks(t *testing.T) {
	env := map[string]string{"bad_key": ""}
	keys := []string{"bad_key", "bad_key"}
	opts := lint.Options{} // all checks disabled
	issues := lint.Run(env, keys, opts)
	if len(issues) != 0 {
		t.Fatalf("expected no issues with all checks disabled, got %d", len(issues))
	}
}

func TestIssue_String(t *testing.T) {
	issue := lint.Issue{Key: "FOO", Message: "value is empty", Severity: lint.SeverityWarning}
	s := issue.String()
	if s == "" {
		t.Fatal("expected non-empty string representation")
	}
}
