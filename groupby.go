package polyfill

// GroupBy groups elements by a key function
//
// Example:
//
//	From([]int{1, 2, 3, 4}).GroupBy(func(n int) any { return n % 2 }) // map[0:[2,4] 1:[1,3]]
func (s *Seq[T]) GroupBy(keyFn func(T) any) map[any][]T {
	result := make(map[any][]T)
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
//	evens, odds := From([]int{1, 2, 3, 4}).Partition(func(n int) bool { return n%2 == 0 })
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
