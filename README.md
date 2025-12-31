# 🚀 Polyfill

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/lofidv/polyfill.svg)](https://pkg.go.dev/github.com/lofidv/polyfill)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A modern, type-safe functional programming library for Go, bringing the power and elegance of JavaScript's array methods to Go slices with full generic support.

## ✨ Features

- 🎯 **Fully Generic** - Works with any type using Go 1.23 generics
- ⛓️ **Chainable API** - Fluent, readable method chaining
- 🚄 **High Performance** - Optimized with Go 1.23's `slices` package
- 🧪 **100% Tested** - Comprehensive test coverage
- 📦 **Zero Dependencies** - Only uses Go standard library
- 🎨 **JavaScript-Inspired** - Familiar API for web developers
- 🔧 **Type-Safe** - Compile-time type checking
- ⚡ **Parallel Support** - Built-in parallel processing

## 📦 Installation

```bash
go get github.com/lofidv/polyfill
```

## 🎯 Quick Start

```go
package main

import (
    "fmt"
    "github.com/lofidv/polyfill"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    // Chain operations fluently
    result := polyfill.From(numbers).
        Filter(func(n int) bool { return n%2 == 0 }).
        Map(func(n int) int { return n * 2 }).
        Slice()
    
    fmt.Println(result) // [4, 8]
}
```

## 🎨 Core Methods

### Transformations

```go
// Map - Transform elements (same type)
doubled := From([]int{1, 2, 3}).
    Map(func(n int) int { return n * 2 }).
    Slice() // [2, 4, 6]

// MapTo - Transform with type change
strings := MapTo(From([]int{1, 2, 3}), 
    func(n int) string { return fmt.Sprint(n) }).
    Slice() // ["1", "2", "3"]

// FlatMap - Map and flatten
words := From([]string{"hello world", "foo bar"}).
    FlatMap(func(s string) []string { return strings.Split(s, " ") }).
    Slice() // ["hello", "world", "foo", "bar"]
```

### Filtering & Searching

```go
// Filter - Keep matching elements
evens := From([]int{1, 2, 3, 4}).
    Filter(func(n int) bool { return n%2 == 0 }).
    Slice() // [2, 4]

// Find - Get first match
first, found := From([]int{1, 2, 3}).
    Find(func(n int) bool { return n > 1 }) // 2, true

// Some/Every - Check conditions
hasEven := From([]int{1, 2, 3}).
    Some(func(n int) bool { return n%2 == 0 }) // true

allPositive := From([]int{1, 2, 3}).
    Every(func(n int) bool { return n > 0 }) // true
```

### Aggregations

```go
// Reduce - Aggregate to single value
sum := From([]int{1, 2, 3, 4}).
    Reduce(0, func(acc, n int) int { return acc + n }) // 10

// MinBy/MaxBy - Find extremes
min, _ := From([]int{3, 1, 4}).
    MinBy(func(a, b int) bool { return a < b }) // 1

max, _ := From([]int{3, 1, 4}).
    MaxBy(func(a, b int) bool { return a < b }) // 4
```

### Utilities

```go
// Sort - Custom sorting
sorted := From([]int{3, 1, 4, 2}).
    Sort(func(a, b int) bool { return a < b }).
    Slice() // [1, 2, 3, 4]

// Reverse - Flip order
reversed := From([]int{1, 2, 3}).
    Reverse().
    Slice() // [3, 2, 1]

// Unique - Remove duplicates
unique := From([]int{1, 2, 2, 3, 3}).
    Unique().
    Slice() // [1, 2, 3]

// Chunk - Split into groups
chunks := From([]int{1, 2, 3, 4, 5}).
    Chunk(2) // [[1, 2], [3, 4], [5]]

// GroupBy - Group by key
byParity := From([]int{1, 2, 3, 4}).
    GroupBy(func(n int) any { return n % 2 })
// map[0:[2, 4] 1:[1, 3]]
```

## 🆕 Go 1.23 Features

```go
// Concat - Combine slices
combined := From([]int{1, 2}).
    Concat([]int{3, 4}, []int{5, 6}).
    Slice() // [1, 2, 3, 4, 5, 6]

// Prepend - Add to beginning
result := From([]int{3, 4}).
    Prepend(1, 2).
    Slice() // [1, 2, 3, 4]

// Clone - Efficient copying (uses slices.Clone)
clone := From([]int{1, 2, 3}).Clone()

// ContainsFunc - Flexible search
hasEven := From([]int{1, 3, 5, 8}).
    ContainsFunc(func(n int) bool { return n%2 == 0 }) // true

// EqualFunc - Custom equality
equal := From([]int{1, 2}).
    EqualFunc([]int{1, 2}, func(a, b int) bool { return a == b }) // true
```

## 🔥 Advanced Examples

### Working with Structs

```go
type Person struct {
    Name string
    Age  int
    City string
}

people := []Person{
    {"Alice", 28, "NYC"},
    {"Bob", 35, "SF"},
    {"Charlie", 22, "NYC"},
}

// Complex filtering and transformation
nycAdults := From(people).
    Filter(func(p Person) bool { 
        return p.City == "NYC" && p.Age >= 25 
    }).
    Sort(func(a, b Person) bool { 
        return a.Age < b.Age 
    }).
    Slice()

// Extract field values
names := MapTo(From(people), 
    func(p Person) string { return p.Name }).
    Slice()

// Group by city
byCity := From(people).
    GroupBy(func(p Person) any { return p.City })

// Calculate average age
totalAge := ReduceTo(From(people), 0,
    func(acc int, p Person) int { return acc + p.Age })
avgAge := float64(totalAge) / float64(len(people))
```

### Parallel Processing

```go
// Process large datasets in parallel
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

result := From(numbers).
    Parallel().
    Map(func(n int) int { 
        // Expensive computation
        return n * n 
    }).
    Slice()
```

### Error Handling

```go
// MapE and ReduceE for error handling
result := From([]string{"1", "2", "abc"}).
    MapE(func(s string) (int, error) {
        return strconv.Atoi(s)
    })

values, err := result.SliceE()
if err != nil {
    // Handle error
}
```

## 📊 Performance

Polyfill is optimized using Go 1.23's standard library:

- **`slices.Clone`** - Efficient slice copying
- **`slices.Reverse`** - Optimized reversal
- **`slices.SortFunc`** - Generic sorting
- **Minimal allocations** - Preallocated slices where possible

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run benchmarks
go test -bench=.
```

## 📚 Complete API Reference

### Creation
- `From[T]([]T) *Seq[T]` - Create from slice
- `New[T](...T) *Seq[T]` - Create from variadic args

### Transformations
- `Map(f func(T) T) *Seq[T]` - Transform same type
- `MapTo[R](f func(T) R) *Seq[R]` - Transform to new type
- `FlatMap(f func(T) []T) *Seq[T]` - Map and flatten
- `FlatMapTo[R](f func(T) []R) *Seq[R]` - FlatMap with type change

### Filtering
- `Filter(f func(T) bool) *Seq[T]` - Keep matching
- `Unique() *Seq[T]` - Remove duplicates
- `UniqueBy(f func(T) any) *Seq[T]` - Unique by key

### Searching
- `Find(f func(T) bool) (T, bool)` - First match
- `FindIndex(f func(T) bool) int` - Index of first match
- `IndexOf(value T) int` - Index of value
- `Includes(value T) bool` - Contains value
- `ContainsFunc(f func(T) bool) bool` - Custom contains

### Aggregation
- `Reduce(initial T, f func(T, T) T) T` - Aggregate
- `ReduceTo[R](initial R, f func(R, T) R) R` - Reduce with type change
- `MinBy(less func(T, T) bool) (T, bool)` - Minimum
- `MaxBy(less func(T, T) bool) (T, bool)` - Maximum

### Validation
- `Some(f func(T) bool) bool` - Any match
- `Every(f func(T) bool) bool` - All match

### Ordering
- `Sort(less func(T, T) bool) *Seq[T]` - Sort
- `Reverse() *Seq[T]` - Reverse order

### Utilities
- `Chunk(size int) [][]T` - Split into chunks
- `GroupBy(f func(T) any) map[any][]T` - Group by key
- `Partition(f func(T) bool) ([]T, []T)` - Split by predicate
- `Concat(...[]T) *Seq[T]` - Combine slices
- `Append(...T) *Seq[T]` - Append elements
- `Prepend(...T) *Seq[T]` - Prepend elements
- `Clone() *Seq[T]` - Clone sequence
- `Take(n int) *Seq[T]` - First n elements
- `Skip(n int) *Seq[T]` - Skip n elements

### Access
- `Slice() []T` - Get underlying slice
- `SliceE() ([]T, error)` - Get slice with error
- `First() (T, bool)` - First element
- `Last() (T, bool)` - Last element
- `At(index int) (T, bool)` - Element at index
- `Len() int` - Length
- `IsEmpty() bool` - Check if empty

### Iteration
- `ForEach(f func(T))` - Iterate
- `ForEachIndexed(f func(int, T))` - Iterate with index

### Parallel
- `Parallel() *ParallelSeq[T]` - Enable parallel processing

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by JavaScript's array methods
- Built with Go 1.23 generics and standard library
- Community feedback and contributions

## 🔗 Links

- [Documentation](https://pkg.go.dev/github.com/lofidv/polyfill)
- [Examples](./example)
- [Go 1.23 Upgrade Guide](./GO_1.23_UPGRADE.md)
- [Issues](https://github.com/lofidv/polyfill/issues)

---

Made with ❤️ for the Go community
