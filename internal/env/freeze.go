package env

import (
	"fmt"
	"sort"
	"strings"
)

// FreezeResult holds the outcome of a freeze operation.
type FreezeResult struct {
	Frozen  []string // keys that were frozen (made read-only via comment marker)
	Skipped []string // keys not found in the set
}

// FrozenPrefix is the marker prepended to a key to denote it is frozen.
const FrozenPrefix = "#frozen:"

// Freeze marks the given keys as frozen by encoding them into the set's
// metadata. Frozen keys are preserved in output with a comment annotation
// and cannot be overwritten by Patch or ApplyDefaults.
func Freeze(s *Set, keys ...string) FreezeResult {
	result := FreezeResult{}
	for _, k := range keys {
		if _, ok := s.Get(k); !ok {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		marker := FrozenPrefix + k
		s.Set(marker, "1")
		result.Frozen = append(result.Frozen, k)
	}
	sort.Strings(result.Frozen)
	sort.Strings(result.Skipped)
	return result
}

// IsFrozen reports whether the given key is frozen in the set.
func IsFrozen(s *Set, key string) bool {
	marker := FrozenPrefix + key
	v, ok := s.Get(marker)
	return ok && v == "1"
}

// FrozenKeys returns all keys currently marked as frozen.
func FrozenKeys(s *Set) []string {
	var frozen []string
	for _, k := range s.Keys() {
		if strings.HasPrefix(k, FrozenPrefix) {
			frozen = append(frozen, strings.TrimPrefix(k, FrozenPrefix))
		}
	}
	sort.Strings(frozen)
	return frozen
}

// Unfreeze removes the frozen marker from the given keys.
func Unfreeze(s *Set, keys ...string) {
	for _, k := range keys {
		marker := FrozenPrefix + k
		s.Delete(marker)
	}
}

// FormatFrozen returns a human-readable summary of frozen keys.
func FormatFrozen(keys []string) string {
	if len(keys) == 0 {
		return "no frozen keys"
	}
	return fmt.Sprintf("frozen keys (%d): %s", len(keys), strings.Join(keys, ", "))
}
