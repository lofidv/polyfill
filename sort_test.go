package polyfill_test

import (
	"github.com/lofidv/polyfill"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	t.Run("sort numbers ascending", func(t *testing.T) {
		nums := []int{3, 1, 4, 2}
		sorted := polyfill.Wrap(nums).
			Sort(func(a, b int) bool { return a < b }).
			Unwrap()

		assert.Equal(t, []int{1, 2, 3, 4}, sorted)
	})

	t.Run("sort strings by length", func(t *testing.T) {
		words := []string{"apple", "banana", "cherry", "date"}
		sorted := polyfill.Wrap(words).
			Sort(func(a, b string) bool { return len(a) < len(b) }).
			Unwrap()

		assert.Equal(t, []string{"date", "apple", "banana", "cherry"}, sorted)
	})

	t.Run("sort empty slice", func(t *testing.T) {
		empty := []float64{}
		sorted := polyfill.Wrap(empty).
			Sort(func(a, b float64) bool { return a < b }).
			Unwrap()

		assert.Empty(t, sorted)
	})
}
