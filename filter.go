package polyfill

func (s Slice[V]) Filter(f func(V) bool) []V {
	slice := Slice[V]{}
	for _, t := range s {
		if f(t) {
			slice = append(slice, t)
		}
	}

	return slice
}
