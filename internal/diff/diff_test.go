package diff_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

func TestCompare_NoDifferences(t *testing.T) {
	left := parser.EnvMap{"KEY": "value", "PORT": "8080"}
	right := parser.EnvMap{"KEY": "value", "PORT": "8080"}

	result := diff.Compare(left, right)
	if result.HasDiff() {
		t.Errorf("expected no differences, got %d", len(result.Entries))
	}
}

func TestCompare_MissingInRight(t *testing.T) {
	left := parser.EnvMap{"KEY": "value", "EXTRA": "only_left"}
	right := parser.EnvMap{"KEY": "value"}

	result := diff.Compare(left, right)
	if !result.HasDiff() {
		t.Fatal("expected differences, got none")
	}
	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
	e := result.Entries[0]
	if e.Type != diff.MissingInRight {
		t.Errorf("expected MissingInRight, got %s", e.Type)
	}
	if e.Key != "EXTRA" {
		t.Errorf("expected key EXTRA, got %s", e.Key)
	}
	if e.LeftValue != "only_left" {
		t.Errorf("unexpected left value: %s", e.LeftValue)
	}
}

func TestCompare_MissingInLeft(t *testing.T) {
	left := parser.EnvMap{"KEY": "value"}
	right := parser.EnvMap{"KEY": "value", "NEW": "only_right"}

	result := diff.Compare(left, right)
	if !result.HasDiff() {
		t.Fatal("expected differences, got none")
	}
	e := result.Entries[0]
	if e.Type != diff.MissingInLeft {
		t.Errorf("expected MissingInLeft, got %s", e.Type)
	}
	if e.Key != "NEW" {
		t.Errorf("expected key NEW, got %s", e.Key)
	}
	if e.RightValue != "only_right" {
		t.Errorf("unexpected right value: %s", e.RightValue)
	}
}

func TestCompare_ValueMismatch(t *testing.T) {
	left := parser.EnvMap{"DB_HOST": "localhost"}
	right := parser.EnvMap{"DB_HOST": "prod.db.example.com"}

	result := diff.Compare(left, right)
	if !result.HasDiff() {
		t.Fatal("expected differences, got none")
	}
	e := result.Entries[0]
	if e.Type != diff.ValueMismatch {
		t.Errorf("expected ValueMismatch, got %s", e.Type)
	}
	if e.LeftValue != "localhost" || e.RightValue != "prod.db.example.com" {
		t.Errorf("unexpected values: left=%s right=%s", e.LeftValue, e.RightValue)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	result := diff.Compare(parser.EnvMap{}, parser.EnvMap{})
	if result.HasDiff() {
		t.Error("expected no diff for two empty maps")
	}
}

func TestCompare_MultipleDiffTypes(t *testing.T) {
	left := parser.EnvMap{"SHARED": "old", "ONLY_LEFT": "gone"}
	right := parser.EnvMap{"SHARED": "new", "ONLY_RIGHT": "added"}

	result := diff.Compare(left, right)
	if !result.HasDiff() {
		t.Fatal("expected differences, got none")
	}
	if len(result.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result.Entries))
	}

	counts := map[diff.EntryType]int{}
	for _, e := range result.Entries {
		counts[e.Type]++
	}
	if counts[diff.ValueMismatch] != 1 {
		t.Errorf("expected 1 ValueMismatch, got %d", counts[diff.ValueMismatch])
	}
	if counts[diff.MissingInRight] != 1 {
		t.Errorf("expected 1 MissingInRight, got %d", counts[diff.MissingInRight])
	}
	if counts[diff.MissingInLeft] != 1 {
		t.Errorf("expected 1 MissingInLeft, got %d", counts[diff.MissingInLeft])
	}
}
