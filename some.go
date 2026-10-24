package polyfill

// Some reports whether at least one element satisfies predicate.
// It returns false on an empty sequence or a failed chain.
//
//	From([]int{1, 2, 3}).Some(func(n int) bool { return n > 2 }) // true
func (s *Seq[T]) Some(predicate func(T) bool) bool {
	if s == nil || s.err != nil {
		return false
	}
	for _, v := range s.elements {
		if predicate(v) {
			return true
		}
	}
	return false
}
