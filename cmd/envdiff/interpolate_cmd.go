package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envdiff/internal/env"
	"github.com/yourorg/envdiff/internal/parser"
)

func runInterpolate(args []string) error {
	fs := flag.NewFlagSet("interpolate", flag.ContinueOnError)
	strict := fs.Bool("strict", false, "exit non-zero if any references are unresolved")
	outFmt := fs.String("format", "dotenv", "output format: dotenv|shell")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: envdiff interpolate [flags] <file>")
	}

	f, err := os.Open(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	entries, err := parser.ParseReader(f)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	s := env.FromEntries(entries)
	out, errs := env.Interpolate(s)

	for _, e := range errs {
		fmt.Fprintf(os.Stderr, "warning: %v\n", e)
	}
	if *strict && len(errs) > 0 {
		return fmt.Errorf("interpolation failed: %d unresolved reference(s)", len(errs))
	}

	if err := writeInterpolated(out, *outFmt); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

// writeInterpolated writes all keys from the interpolated store to stdout
// using the specified output format ("shell" or "dotenv").
func writeInterpolated(out env.Store, outFmt string) error {
	for _, key := range out.Keys() {
		val, _ := out.Get(key)
		switch outFmt {
		case "shell":
			_, err := fmt.Printf("%s=%q\n", key, val)
			if err != nil {
				return err
			}
		default:
			_, err := fmt.Printf("%s=%s\n", key, val)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
