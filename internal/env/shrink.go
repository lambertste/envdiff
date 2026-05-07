package env

import "strings"

// ShrinkOptions controls which entries are removed during shrinking.
type ShrinkOptions struct {
	// RemoveEmpty removes entries with empty values.
	RemoveEmpty bool
	// RemovePrefixes removes entries whose keys start with any of these prefixes.
	RemovePrefixes []string
	// RemoveSuffixes removes entries whose keys end with any of these suffixes.
	RemoveSuffixes []string
	// RemoveKeys removes entries with exactly these keys.
	RemoveKeys []string
}

// DefaultShrinkOptions returns conservative defaults that only strip empty values.
func DefaultShrinkOptions() ShrinkOptions {
	return ShrinkOptions{RemoveEmpty: true}
}

// Shrink returns a new Set with entries removed according to opts.
// The original Set is not mutated.
func Shrink(s *Set, opts ShrinkOptions) (*Set, []string) {
	excluded := map[string]struct{}{}

	for _, k := range opts.RemoveKeys {
		excluded[k] = struct{}{}
	}

	out := New()
	var removed []string

	for _, k := range s.Keys() {
		v, _ := s.Get(k)

		if _, ok := excluded[k]; ok {
			removed = append(removed, k)
			continue
		}
		if opts.RemoveEmpty && v == "" {
			removed = append(removed, k)
			continue
		}
		if matchesAnyPrefix(k, opts.RemovePrefixes) {
			removed = append(removed, k)
			continue
		}
		if matchesAnySuffix(k, opts.RemoveSuffixes) {
			removed = append(removed, k)
			continue
		}

		out.Set(k, v)
	}

	return out, removed
}

// ShrinkReport summarises what was removed.
func ShrinkReport(removed []string) string {
	if len(removed) == 0 {
		return "shrink: nothing removed\n"
	}
	var sb strings.Builder
	sb.WriteString("shrink: removed keys:\n")
	for _, k := range removed {
		sb.WriteString("  - ")
		sb.WriteString(k)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func matchesAnyPrefix(k string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(k, p) {
			return true
		}
	}
	return false
}

func matchesAnySuffix(k string, suffixes []string) bool {
	for _, s := range suffixes {
		if strings.HasSuffix(k, s) {
			return true
		}
	}
	return false
}
