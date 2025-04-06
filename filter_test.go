package polyfill_test

import (
	"github.com/lofidv/polyfill"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	t.Run("filter even numbers", func(t *testing.T) {
		nums := []int{1, 2, 3, 4, 5, 6}
		result := polyfill.Wrap(nums).
			Filter(func(n int) bool { return n%2 == 0 }).
			Unwrap()

		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("filter with no matches", func(t *testing.T) {
		words := []string{"apple", "banana", "cherry"}
		result := polyfill.Wrap(words).
			Filter(func(s string) bool { return len(s) > 10 }).
			Unwrap()

		assert.Empty(t, result)
	})

	t.Run("filter empty slice", func(t *testing.T) {
		empty := []bool{}
		result := polyfill.Wrap(empty).
			Filter(func(b bool) bool { return b }).
			Unwrap()

		assert.Empty(t, result)
	})
}
