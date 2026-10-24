# Polyfill v3

[![Go Reference](https://pkg.go.dev/badge/github.com/lofidv/polyfill/v3.svg)](https://pkg.go.dev/github.com/lofidv/polyfill/v3)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENCE)

Chainable slice pipelines for Go, with a uniform error channel and typed keys.

The purpose is narrow on purpose: **make a multi-step slice transformation read
like a sentence**, and make a failing step stop the pipeline instead of
scattering `if err != nil` through your code.

```go
paid := polyfill.From(orders).Filter(func(o Order) bool { return o.Paid })

report := polyfill.GroupBy(paid, func(o Order) string { return o.Region })
revenue := polyfill.SumBy(paid, func(o Order) float64 { return o.Amount })
```

```sh
go get github.com/lofidv/polyfill/v3
```

Requires Go 1.23 or later. No third-party dependencies.

---

## The reason this library exists: the error channel

A step that parses or validates external data can fail. In plain Go you write a
loop with an early return, plus an error variable. Here the failure travels with
the pipeline and stops it:

```go
users, err := polyfill.MapToE(polyfill.From(csvRows), parseRow).SliceE()
if err != nil {
    // rows accepted before the failure are still available in users
}
```

Three properties make this predictable:

1. **The first error stops the chain.** No step runs after it.
2. **Partial results are kept.** `SliceE` returns the elements produced before
   the failure, which is what you want when importing 10,000 rows and row 3 is
   malformed.
3. **Everything else is a no-op.** After the error, `Filter`, `Sort`, `GroupBy`,
   `Partition` and the rest do nothing, and the value-returning methods give you
   zero values.

That is the whole contract; it is documented in full in
[`doc.go`](./doc.go). This is the feature that distinguishes Polyfill from
`samber/lo` and the standard `slices` package, which have no notion of a failing
chain.

---

## The two rules that explain the API

Go does not allow a method to declare a type parameter of its own (until Go
1.27), so the API splits cleanly:

| Rule | Examples |
|---|---|
| **Same element type → method**, so it chains | `Filter`, `Map`, `Sort`, `Reverse`, `Take`, `Skip`, `Unique`, `Partition`, `Find`, `MinBy`, `MaxBy` |
| **New type parameter → free function**, takes the `Seq` first | `MapTo`, `MapToE`, `ReduceTo`, `GroupBy`, `UniqueBy`, `SortBy`, `KeyBy`, `ToMap`, `SumBy`, `ParallelMapTo` |

`SortBy` is the one same-type operation that is a free function, because the
*key* type is a new type parameter. To sort inside a chain, use `Sort` with a
comparison:

```go
polyfill.From(orders).
    Filter(isPaid).
    Sort(func(a, b Order) bool { return a.Amount > b.Amount }).
    Take(10)
```

---

## Use cases

The full runnable version is in [`example/main.go`](./example/main.go)
(`go run ./example`).

### 1. Import external data with fail-fast validation

```go
users, err := polyfill.MapToE(polyfill.From(csvRows), parseRow).SliceE()
```

Output of the demo: two rows accepted, then it stops at the malformed email and
reports which row failed.

### 2. Build a report

```go
paid, unpaid := polyfill.From(orders).Partition(func(o Order) bool { return o.Paid })

revenue := polyfill.SumBy(paid, func(o Order) float64 { return o.Amount })
average := polyfill.AvgBy(paid, func(o Order) float64 { return o.Amount })

byRegion := polyfill.GroupBy(paid, func(o Order) string { return o.Region })
counts := polyfill.CountBy(paid, func(o Order) string { return o.Region })

top := polyfill.SortBy(paid, func(o Order) float64 { return o.Amount }).Reverse().Take(3)

for i, page := range polyfill.From(orders).Chunk(20) { /* batch or paginate */ }
```

### 3. Serve a hot lookup through the cache

The companion package `polyfill/v3/cache` collapses concurrent misses:

```go
c := cache.NewCache[string, Product](cache.Config{
    DefaultTTL: 5 * time.Minute,
    MaxItems:   10_000,
})

p, err := c.GetOrSet(ctx, sku, func(ctx context.Context) (Product, time.Duration, error) {
    return db.FindProduct(ctx, sku) // runs outside the cache lock
})
```

In the demo, **200 concurrent requests produce 1 repository query**.

---

## API at a glance

**Build and exit:** `From`, `New`, `Slice`, `SliceE`, `Err`, `Len`, `IsEmpty`,
`Clone`, `Push`, `Append`, `Prepend`, `Concat`, `At`, `Get`, `First`, `Last`,
`Take`, `Skip`, `Iter`, `FromIter`.

**Transform:** `Filter`, `Map`, `MapE`, `MapTo`, `MapToE`, `FlatMap`, `FlatMapE`,
`FlatMapTo`, `FlatMapToE`, `Flatten`, `Reverse`, `Sort`, `SortBy`, `Unique`,
`UniqueBy`, `Chunk`, `Partition`, `Parallel`.

**Query:** `Find`, `FindIndex`, `IndexOf`, `Includes`, `ContainsFunc`, `Some`,
`Every`, `MinBy`, `MaxBy`, `ForEach`, `ForEachIndexed`, `EqualFunc`.

**Aggregate and index:** `Reduce`, `ReduceE`, `ReduceRight`, `ReduceTo`,
`ReduceToE`, `SumBy`, `AvgBy`, `GroupBy`, `CountBy`, `KeyBy`, `ToMap`.

**Cache (`polyfill/v3/cache`):** `NewCache`, `Get`, `Peek`, `Set`, `Add`,
`GetOrSet`, `Update`, `Touch`, `Delete`, `Has`, `TTL`, `Len`, `Keys`, `Sweep`,
`Clear`.

---

## Mutability

`From` and `New` do not copy: the Seq shares the caller's array.

Everything that returns a new Seq copies: `Filter`, `Map`, `Sort`, `Reverse`,
`Concat`, `Append`, `Prepend`, `Take`, `Skip`, `Clone`.

`Push` is the only mutating method, and it exists because JavaScript's
`Array.push` mutates. Use `Append` when you want a new sequence.

`Chunk` returns views into the receiver's array. Treat the chunks as read-only.

---

## Go 1.23 iterators

```go
for u := range polyfill.From(users).Iter() {
    fmt.Println(u.Name)
}

total := polyfill.SumBy(polyfill.FromIter(slices.Values(numbers)), func(n int) int { return n })
```

`FromIter` accepts any `iter.Seq`, so the standard `slices` and `maps` packages
can feed a pipeline.

---

## Benchmarks

Go 1.27, Apple M4, `-benchtime=200ms -benchmem`. Reproduce with:

```sh
go test -run '^$' -bench . -benchmem ./...
```

1000 elements:

| Pipeline step | ns/op | Allocs |
|---|---|---|
| `Filter` | 1 253 | 1 |
| equivalent `for` loop | 763 | 0 |
| `Map` | 728 | 1 |
| equivalent `for` loop | 251 | 0 |
| `SumBy` | 533 | 0 |
| equivalent `for` loop | 533 | 0 |
| `SortBy` | 32 790 | 2 |
| `CountBy` | 4 603 | 3 |

A realistic 4-step chain over **10 000** rows (`Filter` → `Map` → `Sort` →
`Take(100)`): **211 µs**, 5 allocations, 165 KB.

Parallel, 50 elements with a 100 µs callback:

| | ns/op |
|---|---|
| sequential loop | 5 910 663 |
| `Parallel(8).Map` | 827 511 |

**And the part most READMEs hide:** `ParallelMap` with a *cheap* callback took
241 µs, versus 0.7 µs for `Map`. That is roughly **300x slower**. Parallelism
only pays when each element costs real time (I/O, hashing, compression).

---

## When not to use this

- **Hot paths over millions of elements.** Each step allocates a slice. Write
  the loop; the compiler will thank you.
- **Two-step transformations.** `for` is shorter.
- **You need lazy evaluation.** `FromIter` + `Iter` give you the bridge to
  `iter.Seq`, but the Seq itself is eager.

Use it where a report, an import, or an API response handler has 20 to 20,000
rows and the *intent* matters more than the last microsecond.

---

## Migrating from v1

| v1 | v3 |
|---|---|
| `polyfill.Wrap(x)` / `.Unwrap()` | `polyfill.From(x)` / `.Slice()` |
| `polyfill.Map(s, f)` | `polyfill.MapTo(s, f)` |
| `Run(seq).Map(...)` same-type | `Seq.Map(...)` is now same-type |
| `Unique(equals)` | `Unique()` and `UniqueBy(key)` |
| `GroupBy(func(T) any)` | `GroupBy(s, func(T) K)` typed |
| `polyfill.NewCache` | `cache.NewCache` (subpackage) |
| Go 1.18 | Go 1.23 |

See [CHANGELOG.md](./CHANGELOG.md) for the full list.

---

## Contributing

```sh
go vet ./...
go test -race ./...
```

## License

MIT. Copyright © 2025 [jpastorm](https://github.com/jpastorm).
See [LICENCE](./LICENCE).
