package polyfill

import "slices"

// Sort returns a sorted copy of the sequence (immutable)
//
// Example:
//
//	From([]int{3, 1, 2}).Sort(func(a, b int) bool { return a < b }).Slice() // [1, 2, 3]
func (s *Seq[T]) Sort(less func(a, b T) bool) *Seq[T] {
	if s.err != nil {
		return s
	}

	// Use slices.Clone for efficient copy
	copy := slices.Clone(s.elements)
	slices.SortFunc(copy, func(a, b T) int {
		if less(a, b) {
			return -1
		}
		if less(b, a) {
			return 1
		}
		return 0
	})

	return From(copy)
}
