package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/export"
	"github.com/user/envdiff/internal/parser"
)

// runExport reads an env file and writes it in the requested format.
// Usage: envdiff export [--format dotenv|json|shell|export] [--sorted] [--omit-empty] <file>
func runExport(args []string) error {
	var (
		formatFlag  = "dotenv"
		sorted      = false
		omitEmpty   = false
		positional  []string
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--format", "-f":
			if i+1 >= len(args) {
				return fmt.Errorf("--format requires a value")
			}
			i++
			formatFlag = args[i]
		case "--sorted", "-s":
			sorted = true
		case "--omit-empty", "-e":
			omitEmpty = true
		default:
			positional = append(positional, args[i])
		}
	}

	if len(positional) < 1 {
		return fmt.Errorf("usage: envdiff export [--format <fmt>] [--sorted] [--omit-empty] <file>")
	}

	filePath := positional[0]
	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", filePath, err)
	}

	s := env.FromEntries(entries)

	opts := export.Options{
		Format:    export.Format(formatFlag),
		Sorted:    sorted,
		OmitEmpty: omitEmpty,
	}

	if err := export.Export(os.Stdout, s, opts); err != nil {
		return fmt.Errorf("exporting: %w", err)
	}
	return nil
}
