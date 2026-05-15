package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/your/envdiff/internal/diff"
)

// Snapshot represents a saved diff result at a point in time.
type Snapshot struct {
	CreatedAt time.Time    `json:"created_at"`
	LeftFile  string       `json:"left_file"`
	RightFile string       `json:"right_file"`
	Results   []diff.Result `json:"results"`
}

// Save writes a snapshot of the given diff results to path.
func Save(path, leftFile, rightFile string, results []diff.Result) error {
	s := Snapshot{
		CreatedAt: time.Now().UTC(),
		LeftFile:  leftFile,
		RightFile: rightFile,
		Results:   results,
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %q: %w", path, err)
	}

	return nil
}

// Load reads a snapshot from path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %q: %w", path, err)
	}

	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}

	return &s, nil
}
