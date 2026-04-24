package diff

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
)

func TestDiff_NoChanges(t *testing.T) {
	left := parser.EnvMap{"KEY": "value"}
	right := parser.EnvMap{"KEY": "value"}
	r := Diff(left, right)
	if r.HasChanges() {
		t.Errorf("expected no changes, got %+v", r.Changes)
	}
}

func TestDiff_Added(t *testing.T) {
	left := parser.EnvMap{}
	right := parser.EnvMap{"NEW_KEY": "hello"}
	r := Diff(left, right)
	if len(r.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(r.Changes))
	}
	if r.Changes[0].Kind != Added || r.Changes[0].Key != "NEW_KEY" {
		t.Errorf("unexpected change: %+v", r.Changes[0])
	}
}

func TestDiff_Removed(t *testing.T) {
	left := parser.EnvMap{"OLD_KEY": "bye"}
	right := parser.EnvMap{}
	r := Diff(left, right)
	if len(r.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(r.Changes))
	}
	if r.Changes[0].Kind != Removed || r.Changes[0].OldValue != "bye" {
		t.Errorf("unexpected change: %+v", r.Changes[0])
	}
}

func TestDiff_Modified(t *testing.T) {
	left := parser.EnvMap{"HOST": "staging.example.com"}
	right := parser.EnvMap{"HOST": "prod.example.com"}
	r := Diff(left, right)
	if len(r.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(r.Changes))
	}
	c := r.Changes[0]
	if c.Kind != Modified || c.OldValue != "staging.example.com" || c.NewValue != "prod.example.com" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiff_SortedOutput(t *testing.T) {
	left := parser.EnvMap{"Z_KEY": "1", "A_KEY": "1"}
	right := parser.EnvMap{"Z_KEY": "2", "A_KEY": "2"}
	r := Diff(left, right)
	if len(r.Changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(r.Changes))
	}
	if r.Changes[0].Key != "A_KEY" || r.Changes[1].Key != "Z_KEY" {
		t.Errorf("changes not sorted: %+v", r.Changes)
	}
}
