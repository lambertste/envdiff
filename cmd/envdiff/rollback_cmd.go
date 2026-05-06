package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"envdiff/internal/env"
	"envdiff/internal/parser"
)

func runRollback(args []string) error {
	fs := flag.NewFlagSet("rollback", flag.ContinueOnError)
	keys := fs.String("keys", "", "comma-separated list of keys to roll back (required)")
	dryRun := fs.Bool("dry-run", false, "print rollback plan without applying")
	outFile := fs.String("out", "", "output file (default: overwrite input)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		return fmt.Errorf("usage: rollback --keys KEY1,KEY2 <before.env> <current.env>")
	}
	if *keys == "" {
		return fmt.Errorf("--keys is required")
	}

	beforePath := fs.Arg(0)
	currentPath := fs.Arg(1)

	beforeEntries, err := parser.ParseFile(beforePath)
	if err != nil {
		return fmt.Errorf("reading before file: %w", err)
	}
	currentEntries, err := parser.ParseFile(currentPath)
	if err != nil {
		return fmt.Errorf("reading current file: %w", err)
	}

	before := env.FromEntries(beforeEntries)
	current := env.FromEntries(currentEntries)

	keyList := strings.Split(*keys, ",")
	for i, k := range keyList {
		keyList[i] = strings.TrimSpace(k)
	}

	plan := env.SnapshotKeys(before, keyList)

	if *dryRun {
		fmt.Print(env.FormatRollback(plan))
		return nil
	}

	result := env.Rollback(current, plan)

	dest := currentPath
	if *outFile != "" {
		dest = *outFile
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("opening output: %w", err)
	}
	defer f.Close()

	for _, k := range result.Keys() {
		v, _ := result.Get(k)
		fmt.Fprintf(f, "%s=%s\n", k, v)
	}
	return nil
}
