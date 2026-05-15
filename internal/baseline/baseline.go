package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/your/envdiff/internal/diff"
)

// Baseline represents a saved reference state for a set of diff results.
type Baseline struct {
	CreatedAt time.Time     `json:"created_at"`
	Label     string        `json:"label,omitempty"`
	Results   []diff.Result `json:"results"`
}

// Save writes a Baseline to the given file path as JSON.
func Save(path, label string, results []diff.Result) error {
	if results == nil {
		results = []diff.Result{}
	}
	b := Baseline{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Results:   results,
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %q: %w", path, err)
	}
	return nil
}

// Load reads a Baseline from the given file path.
func Load(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: read %q: %w", path, err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &b, nil
}

// Compare returns results that are new relative to the baseline — i.e. present
// in current but not in the baseline (matched by key and kind).
func Compare(base *Baseline, current []diff.Result) []diff.Result {
	index := make(map[string]struct{}, len(base.Results))
	for _, r := range base.Results {
		index[indexKey(r)] = struct{}{}
	}
	var novel []diff.Result
	for _, r := range current {
		if _, found := index[indexKey(r)]; !found {
			novel = append(novel, r)
		}
	}
	return novel
}

func indexKey(r diff.Result) string {
	return fmt.Sprintf("%s::%s", r.Key, r.Kind)
}
