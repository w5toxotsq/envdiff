// Package drift compares a previously saved snapshot against a current
// environment map to surface keys that have been added, removed, or
// modified since the snapshot was taken.
//
// Typical usage:
//
//	snap, err := snapshot.Load("baseline.json")
//	if err != nil { ... }
//
//	current, err := parser.ParseFile(".env")
//	if err != nil { ... }
//
//	changes := drift.Detect(snap, current)
//	for _, c := range changes {
//		fmt.Println(c.Kind, c.Key)
//	}
package drift
