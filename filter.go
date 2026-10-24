package polyfill

// Filter keeps the elements that satisfy predicate.
//
//	From([]int{1, 2, 3, 4}).Filter(func(n int) bool { return n%2 == 0 }).Slice() // [2 4]
func (s *Seq[T]) Filter(predicate func(T) bool) *Seq[T] {
	if s == nil || s.err != nil {
		return s
	}
	result := make([]T, 0, len(s.elements))
	for _, v := range s.elements {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return From(result)
}
