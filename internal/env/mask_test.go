package env

import (
	"testing"
)

func baseMaskSet() *Set {
	s := NewSet()
	s.Set("APP_NAME", "myapp")
	s.Set("DB_PASSWORD", "supersecret")
	s.Set("API_KEY", "abc123")
	s.Set("AUTH_TOKEN", "tok-xyz")
	s.Set("PORT", "8080")
	return s
}

func TestMaskSet_MasksSensitiveKeys(t *testing.T) {
	s := baseMaskSet()
	rule := DefaultMaskRule()
	out := MaskSet(s, rule)

	for _, key := range []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN"} {
		val, ok := out.Get(key)
		if !ok {
			t.Fatalf("expected key %q to exist", key)
		}
		if val != "***" {
			t.Errorf("key %q: expected masked value, got %q", key, val)
		}
	}
}

func TestMaskSet_PreservesNonSensitiveKeys(t *testing.T) {
	s := baseMaskSet()
	out := MaskSet(s, DefaultMaskRule())

	for _, tc := range []struct{ key, want string }{
		{"APP_NAME", "myapp"},
		{"PORT", "8080"},
	} {
		val, ok := out.Get(tc.key)
		if !ok {
			t.Fatalf("expected key %q to exist", tc.key)
		}
		if val != tc.want {
			t.Errorf("key %q: got %q, want %q", tc.key, val, tc.want)
		}
	}
}

func TestMaskSet_DoesNotMutateOriginal(t *testing.T) {
	s := baseMaskSet()
	MaskSet(s, DefaultMaskRule())

	val, _ := s.Get("DB_PASSWORD")
	if val != "supersecret" {
		t.Errorf("original set was mutated: got %q", val)
	}
}

func TestMaskSet_CustomMaskString(t *testing.T) {
	s := NewSet()
	s.Set("SECRET_KEY", "topsecret")
	rule := MaskRule{KeyContains: []string{"SECRET"}, MaskWith: "[REDACTED]"}
	out := MaskSet(s, rule)

	val, _ := out.Get("SECRET_KEY")
	if val != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", val)
	}
}

func TestMaskedKeys_ReturnsCorrectKeys(t *testing.T) {
	s := baseMaskSet()
	keys := MaskedKeys(s, DefaultMaskRule())

	expected := map[string]bool{"DB_PASSWORD": true, "API_KEY": true, "AUTH_TOKEN": true}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d masked keys, got %d: %v", len(expected), len(keys), keys)
	}
	for _, k := range keys {
		if !expected[k] {
			t.Errorf("unexpected masked key: %q", k)
		}
	}
}

func TestMaskSet_EmptyMaskWithDefaultsToStars(t *testing.T) {
	s := NewSet()
	s.Set("API_KEY", "val")
	rule := MaskRule{KeyContains: []string{"API_KEY"}, MaskWith: ""}
	out := MaskSet(s, rule)

	val, _ := out.Get("API_KEY")
	if val != "***" {
		t.Errorf("expected default mask ***, got %q", val)
	}
}
