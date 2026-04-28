package env

import (
	"testing"
)

func baseSortSet() *Set {
	s := NewSet()
	s.Set("ZEBRA", "last")
	s.Set("APPLE", "first")
	s.Set("MANGO", "middle")
	s.Set("FIG", "short_val")
	return s
}

func TestSortedKeys_Alpha(t *testing.T) {
	s := baseSortSet()
	keys := SortedKeys(s, SortAlpha)
	want := []string{"APPLE", "FIG", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("index %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSortedKeys_AlphaDesc(t *testing.T) {
	s := baseSortSet()
	keys := SortedKeys(s, SortAlphaDesc)
	want := []string{"ZEBRA", "MANGO", "FIG", "APPLE"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("index %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSortedKeys_ByValue(t *testing.T) {
	s := baseSortSet()
	keys := SortedKeys(s, SortByValue)
	// values: first, last, middle, short_val => APPLE, ZEBRA, MANGO, FIG
	want := []string{"APPLE", "ZEBRA", "MANGO", "FIG"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("index %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSortedKeys_ByLength(t *testing.T) {
	s := baseSortSet()
	keys := SortedKeys(s, SortByLength)
	// lengths: FIG=3, APPLE=5, MANGO=5, ZEBRA=5
	if keys[0] != "FIG" {
		t.Errorf("expected FIG first by length, got %q", keys[0])
	}
}

func TestSortedEntries_ReturnsCorrectPairs(t *testing.T) {
	s := NewSet()
	s.Set("B", "two")
	s.Set("A", "one")
	entries := SortedEntries(s, SortAlpha)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0][0] != "A" || entries[0][1] != "one" {
		t.Errorf("unexpected first entry: %v", entries[0])
	}
	if entries[1][0] != "B" || entries[1][1] != "two" {
		t.Errorf("unexpected second entry: %v", entries[1])
	}
}

func TestSortedKeys_EmptySet(t *testing.T) {
	s := NewSet()
	keys := SortedKeys(s, SortAlpha)
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %v", keys)
	}
}
