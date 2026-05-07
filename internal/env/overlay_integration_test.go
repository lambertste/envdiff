package env_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/env"
)

func TestOverlay_Integration_ThreeLayers(t *testing.T) {
	base := env.NewSet()
	base.Set("APP_ENV", "development")
	base.Set("DB_HOST", "localhost")
	base.Set("LOG_LEVEL", "debug")

	staging := env.NewSet()
	staging.Set("APP_ENV", "staging")
	staging.Set("DB_HOST", "staging-db.internal")

	prod := env.NewSet()
	prod.Set("APP_ENV", "production")
	prod.Set("NEW_RELIC_KEY", "abc123")

	opts := env.DefaultOverlayOptions()
	out, err := env.Overlay([]*env.Set{base, staging, prod}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := map[string]string{
		"APP_ENV":      "production",
		"DB_HOST":      "staging-db.internal",
		"LOG_LEVEL":    "debug",
		"NEW_RELIC_KEY": "abc123",
	}
	for k, want := range cases {
		got, ok := out.Get(k)
		if !ok {
			t.Errorf("key %s missing from output", k)
			continue
		}
		if got != want {
			t.Errorf("key %s: want %q, got %q", k, want, got)
		}
	}
}

func TestOverlay_Integration_SkipEmptyPreservesBase(t *testing.T) {
	base := env.NewSet()
	base.Set("SECRET_KEY", "super-secret")
	base.Set("API_URL", "https://api.example.com")

	patch := env.NewSet()
	patch.Set("SECRET_KEY", "")
	patch.Set("API_URL", "https://api.prod.example.com")

	opts := env.DefaultOverlayOptions()
	opts.SkipEmpty = true

	out, err := env.Overlay([]*env.Set{base, patch}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret, _ := out.Get("SECRET_KEY")
	if secret != "super-secret" {
		t.Errorf("expected secret preserved, got %q", secret)
	}
	api, _ := out.Get("API_URL")
	if api != "https://api.prod.example.com" {
		t.Errorf("expected api url overwritten, got %q", api)
	}
}
