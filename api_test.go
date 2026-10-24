package polyfill_test

import (
	"testing"

	"github.com/lofidv/polyfill/v2"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, polyfill.New(1, 2, 3).Slice())
	assert.True(t, polyfill.New[int]().IsEmpty())
}

func TestFlatMapVariants(t *testing.T) {
	t.Run("FlatMapTo changes type", func(t *testing.T) {
		got := polyfill.FlatMapTo(
			polyfill.From([]string{"ab", "c"}),
			func(s string) []rune { return []rune(s) },
		).Slice()

		assert.Equal(t, []rune{'a', 'b', 'c'}, got)
	})

	t.Run("FlatMapE flattens until the first error", func(t *testing.T) {
		seq := polyfill.From([]int{1, 2, 3, 4}).FlatMapE(func(n int) ([]int, error) {
			if n == 3 {
				return nil, errBadRow
			}
			return []int{n, n}, nil
		})

		got, err := seq.SliceE()
		assert.ErrorIs(t, err, errBadRow)
		assert.Equal(t, []int{1, 1, 2, 2}, got)
	})

	t.Run("FlatMapToE flattens and changes type", func(t *testing.T) {
		seq := polyfill.FlatMapToE(
			polyfill.From([]int{1, 2}),
			func(n int) ([]string, error) { return []string{"a", "a"}, nil },
		)

		assert.Equal(t, []string{"a", "a", "a", "a"}, seq.Slice())
		assert.NoError(t, seq.Err())
	})
}

func TestReduceVariants(t *testing.T) {
	t.Run("ReduceTo changes the accumulator type", func(t *testing.T) {
		got := polyfill.ReduceTo(polyfill.From([]int{1, 2, 3}), "",
			func(acc string, n int) string { return acc + string(rune('0'+n)) })

		assert.Equal(t, "123", got)
	})

	t.Run("ReduceE stops at the first error", func(t *testing.T) {
		got, err := polyfill.From([]int{1, 2, 3, 4}).ReduceE(0, func(acc, n int) (int, error) {
			if n == 3 {
				return acc, errBadRow
			}
			return acc + n, nil
		})

		assert.ErrorIs(t, err, errBadRow)
		assert.Equal(t, 3, got)
	})

	t.Run("ReduceToE accumulates into another type", func(t *testing.T) {
		got, err := polyfill.ReduceToE(polyfill.From([]int{1, 2}), "",
			func(acc string, n int) (string, error) {
				return acc + string(rune('0'+n)), nil
			})

		assert.NoError(t, err)
		assert.Equal(t, "12", got)
	})

	t.Run("ReduceRightTo folds from the end", func(t *testing.T) {
		got := polyfill.ReduceRightTo(polyfill.From([]int{1, 2, 3}), "",
			func(acc string, n int) string { return acc + string(rune('0'+n)) })

		assert.Equal(t, "321", got)
	})

	t.Run("error variants surface a chain error", func(t *testing.T) {
		seq := failingChain()

		_, err := seq.ReduceE(0, func(acc, n int) (int, error) { return acc + n, nil })
		assert.ErrorIs(t, err, errBadRow)

		_, err = polyfill.ReduceToE(seq, "", func(acc string, n int) (string, error) { return acc, nil })
		assert.ErrorIs(t, err, errBadRow)
	})
}

func TestIterationHelpers(t *testing.T) {
	t.Run("ForEach visits every element", func(t *testing.T) {
		var seen []int
		polyfill.From([]int{1, 2, 3}).ForEach(func(n int) { seen = append(seen, n) })

		assert.Equal(t, []int{1, 2, 3}, seen)
	})

	t.Run("ForEachIndexed passes the index", func(t *testing.T) {
		var idx []int
		polyfill.From([]string{"a", "b"}).ForEachIndexed(func(i int, _ string) { idx = append(idx, i) })

		assert.Equal(t, []int{0, 1}, idx)
	})

	t.Run("both are no-ops on a failed chain", func(t *testing.T) {
		calls := 0
		failingChain().ForEach(func(int) { calls++ })
		failingChain().ForEachIndexed(func(int, int) { calls++ })

		assert.Equal(t, 0, calls)
	})
}

func TestNilReceivers(t *testing.T) {
	var seq *polyfill.Seq[int]

	assert.NotPanics(t, func() {
		_ = seq.Clone()
		_ = seq.Append(1, 2)
		_ = seq.Prepend(0)
		_ = seq.Take(1)
		_ = seq.Skip(1)
		_ = seq.Sort(func(a, b int) bool { return a < b })
		_ = seq.Reverse()
		_ = seq.Unique()
		_ = seq.Chunk(2)
		_ = polyfill.SortBy(seq, func(n int) int { return n })
		_ = polyfill.SumBy(seq, func(n int) int { return n })
		_ = polyfill.GroupBy(seq, func(n int) int { return n })
	})

	assert.Equal(t, []int{1, 2}, seq.Append(1, 2).Slice())
	assert.Equal(t, []int{0}, seq.Prepend(0).Slice())
	assert.True(t, seq.Clone().IsEmpty())
}

func TestConcatVariants(t *testing.T) {
	t.Run("EqualFunc", func(t *testing.T) {
		eq := func(a, b int) bool { return a == b }

		assert.True(t, polyfill.From([]int{1, 2}).EqualFunc([]int{1, 2}, eq))
		assert.False(t, polyfill.From([]int{1, 2}).EqualFunc([]int{1}, eq))
		assert.False(t, polyfill.From([]int{1, 2}).EqualFunc([]int{1, 3}, eq))
	})

	t.Run("Clone is independent", func(t *testing.T) {
		original := polyfill.From([]int{1, 2})
		clone := original.Clone()
		clone.Push(3)

		assert.Equal(t, []int{1, 2}, original.Slice())
		assert.Equal(t, []int{1, 2, 3}, clone.Slice())
	})

	t.Run("Concat with several slices", func(t *testing.T) {
		got := polyfill.From([]int{1}).Concat([]int{2}, []int{3, 4}).Slice()

		assert.Equal(t, []int{1, 2, 3, 4}, got)
	})
}

func TestCountByEmpty(t *testing.T) {
	assert.Empty(t, polyfill.CountBy(polyfill.From([]int{}), func(n int) int { return n }))
	assert.Empty(t, polyfill.CountBy(failingChain(), func(n int) int { return n }))
}

func TestIndexOfNilElement(t *testing.T) {
	values := []any{nil, 1, "x"}

	assert.Equal(t, 0, polyfill.From(values).IndexOf(nil))
	assert.Equal(t, 1, polyfill.From(values).IndexOf(1))
	assert.True(t, polyfill.From(values).Includes("x"))
}
