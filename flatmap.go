package polyfill

// FlatMap maps each element to a slice and concatenates the results.
//
//	From([]int{1, 2}).FlatMap(func(n int) []int { return []int{n, n} }).Slice() // [1 1 2 2]
func (s *Seq[T]) FlatMap(f func(T) []T) *Seq[T] {
	if s == nil || s.err != nil {
		return s
	}
	result := make([]T, 0, len(s.elements))
	for _, v := range s.elements {
		result = append(result, f(v)...)
	}
	return From(result)
}

// FlatMapTo maps to another type and concatenates the results.
func FlatMapTo[T any, R any](s *Seq[T], f func(T) []R) *Seq[R] {
	if s == nil || s.err != nil {
		return &Seq[R]{err: errOf(s)}
	}
	result := make([]R, 0, len(s.elements))
	for _, v := range s.elements {
		result = append(result, f(v)...)
	}
	return From(result)
}

// FlatMapE is FlatMap with error handling. It stops at the first error and
// keeps the elements produced so far.
func (s *Seq[T]) FlatMapE(f func(T) ([]T, error)) *Seq[T] {
	if s == nil || s.err != nil {
		return s
	}
	result := make([]T, 0, len(s.elements))
	for _, v := range s.elements {
		items, err := f(v)
		if err != nil {
			return &Seq[T]{elements: result, err: err}
		}
		result = append(result, items...)
	}
	return From(result)
}

// FlatMapToE is FlatMapTo with error handling. It stops at the first error and
// keeps the elements produced so far.
func FlatMapToE[T any, R any](s *Seq[T], f func(T) ([]R, error)) *Seq[R] {
	if s == nil || s.err != nil {
		return &Seq[R]{err: errOf(s)}
	}
	result := make([]R, 0, len(s.elements))
	for _, v := range s.elements {
		items, err := f(v)
		if err != nil {
			return &Seq[R]{elements: result, err: err}
		}
		result = append(result, items...)
	}
	return From(result)
}

// Flatten concatenates a sequence of slices into one sequence.
//
//	Flatten(From([][]int{{1, 2}, {3}})).Slice() // [1 2 3]
func Flatten[T any](s *Seq[[]T]) *Seq[T] {
	if s == nil || s.err != nil {
		return &Seq[T]{err: errOf(s)}
	}
	total := 0
	for _, slice := range s.elements {
		total += len(slice)
	}
	result := make([]T, 0, total)
	for _, slice := range s.elements {
		result = append(result, slice...)
	}
	return From(result)
}
