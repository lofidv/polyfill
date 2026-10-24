package polyfill_test

import (
	"sort"
	"testing"

	"github.com/lofidv/polyfill/v3"
	"github.com/stretchr/testify/assert"
)

type sale struct {
	ID     int
	Region string
	Email  string
	Amount float64
	Units  int
}

func sales() []sale {
	return []sale{
		{1, "eu", "a@x.com", 100, 2},
		{2, "us", "b@x.com", 250, 5},
		{3, "eu", "a@x.com", 50, 1},
		{4, "apac", "c@x.com", 75, 3},
		{5, "us", "d@x.com", 125, 2},
	}
}

func TestGroupByTyped(t *testing.T) {
	t.Run("keys are typed, no assertions needed", func(t *testing.T) {
		byRegion := polyfill.GroupBy(polyfill.From(sales()), func(s sale) string { return s.Region })

		assert.Len(t, byRegion["eu"], 2)
		assert.Len(t, byRegion["us"], 2)
		assert.Len(t, byRegion["apac"], 1)
		// The map is map[string][]sale, so this compiles without a type assertion.
		assert.Equal(t, 100.0, byRegion["eu"][0].Amount)
	})

	t.Run("empty input", func(t *testing.T) {
		assert.Empty(t, polyfill.GroupBy(polyfill.From([]sale{}), func(s sale) string { return s.Region }))
	})
}

func TestCountBy(t *testing.T) {
	counts := polyfill.CountBy(polyfill.From(sales()), func(s sale) string { return s.Region })

	assert.Equal(t, 2, counts["eu"])
	assert.Equal(t, 1, counts["apac"])
}

func TestKeyBy(t *testing.T) {
	byID := polyfill.KeyBy(polyfill.From(sales()), func(s sale) int { return s.ID })

	assert.Len(t, byID, 5)
	assert.Equal(t, "us", byID[2].Region)
}

func TestToMap(t *testing.T) {
	regionByEmail := polyfill.ToMap(
		polyfill.From(sales()),
		func(s sale) string { return s.Email },
		func(s sale) string { return s.Region },
	)

	assert.Equal(t, "eu", regionByEmail["a@x.com"])
	assert.Equal(t, "us", regionByEmail["d@x.com"])
}

func TestPartitionChain(t *testing.T) {
	big, small := polyfill.From(sales()).Partition(func(s sale) bool { return s.Amount >= 100 })

	assert.Len(t, big.Slice(), 3)
	assert.Len(t, small.Slice(), 2)

	// Both sides can chain.
	total := polyfill.SumBy(big.Map(func(s sale) sale { return s }), func(s sale) float64 { return s.Amount })
	assert.Equal(t, 475.0, total)
}

func TestSumByAndAvgBy(t *testing.T) {
	all := polyfill.From(sales())

	t.Run("float sum", func(t *testing.T) {
		assert.Equal(t, 600.0, polyfill.SumBy(all, func(s sale) float64 { return s.Amount }))
	})

	t.Run("int sum", func(t *testing.T) {
		assert.Equal(t, 13, polyfill.SumBy(all, func(s sale) int { return s.Units }))
	})

	t.Run("average", func(t *testing.T) {
		assert.Equal(t, 120.0, polyfill.AvgBy(all, func(s sale) float64 { return s.Amount }))
	})

	t.Run("average of empty is zero", func(t *testing.T) {
		assert.Equal(t, 0.0, polyfill.AvgBy(polyfill.From([]sale{}), func(s sale) float64 { return s.Amount }))
	})
}

func TestSortBy(t *testing.T) {
	t.Run("sorts by string key", func(t *testing.T) {
		sorted := polyfill.SortBy(polyfill.From(sales()), func(s sale) string { return s.Email }).Slice()

		emails := make([]string, len(sorted))
		for i, s := range sorted {
			emails[i] = s.Email
		}
		assert.True(t, sort.StringsAreSorted(emails))
	})

	t.Run("sorts by numeric key", func(t *testing.T) {
		sorted := polyfill.SortBy(polyfill.From(sales()), func(s sale) float64 { return s.Amount }).Slice()

		assert.Equal(t, 50.0, sorted[0].Amount)
		assert.Equal(t, 250.0, sorted[len(sorted)-1].Amount)
	})

	t.Run("does not modify the input", func(t *testing.T) {
		input := sales()
		polyfill.SortBy(polyfill.From(input), func(s sale) float64 { return s.Amount })

		assert.Equal(t, 1, input[0].ID)
	})
}
