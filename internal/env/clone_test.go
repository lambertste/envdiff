package env

import (
	"testing"
)

func baseCloneSet() *Set {
	s := NewSet()
	s.Set("APP_NAME", "envdiff")
	s.Set("APP_ENV", "production")
	s.Set("DB_HOST", "localhost")
	s.Set("DB_PASS", "secret")
	return s
}

func TestClone_AllKeys(t *testing.T) {
	src := baseCloneSet()
	dst := Clone(src)
	for _, k := range src.Keys() {
		v1, _ := src.Get(k)
		v2, ok := dst.Get(k)
		if !ok || v1 != v2 {
			t.Errorf("key %q: expected %q, got %q", k, v1, v2)
		}
	}
}

func TestClone_IsIndependent(t *testing.T) {
	src := baseCloneSet()
	dst := Clone(src)
	dst.Set("APP_NAME", "changed")
	v, _ := src.Get("APP_NAME")
	if v == "changed" {
		t.Error("clone mutation should not affect source")
	}
}

func TestClone_WithKeys(t *testing.T) {
	src := baseCloneSet()
	dst := Clone(src, WithKeys("APP_NAME", "APP_ENV"))
	if dst.Len() != 2 {
		t.Fatalf("expected 2 keys, got %d", dst.Len())
	}
	if _, ok := dst.Get("DB_HOST"); ok {
		t.Error("DB_HOST should not be present")
	}
}

func TestClone_WithoutKeys(t *testing.T) {
	src := baseCloneSet()
	dst := Clone(src, WithoutKeys("DB_PASS"))
	if _, ok := dst.Get("DB_PASS"); ok {
		t.Error("DB_PASS should have been excluded")
	}
	if dst.Len() != src.Len()-1 {
		t.Errorf("expected %d keys, got %d", src.Len()-1, dst.Len())
	}
}

func TestClone_WithKeysAndWithoutKeys(t *testing.T) {
	src := baseCloneSet()
	// WithKeys takes precedence; WithoutKeys is applied after
	dst := Clone(src, WithKeys("APP_NAME", "DB_PASS"), WithoutKeys("DB_PASS"))
	if dst.Len() != 1 {
		t.Fatalf("expected 1 key, got %d", dst.Len())
	}
	if _, ok := dst.Get("APP_NAME"); !ok {
		t.Error("APP_NAME should be present")
	}
}

func TestMergeInto_OverwritesExisting(t *testing.T) {
	dst := NewSet()
	dst.Set("KEY", "old")
	src := NewSet()
	src.Set("KEY", "new")
	src.Set("OTHER", "value")
	MergeInto(dst, src)
	if v, _ := dst.Get("KEY"); v != "new" {
		t.Errorf("expected \"new\", got %q", v)
	}
	if _, ok := dst.Get("OTHER"); !ok {
		t.Error("OTHER should have been merged")
	}
}
