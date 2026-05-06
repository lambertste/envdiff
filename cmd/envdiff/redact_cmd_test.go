package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeRedactEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunRedact_NoArgs_ReturnsError(t *testing.T) {
	if err := runRedact([]string{}); err == nil {
		t.Error("expected error for missing file argument")
	}
}

func TestRunRedact_MissingFile_ReturnsError(t *testing.T) {
	err := runRedact([]string{"/nonexistent/.env"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRunRedact_DefaultOptions_MasksSensitiveKeys(t *testing.T) {
	p := writeRedactEnv(t, "APP_NAME=myapp\nDB_PASSWORD=s3cr3t\nPORT=8080\n")

	// Capture stdout by redirecting to a temp file.
	tmp := filepath.Join(t.TempDir(), "out.env")
	err := runRedact([]string{"-out", tmp, p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	output := string(data)

	if strings.Contains(output, "s3cr3t") {
		t.Error("expected DB_PASSWORD to be redacted")
	}
	if !strings.Contains(output, "APP_NAME") {
		t.Error("expected APP_NAME to be present")
	}
}

func TestRunRedact_ExplicitKeys_RedactsOnlyThose(t *testing.T) {
	p := writeRedactEnv(t, "APP_NAME=myapp\nPORT=8080\n")
	tmp := filepath.Join(t.TempDir(), "out.env")

	err := runRedact([]string{"-keys", "APP_NAME", "-placeholder", "HIDDEN", "-out", tmp, p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	output := string(data)

	if strings.Contains(output, "myapp") {
		t.Error("expected APP_NAME value to be redacted")
	}
	if !strings.Contains(output, "8080") {
		t.Error("expected PORT to be preserved")
	}
}

func TestRunRedact_ListFlag_PrintsKeys(t *testing.T) {
	p := writeRedactEnv(t, "API_KEY=abc\nHOST=localhost\n")
	// -list writes to stdout; just ensure no error and it runs.
	err := runRedact([]string{"-list", p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunRedact_JSONFormat_NoError(t *testing.T) {
	p := writeRedactEnv(t, "SECRET_KEY=abc\nAPP=myapp\n")
	tmp := filepath.Join(t.TempDir(), "out.json")

	err := runRedact([]string{"-format", "json", "-out", tmp, p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	if len(data) == 0 {
		t.Error("expected non-empty JSON output")
	}
}
