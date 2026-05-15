package validate_test

import (
	"testing"

	"github.com/envdiff/envdiff/internal/diff"
	"github.com/envdiff/envdiff/internal/validate"
)

func TestDefaultOptions_Values(t *testing.T) {
	opts := validate.DefaultOptions()
	if !opts.CheckKeyFormat {
		t.Error("expected CheckKeyFormat to be true by default")
	}
	if !opts.CheckEmptyValue {
		t.Error("expected CheckEmptyValue to be true by default")
	}
	if opts.CheckURLValues {
		t.Error("expected CheckURLValues to be false by default")
	}
}

func TestIssue_String(t *testing.T) {
	issue := validate.Issue{Key: "MY_KEY", Message: "some problem"}
	got := issue.String()
	want := "MY_KEY: some problem"
	if got != want {
		t.Errorf("Issue.String() = %q, want %q", got, want)
	}
}

func TestRun_MultipleIssuesSameKey(t *testing.T) {
	results := []diff.Result{
		{Key: "bad-key", Kind: diff.KindMismatch, Left: "", Right: ""},
	}
	opts := validate.DefaultOptions()
	issues := validate.Run(results, opts)
	// Expect: bad key format + left empty + right empty = 3 issues
	if len(issues) != 3 {
		t.Fatalf("expected 3 issues for bad-key with empty values, got %d: %v", len(issues), issues)
	}
}

func TestRun_URLCheckDisabled_NoSchemeIssue(t *testing.T) {
	results := []diff.Result{
		{Key: "SERVICE_URL", Kind: diff.KindMismatch, Left: "http://a", Right: "https://b"},
	}
	opts := validate.DefaultOptions()
	opts.CheckURLValues = false
	issues := validate.Run(results, opts)
	for _, iss := range issues {
		if iss.Key == "SERVICE_URL" && contains(iss.Message, "scheme") {
			t.Errorf("unexpected scheme issue when CheckURLValues is disabled")
		}
	}
}
