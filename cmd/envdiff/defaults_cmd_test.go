package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeDefaultsEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunDefaults_NoFile_ReturnsError(t *testing.T) {
	err := runDefaults([]string{})
	if err == nil || !strings.Contains(err.Error(), "-file is required") {
		t.Fatalf("expected -file required error, got %v", err)
	}
}

func TestRunDefaults_MissingFile_ReturnsError(t *testing.T) {
	err := runDefaults([]string{"-file", "/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunDefaults_FillsMissingKey(t *testing.T) {
	p := writeDefaultsEnv(t, "APP_NAME=myapp\n")
	// capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runDefaults([]string{"-file", p, "-set", "PORT=9090"})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	buf := new(strings.Builder)
	data := make([]byte, 512)
	for {
		n, e := r.Read(data)
		buf.Write(data[:n])
		if e != nil {
			break
		}
	}
	if !strings.Contains(buf.String(), "PORT=9090") {
		t.Errorf("expected PORT=9090 in output, got: %s", buf.String())
	}
}

func TestRunDefaults_MissingFlag_ListsAbsent(t *testing.T) {
	p := writeDefaultsEnv(t, "APP_NAME=myapp\n")
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runDefaults([]string{"-file", p, "-set", "MISSING_KEY=val", "-missing"})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	buf := new(strings.Builder)
	data := make([]byte, 512)
	for {
		n, e := r.Read(data)
		buf.Write(data[:n])
		if e != nil {
			break
		}
	}
	if !strings.Contains(buf.String(), "MISSING_KEY") {
		t.Errorf("expected MISSING_KEY in output, got: %s", buf.String())
	}
}

func TestParseDefaultSpecs_InvalidFormat(t *testing.T) {
	_, err := parseDefaultSpecs([]string{"NOEQUALS"}, false)
	if err == nil || !strings.Contains(err.Error(), "KEY=VALUE") {
		t.Fatalf("expected KEY=VALUE error, got %v", err)
	}
}
