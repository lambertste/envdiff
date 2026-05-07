package env

import "fmt"

// InheritOptions controls how inheritance merges a parent set into a child set.
type InheritOptions struct {
	// OverwriteExisting allows parent values to overwrite keys already in the child.
	OverwriteExisting bool
	// SkipEmpty skips parent keys whose values are empty.
	SkipEmpty bool
}

// DefaultInheritOptions returns conservative defaults: fill gaps, skip empty.
func DefaultInheritOptions() InheritOptions {
	return InheritOptions{
		OverwriteExisting: false,
		SkipEmpty:         true,
	}
}

// InheritResult records what happened during an Inherit call.
type InheritResult struct {
	Inherited []string // keys copied from parent
	Skipped   []string // keys skipped (already present or empty)
}

// Inherit copies keys from parent into child according to opts.
// It never mutates parent.
func Inherit(child, parent *Set, opts InheritOptions) (*Set, InheritResult, error) {
	if child == nil {
		return nil, InheritResult{}, fmt.Errorf("inherit: child set must not be nil")
	}
	if parent == nil {
		return nil, InheritResult{}, fmt.Errorf("inherit: parent set must not be nil")
	}

	out := Clone(child)
	var result InheritResult

	for _, k := range parent.Keys() {
		v, _ := parent.Get(k)
		if opts.SkipEmpty && v == "" {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		_, exists := out.Get(k)
		if exists && !opts.OverwriteExisting {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		out.Set(k, v)
		result.Inherited = append(result.Inherited, k)
	}

	return out, result, nil
}

// FormatInheritResult returns a human-readable summary of an InheritResult.
func FormatInheritResult(r InheritResult) string {
	if len(r.Inherited) == 0 && len(r.Skipped) == 0 {
		return "inherit: nothing to do\n"
	}
	out := ""
	for _, k := range r.Inherited {
		out += fmt.Sprintf("  inherited: %s\n", k)
	}
	for _, k := range r.Skipped {
		out += fmt.Sprintf("  skipped:   %s\n", k)
	}
	return out
}
