package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/parser"
)

// runSchema validates an env file against a simple inline schema.
//
// Usage:
//
//	envdiff schema --file .env --require KEY1,KEY2 --pattern KEY=REGEX
func runSchema(args []string) error {
	fs := flag.NewFlagSet("schema", flag.ContinueOnError)
	filePath := fs.String("file", "", "path to .env file (required)")
	requireFlag := fs.String("require", "", "comma-separated list of required keys")
	patternFlag := fs.String("pattern", "", "KEY=REGEX pairs, comma-separated")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *filePath == "" {
		return fmt.Errorf("--file is required")
	}

	f, err := os.Open(*filePath)
	if err != nil {
		return fmt.Errorf("open %s: %w", *filePath, err)
	}
	defer f.Close()

	entries, err := parser.ParseReader(f)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	set := env.FromEntries(entries)

	schema, err := buildSchema(*requireFlag, *patternFlag)
	if err != nil {
		return err
	}

	violations := schema.Validate(set)
	fmt.Println(env.FormatViolations(violations))
	if len(violations) > 0 {
		return fmt.Errorf("schema validation failed with %d violation(s)", len(violations))
	}
	return nil
}

func buildSchema(requireFlag, patternFlag string) (*env.Schema, error) {
	fieldMap := map[string]*env.SchemaField{}

	if requireFlag != "" {
		for _, key := range strings.Split(requireFlag, ",") {
			key = strings.TrimSpace(key)
			if key == "" {
				continue
			}
			fieldMap[key] = &env.SchemaField{Key: key, Required: true}
		}
	}

	if patternFlag != "" {
		for _, spec := range strings.Split(patternFlag, ",") {
			parts := strings.SplitN(strings.TrimSpace(spec), "=", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid --pattern spec %q: expected KEY=REGEX", spec)
			}
			key, raw := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			re, err := regexp.Compile(raw)
			if err != nil {
				return nil, fmt.Errorf("invalid regex for key %s: %w", key, err)
			}
			if f, ok := fieldMap[key]; ok {
				f.Pattern = re
			} else {
				fieldMap[key] = &env.SchemaField{Key: key, Pattern: re}
			}
		}
	}

	schema := &env.Schema{}
	for _, f := range fieldMap {
		schema.Fields = append(schema.Fields, *f)
	}
	return schema, nil
}
