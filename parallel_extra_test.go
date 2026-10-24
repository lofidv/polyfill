package polyfill_test

import (
	"sync/atomic"
	"testing"

	"github.com/lofidv/polyfill/v3"
	"github.com/stretchr/testify/assert"
)

// Regression: Parallel defaulted to one goroutine per element.
func TestParallelRespectsWorkerLimit(t *testing.T) {
	var inFlight, peak int32

	items := make([]int, 200)
	for i := range items {
		items[i] = i
	}

	got := polyfill.From(items).
		Parallel(4).
		Map(func(n int) int {
			now := atomic.AddInt32(&inFlight, 1)
			for {
				old := atomic.LoadInt32(&peak)
				if now <= old || atomic.CompareAndSwapInt32(&peak, old, now) {
					break
				}
			}
			atomic.AddInt32(&inFlight, -1)
			return n * 2
		}).
		Slice()

	assert.Len(t, got, 200)
	assert.Equal(t, 0, got[0])
	assert.Equal(t, 398, got[199])
	assert.LessOrEqual(t, atomic.LoadInt32(&peak), int32(4))
}

func TestParallelPreservesOrder(t *testing.T) {
	items := make([]int, 100)
	for i := range items {
		items[i] = i
	}

	got := polyfill.From(items).Parallel(8).Map(func(n int) int { return n }).Slice()

	for i := range got {
		assert.Equal(t, i, got[i])
	}
}

// Regression: a panic inside a worker used to kill the process from an
// arbitrary goroutine.
func TestParallelPropagatesPanic(t *testing.T) {
	assert.PanicsWithValue(t, "boom", func() {
		polyfill.From([]int{1, 2, 3}).Parallel(2).Map(func(n int) int {
			if n == 2 {
				panic("boom")
			}
			return n
		})
	})
}

func TestParallelMapTo(t *testing.T) {
	got := polyfill.ParallelMapTo(
		polyfill.From([]int{1, 2, 3}).Parallel(3),
		func(n int) string { return string(rune('a' + n)) },
	).Slice()

	assert.Equal(t, []string{"b", "c", "d"}, got)
}

func TestParallelPropagatesChainError(t *testing.T) {
	seq := failingChain().Parallel(4).Map(func(n int) int { return n })

	assert.ErrorIs(t, seq.Err(), errBadRow)
}
