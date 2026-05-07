package env

import (
	"strings"
	"testing"
)

func baseSquashSet() *Set {
	s := New()
	s.Set("DB_HOST", "localhost")
	s.Set("DB_PORT", "5432")
	s.Set("APP_NAME", "envdiff")
	s.Set("APP_ENV", "production")
	s.Set("LOG_LEVEL", "info")
	s.Set("STANDALONE", "yes")
	return s
}

func TestSquash_KeepLast_OnePerPrefix(t *testing.T) {
	s := baseSquashSet()
	out, _ := Squash(s, DefaultSquashOptions())

	keys := out.Keys()
	prefixCount := map[string]int{}
	for _, k := range keys {
		parts := strings.SplitN(k, "_", 2)
		prefixCount[parts[0]]++
	}
	for prefix, count := range prefixCount {
		if count > 1 {
			t.Errorf("prefix %q appears %d times, want 1", prefix, count)
		}
	}
}

func TestSquash_KeepFirst_RetainsFirstKey(t *testing.T) {
	s := New()
	s.Set("DB_HOST", "first")
	s.Set("DB_PORT", "second")

	opts := DefaultSquashOptions()
	opts.KeepFirst = true
	out, report := Squash(s, opts)

	v, ok := out.Get("DB_HOST")
	if !ok || v != "first" {
		t.Errorf("expected DB_HOST=first, got %q ok=%v", v, ok)
	}
	if len(report.Removed) != 1 || report.Removed[0] != "DB_PORT" {
		t.Errorf("expected DB_PORT removed, got %v", report.Removed)
	}
}

func TestSquash_StandaloneKeyPreserved(t *testing.T) {
	s := baseSquashSet()
	out, _ := Squash(s, DefaultSquashOptions())

	v, ok := out.Get("STANDALONE")
	if !ok || v != "yes" {
		t.Errorf("expected STANDALONE=yes, got %q ok=%v", v, ok)
	}
}

func TestSquash_ReportCounts(t *testing.T) {
	s := baseSquashSet()
	_, report := Squash(s, DefaultSquashOptions())

	// DB and APP each have 2 keys; one from each should be removed.
	if len(report.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d: %v", len(report.Removed), report.Removed)
	}
}

func TestSquash_CustomSeparator(t *testing.T) {
	s := New()
	s.Set("db.host", "localhost")
	s.Set("db.port", "5432")
	s.Set("app.name", "envdiff")

	opts := DefaultSquashOptions()
	opts.Separator = "."
	out, report := Squash(s, opts)

	if len(out.Keys()) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out.Keys()))
	}
	if len(report.Removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(report.Removed))
	}
}

func TestFormatSquashReport_NoneRemoved(t *testing.T) {
	r := SquashReport{}
	out := FormatSquashReport(r)
	if !strings.Contains(out, "nothing removed") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatSquashReport_WithRemoved(t *testing.T) {
	r := SquashReport{Removed: []string{"DB_PORT", "APP_ENV"}}
	out := FormatSquashReport(r)
	if !strings.Contains(out, "DB_PORT") || !strings.Contains(out, "APP_ENV") {
		t.Errorf("expected removed keys in output, got: %q", out)
	}
}
