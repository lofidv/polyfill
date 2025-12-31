package polyfill

// Find returns the first element matching the predicate (like JS array.find)
// Returns zero value and false if not found
//
// Example:
//
//	val, ok := From([]int{1, 2, 3}).Find(func(n int) bool { return n > 1 }) // 2, true
func (s *Seq[T]) Find(f func(T) bool) (T, bool) {
	for _, v := range s.elements {
		if f(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}
