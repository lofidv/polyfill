package polyfill_test

import (
	"github.com/lofidv/polyfill"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	t.Run("sum numbers", func(t *testing.T) {
		nums := []int{1, 2, 3, 4}
		sum := polyfill.Wrap(nums).
			Reduce(0, func(acc any, n int) any {
				return acc.(int) + n
			})

		assert.Equal(t, 10, sum)
	})

	t.Run("concatenate strings", func(t *testing.T) {
		words := []string{"hello", " ", "world"}
		concat := polyfill.Wrap(words).
			Reduce("", func(acc any, s string) any {
				return acc.(string) + s
			})

		assert.Equal(t, "hello world", concat)
	})

	t.Run("empty slice returns initial", func(t *testing.T) {
		empty := []float64{}
		result := polyfill.Wrap(empty).
			Reduce(100.0, func(acc any, f float64) any {
				return acc.(float64) + f
			})

		assert.Equal(t, 100.0, result)
	})
}
