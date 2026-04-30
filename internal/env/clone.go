package env

// CloneOption configures how a Set is cloned.
type CloneOption func(*cloneConfig)

type cloneConfig struct {
	keys      []string
	exclude   []string
	deepCopy  bool
}

// WithKeys restricts the clone to only the specified keys.
func WithKeys(keys ...string) CloneOption {
	return func(c *cloneConfig) {
		c.keys = append(c.keys, keys...)
	}
}

// WithoutKeys excludes the specified keys from the clone.
func WithoutKeys(keys ...string) CloneOption {
	return func(c *cloneConfig) {
		c.exclude = append(c.exclude, keys...)
	}
}

// Clone creates a copy of the Set, applying any CloneOptions provided.
// By default all keys are copied. Order is preserved from the source Set.
func Clone(src *Set, opts ...CloneOption) *Set {
	cfg := &cloneConfig{}
	for _, o := range opts {
		o(cfg)
	}

	excludeSet := make(map[string]bool, len(cfg.exclude))
	for _, k := range cfg.exclude {
		excludeSet[k] = true
	}

	allowSet := make(map[string]bool, len(cfg.keys))
	for _, k := range cfg.keys {
		allowSet[k] = true
	}

	dst := NewSet()
	for _, k := range src.Keys() {
		if excludeSet[k] {
			continue
		}
		if len(allowSet) > 0 && !allowSet[k] {
			continue
		}
		v, _ := src.Get(k)
		dst.Set(k, v)
	}
	return dst
}

// Merge copies all keys from src into dst, overwriting existing values.
// dst is modified in place and returned for convenience.
func MergeInto(dst, src *Set) *Set {
	for _, k := range src.Keys() {
		v, _ := src.Get(k)
		dst.Set(k, v)
	}
	return dst
}
