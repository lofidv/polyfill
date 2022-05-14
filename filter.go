package polyfill

func (s Slice[V, T]) Filter(f func(V) bool) Slice[V, T] {
	slice := Slice[V, T]{}
	for _, t := range s {
		if f(t) {
			slice = append(slice, t)
		}
	}

	return slice
}
