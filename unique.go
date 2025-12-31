package polyfill

// Unique returns a sequence with duplicates removed
//
// Example:
//
//	From([]int{1, 2, 2, 3}).Unique().Slice() // [1, 2, 3]
func (s *Seq[T]) Unique() *Seq[T] {
	if s.err != nil {
		return s
	}

	seen := make(map[any]struct{})
	result := make([]T, 0, len(s.elements))

	for _, v := range s.elements {
		key := any(v)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, v)
		}
	}

	return From(result)
}

// UniqueBy returns a sequence with duplicates removed based on a key function
//
// Example:
//
//	From(items).UniqueBy(func(v Item) any { return v.ID }).Slice()
func (s *Seq[T]) UniqueBy(keyFn func(T) any) *Seq[T] {
	if s.err != nil {
		return s
	}

	seen := make(map[any]struct{})
	result := make([]T, 0, len(s.elements))

	for _, v := range s.elements {
		key := keyFn(v)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, v)
		}
	}

	return From(result)
}
