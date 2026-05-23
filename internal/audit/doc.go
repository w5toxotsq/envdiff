// Package audit provides structured audit-log entries for envdiff operations.
//
// An audit entry captures metadata about a comparison run — which files were
// compared, when the run occurred, and how many differences were found — so
// that teams can track environment drift over time.
//
// Basic usage:
//
//	results := diff.Compare(left, right)
//	entry   := audit.Build(".env.staging", ".env.production", results, audit.DefaultOptions())
//	fmt.Println(entry.Summary())
package audit
