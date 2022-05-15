package polyfill

func (s Slice[V, T]) Reduce(acc func(T, V) T, init T) T {
	r := init
	for _, v := range s {
		r = acc(r, v)
	}

	return r
}
