package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFreezeEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write env: %v", err)
	}
	return p
}

func TestRunFreeze_NoArgs_ReturnsError(t *testing.T) {
	if err := runFreeze([]string{}); err == nil {
		t.Error("expected error for missing file argument")
	}
}

func TestRunFreeze_MissingFile_ReturnsError(t *testing.T) {
	err := runFreeze([]string{"/nonexistent/.env", "--keys=HOST"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRunFreeze_FreezeKeys_NoError(t *testing.T) {
	p := writeFreezeEnv(t, "HOST=localhost\nPORT=8080\n")
	err := runFreeze([]string{p, "--keys=HOST,PORT"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunFreeze_MissingKeysFlag_ReturnsError(t *testing.T) {
	p := writeFreezeEnv(t, "HOST=localhost\n")
	err := runFreeze([]string{p})
	if err == nil {
		t.Error("expected error when --keys not provided")
	}
}

func TestRunFreeze_ListFlag_NoError(t *testing.T) {
	p := writeFreezeEnv(t, "HOST=localhost\n")
	err := runFreeze([]string{p, "--list"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunFreeze_UnfreezeFlag_NoError(t *testing.T) {
	p := writeFreezeEnv(t, "HOST=localhost\n")
	err := runFreeze([]string{p, "--keys=HOST", "--unfreeze"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
