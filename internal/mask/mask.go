// Package mask provides utilities for selectively masking environment variable
// values before they are displayed or stored, based on key patterns and an
// explicit allowlist.
package mask

import "strings"

// Options controls which keys are masked and how.
type Options struct {
	// Enabled toggles masking entirely. When false, all values are returned as-is.
	Enabled bool

	// Patterns is a list of substrings; any key containing one of these
	// substrings (case-insensitive) will have its value masked.
	Patterns []string

	// Allowlist is a set of exact key names that are never masked, even if
	// they match a pattern.
	Allowlist []string

	// Mask is the replacement string used when a value is hidden.
	// Defaults to "***" when empty.
	Mask string
}

// DefaultOptions returns a sensible default configuration that masks common
// secret-like keys.
func DefaultOptions() Options {
	return Options{
		Enabled:  true,
		Patterns: []string{"secret", "password", "token", "key", "auth", "credential"},
		Mask:     "***",
	}
}

// Value returns the masked version of value for the given key, according to
// opts. If masking is disabled, or the key is in the allowlist, the original
// value is returned unchanged.
func Value(key, value string, opts Options) string {
	if !opts.Enabled {
		return value
	}

	mask := opts.Mask
	if mask == "" {
		mask = "***"
	}

	if isAllowed(key, opts.Allowlist) {
		return value
	}

	if matchesPattern(key, opts.Patterns) {
		return mask
	}

	return value
}

// isAllowed reports whether key appears in the allowlist (case-insensitive).
func isAllowed(key string, allowlist []string) bool {
	lower := strings.ToLower(key)
	for _, a := range allowlist {
		if strings.ToLower(a) == lower {
			return true
		}
	}
	return false
}

// matchesPattern reports whether key contains any of the given patterns
// (case-insensitive substring match).
func matchesPattern(key string, patterns []string) bool {
	lower := strings.ToLower(key)
	for _, p := range patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}
