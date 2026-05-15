// Package redact provides utilities for masking sensitive environment variable
// values in output, ensuring secrets are not accidentally exposed in reports
// or logs.
package redact

import "strings"

const defaultMask = "***"

// Options controls redaction behaviour.
type Options struct {
	// Enabled turns redaction on or off.
	Enabled bool
	// Mask is the string used to replace sensitive values.
	// Defaults to "***" when empty.
	Mask string
	// SensitiveSuffixes is a list of key suffixes (case-insensitive) that
	// should always be redacted regardless of the Enabled flag.
	SensitiveSuffixes []string
}

// DefaultOptions returns a sensible default redaction configuration.
func DefaultOptions() Options {
	return Options{
		Enabled: true,
		Mask:    defaultMask,
		SensitiveSuffixes: []string{
			"_SECRET",
			"_PASSWORD",
			"_PASSWD",
			"_TOKEN",
			"_PRIVATE_KEY",
			"_API_KEY",
		},
	}
}

// Value returns the masked value if redaction is enabled or the key matches a
// sensitive suffix; otherwise it returns the original value unchanged.
func Value(key, value string, opts Options) string {
	mask := opts.Mask
	if mask == "" {
		mask = defaultMask
	}

	if isSensitive(key, opts.SensitiveSuffixes) {
		return mask
	}

	if opts.Enabled {
		return mask
	}

	return value
}

// isSensitive reports whether the key ends with any of the given suffixes
// (comparison is case-insensitive).
func isSensitive(key string, suffixes []string) bool {
	upper := strings.ToUpper(key)
	for _, s := range suffixes {
		if strings.HasSuffix(upper, strings.ToUpper(s)) {
			return true
		}
	}
	return false
}
