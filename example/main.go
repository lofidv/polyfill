package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/lofidv/polyfill"
)

type Person struct {
	Name string
	Age  int
}

func (p Person) IsAdult() bool {
	return p.Age >= 18
}

type Pet struct {
	Name string
	Age  int
	Type string
}

func main() {
	// 1. Basic Filter and Map operations
	people := []Person{
		{"Alice", 25},
		{"Bob", 17},
		{"Charlie", 19},
		{"David", 16},
	}

	// Filter adults and map to names - demonstrating chainable API
	adultNames := polyfill.MapTo(
		polyfill.From(people).Filter(func(p Person) bool { return p.IsAdult() }),
		func(p Person) string { return p.Name },
	).Slice()

	fmt.Println("Adult names:", adultNames) // ["Alice", "Charlie"]

	// 2. Parallel Map demonstration
	numbers := polyfill.From([]int{1, 2, 3, 4, 5})
	squared := numbers.Parallel().Map(func(n int) int { return n * n }).Slice()

	fmt.Println("Squared numbers:", squared) // [1, 4, 9, 16, 25]

	// 3. Reduce example with type safety
	sum := polyfill.Reduce(numbers, 0, func(acc int, n int) int { return acc + n })

	fmt.Println("Sum of numbers:", sum) // 15

	// 4. Find and FindIndex
	bob, found := polyfill.From(people).
		Find(func(p Person) bool { return p.Name == "Bob" })

	fmt.Printf("Found Bob: %v (%v)\n", bob, found) // {Bob 17} true

	bobIndex := polyfill.From(people).
		FindIndex(func(p Person) bool { return p.Name == "Bob" })

	fmt.Println("Bob's index:", bobIndex) // 1

	// 5. Some and Every
	hasAdults := polyfill.From(people).
		Some(func(p Person) bool { return p.IsAdult() })

	allAdults := polyfill.From(people).
		Every(func(p Person) bool { return p.IsAdult() })

	fmt.Printf("Has adults: %v, All adults: %v\n", hasAdults, allAdults) // true, false

	// 6. Chunk and Reverse
	chunks := numbers.Chunk(2)
	fmt.Println("Chunks:")
	for _, chunk := range chunks {
		fmt.Println(chunk)
	}
	// [1 2]
	// [3 4]
	// [5]

	reversed := numbers.Reverse().Slice()
	fmt.Println("Reversed numbers:", reversed) // [5, 4, 3, 2, 1]

	// 7. Sort
	unsorted := []int{3, 1, 4, 2}
	sorted := polyfill.From(unsorted).
		Sort(func(a, b int) bool { return a < b }).
		Slice()

	fmt.Println("Sorted numbers:", sorted) // [1, 2, 3, 4]

	// 8. Unique (for comparable types)
	duplicates := []int{1, 2, 2, 3, 4, 4, 5}
	unique := polyfill.Unique(polyfill.From(duplicates)).Slice()

	fmt.Println("Unique numbers:", unique) // [1, 2, 3, 4, 5]

	// 9. String conversion with type safety
	strNumbers := polyfill.From([]string{"1", "2", "3", "4"})
	intNumbers := polyfill.MapTo(strNumbers, func(s string) int {
		n, _ := strconv.Atoi(s)
		return n
	}).Slice()

	fmt.Println("Converted numbers:", intNumbers) // [1, 2, 3, 4]

	// 10. GroupBy - killer feature!
	pets := []Pet{
		{"Fido", 3, "dog"},
		{"Whiskers", 2, "cat"},
		{"Rover", 5, "dog"},
		{"Mittens", 1, "cat"},
	}

	grouped := polyfill.GroupBy(polyfill.From(pets), func(p Pet) string {
		return p.Type
	})

	fmt.Println("Pets grouped by type:")
	for typ, pets := range grouped {
		fmt.Printf("%s: %v\n", typ, pets)
	}
	// dog: [{Fido 3 dog} {Rover 5 dog}]
	// cat: [{Whiskers 2 cat} {Mittens 1 cat}]

	// 11. New utility methods
	fmt.Println("\n--- New Utility Methods ---")

	// Take and Skip
	first3 := numbers.Take(3).Slice()
	fmt.Println("First 3 numbers:", first3) // [1, 2, 3]

	skip2 := numbers.Skip(2).Slice()
	fmt.Println("Skip 2 numbers:", skip2) // [3, 4, 5]

	// At method with negative indices
	if val, ok := numbers.At(2); ok {
		fmt.Printf("Element at index 2: %d\n", val) // 3
	}

	if first, ok := numbers.At(0); ok {
		fmt.Printf("First element: %d\n", first) // 1
	}

	if last, ok := numbers.At(-1); ok {
		fmt.Printf("Last element: %d\n", last) // 5
	}

	// FlatMap
	words := []string{"hello world", "foo bar"}
	allWords := polyfill.FlatMap(polyfill.From(words), func(s string) []string {
		var parts []string
		current := ""
		for _, ch := range s {
			if ch == ' ' {
				if current != "" {
					parts = append(parts, current)
					current = ""
				}
			} else {
				current += string(ch)
			}
		}
		if current != "" {
			parts = append(parts, current)
		}
		return parts
	}).Slice()
	fmt.Println("FlatMap result:", allWords) // [hello world foo bar]

	// 12. Cache example
	c := polyfill.NewCache[string, int](polyfill.Config{
		DefaultTTL: 5 * time.Minute,
		MaxItems:   1000,
		OnEvict: func(k, v any, r polyfill.EvictReason) {
			// log.Printf("evicted %v (%v): %v", k, r, v)
		},
	})

	_ = c.Add("x", 1, -1)         // -1 => no expiration
	c.Set("y", 2, 10*time.Second) // item TTL
	v, ok := c.Get("x")           // 1, true
	_ = c.Update("x", func(p *int) error { *p += 41; return nil })
	val, _ := c.GetOrSet("z", func() (int, time.Duration, error) { return 3, 0, nil })
	_ = val
	_ = ok
	_ = v
}
