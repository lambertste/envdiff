package env

import (
	"strings"
	"testing"
)

func baseRequiredSet() *Set {
	s := NewSet()
	s.Set("HOST", "localhost")
	s.Set("PORT", "8080")
	s.Set("EMPTY_KEY", "")
	return s
}

func TestCheckRequired_AllPresent(t *testing.T) {
	s := baseRequiredSet()
	results := CheckRequired(s, []string{"HOST", "PORT"})
	for _, r := range results {
		if !r.Present || r.Empty {
			t.Errorf("expected %s to be present and non-empty", r.Key)
		}
	}
}

func TestCheckRequired_MissingKey(t *testing.T) {
	s := baseRequiredSet()
	results := CheckRequired(s, []string{"MISSING"})
	if len(results) != 1 || results[0].Present {
		t.Errorf("expected MISSING to be absent")
	}
}

func TestCheckRequired_EmptyValue(t *testing.T) {
	s := baseRequiredSet()
	results := CheckRequired(s, []string{"EMPTY_KEY"})
	if len(results) != 1 || !results[0].Present || !results[0].Empty {
		t.Errorf("expected EMPTY_KEY to be present but empty")
	}
}

func TestMissingRequired_NoneAbsent(t *testing.T) {
	s := baseRequiredSet()
	missing := MissingRequired(s, []string{"HOST", "PORT"})
	if len(missing) != 0 {
		t.Errorf("expected no missing keys, got %v", missing)
	}
}

func TestMissingRequired_SomeMissing(t *testing.T) {
	s := baseRequiredSet()
	missing := MissingRequired(s, []string{"HOST", "MISSING", "EMPTY_KEY"})
	if len(missing) != 2 {
		t.Errorf("expected 2 missing/empty, got %v", missing)
	}
}

func TestFormatRequired_OK(t *testing.T) {
	s := baseRequiredSet()
	results := CheckRequired(s, []string{"HOST"})
	out := FormatRequired(results)
	if !strings.Contains(out, "OK") || !strings.Contains(out, "HOST") {
		t.Errorf("unexpected format output: %q", out)
	}
}

func TestFormatRequired_Missing(t *testing.T) {
	s := baseRequiredSet()
	results := CheckRequired(s, []string{"GHOST"})
	out := FormatRequired(results)
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected MISSING label in output: %q", out)
	}
}

func TestFormatRequired_Empty(t *testing.T) {
	s := baseRequiredSet()
	results := CheckRequired(s, []string{"EMPTY_KEY"})
	out := FormatRequired(results)
	if !strings.Contains(out, "EMPTY") {
		t.Errorf("expected EMPTY label in output: %q", out)
	}
}

func TestFormatRequired_NoKeys(t *testing.T) {
	out := FormatRequired(nil)
	if !strings.Contains(out, "no required keys") {
		t.Errorf("expected fallback message, got %q", out)
	}
}
