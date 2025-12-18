package polyfill_test

import (
	"github.com/lofidv/polyfill"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	t.Run("reverse slice", func(t *testing.T) {
		nums := []int{1, 2, 3, 4}
		reversed := polyfill.From(nums).Reverse().Slice()

		assert.Equal(t, []int{4, 3, 2, 1}, reversed)
	})

	t.Run("reverse empty slice", func(t *testing.T) {
		empty := []string{}
		reversed := polyfill.From(empty).Reverse().Slice()

		assert.Empty(t, reversed)
	})

	t.Run("reverse single element", func(t *testing.T) {
		single := []float64{3.14}
		reversed := polyfill.From(single).Reverse().Slice()

		assert.Equal(t, single, reversed)
	})
}
