// Package suggest provides heuristic-based key suggestions for missing
// or mismatched diff results. Given a set of known keys, it attempts to
// find close matches using edit-distance comparison.
package suggest

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Suggestion pairs a diff result with a list of candidate key names
// from the opposite environment that are likely typos or renames.
type Suggestion struct {
	Result     diff.Result
	Candidates []string
}

// Options controls how suggestions are generated.
type Options struct {
	// MaxDistance is the maximum Levenshtein distance to consider a match.
	// Defaults to 2.
	MaxDistance int
	// MaxCandidates limits how many suggestions are returned per result.
	// Defaults to 3.
	MaxCandidates int
}

// DefaultOptions returns sensible defaults for suggestion generation.
func DefaultOptions() Options {
	return Options{
		MaxDistance:   2,
		MaxCandidates: 3,
	}
}

// Run analyses the given results and, for each missing key, searches
// the pool of known keys for close matches. Only results of kind
// KindMissingInRight or KindMissingInLeft are considered.
func Run(results []diff.Result, knownKeys []string, opts Options) []Suggestion {
	if opts.MaxDistance <= 0 {
		opts.MaxDistance = DefaultOptions().MaxDistance
	}
	if opts.MaxCandidates <= 0 {
		opts.MaxCandidates = DefaultOptions().MaxCandidates
	}

	var suggestions []Suggestion
	for _, r := range results {
		if r.Kind != diff.KindMissingInRight && r.Kind != diff.KindMissingInLeft {
			continue
		}
		candidates := closestKeys(r.Key, knownKeys, opts.MaxDistance, opts.MaxCandidates)
		if len(candidates) > 0 {
			suggestions = append(suggestions, Suggestion{
				Result:     r,
				Candidates: candidates,
			})
		}
	}
	return suggestions
}

// closestKeys returns keys from pool whose Levenshtein distance to target
// is within maxDist, sorted by distance ascending.
func closestKeys(target string, pool []string, maxDist, limit int) []string {
	type scored struct {
		key  string
		dist int
	}
	var matches []scored
	norm := strings.ToLower(target)
	for _, k := range pool {
		if strings.EqualFold(k, target) {
			continue
		}
		d := levenshtein(norm, strings.ToLower(k))
		if d <= maxDist {
			matches = append(matches, scored{k, d})
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].dist < matches[j].dist
	})
	if len(matches) > limit {
		matches = matches[:limit]
	}
	out := make([]string, len(matches))
	for i, m := range matches {
		out[i] = m.key
	}
	return out
}

// levenshtein computes the edit distance between two strings.
func levenshtein(a, b string) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	row := make([]int, lb+1)
	for j := range row {
		row[j] = j
	}
	for i := 1; i <= la; i++ {
		prev := row[0]
		row[0] = i
		for j := 1; j <= lb; j++ {
			old := row[j]
			if a[i-1] == b[j-1] {
				row[j] = prev
			} else {
				row[j] = 1 + min3(prev, row[j], row[j-1])
			}
			prev = old
		}
	}
	return row[lb]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
