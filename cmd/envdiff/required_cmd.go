package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/parser"
)

// runRequired checks that a set of required keys are present and non-empty
// in the given .env file.
//
// Usage: envdiff required <file> --keys KEY1,KEY2 [--strict] [--list-missing]
func runRequired(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff required <file> --keys KEY1,KEY2")
	}

	filePath := args[0]
	var keyList string
	strict := false
	listMissing := false

	for i := 1; i < len(args); i++ {
		switch {
		case args[i] == "--strict":
			strict = true
		case args[i] == "--list-missing":
			listMissing = true
		case strings.HasPrefix(args[i], "--keys="):
			keyList = strings.TrimPrefix(args[i], "--keys=")
		case args[i] == "--keys" && i+1 < len(args):
			i++
			keyList = args[i]
		}
	}

	if keyList == "" {
		return fmt.Errorf("--keys flag is required")
	}

	required := strings.Split(keyList, ",")
	for i, k := range required {
		required[i] = strings.TrimSpace(k)
	}

	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filePath, err)
	}

	s := env.FromEntries(entries)

	if listMissing {
		missing := env.MissingRequired(s, required)
		if len(missing) == 0 {
			fmt.Println("all required keys present")
			return nil
		}
		for _, k := range missing {
			fmt.Println(k)
		}
		if strict {
			return fmt.Errorf("%d required key(s) missing or empty", len(missing))
		}
		return nil
	}

	results := env.CheckRequired(s, required)
	fmt.Fprint(os.Stdout, env.FormatRequired(results))

	if strict {
		missing := env.MissingRequired(s, required)
		if len(missing) > 0 {
			return fmt.Errorf("%d required key(s) missing or empty", len(missing))
		}
	}

	return nil
}
