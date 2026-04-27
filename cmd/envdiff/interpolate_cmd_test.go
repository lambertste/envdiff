package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeInterpEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunInterpolate_BasicExpansion(t *testing.T) {
	p := writeInterpEnv(t, "BASE=/opt\nBIN=${BASE}/bin\n")
	err := runInterpolate([]string{p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunInterpolate_NoFile_ReturnsError(t *testing.T) {
	err := runInterpolate([]string{})
	if err == nil {
		t.Fatal("expected error for missing file argument")
	}
}

func TestRunInterpolate_MissingFile_ReturnsError(t *testing.T) {
	err := runInterpolate([]string{"/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestRunInterpolate_StrictMode_FailsOnMissing(t *testing.T) {
	p := writeInterpEnv(t, "DIR=${UNDEFINED}/bin\n")
	err := runInterpolate([]string{"-strict", p})
	if err == nil {
		t.Fatal("expected strict mode to return error")
	}
	if !strings.Contains(err.Error(), "unresolved") {
		t.Errorf("error should mention 'unresolved', got: %v", err)
	}
}

func TestRunInterpolate_ShellFormat(t *testing.T) {
	p := writeInterpEnv(t, "GREETING=hello\n")
	err := runInterpolate([]string{"-format", "shell", p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
