package env

import (
	"fmt"
	"regexp"
	"strings"
)

var varPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// InterpolateError describes a missing variable during interpolation.
type InterpolateError struct {
	Key     string
	Missing string
}

func (e *InterpolateError) Error() string {
	return fmt.Sprintf("interpolate: key %q references undefined variable %q", e.Key, e.Missing)
}

// Interpolate expands variable references in the values of s using its own
// entries as the source. References take the form $VAR or ${VAR}.
// Returns the expanded Set and a slice of any unresolved reference errors.
func Interpolate(s *Set) (*Set, []error) {
	out := NewSet()
	var errs []error

	for _, key := range s.Keys() {
		val, _ := s.Get(key)
		expanded, err := expand(key, val, s)
		if err != nil {
			errs = append(errs, err)
			expanded = val // keep original on error
		}
		out.Set(key, expanded)
	}

	return out, errs
}

func expand(key, val string, s *Set) (string, error) {
	var firstErr error
	result := varPattern.ReplaceAllStringFunc(val, func(match string) string {
		name := strings.TrimPrefix(strings.Trim(match, "${}"), "$")
		if resolved, ok := s.Get(name); ok {
			return resolved
		}
		if firstErr == nil {
			firstErr = &InterpolateError{Key: key, Missing: name}
		}
		return match
	})
	return result, firstErr
}
