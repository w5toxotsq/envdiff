// Package graph builds a dependency graph between .env files based on
// shared keys, allowing callers to visualise which environments are most
// similar or most divergent from one another.
package graph

import (
	"sort"
)

// Edge represents a weighted relationship between two environment files.
type Edge struct {
	// Left and Right are the names of the two environments being compared.
	Left  string
	Right string
	// SharedKeys is the number of keys present in both environments.
	SharedKeys int
	// TotalKeys is the union of all keys across both environments.
	TotalKeys int
	// Similarity is SharedKeys/TotalKeys expressed as a value in [0, 1].
	Similarity float64
}

// Build constructs a complete set of edges for all pairs of environments
// described by names and their associated key sets. The keys parameter maps
// an environment name to the set of keys it contains (values are ignored).
func Build(keys map[string]map[string]string) []Edge {
	names := sortedNames(keys)
	var edges []Edge

	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			left := names[i]
			right := names[j]
			shared, total := countKeys(keys[left], keys[right])
			sim := 0.0
			if total > 0 {
				sim = float64(shared) / float64(total)
			}
			edges = append(edges, Edge{
				Left:       left,
				Right:      right,
				SharedKeys: shared,
				TotalKeys:  total,
				Similarity: sim,
			})
		}
	}
	return edges
}

// countKeys returns the number of shared keys and the size of the union.
func countKeys(a, b map[string]string) (shared, total int) {
	union := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		union[k] = struct{}{}
	}
	for k := range b {
		union[k] = struct{}{}
	}
	total = len(union)
	for k := range a {
		if _, ok := b[k]; ok {
			shared++
		}
	}
	return shared, total
}

func sortedNames(m map[string]map[string]string) []string {
	names := make([]string, 0, len(m))
	for n := range m {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
