package profile_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"envdiff/internal/profile"
)

func TestRegistry_AddAndGet(t *testing.T) {
	r := profile.NewRegistry()
	p := profile.Profile{Name: "staging", File: ".env.staging"}
	r.Add(p)

	got, ok := r.Get("staging")
	if !ok {
		t.Fatal("expected profile to exist")
	}
	if got.File != ".env.staging" {
		t.Errorf("got file %q, want .env.staging", got.File)
	}
}

func TestRegistry_Remove(t *testing.T) {
	r := profile.NewRegistry()
	r.Add(profile.Profile{Name: "prod", File: ".env.prod"})
	r.Remove("prod")
	_, ok := r.Get("prod")
	if ok {
		t.Error("expected profile to be removed")
	}
}

func TestRegistry_List(t *testing.T) {
	r := profile.NewRegistry()
	r.Add(profile.Profile{Name: "a", File: "a.env"})
	r.Add(profile.Profile{Name: "b", File: "b.env"})
	r.Add(profile.Profile{Name: "c", File: "c.env"})

	names := r.List()
	sort.Strings(names)
	if len(names) != 3 || names[0] != "a" || names[2] != "c" {
		t.Errorf("unexpected list: %v", names)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profiles.json")

	r := profile.NewRegistry()
	r.Add(profile.Profile{
		Name: "dev",
		File: ".env.dev",
		Tags: []string{"local"},
		Meta: map[string]string{"owner": "alice"},
	})

	if err := profile.Save(path, r); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := profile.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	p, ok := loaded.Get("dev")
	if !ok {
		t.Fatal("expected dev profile after reload")
	}
	if p.File != ".env.dev" {
		t.Errorf("file mismatch: %q", p.File)
	}
	if p.Meta["owner"] != "alice" {
		t.Errorf("meta mismatch: %v", p.Meta)
	}
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	r, err := profile.Load("/nonexistent/path/profiles.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(r.List()) != 0 {
		t.Error("expected empty registry")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)
	_, err := profile.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
