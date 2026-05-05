package env

import (
	"testing"
)

// TestChain_Integration_LayeredConfig simulates a real layered config scenario:
// defaults < staging < local overrides.
func TestChain_Integration_LayeredConfig(t *testing.T) {
	defaults := NewSet()
	defaults.Set("LOG_LEVEL", "warn")
	defaults.Set("DB_PORT", "5432")
	defaults.Set("TIMEOUT", "30s")

	staging := NewSet()
	staging.Set("LOG_LEVEL", "info")
	staging.Set("DB_HOST", "staging-db.internal")

	local := NewSet()
	local.Set("LOG_LEVEL", "debug")
	local.Set("DB_HOST", "localhost")

	// Without overwrite: defaults win for conflicts
	safe := Chain([]*Set{defaults, staging, local})
	if v, _ := safe.Get("LOG_LEVEL"); v != "warn" {
		t.Errorf("safe chain: expected warn, got %s", v)
	}

	// With overwrite: last (local) wins
	over := Chain([]*Set{defaults, staging, local}, WithOverwrite())
	if v, _ := over.Get("LOG_LEVEL"); v != "debug" {
		t.Errorf("overwrite chain: expected debug, got %s", v)
	}
	if v, _ := over.Get("DB_HOST"); v != "localhost" {
		t.Errorf("overwrite chain: expected localhost, got %s", v)
	}
	if v, _ := over.Get("DB_PORT"); v != "5432" {
		t.Errorf("overwrite chain: expected 5432, got %s", v)
	}
}

// TestChain_Integration_RoundTripWithFilter ensures Chain output can be
// further filtered without side effects on original sets.
func TestChain_Integration_RoundTripWithFilter(t *testing.T) {
	a := NewSet()
	a.Set("DB_HOST", "db")
	a.Set("DB_PORT", "5432")
	a.Set("CACHE_HOST", "redis")

	b := NewSet()
	b.Set("DB_NAME", "myapp")
	b.Set("CACHE_TTL", "60")

	merged := Chain([]*Set{a, b}, WithOverwrite())
	dbOnly := Filter(merged, WithPrefix("DB_"))

	keys := dbOnly.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 DB_ keys, got %d", len(keys))
	}
	if _, ok := dbOnly.Get("CACHE_HOST"); ok {
		t.Error("CACHE_HOST should not be in filtered set")
	}
}
