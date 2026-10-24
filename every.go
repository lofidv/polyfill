package polyfill

// Every reports whether all elements satisfy predicate.
// It returns false when the chain failed.
//
//	From([]int{2, 4, 6}).Every(func(n int) bool { return n%2 == 0 }) // true
func (s *Seq[T]) Every(predicate func(T) bool) bool {
	if s == nil || s.err != nil {
		return false
	}
	for _, v := range s.elements {
		if !predicate(v) {
			return false
		}
	}
	return true
}
