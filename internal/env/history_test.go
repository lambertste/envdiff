package env

import (
	"strings"
	"testing"
)

func baseHistorySet() *Set {
	s := NewSet()
	s.Set("HOST", "localhost")
	s.Set("PORT", "5432")
	return s
}

func TestHistory_EmptyFormat(t *testing.T) {
	h := &History{}
	if h.Format() != "(no history)" {
		t.Fatalf("expected empty message, got %q", h.Format())
	}
}

func TestHistory_Len(t *testing.T) {
	h := &History{}
	h.Record(ChangeSet, "KEY", "", "val")
	if h.Len() != 1 {
		t.Fatalf("expected 1, got %d", h.Len())
	}
}

func TestTrackSet_NewKey(t *testing.T) {
	s := baseHistorySet()
	h := &History{}
	TrackSet(h, s, "DB", "postgres")

	if v, _ := s.Get("DB"); v != "postgres" {
		t.Fatalf("expected postgres, got %q", v)
	}
	if h.Len() != 1 {
		t.Fatalf("expected 1 history entry")
	}
	e := h.Entries()[0]
	if e.Kind != ChangeSet || e.OldVal != "" || e.NewVal != "postgres" {
		t.Fatalf("unexpected entry: %+v", e)
	}
}

func TestTrackSet_OverwriteKey(t *testing.T) {
	s := baseHistorySet()
	h := &History{}
	TrackSet(h, s, "PORT", "3306")

	e := h.Entries()[0]
	if e.OldVal != "5432" || e.NewVal != "3306" {
		t.Fatalf("unexpected old/new: %q -> %q", e.OldVal, e.NewVal)
	}
}

func TestTrackDelete(t *testing.T) {
	s := baseHistorySet()
	h := &History{}
	TrackDelete(h, s, "HOST")

	if _, ok := s.Get("HOST"); ok {
		t.Fatal("expected HOST to be deleted")
	}
	e := h.Entries()[0]
	if e.Kind != ChangeDelete || e.OldVal != "localhost" {
		t.Fatalf("unexpected entry: %+v", e)
	}
}

func TestHistory_Format_ContainsKey(t *testing.T) {
	s := baseHistorySet()
	h := &History{}
	TrackSet(h, s, "NEW_KEY", "val")
	TrackDelete(h, s, "PORT")

	out := h.Format()
	if !strings.Contains(out, "NEW_KEY") {
		t.Errorf("expected NEW_KEY in output: %s", out)
	}
	if !strings.Contains(out, "del") {
		t.Errorf("expected del in output: %s", out)
	}
}

func TestHistory_Entries_IsCopy(t *testing.T) {
	h := &History{}
	h.Record(ChangeSet, "A", "", "1")
	copy1 := h.Entries()
	copy1[0].Key = "MUTATED"
	copy2 := h.Entries()
	if copy2[0].Key == "MUTATED" {
		t.Fatal("Entries() should return a copy, not a reference")
	}
}
