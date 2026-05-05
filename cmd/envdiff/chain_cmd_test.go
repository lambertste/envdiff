package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeChainEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunChain_NoArgs_ReturnsError(t *testing.T) {
	err := runChain([]string{})
	if err == nil {
		t.Fatal("expected error for no args")
	}
}

func TestRunChain_MissingFile_ReturnsError(t *testing.T) {
	err := runChain([]string{"/nonexistent/path.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunChain_FirstWins_Default(t *testing.T) {
	dir := t.TempDir()
	a := writeChainEnv(t, dir, "a.env", "APP_PORT=8080\nAPP_HOST=localhost\n")
	b := writeChainEnv(t, dir, "b.env", "APP_PORT=9090\nAPP_DEBUG=true\n")

	// Capture stdout via pipe
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runChain([]string{a, b})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf strings.Builder
	buf.ReadFrom(r)
	out := buf.String()

	if !strings.Contains(out, "APP_PORT=8080") {
		t.Errorf("expected APP_PORT=8080 (first wins), got:\n%s", out)
	}
}

func TestRunChain_OverwriteFlag_LastWins(t *testing.T) {
	dir := t.TempDir()
	a := writeChainEnv(t, dir, "a.env", "APP_PORT=8080\n")
	b := writeChainEnv(t, dir, "b.env", "APP_PORT=9090\n")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runChain([]string{"--overwrite", a, b})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf strings.Builder
	buf.ReadFrom(r)
	out := buf.String()

	if !strings.Contains(out, "APP_PORT=9090") {
		t.Errorf("expected APP_PORT=9090 (overwrite), got:\n%s", out)
	}
}
