package exit_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/exit"
)

func TestCode_String_OK(t *testing.T) {
	if got := exit.OK.String(); got != "ok" {
		t.Errorf("expected \"ok\", got %q", got)
	}
}

func TestCode_String_Err(t *testing.T) {
	if got := exit.Err.String(); got != "error" {
		t.Errorf("expected \"error\", got %q", got)
	}
}

func TestCode_String_Differences(t *testing.T) {
	if got := exit.Differences.String(); got != "differences" {
		t.Errorf("expected \"differences\", got %q", got)
	}
}

func TestCode_String_Unknown(t *testing.T) {
	unknown := exit.Code(99)
	if got := unknown.String(); got != "unknown" {
		t.Errorf("expected \"unknown\", got %q", got)
	}
}

func TestFromDiff_NoDifferences(t *testing.T) {
	code := exit.FromDiff(false)
	if code != exit.OK {
		t.Errorf("expected OK (%d), got %d", exit.OK, code)
	}
}

func TestFromDiff_WithDifferences(t *testing.T) {
	code := exit.FromDiff(true)
	if code != exit.Differences {
		t.Errorf("expected Differences (%d), got %d", exit.Differences, code)
	}
}

func TestCode_IntValues(t *testing.T) {
	if int(exit.OK) != 0 {
		t.Errorf("OK should be 0, got %d", exit.OK)
	}
	if int(exit.Err) != 1 {
		t.Errorf("Err should be 1, got %d", exit.Err)
	}
	if int(exit.Differences) != 2 {
		t.Errorf("Differences should be 2, got %d", exit.Differences)
	}
}

func TestCode_String_AllKnown(t *testing.T) {
	tests := []struct {
		code exit.Code
		want string
	}{
		{exit.OK, "ok"},
		{exit.Err, "error"},
		{exit.Differences, "differences"},
	}
	for _, tt := range tests {
		if got := tt.code.String(); got != tt.want {
			t.Errorf("Code(%d).String() = %q, want %q", tt.code, got, tt.want)
		}
	}
}
