package env

import (
	"strings"
	"testing"
)

func baseRedactSet() *Set {
	s := NewSet()
	s.Set("APP_NAME", "myapp")
	s.Set("DB_PASSWORD", "s3cr3t")
	s.Set("API_KEY", "abc123")
	s.Set("PORT", "8080")
	s.Set("JWT_TOKEN", "tok.en.value")
	s.Set("DEBUG", "true")
	return s
}

func TestRedact_DefaultOptions_MasksSensitiveKeys(t *testing.T) {
	s := baseRedactSet()
	out := Redact(s, DefaultRedactOptions())

	for _, k := range []string{"DB_PASSWORD", "API_KEY", "JWT_TOKEN"} {
		v, _ := out.Get(k)
		if v != "***" {
			t.Errorf("expected %s to be redacted, got %q", k, v)
		}
	}
}

func TestRedact_DefaultOptions_PreservesNonSensitiveKeys(t *testing.T) {
	s := baseRedactSet()
	out := Redact(s, DefaultRedactOptions())

	for _, k := range []string{"APP_NAME", "PORT", "DEBUG"} {
		v, _ := out.Get(k)
		original, _ := s.Get(k)
		if v != original {
			t.Errorf("expected %s to be preserved as %q, got %q", k, original, v)
		}
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	s := baseRedactSet()
	Redact(s, DefaultRedactOptions())

	v, _ := s.Get("DB_PASSWORD")
	if v != "s3cr3t" {
		t.Errorf("original set was mutated: DB_PASSWORD = %q", v)
	}
}

func TestRedact_ExplicitKeys(t *testing.T) {
	s := baseRedactSet()
	opts := RedactOptions{
		Keys:        []string{"APP_NAME", "PORT"},
		Placeholder: "REDACTED",
	}
	out := Redact(s, opts)

	for _, k := range []string{"APP_NAME", "PORT"} {
		v, _ := out.Get(k)
		if v != "REDACTED" {
			t.Errorf("expected %s to be REDACTED, got %q", k, v)
		}
	}

	v, _ := out.Get("DEBUG")
	if v != "true" {
		t.Errorf("expected DEBUG to be preserved, got %q", v)
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	s := baseRedactSet()
	opts := DefaultRedactOptions()
	opts.Placeholder = "<hidden>"
	out := Redact(s, opts)

	v, _ := out.Get("API_KEY")
	if v != "<hidden>" {
		t.Errorf("expected <hidden>, got %q", v)
	}
}

func TestRedactedKeys_ReturnsExpectedKeys(t *testing.T) {
	s := baseRedactSet()
	keys := RedactedKeys(s, DefaultRedactOptions())

	if len(keys) != 3 {
		t.Errorf("expected 3 redacted keys, got %d: %v", len(keys), keys)
	}
}

func TestFormatRedacted_Empty(t *testing.T) {
	out := FormatRedacted(nil)
	if out != "no keys redacted" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatRedacted_WithKeys(t *testing.T) {
	out := FormatRedacted([]string{"API_KEY", "DB_PASSWORD"})
	if !strings.Contains(out, "API_KEY") || !strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("expected both keys in output, got: %q", out)
	}
	if !strings.Contains(out, "2 key(s)") {
		t.Errorf("expected count in output, got: %q", out)
	}
}
