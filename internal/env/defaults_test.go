package env

import (
	"sort"
	"testing"
)

func baseDefaultSet() *Set {
	s := NewSet()
	s.Set("APP_NAME", "myapp")
	s.Set("LOG_LEVEL", "")
	return s
}

func TestApplyDefaults_FillsMissing(t *testing.T) {
	s := baseDefaultSet()
	specs := []DefaultSpec{
		{Key: "PORT", Value: "8080"},
		{Key: "APP_NAME", Value: "other"},
	}
	out := ApplyDefaults(s, specs)
	if v, _ := out.Get("PORT"); v != "8080" {
		t.Errorf("expected PORT=8080, got %q", v)
	}
	if v, _ := out.Get("APP_NAME"); v != "myapp" {
		t.Errorf("expected APP_NAME=myapp (unchanged), got %q", v)
	}
}

func TestApplyDefaults_FillsEmpty(t *testing.T) {
	s := baseDefaultSet()
	specs := []DefaultSpec{{Key: "LOG_LEVEL", Value: "info"}}
	out := ApplyDefaults(s, specs)
	if v, _ := out.Get("LOG_LEVEL"); v != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %q", v)
	}
}

func TestApplyDefaults_Override(t *testing.T) {
	s := baseDefaultSet()
	specs := []DefaultSpec{{Key: "APP_NAME", Value: "forced", Override: true}}
	out := ApplyDefaults(s, specs)
	if v, _ := out.Get("APP_NAME"); v != "forced" {
		t.Errorf("expected APP_NAME=forced, got %q", v)
	}
}

func TestApplyDefaults_DoesNotMutateOriginal(t *testing.T) {
	s := baseDefaultSet()
	specs := []DefaultSpec{{Key: "NEW_KEY", Value: "val"}}
	ApplyDefaults(s, specs)
	if _, ok := s.Get("NEW_KEY"); ok {
		t.Error("original set should not be mutated")
	}
}

func TestMissingDefaults_ReturnsAbsentKeys(t *testing.T) {
	s := baseDefaultSet()
	specs := []DefaultSpec{
		{Key: "APP_NAME", Value: "x"},
		{Key: "MISSING_ONE", Value: "y"},
		{Key: "MISSING_TWO", Value: "z"},
	}
	missing := MissingDefaults(s, specs)
	sort.Strings(missing)
	if len(missing) != 2 || missing[0] != "MISSING_ONE" || missing[1] != "MISSING_TWO" {
		t.Errorf("unexpected missing keys: %v", missing)
	}
}

func TestDefaultsFromMap_Converts(t *testing.T) {
	m := map[string]string{"A": "1", "B": "2"}
	specs := DefaultsFromMap(m)
	if len(specs) != 2 {
		t.Errorf("expected 2 specs, got %d", len(specs))
	}
}
