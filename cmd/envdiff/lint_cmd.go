package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/lint"
	"github.com/user/envdiff/internal/parser"
)

func runLint(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff lint <file>")
	}

	filePath := args[0]

	env, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	rules := lint.DefaultRules()
	findings := lint.Lint(env, rules)

	if len(findings) == 0 {
		fmt.Println("No lint findings.")
		return nil
	}

	errCount := 0
	warnCount := 0
	infoCount := 0

	for _, f := range findings {
		switch f.Severity {
		case lint.SeverityError:
			fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", f)
			errCount++
		case lint.SeverityWarning:
			fmt.Fprintf(os.Stderr, "\033[33m%s\033[0m\n", f)
			warnCount++
		default:
			fmt.Println(f)
			infoCount++
		}
	}

	fmt.Fprintf(os.Stderr, "\n%d error(s), %d warning(s), %d info(s)\n", errCount, warnCount, infoCount)

	if errCount > 0 {
		return fmt.Errorf("lint failed with %d error(s)", errCount)
	}
	return nil
}
