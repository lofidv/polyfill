package polyfill

import "sync"

// Map transforms each element using the provided function
func (s *Slice[T]) Map(f func(T) any) *Slice[any] {
	result := make([]any, len(s.elements))
	for i, v := range s.elements {
		result[i] = f(v)
	}
	return Wrap(result)
}

// ParallelMap transforms elements concurrently (use with caution for CPU-bound operations)
func (s *Slice[T]) ParallelMap(f func(T) any) *Slice[any] {
	result := make([]any, len(s.elements))
	var wg sync.WaitGroup
	wg.Add(len(s.elements))

	for i, v := range s.elements {
		go func(idx int, val T) {
			defer wg.Done()
			result[idx] = f(val)
		}(i, v)
	}

	wg.Wait()
	return Wrap(result)
}
