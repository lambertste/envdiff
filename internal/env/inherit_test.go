package env

import (
	"strings"
	"testing"
)

func baseInheritChild() *Set {
	s := NewSet()
	s.Set("APP_NAME", "myapp")
	s.Set("LOG_LEVEL", "debug")
	return s
}

func baseInheritParent() *Set {
	s := NewSet()
	s.Set("APP_NAME", "parent-app")
	s.Set("DB_HOST", "db.prod.internal")
	s.Set("SECRET", "")
	s.Set("TIMEOUT", "30s")
	return s
}

func TestInherit_FillsMissingKeys(t *testing.T) {
	child := baseInheritChild()
	parent := baseInheritParent()
	out, res, err := Inherit(child, parent, DefaultInheritOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := out.Get("DB_HOST")
	if !ok || v != "db.prod.internal" {
		t.Errorf("expected DB_HOST=db.prod.internal, got %q", v)
	}
	if !containsStr(res.Inherited, "DB_HOST") {
		t.Errorf("expected DB_HOST in inherited list")
	}
}

func TestInherit_DoesNotOverwriteByDefault(t *testing.T) {
	child := baseInheritChild()
	parent := baseInheritParent()
	out, _, err := Inherit(child, parent, DefaultInheritOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := out.Get("APP_NAME")
	if v != "myapp" {
		t.Errorf("expected APP_NAME to remain 'myapp', got %q", v)
	}
}

func TestInherit_OverwriteExisting(t *testing.T) {
	child := baseInheritChild()
	parent := baseInheritParent()
	opts := DefaultInheritOptions()
	opts.OverwriteExisting = true
	out, _, err := Inherit(child, parent, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := out.Get("APP_NAME")
	if v != "parent-app" {
		t.Errorf("expected APP_NAME='parent-app', got %q", v)
	}
}

func TestInherit_SkipsEmptyParentValues(t *testing.T) {
	child := baseInheritChild()
	parent := baseInheritParent()
	_, res, err := Inherit(child, parent, DefaultInheritOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsStr(res.Skipped, "SECRET") {
		t.Errorf("expected SECRET skipped due to empty value")
	}
}

func TestInherit_NilChildReturnsError(t *testing.T) {
	_, _, err := Inherit(nil, baseInheritParent(), DefaultInheritOptions())
	if err == nil {
		t.Error("expected error for nil child")
	}
}

func TestInherit_NilParentReturnsError(t *testing.T) {
	_, _, err := Inherit(baseInheritChild(), nil, DefaultInheritOptions())
	if err == nil {
		t.Error("expected error for nil parent")
	}
}

func TestFormatInheritResult_ContainsKeys(t *testing.T) {
	r := InheritResult{
		Inherited: []string{"DB_HOST", "TIMEOUT"},
		Skipped:   []string{"APP_NAME"},
	}
	out := FormatInheritResult(r)
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output")
	}
	if !strings.Contains(out, "skipped") {
		t.Errorf("expected 'skipped' label in output")
	}
}

func containsStr(ss []string, target string) bool {
	for _, s := range ss {
		if s == target {
			return true
		}
	}
	return false
}
