package polyfill

// === TRANSFORMATION METHODS ===

// Map transforms each element (same type T -> T)
//
// Example:
//
//	From([]int{1, 2, 3}).Map(func(n int) int { return n * 2 }).Slice() // [2, 4, 6]
func (s *Seq[T]) Map(f func(T) T) *Seq[T] {
	if s.err != nil {
		return s
	}

	result := make([]T, 0, len(s.elements))
	for _, v := range s.elements {
		result = append(result, f(v))
	}
	return From(result)
}

// MapE transforms elements with error handling (same type)
func (s *Seq[T]) MapE(f func(T) (T, error)) *Seq[T] {
	if s.err != nil {
		return s
	}

	result := make([]T, 0, len(s.elements))
	for _, v := range s.elements {
		val, err := f(v)
		if err != nil {
			return &Seq[T]{err: err}
		}
		result = append(result, val)
	}
	return From(result)
}

// MapTo transforms elements with type change (T -> R)
// Must remain a global function due to Go generic method limitations
//
// Example:
//
//	MapTo(From([]int{1, 2}), func(n int) string { return fmt.Sprint(n) }).Slice() // ["1", "2"]
func MapTo[T any, R any](s *Seq[T], f func(T) R) *Seq[R] {
	if s.err != nil {
		return &Seq[R]{err: s.err}
	}

	result := make([]R, 0, len(s.elements))
	for _, v := range s.elements {
		result = append(result, f(v))
	}
	return From(result)
}

// MapToE transforms elements with type change and error handling
func MapToE[T any, R any](s *Seq[T], f func(T) (R, error)) *Seq[R] {
	if s.err != nil {
		return &Seq[R]{err: s.err}
	}

	result := make([]R, 0, len(s.elements))
	for _, v := range s.elements {
		val, err := f(v)
		if err != nil {
			return &Seq[R]{err: err}
		}
		result = append(result, val)
	}
	return From(result)
}
