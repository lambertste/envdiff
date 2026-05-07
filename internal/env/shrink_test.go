package env

import (
	"strings"
	"testing"
)

func baseShrinkSet() *Set {
	s := New()
	s.Set("APP_NAME", "envdiff")
	s.Set("APP_ENV", "production")
	s.Set("DEBUG_VERBOSE", "true")
	s.Set("LEGACY_FEATURE", "")
	s.Set("TMP_TOKEN", "abc123")
	s.Set("LOG_LEVEL", "info")
	return s
}

func TestShrink_RemoveEmpty(t *testing.T) {
	s := baseShrinkSet()
	out, removed := Shrink(s, ShrinkOptions{RemoveEmpty: true})

	if len(removed) != 1 || removed[0] != "LEGACY_FEATURE" {
		t.Fatalf("expected [LEGACY_FEATURE] removed, got %v", removed)
	}
	if _, ok := out.Get("LEGACY_FEATURE"); ok {
		t.Error("LEGACY_FEATURE should not be present in output")
	}
	if _, ok := out.Get("APP_NAME"); !ok {
		t.Error("APP_NAME should be preserved")
	}
}

func TestShrink_RemoveByPrefix(t *testing.T) {
	s := baseShrinkSet()
	opts := ShrinkOptions{RemovePrefixes: []string{"DEBUG_", "TMP_"}}
	out, removed := Shrink(s, opts)

	if len(removed) != 2 {
		t.Fatalf("expected 2 removed, got %d: %v", len(removed), removed)
	}
	for _, k := range []string{"DEBUG_VERBOSE", "TMP_TOKEN"} {
		if _, ok := out.Get(k); ok {
			t.Errorf("%s should have been removed", k)
		}
	}
}

func TestShrink_RemoveBySuffix(t *testing.T) {
	s := baseShrinkSet()
	opts := ShrinkOptions{RemoveSuffixes: []string{"_LEVEL"}}
	out, removed := Shrink(s, opts)

	if len(removed) != 1 || removed[0] != "LOG_LEVEL" {
		t.Fatalf("expected [LOG_LEVEL], got %v", removed)
	}
	if _, ok := out.Get("LOG_LEVEL"); ok {
		t.Error("LOG_LEVEL should have been removed")
	}
}

func TestShrink_RemoveExplicitKeys(t *testing.T) {
	s := baseShrinkSet()
	opts := ShrinkOptions{RemoveKeys: []string{"APP_ENV", "TMP_TOKEN"}}
	out, removed := Shrink(s, opts)

	if len(removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(removed))
	}
	if _, ok := out.Get("APP_ENV"); ok {
		t.Error("APP_ENV should be removed")
	}
}

func TestShrink_DoesNotMutateOriginal(t *testing.T) {
	s := baseShrinkSet()
	before := s.Keys()
	Shrink(s, DefaultShrinkOptions())
	after := s.Keys()

	if len(before) != len(after) {
		t.Error("original set was mutated")
	}
}

func TestShrinkReport_Empty(t *testing.T) {
	report := ShrinkReport(nil)
	if !strings.Contains(report, "nothing removed") {
		t.Errorf("unexpected report: %q", report)
	}
}

func TestShrinkReport_WithKeys(t *testing.T) {
	report := ShrinkReport([]string{"FOO", "BAR"})
	if !strings.Contains(report, "FOO") || !strings.Contains(report, "BAR") {
		t.Errorf("report missing keys: %q", report)
	}
}
