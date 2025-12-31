package main

import (
	"fmt"
	"strings"

	"github.com/lofidv/polyfill/v2"
)

// Example data structures
type Person struct {
	Name string
	Age  int
	City string
}

type Product struct {
	Name  string
	Price float64
	Stock int
}

func main() {
	printHeader("🚀 Polyfill - Functional Programming for Go")

	// Example 1: Basic transformations
	example1()

	// Example 2: Advanced filtering and mapping
	example2()

	// Example 3: Aggregations
	example3()

	// Example 4: Working with structs
	example4()

	// Example 5: Go 1.23 new features
	example5()

	fmt.Println("\n✨ Explore more at github.com/lofidv/polyfill")
}

func example1() {
	printSection("1. Basic Transformations")

	numbers := []int{1, 2, 3, 4, 5}

	// Map: Transform each element
	doubled := polyfill.From(numbers).
		Map(func(n int) int { return n * 2 }).
		Slice()
	fmt.Printf("   Doubled: %v\n", doubled)

	// Filter: Keep only matching elements
	evens := polyfill.From(numbers).
		Filter(func(n int) bool { return n%2 == 0 }).
		Slice()
	fmt.Printf("   Evens: %v\n", evens)

	// Reverse: Flip the order
	reversed := polyfill.From(numbers).Reverse().Slice()
	fmt.Printf("   Reversed: %v\n", reversed)

	fmt.Println()
}

func example2() {
	printSection("2. Advanced Filtering & Mapping")

	words := []string{"hello", "world", "go", "programming", "is", "awesome"}

	// Chain operations
	longWords := polyfill.From(words).
		Filter(func(w string) bool { return len(w) > 3 }).
		Map(func(w string) string { return strings.ToUpper(w) }).
		Sort(func(a, b string) bool { return a < b }).
		Slice()

	fmt.Printf("   Long words (uppercase, sorted): %v\n", longWords)

	// Unique values
	duplicates := []int{1, 2, 2, 3, 3, 3, 4, 5, 5}
	unique := polyfill.From(duplicates).Unique().Slice()
	fmt.Printf("   Unique: %v\n", unique)

	fmt.Println()
}

func example3() {
	printSection("3. Aggregations & Reductions")

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Sum using Reduce
	sum := polyfill.From(numbers).
		Reduce(0, func(acc, n int) int { return acc + n })
	fmt.Printf("   Sum: %d\n", sum)

	// Min and Max
	min, _ := polyfill.From(numbers).MinBy(func(a, b int) bool { return a < b })
	max, _ := polyfill.From(numbers).MaxBy(func(a, b int) bool { return a < b })
	fmt.Printf("   Min: %d, Max: %d\n", min, max)

	// Find first element matching condition
	firstGt5, found := polyfill.From(numbers).Find(func(n int) bool { return n > 5 })
	if found {
		fmt.Printf("   First > 5: %d\n", firstGt5)
	}

	// Check conditions
	hasEven := polyfill.From(numbers).Some(func(n int) bool { return n%2 == 0 })
	allPositive := polyfill.From(numbers).Every(func(n int) bool { return n > 0 })
	fmt.Printf("   Has even: %v, All positive: %v\n", hasEven, allPositive)

	fmt.Println()
}

func example4() {
	printSection("4. Working with Structs")

	people := []Person{
		{"Alice", 28, "NYC"},
		{"Bob", 35, "SF"},
		{"Charlie", 22, "NYC"},
		{"Diana", 31, "LA"},
		{"Eve", 25, "NYC"},
	}

	// Filter and sort
	nycAdults := polyfill.From(people).
		Filter(func(p Person) bool { return p.City == "NYC" && p.Age >= 25 }).
		Sort(func(a, b Person) bool { return a.Age < b.Age }).
		Slice()

	fmt.Println("   NYC adults (sorted by age):")
	for _, p := range nycAdults {
		fmt.Printf("     - %s (%d)\n", p.Name, p.Age)
	}

	// Extract names using MapTo
	names := polyfill.MapTo(
		polyfill.From(people),
		func(p Person) string { return p.Name },
	).Slice()
	fmt.Printf("   All names: %v\n", names)

	// Group by city
	byCity := polyfill.From(people).GroupBy(func(p Person) any { return p.City })
	fmt.Printf("   Grouped by city: %d cities\n", len(byCity))

	// Average age
	totalAge := polyfill.ReduceTo(
		polyfill.From(people),
		0,
		func(acc int, p Person) int { return acc + p.Age },
	)
	avgAge := float64(totalAge) / float64(len(people))
	fmt.Printf("   Average age: %.1f\n", avgAge)

	fmt.Println()
}

func example5() {
	printSection("5. Go 1.23 Features")

	numbers := []int{5, 2, 8, 1, 9, 3}

	// Concat - combine slices
	combined := polyfill.From(numbers).
		Concat([]int{10, 11, 12}).
		Slice()
	fmt.Printf("   Concatenated: %v\n", combined)

	// Prepend - add to beginning
	withZero := polyfill.From(numbers).Prepend(0).Slice()
	fmt.Printf("   With zero prepended: %v\n", withZero)

	// Clone - efficient copy
	original := polyfill.From(numbers)
	clone := original.Clone()
	fmt.Printf("   Original: %v\n", original.Slice())
	fmt.Printf("   Clone: %v\n", clone.Slice())

	// Chunk - split into groups
	chunks := polyfill.From(numbers).Chunk(3)
	fmt.Printf("   Chunks of 3: %v\n", chunks)

	// ContainsFunc - flexible search
	hasLarge := polyfill.From(numbers).ContainsFunc(func(n int) bool { return n > 7 })
	fmt.Printf("   Has number > 7: %v\n", hasLarge)

	fmt.Println()
}

// Helper functions for pretty printing
func printHeader(title string) {
	border := strings.Repeat("=", len(title)+4)
	fmt.Printf("\n%s\n  %s\n%s\n\n", border, title, border)
}

func printSection(title string) {
	fmt.Printf("📌 %s\n", title)
}
