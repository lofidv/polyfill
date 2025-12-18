package polyfill

// Unique returns a sequence with duplicates removed (for comparable types)
// Automatically handles int, string, bool, etc.
//
// Example:
//
//	unique := polyfill.From([]int{1, 2, 2, 3, 4, 4, 5}).Unique().Slice()
//	// Result: [1, 2, 3, 4, 5]
func Unique[T comparable](s *Seq[T]) *Seq[T] {
	if s.err != nil {
		return s
	}

	seen := make(map[T]struct{})
	result := make([]T, 0, len(s.elements))

	for _, v := range s.elements {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return From(result)
}

// UniqueBy returns a sequence with duplicates removed based on a key function
// Use this for structs or when you need custom uniqueness logic
//
// Example:
//
//	uniquePeople := polyfill.From(people).
//	    UniqueBy(func(p Person) string { return p.Name }).
//	    Slice()
func UniqueBy[T any, K comparable](s *Seq[T], keyFn func(T) K) *Seq[T] {
	if s.err != nil {
		return s
	}

	seen := make(map[K]struct{})
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
