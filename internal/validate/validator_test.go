package validate_test

import (
	"testing"

	"github.com/user/envdiff/internal/validate"
)

func TestValidate_ValidEnv(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
	}
	results := validate.Validate(env, validate.DefaultRules())
	if len(results) != 0 {
		t.Errorf("expected no violations, got %d: %v", len(results), results)
	}
}

func TestValidate_InvalidKeyFormat(t *testing.T) {
	env := map[string]string{
		"invalid-key": "value",
	}
	results := validate.Validate(env, validate.DefaultRules())
	if !hasRule(results, "key-format") {
		t.Error("expected key-format violation")
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	env := map[string]string{
		"SOME_KEY": "",
	}
	results := validate.Validate(env, validate.DefaultRules())
	if !hasRule(results, "no-empty-value") {
		t.Error("expected no-empty-value violation")
	}
}

func TestValidate_WhitespaceInKey(t *testing.T) {
	env := map[string]string{
		"BAD KEY": "value",
	}
	results := validate.Validate(env, validate.DefaultRules())
	if !hasRule(results, "no-whitespace-in-key") {
		t.Error("expected no-whitespace-in-key violation")
	}
}

func TestValidate_MultipleViolations(t *testing.T) {
	env := map[string]string{
		"bad-key":  "",
		"GOOD_KEY": "ok",
	}
	results := validate.Validate(env, validate.DefaultRules())
	if len(results) < 2 {
		t.Errorf("expected at least 2 violations, got %d", len(results))
	}
}

func TestResult_String(t *testing.T) {
	r := validate.Result{Key: "FOO", Rule: "key-format", Message: "some error"}
	s := r.String()
	if s == "" {
		t.Error("expected non-empty string from Result.String()")
	}
}

func hasRule(results []validate.Result, rule string) bool {
	for _, r := range results {
		if r.Rule == rule {
			return true
		}
	}
	return false
}
