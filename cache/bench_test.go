package cache_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lofidv/polyfill/v2/cache"
)

// Precomputed keys keep fmt.Sprintf out of the measurement.
func benchKeys(n int) []string {
	keys := make([]string, n)
	for i := range keys {
		keys[i] = fmt.Sprintf("key-%d", i)
	}
	return keys
}

func BenchmarkSet(b *testing.B) {
	c := cache.NewCache[string, int](cache.Config{})
	keys := benchKeys(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(keys[i%len(keys)], i, time.Minute)
	}
}

func BenchmarkGetSequential(b *testing.B) {
	c := cache.NewCache[string, int](cache.Config{})
	keys := benchKeys(1000)
	for i, k := range keys {
		c.Set(k, i, time.Minute)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = c.Get(keys[i%len(keys)])
	}
}

// BenchmarkGet measures concurrent reads without a capacity limit.
// Get takes a read lock in this mode.
func BenchmarkGet(b *testing.B) {
	c := cache.NewCache[string, int](cache.Config{})
	keys := benchKeys(1000)
	for i, k := range keys {
		c.Set(k, i, time.Minute)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = c.Get(keys[i%len(keys)])
			i++
		}
	})
}

// BenchmarkGetWithLRU measures concurrent reads when LRU is enabled.
// Every Get needs the write lock, because it must move the key to the front.
func BenchmarkGetWithLRU(b *testing.B) {
	c := cache.NewCache[string, int](cache.Config{MaxItems: 1000})
	keys := benchKeys(1000)
	for i, k := range keys {
		c.Set(k, i, time.Minute)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = c.Get(keys[i%len(keys)])
			i++
		}
	})
}

// BenchmarkPeek is the read path that never touches LRU order.
func BenchmarkPeek(b *testing.B) {
	c := cache.NewCache[string, int](cache.Config{MaxItems: 1000})
	keys := benchKeys(1000)
	for i, k := range keys {
		c.Set(k, i, time.Minute)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = c.Peek(keys[i%len(keys)])
			i++
		}
	})
}

func BenchmarkGetOrSetHit(b *testing.B) {
	c := cache.NewCache[string, int](cache.Config{DefaultTTL: time.Minute})
	c.Set("hot", 1, 0)
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = c.GetOrSet(ctx, "hot", func(context.Context) (int, time.Duration, error) {
			b.Fatal("supplier must not run on a hit")
			return 0, 0, nil
		})
	}
}
