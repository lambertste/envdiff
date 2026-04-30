package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/parser"
)

// runPatch applies a JSON patch file to an env file and writes the result.
//
// Usage: envdiff patch <env-file> <patch-file> [--out <output-file>]
//
// The patch file is a JSON array of PatchInstruction objects:
//
//	[{"op":"set","key":"FOO","value":"bar"},{"op":"delete","key":"OLD"}]
func runPatch(args []string, outPath string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff patch <env-file> <patch-file>")
	}

	envFile := args[0]
	patchFile := args[1]

	entries, err := parser.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("parse env file: %w", err)
	}

	s := env.FromEntries(entries)

	patchData, err := os.ReadFile(patchFile)
	if err != nil {
		return fmt.Errorf("read patch file: %w", err)
	}

	type rawInstruction struct {
		Op     string `json:"op"`
		Key    string `json:"key"`
		Value  string `json:"value"`
		NewKey string `json:"new_key"`
	}

	var raw []rawInstruction
	if err := json.Unmarshal(patchData, &raw); err != nil {
		return fmt.Errorf("parse patch JSON: %w", err)
	}

	instructions := make([]env.PatchInstruction, len(raw))
	for i, r := range raw {
		instructions[i] = env.PatchInstruction{
			Op:     env.PatchOp(r.Op),
			Key:    r.Key,
			Value:  r.Value,
			NewKey: r.NewKey,
		}
	}

	result, err := env.Patch(s, instructions)
	if err != nil {
		return fmt.Errorf("apply patch: %w", err)
	}

	w := os.Stdout
	if outPath != "" {
		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	for _, k := range result.Keys() {
		v, _ := result.Get(k)
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
	return nil
}
