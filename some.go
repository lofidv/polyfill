package polyfill

func (s Slice[V, T]) Some(predicate func(V) bool) bool {
	for _, slice := range s {
		if predicate(slice) {
			return true
		}
	}

	return false
}
