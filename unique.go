package polyfill

// Unique returns a new slice with duplicate elements removed
func (s *Slice[T]) Unique(equals func(a, b T) bool) *Slice[T] {
	result := make([]T, 0, len(s.elements))
	seen := make(map[any]struct{})

	for _, v := range s.elements {
		if _, exists := seen[any(v)]; !exists {
			seen[any(v)] = struct{}{}
			result = append(result, v)
		}
	}
	return Wrap(result)
}
