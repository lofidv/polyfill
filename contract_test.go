package polyfill_test

import (
	"errors"
	"testing"

	"github.com/lofidv/polyfill/v2"
	"github.com/stretchr/testify/assert"
)

var errBadRow = errors.New("bad row")

func failingChain() *polyfill.Seq[int] {
	return polyfill.From([]int{1, 2, 3, 4}).MapE(func(n int) (int, error) {
		if n == 3 {
			return 0, errBadRow
		}
		return n * 10, nil
	})
}

func TestErrorContract(t *testing.T) {
	t.Run("keeps the elements processed before the failure", func(t *testing.T) {
		got, err := failingChain().SliceE()

		assert.ErrorIs(t, err, errBadRow)
		assert.Equal(t, []int{10, 20}, got)
	})

	t.Run("later transformations are no-ops", func(t *testing.T) {
		seq := failingChain().
			Filter(func(n int) bool { return true }).
			Reverse().
			Sort(func(a, b int) bool { return a < b }).
			Skip(0)

		assert.ErrorIs(t, seq.Err(), errBadRow)
		assert.Equal(t, []int{10, 20}, seq.Slice())
	})

	t.Run("Len reports the partial length", func(t *testing.T) {
		assert.Equal(t, 2, failingChain().Len())
	})

	t.Run("terminal accessors return zero values", func(t *testing.T) {
		seq := failingChain()

		_, ok := seq.First()
		assert.False(t, ok)
		_, ok = seq.Last()
		assert.False(t, ok)
		_, ok = seq.At(0)
		assert.False(t, ok)
		_, ok = seq.Find(func(int) bool { return true })
		assert.False(t, ok)
		_, ok = seq.MinBy(func(a, b int) bool { return a < b })
		assert.False(t, ok)
		assert.Equal(t, -1, seq.FindIndex(func(int) bool { return true }))
		assert.Equal(t, -1, seq.IndexOf(10))
		assert.False(t, seq.Includes(10))
		assert.False(t, seq.Some(func(int) bool { return true }))
		assert.False(t, seq.Every(func(int) bool { return true }))
		assert.False(t, seq.ContainsFunc(func(int) bool { return true }))
	})

	t.Run("collections are empty", func(t *testing.T) {
		seq := failingChain()

		assert.Empty(t, polyfill.GroupBy(seq, func(n int) int { return n }))
		assert.Empty(t, polyfill.KeyBy(seq, func(n int) int { return n }))
		assert.Empty(t, polyfill.ToMap(seq, func(n int) int { return n }, func(n int) string { return "" }))
		assert.Nil(t, seq.Chunk(2))
		unique := polyfill.UniqueBy(seq, func(n int) int { return n })
		assert.ErrorIs(t, unique.Err(), errBadRow)
		assert.Equal(t, []int{10, 20}, unique.Slice())
	})

	t.Run("folds return the initial value", func(t *testing.T) {
		seq := failingChain()

		assert.Equal(t, 100, seq.Reduce(100, func(acc, n int) int { return acc + n }))
		assert.Equal(t, "x", polyfill.ReduceTo(seq, "x", func(acc string, n int) string { return acc }))
		assert.Equal(t, 0.0, polyfill.SumBy(seq, func(n int) float64 { return float64(n) }))
	})

	t.Run("error survives a type change", func(t *testing.T) {
		seq := polyfill.MapTo(failingChain(), func(n int) string { return "x" })

		assert.ErrorIs(t, seq.Err(), errBadRow)
		assert.Empty(t, seq.Slice())
	})

	t.Run("MapToE keeps what it produced", func(t *testing.T) {
		seq := polyfill.MapToE(polyfill.From([]int{1, 2, 3}), func(n int) (string, error) {
			if n == 3 {
				return "", errBadRow
			}
			return "ok", nil
		})

		got, err := seq.SliceE()
		assert.ErrorIs(t, err, errBadRow)
		assert.Equal(t, []string{"ok", "ok"}, got)
	})

	t.Run("Partition propagates the error to both sides", func(t *testing.T) {
		yes, no := failingChain().Partition(func(n int) bool { return n > 15 })

		assert.ErrorIs(t, yes.Err(), errBadRow)
		assert.ErrorIs(t, no.Err(), errBadRow)
	})

	t.Run("a clean chain has no error", func(t *testing.T) {
		seq := polyfill.From([]int{1, 2}).Map(func(n int) int { return n * 2 })

		assert.NoError(t, seq.Err())
		_, err := seq.SliceE()
		assert.NoError(t, err)
	})
}

func TestNilSeqIsSafe(t *testing.T) {
	var seq *polyfill.Seq[int]

	assert.Equal(t, 0, seq.Len())
	assert.True(t, seq.IsEmpty())
	assert.Nil(t, seq.Slice())
	assert.NoError(t, seq.Err())

	_, ok := seq.First()
	assert.False(t, ok)

	assert.Empty(t, seq.Filter(func(int) bool { return true }).Slice())
	assert.Equal(t, []int{9}, seq.Push(9).Slice())
}
