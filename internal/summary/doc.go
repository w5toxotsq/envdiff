// Package summary computes aggregated statistics over a collection of diff
// results produced by the diff package.
//
// It is intended to provide a quick overview of how many keys are missing,
// extra, or mismatched between two compared .env files without requiring the
// caller to iterate over individual results manually.
//
// Example usage:
//
//	results := diff.Compare(left, right)
//	stats := summary.Compute(results)
//	fmt.Printf("Total differences: %d\n", stats.Total)
//	for label, count := range summary.Breakdown(stats) {
//		fmt.Printf("  %s: %d\n", label, count)
//	}
package summary
