// Command example runs three realistic use cases for the library:
//
//  1. Importing external rows with fail-fast validation.
//  2. Building a sales report: group, aggregate, rank and paginate.
//  3. Serving a hot lookup through the cache with single-flight loading.
//
// Run it with:
//
//	go run ./example
package main

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/lofidv/polyfill/v3"
	"github.com/lofidv/polyfill/v3/cache"
)

func main() {
	useCaseImport()
	useCaseReport()
	useCaseCache()
}

// ---------------------------------------------------------------- use case 1

type User struct {
	ID    int
	Name  string
	Email string
	Plan  string
}

// rawRows simulates a CSV export. Row 4 is malformed on purpose.
var rawRows = []string{
	"1,ana,ana@example.com,pro",
	"2,bob,bob@example.com,free",
	"3,carol,not-an-email,pro",
	"4,dave,dave@example.com,free",
}

// parseRow validates one row. A single bad row must not silently disappear.
func parseRow(row string) (User, error) {
	parts := strings.Split(row, ",")
	if len(parts) != 4 {
		return User{}, fmt.Errorf("expected 4 columns, got %d", len(parts))
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return User{}, fmt.Errorf("invalid id %q: %w", parts[0], err)
	}
	if !strings.Contains(parts[2], "@") {
		return User{}, fmt.Errorf("row %d: invalid email %q", id, parts[2])
	}
	return User{ID: id, Name: parts[1], Email: parts[2], Plan: parts[3]}, nil
}

func useCaseImport() {
	fmt.Println("== Use case 1: import with fail-fast validation ==")

	// Parsing changes the element type, so it is MapToE, not MapE.
	users, err := polyfill.MapToE(polyfill.From(rawRows), parseRow).SliceE()

	if err != nil {
		fmt.Printf("import stopped: %v\n", err)
		fmt.Printf("rows accepted before the failure: %d\n", len(users))
		for _, u := range users {
			fmt.Printf("  #%d %s (%s)\n", u.ID, u.Name, u.Plan)
		}
		fmt.Println()
		return
	}
}

// ---------------------------------------------------------------- use case 2

type Order struct {
	ID       int
	Region   string
	Customer string
	Amount   float64
	Paid     bool
}

func orders() []Order {
	return []Order{
		{1, "eu", "ana", 120.00, true},
		{2, "us", "bob", 80.50, true},
		{3, "eu", "carol", 45.00, false},
		{4, "apac", "dave", 300.00, true},
		{5, "us", "erin", 60.00, true},
		{6, "eu", "frank", 200.00, true},
		{7, "apac", "grace", 40.00, true},
	}
}

func useCaseReport() {
	fmt.Println("== Use case 2: sales report ==")

	all := polyfill.From(orders())

	paid, unpaid := all.Partition(func(o Order) bool { return o.Paid })
	fmt.Printf("paid: %d  unpaid: %d\n", paid.Len(), unpaid.Len())

	fmt.Printf("paid revenue: %.2f\n", polyfill.SumBy(paid, func(o Order) float64 { return o.Amount }))
	fmt.Printf("average ticket: %.2f\n", polyfill.AvgBy(paid, func(o Order) float64 { return o.Amount }))

	byRegion := polyfill.GroupBy(paid, func(o Order) string { return o.Region })
	regions := make([]string, 0, len(byRegion))
	for region := range byRegion {
		regions = append(regions, region)
	}
	sort.Strings(regions)

	fmt.Println("revenue by region:")
	for _, region := range regions {
		fmt.Printf("  %-4s %8.2f (%d orders)\n",
			region,
			polyfill.SumBy(polyfill.From(byRegion[region]), func(o Order) float64 { return o.Amount }),
			len(byRegion[region]),
		)
	}

	top := polyfill.SortBy(paid, func(o Order) float64 { return o.Amount }).
		Reverse().
		Take(3).
		Slice()

	fmt.Println("top 3 orders:")
	for _, o := range top {
		fmt.Printf("  #%d %-5s %7.2f\n", o.ID, o.Customer, o.Amount)
	}

	// Paginate for an API response.
	for i, page := range all.Chunk(3) {
		fmt.Printf("page %d: %d rows\n", i+1, len(page))
	}
	fmt.Println()
}

// ---------------------------------------------------------------- use case 3

type Product struct {
	SKU   string
	Name  string
	Price float64
}

type repository struct {
	queries int64
}

func (r *repository) find(ctx context.Context) (Product, time.Duration, error) {
	atomic.AddInt64(&r.queries, 1)
	select {
	case <-time.After(10 * time.Millisecond): // pretend this is a database call
	case <-ctx.Done():
		return Product{}, 0, ctx.Err()
	}
	return Product{SKU: "SKU-1", Name: "Mechanical Keyboard", Price: 89.90}, 0, nil
}

func useCaseCache() {
	fmt.Println("== Use case 3: hot lookup with single-flight ==")

	repo := &repository{}
	catalog := cache.NewCache[string, Product](cache.Config{
		DefaultTTL: 5 * time.Minute,
		MaxItems:   10_000,
		OnEvict: func(key, _ any, reason cache.EvictReason) {
			fmt.Printf("evicted %v (%s)\n", key, reason)
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 200 concurrent requests for the same SKU.
	done := make(chan struct{}, 200)
	for i := 0; i < 200; i++ {
		go func() {
			defer func() { done <- struct{}{} }()
			_, err := catalog.GetOrSet(ctx, "SKU-1", repo.find)
			if err != nil {
				fmt.Println("lookup failed:", err)
			}
		}()
	}
	for i := 0; i < 200; i++ {
		<-done
	}

	p, ok := catalog.Get("SKU-1")
	fmt.Printf("requests: 200  repository queries: %d  cached: %v\n", atomic.LoadInt64(&repo.queries), ok)
	fmt.Printf("product: %s (%.2f)\n", p.Name, p.Price)

	missing, err := catalog.GetOrSet(ctx, "", func(context.Context) (Product, time.Duration, error) {
		return Product{}, 0, errors.New("not found")
	})
	_ = missing
	fmt.Printf("missing key error is preserved: %v\n", err)
}
