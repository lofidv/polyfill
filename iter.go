package polyfill

import "iter"

// Iter exposes the sequence as a Go 1.23 iterator, so it can feed any API that
// accepts an iter.Seq: standard library helpers, range-over-func loops, or
// another library.
//
//	for u := range From(users).Iter() {
//		fmt.Println(u.Name)
//	}
func (s *Seq[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		if s == nil || s.err != nil {
			return
		}
		for _, v := range s.elements {
			if !yield(v) {
				return
			}
		}
	}
}

// FromIter collects a Go 1.23 iterator into a Seq.
//
//	seq := FromIter(slices.Values([]int{1, 2, 3}))
//	seq := FromIter(maps.Values(index))
func FromIter[T any](source iter.Seq[T]) *Seq[T] {
	var items []T
	for v := range source {
		items = append(items, v)
	}
	return From(items)
}
