package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeHistoryEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeHistoryEnv: %v", err)
	}
	return p
}

func TestRunHistory_NoArgs_ReturnsError(t *testing.T) {
	if err := runHistory([]string{}); err == nil {
		t.Fatal("expected error for missing file argument")
	}
}

func TestRunHistory_MissingFile_ReturnsError(t *testing.T) {
	err := runHistory([]string{"/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunHistory_NoOps_PrintsNoHistory(t *testing.T) {
	p := writeHistoryEnv(t, "HOST=localhost\nPORT=5432\n")
	// Capture stdout via pipe.
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runHistory([]string{p})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	buf := make([]byte, 256)
	n, _ := r.Read(buf)
	if !strings.Contains(string(buf[:n]), "no history") {
		t.Errorf("expected 'no history', got %q", string(buf[:n]))
	}
}

func TestRunHistory_SetOp_PrintsEntry(t *testing.T) {
	p := writeHistoryEnv(t, "HOST=localhost\n")
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runHistory([]string{p, "--set", "HOST=production"})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	buf := make([]byte, 512)
	n, _ := r.Read(buf)
	out := string(buf[:n])
	if !strings.Contains(out, "HOST") {
		t.Errorf("expected HOST in output: %s", out)
	}
}

func TestRunHistory_InvalidSet_ReturnsError(t *testing.T) {
	p := writeHistoryEnv(t, "A=1\n")
	err := runHistory([]string{p, "--set", "NOKEYVALUE"})
	if err == nil {
		t.Fatal("expected error for malformed --set value")
	}
}
