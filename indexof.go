package polyfill

// FindIndex returns the index of the first element matching the predicate (like JS array.findIndex)
// Returns -1 if not found
//
// Example:
//
//	From([]int{1, 2, 3}).FindIndex(func(n int) bool { return n > 1 }) // 1
func (s *Seq[T]) FindIndex(f func(T) bool) int {
	for i, v := range s.elements {
		if f(v) {
			return i
		}
	}
	return -1
}

// IndexOf returns the index of the first occurrence of value (like JS array.indexOf)
// Returns -1 if not found
//
// Example:
//
//	From([]int{1, 2, 3}).IndexOf(2) // 1
func (s *Seq[T]) IndexOf(value T) int {
	for i, v := range s.elements {
		if any(v) == any(value) {
			return i
		}
	}
	return -1
}

// Includes returns true if value is found (like JS array.includes)
//
// Example:
//
//	From([]string{"a", "b"}).Includes("b") // true
func (s *Seq[T]) Includes(value T) bool {
	return s.IndexOf(value) != -1
}
