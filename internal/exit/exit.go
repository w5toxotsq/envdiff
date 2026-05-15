// Package exit defines exit codes used by the envdiff CLI.
//
// Exit codes follow common Unix conventions:
//
//	0 - success, no differences found
//	1 - an error occurred (e.g. file not found, parse error)
//	2 - success, but differences were found
package exit

// Code represents a CLI exit code.
type Code int

const (
	// OK indicates the command ran successfully with no differences.
	OK Code = 0

	// Err indicates a runtime or usage error.
	Err Code = 1

	// Differences indicates the command ran successfully but differences were found.
	Differences Code = 2
)

// String returns a human-readable label for the exit code.
func (c Code) String() string {
	switch c {
	case OK:
		return "ok"
	case Err:
		return "error"
	case Differences:
		return "differences"
	default:
		return "unknown"
	}
}

// FromDiff returns the appropriate exit code based on whether
// any differences were detected.
func FromDiff(hasDifferences bool) Code {
	if hasDifferences {
		return Differences
	}
	return OK
}
