package polyfill

// IndexOf returns the index of the first matching element or -1
func (s *Slice[T]) IndexOf(f func(T) bool) int {
	for i, v := range s.elements {
		if f(v) {
			return i
		}
	}
	return -1
}
