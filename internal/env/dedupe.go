package env

// DedupeStrategy controls how duplicate keys are resolved.
type DedupeStrategy int

const (
	// DedupeKeepFirst retains the first occurrence of a duplicate key.
	DedupeKeepFirst DedupeStrategy = iota
	// DedupeKeepLast retains the last occurrence of a duplicate key.
	DedupeKeepLast
)

// DedupeResult holds the outcome of a deduplication pass.
type DedupeResult struct {
	Set        *Set
	Duplicates []string // keys that had duplicates
}

// Dedupe removes duplicate keys from src according to the given strategy.
// Keys are returned in their original insertion order (first or last seen).
func Dedupe(src *Set, strategy DedupeStrategy) DedupeResult {
	seen := make(map[string]int) // key -> count
	for _, k := range src.Keys() {
		seen[k]++
	}

	duplicates := make([]string, 0)
	for k, count := range seen {
		if count > 1 {
			duplicates = append(duplicates, k)
		}
	}

	// Build result set respecting strategy.
	// Since Set already enforces unique keys (last-write wins internally),
	// we rebuild by iterating entries and skipping already-added keys for
	// KeepFirst, or overwriting for KeepLast.
	out := NewSet()

	switch strategy {
	case DedupeKeepFirst:
		for _, k := range src.Keys() {
			if _, exists := out.Get(k); !exists {
				v, _ := src.Get(k)
				out.Set(k, v)
			}
		}
	case DedupeKeepLast:
		for _, k := range src.Keys() {
			v, _ := src.Get(k)
			out.Set(k, v)
		}
	}

	sortStrings(duplicates)
	return DedupeResult{Set: out, Duplicates: duplicates}
}

// sortStrings sorts a slice of strings in place (simple insertion sort to avoid
// importing sort just for this helper).
func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
