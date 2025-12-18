package polyfill

// GroupBy groups elements by a key function (super useful!)
//
// Example:
//
//	byType := polyfill.From(pets).
//	    GroupBy(func(p Pet) string { return p.Type })
//	// Result: map[string][]Pet{"dog": [...], "cat": [...]}
func GroupBy[T any, K comparable](s *Seq[T], keyFn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, v := range s.elements {
		key := keyFn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// Partition splits the sequence into two based on a predicate
// Returns (matching, not matching)
//
// Example:
//
//	adults, minors := polyfill.From(people).
//	    Partition(func(p Person) bool { return p.Age >= 18 })
func (s *Seq[T]) Partition(f func(T) bool) ([]T, []T) {
	matching := make([]T, 0)
	notMatching := make([]T, 0)

	for _, v := range s.elements {
		if f(v) {
			matching = append(matching, v)
		} else {
			notMatching = append(notMatching, v)
		}
	}

	return matching, notMatching
}
