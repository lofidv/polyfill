package polyfill_test

import (
	"testing"

	"github.com/lofidv/polyfill"
	"github.com/stretchr/testify/assert"
)

// Map is now a method, not a global function
// Tests need to use seq.Map() instead of Map(seq, fn)

func TestSeqMap(t *testing.T) {
	t.Run("basic mapping", func(t *testing.T) {
		nums := []int{1, 2, 3}
		result := polyfill.From(nums).
			Map(func(n int) int { return n * 2 }).
			Slice()

		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		empty := []string{}
		result := polyfill.MapTo(polyfill.From(empty), func(s string) int { return len(s) }).
			Slice()

		assert.Empty(t, result)
	})
}

// ParallelMap is now via Parallel() method
func TestParallelSeqMap(t *testing.T) {
	t.Run("parallel map maintains order", func(t *testing.T) {
		nums := []int{1, 2, 3, 4, 5}
		result := polyfill.From(nums).
			Parallel().
			Map(func(n int) int { return n * n }).
			Slice()

		assert.Equal(t, []int{1, 4, 9, 16, 25}, result)
	})

	t.Run("with empty slice", func(t *testing.T) {
		empty := []float64{}
		result := polyfill.From(empty).
			Parallel().
			Map(func(f float64) float64 { return f * 2 }).
			Slice()

		assert.Empty(t, result)
	})
}
