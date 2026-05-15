// Package lint provides checks for common .env file issues such as
// duplicate keys, empty values, and keys that don't follow naming conventions.
package lint

import (
	"fmt"
	"regexp"
	"strings"
)

// Severity indicates how serious a lint issue is.
type Severity string

const (
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// Issue describes a single lint finding for a specific key.
type Issue struct {
	Key      string
	Message  string
	Severity Severity
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s: %s", i.Severity, i.Key, i.Message)
}

// Options controls which lint checks are performed.
type Options struct {
	CheckNaming     bool // keys must match [A-Z][A-Z0-9_]*
	CheckEmptyValues bool // warn on keys with empty values
	CheckDuplicates bool // error on duplicate keys (from raw lines)
}

// DefaultOptions returns Options with all checks enabled.
func DefaultOptions() Options {
	return Options{
		CheckNaming:      true,
		CheckEmptyValues: true,
		CheckDuplicates:  true,
	}
}

var validKeyRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// Run executes lint checks against the provided env map and returns any issues found.
// keys is an ordered slice of keys (may contain duplicates for duplicate detection).
func Run(env map[string]string, keys []string, opts Options) []Issue {
	var issues []Issue

	if opts.CheckDuplicates {
		seen := make(map[string]int)
		for _, k := range keys {
			seen[k]++
		}
		for k, count := range seen {
			if count > 1 {
				issues = append(issues, Issue{
					Key:      k,
					Message:  fmt.Sprintf("duplicate key appears %d times", count),
					Severity: SeverityError,
				})
			}
		}
	}

	for k, v := range env {
		if opts.CheckNaming && !validKeyRe.MatchString(k) {
			issues = append(issues, Issue{
				Key:      k,
				Message:  "key does not match recommended naming convention [A-Z][A-Z0-9_]*",
				Severity: SeverityWarning,
			})
		}
		if opts.CheckEmptyValues && strings.TrimSpace(v) == "" {
			issues = append(issues, Issue{
				Key:      k,
				Message:  "value is empty",
				Severity: SeverityWarning,
			})
		}
	}

	return issues
}
