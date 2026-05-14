package report

import (
	"encoding/json"
	"io"

	"github.com/user/envdiff/internal/diff"
)

type jsonReport struct {
	Left    string        `json:"left"`
	Right   string        `json:"right"`
	Total   int           `json:"total_differences"`
	Results []jsonResult  `json:"results"`
}

type jsonResult struct {
	Key        string `json:"key"`
	Kind       string `json:"kind"`
	LeftValue  string `json:"left_value,omitempty"`
	RightValue string `json:"right_value,omitempty"`
}

func renderJSON(w io.Writer, results []diff.Result, leftName, rightName string, opts Options) error {
	jResults := make([]jsonResult, 0, len(results))
	for _, r := range results {
		jr := jsonResult{
			Key:  r.Key,
			Kind: kindString(r.Kind),
		}
		if opts.ShowValues {
			jr.LeftValue = r.LeftValue
			jr.RightValue = r.RightValue
		}
		jResults = append(jResults, jr)
	}

	report := jsonReport{
		Left:    leftName,
		Right:   rightName,
		Total:   len(results),
		Results: jResults,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func kindString(k diff.DiffKind) string {
	switch k {
	case diff.MissingInRight:
		return "missing_in_right"
	case diff.MissingInLeft:
		return "missing_in_left"
	case diff.ValueMismatch:
		return "value_mismatch"
	default:
		return "unknown"
	}
}
