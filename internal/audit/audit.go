package audit

import (
	"fmt"
	"strings"
	"time"
)

// EventType classifies what kind of change occurred.
type EventType string

const (
	EventAdded    EventType = "added"
	EventRemoved  EventType = "removed"
	EventModified EventType = "modified"
)

// Event records a single change to an environment variable.
type Event struct {
	Timestamp time.Time
	File      string
	Key       string
	Type      EventType
	OldValue  string
	NewValue  string
}

// Log holds an ordered list of audit events.
type Log struct {
	Events []Event
}

// Record appends a new event to the log.
func (l *Log) Record(file, key string, t EventType, oldVal, newVal string) {
	l.Events = append(l.Events, Event{
		Timestamp: time.Now().UTC(),
		File:      file,
		Key:       key,
		Type:      t,
		OldValue:  oldVal,
		NewValue:  newVal,
	})
}

// Format returns a human-readable summary of all events.
func (l *Log) Format() string {
	if len(l.Events) == 0 {
		return "no changes recorded"
	}
	var sb strings.Builder
	for _, e := range l.Events {
		ts := e.Timestamp.Format(time.RFC3339)
		switch e.Type {
		case EventAdded:
			fmt.Fprintf(&sb, "[%s] %s: + %s=%q\n", ts, e.File, e.Key, e.NewValue)
		case EventRemoved:
			fmt.Fprintf(&sb, "[%s] %s: - %s=%q\n", ts, e.File, e.Key, e.OldValue)
		case EventModified:
			fmt.Fprintf(&sb, "[%s] %s: ~ %s %q -> %q\n", ts, e.File, e.Key, e.OldValue, e.NewValue)
		}
	}
	return sb.String()
}

// Summary returns counts per event type.
func (l *Log) Summary() map[EventType]int {
	counts := map[EventType]int{}
	for _, e := range l.Events {
		counts[e.Type]++
	}
	return counts
}
