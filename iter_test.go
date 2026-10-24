package polyfill_test

import (
	"maps"
	"slices"
	"testing"

	"github.com/lofidv/polyfill/v2"
	"github.com/stretchr/testify/assert"
)

func TestIter(t *testing.T) {
	t.Run("ranges over the sequence", func(t *testing.T) {
		var got []int
		for n := range polyfill.From([]int{1, 2, 3}).Iter() {
			got = append(got, n)
		}

		assert.Equal(t, []int{1, 2, 3}, got)
	})

	t.Run("interoperates with the slices package", func(t *testing.T) {
		got := slices.Collect(polyfill.From([]int{3, 1, 2}).Iter())

		assert.Equal(t, []int{3, 1, 2}, got)
	})

	t.Run("stops when the consumer breaks", func(t *testing.T) {
		var got []int
		for n := range polyfill.From([]int{1, 2, 3}).Iter() {
			if n == 2 {
				break
			}
			got = append(got, n)
		}

		assert.Equal(t, []int{1}, got)
	})

	t.Run("failed chains yield nothing", func(t *testing.T) {
		var count int
		for range failingChain().Iter() {
			count++
		}

		assert.Equal(t, 0, count)
	})
}

func TestFromIter(t *testing.T) {
	t.Run("from slices.Values", func(t *testing.T) {
		seq := polyfill.FromIter(slices.Values([]int{1, 2, 3}))

		assert.Equal(t, []int{1, 2, 3}, seq.Slice())
	})

	t.Run("from a map", func(t *testing.T) {
		index := map[string]int{"a": 1, "b": 2}

		seq := polyfill.FromIter(maps.Keys(index)).Sort(func(a, b string) bool { return a < b })

		assert.Equal(t, []string{"a", "b"}, seq.Slice())
	})

	t.Run("from a generator function", func(t *testing.T) {
		first := func(n int) func(func(int) bool) {
			return func(yield func(int) bool) {
				for i := 0; i < n; i++ {
					if !yield(i) {
						return
					}
				}
			}
		}

		assert.Equal(t, []int{0, 1, 2}, polyfill.FromIter(first(3)).Slice())
	})
}
