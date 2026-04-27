package template

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Entry represents a template entry with an optional description and required flag.
type Entry struct {
	Key         string
	Description string
	Required    bool
	Default     string
}

// Template holds the parsed entries from a .env.template file.
type Template struct {
	Entries []Entry
}

// Parse reads a .env.template file from r and returns a Template.
// Template format:
//   # @required @desc=Some description
//   KEY=default_value
func Parse(r io.Reader) (*Template, error) {
	var tmpl Template
	scanner := bufio.NewScanner(r)

	var pendingDesc string
	var pendingRequired bool

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			pendingDesc = ""
			pendingRequired = false
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			meta := strings.TrimPrefix(trimmed, "#")
			meta = strings.TrimSpace(meta)
			pendingRequired = strings.Contains(meta, "@required")
			if idx := strings.Index(meta, "@desc="); idx != -1 {
				pendingDesc = strings.TrimSpace(meta[idx+6:])
			}
			continue
		}

		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid template line: %q", trimmed)
		}

		entry := Entry{
			Key:         strings.TrimSpace(parts[0]),
			Default:     strings.TrimSpace(parts[1]),
			Description: pendingDesc,
			Required:    pendingRequired,
		}
		tmpl.Entries = append(tmpl.Entries, entry)

		pendingDesc = ""
		pendingRequired = false
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning template: %w", err)
	}

	return &tmpl, nil
}

// Check verifies that all required keys in the template are present in the
// provided env map. Returns a list of missing required keys.
func Check(tmpl *Template, env map[string]string) []string {
	var missing []string
	for _, e := range tmpl.Entries {
		if e.Required {
			if val, ok := env[e.Key]; !ok || strings.TrimSpace(val) == "" {
				missing = append(missing, e.Key)
			}
		}
	}
	return missing
}

// Generate produces a .env file string from the template using the provided
// env map to fill in values, falling back to defaults.
func Generate(tmpl *Template, env map[string]string) string {
	var sb strings.Builder
	for _, e := range tmpl.Entries {
		val, ok := env[e.Key]
		if !ok || val == "" {
			val = e.Default
		}
		if e.Description != "" {
			fmt.Fprintf(&sb, "# %s\n", e.Description)
		}
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, val)
	}
	return sb.String()
}
