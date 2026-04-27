package env

import (
	"strings"
)

// TransformFn is a function that transforms an entry's value given its key.
type TransformFn func(key, value string) string

// Transform applies one or more TransformFns to every entry in the Set,
// returning a new Set with the transformed values. The original Set is
// not modified.
func Transform(s *Set, fns ...TransformFn) *Set {
	out := NewSet()
	for _, key := range s.Keys() {
		val, _ := s.Get(key)
		for _, fn := range fns {
			val = fn(key, val)
		}
		out.Set(key, val)
	}
	return out
}

// UppercaseValues returns a TransformFn that converts every value to uppercase.
func UppercaseValues() TransformFn {
	return func(_, value string) string {
		return strings.ToUpper(value)
	}
}

// TrimValues returns a TransformFn that trims leading/trailing whitespace
// from every value.
func TrimValues() TransformFn {
	return func(_, value string) string {
		return strings.TrimSpace(value)
	}
}

// MaskSecrets returns a TransformFn that replaces the value with a redaction
// placeholder when the key contains any of the given substrings (case-insensitive).
func MaskSecrets(keywords ...string) TransformFn {
	return func(key, value string) string {
		lower := strings.ToLower(key)
		for _, kw := range keywords {
			if strings.Contains(lower, strings.ToLower(kw)) {
				return "***REDACTED***"
			}
		}
		return value
	}
}

// PrefixValues returns a TransformFn that prepends prefix to every value.
func PrefixValues(prefix string) TransformFn {
	return func(_, value string) string {
		return prefix + value
	}
}
