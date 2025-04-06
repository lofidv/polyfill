package polyfill

// Find returns the first element that satisfies the predicate
func (s *Slice[T]) Find(f func(T) bool) (T, bool) {
	for _, v := range s.elements {
		if f(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}
