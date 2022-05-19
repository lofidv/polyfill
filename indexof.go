package polyfill

func (s Slice[V, T]) IndexOf(predicate func(V, T) bool, criteria T, start int) int {
	for i := start; i < len(s); i++ {
		if predicate(s[i], criteria) {
			return i
		}
	}

	return -1
}
