package reconcile_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/reconcile"
)

func entries(items ...diff.Entry) []diff.Entry { return items }

func TestPlan_Empty(t *testing.T) {
	steps := reconcile.Plan(entries())
	if len(steps) != 0 {
		t.Fatalf("expected 0 steps, got %d", len(steps))
	}
}

func TestPlan_AddedEntry(t *testing.T) {
	e := diff.Entry{Key: "NEW_KEY", Status: diff.Added, NewValue: "value1"}
	steps := reconcile.Plan(entries(e))
	if len(steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(steps))
	}
	if steps[0].Action != reconcile.ActionAdd || steps[0].Key != "NEW_KEY" || steps[0].Value != "value1" {
		t.Errorf("unexpected step: %+v", steps[0])
	}
}

func TestPlan_RemovedEntry(t *testing.T) {
	e := diff.Entry{Key: "OLD_KEY", Status: diff.Removed, OldValue: "gone"}
	steps := reconcile.Plan(entries(e))
	if steps[0].Action != reconcile.ActionRemove || steps[0].Key != "OLD_KEY" {
		t.Errorf("unexpected step: %+v", steps[0])
	}
}

func TestPlan_ModifiedEntry(t *testing.T) {
	e := diff.Entry{Key: "DB_HOST", Status: diff.Modified, OldValue: "localhost", NewValue: "prod-db"}
	steps := reconcile.Plan(entries(e))
	if steps[0].Action != reconcile.ActionUpdate || steps[0].Value != "prod-db" {
		t.Errorf("unexpected step: %+v", steps[0])
	}
}

func TestPlan_SortedOutput(t *testing.T) {
	steps := reconcile.Plan(entries(
		diff.Entry{Key: "Z_KEY", Status: diff.Added, NewValue: "z"},
		diff.Entry{Key: "A_KEY", Status: diff.Added, NewValue: "a"},
		diff.Entry{Key: "M_KEY", Status: diff.Removed},
	))
	keys := []string{steps[0].Key, steps[1].Key, steps[2].Key}
	expected := []string{"A_KEY", "M_KEY", "Z_KEY"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("position %d: want %s, got %s", i, expected[i], k)
		}
	}
}

func TestApply_AddAndUpdate(t *testing.T) {
	base := map[string]string{"EXISTING": "old"}
	steps := []reconcile.Step{
		{Action: reconcile.ActionAdd, Key: "NEW", Value: "new_val"},
		{Action: reconcile.ActionUpdate, Key: "EXISTING", Value: "updated"},
	}
	result := reconcile.Apply(base, steps)
	if result["NEW"] != "new_val" {
		t.Errorf("expected NEW=new_val, got %s", result["NEW"])
	}
	if result["EXISTING"] != "updated" {
		t.Errorf("expected EXISTING=updated, got %s", result["EXISTING"])
	}
}

func TestApply_Remove(t *testing.T) {
	base := map[string]string{"TO_REMOVE": "val", "KEEP": "keep"}
	steps := []reconcile.Step{{Action: reconcile.ActionRemove, Key: "TO_REMOVE"}}
	result := reconcile.Apply(base, steps)
	if _, ok := result["TO_REMOVE"]; ok {
		t.Error("expected TO_REMOVE to be deleted")
	}
	if result["KEEP"] != "keep" {
		t.Error("expected KEEP to remain")
	}
}

func TestFormat(t *testing.T) {
	steps := []reconcile.Step{
		{Action: reconcile.ActionAdd, Key: "A", Value: "1"},
		{Action: reconcile.ActionRemove, Key: "B"},
		{Action: reconcile.ActionUpdate, Key: "C", Value: "3"},
	}
	out := reconcile.Format(steps)
	expected := "+ A=1\n- B\n~ C=3"
	if out != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, out)
	}
}
