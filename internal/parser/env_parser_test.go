package parser

import (
	"bufio"
	"strings"
	"testing"
)

func scannerFrom(s string) *bufio.Scanner {
	return bufio.NewScanner(strings.NewReader(s))
}

func TestParseReader_BasicKeyValue(t *testing.T) {
	input := "APP_ENV=production\nDB_HOST=localhost\n"
	env, err := ParseReader(scannerFrom(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["APP_ENV"] != "production" {
		t.Errorf("expected production, got %q", env["APP_ENV"])
	}
	if env["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %q", env["DB_HOST"])
	}
}

func TestParseReader_SkipsCommentsAndBlanks(t *testing.T) {
	input := "# this is a comment\n\nKEY=value\n"
	env, err := ParseReader(scannerFrom(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 entry, got %d", len(env))
	}
}

func TestParseReader_QuotedValues(t *testing.T) {
	input := `SECRET="my secret value"` + "\n" + `TOKEN='abc123'` + "\n"
	env, err := ParseReader(scannerFrom(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["SECRET"] != "my secret value" {
		t.Errorf("expected unquoted value, got %q", env["SECRET"])
	}
	if env["TOKEN"] != "abc123" {
		t.Errorf("expected abc123, got %q", env["TOKEN"])
	}
}

func TestParseReader_MissingEquals(t *testing.T) {
	input := "INVALID_LINE\n"
	_, err := ParseReader(scannerFrom(input))
	if err == nil {
		t.Fatal("expected error for missing '=', got nil")
	}
}

func TestParseReader_EmptyKey(t *testing.T) {
	input := "=value\n"
	_, err := ParseReader(scannerFrom(input))
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestParseReader_ValueWithEquals(t *testing.T) {
	input := "URL=http://example.com?foo=bar\n"
	env, err := ParseReader(scannerFrom(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["URL"] != "http://example.com?foo=bar" {
		t.Errorf("unexpected value: %q", env["URL"])
	}
}
