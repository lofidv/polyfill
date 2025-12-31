package polyfill_test

import (
	"testing"

	"github.com/lofidv/polyfill"
	"github.com/stretchr/testify/assert"
)

func TestConcat(t *testing.T) {
	t.Run("concat multiple slices", func(t *testing.T) {
		result := polyfill.From([]int{1, 2}).
			Concat([]int{3, 4}, []int{5, 6}).
			Slice()
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("concat empty", func(t *testing.T) {
		result := polyfill.From([]int{1, 2}).
			Concat().
			Slice()
		assert.Equal(t, []int{1, 2}, result)
	})
}

func TestAppend(t *testing.T) {
	t.Run("append elements", func(t *testing.T) {
		result := polyfill.From([]int{1, 2}).
			Append(3, 4).
			Slice()
		assert.Equal(t, []int{1, 2, 3, 4}, result)
	})
}

func TestPrepend(t *testing.T) {
	t.Run("prepend elements", func(t *testing.T) {
		result := polyfill.From([]int{3, 4}).
			Prepend(1, 2).
			Slice()
		assert.Equal(t, []int{1, 2, 3, 4}, result)
	})
}

func TestContainsFunc(t *testing.T) {
	t.Run("contains existing", func(t *testing.T) {
		result := polyfill.From([]int{1, 2, 3}).ContainsFunc(func(n int) bool { return n == 2 })
		assert.True(t, result)
	})

	t.Run("contains non-existing", func(t *testing.T) {
		result := polyfill.From([]int{1, 2, 3}).ContainsFunc(func(n int) bool { return n == 5 })
		assert.False(t, result)
	})
}

func TestEqualFunc(t *testing.T) {
	t.Run("equal slices", func(t *testing.T) {
		result := polyfill.From([]int{1, 2, 3}).EqualFunc([]int{1, 2, 3}, func(a, b int) bool { return a == b })
		assert.True(t, result)
	})

	t.Run("not equal slices", func(t *testing.T) {
		result := polyfill.From([]int{1, 2, 3}).EqualFunc([]int{1, 2, 4}, func(a, b int) bool { return a == b })
		assert.False(t, result)
	})
}

func TestClone(t *testing.T) {
	t.Run("clone sequence", func(t *testing.T) {
		original := polyfill.From([]int{1, 2, 3})
		clone := original.Clone()
		assert.Equal(t, original.Slice(), clone.Slice())
	})
}

func TestMinMaxBy(t *testing.T) {
	t.Run("min of integers", func(t *testing.T) {
		min, ok := polyfill.From([]int{3, 1, 4, 1, 5}).MinBy(func(a, b int) bool { return a < b })
		assert.True(t, ok)
		assert.Equal(t, 1, min)
	})

	t.Run("max of integers", func(t *testing.T) {
		max, ok := polyfill.From([]int{3, 1, 4, 1, 5}).MaxBy(func(a, b int) bool { return a < b })
		assert.True(t, ok)
		assert.Equal(t, 5, max)
	})

	t.Run("min of empty slice", func(t *testing.T) {
		_, ok := polyfill.From([]int{}).MinBy(func(a, b int) bool { return a < b })
		assert.False(t, ok)
	})

	t.Run("max of strings", func(t *testing.T) {
		max, ok := polyfill.From([]string{"apple", "zebra", "banana"}).MaxBy(func(a, b string) bool { return a < b })
		assert.True(t, ok)
		assert.Equal(t, "zebra", max)
	})
}
