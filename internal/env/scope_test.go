package env

import (
	"testing"
)

func baseScopeSet() *Set {
	s := NewSet()
	s.Set("DB_HOST", "localhost")
	s.Set("DB_PORT", "5432")
	s.Set("AWS_KEY", "abc")
	s.Set("AWS_SECRET", "xyz")
	s.Set("APP_NAME", "envdiff")
	s.Set("PORT", "8080")
	return s
}

func TestSplitByScope_BasicPartition(t *testing.T) {
	s := baseScopeSet()
	scopes := SplitByScope(s, []string{"DB_", "AWS_"})

	if len(scopes) != 3 {
		t.Fatalf("expected 3 scopes (DB_, AWS_, default), got %d", len(scopes))
	}
	if scopes[0].Name != "DB_" {
		t.Errorf("expected first scope DB_, got %s", scopes[0].Name)
	}
	if scopes[1].Name != "AWS_" {
		t.Errorf("expected second scope AWS_, got %s", scopes[1].Name)
	}
	if scopes[2].Name != "default" {
		t.Errorf("expected third scope default, got %s", scopes[2].Name)
	}
}

func TestSplitByScope_CorrectKeys(t *testing.T) {
	s := baseScopeSet()
	scopes := SplitByScope(s, []string{"DB_", "AWS_"})

	dbScope := scopes[0]
	if v, _ := dbScope.Entries.Get("DB_HOST"); v != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", v)
	}
	if v, _ := dbScope.Entries.Get("DB_PORT"); v != "5432" {
		t.Errorf("expected DB_PORT=5432, got %s", v)
	}
}

func TestSplitByScope_DefaultBucket(t *testing.T) {
	s := baseScopeSet()
	scopes := SplitByScope(s, []string{"DB_", "AWS_"})
	def := scopes[2]

	if _, ok := def.Entries.Get("PORT"); !ok {
		t.Error("expected PORT in default scope")
	}
	if _, ok := def.Entries.Get("APP_NAME"); !ok {
		t.Error("expected APP_NAME in default scope")
	}
}

func TestSplitByScope_NoMatchingPrefix(t *testing.T) {
	s := NewSet()
	s.Set("FOO", "bar")
	scopes := SplitByScope(s, []string{"DB_"})

	if len(scopes) != 1 || scopes[0].Name != "default" {
		t.Errorf("expected single default scope, got %+v", scopes)
	}
}

func TestMergeScopes_RoundTrip(t *testing.T) {
	s := baseScopeSet()
	scopes := SplitByScope(s, []string{"DB_", "AWS_"})
	merged := MergeScopes(scopes)

	for _, key := range s.Keys() {
		origVal, _ := s.Get(key)
		mergedVal, ok := merged.Get(key)
		if !ok {
			t.Errorf("key %s missing after merge", key)
		}
		if origVal != mergedVal {
			t.Errorf("key %s: expected %s, got %s", key, origVal, mergedVal)
		}
	}
}

func TestMergeScopes_LaterOverwrites(t *testing.T) {
	s1 := NewSet()
	s1.Set("KEY", "first")
	s2 := NewSet()
	s2.Set("KEY", "second")

	scopes := []Scope{{Name: "a", Entries: s1}, {Name: "b", Entries: s2}}
	merged := MergeScopes(scopes)

	if v, _ := merged.Get("KEY"); v != "second" {
		t.Errorf("expected 'second', got %s", v)
	}
}
