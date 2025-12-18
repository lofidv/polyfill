package polyfill

// FlatMap applies a function to each element and flattens the result (like JS array.flatMap)
// Useful when the mapping function returns a slice for each element
//
// Example:
//
//	words := polyfill.From([]string{"hello world", "foo bar"}).
//	    FlatMap(func(s string) []string { return strings.Split(s, " ") }).
//	    Slice()
//	// Result: ["hello", "world", "foo", "bar"]
func FlatMap[T any, R any](s *Seq[T], f func(T) []R) *Seq[R] {
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
func FlatMapE[T any, R any](s *Seq[T], f func(T) ([]R, error)) *Seq[R] {
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
//	flat := polyfill.From([][]int{{1, 2}, {3, 4}, {5}}).Flatten().Slice()
//	// Result: [1, 2, 3, 4, 5]
func Flatten[T any](s *Seq[[]T]) *Seq[T] {
	result := make([]T, 0)
	for _, slice := range s.elements {
		result = append(result, slice...)
	}
	return From(result)
}
