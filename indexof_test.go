package polyfill_test

import (
	"github.com/lofidv/polyfill"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindIndex(t *testing.T) {
	t.Run("find existing element", func(t *testing.T) {
		nums := []int{10, 20, 30, 40}
		index := polyfill.From(nums).
			FindIndex(func(n int) bool { return n == 30 })

		assert.Equal(t, 2, index)
	})

	t.Run("find non-existing element", func(t *testing.T) {
		words := []string{"apple", "banana", "cherry"}
		index := polyfill.From(words).
			FindIndex(func(s string) bool { return s == "date" })

		assert.Equal(t, -1, index)
	})

	t.Run("empty slice returns -1", func(t *testing.T) {
		empty := []bool{}
		index := polyfill.From(empty).
			FindIndex(func(b bool) bool { return b })

		assert.Equal(t, -1, index)
	})
}
