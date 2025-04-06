package polyfill_test

import (
	"github.com/lofidv/polyfill"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvery(t *testing.T) {
	t.Run("all elements match", func(t *testing.T) {
		nums := []int{2, 4, 6, 8}
		result := polyfill.Wrap(nums).
			Every(func(n int) bool { return n%2 == 0 })

		assert.True(t, result)
	})

	t.Run("not all elements match", func(t *testing.T) {
		words := []string{"apple", "banana", "cherry", "date"}
		result := polyfill.Wrap(words).
			Every(func(s string) bool { return len(s) > 4 })

		assert.False(t, result)
	})
}
