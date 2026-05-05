package main

import (
	"fmt"
	"os"
	"strings"

	"envdiff/internal/env"
	"envdiff/internal/parser"
)

// runFreeze handles the `envdiff freeze` subcommand.
// Usage:
//
//	envdiff freeze <file> --keys KEY1,KEY2 [--unfreeze] [--list]
func runFreeze(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff freeze <file> [--keys K1,K2] [--unfreeze] [--list]")
	}

	filePath := args[0]
	var keyList string
	unfreeze := false
	list := false

	for i := 1; i < len(args); i++ {
		switch {
		case args[i] == "--unfreeze":
			unfreeze = true
		case args[i] == "--list":
			list = true
		case strings.HasPrefix(args[i], "--keys="):
			keyList = strings.TrimPrefix(args[i], "--keys=")
		case args[i] == "--keys" && i+1 < len(args):
			i++
			keyList = args[i]
		}
	}

	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filePath, err)
	}

	s := env.FromEntries(entries)

	if list {
		frozen := env.FrozenKeys(s)
		fmt.Fprintln(os.Stdout, env.FormatFrozen(frozen))
		return nil
	}

	if keyList == "" {
		return fmt.Errorf("--keys is required unless --list is specified")
	}

	keys := strings.Split(keyList, ",")
	for i, k := range keys {
		keys[i] = strings.TrimSpace(k)
	}

	if unfreeze {
		env.Unfreeze(s, keys...)
		fmt.Fprintf(os.Stdout, "unfroze %d key(s)\n", len(keys))
	} else {
		res := env.Freeze(s, keys...)
		if len(res.Skipped) > 0 {
			fmt.Fprintf(os.Stderr, "warning: skipped missing keys: %s\n", strings.Join(res.Skipped, ", "))
		}
		fmt.Fprintf(os.Stdout, "froze %d key(s)\n", len(res.Frozen))
	}
	return nil
}
