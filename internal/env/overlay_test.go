package env

import (
	"testing"
)

func baseOverlaySet() *Set {
	s := NewSet()
	s.Set("APP_ENV", "staging")
	s.Set("DB_HOST", "localhost")
	s.Set("LOG_LEVEL", "debug")
	return s
}

func TestOverlay_EmptyLayers_ReturnsError(t *testing.T) {
	_, err := Overlay(nil, DefaultOverlayOptions())
	if err == nil {
		t.Fatal("expected error for empty layers")
	}
}

func TestOverlay_SingleLayer_ReturnsClone(t *testing.T) {
	base := baseOverlaySet()
	out, err := Overlay([]*Set{base}, DefaultOverlayOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := out.Get("APP_ENV")
	if v != "staging" {
		t.Errorf("expected staging, got %s", v)
	}
}

func TestOverlay_Overwrite_LaterLayerWins(t *testing.T) {
	base := baseOverlaySet()
	over := NewSet()
	over.Set("APP_ENV", "production")
	opts := DefaultOverlayOptions()
	opts.Overwrite = true
	out, err := Overlay([]*Set{base, over}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := out.Get("APP_ENV")
	if v != "production" {
		t.Errorf("expected production, got %s", v)
	}
}

func TestOverlay_NoOverwrite_BaseWins(t *testing.T) {
	base := baseOverlaySet()
	over := NewSet()
	over.Set("APP_ENV", "production")
	opts := DefaultOverlayOptions()
	opts.Overwrite = false
	out, err := Overlay([]*Set{base, over}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := out.Get("APP_ENV")
	if v != "staging" {
		t.Errorf("expected staging, got %s", v)
	}
}

func TestOverlay_SkipEmpty_EmptyValuesIgnored(t *testing.T) {
	base := baseOverlaySet()
	over := NewSet()
	over.Set("LOG_LEVEL", "")
	opts := DefaultOverlayOptions()
	opts.SkipEmpty = true
	out, err := Overlay([]*Set{base, over}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := out.Get("LOG_LEVEL")
	if v != "debug" {
		t.Errorf("expected debug, got %q", v)
	}
}

func TestOverlayWithReport_Added(t *testing.T) {
	base := baseOverlaySet()
	over := NewSet()
	over.Set("NEW_KEY", "value")
	out, report, err := OverlayWithReport(base, over, DefaultOverlayOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Added) != 1 || report.Added[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY in Added, got %v", report.Added)
	}
	v, _ := out.Get("NEW_KEY")
	if v != "value" {
		t.Errorf("expected value, got %s", v)
	}
}

func TestOverlayWithReport_Overwritten(t *testing.T) {
	base := baseOverlaySet()
	over := NewSet()
	over.Set("APP_ENV", "production")
	_, report, err := OverlayWithReport(base, over, DefaultOverlayOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Overwritten) != 1 || report.Overwritten[0] != "APP_ENV" {
		t.Errorf("expected APP_ENV in Overwritten, got %v", report.Overwritten)
	}
}

func TestOverlayWithReport_Preserved(t *testing.T) {
	base := baseOverlaySet()
	over := NewSet()
	over.Set("APP_ENV", "production")
	opts := DefaultOverlayOptions()
	opts.Overwrite = false
	_, report, err := OverlayWithReport(base, over, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Preserved) != 1 || report.Preserved[0] != "APP_ENV" {
		t.Errorf("expected APP_ENV in Preserved, got %v", report.Preserved)
	}
}
