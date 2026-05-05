package main

import (
	"flag"
	"fmt"
	"os"

	"envdiff/internal/env"
	"envdiff/internal/parser"
)

// runChain merges multiple .env files using a layered chain strategy.
// Usage: envdiff chain [--overwrite] file1.env file2.env ...
func runChain(args []string) error {
	fs := flag.NewFlagSet("chain", flag.ContinueOnError)
	overwrite := fs.Bool("overwrite", false, "later files overwrite earlier keys (default: first wins)")
	format := fs.String("format", "dotenv", "output format: dotenv|shell|json")
	if err := fs.Parse(args); err != nil {
		return err
	}

	paths := fs.Args()
	if len(paths) == 0 {
		return fmt.Errorf("chain: at least one env file required")
	}

	var sets []*env.Set
	for _, p := range paths {
		entries, err := parser.ParseFile(p)
		if err != nil {
			return fmt.Errorf("chain: reading %s: %w", p, err)
		}
		s := env.FromEntries(entries)
		sets = append(sets, s)
	}

	var opts []env.ChainOption
	if *overwrite {
		opts = append(opts, env.WithOverwrite())
	}

	result := env.Chain(sets, opts...)

	switch *format {
	case "shell":
		for _, k := range result.Keys() {
			v, _ := result.Get(k)
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
	case "json":
		fmt.Fprintln(os.Stdout, "{")
		keys := result.Keys()
		for i, k := range keys {
			v, _ := result.Get(k)
			comma := ","
			if i == len(keys)-1 {
				comma = ""
			}
			fmt.Fprintf(os.Stdout, "  %q: %q%s\n", k, v, comma)
		}
		fmt.Fprintln(os.Stdout, "}")
	default: // dotenv
		for _, k := range result.Keys() {
			v, _ := result.Get(k)
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
	}
	return nil
}
