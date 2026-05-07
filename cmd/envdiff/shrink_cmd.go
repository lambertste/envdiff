package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/parser"
)

// runShrink removes entries from an env file based on CLI flags and writes the
// result to stdout (or in-place when --write is set).
func runShrink(args []string, flags map[string]string, write bool, dryRun bool) error {
	if len(args) == 0 {
		return fmt.Errorf("shrink: at least one env file is required")
	}

	filePath := args[0]

	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("shrink: %w", err)
	}

	s := env.FromEntries(entries)

	opts := env.ShrinkOptions{}

	if v, ok := flags["remove-empty"]; ok && v != "false" {
		opts.RemoveEmpty = true
	}
	if v, ok := flags["prefix"]; ok && v != "" {
		opts.RemovePrefixes = splitComma(v)
	}
	if v, ok := flags["suffix"]; ok && v != "" {
		opts.RemoveSuffixes = splitComma(v)
	}
	if v, ok := flags["keys"]; ok && v != "" {
		opts.RemoveKeys = splitComma(v)
	}

	out, removed := env.Shrink(s, opts)

	report := env.ShrinkReport(removed)
	fmt.Fprint(os.Stderr, report)

	if dryRun {
		return nil
	}

	var sb strings.Builder
	for _, k := range out.Keys() {
		v, _ := out.Get(k)
		if strings.ContainsAny(v, " \t") {
			sb.WriteString(fmt.Sprintf("%s=%q\n", k, v))
		} else {
			sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
		}
	}

	if write {
		if err := os.WriteFile(filePath, []byte(sb.String()), 0o644); err != nil {
			return fmt.Errorf("shrink: write %s: %w", filePath, err)
		}
		return nil
	}

	fmt.Print(sb.String())
	return nil
}

func splitComma(s string) []string {
	parts := strings.Split(s, ",")
	out := parts[:0]
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
