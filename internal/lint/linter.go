package lint

import (
	"fmt"
	"strings"
)

// Severity represents the level of a lint finding.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Finding represents a single lint result for a key.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s: %s", f.Severity, f.Key, f.Message)
}

// Rule is a lint rule applied to a key-value pair.
type Rule func(key, value string) *Finding

// DefaultRules returns the built-in set of lint rules.
func DefaultRules() []Rule {
	return []Rule{
		RuleNoTrailingWhitespace,
		RuleNoLowercaseKey,
		RuleNoDuplicateSuffix,
		RuleWarnEmptyValue,
	}
}

// Lint applies the given rules to each key-value pair and returns all findings.
func Lint(env map[string]string, rules []Rule) []Finding {
	findings := []Finding{}
	for k, v := range env {
		for _, rule := range rules {
			if f := rule(k, v); f != nil {
				findings = append(findings, *f)
			}
		}
	}
	return findings
}

// RuleNoTrailingWhitespace flags values with leading or trailing whitespace.
func RuleNoTrailingWhitespace(key, value string) *Finding {
	if value != strings.TrimSpace(value) {
		return &Finding{Key: key, Message: "value has leading or trailing whitespace", Severity: SeverityWarning}
	}
	return nil
}

// RuleNoLowercaseKey flags keys that contain lowercase letters.
func RuleNoLowercaseKey(key, _ string) *Finding {
	if key != strings.ToUpper(key) {
		return &Finding{Key: key, Message: "key contains lowercase letters; prefer ALL_CAPS", Severity: SeverityWarning}
	}
	return nil
}

// RuleNoDuplicateSuffix flags keys that end with common redundant suffixes like _KEY or _VAR.
func RuleNoDuplicateSuffix(key, _ string) *Finding {
	redundant := []string{"_VAR", "_ENV"}
	for _, suffix := range redundant {
		if strings.HasSuffix(key, suffix) {
			return &Finding{Key: key, Message: fmt.Sprintf("key ends with redundant suffix %q", suffix), Severity: SeverityInfo}
		}
	}
	return nil
}

// RuleWarnEmptyValue flags keys whose value is empty.
func RuleWarnEmptyValue(key, value string) *Finding {
	if value == "" {
		return &Finding{Key: key, Message: "value is empty", Severity: SeverityWarning}
	}
	return nil
}
