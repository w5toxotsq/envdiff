package rename_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/rename"
)

func TestApply_RenamesKey(t *testing.T) {
	src := map[string]string{"OLD_KEY": "value1", "KEEP": "value2"}
	mappings := []rename.Mapping{{OldKey: "OLD_KEY", NewKey: "NEW_KEY"}}

	out, err := rename.Apply(src, mappings, rename.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if out["NEW_KEY"] != "value1" {
		t.Errorf("expected NEW_KEY=value1, got %q", out["NEW_KEY"])
	}
	if out["KEEP"] != "value2" {
		t.Errorf("expected KEEP=value2, got %q", out["KEEP"])
	}
}

func TestApply_DoesNotModifyOriginal(t *testing.T) {
	src := map[string]string{"OLD_KEY": "v"}
	mappings := []rename.Mapping{{OldKey: "OLD_KEY", NewKey: "NEW_KEY"}}

	_, _ = rename.Apply(src, mappings, rename.DefaultOptions())

	if _, ok := src["OLD_KEY"]; !ok {
		t.Error("Apply must not modify the original map")
	}
}

func TestApply_DryRun_NoMutation(t *testing.T) {
	src := map[string]string{"OLD_KEY": "v"}
	mappings := []rename.Mapping{{OldKey: "OLD_KEY", NewKey: "NEW_KEY"}}
	opts := rename.Options{DryRun: true}

	out, err := rename.Apply(src, mappings, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["OLD_KEY"]; !ok {
		t.Error("DryRun: OLD_KEY should still be present in output")
	}
	if _, ok := out["NEW_KEY"]; ok {
		t.Error("DryRun: NEW_KEY should not be present in output")
	}
}

func TestApply_MissingKey_NoError_ByDefault(t *testing.T) {
	src := map[string]string{"OTHER": "v"}
	mappings := []rename.Mapping{{OldKey: "MISSING", NewKey: "NEW_KEY"}}

	_, err := rename.Apply(src, mappings, rename.DefaultOptions())
	if err != nil {
		t.Errorf("expected no error by default, got: %v", err)
	}
}

func TestApply_MissingKey_FailOnMissing(t *testing.T) {
	src := map[string]string{"OTHER": "v"}
	mappings := []rename.Mapping{{OldKey: "MISSING", NewKey: "NEW_KEY"}}
	opts := rename.Options{FailOnMissing: true}

	_, err := rename.Apply(src, mappings, opts)
	if err == nil {
		t.Error("expected error when FailOnMissing=true and key is absent")
	}
}

func TestApply_EmptyKeyReturnsError(t *testing.T) {
	src := map[string]string{"A": "1"}
	mappings := []rename.Mapping{{OldKey: "", NewKey: "B"}}

	_, err := rename.Apply(src, mappings, rename.DefaultOptions())
	if err == nil {
		t.Error("expected error for empty OldKey")
	}
}
