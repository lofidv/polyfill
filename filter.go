package polyfill

// Filter returns a new slice with elements that satisfy the predicate
func (s *Slice[T]) Filter(f func(T) bool) *Slice[T] {
	result := make([]T, 0, len(s.elements))
	for _, v := range s.elements {
		if f(v) {
			result = append(result, v)
		}
	}
	return Wrap(result)
}
