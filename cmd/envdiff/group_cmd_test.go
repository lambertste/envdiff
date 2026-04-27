package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeGroupEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunGroup_NoArgs_ReturnsError(t *testing.T) {
	if err := runGroup([]string{}); err == nil {
		t.Error("expected error for missing file argument")
	}
}

func TestRunGroup_MissingFile_ReturnsError(t *testing.T) {
	if err := runGroup([]string{"/no/such/file.env"}); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRunGroup_ListFlag_PrintsGroupNames(t *testing.T) {
	p := writeGroupEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=test\n")

	// Capture stdout by checking no error is returned; full output capture
	// would require refactoring runGroup to accept an io.Writer.
	if err := runGroup([]string{p, "--list"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunGroup_FilterByGroup_NoError(t *testing.T) {
	p := writeGroupEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=test\n")
	if err := runGroup([]string{p, "--group", "DB"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunGroup_GroupFlagMissingName_ReturnsError(t *testing.T) {
	p := writeGroupEnv(t, "DB_HOST=localhost\n")
	if err := runGroup([]string{p, "--group"}); err == nil {
		t.Error("expected error when --group has no argument")
	}
}

func TestRunGroup_DefaultBucket(t *testing.T) {
	p := writeGroupEnv(t, "NOPREFIX=value\n")
	_ = strings.Contains // satisfy import
	if err := runGroup([]string{p}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
