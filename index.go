package polyfill

// FindIndex returns the index of the first element that satisfies predicate,
// or -1. It returns -1 when the chain failed.
//
//	From([]int{1, 2, 3}).FindIndex(func(n int) bool { return n > 1 }) // 1
func (s *Seq[T]) FindIndex(predicate func(T) bool) int {
	if s == nil || s.err != nil {
		return -1
	}
	for i, v := range s.elements {
		if predicate(v) {
			return i
		}
	}
	return -1
}

// IndexOf returns the index of the first occurrence of value, or -1.
//
// It compares with == for comparable types and falls back to
// reflect.DeepEqual for the rest, so it does not panic on slices or maps.
func (s *Seq[T]) IndexOf(value T) int {
	if s == nil || s.err != nil {
		return -1
	}
	for i, v := range s.elements {
		if sameValue(v, value) {
			return i
		}
	}
	return -1
}

// Includes reports whether value is present.
//
//	From([]string{"a", "b"}).Includes("b") // true
func (s *Seq[T]) Includes(value T) bool {
	return s.IndexOf(value) != -1
}

// ContainsFunc reports whether any element satisfies predicate.
//
//	From([]int{1, 2, 3}).ContainsFunc(func(n int) bool { return n == 2 }) // true
func (s *Seq[T]) ContainsFunc(predicate func(T) bool) bool {
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
