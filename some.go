package polyfill

// Some returns true if any element satisfies the predicate (like JS array.some)
//
// Example:
//
//	hasAdults := polyfill.From(people).
//	    Some(func(p Person) bool { return p.Age >= 18 })
func (s *Seq[T]) Some(f func(T) bool) bool {
	for _, v := range s.elements {
		if f(v) {
			return true
		}
	}
	return false
}
