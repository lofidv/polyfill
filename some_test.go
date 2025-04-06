package polyfill_test

import (
	"github.com/lofidv/polyfill"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSome(t *testing.T) {
	t.Run("some elements match", func(t *testing.T) {
		nums := []int{1, 3, 5, 7, 8}
		result := polyfill.Wrap(nums).
			Some(func(n int) bool { return n%2 == 0 })

		assert.True(t, result)
	})

	t.Run("no elements match", func(t *testing.T) {
		words := []string{"apple", "banana", "cherry"}
		result := polyfill.Wrap(words).
			Some(func(s string) bool { return len(s) > 10 })

		assert.False(t, result)
	})
}
