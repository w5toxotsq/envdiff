// Package filter provides utilities for filtering diff results
// based on key patterns, prefixes, or kinds of differences.
package filter

import (
	"strings"

	"github.com/yourusername/envdiff/internal/diff"
)

// Options holds configuration for filtering diff results.
type Options struct {
	// Prefix restricts results to keys with this prefix (case-insensitive).
	Prefix string

	// OnlyMissing restricts results to missing keys (left or right).
	OnlyMissing bool

	// OnlyMismatched restricts results to value mismatches.
	OnlyMismatched bool
}

// Apply returns a new Result containing only the entries that match
// the given Options. The original Result is not modified.
func Apply(result diff.Result, opts Options) diff.Result {
	filtered := make([]diff.Entry, 0, len(result.Entries))

	for _, entry := range result.Entries {
		if !matchesPrefix(entry.Key, opts.Prefix) {
			continue
		}
		if opts.OnlyMissing && entry.Kind == diff.KindMismatch {
			continue
		}
		if opts.OnlyMismatched && entry.Kind != diff.KindMismatch {
			continue
		}
		filtered = append(filtered, entry)
	}

	return diff.Result{Entries: filtered}
}

// matchesPrefix returns true if key starts with prefix (case-insensitive),
// or if prefix is empty.
func matchesPrefix(key, prefix string) bool {
	if prefix == "" {
		return true
	}
	return strings.HasPrefix(strings.ToUpper(key), strings.ToUpper(prefix))
}
