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
//
// Conflict behaviour:
//   - "base"     keeps the value from <base-file> when the same key appears in both files.
//   - "override" keeps the value from <override-file> when the same key appears in both files.
//
// Conflicts are always reported to stderr regardless of the chosen strategy.
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

	strategy, err := parseStrategy(*strategyFlag)
	if err != nil {
		return err
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

// parseStrategy converts the strategy flag string into a merge.Strategy value.
func parseStrategy(s string) (merge.Strategy, error) {
	switch s {
	case "base":
		return merge.PreferBase, nil
	case "override":
		return merge.PreferOverride, nil
	default:
		return 0, fmt.Errorf("unknown strategy %q: use 'base' or 'override'", s)
	}
}
