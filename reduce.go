package polyfill

// Reduce reduces the sequence to a single value (like JS array.reduce)
//
// Example:
//
//	From([]int{1, 2, 3}).Reduce(0, func(a, n int) int { return a + n }) // 6
func (s *Seq[T]) Reduce(initial T, f func(acc T, val T) T) T {
	acc := initial
	for _, v := range s.elements {
		acc = f(acc, v)
	}
	return acc
}

// ReduceTo reduces the sequence to a single value of a different type
//
// Example:
//
//	ReduceTo(From([]int{1, 2}), "", func(a string, n int) string { return a + fmt.Sprint(n) }) // "12"
func ReduceTo[T any, R any](s *Seq[T], initial R, f func(acc R, val T) R) R {
	acc := initial
	for _, v := range s.elements {
		acc = f(acc, v)
	}
	return acc
}

// ReduceE reduces with error handling (same type)
func (s *Seq[T]) ReduceE(initial T, f func(acc T, val T) (T, error)) (T, error) {
	if s.err != nil {
		return initial, s.err
	}

	acc := initial
	for _, v := range s.elements {
		var err error
		acc, err = f(acc, v)
		if err != nil {
			return acc, err
		}
	}
	return acc, nil
}

// ReduceToE reduces with error handling and type change
func ReduceToE[T any, R any](s *Seq[T], initial R, f func(acc R, val T) (R, error)) (R, error) {
	if s.err != nil {
		return initial, s.err
	}

	acc := initial
	for _, v := range s.elements {
		var err error
		acc, err = f(acc, v)
		if err != nil {
			return acc, err
		}
	}
	return acc, nil
}

// ReduceRight reduces the sequence from right to left
func (s *Seq[T]) ReduceRight(initial T, f func(acc T, val T) T) T {
	acc := initial
	for i := len(s.elements) - 1; i >= 0; i-- {
		acc = f(acc, s.elements[i])
	}
	return acc
}

// ReduceRightTo reduces from right to left with type change
func ReduceRightTo[T any, R any](s *Seq[T], initial R, f func(acc R, val T) R) R {
	acc := initial
	for i := len(s.elements) - 1; i >= 0; i-- {
		acc = f(acc, s.elements[i])
	}
	return acc
}
