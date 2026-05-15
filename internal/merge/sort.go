package merge

import "sort"

// sortStrings sorts a slice of strings in ascending order in-place.
func sortStrings(s []string) {
	sort.Strings(s)
}

// sortedKeys returns a sorted copy of the keys from the given map.
func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
