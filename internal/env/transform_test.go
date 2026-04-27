package env

import (
	"testing"
)

func baseTransformSet() *Set {
	s := NewSet()
	s.Set("APP_NAME", "  myapp  ")
	s.Set("DB_PASSWORD", "s3cr3t")
	s.Set("API_KEY", "key-abc")
	s.Set("LOG_LEVEL", "debug")
	return s
}

func TestTransform_NoFns(t *testing.T) {
	s := baseTransformSet()
	out := Transform(s)
	if v, _ := out.Get("LOG_LEVEL"); v != "debug" {
		t.Errorf("expected 'debug', got %q", v)
	}
}

func TestTransform_TrimValues(t *testing.T) {
	s := baseTransformSet()
	out := Transform(s, TrimValues())
	if v, _ := out.Get("APP_NAME"); v != "myapp" {
		t.Errorf("expected 'myapp', got %q", v)
	}
}

func TestTransform_UppercaseValues(t *testing.T) {
	s := baseTransformSet()
	out := Transform(s, UppercaseValues())
	if v, _ := out.Get("LOG_LEVEL"); v != "DEBUG" {
		t.Errorf("expected 'DEBUG', got %q", v)
	}
}

func TestTransform_MaskSecrets(t *testing.T) {
	s := baseTransformSet()
	out := Transform(s, MaskSecrets("password", "key"))

	if v, _ := out.Get("DB_PASSWORD"); v != "***REDACTED***" {
		t.Errorf("DB_PASSWORD: expected redacted, got %q", v)
	}
	if v, _ := out.Get("API_KEY"); v != "***REDACTED***" {
		t.Errorf("API_KEY: expected redacted, got %q", v)
	}
	if v, _ := out.Get("APP_NAME"); v == "***REDACTED***" {
		t.Errorf("APP_NAME should not be redacted")
	}
}

func TestTransform_PrefixValues(t *testing.T) {
	s := NewSet()
	s.Set("ENV", "prod")
	out := Transform(s, PrefixValues("env:"))
	if v, _ := out.Get("ENV"); v != "env:prod" {
		t.Errorf("expected 'env:prod', got %q", v)
	}
}

func TestTransform_ChainedFns(t *testing.T) {
	s := NewSet()
	s.Set("REGION", "  us-east-1  ")
	out := Transform(s, TrimValues(), UppercaseValues())
	if v, _ := out.Get("REGION"); v != "US-EAST-1" {
		t.Errorf("expected 'US-EAST-1', got %q", v)
	}
}

func TestTransform_DoesNotMutateOriginal(t *testing.T) {
	s := NewSet()
	s.Set("X", "original")
	Transform(s, UppercaseValues())
	if v, _ := s.Get("X"); v != "original" {
		t.Errorf("original set was mutated: got %q", v)
	}
}
