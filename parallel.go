package polyfill

import (
	"runtime"
	"sync"
)

// ParallelSeq runs one operation across a bounded pool of goroutines.
// Create it with Seq.Parallel.
type ParallelSeq[T any] struct {
	seq     *Seq[T]
	workers int
}

// Parallel starts a parallel context. The optional argument sets the number of
// workers; it defaults to GOMAXPROCS and is capped at the number of elements.
//
// Use it only when the callback is expensive. For cheap callbacks the
// goroutine and channel overhead costs far more than the work, and a plain Map
// is faster. See the benchmarks in the README.
//
//	From(files).Parallel(8).Map(readAndHash).Slice()
func (s *Seq[T]) Parallel(workers ...int) *ParallelSeq[T] {
	n := runtime.GOMAXPROCS(0)
	if len(workers) > 0 && workers[0] > 0 {
		n = workers[0]
	}
	return &ParallelSeq[T]{seq: s, workers: n}
}

// Map transforms every element concurrently, keeping the order.
// The callback must be safe for concurrent use.
//
// A panic inside the callback is recovered in the worker and re-raised in the
// calling goroutine, so it cannot silently kill an unrelated goroutine.
func (p *ParallelSeq[T]) Map(f func(T) T) *Seq[T] {
	if p == nil || p.seq == nil {
		return nil
	}
	if p.seq.err != nil {
		return p.seq
	}
	return From(runParallel(p.seq.elements, p.workers, f))
}

// ParallelMapTo transforms every element concurrently into another type.
//
//	urls := ParallelMapTo(From(ids).Parallel(16), makeURL).Slice()
func ParallelMapTo[T any, R any](p *ParallelSeq[T], f func(T) R) *Seq[R] {
	if p == nil || p.seq == nil {
		return nil
	}
	if p.seq.err != nil {
		return &Seq[R]{err: p.seq.err}
	}
	return From(runParallel(p.seq.elements, p.workers, f))
}

func runParallel[T any, R any](items []T, workers int, f func(T) R) []R {
	out := make([]R, len(items))
	if len(items) == 0 {
		return out
	}
	if workers > len(items) {
		workers = len(items)
	}
	if workers < 1 {
		workers = 1
	}

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		panicked bool
		panicVal any
	)

	jobs := make(chan int)
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := range jobs {
				func() {
					finished := false
					defer func() {
						r := recover()
						if finished {
							return
						}
						mu.Lock()
						if !panicked {
							panicked, panicVal = true, r
						}
						mu.Unlock()
					}()
					out[i] = f(items[i])
					finished = true
				}()
			}
		}()
	}

	for i := 0; i < len(items); i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	if panicked {
		panic(panicVal)
	}
	return out
}
