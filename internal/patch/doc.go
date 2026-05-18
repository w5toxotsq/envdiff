// Package patch applies diff results to a target .env file.
//
// Given a slice of [diff.Result] values produced by [diff.Compare], Apply
// writes missing or mismatched keys into the destination file.  A dry-run
// mode is available so callers can preview changes before committing them.
//
// Example
//
//	results, _ := diff.Compare(base, target)
//	out, err := patch.Apply(".env.staging", target, results, patch.DefaultOptions())
//	for _, c := range out.Changes {
//	    fmt.Printf("set %s = %s\n", c.Key, c.NewValue)
//	}
package patch
