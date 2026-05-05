package env

import (
	"testing"
)

func baseFlattenSet() *Set {
	s := NewSet()
	s.Set("app.db.host", "localhost")
	s.Set("app.db.port", "5432")
	s.Set("app.name", "myapp")
	s.Set("LOG_LEVEL", "info")
	return s
}

func TestFlatten_DefaultOptions(t *testing.T) {
	s := baseFlattenSet()
	out := Flatten(s, ".", DefaultFlattenOptions())

	if v, ok := out.Get("APP_DB_HOST"); !ok || v != "localhost" {
		t.Errorf("expected APP_DB_HOST=localhost, got %q ok=%v", v, ok)
	}
	if v, ok := out.Get("APP_DB_PORT"); !ok || v != "5432" {
		t.Errorf("expected APP_DB_PORT=5432, got %q ok=%v", v, ok)
	}
	if v, ok := out.Get("APP_NAME"); !ok || v != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %q ok=%v", v, ok)
	}
}

func TestFlatten_PreservesAlreadyFlatKey(t *testing.T) {
	s := baseFlattenSet()
	out := Flatten(s, ".", DefaultFlattenOptions())

	// LOG_LEVEL had no dots; uppercase keeps it the same
	if v, ok := out.Get("LOG_LEVEL"); !ok || v != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %q ok=%v", v, ok)
	}
}

func TestFlatten_WithPrefix(t *testing.T) {
	s := NewSet()
	s.Set("db.host", "pg")

	opts := DefaultFlattenOptions()
	opts.Prefix = "PROD"
	out := Flatten(s, ".", opts)

	if v, ok := out.Get("PROD_DB_HOST"); !ok || v != "pg" {
		t.Errorf("expected PROD_DB_HOST=pg, got %q ok=%v", v, ok)
	}
}

func TestFlatten_CustomSeparator(t *testing.T) {
	s := NewSet()
	s.Set("app.feature.enabled", "true")

	opts := FlattenOptions{Separator: "__", UppercaseKeys: false}
	out := Flatten(s, ".", opts)

	if v, ok := out.Get("app__feature__enabled"); !ok || v != "true" {
		t.Errorf("expected app__feature__enabled=true, got %q ok=%v", v, ok)
	}
}

func TestFlatten_NoSourceDelim_PassesThrough(t *testing.T) {
	s := NewSet()
	s.Set("already_flat", "yes")

	out := Flatten(s, "", DefaultFlattenOptions())

	if v, ok := out.Get("ALREADY_FLAT"); !ok || v != "yes" {
		t.Errorf("expected ALREADY_FLAT=yes, got %q ok=%v", v, ok)
	}
}

func TestFlattenedKeys_ReturnsCorrectNames(t *testing.T) {
	s := NewSet()
	s.Set("x.y", "1")
	s.Set("a.b", "2")

	keys := FlattenedKeys(s, ".", DefaultFlattenOptions())
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	for _, k := range keys {
		if k != "X_Y" && k != "A_B" {
			t.Errorf("unexpected key %q", k)
		}
	}
}

func TestFlatten_EmptySet(t *testing.T) {
	s := NewSet()
	out := Flatten(s, ".", DefaultFlattenOptions())
	if len(out.Keys()) != 0 {
		t.Errorf("expected empty set, got %d keys", len(out.Keys()))
	}
}
