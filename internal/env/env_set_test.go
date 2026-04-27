package env

import (
	"strings"
	"testing"
)

func TestNewSet_Empty(t *testing.T) {
	s := NewSet()
	if s.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", s.Len())
	}
}

func TestSet_SetAndGet(t *testing.T) {
	s := NewSet()
	s.Set("FOO", "bar")
	v, ok := s.Get("FOO")
	if !ok || v != "bar" {
		t.Fatalf("expected bar, got %q ok=%v", v, ok)
	}
}

func TestSet_OverwriteValue(t *testing.T) {
	s := NewSet()
	s.Set("KEY", "old")
	s.Set("KEY", "new")
	v, _ := s.Get("KEY")
	if v != "new" {
		t.Fatalf("expected new, got %q", v)
	}
	if s.Len() != 1 {
		t.Fatalf("expected 1 key, got %d", s.Len())
	}
}

func TestSet_Delete(t *testing.T) {
	s := NewSet()
	s.Set("A", "1")
	s.Set("B", "2")
	s.Delete("A")
	if _, ok := s.Get("A"); ok {
		t.Fatal("expected A to be deleted")
	}
	if s.Len() != 1 {
		t.Fatalf("expected 1 key after delete, got %d", s.Len())
	}
}

func TestSet_InsertionOrder(t *testing.T) {
	s := NewSet()
	s.Set("Z", "z")
	s.Set("A", "a")
	s.Set("M", "m")
	keys := s.Keys()
	expected := []string{"Z", "A", "M"}
	for i, k := range expected {
		if keys[i] != k {
			t.Fatalf("position %d: want %s got %s", i, k, keys[i])
		}
	}
}

func TestSet_SortedKeys(t *testing.T) {
	s := NewSet()
	s.Set("Z", "z")
	s.Set("A", "a")
	s.Set("M", "m")
	keys := s.SortedKeys()
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Fatalf("unexpected sorted order: %v", keys)
	}
}

func TestFromEntries(t *testing.T) {
	entries := []Entry{{Key: "X", Value: "1"}, {Key: "Y", Value: "2"}}
	s := FromEntries(entries)
	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}
}

func TestSet_String(t *testing.T) {
	s := NewSet()
	s.Set("FOO", "bar")
	s.Set("BAZ", "qux")
	out := s.String()
	if !strings.Contains(out, "FOO=bar") || !strings.Contains(out, "BAZ=qux") {
		t.Fatalf("unexpected string output: %q", out)
	}
}
