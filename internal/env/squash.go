package env

import "strings"

// SquashOptions controls how duplicate prefixes are squashed.
type SquashOptions struct {
	// Separator is the delimiter used to split key segments (default "_").
	Separator string
	// KeepFirst retains the first occurrence of a prefix group instead of the last.
	KeepFirst bool
}

// DefaultSquashOptions returns sensible defaults for Squash.
func DefaultSquashOptions() SquashOptions {
	return SquashOptions{
		Separator: "_",
		KeepFirst: false,
	}
}

// SquashReport describes the result of a Squash operation.
type SquashReport struct {
	// Removed holds keys that were dropped during squashing.
	Removed []string
	// Kept holds keys that were retained.
	Kept []string
}

// Squash reduces a Set by collapsing keys that share the same top-level
// prefix segment, keeping only one representative per prefix group.
// Keys with no separator are left untouched.
func Squash(s *Set, opts SquashOptions) (*Set, SquashReport) {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	// Track which prefix groups we have already seen.
	seen := make(map[string]bool)
	out := New()
	report := SquashReport{}

	keys := s.Keys()
	if opts.KeepFirst {
		// iterate forward so first wins
	} else {
		// reverse so that when we mark seen we keep the last occurrence
		reversed := make([]string, len(keys))
		for i, k := range keys {
			reversed[len(keys)-1-i] = k
		}
		keys = reversed
	}

	for _, k := range keys {
		parts := strings.SplitN(k, opts.Separator, 2)
		prefix := parts[0]
		if len(parts) == 1 {
			// no separator — always keep
			v, _ := s.Get(k)
			out.Set(k, v)
			report.Kept = append(report.Kept, k)
			continue
		}
		if seen[prefix] {
			report.Removed = append(report.Removed, k)
			continue
		}
		seen[prefix] = true
		v, _ := s.Get(k)
		out.Set(k, v)
		report.Kept = append(report.Kept, k)
	}

	return out, report
}

// FormatSquashReport returns a human-readable summary of a SquashReport.
func FormatSquashReport(r SquashReport) string {
	if len(r.Removed) == 0 {
		return "squash: nothing removed\n"
	}
	var sb strings.Builder
	sb.WriteString("squash: removed keys:\n")
	for _, k := range r.Removed {
		sb.WriteString("  - " + k + "\n")
	}
	return sb.String()
}
