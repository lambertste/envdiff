package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/parser"
)

// runDefaults applies default values to an env file.
//
// Usage:
//
//	envdiff defaults -file <path> [-set KEY=VALUE]... [-override] [-missing]
func runDefaults(args []string) error {
	fs := flag.NewFlagSet("defaults", flag.ContinueOnError)
	filePath := fs.String("file", "", "path to .env file")
	override := fs.Bool("override", false, "overwrite existing values")
	listMissing := fs.Bool("missing", false, "only list keys absent from the file")
	var rawSpecs []string
	fs.Func("set", "KEY=VALUE default (repeatable)", func(s string) error {
		rawSpecs = append(rawSpecs, s)
		return nil
	})
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *filePath == "" {
		return fmt.Errorf("defaults: -file is required")
	}

	specs, err := parseDefaultSpecs(rawSpecs, *override)
	if err != nil {
		return err
	}

	entries, err := parser.ParseFile(*filePath)
	if err != nil {
		return fmt.Errorf("defaults: %w", err)
	}
	s := env.FromEntries(entries)

	if *listMissing {
		for _, k := range env.MissingDefaults(s, specs) {
			fmt.Println(k)
		}
		return nil
	}

	out := env.ApplyDefaults(s, specs)
	for _, k := range out.Keys() {
		v, _ := out.Get(k)
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
	}
	return nil
}

func parseDefaultSpecs(raw []string, override bool) ([]env.DefaultSpec, error) {
	specs := make([]env.DefaultSpec, 0, len(raw))
	for _, r := range raw {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("defaults: invalid -set value %q (expected KEY=VALUE)", r)
		}
		specs = append(specs, env.DefaultSpec{Key: parts[0], Value: parts[1], Override: override})
	}
	return specs, nil
}
