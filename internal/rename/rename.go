package rename

// Package rename provides utilities for detecting and applying key renames
// across environment files based on diff results.

import (
	"errors"

	"github.com/user/envdiff/internal/diff"
)

// Mapping represents a single key rename from OldKey to NewKey.
type Mapping struct {
	OldKey string `json:"old_key"`
	NewKey string `json:"new_key"`
}

// Options controls the behaviour of Apply.
type Options struct {
	// DryRun prevents any mutations; Apply returns what would change.
	DryRun bool
	// FailOnMissing returns an error if OldKey is not present in src.
	FailOnMissing bool
}

// DefaultOptions returns sensible defaults for Apply.
func DefaultOptions() Options {
	return Options{
		DryRun:        false,
		FailOnMissing: false,
	}
}

// Apply renames keys in src according to the provided mappings and returns
// a new map with the renames applied. The original map is never modified.
// When DryRun is true the returned map reflects what would change but src
// itself remains untouched.
func Apply(src map[string]string, mappings []Mapping, opts Options) (map[string]string, error) {
	out := copyMap(src)

	for _, m := range mappings {
		if m.OldKey == "" || m.NewKey == "" {
			return nil, errors.New("rename: OldKey and NewKey must not be empty")
		}

		val, ok := out[m.OldKey]
		if !ok {
			if opts.FailOnMissing {
				return nil, errors.New("rename: key not found: " + m.OldKey)
			}
			continue
		}

		if opts.DryRun {
			continue
		}

		delete(out, m.OldKey)
		out[m.NewKey] = val
	}

	return out, nil
}

// FromDiff derives rename Mappings from diff results produced by the suggest
// package or similar sources. It pairs a MissingInRight result with a
// MissingInLeft result that shares the same value, treating them as a rename.
func FromDiff(results []diff.Result) []Mapping {
	// index missing-in-left results by value so we can look them up quickly.
	byValue := make(map[string]string) // value -> key
	for _, r := range results {
		if r.Kind == diff.MissingInLeft {
			byValue[r.RightValue] = r.Key
		}
	}

	var mappings []Mapping
	for _, r := range results {
		if r.Kind != diff.MissingInRight {
			continue
		}
		if newKey, ok := byValue[r.LeftValue]; ok {
			mappings = append(mappings, Mapping{OldKey: r.Key, NewKey: newKey})
		}
	}
	return mappings
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
