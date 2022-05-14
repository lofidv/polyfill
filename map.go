package polyfill

func (s Slice[V, T]) Map(f func(V) T) Slice[T, V] {
	slice := make([]T, len(s))
	for i, v := range s {
		slice[i] = f(v)
	}
	return slice
}
