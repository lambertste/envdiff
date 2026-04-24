package merge

import (
	"fmt"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// Strategy defines how conflicting keys are resolved during a merge.
type Strategy int

const (
	// PreferBase keeps the base value on conflict.
	PreferBase Strategy = iota
	// PreferOverride uses the override value on conflict.
	PreferOverride
)

// Result holds the merged environment and any conflicts encountered.
type Result struct {
	Merged    map[string]string
	Conflicts []Conflict
}

// Conflict describes a key that existed in both sources with different values.
type Conflict struct {
	Key       string
	BaseValue string
	OverValue string
	Resolved  string
}

// Merge combines base and override env maps using the given strategy.
// Keys only in base or only in override are always included.
// Conflicting keys are resolved according to strategy and recorded.
func Merge(base, override map[string]string, strategy Strategy) Result {
	merged := make(map[string]string, len(base))
	var conflicts []Conflict

	for k, v := range base {
		merged[k] = v
	}

	for k, ov := range override {
		bv, exists := merged[k]
		if !exists {
			merged[k] = ov
			continue
		}
		if bv == ov {
			continue
		}
		// Conflict detected
		resolved := bv
		if strategy == PreferOverride {
			resolved = ov
		}
		merged[k] = resolved
		conflicts = append(conflicts, Conflict{
			Key:       k,
			BaseValue: bv,
			OverValue: ov,
			Resolved:  resolved,
		})
	}

	sort.Slice(conflicts, func(i, j int) bool {
		return conflicts[i].Key < conflicts[j].Key
	})

	return Result{Merged: merged, Conflicts: conflicts}
}

// FormatConflicts returns a human-readable summary of merge conflicts.
func FormatConflicts(conflicts []Conflict) string {
	if len(conflicts) == 0 {
		return ""
	}
	out := fmt.Sprintf("%d conflict(s) resolved:\n", len(conflicts))
	for _, c := range conflicts {
		out += fmt.Sprintf("  %s: base=%q override=%q -> resolved=%q\n",
			c.Key, c.BaseValue, c.OverValue, c.Resolved)
	}
	return out
}

// ToEntries converts a merged map to a sorted slice of diff.Entry with kind Unchanged.
func ToEntries(merged map[string]string) []diff.Entry {
	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]diff.Entry, 0, len(merged))
	for _, k := range keys {
		entries = append(entries, diff.Entry{
			Key:      k,
			NewValue: merged[k],
			Kind:     diff.Unchanged,
		})
	}
	return entries
}
