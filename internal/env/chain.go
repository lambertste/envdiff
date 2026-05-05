package env

// ChainOption configures a Chain operation.
type ChainOption func(*chainConfig)

type chainConfig struct {
	overwrite bool
}

// WithOverwrite allows later sets to overwrite keys from earlier ones.
func WithOverwrite() ChainOption {
	return func(c *chainConfig) {
		c.overwrite = true
	}
}

// Chain merges multiple Sets in order. By default, the first definition of a
// key wins (safe merge). Pass WithOverwrite to let later sets win instead.
func Chain(sets []*Set, opts ...ChainOption) *Set {
	cfg := &chainConfig{overwrite: false}
	for _, o := range opts {
		o(cfg)
	}

	out := NewSet()
	for _, s := range sets {
		for _, k := range s.Keys() {
			v, _ := s.Get(k)
			_, exists := out.Get(k)
			if !exists || cfg.overwrite {
				out.Set(k, v)
			}
		}
	}
	return out
}

// ChainKeys returns the ordered list of unique keys across all sets,
// preserving first-seen order.
func ChainKeys(sets []*Set) []string {
	seen := make(map[string]struct{})
	var keys []string
	for _, s := range sets {
		for _, k := range s.Keys() {
			if _, ok := seen[k]; !ok {
				seen[k] = struct{}{}
				keys = append(keys, k)
			}
		}
	}
	return keys
}
