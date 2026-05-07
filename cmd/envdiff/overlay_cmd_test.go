package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeOverlayEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeOverlayEnv: %v", err)
	}
	return p
}

func TestRunOverlay_NoArgs_ReturnsError(t *testing.T) {
	err := runOverlay(nil)
	if err == nil {
		t.Fatal("expected error for no args")
	}
}

func TestRunOverlay_SingleFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	f := writeOverlayEnv(t, dir, "base.env", "APP_ENV=dev\n")
	err := runOverlay([]string{f})
	if err == nil {
		t.Fatal("expected error for single file")
	}
}

func TestRunOverlay_MissingFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	f := writeOverlayEnv(t, dir, "base.env", "APP_ENV=dev\n")
	err := runOverlay([]string{f, "/nonexistent/file.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunOverlay_TwoFiles_MergesCorrectly(t *testing.T) {
	dir := t.TempDir()
	base := writeOverlayEnv(t, dir, "base.env", "APP_ENV=staging\nDB_HOST=localhost\n")
	over := writeOverlayEnv(t, dir, "over.env", "APP_ENV=production\nNEW_KEY=value\n")

	// capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runOverlay([]string{base, over})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var sb strings.Builder
	buf := make([]byte, 4096)
	for {
		n, e := r.Read(buf)
		sb.Write(buf[:n])
		if e != nil {
			break
		}
	}
	out := sb.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got:\n%s", out)
	}
	if !strings.Contains(out, "NEW_KEY=value") {
		t.Errorf("expected NEW_KEY=value in output, got:\n%s", out)
	}
}

func TestRunOverlay_NoOverwrite_BaseWins(t *testing.T) {
	dir := t.TempDir()
	base := writeOverlayEnv(t, dir, "base.env", "APP_ENV=staging\n")
	over := writeOverlayEnv(t, dir, "over.env", "APP_ENV=production\n")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runOverlay([]string{"--no-overwrite", base, over})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var sb strings.Builder
	buf := make([]byte, 4096)
	for {
		n, e := r.Read(buf)
		sb.Write(buf[:n])
		if e != nil {
			break
		}
	}
	if !strings.Contains(sb.String(), "APP_ENV=staging") {
		t.Errorf("expected base value preserved, got:\n%s", sb.String())
	}
}
