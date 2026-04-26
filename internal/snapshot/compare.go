package snapshot

import (
	"fmt"
	"sort"
	"strings"
)

// DiffResult holds the result of comparing two snapshots.
type DiffResult struct {
	Added    []string
	Removed  []string
	Modified []string
	Unchanged []string
}

// Compare diffs two snapshots and returns a DiffResult.
func Compare(base, other *Snapshot) DiffResult {
	var result DiffResult

	for k, v := range other.Entries {
		baseVal, exists := base.Entries[k]
		if !exists {
			result.Added = append(result.Added, k)
		} else if baseVal != v {
			result.Modified = append(result.Modified, k)
		} else {
			result.Unchanged = append(result.Unchanged, k)
		}
	}

	for k := range base.Entries {
		if _, exists := other.Entries[k]; !exists {
			result.Removed = append(result.Removed, k)
		}
	}

	sort.Strings(result.Added)
	sort.Strings(result.Removed)
	sort.Strings(result.Modified)
	sort.Strings(result.Unchanged)

	return result
}

// Format renders a DiffResult as a human-readable string.
func FormatDiff(base, other *Snapshot, r DiffResult) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Snapshot diff: %q → %q\n", base.Label, other.Label)
	fmt.Fprintf(&sb, "  Added:    %d\n", len(r.Added))
	fmt.Fprintf(&sb, "  Removed:  %d\n", len(r.Removed))
	fmt.Fprintf(&sb, "  Modified: %d\n", len(r.Modified))
	for _, k := range r.Added {
		fmt.Fprintf(&sb, "  + %s=%s\n", k, other.Entries[k])
	}
	for _, k := range r.Removed {
		fmt.Fprintf(&sb, "  - %s=%s\n", k, base.Entries[k])
	}
	for _, k := range r.Modified {
		fmt.Fprintf(&sb, "  ~ %s: %q → %q\n", k, base.Entries[k], other.Entries[k])
	}
	return sb.String()
}
