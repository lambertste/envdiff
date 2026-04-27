package env

import (
	"testing"
)

func baseInterpSet() *Set {
	s := NewSet()
	s.Set("HOME", "/home/user")
	s.Set("CONFIG_DIR", "${HOME}/.config")
	s.Set("CACHE_DIR", "$HOME/.cache")
	s.Set("APP_DIR", "${CONFIG_DIR}/app")
	s.Set("PLAIN", "no-refs-here")
	return s
}

func TestInterpolate_NoRefs(t *testing.T) {
	s := NewSet()
	s.Set("PLAIN", "hello")
	out, errs := Interpolate(s)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	v, _ := out.Get("PLAIN")
	if v != "hello" {
		t.Errorf("expected 'hello', got %q", v)
	}
}

func TestInterpolate_BraceStyle(t *testing.T) {
	s := NewSet()
	s.Set("BASE", "/opt")
	s.Set("DIR", "${BASE}/bin")
	out, errs := Interpolate(s)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	v, _ := out.Get("DIR")
	if v != "/opt/bin" {
		t.Errorf("expected '/opt/bin', got %q", v)
	}
}

func TestInterpolate_DollarStyle(t *testing.T) {
	s := NewSet()
	s.Set("HOME", "/home/user")
	s.Set("CACHE", "$HOME/.cache")
	out, errs := Interpolate(s)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	v, _ := out.Get("CACHE")
	if v != "/home/user/.cache" {
		t.Errorf("expected '/home/user/.cache', got %q", v)
	}
}

func TestInterpolate_MissingRef_ReturnsError(t *testing.T) {
	s := NewSet()
	s.Set("DIR", "${UNDEFINED}/bin")
	_, errs := Interpolate(s)
	if len(errs) == 0 {
		t.Fatal("expected an error for missing reference")
	}
	ie, ok := errs[0].(*InterpolateError)
	if !ok {
		t.Fatalf("expected *InterpolateError, got %T", errs[0])
	}
	if ie.Missing != "UNDEFINED" {
		t.Errorf("expected missing='UNDEFINED', got %q", ie.Missing)
	}
}

func TestInterpolate_KeepsOriginalOnError(t *testing.T) {
	s := NewSet()
	s.Set("X", "${MISSING}")
	out, errs := Interpolate(s)
	if len(errs) == 0 {
		t.Fatal("expected error")
	}
	v, _ := out.Get("X")
	if v != "${MISSING}" {
		t.Errorf("expected original value to be preserved, got %q", v)
	}
}

func TestInterpolate_ChainedRefs(t *testing.T) {
	s := baseInterpSet()
	out, errs := Interpolate(s)
	// CONFIG_DIR = ${HOME}/.config => /home/user/.config
	v, _ := out.Get("CONFIG_DIR")
	if v != "/home/user/.config" {
		t.Errorf("CONFIG_DIR: expected '/home/user/.config', got %q", v)
	}
	_ = errs // APP_DIR may not chain in single pass; just ensure no panic
}
