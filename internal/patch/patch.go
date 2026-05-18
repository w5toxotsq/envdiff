package patch

import (
	"fmt"
	"os"
	"strings"

	"github.com/your/envdiff/internal/diff"
)

// Options controls how a patch is applied.
type Options struct {
	// DryRun reports what would change without writing to disk.
	DryRun bool
	// OnlyMissing applies only missing-key entries (skips mismatches).
	OnlyMissing bool
}

// DefaultOptions returns an Options with safe defaults.
func DefaultOptions() Options {
	return Options{
		DryRun:      false,
		OnlyMissing: false,
	}
}

// Change describes a single key/value write operation.
type Change struct {
	Key      string
	OldValue string // empty when the key did not previously exist
	NewValue string
	Kind     diff.Kind
}

// Result is returned by Apply.
type Result struct {
	Changes []Change
	DryRun  bool
}

// Apply writes diff results into dst (a map representing the target .env
// contents) and, unless DryRun is set, serialises the updated map to path.
func Apply(path string, dst map[string]string, results []diff.Result, opts Options) (Result, error) {
	var changes []Change

	for _, r := range results {
		switch r.Kind {
		case diff.MissingInRight:
			changes = append(changes, Change{Key: r.Key, OldValue: "", NewValue: r.LeftValue, Kind: r.Kind})
		case diff.Mismatch:
			if opts.OnlyMissing {
				continue
			}
			changes = append(changes, Change{Key: r.Key, OldValue: r.RightValue, NewValue: r.LeftValue, Kind: r.Kind})
		}
	}

	if opts.DryRun {
		return Result{Changes: changes, DryRun: true}, nil
	}

	updated := copyMap(dst)
	for _, c := range changes {
		updated[c.Key] = c.NewValue
	}

	if err := writeEnvFile(path, updated); err != nil {
		return Result{}, fmt.Errorf("patch: write %s: %w", path, err)
	}

	return Result{Changes: changes, DryRun: false}, nil
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func writeEnvFile(path string, m map[string]string) error {
	var sb strings.Builder
	for k, v := range m {
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(v)
		sb.WriteByte('\n')
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
