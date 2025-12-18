<div align="center">

# 🚀 Polyfill

### *JavaScript-inspired functional programming for Go*

[![Go Report](https://goreportcard.com/badge/github.com/lofidv/polyfill)](https://goreportcard.com/report/github.com/lofidv/polyfill)
[![Go Reference](https://pkg.go.dev/badge/github.com/lofidv/polyfill.svg)](https://pkg.go.dev/github.com/lofidv/polyfill)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

*Bring the power and elegance of JavaScript's array methods to Go with full type safety and zero reflection overhead.*

[Features](#-features) • [Installation](#-installation) • [Quick Start](#-quick-start) • [Documentation](#-documentation) • [Examples](#-examples)

</div>

---

## ✨ Features

<table>
<tr>
<td>

🎯 **JavaScript-Inspired**
<br/>Familiar API from JS/TS with Go's type safety

</td>
<td>

⚡ **Zero Overhead**
<br/>Pure generics, no reflection, blazing fast

</td>
</tr>
<tr>
<td>

🔗 **Chainable**
<br/>Fluent interface for elegant pipelines

</td>
<td>

🛡️ **Type Safe**
<br/>Compile-time guarantees, no runtime panics

</td>
</tr>
<tr>
<td>

🚄 **Parallel Ready**
<br/>Built-in concurrent processing

</td>
<td>

💎 **Rich API**
<br/>50+ methods for array manipulation

</td>
</tr>
</table>

## 🎯 Why Polyfill?

```go
// Before: Verbose Go loops
var adults []string
for _, p := range people {
    if p.Age >= 18 {
        adults = append(adults, p.Name)
    }
}

// After: Expressive Polyfill
adults := polyfill.MapTo(
    polyfill.From(people).Filter(func(p Person) bool { 
        return p.Age >= 18 
    }),
    func(p Person) string { return p.Name },
).Slice()
```

## 📦 Installation

```bash
go get github.com/lofidv/polyfill
```

**Requirements:** Go 1.18+ (generics support)

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "github.com/lofidv/polyfill"
)

type Person struct {
    Name string
    Age  int
}

func main() {
    people := []Person{
        {"Alice", 25},
        {"Bob", 17},
        {"Charlie", 19},
    }

    // Filter adults and get their names
    adultNames := polyfill.MapTo(
        polyfill.From(people).Filter(func(p Person) bool { 
            return p.Age >= 18 
        }),
        func(p Person) string { return p.Name },
    ).Slice()

    fmt.Println(adultNames) // [Alice Charlie]
}
```

## 📚 Documentation

### Core API

#### Creating Sequences

```go
// From slice
seq := polyfill.From([]int{1, 2, 3, 4, 5})

// Back to slice
numbers := seq.Slice()

// With error handling
numbers, err := seq.SliceE()
```

### 🔍 Filtering & Searching

<details>
<summary><b>Filter, Find, Some, Every</b></summary>

```go
// Filter elements
adults := polyfill.From(people).
    Filter(func(p Person) bool { return p.Age >= 18 }).
    Slice()

// Find first match
person, found := polyfill.From(people).
    Find(func(p Person) bool { return p.Name == "Alice" })

// Find index
index := polyfill.From(people).
    FindIndex(func(p Person) bool { return p.Name == "Bob" })

// Check if any element matches
hasAdults := polyfill.From(people).
    Some(func(p Person) bool { return p.Age >= 18 })

// Check if all elements match
allAdults := polyfill.From(people).
    Every(func(p Person) bool { return p.Age >= 18 })
```

</details>

### 🔄 Transforming

<details>
<summary><b>Map, FlatMap, Reduce</b></summary>

**Important:** Due to Go's limitations on method type parameters:
- Use **methods** for same-type transformations (T → T)
- Use **functions** for type-changing transformations (T → R)

```go
// Map same type (method)
doubled := polyfill.From([]int{1, 2, 3}).
    Map(func(n int) int { return n * 2 }).
    Slice()

// Map with type change (function)
names := polyfill.MapTo(
    polyfill.From(people),
    func(p Person) string { return p.Name },
).Slice()

// Map with error handling
numbers, err := polyfill.MapToE(
    polyfill.From([]string{"1", "2", "3"}),
    strconv.Atoi,
).SliceE()

// FlatMap
words := polyfill.FlatMap(
    polyfill.From([]string{"hello world", "foo bar"}),
    func(s string) []string { return strings.Split(s, " ") },
).Slice()

// Reduce
sum := polyfill.Reduce(
    polyfill.From([]int{1, 2, 3, 4}),
    0,
    func(acc, n int) int { return acc + n },
) // 10
```

</details>

### 🎯 Grouping & Partitioning

<details>
<summary><b>GroupBy, Partition, Unique</b></summary>

```go
// GroupBy - killer feature!
grouped := polyfill.GroupBy(
    polyfill.From(pets),
    func(p Pet) string { return p.Type },
)
// Returns: map[string][]Pet

// Partition into two groups
adults, minors := polyfill.From(people).
    Partition(func(p Person) bool { return p.Age >= 18 })

// Unique for comparable types
unique := polyfill.Unique(
    polyfill.From([]int{1, 2, 2, 3, 3, 4}),
).Slice() // [1 2 3 4]

// UniqueBy with custom key
uniquePeople := polyfill.UniqueBy(
    polyfill.From(people),
    func(p Person) string { return p.Name },
).Slice()
```

</details>

### 🔀 Sorting & Reversing

<details>
<summary><b>Sort, Reverse (Immutable)</b></summary>

```go
// Sort (returns new sequence)
sorted := polyfill.From([]int{3, 1, 4, 2}).
    Sort(func(a, b int) bool { return a < b }).
    Slice() // [1 2 3 4]

// Reverse (returns new sequence)
reversed := polyfill.From([]int{1, 2, 3}).
    Reverse().
    Slice() // [3 2 1]

// Original unchanged - immutable operations!
```

</details>

### ⚡ Parallel Execution

<details>
<summary><b>Process data concurrently</b></summary>

```go
// Basic parallel map
squared := polyfill.From(numbers).
    Parallel().
    Map(func(n int) int { return n * n }).
    Slice()

// Configure workers and ordering
results := polyfill.From(largeDataset).
    Parallel().
    Workers(8).              // 8 concurrent workers
    Unordered().             // Don't preserve order (faster)
    Map(expensiveFunc).
    Slice()

// Parallel with type change
results := polyfill.From(urls).
    Parallel().
    ParallelMapTo(fetchData).
    Slice()
```

</details>

### 🔧 Utility Methods

<details>
<summary><b>Take, Skip, At, Chunk, and more</b></summary>

```go
// Take first N
first3 := polyfill.From(numbers).Take(3).Slice()

// Skip first N
rest := polyfill.From(numbers).Skip(2).Slice()

// Access with negative indices (Python-style)
last, ok := polyfill.From(numbers).At(-1)
first, ok := seq.At(0)

// Chunk into groups
chunks := polyfill.From([]int{1, 2, 3, 4, 5}).Chunk(2)
// [[1 2] [3 4] [5]]

// Check emptiness
if seq.IsEmpty() { }

// Get length
length := seq.Len()
```

</details>

## 💡 Examples

### Data Pipeline

```go
type Product struct {
    Name  string
    Price float64
    Stock int
}

// Complex transformation pipeline
expensive := polyfill.MapTo(
    polyfill.From(products).
        Filter(func(p Product) bool { return p.Stock > 0 }).
        Sort(func(a, b Product) bool { return a.Price > b.Price }).
        Take(5),
    func(p Product) string { return p.Name },
).Slice()
```

### Error Handling

```go
// Parse strings to integers with error propagation
numbers, err := polyfill.MapToE(
    polyfill.From([]string{"1", "2", "invalid", "4"}),
    strconv.Atoi,
).SliceE()

if err != nil {
    log.Fatal(err) // Catches parse error
}
```

### Grouping Data

```go
type Order struct {
    Customer string
    Amount   float64
}

// Group orders by customer
byCustomer := polyfill.GroupBy(
    polyfill.From(orders),
    func(o Order) string { return o.Customer },
)

// Calculate totals per customer
for customer, orders := range byCustomer {
    total := polyfill.Reduce(
        polyfill.From(orders),
        0.0,
        func(sum float64, o Order) float64 { return sum + o.Amount },
    )
    fmt.Printf("%s: $%.2f\n", customer, total)
}
```

### Parallel Processing

```go
// Process large dataset concurrently
results := polyfill.From(imageURLs).
    Parallel().
    Workers(runtime.NumCPU()).
    ParallelMapTo(downloadAndProcess).
    Slice()
```

## 🎨 Design Philosophy

### Method vs Function Pattern

Go doesn't allow methods with additional type parameters. Our solution:

```go
// ✅ Same type (T → T) - use METHOD
seq.Map(func(n int) int { return n * 2 })

// ✅ Type change (T → R) - use FUNCTION
polyfill.MapTo(seq, func(n int) string { return fmt.Sprint(n) })
```

This maintains clean, readable code while respecting Go's constraints.

### Immutability

Operations like `Sort()` and `Reverse()` return new sequences without modifying originals:

```go
original := polyfill.From([]int{3, 1, 2})
sorted := original.Sort(less)

// original: [3 1 2] ✅ unchanged
// sorted:   [1 2 3] ✅ new sequence
```

## 🎯 Key Methods at a Glance

| Category | Methods |
|----------|---------|
| **Creation** | `From()`, `Slice()`, `SliceE()` |
| **Filtering** | `Filter()`, `Find()`, `FindIndex()`, `Some()`, `Every()` |
| **Transform** | `Map()`, `MapTo()`, `FlatMap()`, `Flatten()` |
| **Reduce** | `Reduce()`, `ReduceE()`, `ReduceRight()` |
| **Grouping** | `GroupBy()`, `Partition()`, `Unique()`, `UniqueBy()` |
| **Sorting** | `Sort()`, `Reverse()` |
| **Slicing** | `Take()`, `Skip()`, `At()`, `Chunk()` |
| **Parallel** | `Parallel()`, `Workers()`, `Unordered()` |
| **Utility** | `IsEmpty()`, `Len()`, `IndexOf()` |

## 📊 Performance

- ✅ **Zero reflection** - pure generics
- ✅ **Inline-friendly** - compiler optimizations
- ✅ **Memory efficient** - minimal allocations
- ✅ **Parallel ready** - leverage all cores

Benchmarks show performance comparable to hand-written loops with much better readability.

## 🎓 Advanced Features

### Cache (Bonus)

Thread-safe in-memory cache with TTL:

```go
cache := polyfill.NewCache[string, User](polyfill.Config{
    DefaultTTL: 5 * time.Minute,
    MaxItems:   1000,
})

cache.Set("user:123", user, 0)
if user, ok := cache.Get("user:123"); ok {
    // Use cached user
}
```

## 🤝 Contributing

Contributions are welcome! Feel free to:

- 🐛 Report bugs
- 💡 Suggest features  
- 🔧 Submit pull requests
- ⭐ Star the project

Visit [github.com/lofidv/polyfill](https://github.com/lofidv/polyfill)

## 📄 License

MIT License - see [LICENCE](./LICENCE) for details

---

<div align="center">

**Made with ❤️ by [jpastorm](https://github.com/lofidv)**

If this project helped you, consider [buying me a coffee](https://ko-fi.com/lofidev) ☕

⭐ **Star us on GitHub** — it motivates us a lot!

</div>
