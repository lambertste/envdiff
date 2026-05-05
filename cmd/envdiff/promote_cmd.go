package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/parser"
)

func runPromote(args []string) error {
	fs := flag.NewFlagSet("promote", flag.ContinueOnError)
	dryRun := fs.Bool("dry-run", false, "preview changes without applying them")
	skipExisting := fs.Bool("skip-existing", false, "do not overwrite keys already present in target")
	keysFlag := fs.String("keys", "", "comma-separated list of keys to promote (default: all)")
	output := fs.String("output", "", "write updated target to this file (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 2 {
		return fmt.Errorf("usage: envdiff promote [flags] <src> <dst>")
	}

	srcPath := fs.Arg(0)
	dstPath := fs.Arg(1)

	srcEntries, err := parser.ParseFile(srcPath)
	if err != nil {
		return fmt.Errorf("reading src %s: %w", srcPath, err)
	}

	dstEntries, err := parser.ParseFile(dstPath)
	if err != nil {
		return fmt.Errorf("reading dst %s: %w", dstPath, err)
	}

	srcSet := env.FromEntries(srcEntries)
	dstSet := env.FromEntries(dstEntries)

	var keys []string
	if *keysFlag != "" {
		for _, k := range strings.Split(*keysFlag, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				keys = append(keys, k)
			}
		}
	}

	opts := env.PromoteOptions{
		DryRun:       *dryRun,
		SkipExisting: *skipExisting,
		Keys:         keys,
	}

	results, err := env.Promote(dstSet, srcSet, opts)
	if err != nil {
		return err
	}

	fmt.Print(env.FormatPromoteResults(results))

	if *dryRun {
		return nil
	}

	w := os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			return fmt.Errorf("opening output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	for _, k := range dstSet.Keys() {
		v, _ := dstSet.Get(k)
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}

	return nil
}
