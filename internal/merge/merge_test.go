package merge

import (
	"strings"
	"testing"
)

func TestMerge_NoConflicts(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	over := map[string]string{"C": "3"}

	res := Merge(base, over, PreferBase)

	if len(res.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %d", len(res.Conflicts))
	}
	if res.Merged["A"] != "1" || res.Merged["B"] != "2" || res.Merged["C"] != "3" {
		t.Errorf("unexpected merged map: %v", res.Merged)
	}
}

func TestMerge_PreferBase(t *testing.T) {
	base := map[string]string{"KEY": "base_val"}
	over := map[string]string{"KEY": "over_val"}

	res := Merge(base, over, PreferBase)

	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	if res.Merged["KEY"] != "base_val" {
		t.Errorf("expected base_val, got %q", res.Merged["KEY"])
	}
	if res.Conflicts[0].Resolved != "base_val" {
		t.Errorf("conflict resolved to wrong value: %q", res.Conflicts[0].Resolved)
	}
}

func TestMerge_PreferOverride(t *testing.T) {
	base := map[string]string{"KEY": "base_val"}
	over := map[string]string{"KEY": "over_val"}

	res := Merge(base, over, PreferOverride)

	if res.Merged["KEY"] != "over_val" {
		t.Errorf("expected over_val, got %q", res.Merged["KEY"])
	}
	if res.Conflicts[0].Resolved != "over_val" {
		t.Errorf("conflict resolved to wrong value: %q", res.Conflicts[0].Resolved)
	}
}

func TestMerge_SameValueNoConflict(t *testing.T) {
	base := map[string]string{"KEY": "same"}
	over := map[string]string{"KEY": "same"}

	res := Merge(base, over, PreferBase)

	if len(res.Conflicts) != 0 {
		t.Errorf("identical values should not produce conflicts")
	}
}

func TestMerge_ConflictsSorted(t *testing.T) {
	base := map[string]string{"Z": "1", "A": "1", "M": "1"}
	over := map[string]string{"Z": "2", "A": "2", "M": "2"}

	res := Merge(base, over, PreferBase)

	if len(res.Conflicts) != 3 {
		t.Fatalf("expected 3 conflicts, got %d", len(res.Conflicts))
	}
	if res.Conflicts[0].Key != "A" || res.Conflicts[1].Key != "M" || res.Conflicts[2].Key != "Z" {
		t.Errorf("conflicts not sorted: %v", res.Conflicts)
	}
}

func TestFormatConflicts_Empty(t *testing.T) {
	out := FormatConflicts(nil)
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestFormatConflicts_NonEmpty(t *testing.T) {
	conflicts := []Conflict{
		{Key: "DB_PASS", BaseValue: "secret", OverValue: "other", Resolved: "secret"},
	}
	out := FormatConflicts(conflicts)
	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected key in output, got: %s", out)
	}
	if !strings.Contains(out, "1 conflict") {
		t.Errorf("expected conflict count in output, got: %s", out)
	}
}

func TestToEntries_Sorted(t *testing.T) {
	merged := map[string]string{"Z": "z", "A": "a", "M": "m"}
	entries := ToEntries(merged)

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Key != "A" || entries[1].Key != "M" || entries[2].Key != "Z" {
		t.Errorf("entries not sorted: %v", entries)
	}
}
