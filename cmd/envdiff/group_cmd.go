package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/parser"
)

// runGroup reads an env file and groups keys by a common prefix (split on "_").
// Usage: envdiff group <file> [--list] [--group <name>]
func runGroup(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff group <file> [--list] [--group <name>]")
	}

	filePath := args[0]
	listOnly := false
	filterGroup := ""

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--list":
			listOnly = true
		case "--group":
			if i+1 >= len(args) {
				return fmt.Errorf("--group requires a name argument")
			}
			i++
			filterGroup = args[i]
		}
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open %s: %w", filePath, err)
	}
	defer f.Close()

	entries, err := parser.ParseReader(f)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filePath, err)
	}

	s := env.FromEntries(entries)
	gr := env.GroupBy(s, func(key string) string {
		if idx := strings.Index(key, "_"); idx > 0 {
			return key[:idx]
		}
		return ""
	})

	if listOnly {
		for _, name := range env.GroupNames(gr) {
			fmt.Println(name)
		}
		return nil
	}

	for _, name := range env.GroupNames(gr) {
		if filterGroup != "" && name != filterGroup {
			continue
		}
		fmt.Printf("[%s]\n", name)
		g := gr[name]
		for _, k := range g.Keys() {
			v, _ := g.Get(k)
			fmt.Printf("  %s=%s\n", k, v)
		}
	}
	return nil
}
