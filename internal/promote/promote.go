// Package promote provides utilities for promoting env values
// from one environment to another, filling in missing keys.
package promote

import (
	"fmt"

	"github.com/user/envdiff/internal/diff"
)

// Options controls the behaviour of a promotion run.
type Options struct {
	// DryRun reports what would change without modifying Dst.
	DryRun bool

	// Overwrite replaces keys that already exist in Dst.
	Overwrite bool
}

// DefaultOptions returns the recommended Options for promote.Run.
func DefaultOptions() Options {
	return Options{
		DryRun:    false,
		Overwrite: false,
	}
}

// Change describes a single key that was (or would be) promoted.
type Change struct {
	// Key is the env key being promoted.
	Key string
	// Value is the value written into Dst.
	Value string
	// Overwritten is true when a pre-existing value in Dst was replaced.
	Overwritten bool
}

// Result holds the outcome of a Promote call.
type Result struct {
	Changes []Change
}

// Promote copies keys from src into dst according to the diff results.
// Only entries whose Kind is KindMissingInRight (present in src, absent in dst)
// are promoted unless opts.Overwrite is also set for KindMismatch entries.
// When opts.DryRun is true, dst is never modified.
func Promote(src, dst map[string]string, results []diff.Result, opts Options) (Result, error) {
	if src == nil {
		return Result{}, fmt.Errorf("promote: src map must not be nil")
	}
	if dst == nil {
		return Result{}, fmt.Errorf("promote: dst map must not be nil")
	}

	var changes []Change

	for _, r := range results {
		switch r.Kind {
		case diff.KindMissingInRight:
			v := src[r.Key]
			if !opts.DryRun {
				dst[r.Key] = v
			}
			changes = append(changes, Change{Key: r.Key, Value: v})

		case diff.KindMismatch:
			if !opts.Overwrite {
				continue
			}
			v := src[r.Key]
			if !opts.DryRun {
				dst[r.Key] = v
			}
			changes = append(changes, Change{Key: r.Key, Value: v, Overwritten: true})
		}
	}

	return Result{Changes: changes}, nil
}
