package env

// DefaultSpec describes a single default value entry.
type DefaultSpec struct {
	Key      string
	Value    string
	Override bool // if true, overwrite existing value
}

// ApplyDefaults sets default values on the given Set.
// By default it only fills in keys that are missing or empty.
// If spec.Override is true the existing value is replaced unconditionally.
func ApplyDefaults(s *Set, specs []DefaultSpec) *Set {
	out := Clone(s)
	for _, spec := range specs {
		existing, ok := out.Get(spec.Key)
		if !ok || existing == "" || spec.Override {
			out.Set(spec.Key, spec.Value)
		}
	}
	return out
}

// MissingDefaults returns the keys from specs that are absent in s.
func MissingDefaults(s *Set, specs []DefaultSpec) []string {
	var missing []string
	for _, spec := range specs {
		_, ok := s.Get(spec.Key)
		if !ok {
			missing = append(missing, spec.Key)
		}
	}
	return missing
}

// DefaultsFromMap converts a plain map into a slice of DefaultSpec.
func DefaultsFromMap(m map[string]string) []DefaultSpec {
	specs := make([]DefaultSpec, 0, len(m))
	for k, v := range m {
		specs = append(specs, DefaultSpec{Key: k, Value: v})
	}
	return specs
}
