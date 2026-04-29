package env

import (
	"testing"
)

func baseDiffSet(pairs ...string) *Set {
	s := NewSet()
	for i := 0; i+1 < len(pairs); i += 2 {
		s.Set(pairs[i], pairs[i+1])
	}
	return s
}

func TestDiffSets_NoChanges(t *testing.T) {
	a := baseDiffSet("A", "1", "B", "2")
	b := baseDiffSet("A", "1", "B", "2")
	entries := DiffSets(a, b)
	for _, e := range entries {
		if e.Kind != DiffUnchanged {
			t.Errorf("expected unchanged for key %s, got %s", e.Key, e.Kind)
		}
	}
}

func TestDiffSets_Added(t *testing.T) {
	a := baseDiffSet("A", "1")
	b := baseDiffSet("A", "1", "B", "2")
	entries := DiffSets(a, b)
	found := FilterDiff(entries, DiffAdded)
	if len(found) != 1 || found[0].Key != "B" {
		t.Errorf("expected added key B, got %+v", found)
	}
}

func TestDiffSets_Removed(t *testing.T) {
	a := baseDiffSet("A", "1", "B", "2")
	b := baseDiffSet("A", "1")
	entries := DiffSets(a, b)
	found := FilterDiff(entries, DiffRemoved)
	if len(found) != 1 || found[0].Key != "B" {
		t.Errorf("expected removed key B, got %+v", found)
	}
	if found[0].OldValue != "2" {
		t.Errorf("expected OldValue=2, got %s", found[0].OldValue)
	}
}

func TestDiffSets_Modified(t *testing.T) {
	a := baseDiffSet("A", "old")
	b := baseDiffSet("A", "new")
	entries := DiffSets(a, b)
	found := FilterDiff(entries, DiffModified)
	if len(found) != 1 || found[0].Key != "A" {
		t.Errorf("expected modified key A, got %+v", found)
	}
	if found[0].OldValue != "old" || found[0].NewValue != "new" {
		t.Errorf("unexpected values: %+v", found[0])
	}
}

func TestFilterDiff_MultipleKinds(t *testing.T) {
	a := baseDiffSet("A", "1", "B", "2")
	b := baseDiffSet("A", "changed", "C", "3")
	entries := DiffSets(a, b)
	active := FilterDiff(entries, DiffAdded, DiffRemoved, DiffModified)
	if len(active) != 3 {
		t.Errorf("expected 3 active changes, got %d: %+v", len(active), active)
	}
}
