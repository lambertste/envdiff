package env

import "sort"

// CompareResult holds the outcome of comparing two EnvSets.
type CompareResult struct {
	Added    []string // keys present in b but not a
	Removed  []string // keys present in a but not b
	Modified []string // keys present in both but with different values
	Unchanged []string // keys present in both with identical values
}

// Compare returns a CompareResult describing the difference between sets a and b.
func Compare(a, b *Set) CompareResult {
	result := CompareResult{}

	aKeys := a.Keys()
	bKeys := b.Keys()

	aMap := make(map[string]struct{}, len(aKeys))
	for _, k := range aKeys {
		aMap[k] = struct{}{}
	}

	bMap := make(map[string]struct{}, len(bKeys))
	for _, k := range bKeys {
		bMap[k] = struct{}{}
	}

	for _, k := range aKeys {
		if _, ok := bMap[k]; !ok {
			result.Removed = append(result.Removed, k)
			continue
		}
		av, _ := a.Get(k)
		bv, _ := b.Get(k)
		if av == bv {
			result.Unchanged = append(result.Unchanged, k)
		} else {
			result.Modified = append(result.Modified, k)
		}
	}

	for _, k := range bKeys {
		if _, ok := aMap[k]; !ok {
			result.Added = append(result.Added, k)
		}
	}

	sort.Strings(result.Added)
	sort.Strings(result.Removed)
	sort.Strings(result.Modified)
	sort.Strings(result.Unchanged)

	return result
}

// HasChanges returns true if there are any additions, removals, or modifications.
func (r CompareResult) HasChanges() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Modified) > 0
}

// Summary returns a brief human-readable description of the comparison.
func (r CompareResult) Summary() string {
	if !r.HasChanges() {
		return "no changes"
	}
	return fmt.Sprintf("+%d added, -%d removed, ~%d modified",
		len(r.Added), len(r.Removed), len(r.Modified))
}
