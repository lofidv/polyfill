package polyfill_test

import (
	"github.com/lofidv/polyfill"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Run("basic map operation", func(t *testing.T) {
		nums := []int{1, 2, 3}
		result := polyfill.Wrap(nums).
			Map(func(n int) any { return n * 2 }).
			Unwrap()

		assert.Equal(t, []any{2, 4, 6}, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		empty := []string{}
		result := polyfill.Wrap(empty).
			Map(func(s string) any { return len(s) }).
			Unwrap()

		assert.Empty(t, result)
	})
}

func TestParallelMap(t *testing.T) {
	t.Run("parallel map maintains order", func(t *testing.T) {
		nums := []int{1, 2, 3, 4, 5}
		result := polyfill.Wrap(nums).
			ParallelMap(func(n int) any { return n * n }).
			Unwrap()

		assert.Equal(t, []any{1, 4, 9, 16, 25}, result)
	})

	t.Run("with empty slice", func(t *testing.T) {
		empty := []float64{}
		result := polyfill.Wrap(empty).
			ParallelMap(func(f float64) any { return f * 2 }).
			Unwrap()

		assert.Empty(t, result)
	})
}
