// Package sorter provides utilities for sorting diff results
// in a consistent and deterministic order for display and comparison.
package sorter

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// SortOrder defines the order in which results are sorted.
type SortOrder int

const (
	// ByKey sorts results alphabetically by key name.
	ByKey SortOrder = iota
	// ByKind sorts results grouped by their diff kind (missing, mismatch, etc.).
	ByKind
)

// Options configures how results are sorted.
type Options struct {
	Order SortOrder
}

// DefaultOptions returns the default sort options.
func DefaultOptions() Options {
	return Options{
		Order: ByKey,
	}
}

// Sort returns a new slice of diff results sorted according to the given options.
// The original slice is not modified.
func Sort(results []diff.Result, opts Options) []diff.Result {
	if len(results) == 0 {
		return results
	}

	out := make([]diff.Result, len(results))
	copy(out, results)

	switch opts.Order {
	case ByKind:
		sort.SliceStable(out, func(i, j int) bool {
			if out[i].Kind != out[j].Kind {
				return out[i].Kind < out[j].Kind
			}
			return out[i].Key < out[j].Key
		})
	default: // ByKey
		sort.SliceStable(out, func(i, j int) bool {
			return out[i].Key < out[j].Key
		})
	}

	return out
}
