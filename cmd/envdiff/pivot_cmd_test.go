package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writePivotEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writePivotEnv: %v", err)
	}
	return p
}

func TestRunPivot_NoKey_ReturnsError(t *testing.T) {
	err := runPivot("", []string{"a.env", "b.env"})
	if err == nil || !strings.Contains(err.Error(), "--key") {
		t.Fatalf("expected --key error, got %v", err)
	}
}

func TestRunPivot_TooFewFiles_ReturnsError(t *testing.T) {
	err := runPivot("ENV", []string{"a.env"})
	if err == nil || !strings.Contains(err.Error(), "two env files") {
		t.Fatalf("expected file count error, got %v", err)
	}
}

func TestRunPivot_MissingFile_ReturnsError(t *testing.T) {
	err := runPivot("ENV", []string{"/no/such/a.env", "/no/such/b.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunPivot_ValidFiles_NoError(t *testing.T) {
	dir := t.TempDir()
	a := writePivotEnv(t, dir, "staging.env", "ENV=staging\nDB_HOST=db1\nAPI_KEY=abc\n")
	b := writePivotEnv(t, dir, "prod.env", "ENV=prod\nDB_HOST=db2\nAPI_KEY=xyz\n")

	if err := runPivot("ENV", []string{a, b}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunPivot_MissingPivotKeyInFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	a := writePivotEnv(t, dir, "a.env", "DB_HOST=db1\n")
	b := writePivotEnv(t, dir, "b.env", "DB_HOST=db2\n")

	err := runPivot("ENV", []string{a, b})
	if err == nil {
		t.Fatal("expected error when pivot key absent from file")
	}
}
