// Package promote copies env keys from one environment map into another,
// driven by a slice of diff.Result values.
//
// Typical usage:
//
//	results, _ := diff.Compare(staging, production)
//	promoted, err := promote.Promote(staging, production, results, promote.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, c := range promoted.Changes {
//		fmt.Printf("promoted %s\n", c.Key)
//	}
//
// Dry-run mode reports what would change without touching the destination map.
package promote
