package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/parser"
)

// runSort reads an env file and prints its keys sorted according to the
// specified strategy. Supported strategies: alpha (default), alpha-desc,
// by-value, by-length.
func runSort(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff sort <file> [--by alpha|alpha-desc|by-value|by-length]")
	}

	filePath := args[0]
	strategy := "alpha"

	// Parse optional --by flag.
	for i := 1; i < len(args)-1; i++ {
		if args[i] == "--by" {
			strategy = args[i+1]
		}
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open %s: %w", filePath, err)
	}
	defer f.Close()

	entries, err := parser.ParseReader(f)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filePath, err)
	}

	set := env.FromEntries(entries)

	var order env.SortOrder
	switch strings.ToLower(strategy) {
	case "alpha", "":
		order = env.Alpha
	case "alpha-desc":
		order = env.AlphaDesc
	case "by-value":
		order = env.ByValue
	case "by-length":
		order = env.ByLength
	default:
		return fmt.Errorf("unknown sort strategy %q: want alpha, alpha-desc, by-value, or by-length", strategy)
	}

	sortedEntries := env.SortedEntries(set, order)
	for _, e := range sortedEntries {
		fmt.Printf("%s=%s\n", e.Key, e.Value)
	}

	return nil
}
