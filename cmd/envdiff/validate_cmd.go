package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/validate"
)

// runValidate loads an env file and runs all default rules plus the
// no-plaintext-secret rule, printing any violations to stdout.
// Returns a non-zero exit code if violations are found.
func runValidate(path string, requiredKeys []string) int {
	env, err := parser.ParseFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", path, err)
		return 2
	}

	rules := append(validate.DefaultRules(), validate.NoSecretInPlaintextRule())
	results := validate.Validate(env, rules)

	if len(requiredKeys) > 0 {
		checker := validate.RequiredKeysRule(requiredKeys)
		results = append(results, checker(env)...)
	}

	if len(results) == 0 {
		fmt.Printf("✓ %s passed all validation rules\n", path)
		return 0
	}

	fmt.Printf("✗ %s has %d violation(s):\n", path, len(results))
	for _, r := range results {
		fmt.Printf("  %s\n", r.String())
	}
	return 1
}
