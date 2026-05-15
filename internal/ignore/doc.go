// Package ignore provides key-based and pattern-based ignore rules for
// envdiff diff results.
//
// It allows users to exclude specific environment variable keys or groups
// of keys (by prefix) from comparison output. This is useful when certain
// keys are intentionally different across environments, such as secrets or
// machine-specific identifiers.
//
// Usage:
//
//	opts := ignore.Options{
//		Keys:     []string{"SECRET_KEY", "MACHINE_ID"},
//		Patterns: []string{"CI_", "INTERNAL_"},
//	}
//	filtered := ignore.Apply(results, opts)
package ignore
