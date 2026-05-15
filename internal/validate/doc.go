// Package validate provides post-diff validation checks for environment
// variable results. It inspects diff.Result values and reports issues such
// as non-standard key naming, empty values, or URL scheme mismatches.
//
// Validations are opt-in via Options so callers can enable only the checks
// relevant to their workflow.
package validate
