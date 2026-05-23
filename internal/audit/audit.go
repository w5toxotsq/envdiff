package audit

import (
	"fmt"
	"time"

	"github.com/your/envdiff/internal/diff"
)

// Entry represents a single audit log entry recording a diff operation.
type Entry struct {
	Timestamp time.Time    `json:"timestamp"`
	LeftFile  string       `json:"left_file"`
	RightFile string       `json:"right_file"`
	Total     int          `json:"total"`
	Missing   int          `json:"missing"`
	Mismatch  int          `json:"mismatch"`
	Results   []diff.Result `json:"results,omitempty"`
}

// Options controls audit log behaviour.
type Options struct {
	// IncludeResults embeds the full diff results in each entry.
	IncludeResults bool
}

// DefaultOptions returns sensible defaults for audit logging.
func DefaultOptions() Options {
	return Options{
		IncludeResults: false,
	}
}

// Build constructs an Entry from a diff result set.
func Build(left, right string, results []diff.Result, opts Options) Entry {
	e := Entry{
		Timestamp: time.Now().UTC(),
		LeftFile:  left,
		RightFile: right,
		Total:     len(results),
	}

	for _, r := range results {
		switch r.Kind {
		case diff.MissingInLeft, diff.MissingInRight:
			e.Missing++
		case diff.ValueMismatch:
			e.Mismatch++
		}
	}

	if opts.IncludeResults {
		e.Results = results
	}

	return e
}

// Summary returns a human-readable one-line summary of the entry.
func (e Entry) Summary() string {
	return fmt.Sprintf(
		"[%s] %s vs %s — total:%d missing:%d mismatch:%d",
		e.Timestamp.Format(time.RFC3339),
		e.LeftFile,
		e.RightFile,
		e.Total,
		e.Missing,
		e.Mismatch,
	)
}
