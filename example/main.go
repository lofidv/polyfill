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

	wrappedPeople := polyfill.Wrap(people)
	adults := wrappedPeople.Filter(func(p Person) bool {
		return p.IsAdult()
	})

	// Filter adults and map to names
	adultNames := polyfill.Map(adults, func(p Person) string {
		return p.Name
	}).Unwrap()

	fmt.Println("Adult names:", adultNames) // ["Alice", "Charlie"]

	// 2. Parallel Map demonstration
	numbers := polyfill.Wrap([]int{1, 2, 3, 4, 5})
	squared := polyfill.ParallelMap(numbers, func(n int) any { return n * n }).Unwrap()

	fmt.Println("Squared numbers:", squared) // [1, 4, 9, 16, 25]

	// 3. Reduce example
	sum := numbers.Reduce(0, func(acc any, n int) any { return acc.(int) + n })

	fmt.Println("Sum of numbers:", sum) // 15

	// 4. Find and IndexOf
	bob, found := polyfill.Wrap(people).
		Find(func(p Person) bool { return p.Name == "Bob" })

	fmt.Printf("Found Bob: %v (%v)\n", bob, found) // {Bob 17} true

	bobIndex := polyfill.Wrap(people).
		IndexOf(func(p Person) bool { return p.Name == "Bob" })

	fmt.Println("Bob's index:", bobIndex) // 1

	// 5. Some and Every
	hasAdults := polyfill.Wrap(people).
		Some(func(p Person) bool { return p.IsAdult() })

	allAdults := polyfill.Wrap(people).
		Every(func(p Person) bool { return p.IsAdult() })

	fmt.Printf("Has adults: %v, All adults: %v\n", hasAdults, allAdults) // true, false

	// 6. Chunk and Reverse
	chunks := numbers.Chunk(2)
	fmt.Println("Chunks:")
	for _, chunk := range chunks {
		fmt.Println(chunk.Unwrap())
	}
	// [1 2]
	// [3 4]
	// [5]

	reversed := numbers.Reverse().Unwrap()
	fmt.Println("Reversed numbers:", reversed) // [5, 4, 3, 2, 1]

	// 7. Sort
	unsorted := []int{3, 1, 4, 2}
	sorted := polyfill.Wrap(unsorted).
		Sort(func(a, b int) bool { return a < b }).
		Unwrap()

	fmt.Println("Sorted numbers:", sorted) // [1, 2, 3, 4]

	// 8. Unique with custom equality
	duplicates := []int{1, 2, 2, 3, 4, 4, 5}
	unique := polyfill.Wrap(duplicates).
		Unique(func(a, b int) bool { return a == b }).
		Unwrap()

	fmt.Println("Unique numbers:", unique) // [1, 2, 3, 4, 5]

	// 9. String conversion
	strNumbers := polyfill.Wrap([]string{"1", "2", "3", "4"})
	intNumbers := polyfill.Map(strNumbers, func(s string) any {
		n, _ := strconv.Atoi(s)
		return n
	}).
		Unwrap()

	fmt.Println("Converted numbers:", intNumbers) // [1, 2, 3, 4]

	// 10. Complex Reduce - group pets by type
	pets := []Pet{
		{"Fido", 3, "dog"},
		{"Whiskers", 2, "cat"},
		{"Rover", 5, "dog"},
		{"Mittens", 1, "cat"},
	}

	type PetGroup map[string][]Pet
	grouped := polyfill.Wrap(pets).
		Reduce(make(PetGroup), func(acc any, p Pet) any {
			group := acc.(PetGroup)
			group[p.Type] = append(group[p.Type], p)
			return group
		}).(PetGroup)

	fmt.Println("Pets grouped by type:")
	for typ, pets := range grouped {
		fmt.Printf("%s: %v\n", typ, pets)
	}
	// dog: [{Fido 3 dog} {Rover 5 dog}]
	// cat: [{Whiskers 2 cat} {Mittens 1 cat}]

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
