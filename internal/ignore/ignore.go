// Package ignore provides functionality for filtering out specific keys
// from diff results based on user-defined ignore rules.
package ignore

import "github.com/your-org/envdiff/internal/diff"

// Options holds configuration for the ignore filter.
type Options struct {
	// Keys is a list of exact key names to ignore.
	Keys []string
	// Patterns is a list of glob-style prefix patterns to ignore.
	Patterns []string
}

// Apply removes any diff.Result entries whose key matches an ignored key
// or an ignored pattern prefix. The original slice is not modified.
func Apply(results []diff.Result, opts Options) []diff.Result {
	if len(opts.Keys) == 0 && len(opts.Patterns) == 0 {
		return results
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	out := make([]diff.Result, 0, len(results))
	for _, r := range results {
		if _, ignored := keySet[r.Key]; ignored {
			continue
		}
		if matchesAnyPattern(r.Key, opts.Patterns) {
			continue
		}
		out = append(out, r)
	}
	return out
}

// matchesAnyPattern returns true if key starts with any of the given patterns.
func matchesAnyPattern(key string, patterns []string) bool {
	for _, p := range patterns {
		if len(key) >= len(p) && key[:len(p)] == p {
			return true
		}
	}
	return false
}
