// Package report provides formatting and rendering of envdiff comparison results.
//
// It supports multiple output formats including plain text and JSON.
// The text format is designed for human consumption in a terminal, while
// the JSON format is suitable for machine processing or CI pipelines.
//
// Usage:
//
//	results := diff.Compare(left, right)
//	opts := report.DefaultOptions()
//	err := report.Render(os.Stdout, results, ".env.dev", ".env.prod", opts)
//
// To enable JSON output:
//
//	opts.Format = report.FormatJSON
//
// To reveal actual values in the report (use with caution in CI logs):
//
//	opts.ShowValues = true
package report
