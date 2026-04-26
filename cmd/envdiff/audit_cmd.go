package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/audit"
	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// runAudit diffs two env files and emits a structured audit log.
func runAudit(baseFile, overrideFile string) error {
	baseEntries, err := parser.ParseFile(baseFile)
	if err != nil {
		return fmt.Errorf("parsing base file: %w", err)
	}

	overrideEntries, err := parser.ParseFile(overrideFile)
	if err != nil {
		return fmt.Errorf("parsing override file: %w", err)
	}

	diffEntries := diff.Diff(baseEntries, overrideEntries)

	l := &audit.Log{}
	for _, e := range diffEntries {
		switch e.Status {
		case diff.Added:
			l.Record(overrideFile, e.Key, audit.EventAdded, "", e.Value)
		case diff.Removed:
			l.Record(baseFile, e.Key, audit.EventRemoved, e.Value, "")
		case diff.Modified:
			l.Record(overrideFile, e.Key, audit.EventModified, e.OldValue, e.Value)
		}
	}

	fmt.Fprint(os.Stdout, l.Format())

	summary := l.Summary()
	fmt.Fprintf(os.Stdout, "\nsummary: +%d -%d ~%d\n",
		summary[audit.EventAdded],
		summary[audit.EventRemoved],
		summary[audit.EventModified],
	)
	return nil
}
