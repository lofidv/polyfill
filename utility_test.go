package polyfill_test

import (
	"testing"

	"github.com/lofidv/polyfill"

	"github.com/stretchr/testify/assert"
)

func TestUtilityMethods(t *testing.T) {
	t.Run("IsEmpty", func(t *testing.T) {
		assert.True(t, polyfill.From([]int{}).IsEmpty())
		assert.False(t, polyfill.From([]int{1, 2, 3}).IsEmpty())
	})

	t.Run("Get", func(t *testing.T) {
		s := polyfill.From([]int{10, 20, 30})

		val, ok := s.Get(1)
		assert.True(t, ok)
		assert.Equal(t, 20, val)

		val, ok = s.Get(10)
		assert.False(t, ok)
		assert.Equal(t, 0, val)

		val, ok = s.Get(-1)
		assert.False(t, ok)
	})

	t.Run("First", func(t *testing.T) {
		s := polyfill.From([]string{"a", "b", "c"})

		val, ok := s.First()
		assert.True(t, ok)
		assert.Equal(t, "a", val)

		empty := polyfill.From([]string{})
		val, ok = empty.First()
		assert.False(t, ok)
		assert.Equal(t, "", val)
	})

	t.Run("Last", func(t *testing.T) {
		s := polyfill.From([]string{"a", "b", "c"})

		val, ok := s.Last()
		assert.True(t, ok)
		assert.Equal(t, "c", val)

		empty := polyfill.From([]string{})
		val, ok = empty.Last()
		assert.False(t, ok)
	})

	t.Run("Take", func(t *testing.T) {
		s := polyfill.From([]int{1, 2, 3, 4, 5})

		result := s.Take(3).Slice()
		assert.Equal(t, []int{1, 2, 3}, result)

		result = s.Take(0).Slice()
		assert.Equal(t, []int{}, result)

		result = s.Take(10).Slice()
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("Skip", func(t *testing.T) {
		s := polyfill.From([]int{1, 2, 3, 4, 5})

		result := s.Skip(2).Slice()
		assert.Equal(t, []int{3, 4, 5}, result)

		result = s.Skip(0).Slice()
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)

		result = s.Skip(10).Slice()
		assert.Equal(t, []int{}, result)
	})

	t.Run("ForEach", func(t *testing.T) {
		s := polyfill.From([]int{1, 2, 3})
		sum := 0
		s.ForEach(func(n int) {
			sum += n
		})
		assert.Equal(t, 6, sum)
	})

	t.Run("ForEachIndexed", func(t *testing.T) {
		s := polyfill.From([]string{"a", "b", "c"})
		result := make(map[int]string)
		s.ForEachIndexed(func(i int, val string) {
			result[i] = val
		})
		assert.Equal(t, map[int]string{0: "a", 1: "b", 2: "c"}, result)
	})
}
