package snapshot_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/envdiff/internal/snapshot"
)

func makeSnap(label string, entries map[string]string) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		Timestamp: time.Now(),
		Label:     label,
		Entries:   entries,
	}
}

func TestCompare_NoChanges(t *testing.T) {
	base := makeSnap("base", map[string]string{"A": "1", "B": "2"})
	other := makeSnap("other", map[string]string{"A": "1", "B": "2"})
	r := snapshot.Compare(base, other)
	if len(r.Added) != 0 || len(r.Removed) != 0 || len(r.Modified) != 0 {
		t.Errorf("expected no changes, got %+v", r)
	}
	if len(r.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged, got %d", len(r.Unchanged))
	}
}

func TestCompare_Added(t *testing.T) {
	base := makeSnap("base", map[string]string{"A": "1"})
	other := makeSnap("other", map[string]string{"A": "1", "NEW": "val"})
	r := snapshot.Compare(base, other)
	if len(r.Added) != 1 || r.Added[0] != "NEW" {
		t.Errorf("expected Added=[NEW], got %v", r.Added)
	}
}

func TestCompare_Removed(t *testing.T) {
	base := makeSnap("base", map[string]string{"A": "1", "OLD": "gone"})
	other := makeSnap("other", map[string]string{"A": "1"})
	r := snapshot.Compare(base, other)
	if len(r.Removed) != 1 || r.Removed[0] != "OLD" {
		t.Errorf("expected Removed=[OLD], got %v", r.Removed)
	}
}

func TestCompare_Modified(t *testing.T) {
	base := makeSnap("base", map[string]string{"A": "old"})
	other := makeSnap("other", map[string]string{"A": "new"})
	r := snapshot.Compare(base, other)
	if len(r.Modified) != 1 || r.Modified[0] != "A" {
		t.Errorf("expected Modified=[A], got %v", r.Modified)
	}
}

func TestFormatDiff_ContainsLabels(t *testing.T) {
	base := makeSnap("staging", map[string]string{"X": "1"})
	other := makeSnap("prod", map[string]string{"X": "2", "Y": "3"})
	r := snapshot.Compare(base, other)
	out := snapshot.FormatDiff(base, other, r)
	if !strings.Contains(out, "staging") {
		t.Error("expected 'staging' in output")
	}
	if !strings.Contains(out, "prod") {
		t.Error("expected 'prod' in output")
	}
	if !strings.Contains(out, "+ Y") {
		t.Error("expected added key Y in output")
	}
	if !strings.Contains(out, "~ X") {
		t.Error("expected modified key X in output")
	}
}
