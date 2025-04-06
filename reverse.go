package polyfill

// Reverse returns a new slice with elements in reverse order
func (s *Slice[T]) Reverse() *Slice[T] {
	result := make([]T, len(s.elements))
	for i := 0; i < len(s.elements); i++ {
		result[len(s.elements)-1-i] = s.elements[i]
	}
	return Wrap(result)
}
