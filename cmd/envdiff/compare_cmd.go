package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/parser"
)

// runCompare diffs two env files and prints a structured summary.
func runCompare(baseFile, overrideFile string, showUnchanged bool) error {
	if baseFile == "" || overrideFile == "" {
		return fmt.Errorf("compare requires two file arguments")
	}

	aEntries, err := parser.ParseFile(baseFile)
	if err != nil {
		return fmt.Errorf("reading base file: %w", err)
	}

	bEntries, err := parser.ParseFile(overrideFile)
	if err != nil {
		return fmt.Errorf("reading override file: %w", err)
	}

	a := env.FromEntries(aEntries)
	b := env.FromEntries(bEntries)

	result := env.Compare(a, b)

	if !result.HasChanges() {
		fmt.Println("No differences found.")
		return nil
	}

	for _, k := range result.Added {
		v, _ := b.Get(k)
		fmt.Fprintf(os.Stdout, "+ %s=%s\n", k, v)
	}

	for _, k := range result.Removed {
		v, _ := a.Get(k)
		fmt.Fprintf(os.Stdout, "- %s=%s\n", k, v)
	}

	for _, k := range result.Modified {
		av, _ := a.Get(k)
		bv, _ := b.Get(k)
		fmt.Fprintf(os.Stdout, "~ %s: %s -> %s\n", k, av, bv)
	}

	if showUnchanged {
		for _, k := range result.Unchanged {
			v, _ := a.Get(k)
			fmt.Fprintf(os.Stdout, "  %s=%s\n", k, v)
		}
	}

	fmt.Println()
	fmt.Println(result.Summary())
	return nil
}
