package env_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/env"
	"github.com/user/envdiff/internal/parser"
)

const groupIntegEnv = `
DB_HOST=db.example.com
DB_PORT=5432
DB_NAME=appdb
APP_PORT=8080
APP_DEBUG=false
CACHE_URL=redis://localhost
`

func prefixFn(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return ""
}

func TestGroupBy_Integration_CorrectPartition(t *testing.T) {
	entries, err := parser.ParseReader(strings.NewReader(groupIntegEnv))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := env.FromEntries(entries)
	gr := env.GroupBy(s, prefixFn)

	expected := map[string]int{"DB": 3, "APP": 2, "CACHE": 1}
	for grp, count := range expected {
		g, ok := gr[grp]
		if !ok {
			t.Errorf("group %q not found", grp)
			continue
		}
		if len(g.Keys()) != count {
			t.Errorf("group %q: expected %d keys, got %d", grp, count, len(g.Keys()))
		}
	}
}

func TestMergeGroups_Integration_RoundTrip(t *testing.T) {
	entries, err := parser.ParseReader(strings.NewReader(groupIntegEnv))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := env.FromEntries(entries)
	gr := env.GroupBy(s, prefixFn)
	merged := env.MergeGroups(gr)

	for _, k := range s.Keys() {
		origVal, _ := s.Get(k)
		mergedVal, ok := merged.Get(k)
		if !ok {
			t.Errorf("key %q missing after MergeGroups", k)
			continue
		}
		if origVal != mergedVal {
			t.Errorf("key %q: original=%q merged=%q", k, origVal, mergedVal)
		}
	}
}
