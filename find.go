package polyfill

func (s Slice[V, T]) Find(f func(V) bool) V {
	var slice V
	for _, t := range s {
		if f(t) {
			slice = t
		}
	}

	return slice
}
