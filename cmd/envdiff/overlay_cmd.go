package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/envdiff/internal/env"
	"github.com/yourorg/envdiff/internal/parser"
)

// runOverlay merges multiple .env files as ordered layers into a single output.
// Usage: envdiff overlay [--no-overwrite] [--skip-empty] [--report] <base> <layer1> [layer2...]
func runOverlay(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("overlay: at least two files required (base + one layer)")
	}

	noOverwrite := false
	skipEmpty := false
	showReport := false
	var files []string

	for _, a := range args {
		switch a {
		case "--no-overwrite":
			noOverwrite = true
		case "--skip-empty":
			skipEmpty = true
		case "--report":
			showReport = true
		default:
			files = append(files, a)
		}
	}

	if len(files) < 2 {
		return fmt.Errorf("overlay: at least two file paths required")
	}

	layers := make([]*env.Set, 0, len(files))
	for _, f := range files {
		entries, err := parser.ParseFile(f)
		if err != nil {
			return fmt.Errorf("overlay: reading %s: %w", f, err)
		}
		layers = append(layers, env.FromEntries(entries))
	}

	opts := env.DefaultOverlayOptions()
	opts.Overwrite = !noOverwrite
	opts.SkipEmpty = skipEmpty

	if showReport && len(layers) == 2 {
		out, report, err := env.OverlayWithReport(layers[0], layers[1], opts)
		if err != nil {
			return err
		}
		printOverlayReport(report)
		return printSet(out)
	}

	out, err := env.Overlay(layers, opts)
	if err != nil {
		return err
	}
	return printSet(out)
}

func printOverlayReport(r env.OverlayReport) {
	if len(r.Added) > 0 {
		fmt.Fprintf(os.Stderr, "+ added:      %s\n", strings.Join(r.Added, ", "))
	}
	if len(r.Overwritten) > 0 {
		fmt.Fprintf(os.Stderr, "~ overwritten: %s\n", strings.Join(r.Overwritten, ", "))
	}
	if len(r.Preserved) > 0 {
		fmt.Fprintf(os.Stderr, "= preserved:  %s\n", strings.Join(r.Preserved, ", "))
	}
}

func printSet(s *env.Set) error {
	for _, k := range s.Keys() {
		v, _ := s.Get(k)
		fmt.Printf("%s=%s\n", k, v)
	}
	return nil
}
