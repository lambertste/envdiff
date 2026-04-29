package env

import (
	"strings"
)

// MaskRule defines a strategy for masking values in an EnvSet.
type MaskRule struct {
	// KeyContains masks any key that contains one of these substrings (case-insensitive).
	KeyContains []string
	// MaskWith replaces the matched value with this string. Defaults to "***".
	MaskWith string
}

// DefaultMaskRule returns a MaskRule that masks common secret key patterns.
func DefaultMaskRule() MaskRule {
	return MaskRule{
		KeyContains: []string{"SECRET", "PASSWORD", "PASSWD", "TOKEN", "PRIVATE_KEY", "API_KEY", "CREDENTIALS"},
		MaskWith:    "***",
	}
}

// MaskSet returns a new EnvSet with values masked according to the provided rule.
// Original set is not modified.
func MaskSet(s *Set, rule MaskRule) *Set {
	mask := rule.MaskWith
	if mask == "" {
		mask = "***"
	}

	out := NewSet()
	for _, key := range s.Keys() {
		val, _ := s.Get(key)
		if shouldMask(key, rule.KeyContains) {
			out.Set(key, mask)
		} else {
			out.Set(key, val)
		}
	}
	return out
}

// MaskedKeys returns the list of keys that would be masked by the given rule.
func MaskedKeys(s *Set, rule MaskRule) []string {
	var masked []string
	for _, key := range s.Keys() {
		if shouldMask(key, rule.KeyContains) {
			masked = append(masked, key)
		}
	}
	return masked
}

func shouldMask(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}
