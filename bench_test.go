package polyfill_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/lofidv/polyfill/v3"
)

var (
	sinkInts  []int
	sinkFloat float64
)

func benchInts(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = rand.Intn(1000)
	}
	return out
}

// Each helper is compared with the equivalent for loop so the README can state
// a real cost instead of a vague claim.

func BenchmarkFilter(b *testing.B) {
	data := benchInts(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInts = polyfill.From(data).Filter(func(n int) bool { return n%2 == 0 }).Slice()
	}
}

func BenchmarkFilterLoop(b *testing.B) {
	data := benchInts(1000)
	out := make([]int, 0, len(data))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out = out[:0]
		for _, n := range data {
			if n%2 == 0 {
				out = append(out, n)
			}
		}
		sinkInts = out
	}
}

func BenchmarkMap(b *testing.B) {
	data := benchInts(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInts = polyfill.From(data).Map(func(n int) int { return n * 2 }).Slice()
	}
}

func BenchmarkMapLoop(b *testing.B) {
	data := benchInts(1000)
	out := make([]int, len(data))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j, n := range data {
			out[j] = n * 2
		}
		sinkInts = out
	}
}

func BenchmarkSortBy(b *testing.B) {
	data := benchInts(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInts = polyfill.SortBy(polyfill.From(data), func(n int) int { return n }).Slice()
	}
}

func BenchmarkCountBy(b *testing.B) {
	data := benchInts(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = polyfill.CountBy(polyfill.From(data), func(n int) int { return n % 10 })
	}
}

func BenchmarkSumBy(b *testing.B) {
	data := benchInts(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkFloat = polyfill.SumBy(polyfill.From(data), func(n int) float64 { return float64(n) })
	}
}

func BenchmarkSumByLoop(b *testing.B) {
	data := benchInts(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var total float64
		for _, n := range data {
			total += float64(n)
		}
		sinkFloat = total
	}
}

// BenchmarkChain is the real-world shape: four steps over 10k rows.
func BenchmarkChain(b *testing.B) {
	data := benchInts(10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInts = polyfill.From(data).
			Filter(func(n int) bool { return n%2 == 0 }).
			Map(func(n int) int { return n * 2 }).
			Sort(func(a, b int) bool { return a < b }).
			Take(100).
			Slice()
	}
}

const (
	expensiveItems = 50
	expensiveCost  = 100 * time.Microsecond
)

// ParallelMap wins only when the callback is expensive. These two benchmarks
// show both sides.

func BenchmarkExpensiveSequential(b *testing.B) {
	data := make([]int, expensiveItems)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out := make([]int, len(data))
		for j, n := range data {
			time.Sleep(expensiveCost)
			out[j] = n * 2
		}
		sinkInts = out
	}
}

func BenchmarkExpensiveParallel(b *testing.B) {
	data := make([]int, expensiveItems)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInts = polyfill.From(data).
			Parallel(8).
			Map(func(n int) int {
				time.Sleep(expensiveCost)
				return n * 2
			}).
			Slice()
	}
}

func BenchmarkCheapParallel(b *testing.B) {
	data := benchInts(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInts = polyfill.From(data).Parallel().Map(func(n int) int { return n * n }).Slice()
	}
}
