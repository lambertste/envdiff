package env

import (
	"fmt"
	"strings"
)

// RedactOptions controls how redaction is applied.
type RedactOptions struct {
	// Keys is the explicit list of keys to redact.
	Keys []string
	// Patterns are substring patterns; any key containing a pattern is redacted.
	Patterns []string
	// Placeholder replaces the redacted value (default "***").
	Placeholder string
}

// defaultRedactPatterns are common sensitive key substrings.
var defaultRedactPatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE_KEY", "CREDENTIAL",
}

// DefaultRedactOptions returns options that redact common secret keys.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		Patterns:    defaultRedactPatterns,
		Placeholder: "***",
	}
}

// Redact returns a new Set with sensitive values replaced by the placeholder.
// The original Set is not mutated.
func Redact(s *Set, opts RedactOptions) *Set {
	if opts.Placeholder == "" {
		opts.Placeholder = "***"
	}

	explicit := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		explicit[k] = struct{}{}
	}

	out := NewSet()
	for _, k := range s.Keys() {
		v, _ := s.Get(k)
		if shouldRedact(k, explicit, opts.Patterns) {
			v = opts.Placeholder
		}
		out.Set(k, v)
	}
	return out
}

// RedactedKeys returns the list of keys that would be redacted by the given options.
func RedactedKeys(s *Set, opts RedactOptions) []string {
	explicit := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		explicit[k] = struct{}{}
	}

	var redacted []string
	for _, k := range s.Keys() {
		if shouldRedact(k, explicit, opts.Patterns) {
			redacted = append(redacted, k)
		}
	}
	return redacted
}

// FormatRedacted returns a human-readable summary of which keys were redacted.
func FormatRedacted(keys []string) string {
	if len(keys) == 0 {
		return "no keys redacted"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "redacted %d key(s):\n", len(keys))
	for _, k := range keys {
		fmt.Fprintf(&sb, "  - %s\n", k)
	}
	return strings.TrimRight(sb.String(), "\n")
}

func shouldRedact(key string, explicit map[string]struct{}, patterns []string) bool {
	if _, ok := explicit[key]; ok {
		return true
	}
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}
