package env

import (
	"testing"
)

func baseRenameSet() *Set {
	s := NewSet()
	s.Set("APP_HOST", "localhost")
	s.Set("APP_PORT", "8080")
	s.Set("DB_HOST", "db.local")
	return s
}

func TestRename_AddPrefix(t *testing.T) {
	s := baseRenameSet()
	out := Rename(s, AddPrefix("PROD_"))

	if _, ok := out.Get("PROD_APP_HOST"); !ok {
		t.Error("expected PROD_APP_HOST to exist")
	}
	if _, ok := out.Get("PROD_DB_HOST"); !ok {
		t.Error("expected PROD_DB_HOST to exist")
	}
	if len(out.Keys()) != 3 {
		t.Errorf("expected 3 keys, got %d", len(out.Keys()))
	}
}

func TestRename_StripPrefix(t *testing.T) {
	s := baseRenameSet()
	out := Rename(s, StripPrefix("APP_"))

	if v, ok := out.Get("HOST"); !ok || v != "localhost" {
		t.Errorf("expected HOST=localhost, got %q ok=%v", v, ok)
	}
	if v, ok := out.Get("PORT"); !ok || v != "8080" {
		t.Errorf("expected PORT=8080, got %q ok=%v", v, ok)
	}
	// DB_HOST has no APP_ prefix — key should be unchanged
	if _, ok := out.Get("DB_HOST"); !ok {
		t.Error("expected DB_HOST to remain unchanged")
	}
}

func TestRename_UppercaseKeys(t *testing.T) {
	s := NewSet()
	s.Set("app_host", "localhost")
	s.Set("app_port", "8080")

	out := Rename(s, UppercaseKeys())

	if _, ok := out.Get("APP_HOST"); !ok {
		t.Error("expected APP_HOST after uppercase rename")
	}
	if _, ok := out.Get("APP_PORT"); !ok {
		t.Error("expected APP_PORT after uppercase rename")
	}
}

func TestRename_ReplaceInKey(t *testing.T) {
	s := NewSet()
	s.Set("APP-HOST", "localhost")
	s.Set("APP-PORT", "8080")

	out := Rename(s, ReplaceInKey("-", "_"))

	if _, ok := out.Get("APP_HOST"); !ok {
		t.Error("expected APP_HOST after dash replacement")
	}
	if _, ok := out.Get("APP_PORT"); !ok {
		t.Error("expected APP_PORT after dash replacement")
	}
}

func TestRename_PreservesValues(t *testing.T) {
	s := NewSet()
	s.Set("KEY", "secret-value")

	out := Rename(s, AddPrefix("X_"))

	if v, ok := out.Get("X_KEY"); !ok || v != "secret-value" {
		t.Errorf("expected X_KEY=secret-value, got %q ok=%v", v, ok)
	}
}

func TestRename_CollisionLastWins(t *testing.T) {
	s := NewSet()
	s.Set("FOO", "first")
	s.Set("BAR", "second")

	// Both keys become "SAME" — last insertion order wins
	out := Rename(s, func(_ string) string { return "SAME" })

	if len(out.Keys()) != 1 {
		t.Errorf("expected 1 key after collision, got %d", len(out.Keys()))
	}
}
