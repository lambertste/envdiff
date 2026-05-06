package env

import (
	"regexp"
	"strings"
	"testing"
)

func baseSchemaSet() *Set {
	s := NewSet()
	s.Set("APP_ENV", "production")
	s.Set("PORT", "8080")
	s.Set("DATABASE_URL", "postgres://localhost/mydb")
	return s
}

func TestSchema_NoViolations(t *testing.T) {
	schema := &Schema{
		Fields: []SchemaField{
			{Key: "APP_ENV", Required: true},
			{Key: "PORT", Required: true},
		},
	}
	violations := schema.Validate(baseSchemaSet())
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestSchema_MissingRequiredKey(t *testing.T) {
	schema := &Schema{
		Fields: []SchemaField{
			{Key: "MISSING_KEY", Required: true},
		},
	}
	violations := schema.Validate(baseSchemaSet())
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "MISSING_KEY" {
		t.Errorf("expected key MISSING_KEY, got %s", violations[0].Key)
	}
}

func TestSchema_PatternMismatch(t *testing.T) {
	schema := &Schema{
		Fields: []SchemaField{
			{Key: "PORT", Required: true, Pattern: regexp.MustCompile(`^\d{4}$`)},
		},
	}
	s := NewSet()
	s.Set("PORT", "notaport")
	violations := schema.Validate(s)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "does not match pattern") {
		t.Errorf("unexpected message: %s", violations[0].Message)
	}
}

func TestSchema_PatternMatch(t *testing.T) {
	schema := &Schema{
		Fields: []SchemaField{
			{Key: "PORT", Required: true, Pattern: regexp.MustCompile(`^\d+$`)},
		},
	}
	violations := schema.Validate(baseSchemaSet())
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestSchema_OptionalMissingKeyNoViolation(t *testing.T) {
	schema := &Schema{
		Fields: []SchemaField{
			{Key: "OPTIONAL_KEY", Required: false},
		},
	}
	violations := schema.Validate(baseSchemaSet())
	if len(violations) != 0 {
		t.Fatalf("expected no violations for optional key, got %v", violations)
	}
}

func TestFormatViolations_Empty(t *testing.T) {
	out := FormatViolations(nil)
	if out != "schema: all checks passed" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatViolations_NonEmpty(t *testing.T) {
	v := []SchemaViolation{
		{Key: "FOO", Message: "required key is missing or empty"},
		{Key: "BAR", Message: "value \"x\" does not match pattern"},
	}
	out := FormatViolations(v)
	if !strings.Contains(out, "2 violation(s)") {
		t.Errorf("expected violation count in output: %s", out)
	}
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "BAR") {
		t.Errorf("expected keys in output: %s", out)
	}
}
