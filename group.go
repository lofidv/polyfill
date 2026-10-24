package polyfill

// GroupBy groups elements by a comparable key.
//
// The keys are typed, so callers do not need type assertions, and the order
// inside each group is the input order.
//
//	byRegion := GroupBy(From(orders), func(o Order) string { return o.Region })
//	for region, orders := range byRegion { ... }
func GroupBy[T any, K comparable](s *Seq[T], key func(T) K) map[K][]T {
	result := make(map[K][]T)
	if s == nil || s.err != nil {
		return result
	}
	for _, v := range s.elements {
		k := key(v)
		result[k] = append(result[k], v)
	}
	return result
}

// CountBy counts elements per key. It is GroupBy plus len, for reports.
//
//	totals := CountBy(From(orders), func(o Order) string { return o.Status })
func CountBy[T any, K comparable](s *Seq[T], key func(T) K) map[K]int {
	result := make(map[K]int)
	if s == nil || s.err != nil {
		return result
	}
	for _, v := range s.elements {
		result[key(v)]++
	}
	return result
}

// KeyBy builds an index from key to element. When two elements share a key,
// the last one wins.
//
//	byID := KeyBy(From(users), func(u User) int { return u.ID })
func KeyBy[T any, K comparable](s *Seq[T], key func(T) K) map[K]T {
	result := make(map[K]T, len(s.Slice()))
	if s == nil || s.err != nil {
		return result
	}
	for _, v := range s.elements {
		result[key(v)] = v
	}
	return result
}

// ToMap builds a map from a key function and a value function.
//
//	byID := ToMap(From(users), func(u User) int { return u.ID }, func(u User) string { return u.Email })
func ToMap[T any, K comparable, V any](s *Seq[T], key func(T) K, value func(T) V) map[K]V {
	result := make(map[K]V, len(s.Slice()))
	if s == nil || s.err != nil {
		return result
	}
	for _, v := range s.elements {
		result[key(v)] = value(v)
	}
	return result
}

// Partition splits the sequence in two: elements that satisfy predicate and
// elements that do not. Both parts preserve input order and can chain.
//
//	paid, unpaid := From(orders).Partition(func(o Order) bool { return o.Paid })
func (s *Seq[T]) Partition(predicate func(T) bool) (matching, rest *Seq[T]) {
	if s == nil || s.err != nil {
		return &Seq[T]{err: errOf(s)}, &Seq[T]{err: errOf(s)}
	}
	yes := make([]T, 0, len(s.elements)/2)
	no := make([]T, 0, len(s.elements)/2)
	for _, v := range s.elements {
		if predicate(v) {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	return From(yes), From(no)
}
