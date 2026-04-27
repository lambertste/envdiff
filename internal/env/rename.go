package env

import "strings"

// RenameFunc is a function that transforms a key name.
type RenameFunc func(key string) string

// Rename applies the given RenameFunc to all keys in the Set, returning a new
// Set with renamed keys. Values are preserved. If two keys collide after
// renaming, the last one (in insertion order) wins.
func Rename(s *Set, fn RenameFunc) *Set {
	out := NewSet()
	for _, key := range s.Keys() {
		val, _ := s.Get(key)
		newKey := fn(key)
		out.Set(newKey, val)
	}
	return out
}

// AddPrefix returns a RenameFunc that prepends prefix to every key.
func AddPrefix(prefix string) RenameFunc {
	return func(key string) string {
		return prefix + key
	}
}

// StripPrefix returns a RenameFunc that removes prefix from keys that have it.
// Keys without the prefix are left unchanged.
func StripPrefix(prefix string) RenameFunc {
	return func(key string) string {
		return strings.TrimPrefix(key, prefix)
	}
}

// UppercaseKeys returns a RenameFunc that converts keys to uppercase.
func UppercaseKeys() RenameFunc {
	return func(key string) string {
		return strings.ToUpper(key)
	}
}

// ReplaceInKey returns a RenameFunc that replaces all occurrences of old with
// new within each key name.
func ReplaceInKey(old, newStr string) RenameFunc {
	return func(key string) string {
		return strings.ReplaceAll(key, old, newStr)
	}
}
