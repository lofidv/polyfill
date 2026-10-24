package polyfill

// Seq is a chainable wrapper around a Go slice.
//
// A Seq may carry an error. See the package documentation for the error
// contract: once the error is set, every operation is a no-op and the terminal
// operations return zero values.
type Seq[T any] struct {
	elements []T
	err      error
}

// From creates a Seq from an existing slice.
// It does not copy: the Seq shares the caller's backing array.
//
//	From([]int{1, 2, 3}).Filter(func(n int) bool { return n > 1 }).Slice()
func From[T any](items []T) *Seq[T] {
	return &Seq[T]{elements: items}
}

// New creates a Seq from variadic arguments.
//
//	New(1, 2, 3).Map(func(n int) int { return n * 2 }).Slice()
func New[T any](items ...T) *Seq[T] {
	return From(items)
}

// Slice returns the underlying slice. It is the primary way to exit a chain.
//
// When the chain failed, Slice returns the elements processed before the
// failure. Use SliceE when you need the error.
func (s *Seq[T]) Slice() []T {
	if s == nil {
		return nil
	}
	return s.elements
}

// SliceE returns the elements processed so far and the chain error.
//
//	rows, err := From(raw).MapE(parse).Filter(valid).SliceE()
func (s *Seq[T]) SliceE() ([]T, error) {
	if s == nil {
		return nil, nil
	}
	return s.elements, s.err
}

// Err returns the error that stopped the chain, or nil.
func (s *Seq[T]) Err() error {
	if s == nil {
		return nil
	}
	return s.err
}

// Len returns the number of elements. On a failed chain it returns the number
// of elements processed before the failure.
func (s *Seq[T]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.elements)
}

// IsEmpty reports whether the sequence has no elements.
func (s *Seq[T]) IsEmpty() bool {
	return s == nil || len(s.elements) == 0
}

// Push appends items to the sequence and returns the receiver.
//
// Push mutates, like JavaScript's Array.push. Use Append for the immutable
// version that returns a new Seq.
func (s *Seq[T]) Push(items ...T) *Seq[T] {
	if s == nil {
		return &Seq[T]{elements: items}
	}
	if s.err != nil {
		return s
	}
	s.elements = append(s.elements, items...)
	return s
}

// At returns the element at index, like JavaScript's Array.at.
// A negative index counts from the end: -1 is the last element.
// It returns false when the index is out of bounds or the chain failed.
func (s *Seq[T]) At(index int) (T, bool) {
	var zero T
	if s == nil || s.err != nil {
		return zero, false
	}
	if index < 0 {
		index += len(s.elements)
	}
	if index < 0 || index >= len(s.elements) {
		return zero, false
	}
	return s.elements[index], true
}

// Get returns the element at index. Negative indices are not supported; use At
// for that. It returns false when the index is out of bounds or the chain
// failed.
func (s *Seq[T]) Get(index int) (T, bool) {
	var zero T
	if s == nil || s.err != nil || index < 0 || index >= len(s.elements) {
		return zero, false
	}
	return s.elements[index], true
}

// First returns the first element.
func (s *Seq[T]) First() (T, bool) {
	var zero T
	if s == nil || s.err != nil || len(s.elements) == 0 {
		return zero, false
	}
	return s.elements[0], true
}

// Last returns the last element.
func (s *Seq[T]) Last() (T, bool) {
	var zero T
	if s == nil || s.err != nil || len(s.elements) == 0 {
		return zero, false
	}
	return s.elements[len(s.elements)-1], true
}

// Take returns a new Seq with the first n elements.
// It copies, so mutating the result cannot affect the receiver.
func (s *Seq[T]) Take(n int) *Seq[T] {
	if s == nil {
		return &Seq[T]{}
	}
	if s.err != nil {
		return s
	}
	if n <= 0 {
		return From([]T{})
	}
	if n >= len(s.elements) {
		n = len(s.elements)
	}
	out := make([]T, n)
	copy(out, s.elements[:n])
	return From(out)
}

// Skip returns a new Seq without the first n elements.
// It copies, so mutating the result cannot affect the receiver.
func (s *Seq[T]) Skip(n int) *Seq[T] {
	if s == nil {
		return &Seq[T]{}
	}
	if s.err != nil {
		return s
	}
	if n <= 0 {
		n = 0
	}
	if n >= len(s.elements) {
		return From([]T{})
	}
	out := make([]T, len(s.elements)-n)
	copy(out, s.elements[n:])
	return From(out)
}

// ForEach calls f for every element.
func (s *Seq[T]) ForEach(f func(T)) {
	if s == nil || s.err != nil {
		return
	}
	for _, v := range s.elements {
		f(v)
	}
}

// ForEachIndexed calls f with the index and the element.
func (s *Seq[T]) ForEachIndexed(f func(int, T)) {
	if s == nil || s.err != nil {
		return
	}
	for i, v := range s.elements {
		f(i, v)
	}
}

// MinBy returns the minimum element according to less.
// It returns false when the sequence is empty or the chain failed.
func (s *Seq[T]) MinBy(less func(a, b T) bool) (T, bool) {
	var zero T
	if s == nil || s.err != nil || len(s.elements) == 0 {
		return zero, false
	}
	minVal := s.elements[0]
	for i := 1; i < len(s.elements); i++ {
		if less(s.elements[i], minVal) {
			minVal = s.elements[i]
		}
	}
	return minVal, true
}

// MaxBy returns the maximum element according to less.
// It returns false when the sequence is empty or the chain failed.
func (s *Seq[T]) MaxBy(less func(a, b T) bool) (T, bool) {
	var zero T
	if s == nil || s.err != nil || len(s.elements) == 0 {
		return zero, false
	}
	maxVal := s.elements[0]
	for i := 1; i < len(s.elements); i++ {
		if less(maxVal, s.elements[i]) {
			maxVal = s.elements[i]
		}
	}
	return maxVal, true
}
