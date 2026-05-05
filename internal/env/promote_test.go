package env

import (
	"strings"
	"testing"
)

func basePromoteSrc() *Set {
	s := NewSet()
	s.Set("APP_ENV", "production")
	s.Set("DB_HOST", "prod-db.internal")
	s.Set("LOG_LEVEL", "warn")
	return s
}

func basePromoteDst() *Set {
	s := NewSet()
	s.Set("APP_ENV", "staging")
	s.Set("CACHE_TTL", "300")
	return s
}

func TestPromote_AddsNewKeys(t *testing.T) {
	dst := basePromoteDst()
	src := basePromoteSrc()

	results, err := Promote(dst, src, PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val, ok := dst.Get("DB_HOST")
	if !ok || val != "prod-db.internal" {
		t.Errorf("expected DB_HOST=prod-db.internal, got %q", val)
	}

	var added []string
	for _, r := range results {
		if r.Action == "added" {
			added = append(added, r.Key)
		}
	}
	if len(added) == 0 {
		t.Error("expected at least one added result")
	}
}

func TestPromote_UpdatesExistingKeys(t *testing.T) {
	dst := basePromoteDst()
	src := basePromoteSrc()

	_, err := Promote(dst, src, PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val, _ := dst.Get("APP_ENV")
	if val != "production" {
		t.Errorf("expected APP_ENV=production, got %q", val)
	}
}

func TestPromote_SkipExisting(t *testing.T) {
	dst := basePromoteDst()
	src := basePromoteSrc()

	_, err := Promote(dst, src, PromoteOptions{SkipExisting: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val, _ := dst.Get("APP_ENV")
	if val != "staging" {
		t.Errorf("expected APP_ENV to remain staging, got %q", val)
	}
}

func TestPromote_DryRun_DoesNotMutate(t *testing.T) {
	dst := basePromoteDst()
	src := basePromoteSrc()

	_, err := Promote(dst, src, PromoteOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := dst.Get("DB_HOST")
	if ok {
		t.Error("dry run should not have added DB_HOST to dst")
	}
}

func TestPromote_FilterByKeys(t *testing.T) {
	dst := basePromoteDst()
	src := basePromoteSrc()

	results, err := Promote(dst, src, PromoteOptions{Keys: []string{"LOG_LEVEL"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 || results[0].Key != "LOG_LEVEL" {
		t.Errorf("expected only LOG_LEVEL in results, got %+v", results)
	}

	_, ok := dst.Get("DB_HOST")
	if ok {
		t.Error("DB_HOST should not have been promoted")
	}
}

func TestPromote_NilSrcReturnsError(t *testing.T) {
	dst := NewSet()
	_, err := Promote(dst, nil, PromoteOptions{})
	if err == nil {
		t.Error("expected error for nil src")
	}
}

func TestFormatPromoteResults_ContainsActions(t *testing.T) {
	results := []PromoteResult{
		{Key: "A", OldValue: "", NewValue: "1", Action: "added"},
		{Key: "B", OldValue: "old", NewValue: "new", Action: "updated"},
		{Key: "C", OldValue: "x", NewValue: "x", Action: "skipped"},
	}
	out := FormatPromoteResults(results)
	if !strings.Contains(out, "+ A") {
		t.Error("expected added marker for A")
	}
	if !strings.Contains(out, "~ B") {
		t.Error("expected updated marker for B")
	}
	if !strings.Contains(out, "skipped") {
		t.Error("expected skipped label for C")
	}
}
