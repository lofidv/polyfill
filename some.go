package polyfill

// Some returns true if any element satisfies the predicate (like JS array.some)
//
// Example:
//
//	From([]int{1, 2, 3}).Some(func(n int) bool { return n > 2 }) // true
func (s *Seq[T]) Some(f func(T) bool) bool {
	for _, v := range s.elements {
		if f(v) {
			return true
		}
	}
	return false
}
