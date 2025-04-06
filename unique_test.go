package polyfill_test

import (
	"github.com/lofidv/polyfill"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnique(t *testing.T) {
	t.Run("remove duplicates", func(t *testing.T) {
		nums := []int{1, 2, 2, 3, 4, 4, 4, 5}
		unique := polyfill.Wrap(nums).
			Unique(func(a, b int) bool { return a == b }).
			Unwrap()

		assert.Equal(t, []int{1, 2, 3, 4, 5}, unique)
	})

	t.Run("all unique", func(t *testing.T) {
		words := []string{"a", "b", "c"}
		unique := polyfill.Wrap(words).
			Unique(func(a, b string) bool { return a == b }).
			Unwrap()

		assert.Equal(t, words, unique)
	})

	t.Run("empty slice", func(t *testing.T) {
		empty := []bool{}
		unique := polyfill.Wrap(empty).
			Unique(func(a, b bool) bool { return a == b }).
			Unwrap()

		assert.Empty(t, unique)
	})
}
