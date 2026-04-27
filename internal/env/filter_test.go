package env

import "testing"

func baseSet() *Set {
	s := NewSet()
	s.Set("APP_HOST", "localhost")
	s.Set("APP_PORT", "8080")
	s.Set("DB_HOST", "db.internal")
	s.Set("DB_PASS", "")
	s.Set("LOG_LEVEL", "info")
	return s
}

func TestFilter_WithPrefix(t *testing.T) {
	out := Filter(baseSet(), WithPrefix("APP_"))
	if out.Len() != 2 {
		t.Fatalf("expected 2, got %d", out.Len())
	}
	if _, ok := out.Get("APP_HOST"); !ok {
		t.Error("expected APP_HOST")
	}
	if _, ok := out.Get("DB_HOST"); ok {
		t.Error("did not expect DB_HOST")
	}
}

func TestFilter_WithSuffix(t *testing.T) {
	out := Filter(baseSet(), WithSuffix("_HOST"))
	if out.Len() != 2 {
		t.Fatalf("expected 2, got %d", out.Len())
	}
}

func TestFilter_NonEmpty(t *testing.T) {
	out := Filter(baseSet(), NonEmpty())
	if _, ok := out.Get("DB_PASS"); ok {
		t.Error("DB_PASS should be excluded (empty value)")
	}
	if out.Len() != 4 {
		t.Fatalf("expected 4, got %d", out.Len())
	}
}

func TestFilter_ExcludeKeys(t *testing.T) {
	out := Filter(baseSet(), ExcludeKeys("APP_PORT", "LOG_LEVEL"))
	if out.Len() != 3 {
		t.Fatalf("expected 3, got %d", out.Len())
	}
	if _, ok := out.Get("APP_PORT"); ok {
		t.Error("APP_PORT should be excluded")
	}
}

func TestFilter_CombinedPredicates(t *testing.T) {
	out := Filter(baseSet(), WithPrefix("DB_"), NonEmpty())
	// DB_HOST passes, DB_PASS fails NonEmpty
	if out.Len() != 1 {
		t.Fatalf("expected 1, got %d", out.Len())
	}
	if _, ok := out.Get("DB_HOST"); !ok {
		t.Error("expected DB_HOST")
	}
}

func TestFilter_NoPredicates_ReturnsAll(t *testing.T) {
	out := Filter(baseSet())
	if out.Len() != baseSet().Len() {
		t.Fatalf("expected all entries, got %d", out.Len())
	}
}
