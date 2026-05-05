package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/parser"
)

// runNormalize reads an env file, applies the requested normalization passes,
// and writes the result to stdout (or back to the file when --in-place is set).
func runNormalize(args []string, trimKeys, trimValues, uppercaseKeys, collapseEmpty, inPlace, listChanged bool) error {
	if len(args) < 1 {
		return fmt.Errorf("normalize: env file path required")
	}
	path := args[0]

	entries, err := parser.ParseFile(path)
	if err != nil {
		return fmt.Errorf("normalize: %w", err)
	}

	s := env.FromEntries(entries)

	var opts []env.NormalizeOption
	if trimKeys {
		opts = append(opts, env.NormalizeTrimKeys)
	}
	if trimValues {
		opts = append(opts, env.NormalizeTrimValues)
	}
	if uppercaseKeys {
		opts = append(opts, env.NormalizeUppercaseKeys)
	}
	if collapseEmpty {
		opts = append(opts, env.NormalizeCollapseEmptyValues)
	}

	if listChanged {
		changed := env.NormalizedKeys(s, opts...)
		if len(changed) == 0 {
			fmt.Println("no keys would change")
			return nil
		}
		for _, k := range changed {
			fmt.Println(k)
		}
		return nil
	}

	out := env.Normalize(s, opts...)

	var sb strings.Builder
	for _, k := range out.Keys() {
		v, _ := out.Get(k)
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(v)
		sb.WriteByte('\n')
	}

	if inPlace {
		if err := os.WriteFile(path, []byte(sb.String()), 0644); err != nil {
			return fmt.Errorf("normalize: write %s: %w", path, err)
		}
		fmt.Fprintf(os.Stderr, "normalized %s\n", path)
		return nil
	}

	fmt.Print(sb.String())
	return nil
}
