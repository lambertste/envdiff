package env

import (
	"strings"
	"testing"
)

func makePivotSet(t *testing.T, pairs ...string) *Set {
	t.Helper()
	s := NewSet()
	for i := 0; i+1 < len(pairs); i += 2 {
		s.Set(pairs[i], pairs[i+1])
	}
	return s
}

func TestPivot_BasicTwoColumns(t *testing.T) {
	staging := makePivotSet(t, "ENV", "staging", "DB_HOST", "db1", "API_KEY", "abc")
	prod := makePivotSet(t, "ENV", "prod", "DB_HOST", "db2", "API_KEY", "xyz")

	pr, err := Pivot([]*Set{staging, prod}, "ENV")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pr.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(pr.Columns))
	}
	if pr.Rows["DB_HOST"]["staging"] != "db1" {
		t.Errorf("expected db1, got %q", pr.Rows["DB_HOST"]["staging"])
	}
	if pr.Rows["DB_HOST"]["prod"] != "db2" {
		t.Errorf("expected db2, got %q", pr.Rows["DB_HOST"]["prod"])
	}
	if _, ok := pr.Rows["ENV"]; ok {
		t.Error("pivot key should not appear in rows")
	}
}

func TestPivot_MissingPivotKey_ReturnsError(t *testing.T) {
	s := makePivotSet(t, "DB_HOST", "db1")
	_, err := Pivot([]*Set{s}, "ENV")
	if err == nil {
		t.Fatal("expected error for missing pivot key")
	}
}

func TestPivot_EmptySlice(t *testing.T) {
	pr, err := Pivot([]*Set{}, "ENV")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pr.Rows) != 0 {
		t.Errorf("expected empty rows")
	}
}

func TestPivot_MissingValueInOneColumn(t *testing.T) {
	staging := makePivotSet(t, "ENV", "staging", "DB_HOST", "db1")
	prod := makePivotSet(t, "ENV", "prod", "DB_HOST", "db2", "EXTRA", "only-prod")

	pr, err := Pivot([]*Set{staging, prod}, "ENV")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pr.Rows["EXTRA"]["staging"] != "" {
		t.Errorf("expected empty string for missing column value")
	}
	if pr.Rows["EXTRA"]["prod"] != "only-prod" {
		t.Errorf("expected only-prod")
	}
}

func TestFormatPivot_ContainsColumnsAndKeys(t *testing.T) {
	staging := makePivotSet(t, "ENV", "staging", "DB_HOST", "db1")
	prod := makePivotSet(t, "ENV", "prod", "DB_HOST", "db2")

	pr, _ := Pivot([]*Set{staging, prod}, "ENV")
	out := FormatPivot(pr)

	for _, want := range []string{"staging", "prod", "DB_HOST", "db1", "db2"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestFormatPivot_Empty(t *testing.T) {
	pr := &PivotResult{Rows: make(map[string]map[string]string)}
	out := FormatPivot(pr)
	if !strings.Contains(out, "empty") {
		t.Errorf("expected empty marker, got %q", out)
	}
}
