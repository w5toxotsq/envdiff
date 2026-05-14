// Package summary provides aggregated statistics over a set of diff results.
package summary

import "github.com/user/envdiff/internal/diff"

// Stats holds counts of each difference kind found across compared env files.
type Stats struct {
	Total      int
	Missing    int
	Extra      int
	Mismatched int
}

// HasDifferences reports whether any differences were found.
func (s Stats) HasDifferences() bool {
	return s.Total > 0
}

// Compute calculates summary statistics from a slice of diff results.
func Compute(results []diff.Result) Stats {
	var s Stats
	for _, r := range results {
		switch r.Kind {
		case diff.KindMissingInRight:
			s.Missing++
		case diff.KindMissingInLeft:
			s.Extra++
		case diff.KindValueMismatch:
			s.Mismatched++
		}
	}
	s.Total = s.Missing + s.Extra + s.Mismatched
	return s
}

// Breakdown returns a map of kind label to count for non-zero categories.
func Breakdown(s Stats) map[string]int {
	m := make(map[string]int)
	if s.Missing > 0 {
		m["missing"] = s.Missing
	}
	if s.Extra > 0 {
		m["extra"] = s.Extra
	}
	if s.Mismatched > 0 {
		m["mismatched"] = s.Mismatched
	}
	return m
}
