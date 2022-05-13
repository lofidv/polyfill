package polyfill

type Constraint interface {
	any
}

type Slice[V Constraint] []V

func (s Slice[V]) Len() int {
	return len(s)
}

func (s *Slice[V]) Add(items ...V) Slice[V] {
	for _, i := range items {
		*s = append(*s, i)
	}

	return *s
}
