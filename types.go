package polyfill

// Slice represents a functional wrapper around a Go slice
type Slice[T any] struct {
	elements []T
}

// Wrap creates a new Slice from an existing slice
func Wrap[T any](items []T) *Slice[T] {
	return &Slice[T]{elements: items}
}

// New creates a new Slice from variadic arguments
func New[T any](items ...T) *Slice[T] {
	return Wrap(items)
}

// Unwrap returns the underlying slice
func (s *Slice[T]) Unwrap() []T {
	return s.elements
}

// Len returns the length of the slice
func (s *Slice[T]) Len() int {
	return len(s.elements)
}

// Add appends items to the slice
func (s *Slice[T]) Add(items ...T) *Slice[T] {
	s.elements = append(s.elements, items...)
	return s
}
