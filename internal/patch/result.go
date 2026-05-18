package patch

import "fmt"

// Summary returns a short human-readable description of the patch result.
func (r Result) Summary() string {
	if r.DryRun {
		return fmt.Sprintf("dry-run: %d change(s) would be applied", len(r.Changes))
	}
	return fmt.Sprintf("%d change(s) applied", len(r.Changes))
}

// HasChanges reports whether any changes were recorded.
func (r Result) HasChanges() bool {
	return len(r.Changes) > 0
}

// ByKind returns only the changes that match the given diff.Kind.
func (r Result) ByKind(k interface{ String() string }) []Change {
	// We accept interface{} to avoid a cyclic import concern; callers pass diff.Kind.
	target := k.String()
	var out []Change
	for _, c := range r.Changes {
		if c.Kind.String() == target {
			out = append(out, c)
		}
	}
	return out
}
