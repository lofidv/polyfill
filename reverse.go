package polyfill

// Reverse returns a reversed copy of the sequence (immutable)
// Does NOT mutate the original sequence
//
// Example:
//
//	reversed := polyfill.From([]int{1, 2, 3}).Reverse().Slice()
//	// Result: [3, 2, 1]
func (s *Seq[T]) Reverse() *Seq[T] {
	if s.err != nil {
		return s
	}

	result := make([]T, len(s.elements))
	for i := 0; i < len(s.elements); i++ {
		result[len(s.elements)-1-i] = s.elements[i]
	}
	return From(result)
}
