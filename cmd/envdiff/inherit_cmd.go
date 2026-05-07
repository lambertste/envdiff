package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/parser"
)

// runInherit merges a parent env file into a child env file, filling missing keys.
//
// Usage: envdiff inherit [--overwrite] [--include-empty] <child> <parent>
func runInherit(args []string, overwrite, includeEmpty bool, output string) error {
	if len(args) < 2 {
		return fmt.Errorf("inherit: requires <child> and <parent> file arguments")
	}

	childEntries, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("inherit: reading child %q: %w", args[0], err)
	}
	parentEntries, err := parser.ParseFile(args[1])
	if err != nil {
		return fmt.Errorf("inherit: reading parent %q: %w", args[1], err)
	}

	child := env.FromEntries(childEntries)
	parent := env.FromEntries(parentEntries)

	opts := env.DefaultInheritOptions()
	opts.OverwriteExisting = overwrite
	opts.SkipEmpty = !includeEmpty

	result, res, err := env.Inherit(child, parent, opts)
	if err != nil {
		return fmt.Errorf("inherit: %w", err)
	}

	// Print summary to stderr.
	fmt.Fprint(os.Stderr, env.FormatInheritResult(res))

	// Write result.
	w := os.Stdout
	if output != "" && output != "-" {
		f, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("inherit: opening output %q: %w", output, err)
		}
		defer f.Close()
		w = f
	}

	for _, k := range result.Keys() {
		v, _ := result.Get(k)
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
	return nil
}
