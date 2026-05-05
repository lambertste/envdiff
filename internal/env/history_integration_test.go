package env

import (
	"strings"
	"testing"
)

// TestHistory_Integration_MultipleOps exercises a realistic sequence of
// tracked mutations and verifies the full formatted output.
func TestHistory_Integration_MultipleOps(t *testing.T) {
	s := NewSet()
	s.Set("APP_ENV", "development")
	s.Set("LOG_LEVEL", "debug")
	s.Set("SECRET", "plain")

	h := &History{}
	TrackSet(h, s, "APP_ENV", "production")
	TrackSet(h, s, "NEW_VAR", "added")
	TrackDelete(h, s, "SECRET")

	if h.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", h.Len())
	}

	out := h.Format()
	for _, want := range []string{"APP_ENV", "NEW_VAR", "SECRET", "del"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in format output:\n%s", want, out)
		}
	}

	// Verify final Set state.
	if v, _ := s.Get("APP_ENV"); v != "production" {
		t.Errorf("APP_ENV should be production, got %q", v)
	}
	if v, _ := s.Get("NEW_VAR"); v != "added" {
		t.Errorf("NEW_VAR should be added, got %q", v)
	}
	if _, ok := s.Get("SECRET"); ok {
		t.Error("SECRET should have been deleted")
	}
}

// TestHistory_Integration_ReplayOnFreshSet verifies that replaying history
// entries onto a new Set reproduces the same final state.
func TestHistory_Integration_ReplayOnFreshSet(t *testing.T) {
	src := NewSet()
	src.Set("A", "1")
	src.Set("B", "2")

	h := &History{}
	TrackSet(h, src, "A", "10")
	TrackSet(h, src, "C", "3")
	TrackDelete(h, src, "B")

	// Replay onto a new Set.
	dst := NewSet()
	dst.Set("A", "1")
	dst.Set("B", "2")
	for _, e := range h.Entries() {
		switch e.Kind {
		case ChangeSet:
			dst.Set(e.Key, e.NewVal)
		case ChangeDelete:
			dst.Delete(e.Key)
		}
	}

	if v, _ := dst.Get("A"); v != "10" {
		t.Errorf("A: want 10, got %q", v)
	}
	if v, _ := dst.Get("C"); v != "3" {
		t.Errorf("C: want 3, got %q", v)
	}
	if _, ok := dst.Get("B"); ok {
		t.Error("B should be absent after replay")
	}
}
