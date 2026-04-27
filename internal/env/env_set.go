package env

import (
	"fmt"
	"sort"
	"strings"
)

// Entry represents a single key-value pair from an env file.
type Entry struct {
	Key   string
	Value string
}

// Set is an ordered, keyed collection of env entries.
type Set struct {
	entries map[string]string
	keys    []string
}

// NewSet creates an empty Set.
func NewSet() *Set {
	return &Set{entries: make(map[string]string)}
}

// FromEntries builds a Set from a slice of Entry.
func FromEntries(entries []Entry) *Set {
	s := NewSet()
	for _, e := range entries {
		s.Set(e.Key, e.Value)
	}
	return s
}

// Set adds or updates a key.
func (s *Set) Set(key, value string) {
	if _, exists := s.entries[key]; !exists {
		s.keys = append(s.keys, key)
	}
	s.entries[key] = value
}

// Get returns the value and whether the key exists.
func (s *Set) Get(key string) (string, bool) {
	v, ok := s.entries[key]
	return v, ok
}

// Delete removes a key from the set.
func (s *Set) Delete(key string) {
	delete(s.entries, key)
	for i, k := range s.keys {
		if k == key {
			s.keys = append(s.keys[:i], s.keys[i+1:]...)
			break
		}
	}
}

// Keys returns all keys in insertion order.
func (s *Set) Keys() []string {
	out := make([]string, len(s.keys))
	copy(out, s.keys)
	return out
}

// SortedKeys returns all keys in alphabetical order.
func (s *Set) SortedKeys() []string {
	out := s.Keys()
	sort.Strings(out)
	return out
}

// Entries returns all entries in insertion order.
func (s *Set) Entries() []Entry {
	out := make([]Entry, 0, len(s.keys))
	for _, k := range s.keys {
		out = append(out, Entry{Key: k, Value: s.entries[k]})
	}
	return out
}

// Len returns the number of entries.
func (s *Set) Len() int { return len(s.keys) }

// String serialises the set as dotenv lines.
func (s *Set) String() string {
	var sb strings.Builder
	for _, k := range s.keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, s.entries[k])
	}
	return sb.String()
}
