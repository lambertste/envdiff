package env

import (
	"strings"
	"testing"
)

func baseRollbackSet() *Set {
	s := NewSet()
	s.Set("HOST", "localhost")
	s.Set("PORT", "5432")
	s.Set("DEBUG", "false")
	return s
}

func TestRollback_RestoresModifiedKey(t *testing.T) {
	s := baseRollbackSet()
	plan := SnapshotKeys(s, []string{"PORT"})
	s.Set("PORT", "9999")
	out := Rollback(s, plan)
	v, _ := out.Get("PORT")
	if v != "5432" {
		t.Errorf("expected 5432, got %s", v)
	}
}

func TestRollback_DeletesNewKey(t *testing.T) {
	s := baseRollbackSet()
	plan := SnapshotKeys(s, []string{"NEWKEY"})
	s.Set("NEWKEY", "added")
	out := Rollback(s, plan)
	_, ok := out.Get("NEWKEY")
	if ok {
		t.Error("expected NEWKEY to be deleted after rollback")
	}
}

func TestRollback_PreservesUntouchedKeys(t *testing.T) {
	s := baseRollbackSet()
	plan := SnapshotKeys(s, []string{"PORT"})
	s.Set("PORT", "9999")
	out := Rollback(s, plan)
	v, _ := out.Get("HOST")
	if v != "localhost" {
		t.Errorf("expected HOST=localhost, got %s", v)
	}
}

func TestSnapshotKeys_CapturesHadKey(t *testing.T) {
	s := baseRollbackSet()
	plan := SnapshotKeys(s, []string{"HOST", "MISSING"})
	if len(plan.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(plan.Entries))
	}
	if !plan.Entries[0].HadKey {
		t.Error("HOST should have HadKey=true")
	}
	if plan.Entries[1].HadKey {
		t.Error("MISSING should have HadKey=false")
	}
}

func TestFormatRollback_Empty(t *testing.T) {
	out := FormatRollback(RollbackPlan{})
	if !strings.Contains(out, "no rollback") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatRollback_ShowsActions(t *testing.T) {
	s := baseRollbackSet()
	plan := SnapshotKeys(s, []string{"PORT", "NEWKEY"})
	out := FormatRollback(plan)
	if !strings.Contains(out, "restore PORT") {
		t.Errorf("expected restore PORT in output: %s", out)
	}
	if !strings.Contains(out, "delete  NEWKEY") {
		t.Errorf("expected delete NEWKEY in output: %s", out)
	}
}
