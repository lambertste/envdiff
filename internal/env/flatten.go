package env

import (
	"fmt"
	"strings"
)

// FlattenOptions controls how nested key segments are joined.
type FlattenOptions struct {
	// Separator is placed between key segments (default "_").
	Separator string
	// UppercaseKeys forces all resulting keys to uppercase.
	UppercaseKeys bool
	// Prefix is prepended to every flattened key.
	Prefix string
}

// DefaultFlattenOptions returns sensible defaults for FlattenOptions.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Separator:     "_",
		UppercaseKeys: true,
	}
}

// Flatten takes a Set whose keys may contain a delimiter (e.g. ".") and
// rewrites them using opts.Separator, optionally uppercasing and prefixing.
// The source delimiter is the first argument; keys that do not contain it
// are passed through unchanged (subject to casing / prefix rules).
func Flatten(s *Set, sourceDelim string, opts FlattenOptions) *Set {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	out := NewSet()
	for _, key := range s.Keys() {
		val, _ := s.Get(key)

		newKey := key
		if sourceDelim != "" && sourceDelim != opts.Separator {
			newKey = strings.ReplaceAll(key, sourceDelim, opts.Separator)
		}
		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		}
		if opts.Prefix != "" {
			newKey = fmt.Sprintf("%s%s%s", opts.Prefix, opts.Separator, newKey)
		}

		out.Set(newKey, val)
	}
	return out
}

// FlattenedKeys returns only the remapped key names without building a new Set.
func FlattenedKeys(s *Set, sourceDelim string, opts FlattenOptions) []string {
	if opts.Separator == "" {
		opts.Separator = "_"
	}
	keys := s.Keys()
	result := make([]string, len(keys))
	for i, key := range keys {
		newKey := key
		if sourceDelim != "" && sourceDelim != opts.Separator {
			newKey = strings.ReplaceAll(key, sourceDelim, opts.Separator)
		}
		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		}
		if opts.Prefix != "" {
			newKey = fmt.Sprintf("%s%s%s", opts.Prefix, opts.Separator, newKey)
		}
		result[i] = newKey
	}
	return result
}
