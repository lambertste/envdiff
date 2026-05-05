package env

import (
	"testing"
)

func baseFreezeSet() *Set {
	s := NewSet()
	s.Set("HOST", "localhost")
	s.Set("PORT", "8080")
	s.Set("SECRET", "s3cr3t")
	return s
}

func TestFreeze_MarksExistingKeys(t *testing.T) {
	s := baseFreezeSet()
	res := Freeze(s, "HOST", "PORT")
	if len(res.Frozen) != 2 {
		t.Fatalf("expected 2 frozen, got %d", len(res.Frozen))
	}
	if res.Frozen[0] != "HOST" || res.Frozen[1] != "PORT" {
		t.Errorf("unexpected frozen keys: %v", res.Frozen)
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected no skipped, got %v", res.Skipped)
	}
}

func TestFreeze_SkipsMissingKeys(t *testing.T) {
	s := baseFreezeSet()
	res := Freeze(s, "MISSING")
	if len(res.Frozen) != 0 {
		t.Errorf("expected no frozen, got %v", res.Frozen)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING" {
		t.Errorf("expected MISSING in skipped, got %v", res.Skipped)
	}
}

func TestIsFrozen_TrueAfterFreeze(t *testing.T) {
	s := baseFreezeSet()
	Freeze(s, "SECRET")
	if !IsFrozen(s, "SECRET") {
		t.Error("expected SECRET to be frozen")
	}
}

func TestIsFrozen_FalseBeforeFreeze(t *testing.T) {
	s := baseFreezeSet()
	if IsFrozen(s, "HOST") {
		t.Error("expected HOST to not be frozen")
	}
}

func TestFrozenKeys_ReturnsAll(t *testing.T) {
	s := baseFreezeSet()
	Freeze(s, "HOST", "SECRET")
	keys := FrozenKeys(s)
	if len(keys) != 2 {
		t.Fatalf("expected 2 frozen keys, got %d", len(keys))
	}
	if keys[0] != "HOST" || keys[1] != "SECRET" {
		t.Errorf("unexpected frozen keys: %v", keys)
	}
}

func TestUnfreeze_RemovesMarker(t *testing.T) {
	s := baseFreezeSet()
	Freeze(s, "HOST")
	Unfreeze(s, "HOST")
	if IsFrozen(s, "HOST") {
		t.Error("expected HOST to be unfrozen after Unfreeze")
	}
}

func TestFormatFrozen_Empty(t *testing.T) {
	out := FormatFrozen([]string{})
	if out != "no frozen keys" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatFrozen_WithKeys(t *testing.T) {
	out := FormatFrozen([]string{"HOST", "PORT"})
	if out == "" {
		t.Error("expected non-empty output")
	}
}
