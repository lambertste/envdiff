package export_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/export"
)

func makeSet(pairs ...string) *env.Set {
	s := env.NewSet()
	for i := 0; i+1 < len(pairs); i += 2 {
		s.Set(pairs[i], pairs[i+1])
	}
	return s
}

func TestExport_Dotenv(t *testing.T) {
	s := makeSet("FOO", "bar", "BAZ", "qux")
	var buf strings.Builder
	if err := export.Export(&buf, s, export.Options{Format: export.FormatDotenv, Sorted: true}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "BAZ=qux") || !strings.Contains(out, "FOO=bar") {
		t.Errorf("unexpected dotenv output: %q", out)
	}
}

func TestExport_JSON(t *testing.T) {
	s := makeSet("KEY", "value")
	var buf strings.Builder
	if err := export.Export(&buf, s, export.Options{Format: export.FormatJSON}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, `"KEY"`) || !strings.Contains(out, `"value"`) {
		t.Errorf("unexpected json output: %q", out)
	}
}

func TestExport_Shell(t *testing.T) {
	s := makeSet("MY_VAR", "hello world")
	var buf strings.Builder
	if err := export.Export(&buf, s, export.Options{Format: export.FormatShell}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "MY_VAR=") {
		t.Errorf("unexpected shell output: %q", out)
	}
}

func TestExport_ExportPrefix(t *testing.T) {
	s := makeSet("PORT", "8080")
	var buf strings.Builder
	if err := export.Export(&buf, s, export.Options{Format: export.FormatExport}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "export PORT=") {
		t.Errorf("expected 'export' prefix, got: %q", out)
	}
}

func TestExport_OmitEmpty(t *testing.T) {
	s := makeSet("FILLED", "yes", "EMPTY", "")
	var buf strings.Builder
	if err := export.Export(&buf, s, export.Options{Format: export.FormatDotenv, OmitEmpty: true}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, "EMPTY") {
		t.Errorf("expected EMPTY to be omitted, got: %q", out)
	}
	if !strings.Contains(out, "FILLED=yes") {
		t.Errorf("expected FILLED in output, got: %q", out)
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	s := makeSet("X", "y")
	var buf strings.Builder
	err := export.Export(&buf, s, export.Options{Format: "xml"})
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestExport_SortedOutput(t *testing.T) {
	s := makeSet("ZEBRA", "1", "ALPHA", "2", "MANGO", "3")
	var buf strings.Builder
	if err := export.Export(&buf, s, export.Options{Format: export.FormatDotenv, Sorted: true}); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 || !strings.HasPrefix(lines[0], "ALPHA") {
		t.Errorf("expected sorted output, got: %v", lines)
	}
}
