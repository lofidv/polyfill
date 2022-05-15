package polyfill

func (s Slice[V, T]) IndexOf(predicate func(V, int) bool, pos int, start int) int {
	for i := start; i < len(s); i++ {
		if predicate(s[i], pos) {
			return i
		}
	}

	return -1
}
