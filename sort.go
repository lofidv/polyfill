package polyfill

import (
	"cmp"
	"slices"
)

// Sort returns a sorted copy of the sequence. The receiver is not modified.
//
//	From([]int{3, 1, 2}).Sort(func(a, b int) bool { return a < b }).Slice() // [1 2 3]
func (s *Seq[T]) Sort(less func(a, b T) bool) *Seq[T] {
	if s == nil || s.err != nil {
		return s
	}
	out := slices.Clone(s.elements)
	slices.SortFunc(out, func(a, b T) int {
		if less(a, b) {
			return -1
		}
		if less(b, a) {
			return 1
		}
		return 0
	})
	return From(out)
}

// SortBy returns a sorted copy ordered by a comparable key.
//
// This is the version most code wants: sorting by a field without writing a
// comparison function.
//
//	SortBy(From(users), func(u User) string { return u.Name }).Slice()
//	SortBy(From(orders), func(o Order) time.Time { return o.CreatedAt }).Slice()
func SortBy[T any, K cmp.Ordered](s *Seq[T], key func(T) K) *Seq[T] {
	if s == nil || s.err != nil {
		return s
	}
	out := slices.Clone(s.elements)
	slices.SortFunc(out, func(a, b T) int {
		return cmp.Compare(key(a), key(b))
	})
	return From(out)
}
