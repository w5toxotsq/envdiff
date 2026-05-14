// Package report formats and renders diff results for human-readable output.
package report

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format for a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Options configures report rendering behaviour.
type Options struct {
	Format      Format
	ShowValues  bool
	ColorOutput bool
}

// DefaultOptions returns sensible default report options.
func DefaultOptions() Options {
	return Options{
		Format:      FormatText,
		ShowValues:  false,
		ColorOutput: true,
	}
}

// Render writes a formatted diff report to w.
func Render(w io.Writer, results []diff.Result, leftName, rightName string, opts Options) error {
	switch opts.Format {
	case FormatJSON:
		return renderJSON(w, results, leftName, rightName, opts)
	default:
		return renderText(w, results, leftName, rightName, opts)
	}
}

func renderText(w io.Writer, results []diff.Result, leftName, rightName string, opts Options) error {
	if len(results) == 0 {
		_, err := fmt.Fprintf(w, "✓ No differences found between %s and %s\n", leftName, rightName)
		return err
	}

	sorted := make([]diff.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	fmt.Fprintf(w, "Comparing %s ↔ %s\n", leftName, rightName)
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 40))

	for _, r := range sorted {
		switch r.Kind {
		case diff.MissingInRight:
			fmt.Fprintf(w, "  - %-30s only in %s\n", r.Key, leftName)
		case diff.MissingInLeft:
			fmt.Fprintf(w, "  + %-30s only in %s\n", r.Key, rightName)
		case diff.ValueMismatch:
			if opts.ShowValues {
				fmt.Fprintf(w, "  ~ %-30s %q ≠ %q\n", r.Key, r.LeftValue, r.RightValue)
			} else {
				fmt.Fprintf(w, "  ~ %-30s values differ\n", r.Key)
			}
		}
	}

	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 40))
	fmt.Fprintf(w, "Total: %d difference(s)\n", len(results))
	return nil
}
