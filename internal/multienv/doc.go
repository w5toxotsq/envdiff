// Package multienv provides utilities for loading and comparing multiple
// .env files in a single operation.
//
// Instead of comparing exactly two files, multienv allows loading an
// arbitrary set of named environments (e.g. "dev", "staging", "prod")
// and querying the union of their keys.
//
// Example:
//
//	set, err := multienv.Load(
//		[]string{"dev", "prod"},
//		[]string{".env.dev", ".env.prod"},
//		multienv.LoadOptions{},
//	)
//	if err != nil { ... }
//
//	for _, key := range set.AllKeys() {
//		// inspect each key across environments
//	}
package multienv
