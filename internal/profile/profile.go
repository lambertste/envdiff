package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Profile represents a named environment configuration profile.
type Profile struct {
	Name    string            `json:"name"`
	File    string            `json:"file"`
	Tags    []string          `json:"tags,omitempty"`
	Meta    map[string]string `json:"meta,omitempty"`
}

// Registry holds a collection of named profiles.
type Registry struct {
	Profiles map[string]Profile `json:"profiles"`
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{Profiles: make(map[string]Profile)}
}

// Add inserts or replaces a profile in the registry.
func (r *Registry) Add(p Profile) {
	r.Profiles[p.Name] = p
}

// Get retrieves a profile by name.
func (r *Registry) Get(name string) (Profile, bool) {
	p, ok := r.Profiles[name]
	return p, ok
}

// Remove deletes a profile by name.
func (r *Registry) Remove(name string) {
	delete(r.Profiles, name)
}

// List returns all profile names in the registry.
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.Profiles))
	for k := range r.Profiles {
		names = append(names, k)
	}
	return names
}

// Save persists the registry to a JSON file at the given path.
func Save(path string, r *Registry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("profile: mkdir: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("profile: create: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(r); err != nil {
		return fmt.Errorf("profile: encode: %w", err)
	}
	return nil
}

// Load reads a registry from a JSON file at the given path.
func Load(path string) (*Registry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewRegistry(), nil
		}
		return nil, fmt.Errorf("profile: open: %w", err)
	}
	defer f.Close()
	var r Registry
	if err := json.NewDecoder(f).Decode(&r); err != nil {
		return nil, fmt.Errorf("profile: decode: %w", err)
	}
	if r.Profiles == nil {
		r.Profiles = make(map[string]Profile)
	}
	return &r, nil
}
