package polyfill

import (
	"sort"
)

// Sort sorts the slice in place using the provided less function
func (s *Slice[T]) Sort(less func(a, b T) bool) *Slice[T] {
	sort.Slice(s.elements, func(i, j int) bool {
		return less(s.elements[i], s.elements[j])
	})
	return s
}
