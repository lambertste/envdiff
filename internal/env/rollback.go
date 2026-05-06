package env

import (
	"fmt"
	"strings"
)

// RollbackEntry records a single key's previous state for rollback purposes.
type RollbackEntry struct {
	Key      string
	OldValue string
	HadKey   bool
}

// RollbackPlan holds all the changes needed to revert a set of mutations.
type RollbackPlan struct {
	Entries []RollbackEntry
}

// Rollback reverts the given Set to its state before the mutations described
// by the plan. Keys that did not exist before are deleted; others are restored.
func Rollback(s *Set, plan RollbackPlan) *Set {
	out := Clone(s)
	for _, e := range plan.Entries {
		if !e.HadKey {
			out.Delete(e.Key)
		} else {
			out.Set(e.Key, e.OldValue)
		}
	}
	return out
}

// SnapshotKeys captures the current values of the given keys from s so that a
// RollbackPlan can be constructed before any mutations are applied.
func SnapshotKeys(s *Set, keys []string) RollbackPlan {
	plan := RollbackPlan{}
	for _, k := range keys {
		v, ok := s.Get(k)
		plan.Entries = append(plan.Entries, RollbackEntry{
			Key:      k,
			OldValue: v,
			HadKey:   ok,
		})
	}
	return plan
}

// FormatRollback returns a human-readable summary of what would be rolled back.
func FormatRollback(plan RollbackPlan) string {
	if len(plan.Entries) == 0 {
		return "(no rollback entries)"
	}
	var sb strings.Builder
	for _, e := range plan.Entries {
		if !e.HadKey {
			fmt.Fprintf(&sb, "  delete  %s\n", e.Key)
		} else {
			fmt.Fprintf(&sb, "  restore %s=%s\n", e.Key, e.OldValue)
		}
	}
	return sb.String()
}
