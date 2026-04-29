package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeDedupeEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeDedupeEnv: %v", err)
	}
	return p
}

func TestRunDedupe_NoArgs_ReturnsError(t *testing.T) {
	if err := runDedupe([]string{}); err == nil {
		t.Fatal("expected error for missing file argument")
	}
}

func TestRunDedupe_MissingFile_ReturnsError(t *testing.T) {
	err := runDedupe([]string{"/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunDedupe_NoDuplicates_PrintsAllKeys(t *testing.T) {
	p := writeDedupeEnv(t, "ALPHA=1\nBETA=2\nGAMMA=3\n")
	if err := runDedupe([]string{p}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunDedupe_ListDupes_NoDuplicates(t *testing.T) {
	p := writeDedupeEnv(t, "ALPHA=1\nBETA=2\n")
	// Redirect stdout capture via running function — just ensure no error.
	if err := runDedupe([]string{"--list-dupes", p}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunDedupe_StrategyFirst_Flag(t *testing.T) {
	p := writeDedupeEnv(t, "KEY=first\nOTHER=x\nKEY=second\n")
	if err := runDedupe([]string{"--strategy=first", p}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunDedupe_StrategyLast_Flag(t *testing.T) {
	p := writeDedupeEnv(t, "KEY=first\nOTHER=x\nKEY=second\n")
	if err := runDedupe([]string{"--strategy=last", p}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunDedupe_OnlyFlagsNoFile_ReturnsError(t *testing.T) {
	err := runDedupe([]string{"--strategy=first", "--list-dupes"})
	if err == nil {
		t.Fatal("expected error when no file path given")
	}
	if !strings.Contains(err.Error(), "no input file") {
		t.Errorf("unexpected error message: %v", err)
	}
}
