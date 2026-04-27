package lint

import (
	"testing"
)

func TestLint_NoFindings(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "envdiff",
		"LOG_LEVEL": "info",
	}
	findings := Lint(env, DefaultRules())
	for _, f := range findings {
		if f.Severity == SeverityError {
			t.Errorf("unexpected error finding: %s", f)
		}
	}
}

func TestRuleNoTrailingWhitespace_Triggered(t *testing.T) {
	f := RuleNoTrailingWhitespace("KEY", "value ")
	if f == nil {
		t.Fatal("expected finding, got nil")
	}
	if f.Severity != SeverityWarning {
		t.Errorf("expected warning, got %s", f.Severity)
	}
}

func TestRuleNoTrailingWhitespace_Clean(t *testing.T) {
	f := RuleNoTrailingWhitespace("KEY", "value")
	if f != nil {
		t.Errorf("expected nil, got %s", f)
	}
}

func TestRuleNoLowercaseKey_Triggered(t *testing.T) {
	f := RuleNoLowercaseKey("app_name", "x")
	if f == nil {
		t.Fatal("expected finding, got nil")
	}
	if f.Severity != SeverityWarning {
		t.Errorf("expected warning, got %s", f.Severity)
	}
}

func TestRuleNoLowercaseKey_Clean(t *testing.T) {
	f := RuleNoLowercaseKey("APP_NAME", "x")
	if f != nil {
		t.Errorf("expected nil, got %s", f)
	}
}

func TestRuleNoDuplicateSuffix_Triggered(t *testing.T) {
	f := RuleNoDuplicateSuffix("DATABASE_ENV", "prod")
	if f == nil {
		t.Fatal("expected finding, got nil")
	}
	if f.Severity != SeverityInfo {
		t.Errorf("expected info, got %s", f.Severity)
	}
}

func TestRuleNoDuplicateSuffix_Clean(t *testing.T) {
	f := RuleNoDuplicateSuffix("DATABASE_URL", "postgres://")
	if f != nil {
		t.Errorf("expected nil, got %s", f)
	}
}

func TestRuleWarnEmptyValue_Triggered(t *testing.T) {
	f := RuleWarnEmptyValue("SOME_KEY", "")
	if f == nil {
		t.Fatal("expected finding, got nil")
	}
}

func TestFindingString(t *testing.T) {
	f := Finding{Key: "FOO", Message: "bad value", Severity: SeverityError}
	got := f.String()
	expected := "[error] FOO: bad value"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
