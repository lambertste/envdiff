package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule represents a validation rule applied to env var keys or values.
type Rule struct {
	Name    string
	Check   func(key, value string) error
}

// Result holds the outcome of validating a single entry.
type Result struct {
	Key     string
	Rule    string
	Message string
}

func (r Result) String() string {
	return fmt.Sprintf("[%s] %s: %s", r.Rule, r.Key, r.Message)
}

var keyPattern = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// DefaultRules returns the standard set of validation rules.
func DefaultRules() []Rule {
	return []Rule{
		{
			Name: "key-format",
			Check: func(key, _ string) error {
				if !keyPattern.MatchString(key) {
					return fmt.Errorf("key %q must match [A-Z_][A-Z0-9_]*", key)
				}
				return nil
			},
		},
		{
			Name: "no-empty-value",
			Check: func(key, value string) error {
				if strings.TrimSpace(value) == "" {
					return fmt.Errorf("key %q has an empty value", key)
				}
				return nil
			},
		},
		{
			Name: "no-whitespace-in-key",
			Check: func(key, _ string) error {
				if strings.ContainsAny(key, " \t") {
					return fmt.Errorf("key %q contains whitespace", key)
				}
				return nil
			},
		},
	}
}

// Validate runs all rules against the provided key-value map and returns
// a slice of Results for every violation found.
func Validate(env map[string]string, rules []Rule) []Result {
	var results []Result
	for key, value := range env {
		for _, rule := range rules {
			if err := rule.Check(key, value); err != nil {
				results = append(results, Result{
					Key:     key,
					Rule:    rule.Name,
					Message: err.Error(),
				})
			}
		}
	}
	return results
}
