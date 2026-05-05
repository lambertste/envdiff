package env

import (
	"testing"
)

func baseChainSets() (*Set, *Set, *Set) {
	a := NewSet()
	a.Set("APP_HOST", "localhost")
	a.Set("APP_PORT", "8080")

	b := NewSet()
	b.Set("APP_PORT", "9090") // conflicts with a
	b.Set("APP_DEBUG", "true")

	c := NewSet()
	c.Set("APP_LOG", "info")
	return a, b, c
}

func TestChain_FirstWins(t *testing.T) {
	a, b, c := baseChainSets()
	out := Chain([]*Set{a, b, c})

	v, ok := out.Get("APP_PORT")
	if !ok {
		t.Fatal("expected APP_PORT to be present")
	}
	if v != "8080" {
		t.Errorf("expected 8080 (first wins), got %s", v)
	}
}

func TestChain_WithOverwrite_LastWins(t *testing.T) {
	a, b, c := baseChainSets()
	out := Chain([]*Set{a, b, c}, WithOverwrite())

	v, _ := out.Get("APP_PORT")
	if v != "9090" {
		t.Errorf("expected 9090 (overwrite), got %s", v)
	}
}

func TestChain_MergesAllKeys(t *testing.T) {
	a, b, c := baseChainSets()
	out := Chain([]*Set{a, b, c})

	expected := []string{"APP_HOST", "APP_PORT", "APP_DEBUG", "APP_LOG"}
	for _, k := range expected {
		if _, ok := out.Get(k); !ok {
			t.Errorf("expected key %s to be present", k)
		}
	}
}

func TestChain_EmptySlice(t *testing.T) {
	out := Chain([]*Set{})
	if len(out.Keys()) != 0 {
		t.Errorf("expected empty set, got %d keys", len(out.Keys()))
	}
}

func TestChainKeys_UniqueOrder(t *testing.T) {
	a, b, c := baseChainSets()
	keys := ChainKeys([]*Set{a, b, c})

	seen := make(map[string]int)
	for _, k := range keys {
		seen[k]++
	}
	for k, count := range seen {
		if count > 1 {
			t.Errorf("key %s appeared %d times, expected 1", k, count)
		}
	}
	if len(keys) != 4 {
		t.Errorf("expected 4 unique keys, got %d", len(keys))
	}
}

func TestChain_SingleSet(t *testing.T) {
	a, _, _ := baseChainSets()
	out := Chain([]*Set{a})

	v, ok := out.Get("APP_HOST")
	if !ok || v != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %s", v)
	}
}
