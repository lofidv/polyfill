package polyfill

import "slices"

// Concat returns a new sequence with the given slices appended.
//
//	From([]int{1, 2}).Concat([]int{3, 4}).Slice() // [1 2 3 4]
func (s *Seq[T]) Concat(others ...[]T) *Seq[T] {
	if s == nil || s.err != nil {
		return s
	}
	total := len(s.elements)
	for _, other := range others {
		total += len(other)
	}
	result := make([]T, 0, total)
	result = append(result, s.elements...)
	for _, other := range others {
		result = append(result, other...)
	}
	return From(result)
}

// Append returns a new sequence with the items appended. It never modifies the
// receiver. Use Push when you want the mutating version.
//
//	From([]int{1, 2}).Append(3, 4).Slice() // [1 2 3 4]
func (s *Seq[T]) Append(items ...T) *Seq[T] {
	if s == nil {
		return From(slices.Clone(items))
	}
	if s.err != nil {
		return s
	}
	result := make([]T, 0, len(s.elements)+len(items))
	result = append(result, s.elements...)
	result = append(result, items...)
	return From(result)
}

// Prepend returns a new sequence with the items inserted at the beginning.
//
//	From([]int{3, 4}).Prepend(1, 2).Slice() // [1 2 3 4]
func (s *Seq[T]) Prepend(items ...T) *Seq[T] {
	if s == nil {
		return From(slices.Clone(items))
	}
	if s.err != nil {
		return s
	}
	result := make([]T, 0, len(items)+len(s.elements))
	result = append(result, items...)
	result = append(result, s.elements...)
	return From(result)
}

// EqualFunc reports whether both slices have the same length and every pair of
// elements is equal according to eq.
//
//	From([]int{1, 2}).EqualFunc([]int{1, 2}, func(a, b int) bool { return a == b }) // true
func (s *Seq[T]) EqualFunc(other []T, eq func(T, T) bool) bool {
	if s == nil || s.err != nil {
		return false
	}
	if len(s.elements) != len(other) {
		return false
	}
	for i := range s.elements {
		if !eq(s.elements[i], other[i]) {
			return false
		}
	}
	return true
}

// Clone returns a sequence with its own copy of the backing array.
func (s *Seq[T]) Clone() *Seq[T] {
	if s == nil {
		return &Seq[T]{}
	}
	if s.err != nil {
		return s
	}
	return From(slices.Clone(s.elements))
}
