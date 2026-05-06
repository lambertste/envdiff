package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/envdiff/internal/env"
	"github.com/yourorg/envdiff/internal/export"
	"github.com/yourorg/envdiff/internal/parser"
)

func runRedact(args []string) error {
	fs := flag.NewFlagSet("redact", flag.ContinueOnError)
	keys := fs.String("keys", "", "comma-separated list of keys to explicitly redact")
	patterns := fs.String("patterns", "", "comma-separated substrings; matching keys are redacted")
	placeholder := fs.String("placeholder", "***", "replacement value for redacted entries")
	listFlag := fs.Bool("list", false, "print which keys would be redacted, then exit")
	format := fs.String("format", "dotenv", "output format: dotenv|shell|json")
	outFile := fs.String("out", "", "write output to file instead of stdout")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: envdiff redact [flags] <file>")
	}

	f, err := os.Open(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("open %s: %w", fs.Arg(0), err)
	}
	defer f.Close()

	entries, err := parser.ParseReader(f)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	s := env.FromEntries(entries)

	opts := env.DefaultRedactOptions()
	opts.Placeholder = *placeholder

	if *keys != "" {
		for _, k := range strings.Split(*keys, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				opts.Keys = append(opts.Keys, k)
			}
		}
	}
	if *patterns != "" {
		opts.Patterns = nil
		for _, p := range strings.Split(*patterns, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				opts.Patterns = append(opts.Patterns, p)
			}
		}
	}

	if *listFlag {
		redacted := env.RedactedKeys(s, opts)
		fmt.Println(env.FormatRedacted(redacted))
		return nil
	}

	out := env.Redact(s, opts)

	w := os.Stdout
	if *outFile != "" {
		w, err = os.Create(*outFile)
		if err != nil {
			return fmt.Errorf("create %s: %w", *outFile, err)
		}
		defer w.Close()
	}

	return export.Export(out, *format, w)
}
