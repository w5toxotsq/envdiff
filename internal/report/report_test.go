package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/report"
)

func TestRenderText_NoDifferences(t *testing.T) {
	var buf bytes.Buffer
	opts := report.DefaultOptions()
	err := report.Render(&buf, []diff.Result{}, ".env.dev", ".env.prod", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected 'No differences' in output, got: %s", buf.String())
	}
}

func TestRenderText_MissingInRight(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", Kind: diff.MissingInRight},
	}
	var buf bytes.Buffer
	opts := report.DefaultOptions()
	err := report.Render(&buf, results, "left", "right", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected key DB_HOST in output")
	}
	if !strings.Contains(out, "-") {
		t.Errorf("expected '-' marker for missing_in_right")
	}
}

func TestRenderText_ValueMismatch_HiddenValues(t *testing.T) {
	results := []diff.Result{
		{Key: "SECRET", Kind: diff.ValueMismatch, LeftValue: "abc", RightValue: "xyz"},
	}
	var buf bytes.Buffer
	opts := report.DefaultOptions()
	opts.ShowValues = false
	report.Render(&buf, results, "l", "r", opts)
	out := buf.String()
	if strings.Contains(out, "abc") || strings.Contains(out, "xyz") {
		t.Errorf("values should be hidden when ShowValues=false")
	}
	if !strings.Contains(out, "values differ") {
		t.Errorf("expected 'values differ' text")
	}
}

func TestRenderText_ValueMismatch_ShownValues(t *testing.T) {
	results := []diff.Result{
		{Key: "PORT", Kind: diff.ValueMismatch, LeftValue: "3000", RightValue: "8080"},
	}
	var buf bytes.Buffer
	opts := report.DefaultOptions()
	opts.ShowValues = true
	report.Render(&buf, results, "l", "r", opts)
	out := buf.String()
	if !strings.Contains(out, "3000") || !strings.Contains(out, "8080") {
		t.Errorf("expected values to be shown when ShowValues=true, got: %s", out)
	}
}

func TestRenderText_SortedKeys(t *testing.T) {
	results := []diff.Result{
		{Key: "Z_KEY", Kind: diff.MissingInRight},
		{Key: "A_KEY", Kind: diff.MissingInRight},
		{Key: "M_KEY", Kind: diff.MissingInLeft},
	}
	var buf bytes.Buffer
	opts := report.DefaultOptions()
	report.Render(&buf, results, "l", "r", opts)
	out := buf.String()
	aIdx := strings.Index(out, "A_KEY")
	mIdx := strings.Index(out, "M_KEY")
	zIdx := strings.Index(out, "Z_KEY")
	if !(aIdx < mIdx && mIdx < zIdx) {
		t.Errorf("expected keys to be sorted alphabetically")
	}
}
