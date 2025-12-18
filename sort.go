package polyfill

import (
	"sort"
)

// Sort returns a sorted copy of the sequence (immutable)
// Does NOT mutate the original sequence
//
// Example:
//
//	sorted := polyfill.From([]int{3, 1, 4, 2}).
//	    Sort(func(a, b int) bool { return a < b }).
//	    Slice()
func (s *Seq[T]) Sort(less func(a, b T) bool) *Seq[T] {
	if s.err != nil {
		return s
	}

	// Create a copy to avoid mutating original
	copy := make([]T, len(s.elements))
	for i, v := range s.elements {
		copy[i] = v
	}

	sort.Slice(copy, func(i, j int) bool {
		return less(copy[i], copy[j])
	})

	return From(copy)
}
