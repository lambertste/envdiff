package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func captureExport(t *testing.T, args []string) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runExport(args)

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	b := make([]byte, 4096)
	for {
		n, e := r.Read(b)
		if n > 0 {
			buf.Write(b[:n])
		}
		if e != nil {
			break
		}
	}
	if err != nil {
		t.Fatalf("runExport error: %v", err)
	}
	return buf.String()
}

func TestRunExport_DotenvDefault(t *testing.T) {
	p := writeEnvFile(t, "FOO=bar\nBAZ=qux\n")
	out := captureExport(t, []string{p})
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got: %q", out)
	}
}

func TestRunExport_JSONFormat(t *testing.T) {
	p := writeEnvFile(t, "KEY=val\n")
	out := captureExport(t, []string{"--format", "json", p})
	if !strings.Contains(out, `"KEY"`) {
		t.Errorf("expected JSON output, got: %q", out)
	}
}

func TestRunExport_ExportFormat(t *testing.T) {
	p := writeEnvFile(t, "PORT=9000\n")
	out := captureExport(t, []string{"-f", "export", p})
	if !strings.HasPrefix(out, "export PORT=") {
		t.Errorf("expected export prefix, got: %q", out)
	}
}

func TestRunExport_MissingFile(t *testing.T) {
	err := runExport([]string{"/nonexistent/.env"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRunExport_NoArgs(t *testing.T) {
	err := runExport([]string{})
	if err == nil {
		t.Error("expected usage error")
	}
}
