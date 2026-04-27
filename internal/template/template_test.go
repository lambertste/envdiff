package template

import (
	"strings"
	"testing"
)

func readerFrom(s string) *strings.Reader {
	return strings.NewReader(s)
}

func TestParse_BasicEntries(t *testing.T) {
	input := "APP_ENV=production\nPORT=8080\n"
	tmpl, err := Parse(readerFrom(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tmpl.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(tmpl.Entries))
	}
	if tmpl.Entries[0].Key != "APP_ENV" || tmpl.Entries[0].Default != "production" {
		t.Errorf("unexpected first entry: %+v", tmpl.Entries[0])
	}
}

func TestParse_RequiredAndDescription(t *testing.T) {
	input := "# @required @desc=Database URL\nDATABASE_URL=\n"
	tmpl, err := Parse(readerFrom(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tmpl.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(tmpl.Entries))
	}
	e := tmpl.Entries[0]
	if !e.Required {
		t.Error("expected Required=true")
	}
	if e.Description != "Database URL" {
		t.Errorf("expected description 'Database URL', got %q", e.Description)
	}
}

func TestParse_InvalidLine(t *testing.T) {
	_, err := Parse(readerFrom("NODEQUALS\n"))
	if err == nil {
		t.Error("expected error for missing '='")
	}
}

func TestCheck_AllPresent(t *testing.T) {
	tmpl := &Template{
		Entries: []Entry{
			{Key: "APP_ENV", Required: true},
			{Key: "PORT", Required: false},
		},
	}
	env := map[string]string{"APP_ENV": "staging", "PORT": "8080"}
	missing := Check(tmpl, env)
	if len(missing) != 0 {
		t.Errorf("expected no missing keys, got %v", missing)
	}
}

func TestCheck_MissingRequired(t *testing.T) {
	tmpl := &Template{
		Entries: []Entry{
			{Key: "SECRET_KEY", Required: true},
		},
	}
	missing := Check(tmpl, map[string]string{})
	if len(missing) != 1 || missing[0] != "SECRET_KEY" {
		t.Errorf("expected [SECRET_KEY], got %v", missing)
	}
}

func TestGenerate_UsesEnvOverDefault(t *testing.T) {
	tmpl := &Template{
		Entries: []Entry{
			{Key: "PORT", Default: "3000"},
		},
	}
	result := Generate(tmpl, map[string]string{"PORT": "9090"})
	if !strings.Contains(result, "PORT=9090") {
		t.Errorf("expected PORT=9090 in output, got: %s", result)
	}
}

func TestGenerate_FallsBackToDefault(t *testing.T) {
	tmpl := &Template{
		Entries: []Entry{
			{Key: "LOG_LEVEL", Default: "info", Description: "Logging level"},
		},
	}
	result := Generate(tmpl, map[string]string{})
	if !strings.Contains(result, "LOG_LEVEL=info") {
		t.Errorf("expected LOG_LEVEL=info in output, got: %s", result)
	}
	if !strings.Contains(result, "# Logging level") {
		t.Errorf("expected description comment in output, got: %s", result)
	}
}
