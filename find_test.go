package polyfill_test

import (
	"github.com/lofidv/polyfill"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	people := []person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 35},
	}

	t.Run("find existing element", func(t *testing.T) {
		p, found := polyfill.Wrap(people).
			Find(func(p person) bool { return p.Name == "Bob" })

		assert.True(t, found)
		assert.Equal(t, "Bob", p.Name)
	})

	t.Run("find non-existing element", func(t *testing.T) {
		_, found := polyfill.Wrap(people).
			Find(func(p person) bool { return p.Name == "David" })

		assert.False(t, found)
	})

	t.Run("find in empty slice", func(t *testing.T) {
		empty := []int{}
		_, found := polyfill.Wrap(empty).
			Find(func(n int) bool { return n > 10 })

		assert.False(t, found)
	})
}
