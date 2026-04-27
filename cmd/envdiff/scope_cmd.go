package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/parser"
)

// runScope reads an env file and prints entries grouped by the given prefixes.
// Usage: envdiff scope <file> [prefix1 prefix2 ...]
func runScope(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff scope <file> [prefix ...]")
	}

	filePath := args[0]
	prefixes := args[1:]

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open %s: %w", filePath, err)
	}
	defer f.Close()

	entries, err := parser.ParseReader(f)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filePath, err)
	}

	s := env.FromEntries(entries)

	if len(prefixes) == 0 {
		// No prefixes given: print all keys flat.
		for _, key := range s.Keys() {
			val, _ := s.Get(key)
			fmt.Printf("%s=%s\n", key, val)
		}
		return nil
	}

	scopes := env.SplitByScope(s, prefixes)
	for _, sc := range scopes {
		header := sc.Name
		if header == "default" {
			header = "(default)"
		}
		fmt.Printf("[%s]\n", strings.TrimRight(header, "_"))
		for _, key := range sc.Entries.Keys() {
			val, _ := sc.Entries.Get(key)
			fmt.Printf("  %s=%s\n", key, val)
		}
		fmt.Println()
	}
	return nil
}
