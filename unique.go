package polyfill

import "reflect"

// Unique returns a new sequence with duplicates removed, keeping the first
// occurrence and the original order.
//
// It works for every element type. Comparable elements use an O(n) map;
// non-comparable elements (a slice field, a map, a struct that contains one)
// fall back to an O(n^2) comparison instead of panicking.
//
// Prefer UniqueBy when a cheap, comparable identity exists: it is always O(n).
func (s *Seq[T]) Unique() *Seq[T] {
	if s == nil || s.err != nil {
		return s
	}
	result := make([]T, 0, len(s.elements))

	if allComparable(s.elements) {
		seen := make(map[any]struct{}, len(s.elements))
		for _, v := range s.elements {
			k := any(v)
			if _, exists := seen[k]; exists {
				continue
			}
			seen[k] = struct{}{}
			result = append(result, v)
		}
		return From(result)
	}

	for _, v := range s.elements {
		duplicate := false
		for _, kept := range result {
			if reflect.DeepEqual(kept, v) {
				duplicate = true
				break
			}
		}
		if !duplicate {
			result = append(result, v)
		}
	}
	return From(result)
}

// UniqueBy removes duplicates using a comparable key derived from each element.
// It keeps the first occurrence, preserves order, and runs in O(n).
//
//	UniqueBy(From(users), func(u User) int { return u.ID }).Slice()
//	UniqueBy(From(orders), func(o Order) string { return o.Email }).Slice()
func UniqueBy[T any, K comparable](s *Seq[T], key func(T) K) *Seq[T] {
	if s == nil || s.err != nil {
		return s
	}
	seen := make(map[K]struct{}, len(s.elements))
	result := make([]T, 0, len(s.elements))
	for _, v := range s.elements {
		k := key(v)
		if _, exists := seen[k]; exists {
			continue
		}
		seen[k] = struct{}{}
		result = append(result, v)
	}
	return From(result)
}

func allComparable[T any](items []T) bool {
	for _, v := range items {
		x := any(v)
		if x == nil {
			continue
		}
		if !reflect.TypeOf(x).Comparable() {
			return false
		}
	}
	return true
}

// sameValue compares two values of any type without panicking on
// non-comparable types.
func sameValue[T any](a, b T) bool {
	av, bv := any(a), any(b)
	if av == nil || bv == nil {
		return av == bv
	}
	if reflect.TypeOf(av).Comparable() && reflect.TypeOf(bv).Comparable() {
		return av == bv
	}
	return reflect.DeepEqual(a, b)
}
