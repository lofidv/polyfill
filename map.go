package polyfill

// === TRANSFORMATION METHODS ===

// Map transforms each element (same type T -> T)
func (s *Seq[T]) Map(f func(T) T) *Seq[T] {
	if s.err != nil {
		return s
	}

	result := make([]T, len(s.elements))
	for i, v := range s.elements {
		result[i] = f(v)
	}
	return From(result)
}

// MapE transforms elements with error handling (same type)
func (s *Seq[T]) MapE(f func(T) (T, error)) *Seq[T] {
	if s.err != nil {
		return s
	}

	result := make([]T, len(s.elements))
	for i, v := range s.elements {
		var err error
		result[i], err = f(v)
		if err != nil {
			return &Seq[T]{err: err}
		}
	}
	return From(result)
}

// MapTo transforms elements with type change (T -> R)
// Must be a global function due to Go generic limitations
func MapTo[T any, R any](s *Seq[T], f func(T) R) *Seq[R] {
	if s.err != nil {
		return &Seq[R]{err: s.err}
	}

	result := make([]R, len(s.elements))
	for i, v := range s.elements {
		result[i] = f(v)
	}
	return From(result)
}

// MapToE transforms elements with type change and error handling
func MapToE[T any, R any](s *Seq[T], f func(T) (R, error)) *Seq[R] {
	if s.err != nil {
		return &Seq[R]{err: s.err}
	}

	result := make([]R, len(s.elements))
	for i, v := range s.elements {
		var err error
		result[i], err = f(v)
		if err != nil {
			return &Seq[R]{err: err}
		}
	}
	return From(result)
}
