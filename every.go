package polyfill

// Every tests whether all elements satisfy the predicate
func (s *Slice[T]) Every(f func(T) bool) bool {
	for _, v := range s.elements {
		if !f(v) {
			return false
		}
	}
	return true
}
