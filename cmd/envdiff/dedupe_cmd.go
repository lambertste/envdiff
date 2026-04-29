package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/parser"
)

// runDedupe implements the `envdiff dedupe` sub-command.
// Usage: envdiff dedupe [--strategy=first|last] [--list-dupes] <file>
func runDedupe(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("dedupe: no input file specified")
	}

	strategy := env.DedupeKeepLast
	listDupes := false
	filePath := ""

	for _, arg := range args {
		switch {
		case arg == "--strategy=first":
			strategy = env.DedupeKeepFirst
		case arg == "--strategy=last":
			strategy = env.DedupeKeepLast
		case arg == "--list-dupes":
			listDupes = true
		case !strings.HasPrefix(arg, "--"):
			filePath = arg
		}
	}

	if filePath == "" {
		return fmt.Errorf("dedupe: no input file specified")
	}

	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("dedupe: parse error: %w", err)
	}

	set := env.FromEntries(entries)
	result := env.Dedupe(set, strategy)

	if listDupes {
		if len(result.Duplicates) == 0 {
			fmt.Fprintln(os.Stdout, "no duplicate keys found")
		} else {
			fmt.Fprintln(os.Stdout, "duplicate keys:")
			for _, k := range result.Duplicates {
				fmt.Fprintf(os.Stdout, "  %s\n", k)
			}
		}
		return nil
	}

	for _, k := range result.Set.Keys() {
		v, _ := result.Set.Get(k)
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
	}
	return nil
}
