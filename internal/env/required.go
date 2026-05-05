package env

import "fmt"

// RequiredResult holds the outcome of a required-keys check.
type RequiredResult struct {
	Key     string
	Present bool
	Empty   bool
}

// CheckRequired verifies that every key in required exists in s and is non-empty.
// It returns one result per required key.
func CheckRequired(s *Set, required []string) []RequiredResult {
	results := make([]RequiredResult, 0, len(required))
	for _, key := range required {
		val, ok := s.Get(key)
		results = append(results, RequiredResult{
			Key:     key,
			Present: ok,
			Empty:   ok && val == "",
		})
	}
	return results
}

// MissingRequired returns the subset of required keys that are absent or empty.
func MissingRequired(s *Set, required []string) []string {
	var missing []string
	for _, r := range CheckRequired(s, required) {
		if !r.Present || r.Empty {
			missing = append(missing, r.Key)
		}
	}
	return missing
}

// FormatRequired produces a human-readable summary of required-key results.
func FormatRequired(results []RequiredResult) string {
	if len(results) == 0 {
		return "no required keys specified\n"
	}
	out := ""
	for _, r := range results {
		switch {
		case !r.Present:
			out += fmt.Sprintf("MISSING  %s\n", r.Key)
		case r.Empty:
			out += fmt.Sprintf("EMPTY    %s\n", r.Key)
		default:
			out += fmt.Sprintf("OK       %s\n", r.Key)
		}
	}
	return out
}
