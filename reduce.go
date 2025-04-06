package polyfill

// Reduce reduces the slice to a single value
func (s *Slice[T]) Reduce(initial any, f func(acc any, val T) any) any {
	acc := initial
	for _, v := range s.elements {
		acc = f(acc, v)
	}
	return acc
}
