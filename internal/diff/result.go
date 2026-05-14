// Package diff provides functionality for comparing two environment maps
// and producing structured difference results.
package diff

// Kind represents the type of difference found between two env maps.
type Kind int

const (
	// MissingInRight indicates the key exists in the left map but not the right.
	MissingInRight Kind = iota
	// MissingInLeft indicates the key exists in the right map but not the left.
	MissingInLeft
	// ValueMismatch indicates the key exists in both maps but with different values.
	ValueMismatch
)

// String returns a human-readable label for the Kind.
func (k Kind) String() string {
	switch k {
	case MissingInRight:
		return "missing_in_right"
	case MissingInLeft:
		return "missing_in_left"
	case ValueMismatch:
		return "value_mismatch"
	default:
		return "unknown"
	}
}

// Entry represents a single difference between two environment maps.
type Entry struct {
	// Key is the environment variable name.
	Key string
	// Kind describes the nature of the difference.
	Kind Kind
	// LeftValue is the value from the left (base) env map. Empty if missing.
	LeftValue string
	// RightValue is the value from the right (target) env map. Empty if missing.
	RightValue string
}

// Result holds the full output of a diff comparison.
type Result struct {
	// Entries contains all detected differences, in the order they were found.
	Entries []Entry
}

// HasDifferences returns true if any differences were found.
func (r Result) HasDifferences() bool {
	return len(r.Entries) > 0
}

// ByKind returns only the entries matching the given Kind.
func (r Result) ByKind(k Kind) []Entry {
	var out []Entry
	for _, e := range r.Entries {
		if e.Kind == k {
			out = append(out, e)
		}
	}
	return out
}
