package export

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/env"
)

// Format represents an export output format.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatShell  Format = "shell"
	FormatExport Format = "export"
)

// Options controls export behaviour.
type Options struct {
	Format  Format
	Sorted  bool
	OmitEmpty bool
}

// Export writes the entries from s to w in the requested format.
func Export(w io.Writer, s *env.Set, opts Options) error {
	entries := s.Entries()
	if opts.OmitEmpty {
		filtered := entries[:0]
		for _, e := range entries {
			if e.Value != "" {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}
	if opts.Sorted {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Key < entries[j].Key
		})
	}

	switch opts.Format {
	case FormatDotenv, "":
		return writeDotenv(w, entries)
	case FormatJSON:
		return writeJSON(w, entries)
	case FormatShell:
		return writeShell(w, entries)
	case FormatExport:
		return writeExport(w, entries)
	default:
		return fmt.Errorf("unknown export format: %q", opts.Format)
	}
}

func writeDotenv(w io.Writer, entries []env.Entry) error {
	for _, e := range entries {
		val := quoteIfNeeded(e.Value)
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, val); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, entries []env.Entry) error {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(m)
}

func writeShell(w io.Writer, entries []env.Entry) error {
	for _, e := range entries {
		val := shellEscape(e.Value)
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, val); err != nil {
			return err
		}
	}
	return nil
}

func writeExport(w io.Writer, entries []env.Entry) error {
	for _, e := range entries {
		val := shellEscape(e.Value)
		if _, err := fmt.Fprintf(w, "export %s=%s\n", e.Key, val); err != nil {
			return err
		}
	}
	return nil
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t#") {
		return `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
	}
	return v
}

func shellEscape(v string) string {
	return "'" + strings.ReplaceAll(v, "'", `'\''`) + "'"
}
