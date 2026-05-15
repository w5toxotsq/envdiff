// Package suggest analyses diff results and recommends candidate key names
// that may be typos or renames of missing keys.
//
// It uses Levenshtein edit-distance to rank suggestions and returns only
// those within a configurable distance threshold, making it easy to spot
// common mistakes such as DB_PASWORD vs DB_PASSWORD.
//
// Usage:
//
//	opts := suggest.DefaultOptions()
//	suggestions := suggest.Run(results, allKnownKeys, opts)
//	for _, s := range suggestions {
//		fmt.Printf("missing %q — did you mean: %v?\n", s.Result.Key, s.Candidates)
//	}
package suggest
