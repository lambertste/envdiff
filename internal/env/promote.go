package env

import "fmt"

// PromoteResult holds the outcome of a promotion between two environments.
type PromoteResult struct {
	Key      string
	OldValue string
	NewValue string
	Action   string // "added", "updated", "skipped"
}

// PromoteOptions controls how promotion behaves.
type PromoteOptions struct {
	// DryRun reports what would change without modifying the target.
	DryRun bool
	// Keys restricts promotion to a specific subset of keys. Empty means all keys.
	Keys []string
	// SkipExisting skips keys that already exist in the target.
	SkipExisting bool
}

// Promote copies key-value pairs from src into dst according to opts.
// It returns the list of results describing each key's outcome.
func Promote(dst, src *Set, opts PromoteOptions) ([]PromoteResult, error) {
	if dst == nil || src == nil {
		return nil, fmt.Errorf("promote: dst and src must not be nil")
	}

	allowedKeys := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		allowedKeys[k] = true
	}

	var results []PromoteResult

	for _, key := range src.Keys() {
		if len(allowedKeys) > 0 && !allowedKeys[key] {
			continue
		}

		srcVal, _ := src.Get(key)
		dstVal, exists := dst.Get(key)

		if exists && opts.SkipExisting {
			results = append(results, PromoteResult{
				Key:      key,
				OldValue: dstVal,
				NewValue: dstVal,
				Action:   "skipped",
			})
			continue
		}

		action := "updated"
		if !exists {
			action = "added"
		}

		if !opts.DryRun {
			dst.Set(key, srcVal)
		}

		results = append(results, PromoteResult{
			Key:      key,
			OldValue: dstVal,
			NewValue: srcVal,
			Action:   action,
		})
	}

	return results, nil
}

// FormatPromoteResults returns a human-readable summary of promotion results.
func FormatPromoteResults(results []PromoteResult) string {
	if len(results) == 0 {
		return "nothing to promote\n"
	}
	out := ""
	for _, r := range results {
		switch r.Action {
		case "added":
			out += fmt.Sprintf("+ %s = %s\n", r.Key, r.NewValue)
		case "updated":
			out += fmt.Sprintf("~ %s: %s -> %s\n", r.Key, r.OldValue, r.NewValue)
		case "skipped":
			out += fmt.Sprintf("  %s (skipped)\n", r.Key)
		}
	}
	return out
}
