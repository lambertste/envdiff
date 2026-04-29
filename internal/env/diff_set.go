package env

// DiffKind represents the type of change between two sets.
type DiffKind string

const (
	DiffAdded    DiffKind = "added"
	DiffRemoved  DiffKind = "removed"
	DiffModified DiffKind = "modified"
	DiffUnchanged DiffKind = "unchanged"
)

// DiffEntry holds a key and its change status between two Sets.
type DiffEntry struct {
	Key      string
	Kind     DiffKind
	OldValue string
	NewValue string
}

// DiffSets compares base and other, returning a slice of DiffEntry
// describing every key present in either set.
func DiffSets(base, other *Set) []DiffEntry {
	seen := make(map[string]bool)
	var results []DiffEntry

	// Walk keys in insertion order of base first.
	for _, k := range base.Keys() {
		seen[k] = true
		baseVal, _ := base.Get(k)
		otherVal, otherOK := other.Get(k)
		switch {
		case !otherOK:
			results = append(results, DiffEntry{Key: k, Kind: DiffRemoved, OldValue: baseVal})
		case baseVal != otherVal:
			results = append(results, DiffEntry{Key: k, Kind: DiffModified, OldValue: baseVal, NewValue: otherVal})
		default:
			results = append(results, DiffEntry{Key: k, Kind: DiffUnchanged, OldValue: baseVal, NewValue: otherVal})
		}
	}

	// Keys only in other.
	for _, k := range other.Keys() {
		if seen[k] {
			continue
		}
		v, _ := other.Get(k)
		results = append(results, DiffEntry{Key: k, Kind: DiffAdded, NewValue: v})
	}

	return results
}

// FilterDiff returns only the DiffEntry items matching any of the given kinds.
func FilterDiff(entries []DiffEntry, kinds ...DiffKind) []DiffEntry {
	kindSet := make(map[DiffKind]bool, len(kinds))
	for _, k := range kinds {
		kindSet[k] = true
	}
	var out []DiffEntry
	for _, e := range entries {
		if kindSet[e.Kind] {
			out = append(out, e)
		}
	}
	return out
}
