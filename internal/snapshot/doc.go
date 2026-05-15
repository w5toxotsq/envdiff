// Package snapshot provides functionality for saving and loading
// envdiff comparison results to disk as JSON snapshots.
//
// Snapshots capture the diff results between two .env files at a
// specific point in time, allowing users to track changes over time
// or compare results across CI runs.
//
// Example usage:
//
//	err := snapshot.Save(".envdiff-snapshot.json", ".env", ".env.prod", results)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	snap, err := snapshot.Load(".envdiff-snapshot.json")
//	if err != nil {
//		log.Fatal(err)
//	}
package snapshot
