package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/parser"
)

// runPivot implements: envdiff pivot --key ENV file1.env file2.env ...
//
// It loads each env file as a Set, pivots them on the given key, and prints
// a comparison table to stdout.
func runPivot(pivotKey string, paths []string) error {
	if pivotKey == "" {
		return fmt.Errorf("pivot: --key flag is required")
	}
	if len(paths) < 2 {
		return fmt.Errorf("pivot: at least two env files are required")
	}

	sets := make([]*env.Set, 0, len(paths))
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return fmt.Errorf("pivot: open %q: %w", p, err)
		}
		entries, err := parser.ParseReader(f)
		f.Close()
		if err != nil {
			return fmt.Errorf("pivot: parse %q: %w", p, err)
		}
		s := env.FromEntries(entries)
		sets = append(sets, s)
	}

	pr, err := env.Pivot(sets, pivotKey)
	if err != nil {
		return fmt.Errorf("pivot: %w", err)
	}

	fmt.Print(env.FormatPivot(pr))
	return nil
}
