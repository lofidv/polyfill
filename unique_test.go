package polyfill_test

import (
	"testing"

	"github.com/lofidv/polyfill"

	"github.com/stretchr/testify/assert"
)

func TestUnique(t *testing.T) {
	t.Run("remove duplicates", func(t *testing.T) {
		nums := []int{1, 2, 2, 3, 4, 4, 4, 5}
		unique := polyfill.From(nums).Unique().Slice()

		assert.Equal(t, []int{1, 2, 3, 4, 5}, unique)
	})

	t.Run("all unique", func(t *testing.T) {
		words := []string{"a", "b", "c"}
		unique := polyfill.From(words).Unique().Slice()

		assert.Equal(t, words, unique)
	})

	t.Run("empty slice", func(t *testing.T) {
		empty := []bool{}
		unique := polyfill.From(empty).Unique().Slice()

		assert.Empty(t, unique)
	})
}
