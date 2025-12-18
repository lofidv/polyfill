package polyfill

// Reduce reduces the sequence to a single value (like JS array.reduce)
// The accumulator and return value are both of type R
//
// Example:
//
//	sum := polyfill.From([]int{1, 2, 3, 4}).
//	    Reduce(0, func(acc int, val int) int { return acc + val })
//	// Result: 10
func Reduce[T any, R any](s *Seq[T], initial R, f func(acc R, val T) R) R {
	acc := initial
	for _, v := range s.elements {
		acc = f(acc, v)
	}
	return acc
}

// ReduceE reduces with error handling
func ReduceE[T any, R any](s *Seq[T], initial R, f func(acc R, val T) (R, error)) (R, error) {
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
func ReduceRight[T any, R any](s *Seq[T], initial R, f func(acc R, val T) R) R {
	acc := initial
	for i := len(s.elements) - 1; i >= 0; i-- {
		acc = f(acc, s.elements[i])
	}
	return acc
}
