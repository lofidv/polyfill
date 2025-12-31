package polyfill

import "slices"

// Concat concatenates multiple sequences into one (like JS array.concat)
//
// Example:
//
//	From([]int{1, 2}).Concat([]int{3, 4}).Slice() // [1, 2, 3, 4]
func (s *Seq[T]) Concat(others ...[]T) *Seq[T] {
	if s.err != nil {
		return s
	}

	result := slices.Clone(s.elements)
	for _, other := range others {
		result = append(result, other...)
	}
	return From(result)
}

// Append appends elements to the sequence (returns new sequence, immutable)
//
// Example:
//
//	From([]int{1, 2}).Append(3, 4).Slice() // [1, 2, 3, 4]
func (s *Seq[T]) Append(items ...T) *Seq[T] {
	if s.err != nil {
		return s
	}

	result := slices.Clone(s.elements)
	result = append(result, items...)
	return From(result)
}

// Prepend adds elements to the beginning (returns new sequence, immutable)
//
// Example:
//
//	From([]int{3, 4}).Prepend(1, 2).Slice() // [1, 2, 3, 4]
func (s *Seq[T]) Prepend(items ...T) *Seq[T] {
	if s.err != nil {
		return s
	}

	result := make([]T, 0, len(items)+len(s.elements))
	result = append(result, items...)
	result = append(result, s.elements...)
	return From(result)
}

// ContainsFunc checks if the sequence contains a value matching the predicate
//
// Example:
//
//	From([]int{1, 2, 3}).ContainsFunc(func(n int) bool { return n == 2 }) // true
func (s *Seq[T]) ContainsFunc(predicate func(T) bool) bool {
	for _, v := range s.elements {
		if predicate(v) {
			return true
		}
	}
	return false
}

// EqualFunc checks if two sequences are equal using a comparison function
//
// Example:
//
//	From([]int{1, 2}).EqualFunc([]int{1, 2}, func(a, b int) bool { return a == b }) // true
func (s *Seq[T]) EqualFunc(other []T, eq func(T, T) bool) bool {
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

// Clone creates a shallow copy of the sequence
//
// Example:
//
//	copy := From([]int{1, 2, 3}).Clone()
func (s *Seq[T]) Clone() *Seq[T] {
	return From(slices.Clone(s.elements))
}
