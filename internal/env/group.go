package env

import "sort"

// GroupResult holds named groups of env entries.
type GroupResult map[string]*Set

// GroupBy partitions the Set into named buckets using the provided key function.
// The key function returns the group name for a given key; entries returning
// an empty string are placed in the "_default" bucket.
func GroupBy(s *Set, keyFn func(key string) string) GroupResult {
	result := make(GroupResult)
	for _, k := range s.Keys() {
		v, _ := s.Get(k)
		group := keyFn(k)
		if group == "" {
			group = "_default"
		}
		if _, ok := result[group]; !ok {
			result[group] = NewSet()
		}
		result[group].Set(k, v)
	}
	return result
}

// GroupNames returns a sorted list of group names present in a GroupResult.
func GroupNames(gr GroupResult) []string {
	names := make([]string, 0, len(gr))
	for name := range gr {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// MergeGroups combines multiple named groups back into a single Set.
// Keys from later groups in sorted order overwrite earlier ones on collision.
func MergeGroups(gr GroupResult) *Set {
	out := NewSet()
	for _, name := range GroupNames(gr) {
		g := gr[name]
		for _, k := range g.Keys() {
			v, _ := g.Get(k)
			out.Set(k, v)
		}
	}
	return out
}
