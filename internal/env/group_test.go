package env

import (
	"strings"
	"testing"
)

func baseGroupSet() *Set {
	s := NewSet()
	s.Set("DB_HOST", "localhost")
	s.Set("DB_PORT", "5432")
	s.Set("APP_NAME", "envdiff")
	s.Set("APP_ENV", "production")
	s.Set("LOG_LEVEL", "info")
	return s
}

func prefixGroupFn(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return ""
}

func TestGroupBy_BasicPartition(t *testing.T) {
	s := baseGroupSet()
	gr := GroupBy(s, prefixGroupFn)

	if len(gr) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(gr))
	}
	if _, ok := gr["DB"]; !ok {
		t.Error("expected group DB")
	}
	if _, ok := gr["APP"]; !ok {
		t.Error("expected group APP")
	}
	if _, ok := gr["LOG"]; !ok {
		t.Error("expected group LOG")
	}
}

func TestGroupBy_CorrectKeys(t *testing.T) {
	s := baseGroupSet()
	gr := GroupBy(s, prefixGroupFn)

	dbKeys := gr["DB"].Keys()
	if len(dbKeys) != 2 {
		t.Fatalf("expected 2 DB keys, got %d", len(dbKeys))
	}
}

func TestGroupBy_DefaultBucket(t *testing.T) {
	s := NewSet()
	s.Set("NOPREFIX", "val")
	gr := GroupBy(s, prefixGroupFn)

	if _, ok := gr["_default"]; !ok {
		t.Error("expected _default bucket for key without prefix")
	}
}

func TestGroupNames_Sorted(t *testing.T) {
	s := baseGroupSet()
	gr := GroupBy(s, prefixGroupFn)
	names := GroupNames(gr)

	for i := 1; i < len(names); i++ {
		if names[i] < names[i-1] {
			t.Errorf("group names not sorted: %v", names)
		}
	}
}

func TestMergeGroups_ReturnsAllKeys(t *testing.T) {
	s := baseGroupSet()
	gr := GroupBy(s, prefixGroupFn)
	merged := MergeGroups(gr)

	for _, k := range s.Keys() {
		if _, ok := merged.Get(k); !ok {
			t.Errorf("key %q missing from merged set", k)
		}
	}
}

func TestMergeGroups_EmptyInput(t *testing.T) {
	merged := MergeGroups(GroupResult{})
	if len(merged.Keys()) != 0 {
		t.Error("expected empty set from empty GroupResult")
	}
}
