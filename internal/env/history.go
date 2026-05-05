package env

import (
	"fmt"
	"strings"
	"time"
)

// ChangeKind describes what happened to an entry.
type ChangeKind string

const (
	ChangeSet    ChangeKind = "set"
	ChangeDelete ChangeKind = "delete"
	ChangeRename ChangeKind = "rename"
)

// HistoryEntry records a single mutation applied to a Set.
type HistoryEntry struct {
	At      time.Time
	Kind    ChangeKind
	Key     string
	OldVal  string
	NewVal  string
}

// History holds an ordered log of mutations.
type History struct {
	entries []HistoryEntry
}

// Record appends a new entry to the history.
func (h *History) Record(kind ChangeKind, key, oldVal, newVal string) {
	h.entries = append(h.entries, HistoryEntry{
		At:     time.Now().UTC(),
		Kind:   kind,
		Key:    key,
		OldVal: oldVal,
		NewVal: newVal,
	})
}

// Entries returns a copy of all recorded history entries.
func (h *History) Entries() []HistoryEntry {
	out := make([]HistoryEntry, len(h.entries))
	copy(out, h.entries)
	return out
}

// Len returns the number of recorded entries.
func (h *History) Len() int { return len(h.entries) }

// Format returns a human-readable summary of the history.
func (h *History) Format() string {
	if len(h.entries) == 0 {
		return "(no history)"
	}
	var sb strings.Builder
	for _, e := range h.entries {
		ts := e.At.Format(time.RFC3339)
		switch e.Kind {
		case ChangeSet:
			if e.OldVal == "" {
				fmt.Fprintf(&sb, "[%s] set   %s = %q\n", ts, e.Key, e.NewVal)
			} else {
				fmt.Fprintf(&sb, "[%s] set   %s: %q -> %q\n", ts, e.Key, e.OldVal, e.NewVal)
			}
		case ChangeDelete:
			fmt.Fprintf(&sb, "[%s] del   %s (was %q)\n", ts, e.Key, e.OldVal)
		case ChangeRename:
			fmt.Fprintf(&sb, "[%s] rename %s -> %s\n", ts, e.OldVal, e.NewVal)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

// TrackSet records a set operation, capturing any previous value from s.
func TrackSet(h *History, s *Set, key, newVal string) {
	old, _ := s.Get(key)
	h.Record(ChangeSet, key, old, newVal)
	s.Set(key, newVal)
}

// TrackDelete records a delete operation.
func TrackDelete(h *History, s *Set, key string) {
	old, _ := s.Get(key)
	h.Record(ChangeDelete, key, old, "")
	s.Delete(key)
}
