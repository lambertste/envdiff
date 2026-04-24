package diff

import "sort"

// Status represents the change status of an environment variable.
type Status int

const (
	Unchanged Status = iota
	Added
	Removed
	Modified
)

// Entry represents a single key comparison result between two env maps.
type Entry struct {
	Key      string
	Status   Status
	OldValue string
	NewValue string
}

// Diff compares two env maps (source, target) and returns a slice of Entry
// describing keys that were added, removed, or modified in target relative to source.
// Unchanged keys are omitted. Results are sorted by key.
func Diff(source, target map[string]string) []Entry {
	seen := make(map[string]bool)
	var entries []Entry

	for k, sv := range source {
		seen[k] = true
		if tv, ok := target[k]; !ok {
			entries = append(entries, Entry{Key: k, Status: Removed, OldValue: sv})
		} else if sv != tv {
			entries = append(entries, Entry{Key: k, Status: Modified, OldValue: sv, NewValue: tv})
		}
	}

	for k, tv := range target {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, Status: Added, NewValue: tv})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return entries
}
