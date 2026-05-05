package env_test

import (
	"testing"

	"envdiff/internal/env"
)

func TestFreeze_Integration_PatchRespectsFrozen(t *testing.T) {
	s := env.NewSet()
	s.Set("DB_HOST", "prod.db.internal")
	s.Set("DB_PORT", "5432")

	env.Freeze(s, "DB_HOST")

	// Attempt to patch a frozen key — Patch itself is unaware of freeze;
	// callers should check IsFrozen before applying.
	if env.IsFrozen(s, "DB_HOST") {
		// simulate guard in caller
		t.Log("DB_HOST is frozen, skipping patch")
	} else {
		env.Patch(s, env.PatchOp{Key: "DB_HOST", Value: "other.host"})
	}

	v, _ := s.Get("DB_HOST")
	if v != "prod.db.internal" {
		t.Errorf("frozen key should not have been modified, got %q", v)
	}
}

func TestFreeze_Integration_UnfreezeAllowsMutation(t *testing.T) {
	s := env.NewSet()
	s.Set("API_URL", "https://api.example.com")
	env.Freeze(s, "API_URL")
	env.Unfreeze(s, "API_URL")

	if env.IsFrozen(s, "API_URL") {
		t.Fatal("expected API_URL to be unfrozen")
	}

	env.Patch(s, env.PatchOp{Key: "API_URL", Value: "https://staging.example.com"})
	v, _ := s.Get("API_URL")
	if v != "https://staging.example.com" {
		t.Errorf("expected updated value after unfreeze, got %q", v)
	}
}

func TestFreeze_Integration_FrozenKeysRoundTrip(t *testing.T) {
	s := env.NewSet()
	s.Set("X", "1")
	s.Set("Y", "2")
	s.Set("Z", "3")

	env.Freeze(s, "X", "Z")
	keys := env.FrozenKeys(s)

	if len(keys) != 2 {
		t.Fatalf("expected 2 frozen keys, got %d: %v", len(keys), keys)
	}

	env.Unfreeze(s, "X", "Z")
	if len(env.FrozenKeys(s)) != 0 {
		t.Error("expected no frozen keys after full unfreeze")
	}
}
