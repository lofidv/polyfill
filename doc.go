// Package polyfill provides chainable, error-aware pipelines over Go slices.
//
// It exists for one job: make a multi-step slice transformation read like a
// sentence instead of a nest of loops, without giving up Go's type safety.
//
//	paid := polyfill.From(orders).Filter(func(o Order) bool { return o.Paid })
//	report := polyfill.GroupBy(paid, func(o Order) string { return o.Region })
//
// # Why not just use a for loop?
//
// For a hot path over millions of elements, use a for loop: the pipeline
// allocates one slice per step and a loop does not. For the rest of the code,
// where an order pipeline or a report has 20 to 20,000 rows, the pipeline is
// easier to read, easier to test and harder to get wrong at the boundaries.
//
// See the benchmarks in the README for the real cost.
//
// # The error contract
//
// A Seq carries an optional error. The rule is uniform across the package:
//
//  1. Once the error is set, every operation is a no-op. Transformations
//     return the receiver unchanged.
//  2. Operations that can fail stop at the first failure. The elements
//     processed so far are kept.
//  3. Terminal operations that return a value, a bool or a collection return
//     the zero value, false or an empty result. They do not report the error.
//  4. To observe an error, use SliceE, which returns both the partial result
//     and the error, or Err.
//
// Example of a pipeline that fails on bad input:
//
//	rows, err := polyfill.From(rawRows).
//		MapE(parseRow).          // stops at the first malformed row
//		Filter(func(r Row) bool { return r.Active }).
//		SliceE()
//
// # Mutability
//
// From and New do not copy the input slice. Filter, Map, Sort, Reverse, Concat,
// Append, Prepend, Take, Skip and Clone return new slices. Push is the only
// method that mutates the receiver, which matches JavaScript's Array.push.
//
// # Type changes
//
// A method cannot declare a type parameter of its own before Go 1.27, so
// operations that change the element type are free functions that take the Seq
// as their first argument:
//
//	polyfill.MapTo(From(users), func(u User) string { return u.Email })
//	polyfill.GroupBy(From(orders), func(o Order) string { return o.Region })
//
// Operations that keep the type are methods, so they chain.
//
// The package requires Go 1.23, which is also where Iter and FromIter bridge to
// the standard range-over-func iterators.
package polyfill
