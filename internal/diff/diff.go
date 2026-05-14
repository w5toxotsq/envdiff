// Package diff provides functionality for comparing two EnvMaps and
// reporting missing keys, extra keys, and mismatched values.
package diff

import "github.com/user/envdiff/internal/parser"

// ResultType indicates the kind of difference found between two env files.
type ResultType string

const (
	// MissingInRight indicates a key present in the left file but absent in the right.
	MissingInRight ResultType = "missing_in_right"
	// MissingInLeft indicates a key present in the right file but absent in the left.
	MissingInLeft ResultType = "missing_in_left"
	// ValueMismatch indicates a key present in both files but with different values.
	ValueMismatch ResultType = "value_mismatch"
)

// Entry represents a single difference between two env maps.
type Entry struct {
	Key        string
	Type       ResultType
	LeftValue  string
	RightValue string
}

// Result holds all differences found between two env maps.
type Result struct {
	Entries []Entry
}

// HasDiff returns true if any differences were found.
func (r *Result) HasDiff() bool {
	return len(r.Entries) > 0
}

// Compare compares two EnvMaps and returns a Result containing all differences.
// Keys present only in left are reported as MissingInRight.
// Keys present only in right are reported as MissingInLeft.
// Keys present in both but with different values are reported as ValueMismatch.
func Compare(left, right parser.EnvMap) *Result {
	result := &Result{}

	for key, leftVal := range left {
		rightVal, ok := right[key]
		if !ok {
			result.Entries = append(result.Entries, Entry{
				Key:       key,
				Type:      MissingInRight,
				LeftValue: leftVal,
			})
			continue
		}
		if leftVal != rightVal {
			result.Entries = append(result.Entries, Entry{
				Key:        key,
				Type:       ValueMismatch,
				LeftValue:  leftVal,
				RightValue: rightVal,
			})
		}
	}

	for key, rightVal := range right {
		if _, ok := left[key]; !ok {
			result.Entries = append(result.Entries, Entry{
				Key:        key,
				Type:       MissingInLeft,
				RightValue: rightVal,
			})
		}
	}

	return result
}
