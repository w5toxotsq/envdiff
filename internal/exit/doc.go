// Package exit defines standardised exit codes for the envdiff command-line
// tool.
//
// Using named constants instead of raw integers keeps the main package readable
// and ensures that the same values are used consistently across the codebase
// and in documentation.
//
// Exit code summary:
//
//	0  OK           — ran successfully, no differences detected
//	1  Err          — a runtime or usage error prevented completion
//	2  Differences  — ran successfully, one or more differences detected
//
// The helper [FromDiff] converts a boolean "has differences" flag into the
// correct [Code] so callers do not need to write the conditional themselves.
package exit
