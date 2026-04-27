package template_test

import (
	"os"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/template"
)

func TestParse_SampleTemplateFile(t *testing.T) {
	f, err := os.Open("testdata/sample.env.template")
	if err != nil {
		t.Fatalf("opening sample template: %v", err)
	}
	defer f.Close()

	tmpl, err := template.Parse(f)
	if err != nil {
		t.Fatalf("parsing sample template: %v", err)
	}

	if len(tmpl.Entries) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(tmpl.Entries))
	}

	requiredKeys := map[string]bool{}
	for _, e := range tmpl.Entries {
		if e.Required {
			requiredKeys[e.Key] = true
		}
	}

	for _, k := range []string{"APP_ENV", "DATABASE_URL", "SECRET_KEY"} {
		if !requiredKeys[k] {
			t.Errorf("expected %s to be required", k)
		}
	}
}

func TestCheck_SampleTemplate_MissingKeys(t *testing.T) {
	f, err := os.Open("testdata/sample.env.template")
	if err != nil {
		t.Fatalf("opening sample template: %v", err)
	}
	defer f.Close()

	tmpl, _ := template.Parse(f)

	env := map[string]string{
		"APP_ENV": "staging",
		"PORT":    "9000",
		// DATABASE_URL and SECRET_KEY intentionally missing
	}

	missing := template.Check(tmpl, env)
	if len(missing) != 2 {
		t.Fatalf("expected 2 missing keys, got %d: %v", len(missing), missing)
	}
}

func TestGenerate_SampleTemplate_OutputContainsAllKeys(t *testing.T) {
	f, err := os.Open("testdata/sample.env.template")
	if err != nil {
		t.Fatalf("opening sample template: %v", err)
	}
	defer f.Close()

	tmpl, _ := template.Parse(f)

	env := map[string]string{
		"APP_ENV":      "production",
		"DATABASE_URL": "postgres://localhost/mydb",
		"SECRET_KEY":   "supersecret",
	}

	output := template.Generate(tmpl, env)

	for _, key := range []string{"APP_ENV", "PORT", "DATABASE_URL", "LOG_LEVEL", "SECRET_KEY"} {
		if !strings.Contains(output, key+"=") {
			t.Errorf("expected output to contain %s=, got:\n%s", key, output)
		}
	}

	if !strings.Contains(output, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output")
	}
	if !strings.Contains(output, "PORT=8080") {
		t.Errorf("expected PORT=8080 (default) in output")
	}
}
