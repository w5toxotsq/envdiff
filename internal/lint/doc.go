// Package lint provides static analysis checks for .env file contents.
//
// It detects common problems such as:
//   - Keys that do not follow the recommended uppercase naming convention
//   - Keys with empty or whitespace-only values
//   - Duplicate keys within the same file
//
// Usage:
//
//	opts := lint.DefaultOptions()
//	issues := lint.Run(envMap, orderedKeys, opts)
//	for _, issue := range issues {
//		fmt.Println(issue)
//	}
package lint
