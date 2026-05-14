// Package filter provides post-processing utilities for diff.Result values.
//
// After comparing two .env files with the diff package, callers may want to
// narrow the output to a specific subset of keys — for example, only keys
// that share a common prefix (e.g. "DB_"), or only keys that are outright
// missing rather than having mismatched values.
//
// Usage:
//
//	result := diff.Compare(left, right)
//	filtered := filter.Apply(result, filter.Options{
//		Prefix:      "DB_",
//		OnlyMissing: true,
//	})
//
// Options may be combined; all conditions must be satisfied for an entry to
// appear in the filtered output.
package filter
