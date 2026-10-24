package cache_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/lofidv/polyfill/v2/cache"
)

type product struct {
	SKU   string
	Name  string
	Price float64
}

// GetOrSet collapses concurrent misses for one key into a single call to the
// repository, and it runs that call outside the cache lock.
func ExampleCache_GetOrSet() {
	var queries int32

	db := func(context.Context) (product, time.Duration, error) {
		atomic.AddInt32(&queries, 1)
		time.Sleep(5 * time.Millisecond) // pretend this is I/O
		return product{"SKU-1", "Mechanical Keyboard", 89.90}, 0, nil
	}

	c := cache.NewCache[string, product](cache.Config{
		DefaultTTL: 5 * time.Minute,
		MaxItems:   10_000,
	})

	done := make(chan struct{})
	for i := 0; i < 20; i++ {
		go func() {
			_, _ = c.GetOrSet(context.Background(), "SKU-1", db)
			done <- struct{}{}
		}()
	}
	for i := 0; i < 20; i++ {
		<-done
	}

	p, _ := c.Get("SKU-1")
	fmt.Printf("queries=%d name=%s\n", atomic.LoadInt32(&queries), p.Name)
	// Output: queries=1 name=Mechanical Keyboard
}

// Set with a negative TTL stores a value that never expires.
func ExampleCache_Set() {
	c := cache.NewCache[string, string](cache.Config{})
	c.Set("build", "v0.2.0", -1)

	version, ok := c.Get("build")
	fmt.Println(version, ok)
	// Output: v0.2.0 true
}
