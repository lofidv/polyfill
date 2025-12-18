package polyfill

// Every returns true if all elements satisfy the predicate (like JS array.every)
//
// Example:
//
//	allAdults := polyfill.From(people).
//	    Every(func(p Person) bool { return p.Age >= 18 })
func (s *Seq[T]) Every(f func(T) bool) bool {
	for _, v := range s.elements {
		if !f(v) {
			return false
		}
	}
	return true
}
