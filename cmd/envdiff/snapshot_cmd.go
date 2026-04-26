package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envdiff/internal/parser"
	"github.com/yourorg/envdiff/internal/snapshot"
)

func runSnapshot(args []string) error {
	fs := flag.NewFlagSet("snapshot", flag.ContinueOnError)
	label := fs.String("label", "", "label for this snapshot")
	output := fs.String("out", "", "output file path for the snapshot (required)")
	compareWith := fs.String("compare", "", "path to an existing snapshot to compare against")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("snapshot: env file path required")
	}
	envFile := fs.Arg(0)

	envEntries, err := parser.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("snapshot: parse error: %w", err)
	}

	entryMap := make(map[string]string, len(envEntries))
	for _, e := range envEntries {
		entryMap[e.Key] = e.Value
	}

	if *compareWith != "" {
		base, err := snapshot.Load(*compareWith)
		if err != nil {
			return fmt.Errorf("snapshot: load base failed: %w", err)
		}
		current := &snapshot.Snapshot{Label: *label, Entries: entryMap}
		result := snapshot.Compare(base, current)
		fmt.Print(snapshot.FormatDiff(base, current, result))
		return nil
	}

	if *output == "" {
		return fmt.Errorf("snapshot: --out is required when saving a snapshot")
	}
	if *label == "" {
		*label = envFile
	}
	if err := snapshot.Save(*output, *label, entryMap); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Snapshot saved to %s (label: %s, %d entries)\n", *output, *label, len(entryMap))
	return nil
}
