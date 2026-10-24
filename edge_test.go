package polyfill_test

import (
	"testing"

	"github.com/lofidv/polyfill/v2"
	"github.com/stretchr/testify/assert"
)

// Regression: Unique and UniqueBy compared values through map[any], so a
// non-comparable element type panicked at runtime.
func TestUniqueNonComparable(t *testing.T) {
	rows := [][]int{{1, 2}, {1, 2}, {3, 4}}

	assert.NotPanics(t, func() {
		got := polyfill.From(rows).Unique().Slice()
		assert.Equal(t, [][]int{{1, 2}, {3, 4}}, got)
	})
}

func TestUniqueByTyped(t *testing.T) {
	got := polyfill.UniqueBy(polyfill.From(sales()), func(s sale) string { return s.Email }).Slice()

	assert.Len(t, got, 4) // a@x.com appears twice
	assert.Equal(t, 1, got[0].ID)
}

// Regression: IndexOf used any(v) == any(value), which panics for slices.
func TestIndexOfNonComparable(t *testing.T) {
	rows := [][]int{{1, 2}, {3, 4}}

	assert.NotPanics(t, func() {
		assert.Equal(t, 1, polyfill.From(rows).IndexOf([]int{3, 4}))
		assert.True(t, polyfill.From(rows).Includes([]int{1, 2}))
		assert.False(t, polyfill.From(rows).Includes([]int{9}))
	})
}

func TestChunkRejectsInvalidSize(t *testing.T) {
	assert.PanicsWithValue(t, "polyfill: Chunk size must be greater than zero", func() {
		polyfill.From([]int{1, 2, 3}).Chunk(0)
	})
	assert.PanicsWithValue(t, "polyfill: Chunk size must be greater than zero", func() {
		polyfill.From([]int{1, 2, 3}).Chunk(-1)
	})
}

func TestChunkEmpty(t *testing.T) {
	assert.Empty(t, polyfill.From([]int{}).Chunk(3))
}

// Regression: Take and Skip returned a view, so Push on the result wrote into
// the caller's backing array.
func TestTakeAndSkipDoNotAlias(t *testing.T) {
	original := []int{1, 2, 3, 4}

	taken := polyfill.From(original).Take(2).Push(99)
	assert.Equal(t, []int{1, 2, 99}, taken.Slice())
	assert.Equal(t, []int{1, 2, 3, 4}, original)

	skipped := polyfill.From(original).Skip(2).Push(99)
	assert.Equal(t, []int{3, 4, 99}, skipped.Slice())
	assert.Equal(t, []int{1, 2, 3, 4}, original)
}

func TestPushMutates(t *testing.T) {
	seq := polyfill.From([]int{1, 2})

	seq.Push(3)

	assert.Equal(t, []int{1, 2, 3}, seq.Slice())
}

func TestAtNegativeIndex(t *testing.T) {
	seq := polyfill.From([]string{"a", "b", "c"})

	last, ok := seq.At(-1)
	assert.True(t, ok)
	assert.Equal(t, "c", last)

	_, ok = seq.At(-4)
	assert.False(t, ok)
}
