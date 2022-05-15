package polyfill

type Slice[V, T any] []V

func (s Slice[V, T]) Len() int {
	return len(s)
}

func (s *Slice[V, T]) Add(items ...V) Slice[V, T] {
	*s = append(*s, items...)
	return *s
}

func NewSlice[V, T any](items ...V) Slice[V, T] {
	return append([]V{}, items...)
}
