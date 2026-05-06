package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeRollbackEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunRollback_NoArgs_ReturnsError(t *testing.T) {
	err := runRollback([]string{"--keys", "PORT"})
	if err == nil {
		t.Fatal("expected error for missing file args")
	}
}

func TestRunRollback_NoKeys_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	before := writeRollbackEnv(t, dir, "before.env", "PORT=5432\n")
	current := writeRollbackEnv(t, dir, "current.env", "PORT=9999\n")
	err := runRollback([]string{before, current})
	if err == nil || !strings.Contains(err.Error(), "--keys") {
		t.Fatalf("expected --keys error, got %v", err)
	}
}

func TestRunRollback_MissingFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	before := writeRollbackEnv(t, dir, "before.env", "PORT=5432\n")
	err := runRollback([]string{"--keys", "PORT", before, "/no/such/file.env"})
	if err == nil {
		t.Fatal("expected error for missing current file")
	}
}

func TestRunRollback_DryRun_PrintsPlan(t *testing.T) {
	dir := t.TempDir()
	before := writeRollbackEnv(t, dir, "before.env", "PORT=5432\nHOST=localhost\n")
	current := writeRollbackEnv(t, dir, "current.env", "PORT=9999\nHOST=localhost\n")
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := runRollback([]string{"--keys", "PORT", "--dry-run", before, current})
	w.Close()
	os.Stdout = old
	var buf strings.Builder
	buf.ReadFrom(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "PORT") {
		t.Errorf("expected PORT in dry-run output, got: %s", buf.String())
	}
}

func TestRunRollback_AppliesRollback_WritesFile(t *testing.T) {
	dir := t.TempDir()
	before := writeRollbackEnv(t, dir, "before.env", "PORT=5432\nHOST=localhost\n")
	current := writeRollbackEnv(t, dir, "current.env", "PORT=9999\nHOST=localhost\n")
	out := filepath.Join(dir, "out.env")
	err := runRollback([]string{"--keys", "PORT", "--out", out, before, current})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "PORT=5432") {
		t.Errorf("expected PORT=5432 in output, got: %s", string(data))
	}
}
