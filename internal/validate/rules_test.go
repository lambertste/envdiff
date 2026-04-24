package validate_test

import (
	"testing"

	"github.com/user/envdiff/internal/validate"
)

func TestRequiredKeysRule_AllPresent(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost",
		"PORT":         "8080",
	}
	checker := validate.RequiredKeysRule([]string{"DATABASE_URL", "PORT"})
	results := checker(env)
	if len(results) != 0 {
		t.Errorf("expected no violations, got %v", results)
	}
}

func TestRequiredKeysRule_MissingKey(t *testing.T) {
	env := map[string]string{
		"PORT": "8080",
	}
	checker := validate.RequiredKeysRule([]string{"DATABASE_URL", "PORT"})
	results := checker(env)
	if len(results) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(results))
	}
	if results[0].Key != "DATABASE_URL" {
		t.Errorf("expected missing key DATABASE_URL, got %s", results[0].Key)
	}
	if results[0].Rule != "required-key" {
		t.Errorf("unexpected rule name: %s", results[0].Rule)
	}
}

func TestNoSecretInPlaintextRule_SafeReference(t *testing.T) {
	rule := validate.NoSecretInPlaintextRule()
	if err := rule.Check("API_SECRET", "vault:secret/api#key"); err != nil {
		t.Errorf("expected no error for vault reference, got: %v", err)
	}
}

func TestNoSecretInPlaintextRule_PlaintextSecret(t *testing.T) {
	rule := validate.NoSecretInPlaintextRule()
	if err := rule.Check("DB_PASSWORD", "hunter2"); err == nil {
		t.Error("expected error for plaintext secret value")
	}
}

func TestNoSecretInPlaintextRule_NonSecretKey(t *testing.T) {
	rule := validate.NoSecretInPlaintextRule()
	if err := rule.Check("APP_NAME", "myapp"); err != nil {
		t.Errorf("expected no error for non-secret key, got: %v", err)
	}
}
