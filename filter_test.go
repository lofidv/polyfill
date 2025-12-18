package polyfill_test

import (
	"github.com/lofidv/polyfill"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	t.Run("filter even numbers", func(t *testing.T) {
		nums := []int{1, 2, 3, 4, 5, 6}
		result := polyfill.From(nums).
			Filter(func(n int) bool { return n%2 == 0 }).
			Slice()

		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("filter with no matches", func(t *testing.T) {
		words := []string{"apple", "banana", "cherry"}
		result := polyfill.From(words).
			Filter(func(s string) bool { return len(s) > 10 }).
			Slice()

		assert.Empty(t, result)
	})

	t.Run("filter empty slice", func(t *testing.T) {
		empty := []bool{}
		result := polyfill.From(empty).
			Filter(func(b bool) bool { return b }).
			Slice()

		assert.Empty(t, result)
	})
}
