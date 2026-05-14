package report_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/report"
)

func TestRenderJSON_Structure(t *testing.T) {
	results := []diff.Result{
		{Key: "API_KEY", Kind: diff.MissingInLeft},
		{Key: "DB_URL", Kind: diff.ValueMismatch, LeftValue: "local", RightValue: "prod"},
	}
	var buf bytes.Buffer
	opts := report.Options{
		Format:     report.FormatJSON,
		ShowValues: true,
	}
	err := report.Render(&buf, results, ".env.dev", ".env.prod", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out struct {
		Left    string `json:"left"`
		Right   string `json:"right"`
		Total   int    `json:"total_differences"`
		Results []struct {
			Key        string `json:"key"`
			Kind       string `json:"kind"`
			LeftValue  string `json:"left_value"`
			RightValue string `json:"right_value"`
		} `json:"results"`
	}
	if err := json.NewDecoder(&buf).Decode(&out); err != nil {
		t.Fatalf("failed to decode JSON output: %v", err)
	}
	if out.Left != ".env.dev" {
		t.Errorf("expected left='.env.dev', got %q", out.Left)
	}
	if out.Total != 2 {
		t.Errorf("expected total=2, got %d", out.Total)
	}
	if len(out.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out.Results))
	}
}

func TestRenderJSON_HiddenValues(t *testing.T) {
	results := []diff.Result{
		{Key: "SECRET", Kind: diff.ValueMismatch, LeftValue: "s3cr3t", RightValue: "other"},
	}
	var buf bytes.Buffer
	opts := report.Options{
		Format:     report.FormatJSON,
		ShowValues: false,
	}
	report.Render(&buf, results, "l", "r", opts)
	raw := buf.String()
	if bytes.Contains(buf.Bytes(), []byte("s3cr3t")) {
		t.Errorf("secret value should not appear in JSON when ShowValues=false, got: %s", raw)
	}
}
