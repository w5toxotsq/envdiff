package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// buildBinary compiles the envdiff binary into a temp dir and returns its path.
func buildBinary(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "envdiff")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, out)
	}
	return binPath
}

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	return p
}

func TestMain_NoDifferences(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()
	a := writeFile(t, dir, ".env.a", "FOO=bar\nBAZ=qux\n")
	b := writeFile(t, dir, ".env.b", "FOO=bar\nBAZ=qux\n")

	cmd := exec.Command(bin, a, b)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v\noutput: %s", err, out)
	}
}

func TestMain_Differences_ExitCode2(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()
	a := writeFile(t, dir, ".env.a", "FOO=bar\n")
	b := writeFile(t, dir, ".env.b", "FOO=bar\nEXTRA=value\n")

	cmd := exec.Command(bin, a, b)
	out, err := cmd.CombinedOutput()
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 2 {
			t.Fatalf("expected exit code 2, got %d\noutput: %s", exitErr.ExitCode(), out)
		}
	} else if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	} else {
		t.Fatalf("expected exit code 2 but got 0\noutput: %s", out)
	}
}

func TestMain_JSONFormat(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()
	a := writeFile(t, dir, ".env.a", "FOO=bar\n")
	b := writeFile(t, dir, ".env.b", "FOO=different\n")

	cmd := exec.Command(bin, "-format=json", a, b)
	out, _ := cmd.CombinedOutput()
	if len(out) == 0 {
		t.Fatal("expected non-empty JSON output")
	}
	if out[0] != '{' {
		t.Fatalf("expected JSON object, got: %s", out)
	}
}

func TestMain_MissingArgs(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for missing args")
	}
}
