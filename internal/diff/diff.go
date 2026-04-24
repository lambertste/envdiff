package diff

import (
	"sort"

	"github.com/user/envdiff/internal/parser"
)

// ChangeKind classifies the type of difference between two env files.
type ChangeKind string

const (
	Added    ChangeKind = "added"    // key exists in right but not left
	Removed  ChangeKind = "removed"  // key exists in left but not right
	Modified ChangeKind = "modified" // key exists in both but values differ
)

// Change represents a single diffed entry.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string // empty for Added
	NewValue string // empty for Removed
}

// Result holds the full diff output between two env maps.
type Result struct {
	Changes []Change
}

// HasChanges returns true if there are any differences.
func (r *Result) HasChanges() bool {
	return len(r.Changes) > 0
}

// Diff computes the difference between a base (left) and target (right) EnvMap.
func Diff(left, right parser.EnvMap) *Result {
	result := &Result{}

	// Detect removed and modified keys.
	for k, lv := range left {
		if rv, ok := right[k]; !ok {
			result.Changes = append(result.Changes, Change{
				Key:      k,
				Kind:     Removed,
				OldValue: lv,
			})
		} else if lv != rv {
			result.Changes = append(result.Changes, Change{
				Key:      k,
				Kind:     Modified,
				OldValue: lv,
				NewValue: rv,
			})
		}
	}

	// Detect added keys.
	for k, rv := range right {
		if _, ok := left[k]; !ok {
			result.Changes = append(result.Changes, Change{
				Key:      k,
				Kind:     Added,
				NewValue: rv,
			})
		}
	}

	// Sort changes for deterministic output.
	sort.Slice(result.Changes, func(i, j int) bool {
		return result.Changes[i].Key < result.Changes[j].Key
	})

	return result
}
