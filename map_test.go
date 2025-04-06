package polyfill_test

import (
	"testing"

	"github.com/lofidv/polyfill"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Run("operación básica de mapeo", func(t *testing.T) {
		nums := []int{1, 2, 3}
		result := polyfill.Map(polyfill.Wrap(nums), func(n int) int { return n * 2 }).Unwrap()

		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("slice vacío", func(t *testing.T) {
		empty := []string{}
		result := polyfill.Map(polyfill.Wrap(empty), func(s string) int { return len(s) }).Unwrap()

		assert.Empty(t, result)
	})
}

func TestParallelMap(t *testing.T) {
	t.Run("parallel map mantiene el orden", func(t *testing.T) {
		nums := []int{1, 2, 3, 4, 5}
		result := polyfill.ParallelMap(polyfill.Wrap(nums), func(n int) int { return n * n }).Unwrap()

		assert.Equal(t, []int{1, 4, 9, 16, 25}, result)
	})

	t.Run("con slice vacío", func(t *testing.T) {
		empty := []float64{}
		result := polyfill.ParallelMap(polyfill.Wrap(empty), func(f float64) float64 { return f * 2 }).Unwrap()

		assert.Empty(t, result)
	})
}
