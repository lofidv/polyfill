# Changelog

All notable changes to this project are documented here.
This project follows [Semantic Versioning](https://semver.org/).

## [3.0.0]

First release on the v3 module path. The v2 line stays available at `v2.5.0`.

The library was rewritten around one idea: a slice pipeline that carries its own
error, so a failing step stops the chain and keeps the work already done.

### Added

- **Uniform error contract.** `MapE`, `MapToE`, `FlatMapE`, `FlatMapToE`,
  `ReduceE`, `ReduceToE` stop at the first error and keep partial results.
  Every other operation short-circuits once the error is set.
- `doc.go` with the full contract, the two API rules and the mutability rules.
- Typed helpers: `GroupBy`, `UniqueBy`, `SortBy`, `KeyBy`, `ToMap`, `CountBy`.
- Aggregations: `SumBy`, `AvgBy` over a `Number` constraint.
- Go 1.23 iterator bridge: `Seq.Iter` and `FromIter`.
- `Partition` now returns two `*Seq`, so both sides can keep chaining.
- Runnable godoc examples (`example_test.go`) executed by `go test`.
- Three real use cases in `example/main.go`: fail-fast import, sales report,
  and a cache-backed lookup with single-flight.
- Benchmarks for each helper against the equivalent `for` loop.
- CI: `go vet`, `go test -race` and a benchmark smoke test on Go 1.23, 1.25
  and 1.27.

### Fixed

- `Unique` and `IndexOf` compared values with `map[any]` and `==`, so a
  non-comparable element type (a slice, a map, a struct containing one) panicked
  at runtime. `Unique` now falls back to a deep comparison and `IndexOf` uses a
  safe comparison.
- `Chunk` silently returned the whole slice for a size of zero or less. It now
  panics with an explicit message.
- `Parallel` defaulted to one goroutine per element and let a panic in the
  callback kill the process from an arbitrary goroutine. It is now bounded by
  `GOMAXPROCS`, accepts an explicit worker count, and re-raises the panic in the
  caller's goroutine.
- `Take` and `Skip` returned views over the caller's array, so `Push` on the
  result wrote into the original slice. Both now copy.
- The error channel was implemented inconsistently: `GroupBy`, `Partition`,
  `Find`, `Some`, `Every`, `MinBy`, `MaxBy`, `Len`, `First`, `Last`, `Chunk`,
  `Flatten` and the `Parallel` context ignored the chain error. All of them now
  respect it.
- `MapE`, `MapToE`, `FlatMapE` and `FlatMapToE` discarded the elements produced
  before the failure. They now keep them for `SliceE`.
- `ParallelOptions.Ordered` was declared but never used. Replaced by the worker
  count argument.

### Changed (breaking)

- Module path is `github.com/lofidv/polyfill/v3` and the minimum Go version is
  1.23.
- The cache moved to the `polyfill/v3/cache` subpackage. `NewCache` takes a
  `cache.Config`, and `GetOrSet` takes a `context.Context` and runs the supplier
  outside the cache lock.
- `GroupBy`, `UniqueBy`, `SortBy`, `KeyBy`, `ToMap`, `SumBy` and `AvgBy` are
  free functions, because their key or value type is a new type parameter.
- `UniqueBy` takes a comparable key instead of `func(T) any`.

## [1.1.0] and earlier

The v1 line: `Wrap`/`Unwrap`, `Filter`, `Map`, `ParallelMap`, `Reduce`,
`Find`, `IndexOf`, `Some`, `Every`, `Chunk`, `Reverse`, `Sort`, `Unique`, and a
cache with TTL and LRU in the root package.
