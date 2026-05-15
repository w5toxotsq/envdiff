package merge_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/merge"
)

func TestMerge_AddsNewKeys(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"B": "2"}

	r := merge.Merge(dst, src, merge.DefaultOptions())

	if r.Merged["A"] != "1" {
		t.Errorf("expected A=1, got %s", r.Merged["A"])
	}
	if r.Merged["B"] != "2" {
		t.Errorf("expected B=2, got %s", r.Merged["B"])
	}
	if len(r.Added) != 1 || r.Added[0] != "B" {
		t.Errorf("expected Added=[B], got %v", r.Added)
	}
}

func TestMerge_SkipsExistingWithoutOverwrite(t *testing.T) {
	dst := map[string]string{"A": "original"}
	src := map[string]string{"A": "new"}

	r := merge.Merge(dst, src, merge.DefaultOptions())

	if r.Merged["A"] != "original" {
		t.Errorf("expected A=original, got %s", r.Merged["A"])
	}
	if len(r.Skipped) != 1 || r.Skipped[0] != "A" {
		t.Errorf("expected Skipped=[A], got %v", r.Skipped)
	}
	if len(r.Overwritten) != 0 {
		t.Errorf("expected no overwritten keys, got %v", r.Overwritten)
	}
}

func TestMerge_OverwritesExistingWhenEnabled(t *testing.T) {
	dst := map[string]string{"A": "original"}
	src := map[string]string{"A": "new"}

	opts := merge.MergeOptions{Overwrite: true}
	r := merge.Merge(dst, src, opts)

	if r.Merged["A"] != "new" {
		t.Errorf("expected A=new, got %s", r.Merged["A"])
	}
	if len(r.Overwritten) != 1 || r.Overwritten[0] != "A" {
		t.Errorf("expected Overwritten=[A], got %v", r.Overwritten)
	}
	if len(r.Skipped) != 0 {
		t.Errorf("expected no skipped keys, got %v", r.Skipped)
	}
}

func TestMerge_DoesNotModifyDst(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"A": "2", "B": "3"}

	opts := merge.MergeOptions{Overwrite: true}
	merge.Merge(dst, src, opts)

	if dst["A"] != "1" {
		t.Errorf("dst was modified: A=%s", dst["A"])
	}
	if _, ok := dst["B"]; ok {
		t.Error("dst was modified: B should not exist")
	}
}

func TestMerge_EmptyMaps(t *testing.T) {
	r := merge.Merge(map[string]string{}, map[string]string{}, merge.DefaultOptions())

	if len(r.Merged) != 0 {
		t.Errorf("expected empty merged, got %v", r.Merged)
	}
	if len(r.Added) != 0 || len(r.Skipped) != 0 || len(r.Overwritten) != 0 {
		t.Error("expected all slices empty")
	}
}
