package polyfill_test

import (
	"testing"

	"github.com/lofidv/polyfill"

	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	t.Run("sum numbers", func(t *testing.T) {
		nums := []int{1, 2, 3, 4}
		sum := polyfill.Reduce(polyfill.From(nums), 0, func(acc int, n int) int {
			return acc + n
		})

		assert.Equal(t, 10, sum)
	})

	t.Run("concatenate strings", func(t *testing.T) {
		words := []string{"hello", " ", "world"}
		concat := polyfill.Reduce(polyfill.From(words), "", func(acc string, s string) string {
			return acc + s
		})

		assert.Equal(t, "hello world", concat)
	})

	t.Run("empty slice returns initial", func(t *testing.T) {
		empty := []float64{}
		result := polyfill.Reduce(polyfill.From(empty), 100.0, func(acc float64, f float64) float64 {
			return acc + f
		})

		assert.Equal(t, 100.0, result)
	})
}

func TestReduceRight(t *testing.T) {
	t.Run("reverse concatenation", func(t *testing.T) {
		words := []string{"a", "b", "c"}
		result := polyfill.ReduceRight(polyfill.From(words), "", func(acc string, s string) string {
			return acc + s
		})

		assert.Equal(t, "cba", result)
	})
}
