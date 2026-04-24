package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "APP_ENV", Kind: diff.Added, NewValue: "production"},
		{Key: "DB_PASS", Kind: diff.Removed, OldValue: "secret"},
		{Key: "PORT", Kind: diff.Modified, OldValue: "3000", NewValue: "8080"},
	}
}

func TestWrite_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, makeEntries(), FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ APP_ENV=production") {
		t.Errorf("expected added line, got:\n%s", out)
	}
	if !strings.Contains(out, "- DB_PASS=secret") {
		t.Errorf("expected removed line, got:\n%s", out)
	}
	if !strings.Contains(out, "~ PORT: 3000 -> 8080") {
		t.Errorf("expected modified line, got:\n%s", out)
	}
}

func TestWrite_DotenvFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, makeEntries(), FormatDotenv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected dotenv added line, got:\n%s", out)
	}
	if !strings.Contains(out, "# REMOVED: DB_PASS") {
		t.Errorf("expected dotenv removed comment, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected dotenv modified line with new value, got:\n%s", out)
	}
}

func TestWrite_ColorFormat_ContainsKey(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, makeEntries(), FormatColor); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	// Keys should still appear even with color codes
	for _, key := range []string{"APP_ENV", "DB_PASS", "PORT"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %q in color output, got:\n%s", key, out)
		}
	}
}

func TestWrite_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, []diff.Entry{}, FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for no entries")
	}
}
