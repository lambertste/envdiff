package env_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/env"
	"github.com/yourorg/envdiff/internal/parser"
)

const interpFixture = `
HOME=/home/ci
CONFIG=${HOME}/.config
APP_CONFIG=${CONFIG}/myapp
LOG_DIR=$HOME/logs
PLAIN=static-value
`

func TestInterpolate_Integration_FullExpansion(t *testing.T) {
	entries, err := parser.ParseReader(strings.NewReader(interpFixture))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := env.FromEntries(entries)
	out, errs := env.Interpolate(s)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}

	cases := map[string]string{
		"HOME":       "/home/ci",
		"CONFIG":     "/home/ci/.config",
		"LOG_DIR":    "/home/ci/logs",
		"PLAIN":      "static-value",
	}
	for key, want := range cases {
		got, ok := out.Get(key)
		if !ok {
			t.Errorf("key %q not found in output", key)
			continue
		}
		if got != want {
			t.Errorf("%s: expected %q, got %q", key, want, got)
		}
	}
}

func TestInterpolate_Integration_PartialMissing(t *testing.T) {
	const src = `
GOOD=present
BAD=${MISSING_VAR}/suffix
`
	entries, _ := parser.ParseReader(strings.NewReader(src))
	s := env.FromEntries(entries)
	_, errs := env.Interpolate(s)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
}
