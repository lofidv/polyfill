package polyfill

// Find returns the first element matching the predicate (like JS array.find)
// Returns zero value and false if not found
//
// Example:
//
//	bob, found := polyfill.From(people).
//	    Find(func(p Person) bool { return p.Name == "Bob" })
func (s *Seq[T]) Find(f func(T) bool) (T, bool) {
	for _, v := range s.elements {
		if f(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}
