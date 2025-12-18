package polyfill

// Seq represents a functional sequence wrapper around a Go slice
// providing a chainable, fluent API for slice operations inspired by JavaScript
type Seq[T any] struct {
	elements []T
	err      error // stores error for chainable error handling
}

// From creates a new Seq from an existing slice
// This is the primary way to start a chain of operations
//
// Example:
//
//	polyfill.From([]int{1, 2, 3}).Map(fn).Slice()
func From[T any](items []T) *Seq[T] {
	return &Seq[T]{elements: items}
}

// New creates a new Seq from variadic arguments
//
// Example:
//
//	polyfill.New(1, 2, 3, 4, 5).Filter(fn).Slice()
func New[T any](items ...T) *Seq[T] {
	return From(items)
}

// Slice returns the underlying slice
// This is the primary way to exit a chain of operations
func (s *Seq[T]) Slice() []T {
	return s.elements
}

// SliceE returns the underlying slice and any error that occurred during chaining
func (s *Seq[T]) SliceE() ([]T, error) {
	return s.elements, s.err
}

// Err returns any error that occurred during the chain
func (s *Seq[T]) Err() error {
	return s.err
}

// Len returns the length of the sequence
func (s *Seq[T]) Len() int {
	return len(s.elements)
}

// IsEmpty returns true if the sequence has no elements
func (s *Seq[T]) IsEmpty() bool {
	return len(s.elements) == 0
}

// Push appends items to the sequence (mutates)
func (s *Seq[T]) Push(items ...T) *Seq[T] {
	s.elements = append(s.elements, items...)
	return s
}

// At returns the element at the given index (like JS array.at)
// Returns zero value and false if index is out of bounds
// Supports negative indices: -1 is last element
func (s *Seq[T]) At(index int) (T, bool) {
	if index < 0 {
		index = len(s.elements) + index
	}
	if index < 0 || index >= len(s.elements) {
		var zero T
		return zero, false
	}
	return s.elements[index], true
}

// Get returns the element at the given index (alias for At with positive indices only)
func (s *Seq[T]) Get(index int) (T, bool) {
	if index < 0 || index >= len(s.elements) {
		var zero T
		return zero, false
	}
	return s.elements[index], true
}

// First returns the first element
// Returns zero value and false if sequence is empty
func (s *Seq[T]) First() (T, bool) {
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	return s.elements[0], true
}

// Last returns the last element
// Returns zero value and false if sequence is empty
func (s *Seq[T]) Last() (T, bool) {
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	return s.elements[len(s.elements)-1], true
}

// Take returns a new Seq with the first n elements (like JS slice)
func (s *Seq[T]) Take(n int) *Seq[T] {
	if n <= 0 {
		return From([]T{})
	}
	if n >= len(s.elements) {
		return From(s.elements)
	}
	return From(s.elements[:n])
}

// Skip returns a new Seq with the first n elements removed
func (s *Seq[T]) Skip(n int) *Seq[T] {
	if n <= 0 {
		return From(s.elements)
	}
	if n >= len(s.elements) {
		return From([]T{})
	}
	return From(s.elements[n:])
}

// ForEach executes a function for each element (like JS forEach)
func (s *Seq[T]) ForEach(f func(T)) {
	for _, v := range s.elements {
		f(v)
	}
}

// ForEachIndexed executes a function for each element with its index
func (s *Seq[T]) ForEachIndexed(f func(int, T)) {
	for i, v := range s.elements {
		f(i, v)
	}
}
