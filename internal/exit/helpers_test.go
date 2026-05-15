package exit_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/exit"
)

// TestFromDiff_ReturnsCode verifies that FromDiff always returns a Code
// whose integer value matches the expected Unix convention.
func TestFromDiff_ReturnsCode(t *testing.T) {
	tests := []struct {
		name            string
		hasDifferences  bool
		wantCode        exit.Code
		wantInt         int
	}{
		{"no differences", false, exit.OK, 0},
		{"has differences", true, exit.Differences, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := exit.FromDiff(tt.hasDifferences)
			if got != tt.wantCode {
				t.Errorf("FromDiff(%v): got code %d, want %d", tt.hasDifferences, got, tt.wantCode)
			}
			if int(got) != tt.wantInt {
				t.Errorf("FromDiff(%v): got int %d, want %d", tt.hasDifferences, int(got), tt.wantInt)
			}
		})
	}
}

// TestCode_String_AllCodes ensures every defined constant has a non-empty,
// non-"unknown" string representation.
func TestCode_String_AllCodes(t *testing.T) {
	known := []exit.Code{exit.OK, exit.Err, exit.Differences}
	for _, c := range known {
		s := c.String()
		if s == "" || s == "unknown" {
			t.Errorf("Code(%d).String() = %q; want a descriptive label", c, s)
		}
	}
}
