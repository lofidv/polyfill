package polyfill

// FlatMap applies a function to each element and flattens the result (like JS array.flatMap)
//
// Example:
//
//	From([]int{1, 2}).FlatMap(func(n int) []int { return []int{n, n} }).Slice() // [1, 1, 2, 2]
func (s *Seq[T]) FlatMap(f func(T) []T) *Seq[T] {
	if s.err != nil {
		return s
	}

	result := make([]T, 0)
	for _, v := range s.elements {
		result = append(result, f(v)...)
	}
	return From(result)
}

// FlatMapTo applies a function with type change and flattens the result
//
// Example:
//
//	FlatMapTo(From([]int{1, 2}), func(n int) []string { return []string{fmt.Sprint(n)} }).Slice()
func FlatMapTo[T any, R any](s *Seq[T], f func(T) []R) *Seq[R] {
	if s.err != nil {
		return &Seq[R]{err: s.err}
	}

	result := make([]R, 0)
	for _, v := range s.elements {
		result = append(result, f(v)...)
	}
	return From(result)
}

// FlatMapE applies a function with error handling and flattens the result
func (s *Seq[T]) FlatMapE(f func(T) ([]T, error)) *Seq[T] {
	if s.err != nil {
		return s
	}

	result := make([]T, 0)
	for _, v := range s.elements {
		items, err := f(v)
		if err != nil {
			return &Seq[T]{err: err}
		}
		result = append(result, items...)
	}
	return From(result)
}

// FlatMapToE applies a function with error handling and type change, then flattens
func FlatMapToE[T any, R any](s *Seq[T], f func(T) ([]R, error)) *Seq[R] {
	if s.err != nil {
		return &Seq[R]{err: s.err}
	}

	result := make([]R, 0)
	for _, v := range s.elements {
		items, err := f(v)
		if err != nil {
			return &Seq[R]{err: err}
		}
		result = append(result, items...)
	}
	return From(result)
}

// Flatten flattens a sequence of slices into a single sequence
//
// Example:
//
//	Flatten(From([][]int{{1, 2}, {3}})).Slice() // [1, 2, 3]
func Flatten[T any](s *Seq[[]T]) *Seq[T] {
	result := make([]T, 0)
	for _, slice := range s.elements {
		result = append(result, slice...)
	}
	return From(result)
}
