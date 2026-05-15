// Package multienv provides utilities for loading and comparing
// multiple .env files across different environments in a single pass.
package multienv

import (
	"fmt"

	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/parser"
)

// EnvSet holds named environment maps loaded from files.
type EnvSet map[string]parser.EnvMap

// LoadOptions configures how files are loaded.
type LoadOptions struct {
	// SkipExtensionCheck disables the .env extension validation.
	SkipExtensionCheck bool
}

// Load reads multiple named env files and returns an EnvSet.
// The names slice must correspond 1-to-1 with the paths slice.
func Load(names, paths []string, opts LoadOptions) (EnvSet, error) {
	if len(names) != len(paths) {
		return nil, fmt.Errorf("multienv: names and paths must have the same length")
	}

	set := make(EnvSet, len(paths))

	for i, p := range paths {
		m, err := loader.Load(p, loader.Options{
			SkipExtensionCheck: opts.SkipExtensionCheck,
		})
		if err != nil {
			return nil, fmt.Errorf("multienv: loading %q: %w", p, err)
		}
		set[names[i]] = m
	}

	return set, nil
}

// Names returns a sorted list of environment names in the set.
func (s EnvSet) Names() []string {
	out := make([]string, 0, len(s))
	for k := range s {
		out = append(out, k)
	}
	sortStrings(out)
	return out
}

// AllKeys returns the union of all keys across every environment in the set.
func (s EnvSet) AllKeys() []string {
	seen := make(map[string]struct{})
	for _, m := range s {
		for k := range m {
			seen[k] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sortStrings(out)
	return out
}
