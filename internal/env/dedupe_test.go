package env

import (
	"testing"
)

// buildSetWithDupes creates a Set and forcibly inserts duplicate keys by
// calling Set multiple times (the Set type keeps last value; we test that
// Dedupe can distinguish strategies when a caller pre-loads values).
func buildSetWithDupes() *Set {
	s := NewSet()
	s.Set("ALPHA", "first")
	s.Set("BETA", "one")
	s.Set("ALPHA", "second") // duplicate — KeepLast should keep "second"
	s.Set("GAMMA", "g")
	s.Set("BETA", "two") // duplicate — KeepLast should keep "two"
	return s
}

func TestDedupe_KeepLast_NoDuplicatesReported(t *testing.T) {
	s := buildSetWithDupes()
	res := Dedupe(s, DedupeKeepLast)
	// After deduplication the result set should have exactly 3 keys.
	if len(res.Set.Keys()) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(res.Set.Keys()))
	}
}

func TestDedupe_KeepLast_ValuesAreLatest(t *testing.T) {
	s := buildSetWithDupes()
	res := Dedupe(s, DedupeKeepLast)
	v, _ := res.Set.Get("ALPHA")
	if v != "second" {
		t.Errorf("expected ALPHA=second, got %q", v)
	}
	v, _ = res.Set.Get("BETA")
	if v != "two" {
		t.Errorf("expected BETA=two, got %q", v)
	}
}

func TestDedupe_KeepFirst_ValuesAreEarliest(t *testing.T) {
	s := buildSetWithDupes()
	res := Dedupe(s, DedupeKeepFirst)
	v, _ := res.Set.Get("ALPHA")
	if v != "first" {
		t.Errorf("expected ALPHA=first, got %q", v)
	}
	v, _ = res.Set.Get("BETA")
	if v != "one" {
		t.Errorf("expected BETA=one, got %q", v)
	}
}

func TestDedupe_NoDuplicates_EmptyDuplicatesList(t *testing.T) {
	s := NewSet()
	s.Set("A", "1")
	s.Set("B", "2")
	res := Dedupe(s, DedupeKeepFirst)
	if len(res.Duplicates) != 0 {
		t.Errorf("expected no duplicates, got %v", res.Duplicates)
	}
}

func TestDedupe_DuplicatesSorted(t *testing.T) {
	s := NewSet()
	s.Set("ZEBRA", "z1")
	s.Set("APPLE", "a1")
	s.Set("ZEBRA", "z2")
	s.Set("APPLE", "a2")
	res := Dedupe(s, DedupeKeepLast)
	if len(res.Duplicates) != 2 {
		t.Fatalf("expected 2 duplicates, got %d", len(res.Duplicates))
	}
	if res.Duplicates[0] != "APPLE" || res.Duplicates[1] != "ZEBRA" {
		t.Errorf("expected sorted duplicates [APPLE ZEBRA], got %v", res.Duplicates)
	}
}

func TestDedupe_EmptySet(t *testing.T) {
	s := NewSet()
	res := Dedupe(s, DedupeKeepFirst)
	if len(res.Set.Keys()) != 0 {
		t.Errorf("expected empty result set")
	}
	if len(res.Duplicates) != 0 {
		t.Errorf("expected no duplicates")
	}
}
