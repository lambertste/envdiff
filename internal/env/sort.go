package env

import (
	"sort"
	"strings"
)

// SortOrder defines how entries should be sorted.
type SortOrder int

const (
	SortAlpha      SortOrder = iota // alphabetical by key
	SortAlphaDesc                   // reverse alphabetical by key
	SortByValue                     // alphabetical by value
	SortByLength                    // ascending key length
)

// SortedKeys returns the keys of the set in the specified order.
func SortedKeys(s *Set, order SortOrder) []string {
	keys := s.Keys()

	switch order {
	case SortAlpha:
		sort.Strings(keys)
	case SortAlphaDesc:
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	case SortByValue:
		sort.Slice(keys, func(i, j int) bool {
			vi, _ := s.Get(keys[i])
			vj, _ := s.Get(keys[j])
			return strings.ToLower(vi) < strings.ToLower(vj)
		})
	case SortByLength:
		sort.Slice(keys, func(i, j int) bool {
			return len(keys[i]) < len(keys[j])
		})
	}

	return keys
}

// SortedEntries returns key-value pairs from the set in the specified order.
func SortedEntries(s *Set, order SortOrder) [][2]string {
	keys := SortedKeys(s, order)
	result := make([][2]string, 0, len(keys))
	for _, k := range keys {
		v, _ := s.Get(k)
		result = append(result, [2]string{k, v})
	}
	return result
}
