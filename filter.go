package polyfill

// Filter returns elements that satisfy the predicate (like JS array.filter)
//
// Example:
//
//	adults := polyfill.From(people).
//	    Filter(func(p Person) bool { return p.Age >= 18 }).
//	    Slice()
func (s *Seq[T]) Filter(f func(T) bool) *Seq[T] {
	if s.err != nil {
		return s
	}

	result := make([]T, 0, len(s.elements))
	for _, v := range s.elements {
		if f(v) {
			result = append(result, v)
		}
	}
	return From(result)
}
