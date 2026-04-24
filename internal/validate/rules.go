package validate

import (
	"fmt"
	"strings"
)

// RequiredKeysRule returns a Rule that fails if any of the given keys are absent
// from the env map. This is intended to be used alongside Validate by passing
// the env map keys as a set before calling the rule check per-entry.
func RequiredKeysRule(required []string) func(env map[string]string) []Result {
	return func(env map[string]string) []Result {
		var results []Result
		for _, key := range required {
			if _, ok := env[key]; !ok {
				results = append(results, Result{
					Key:     key,
					Rule:    "required-key",
					Message: fmt.Sprintf("required key %q is missing", key),
				})
			}
		}
		return results
	}
}

// NoSecretInPlaintextRule returns a Rule that warns when a key name suggests
// a secret but the value does not look like a placeholder or reference.
func NoSecretInPlaintextRule() Rule {
	secretKeywords := []string{"SECRET", "PASSWORD", "PASSWD", "TOKEN", "PRIVATE_KEY"}
	return Rule{
		Name: "no-plaintext-secret",
		Check: func(key, value string) error {
			for _, kw := range secretKeywords {
				if strings.Contains(strings.ToUpper(key), kw) {
					if !strings.HasPrefix(value, "${{") && !strings.HasPrefix(value, "vault:") {
						return fmt.Errorf(
							"key %q looks like a secret but value does not appear to be a reference",
							key,
						)
					}
				}
			}
			return nil
		},
	}
}
