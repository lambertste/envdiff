package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func withProfileDir(t *testing.T) func() {
	t.Helper()
	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	return func() { os.Chdir(origDir) }
}

func TestRunProfileAdd_CreatesProfile(t *testing.T) {
	defer withProfileDir(t)()

	if err := runProfileAdd("staging", ".env.staging", []string{"ci"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(".envdiff", "profiles.json")); err != nil {
		t.Error("expected profiles.json to be created")
	}
}

func TestRunProfileRemove_ExistingProfile(t *testing.T) {
	defer withProfileDir(t)()

	runProfileAdd("prod", ".env.prod", nil)
	if err := runProfileRemove("prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunProfileRemove_NotFound(t *testing.T) {
	defer withProfileDir(t)()

	err := runProfileRemove("ghost")
	if err == nil {
		t.Error("expected error for missing profile")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunProfileList_Empty(t *testing.T) {
	defer withProfileDir(t)()

	if err := runProfileList(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunProfileList_WithEntries(t *testing.T) {
	defer withProfileDir(t)()

	runProfileAdd("dev", ".env.dev", []string{"local"})
	runProfileAdd("prod", ".env.prod", nil)

	if err := runProfileList(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunProfileShow_Found(t *testing.T) {
	defer withProfileDir(t)()

	runProfileAdd("dev", ".env.dev", []string{"local"})
	if err := runProfileShow("dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunProfileShow_NotFound(t *testing.T) {
	defer withProfileDir(t)()

	err := runProfileShow("missing")
	if err == nil {
		t.Error("expected error")
	}
}
