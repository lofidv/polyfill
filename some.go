package polyfill

// Some tests whether at least one element satisfies the predicate
func (s *Slice[T]) Some(f func(T) bool) bool {
	for _, v := range s.elements {
		if f(v) {
			return true
		}
	}
	return false
}
