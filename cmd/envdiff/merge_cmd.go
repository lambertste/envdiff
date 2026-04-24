package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/merge"
	"github.com/user/envdiff/internal/output"
	"github.com/user/envdiff/internal/parser"
)

// runMerge executes the merge subcommand.
// Usage: envdiff merge [--strategy=base|override] [--format=text|dotenv|color] <base-file> <override-file>
func runMerge(args []string) error {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	strategyFlag := fs.String("strategy", "base", "conflict resolution strategy: base or override")
	formatFlag := fs.String("format", "text", "output format: text, dotenv, or color")

	if err := fs.Parse(args); err != nil {
		return err
	}

	positional := fs.Args()
	if len(positional) != 2 {
		return fmt.Errorf("merge requires exactly two file arguments: <base> <override>")
	}

	baseFile, overFile := positional[0], positional[1]

	baseEnv, err := parser.ParseFile(baseFile)
	if err != nil {
		return fmt.Errorf("reading base file %q: %w", baseFile, err)
	}

	overEnv, err := parser.ParseFile(overFile)
	if err != nil {
		return fmt.Errorf("reading override file %q: %w", overFile, err)
	}

	var strategy merge.Strategy
	switch *strategyFlag {
	case "base":
		strategy = merge.PreferBase
	case "override":
		strategy = merge.PreferOverride
	default:
		return fmt.Errorf("unknown strategy %q: use 'base' or 'override'", *strategyFlag)
	}

	result := merge.Merge(baseEnv, overEnv, strategy)

	if len(result.Conflicts) > 0 {
		fmt.Fprint(os.Stderr, merge.FormatConflicts(result.Conflicts))
	}

	entries := merge.ToEntries(result.Merged)

	if err := output.Write(os.Stdout, entries, *formatFlag); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	return nil
}
