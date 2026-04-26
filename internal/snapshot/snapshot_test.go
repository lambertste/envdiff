package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envdiff/internal/snapshot"
)

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	entries := map[string]string{
		"APP_ENV": "production",
		"DB_HOST": "localhost",
	}

	if err := snapshot.Save(path, "test-label", entries); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	s, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if s.Label != "test-label" {
		t.Errorf("expected label 'test-label', got %q", s.Label)
	}
	if s.Entries["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", s.Entries["APP_ENV"])
	}
	if s.Entries["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", s.Entries["DB_HOST"])
	}
	if s.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0644)

	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestToEntries(t *testing.T) {
	s := &snapshot.Snapshot{
		Timestamp: time.Now(),
		Label:     "x",
		Entries:   map[string]string{"KEY": "val"},
	}
	lines := s.ToEntries()
	if len(lines) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(lines))
	}
	if lines[0] != "KEY=val" {
		t.Errorf("unexpected entry: %q", lines[0])
	}
}
