package diff

import "sort"

// Kind describes the type of change for a key.
type Kind string

const (
	Added    Kind = "added"
	Removed  Kind = "removed"
	Modified Kind = "modified"
	Unchanged Kind = "unchanged"
)

// Entry represents a single diff result for one environment variable key.
type Entry struct {
	Key      string
	Kind     Kind
	OldValue string
	NewValue string
}

// Diff compares two maps of environment variables (base vs target) and returns
// a sorted slice of Entry describing the differences.
func Diff(base, target map[string]string) []Entry {
	seen := make(map[string]bool)
	var results []Entry

	for k, targetVal := range target {
		seen[k] = true
		baseVal, exists := base[k]
		if !exists {
			results = append(results, Entry{Key: k, Kind: Added, NewValue: targetVal})
		} else if baseVal != targetVal {
			results = append(results, Entry{Key: k, Kind: Modified, OldValue: baseVal, NewValue: targetVal})
		}
	}

	for k, baseVal := range base {
		if !seen[k] {
			results = append(results, Entry{Key: k, Kind: Removed, OldValue: baseVal})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	return results
}
