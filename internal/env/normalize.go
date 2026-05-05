package env

import (
	"strings"
)

// NormalizeOption controls how normalization is applied.
type NormalizeOption func(key, value string) (string, string)

// NormalizeTrimKeys trims whitespace from all keys.
func NormalizeTrimKeys(key, value string) (string, string) {
	return strings.TrimSpace(key), value
}

// NormalizeTrimValues trims whitespace from all values.
func NormalizeTrimValues(key, value string) (string, string) {
	return key, strings.TrimSpace(value)
}

// NormalizeUppercaseKeys converts all keys to uppercase.
func NormalizeUppercaseKeys(key, value string) (string, string) {
	return strings.ToUpper(key), value
}

// NormalizeLowercaseValues converts all values to lowercase.
func NormalizeLowercaseValues(key, value string) (string, string) {
	return key, strings.ToLower(value)
}

// NormalizeCollapseEmptyValues replaces whitespace-only values with empty string.
func NormalizeCollapseEmptyValues(key, value string) (string, string) {
	if strings.TrimSpace(value) == "" {
		return key, ""
	}
	return key, value
}

// Normalize applies a sequence of NormalizeOption functions to every entry in
// the set and returns a new Set with the transformed key/value pairs.
// If two keys collide after normalization the last one wins.
func Normalize(s *Set, opts ...NormalizeOption) *Set {
	out := NewSet()
	for _, k := range s.Keys() {
		v, _ := s.Get(k)
		for _, opt := range opts {
			k, v = opt(k, v)
		}
		out.Set(k, v)
	}
	return out
}

// NormalizedKeys returns the list of keys that would change under the given
// options without constructing a full new Set.
func NormalizedKeys(s *Set, opts ...NormalizeOption) []string {
	var changed []string
	for _, k := range s.Keys() {
		v, _ := s.Get(k)
		nk, nv := k, v
		for _, opt := range opts {
			nk, nv = opt(nk, nv)
		}
		if nk != k || nv != v {
			changed = append(changed, k)
		}
	}
	return changed
}
