package polyfill

import "slices"

// Reverse returns a reversed copy of the sequence. The receiver is not modified.
//
//	From([]int{1, 2, 3}).Reverse().Slice() // [3 2 1]
func (s *Seq[T]) Reverse() *Seq[T] {
	if s == nil || s.err != nil {
		return s
	}
	result := slices.Clone(s.elements)
	slices.Reverse(result)
	return From(result)
}
