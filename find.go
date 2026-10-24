package polyfill

// Find returns the first element that satisfies predicate.
// It returns false when nothing matches or the chain failed.
//
//	if user, ok := From(users).Find(isAdmin); ok { ... }
func (s *Seq[T]) Find(predicate func(T) bool) (T, bool) {
	var zero T
	if s == nil || s.err != nil {
		return zero, false
	}
	for _, v := range s.elements {
		if predicate(v) {
			return v, true
		}
	}
	return zero, false
}
