package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/envdiff/envdiff/internal/diff"
)

// Issue represents a single validation problem found in a diff result.
type Issue struct {
	Key     string
	Message string
}

func (i Issue) String() string {
	return fmt.Sprintf("%s: %s", i.Key, i.Message)
}

// Options controls which validations are performed.
type Options struct {
	CheckKeyFormat  bool // enforce UPPER_SNAKE_CASE
	CheckEmptyValue bool // warn when a matched key has an empty value in either env
	CheckURLValues  bool // warn when a value looks like a URL but schemes differ
}

// DefaultOptions returns a sensible default set of validation options.
func DefaultOptions() Options {
	return Options{
		CheckKeyFormat:  true,
		CheckEmptyValue: true,
		CheckURLValues:  false,
	}
}

var keyFormatRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// Run validates the provided diff results and returns any issues found.
func Run(results []diff.Result, opts Options) []Issue {
	var issues []Issue

	for _, r := range results {
		if opts.CheckKeyFormat {
			if !keyFormatRe.MatchString(r.Key) {
				issues = append(issues, Issue{
					Key:     r.Key,
					Message: "key does not follow UPPER_SNAKE_CASE convention",
				})
			}
		}

		if opts.CheckEmptyValue && r.Kind == diff.KindMismatch {
			if r.Left == "" {
				issues = append(issues, Issue{Key: r.Key, Message: "left value is empty"})
			}
			if r.Right == "" {
				issues = append(issues, Issue{Key: r.Key, Message: "right value is empty"})
			}
		}

		if opts.CheckURLValues && r.Kind == diff.KindMismatch {
			if looksLikeURL(r.Left) && looksLikeURL(r.Right) {
				leftScheme := urlScheme(r.Left)
				rightScheme := urlScheme(r.Right)
				if leftScheme != rightScheme {
					issues = append(issues, Issue{
						Key:     r.Key,
						Message: fmt.Sprintf("URL scheme mismatch: %q vs %q", leftScheme, rightScheme),
					})
				}
			}
		}
	}

	return issues
}

func looksLikeURL(v string) bool {
	return strings.Contains(v, "://")
}

func urlScheme(v string) string {
	if idx := strings.Index(v, "://"); idx >= 0 {
		return v[:idx]
	}
	return ""
}
