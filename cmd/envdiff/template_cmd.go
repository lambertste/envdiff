package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/template"
)

// runTemplate handles the `template` subcommand.
// Usage:
//   envdiff template check  <template-file> <env-file>
//   envdiff template gen    <template-file> [env-file]
func runTemplate(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff template <check|gen> <template-file> [env-file]")
	}

	subcmd := args[0]
	tmplPath := args[1]

	tf, err := os.Open(tmplPath)
	if err != nil {
		return fmt.Errorf("opening template file: %w", err)
	}
	defer tf.Close()

	tmpl, err := template.Parse(tf)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	env := map[string]string{}
	if len(args) >= 3 {
		envPath := args[2]
		entries, err := parser.ParseFile(envPath)
		if err != nil {
			return fmt.Errorf("parsing env file: %w", err)
		}
		for _, e := range entries {
			env[e.Key] = e.Value
		}
	}

	switch subcmd {
	case "check":
		return runTemplateCheck(tmpl, env)
	case "gen":
		return runTemplateGen(tmpl, env)
	default:
		return fmt.Errorf("unknown template subcommand: %q", subcmd)
	}
}

func runTemplateCheck(tmpl *template.Template, env map[string]string) error {
	missing := template.Check(tmpl, env)
	if len(missing) == 0 {
		fmt.Println("OK: all required keys are present")
		return nil
	}
	fmt.Fprintf(os.Stderr, "missing required keys:\n  %s\n", strings.Join(missing, "\n  "))
	return fmt.Errorf("%d required key(s) missing", len(missing))
}

func runTemplateGen(tmpl *template.Template, env map[string]string) error {
	output := template.Generate(tmpl, env)
	fmt.Print(output)
	return nil
}
