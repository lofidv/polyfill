package polyfill

// Every returns true if all elements satisfy the predicate (like JS array.every)
//
// Example:
//
//	From([]int{2, 4, 6}).Every(func(n int) bool { return n%2 == 0 }) // true
func (s *Seq[T]) Every(f func(T) bool) bool {
	for _, v := range s.elements {
		if !f(v) {
			return false
		}
	}
	return true
}
