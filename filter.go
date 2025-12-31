package polyfill

// Filter returns elements that satisfy the predicate (like JS array.filter)
//
// Example:
//
//	From([]int{1, 2, 3, 4}).Filter(func(n int) bool { return n > 2 }).Slice() // [3, 4]
func (s *Seq[T]) Filter(f func(T) bool) *Seq[T] {
	if s.err != nil {
		return s
	}

	result := make([]T, 0, len(s.elements))
	for _, v := range s.elements {
		if f(v) {
			result = append(result, v)
		}
	}
	return From(result)
}
