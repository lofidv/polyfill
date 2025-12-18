package polyfill

// FindIndex returns the index of the first element matching the predicate (like JS array.findIndex)
// Returns -1 if not found
//
// Example:
//
//	idx := polyfill.From(people).
//	    FindIndex(func(p Person) bool { return p.Name == "Bob" })
func (s *Seq[T]) FindIndex(f func(T) bool) int {
	for i, v := range s.elements {
		if f(v) {
			return i
		}
	}
	return -1
}

// IndexOf returns the index of the first occurrence of value (like JS array.indexOf)
// Only works with comparable types
// Returns -1 if not found
//
// Example:
//
//	idx := polyfill.From([]int{1, 2, 3}).IndexOf(2) // returns 1
func IndexOf[T comparable](s *Seq[T], value T) int {
	for i, v := range s.elements {
		if v == value {
			return i
		}
	}
	return -1
}

// Includes returns true if value is found (like JS array.includes)
func Includes[T comparable](s *Seq[T], value T) bool {
	return IndexOf(s, value) != -1
}
