package env

import "strings"

// Scope represents a named grouping of env entries by key prefix.
type Scope struct {
	Name    string
	Entries *Set
}

// SplitByScope partitions the entries in s into groups based on known prefixes.
// Keys that match no prefix are placed in a scope named "default".
func SplitByScope(s *Set, prefixes []string) []Scope {
	buckets := make(map[string]*Set)

	for _, key := range s.Keys() {
		val, _ := s.Get(key)
		matched := false
		for _, p := range prefixes {
			if strings.HasPrefix(key, p) {
				if buckets[p] == nil {
					buckets[p] = NewSet()
				}
				buckets[p].Set(key, val)
				matched = true
				break
			}
		}
		if !matched {
			if buckets["default"] == nil {
				buckets["default"] = NewSet()
			}
			buckets["default"].Set(key, val)
		}
	}

	// Build ordered result: explicit prefixes first, then default.
	var scopes []Scope
	for _, p := range prefixes {
		if set, ok := buckets[p]; ok {
			scopes = append(scopes, Scope{Name: p, Entries: set})
		}
	}
	if def, ok := buckets["default"]; ok {
		scopes = append(scopes, Scope{Name: "default", Entries: def})
	}
	return scopes
}

// MergeScopes combines all scopes back into a single Set.
// Later scopes overwrite earlier ones on key collision.
func MergeScopes(scopes []Scope) *Set {
	out := NewSet()
	for _, sc := range scopes {
		for _, key := range sc.Entries.Keys() {
			val, _ := sc.Entries.Get(key)
			out.Set(key, val)
		}
	}
	return out
}
