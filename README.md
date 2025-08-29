# Polyfill

[![Go report](https://goreportcard.com/badge/github.com/lofidv/polyfill)](https://goreportcard.com/report/github.com/lofidv/polyfill)

✨ **`lofidv/polyfill` is a Go library based on Go 1.18+ Generics.**

This project started as an experiment with the implementation of new generics. It may look like simple js functions but it is fully integrated and functional with arrays and structures.

As expected, generics will be much faster than implementations based on the "reflect" package. Benchmarks also show similar performance gains compared to pure `for` loops.

I feel this library is legitimate and offers many more valuable abstractions.

**Why this name?**

I wanted a **popular name**, similar to "js" and no Go package currently uses this name.

## 🚀 Install

```sh
go get github.com/lofidv/polyfill
```

This library is on beta.

## 💡 Usage

You can import `polyfill` using:

```go
import (
"github.com/lofidv/polyfill"
)
```

Then use it like this:

```go
package main

import (
	"fmt"
	"github.com/lofidv/polyfill"
	"strconv"
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

	// Filter adults and map to names
	adultNames := polyfill.Wrap(people).
		Filter(func(p Person) bool { return p.IsAdult() }).
		Map(func(person Person) any {
			return person.Name
		}).
		Unwrap()

	fmt.Println("Adult names:", adultNames) // ["Alice", "Charlie"]

	// 2. Parallel Map demonstration
	numbers := []int{1, 2, 3, 4, 5}
	squared := polyfill.Wrap(numbers).
		ParallelMap(func(n int) any { return n * n }).
		Unwrap()

	fmt.Println("Squared numbers:", squared) // [1, 4, 9, 16, 25]

	// 3. Reduce example
	sum := polyfill.Wrap(numbers).
		Reduce(0, func(acc any, n int) any { return acc.(int) + n })

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
	chunks := polyfill.Wrap(numbers).Chunk(2)
	fmt.Println("Chunks:")
	for _, chunk := range chunks {
		fmt.Println(chunk.Unwrap())
	}
	// [1 2]
	// [3 4]
	// [5]

	reversed := polyfill.Wrap(numbers).Reverse().Unwrap()
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
	strNumbers := []string{"1", "2", "3", "4"}
	intNumbers := polyfill.Wrap(strNumbers).
		Map(func(s string) any {
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

	_ = c.Add("x", 1, -1)               // -1 => no expiration
	c.Set("y", 2, 10*time.Second)       // item TTL
	v, ok := c.Get("x")                 // 1, true
	_ = c.Update("x", func(p *int) error { *p += 41; return nil })
	val, _ := c.GetOrSet("z", func() (int, time.Duration, error) { return 3, 0, nil })
	_ = val
	_ = ok
	_ = v
}
```

## 🤝 Contributing

- Fork the [project](https://github.com/lofidv/polyfill)
- Fix [open issues](https://github.com/lofidv/polyfill/issues) or request new features

Don't hesitate ;)

## 👤 Authors

- jpastorm

## 🌛 Show your support

Give a ⭐️ if this project helped you!

<p><a href="https://ko-fi.com/lofidev"> <img align="left" src="https://cdn.ko-fi.com/cdn/kofi3.png?v=3" height="50" width="210" alt="lofidev" /></a></p><br><br>

## 📝 License

Copyright © 2025 [jpastorm](https://github.com/lofidv).

This project is [MIT](./LICENSE) licensed.
