package polyfill_test

import (
	"fmt"

	"github.com/lofidv/polyfill/v3"
)

// Filter keeps the element type, so it is a method and it chains.
func ExampleSeq_Filter() {
	type user struct {
		Name   string
		Active bool
	}
	users := []user{{"ana", true}, {"bob", false}, {"carla", true}}

	active := polyfill.From(users).
		Filter(func(u user) bool { return u.Active }).
		Slice()

	fmt.Println(len(active), active[0].Name)
	// Output: 2 ana
}

// MapTo changes the element type, so it is a free function: Go 1.18 to 1.26 do
// not allow a method to declare a type parameter of its own.
func ExampleMapTo() {
	prices := []int{1999, 500} // cents

	formatted := polyfill.MapTo(polyfill.From(prices), func(cents int) string {
		return fmt.Sprintf("%d.%02d", cents/100, cents%100)
	}).Slice()

	fmt.Println(formatted)
	// Output: [19.99 5.00]
}

// MapE stops at the first error and keeps the rows accepted before it. This is
// the whole reason the error channel exists.
func ExampleSeq_MapE() {
	raw := []string{"1,ana", "2,bob", "broken", "4,dave"}

	names, err := polyfill.From(raw).
		MapE(func(row string) (string, error) {
			parts := splitOnce(row, ',')
			if parts == nil {
				return "", fmt.Errorf("malformed row %q", row)
			}
			return parts[1], nil
		}).
		SliceE()

	fmt.Println(names)
	fmt.Println(err)
	// Output:
	// [ana bob]
	// malformed row "broken"
}

// GroupBy has typed keys, so the caller does not need type assertions.
func ExampleGroupBy() {
	type order struct {
		Region string
		Amount float64
	}
	orders := []order{{"eu", 100}, {"us", 50}, {"eu", 25}}

	byRegion := polyfill.GroupBy(polyfill.From(orders), func(o order) string { return o.Region })

	fmt.Println(len(byRegion["eu"]), byRegion["eu"][1].Amount)
	// Output: 2 25
}

// SortBy sorts by a field without a comparison function, and it does not modify
// the input.
func ExampleSortBy() {
	type ticket struct {
		ID       int
		Priority int
	}
	tickets := []ticket{{1, 3}, {2, 1}, {3, 2}}

	sorted := polyfill.SortBy(polyfill.From(tickets), func(t ticket) int { return t.Priority }).Slice()

	for _, t := range sorted {
		fmt.Print(t.ID, " ")
	}
	fmt.Println()
	// Output: 2 3 1
}

// SumBy and AvgBy replace the accumulator boilerplate in reports.
func ExampleSumBy() {
	type order struct{ Amount float64 }
	orders := []order{{100}, {50}, {25}}

	fmt.Println(polyfill.SumBy(polyfill.From(orders), func(o order) float64 { return o.Amount }))
	fmt.Println(polyfill.AvgBy(polyfill.From(orders), func(o order) float64 { return o.Amount }))
	// Output:
	// 175
	// 58.333333333333336
}

// Partition returns two chains, not two raw slices, so both sides keep working.
func ExampleSeq_Partition() {
	paid, unpaid := polyfill.From([]int{10, 3, 8, 5}).Partition(func(n int) bool { return n >= 5 })

	fmt.Println(paid.Slice(), unpaid.Slice())
	// Output: [10 8 5] [3]
}

// UniqueBy deduplicates by a comparable key in O(n). Unique handles the general
// case and falls back to a deep comparison for non-comparable types.
func ExampleUniqueBy() {
	type user struct {
		ID    int
		Email string
	}
	users := []user{{1, "a@x.com"}, {1, "a@x.com"}, {2, "b@x.com"}}

	unique := polyfill.UniqueBy(polyfill.From(users), func(u user) string { return u.Email }).Slice()

	fmt.Println(len(unique), unique[1].ID)
	// Output: 2 2
}

// FromIter and Iter bridge to the Go 1.23 range-over-func iterators.
func ExampleFromIter() {
	first := func(n int) func(func(int) bool) {
		return func(yield func(int) bool) {
			for i := 0; i < n; i++ {
				if !yield(i) {
					return
				}
			}
		}
	}

	total := polyfill.SumBy(polyfill.FromIter(first(5)), func(n int) int { return n })

	fmt.Println(total)
	// Output: 10
}

func splitOnce(s string, sep byte) []string {
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			return []string{s[:i], s[i+1:]}
		}
	}
	return nil
}
