package env

import (
	"testing"
)

func baseCompareSet() *Set {
	s := NewSet()
	s.Set("APP_NAME", "myapp")
	s.Set("APP_ENV", "staging")
	s.Set("DB_HOST", "localhost")
	return s
}

func TestCompare_NoChanges(t *testing.T) {
	a := baseCompareSet()
	b := baseCompareSet()
	r := Compare(a, b)
	if r.HasChanges() {
		t.Errorf("expected no changes, got %+v", r)
	}
	if len(r.Unchanged) != 3 {
		t.Errorf("expected 3 unchanged, got %d", len(r.Unchanged))
	}
}

func TestCompare_Added(t *testing.T) {
	a := baseCompareSet()
	b := baseCompareSet()
	b.Set("NEW_KEY", "value")
	r := Compare(a, b)
	if len(r.Added) != 1 || r.Added[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY added, got %v", r.Added)
	}
	if len(r.Removed) != 0 {
		t.Errorf("expected no removals, got %v", r.Removed)
	}
}

func TestCompare_Removed(t *testing.T) {
	a := baseCompareSet()
	b := baseCompareSet()
	b.Delete("DB_HOST")
	r := Compare(a, b)
	if len(r.Removed) != 1 || r.Removed[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST removed, got %v", r.Removed)
	}
	if len(r.Added) != 0 {
		t.Errorf("expected no additions, got %v", r.Added)
	}
}

func TestCompare_Modified(t *testing.T) {
	a := baseCompareSet()
	b := baseCompareSet()
	b.Set("APP_ENV", "production")
	r := Compare(a, b)
	if len(r.Modified) != 1 || r.Modified[0] != "APP_ENV" {
		t.Errorf("expected APP_ENV modified, got %v", r.Modified)
	}
	if len(r.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged, got %d", len(r.Unchanged))
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	a := NewSet()
	b := NewSet()
	for _, k := range []string{"Z_KEY", "A_KEY", "M_KEY"} {
		b.Set(k, "v")
	}
	r := Compare(a, b)
	if len(r.Added) != 3 {
		t.Fatalf("expected 3 added, got %d", len(r.Added))
	}
	if r.Added[0] != "A_KEY" || r.Added[1] != "M_KEY" || r.Added[2] != "Z_KEY" {
		t.Errorf("expected sorted added keys, got %v", r.Added)
	}
}

func TestCompare_Summary_NoChanges(t *testing.T) {
	a := baseCompareSet()
	b := baseCompareSet()
	r := Compare(a, b)
	if r.Summary() != "no changes" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}
