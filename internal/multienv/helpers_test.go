package multienv_test

import (
	"testing"

	"github.com/user/envdiff/internal/multienv"
)

func TestEnvSet_Names_Sorted(t *testing.T) {
	p1 := writeTempEnv(t, ".env", "A=1\n")
	p2 := writeTempEnv(t, ".env", "B=2\n")
	p3 := writeTempEnv(t, ".env", "C=3\n")

	set, err := multienv.Load(
		[]string{"prod", "dev", "staging"},
		[]string{p1, p2, p3},
		multienv.LoadOptions{SkipExtensionCheck: true},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	names := set.Names()
	expected := []string{"dev", "prod", "staging"}
	for i, n := range names {
		if n != expected[i] {
			t.Errorf("names[%d]: got %q, want %q", i, n, expected[i])
		}
	}
}

func TestEnvSet_AllKeys_Union(t *testing.T) {
	p1 := writeTempEnv(t, ".env", "FOO=1\nSHARED=x\n")
	p2 := writeTempEnv(t, ".env", "BAR=2\nSHARED=y\n")

	set, err := multienv.Load(
		[]string{"a", "b"},
		[]string{p1, p2},
		multienv.LoadOptions{SkipExtensionCheck: true},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	keys := set.AllKeys()
	expected := []string{"BAR", "FOO", "SHARED"}
	if len(keys) != len(expected) {
		t.Fatalf("AllKeys: got %v, want %v", keys, expected)
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("keys[%d]: got %q, want %q", i, k, expected[i])
		}
	}
}
