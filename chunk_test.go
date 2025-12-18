package polyfill_test

import (
	"testing"

	"github.com/lofidv/polyfill"

	"github.com/stretchr/testify/assert"
)

func TestChunk(t *testing.T) {
	t.Run("even chunks", func(t *testing.T) {
		nums := []int{1, 2, 3, 4, 5, 6}
		chunks := polyfill.From(nums).Chunk(2)

		assert.Len(t, chunks, 3)
		assert.Equal(t, []int{1, 2}, chunks[0])
		assert.Equal(t, []int{3, 4}, chunks[1])
		assert.Equal(t, []int{5, 6}, chunks[2])
	})

	t.Run("uneven chunks", func(t *testing.T) {
		words := []string{"a", "b", "c", "d", "e"}
		chunks := polyfill.From(words).Chunk(2)

		assert.Len(t, chunks, 3)
		assert.Equal(t, []string{"a", "b"}, chunks[0])
		assert.Equal(t, []string{"c", "d"}, chunks[1])
		assert.Equal(t, []string{"e"}, chunks[2])
	})

	t.Run("chunk size larger than slice", func(t *testing.T) {
		nums := []float64{1.1, 2.2}
		chunks := polyfill.From(nums).Chunk(5)

		assert.Len(t, chunks, 1)
		assert.Equal(t, nums, chunks[0])
	})
}
