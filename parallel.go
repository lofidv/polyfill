package polyfill

import "sync"

// === PARALLEL EXECUTION ===

// ParallelOptions configures parallel execution
type ParallelOptions struct {
	Workers int  // number of goroutines (0 = number of elements)
	Ordered bool // preserve order (default true)
}

// ParallelSeq wraps a Seq for parallel operations
type ParallelSeq[T any] struct {
	seq  *Seq[T]
	opts ParallelOptions
}

// Parallel returns a parallel execution context
//
// Example:
//
//	squared := polyfill.From(numbers).
//	    Parallel().
//	    Map(func(n int) int { return n * n }).
//	    Slice()
func (s *Seq[T]) Parallel(opts ...ParallelOptions) *ParallelSeq[T] {
	opt := ParallelOptions{Workers: len(s.elements), Ordered: true}
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt.Workers <= 0 {
		opt.Workers = len(s.elements)
	}
	return &ParallelSeq[T]{seq: s, opts: opt}
}

// Map transforms elements concurrently (same type)
func (p *ParallelSeq[T]) Map(f func(T) T) *Seq[T] {
	result := make([]T, len(p.seq.elements))
	var wg sync.WaitGroup

	workers := p.opts.Workers
	if workers > len(p.seq.elements) {
		workers = len(p.seq.elements)
	}

	jobs := make(chan int, len(p.seq.elements))

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				result[idx] = f(p.seq.elements[idx])
			}
		}()
	}

	for i := range p.seq.elements {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	return From(result)
}

// Slice returns the result as a slice
func (p *ParallelSeq[T]) Slice() []T {
	return p.seq.elements
}

// ParallelMapTo transforms elements concurrently with type change
func ParallelMapTo[T any, R any](p *ParallelSeq[T], f func(T) R) *Seq[R] {
	result := make([]R, len(p.seq.elements))
	var wg sync.WaitGroup

	workers := p.opts.Workers
	if workers > len(p.seq.elements) {
		workers = len(p.seq.elements)
	}

	jobs := make(chan int, len(p.seq.elements))

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				result[idx] = f(p.seq.elements[idx])
			}
		}()
	}

	for i := range p.seq.elements {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	return From(result)
}
