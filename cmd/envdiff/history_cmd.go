package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/parser"
)

// runHistory loads an env file, applies a sequence of patch-style operations
// while recording history, then prints the history log.
//
// Usage:
//
//	envdiff history <file> --set KEY=VALUE --del KEY
func runHistory(args []string) error {
	fs := flag.NewFlagSet("history", flag.ContinueOnError)
	var sets sliceFlag
	var dels sliceFlag
	fs.Var(&sets, "set", "set KEY=VALUE (repeatable)")
	fs.Var(&dels, "del", "delete KEY (repeatable)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: envdiff history <file> [--set K=V] [--del K]")
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
	h := &env.History{}

	for _, kv := range sets {
		key, val, ok := splitKV(kv)
		if !ok {
			return fmt.Errorf("invalid --set value %q: expected KEY=VALUE", kv)
		}
		env.TrackSet(h, s, key, val)
	}

	for _, k := range dels {
		env.TrackDelete(h, s, k)
	}

	fmt.Println(h.Format())
	return nil
}

// splitKV splits "KEY=VALUE" into ("KEY", "VALUE", true).
func splitKV(s string) (string, string, bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			return s[:i], s[i+1:], true
		}
	}
	return "", "", false
}

// sliceFlag is a flag.Value that accumulates repeated string flags.
type sliceFlag []string

func (f *sliceFlag) String() string { return fmt.Sprint([]string(*f)) }
func (f *sliceFlag) Set(v string) error {
	*f = append(*f, v)
	return nil
}
