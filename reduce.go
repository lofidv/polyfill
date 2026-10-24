package polyfill

// Reduce folds the sequence into a single value of the same type.
//
//	From([]int{1, 2, 3}).Reduce(0, func(acc, n int) int { return acc + n }) // 6
func (s *Seq[T]) Reduce(initial T, f func(acc, val T) T) T {
	if s == nil || s.err != nil {
		return initial
	}
	acc := initial
	for _, v := range s.elements {
		acc = f(acc, v)
	}
	return acc
}

// ReduceTo folds the sequence into a value of a different type.
//
//	ReduceTo(From([]int{1, 2}), "", func(acc string, n int) string {
//		return acc + strconv.Itoa(n)
//	}) // "12"
func ReduceTo[T any, R any](s *Seq[T], initial R, f func(acc R, val T) R) R {
	if s == nil || s.err != nil {
		return initial
	}
	acc := initial
	for _, v := range s.elements {
		acc = f(acc, v)
	}
	return acc
}

// ReduceE folds with error handling and stops at the first error.
func (s *Seq[T]) ReduceE(initial T, f func(acc, val T) (T, error)) (T, error) {
	if s == nil {
		return initial, nil
	}
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

// ReduceToE folds into another type with error handling and stops at the first
// error.
func ReduceToE[T any, R any](s *Seq[T], initial R, f func(acc R, val T) (R, error)) (R, error) {
	if s == nil {
		return initial, nil
	}
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

// ReduceRight folds from the last element to the first.
func (s *Seq[T]) ReduceRight(initial T, f func(acc, val T) T) T {
	if s == nil || s.err != nil {
		return initial
	}
	acc := initial
	for i := len(s.elements) - 1; i >= 0; i-- {
		acc = f(acc, s.elements[i])
	}
	return acc
}

// ReduceRightTo folds from the last element to the first into another type.
func ReduceRightTo[T any, R any](s *Seq[T], initial R, f func(acc R, val T) R) R {
	if s == nil || s.err != nil {
		return initial
	}
	acc := initial
	for i := len(s.elements) - 1; i >= 0; i-- {
		acc = f(acc, s.elements[i])
	}
	return acc
}
