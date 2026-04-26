package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved state of an environment file at a point in time.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Label     string            `json:"label"`
	Entries   map[string]string `json:"entries"`
}

// Save writes a snapshot to the given file path as JSON.
func Save(path string, label string, entries map[string]string) error {
	s := Snapshot{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Entries:   entries,
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}
	return &s, nil
}

// ToEntries converts the snapshot's entries map to a slice of key=value strings.
func (s *Snapshot) ToEntries() []string {
	out := make([]string, 0, len(s.Entries))
	for k, v := range s.Entries {
		out = append(out, k+"="+v)
	}
	return out
}
