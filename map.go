package polyfill

// Map transforms each element, keeping the element type.
//
//	From([]int{1, 2, 3}).Map(func(n int) int { return n * 2 }).Slice() // [2 4 6]
func (s *Seq[T]) Map(f func(T) T) *Seq[T] {
	if s == nil || s.err != nil {
		return s
	}
	result := make([]T, len(s.elements))
	for i, v := range s.elements {
		result[i] = f(v)
	}
	return From(result)
}

// MapE transforms each element and stops at the first error.
//
// The elements transformed before the failure are kept, so SliceE returns a
// partial result together with the error. This is the method to use when a
// pipeline parses or validates external data.
//
//	rows, err := From(raw).MapE(parseRow).SliceE()
func (s *Seq[T]) MapE(f func(T) (T, error)) *Seq[T] {
	if s == nil || s.err != nil {
		return s
	}
	result := make([]T, 0, len(s.elements))
	for _, v := range s.elements {
		val, err := f(v)
		if err != nil {
			return &Seq[T]{elements: result, err: err}
		}
		result = append(result, val)
	}
	return From(result)
}

// MapTo transforms elements into another type.
//
// MapTo is a free function because Go 1.18 to 1.26 do not allow a method to
// declare a type parameter of its own.
//
//	MapTo(From([]int{1, 2}), strconv.Itoa).Slice() // ["1" "2"]
func MapTo[T any, R any](s *Seq[T], f func(T) R) *Seq[R] {
	if s == nil || s.err != nil {
		return &Seq[R]{err: errOf(s)}
	}
	result := make([]R, len(s.elements))
	for i, v := range s.elements {
		result[i] = f(v)
	}
	return From(result)
}

// MapToE is MapTo with error handling. It stops at the first error and keeps
// the elements produced so far.
func MapToE[T any, R any](s *Seq[T], f func(T) (R, error)) *Seq[R] {
	if s == nil || s.err != nil {
		return &Seq[R]{err: errOf(s)}
	}
	result := make([]R, 0, len(s.elements))
	for _, v := range s.elements {
		val, err := f(v)
		if err != nil {
			return &Seq[R]{elements: result, err: err}
		}
		result = append(result, val)
	}
	return From(result)
}

func errOf[T any](s *Seq[T]) error {
	if s == nil {
		return nil
	}
	return s.err
}
