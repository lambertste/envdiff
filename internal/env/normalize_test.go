package env

import (
	"testing"
)

func baseNormalizeSet() *Set {
	s := NewSet()
	s.Set("db_host", "  localhost  ")
	s.Set("  API_KEY  ", "abc123")
	s.Set("EMPTY_VAL", "   ")
	s.Set("Mixed_Case", "Hello")
	return s
}

func TestNormalize_TrimKeys(t *testing.T) {
	s := baseNormalizeSet()
	out := Normalize(s, NormalizeTrimKeys)
	if _, ok := out.Get("  API_KEY  "); ok {
		t.Error("expected trimmed key, found original")
	}
	if v, ok := out.Get("API_KEY"); !ok || v != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q ok=%v", v, ok)
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	s := baseNormalizeSet()
	out := Normalize(s, NormalizeTrimValues)
	v, _ := out.Get("db_host")
	if v != "localhost" {
		t.Errorf("expected trimmed value, got %q", v)
	}
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	s := baseNormalizeSet()
	out := Normalize(s, NormalizeUppercaseKeys)
	if _, ok := out.Get("Mixed_Case"); ok {
		t.Error("expected uppercase key, found original")
	}
	if _, ok := out.Get("MIXED_CASE"); !ok {
		t.Error("expected MIXED_CASE key after normalization")
	}
}

func TestNormalize_CollapseEmptyValues(t *testing.T) {
	s := baseNormalizeSet()
	out := Normalize(s, NormalizeCollapseEmptyValues)
	v, ok := out.Get("EMPTY_VAL")
	if !ok {
		t.Fatal("expected EMPTY_VAL to exist")
	}
	if v != "" {
		t.Errorf("expected empty string, got %q", v)
	}
}

func TestNormalize_ChainedOptions(t *testing.T) {
	s := baseNormalizeSet()
	out := Normalize(s, NormalizeTrimKeys, NormalizeTrimValues, NormalizeUppercaseKeys)
	v, ok := out.Get("API_KEY")
	if !ok {
		t.Fatal("expected API_KEY after chained normalization")
	}
	if v != "abc123" {
		t.Errorf("unexpected value %q", v)
	}
	v2, ok2 := out.Get("DB_HOST")
	if !ok2 || v2 != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q ok=%v", v2, ok2)
	}
}

func TestNormalize_DoesNotMutateOriginal(t *testing.T) {
	s := baseNormalizeSet()
	Normalize(s, NormalizeUppercaseKeys, NormalizeTrimValues)
	if _, ok := s.Get("db_host"); !ok {
		t.Error("original set was mutated")
	}
}

func TestNormalizedKeys_ReportsChangedKeys(t *testing.T) {
	s := baseNormalizeSet()
	changed := NormalizedKeys(s, NormalizeTrimKeys)
	found := false
	for _, k := range changed {
		if k == "  API_KEY  " {
			found = true
		}
	}
	if !found {
		t.Error("expected '  API_KEY  ' in changed keys")
	}
}
