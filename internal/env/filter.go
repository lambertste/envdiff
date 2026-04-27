package env

import "strings"

// FilterFunc is a predicate over an Entry.
type FilterFunc func(e Entry) bool

// Filter returns a new Set containing only entries that satisfy all predicates.
func Filter(s *Set, predicates ...FilterFunc) *Set {
	out := NewSet()
	for _, e := range s.Entries() {
		pass := true
		for _, fn := range predicates {
			if !fn(e) {
				pass = false
				break
			}
		}
		if pass {
			out.Set(e.Key, e.Value)
		}
	}
	return out
}

// WithPrefix returns a FilterFunc that keeps entries whose key starts with prefix.
func WithPrefix(prefix string) FilterFunc {
	return func(e Entry) bool {
		return strings.HasPrefix(e.Key, prefix)
	}
}

// WithSuffix returns a FilterFunc that keeps entries whose key ends with suffix.
func WithSuffix(suffix string) FilterFunc {
	return func(e Entry) bool {
		return strings.HasSuffix(e.Key, suffix)
	}
}

// NonEmpty returns a FilterFunc that excludes entries with empty values.
func NonEmpty() FilterFunc {
	return func(e Entry) bool {
		return strings.TrimSpace(e.Value) != ""
	}
}

// ExcludeKeys returns a FilterFunc that drops entries whose key is in the given set.
func ExcludeKeys(keys ...string) FilterFunc {
	skip := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		skip[k] = struct{}{}
	}
	return func(e Entry) bool {
		_, found := skip[e.Key]
		return !found
	}
}
