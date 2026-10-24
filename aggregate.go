package polyfill

// Number is the set of numeric types accepted by SumBy and AvgBy.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// SumBy adds a numeric field over the sequence.
//
//	revenue := SumBy(From(orders), func(o Order) float64 { return o.Total })
//	units := SumBy(From(orders), func(o Order) int { return o.Quantity })
func SumBy[T any, N Number](s *Seq[T], value func(T) N) N {
	var total N
	if s == nil || s.err != nil {
		return total
	}
	for _, v := range s.elements {
		total += value(v)
	}
	return total
}

// AvgBy returns the average of a numeric field, or 0 for an empty sequence.
//
//	avgTicket := AvgBy(From(orders), func(o Order) float64 { return o.Total })
func AvgBy[T any, N Number](s *Seq[T], value func(T) N) float64 {
	if s == nil || s.err != nil || len(s.elements) == 0 {
		return 0
	}
	var total float64
	for _, v := range s.elements {
		total += float64(value(v))
	}
	return total / float64(len(s.elements))
}
