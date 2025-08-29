// cache_test.go
package polyfill

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestSetGet_NoTTL(t *testing.T) {
	c := NewCache[string, int](Config{})
	c.Set("a", 1, -1) // immortal
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("expected (1,true), got (%v,%v)", v, ok)
	}
	if n := c.Len(); n != 1 {
		t.Fatalf("expected len=1, got %d", n)
	}
}

func TestSetGet_WithDefaultTTL(t *testing.T) {
	c := NewCache[string, int](Config{DefaultTTL: time.Second})
	start := time.Now()
	c.now = func() time.Time { return start }

	c.Set("k", 42, 0) // uses DefaultTTL=1s
	if v, ok := c.Get("k"); !ok || v != 42 {
		t.Fatalf("expected (42,true), got (%v,%v)", v, ok)
	}
	rem, ok := c.TTL("k")
	if !ok || rem <= 0 || rem > time.Second {
		t.Fatalf("unexpected TTL: rem=%v ok=%v", rem, ok)
	}

	// advance beyond expiry
	c.now = func() time.Time { return start.Add(1100 * time.Millisecond) }
	if _, ok := c.Get("k"); ok {
		t.Fatalf("expected expired")
	}
}

func TestAdd_Exists(t *testing.T) {
	c := NewCache[string, string](Config{})
	if err := c.Add("x", "1", -1); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if err := c.Add("x", "2", -1); !errors.Is(err, ErrExists) {
		t.Fatalf("expected ErrExists, got %v", err)
	}
}

func TestAdd_ReplacesExpired(t *testing.T) {
	c := NewCache[string, string](Config{})
	start := time.Now()
	c.now = func() time.Time { return start }

	if err := c.Add("x", "old", 10*time.Millisecond); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	c.now = func() time.Time { return start.Add(20 * time.Millisecond) }

	// now expired; Add should succeed
	if err := c.Add("x", "new", -1); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	v, ok := c.Get("x")
	if !ok || v != "new" {
		t.Fatalf("expected new,true; got %v,%v", v, ok)
	}
}

func TestUpdate(t *testing.T) {
	c := NewCache[string, int](Config{})
	c.Set("a", 10, -1)
	err := c.Update("a", func(p *int) error { *p += 5; return nil })
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if v, _ := c.Get("a"); v != 15 {
		t.Fatalf("expected 15, got %d", v)
	}
	if err := c.Update("missing", func(*int) error { return nil }); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetOrSet(t *testing.T) {
	c := NewCache[string, int](Config{DefaultTTL: time.Second})
	start := time.Now()
	c.now = func() time.Time { return start }

	builds := 0
	supplier := func() (int, time.Duration, error) {
		builds++
		return 7, 0, nil // ttl=DefaultTTL
	}
	v, err := c.GetOrSet("k", supplier)
	if err != nil || v != 7 {
		t.Fatalf("unexpected: v=%v err=%v", v, err)
	}
	v, err = c.GetOrSet("k", supplier)
	if err != nil || v != 7 {
		t.Fatalf("unexpected second get: v=%v err=%v", v, err)
	}
	if builds != 1 {
		t.Fatalf("supplier should run once, ran %d", builds)
	}

	// expire and ensure supplier runs again
	c.now = func() time.Time { return start.Add(2 * time.Second) }
	v, err = c.GetOrSet("k", supplier)
	if err != nil || v != 7 {
		t.Fatalf("unexpected after expire: v=%v err=%v", v, err)
	}
	if builds != 2 {
		t.Fatalf("supplier should run twice, ran %d", builds)
	}
}

func TestTouch_TTL(t *testing.T) {
	c := NewCache[string, string](Config{DefaultTTL: 100 * time.Millisecond})
	base := time.Now()
	c.now = func() time.Time { return base }

	c.Set("a", "v", 0) // uses DefaultTTL
	if ok := c.Touch("a", 200*time.Millisecond); !ok {
		t.Fatalf("touch should succeed")
	}
	// still alive after 150ms
	c.now = func() time.Time { return base.Add(150 * time.Millisecond) }
	if _, ok := c.Get("a"); !ok {
		t.Fatalf("expected alive after touch")
	}
	// dead after 300ms
	c.now = func() time.Time { return base.Add(300 * time.Millisecond) }
	if _, ok := c.Get("a"); ok {
		t.Fatalf("expected expired after extended ttl")
	}
}

func TestDelete(t *testing.T) {
	c := NewCache[string, int](Config{})
	c.Set("x", 1, -1)
	if !c.Delete("x") {
		t.Fatalf("expected delete true")
	}
	if c.Delete("x") {
		t.Fatalf("expected delete false second time")
	}
}

func TestSweep(t *testing.T) {
	c := NewCache[string, int](Config{})
	start := time.Now()
	c.now = func() time.Time { return start }

	c.Set("a", 1, 10*time.Millisecond)
	c.Set("b", 2, -1)
	c.Set("c", 3, 10*time.Millisecond)

	// advance past expiry
	c.now = func() time.Time { return start.Add(20 * time.Millisecond) }
	n := c.Sweep()
	if n != 2 {
		t.Fatalf("expected 2 swept, got %d", n)
	}
	if c.Len() != 1 {
		t.Fatalf("expected len=1 after sweep, got %d", c.Len())
	}
}

func TestClear_OnEvict(t *testing.T) {
	var calls []struct {
		key    any
		val    any
		reason EvictReason
	}
	c := NewCache[string, int](Config{
		OnEvict: func(k, v any, r EvictReason) {
			calls = append(calls, struct {
				key    any
				val    any
				reason EvictReason
			}{k, v, r})
		},
	})
	c.Set("a", 1, -1)
	c.Set("b", 2, -1)
	c.Clear()
	if len(calls) != 2 {
		t.Fatalf("expected 2 OnEvict calls, got %d", len(calls))
	}
	for _, call := range calls {
		if call.reason != EvictClear {
			t.Fatalf("expected reason=clear, got %v", call.reason)
		}
	}
	if c.Len() != 0 {
		t.Fatalf("expected empty after clear")
	}
}

func TestLRUEvictionOrder(t *testing.T) {
	var evicted []any
	c := NewCache[string, int](Config{
		MaxItems: 3,
		OnEvict: func(k, _ any, r EvictReason) {
			if r == EvictCapacity {
				evicted = append(evicted, k)
			}
		},
	})
	// insert a,b,c (capacity 3)
	c.Set("a", 1, -1)
	c.Set("b", 2, -1)
	c.Set("c", 3, -1)

	// access order: touch a and b (makes c the LRU)
	if _, _ = c.Get("a"); false { // keep linter quiet
	}
	if _, _ = c.Get("b"); false {
	}

	// insert d -> should evict c (LRU)
	c.Set("d", 4, -1)

	if len(evicted) != 1 || evicted[0] != "c" {
		t.Fatalf("expected capacity-evict 'c', got %v", evicted)
	}

	// Now order MRU->LRU roughly: d (new), b, a
	// Insert e -> should evict a
	c.Set("e", 5, -1)
	if len(evicted) != 2 || evicted[1] != "a" {
		t.Fatalf("expected capacity-evict 'a', got %v", evicted)
	}
}

func TestHasAndKeys(t *testing.T) {
	c := NewCache[string, int](Config{})
	c.Set("x", 1, -1)
	c.Set("y", 2, -1)
	if !c.Has("x") || !c.Has("y") || c.Has("z") {
		t.Fatalf("unexpected Has results")
	}
	keys := c.Keys()
	// order not guaranteed
	want := map[string]bool{"x": true, "y": true}
	got := map[string]bool{}
	for _, k := range keys {
		got[k] = true
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestExpirationOnAccessRemoves(t *testing.T) {
	c := NewCache[string, int](Config{})
	start := time.Now()
	c.now = func() time.Time { return start }

	c.Set("die", 1, 5*time.Millisecond)
	c.Set("live", 2, -1)

	// advance time; accessing "die" should auto-remove
	c.now = func() time.Time { return start.Add(10 * time.Millisecond) }
	if _, ok := c.Get("die"); ok {
		t.Fatalf("expected die to be expired")
	}
	if c.Len() != 1 {
		t.Fatalf("expected len=1 after access-expire, got %d", c.Len())
	}
}

func TestTouchToImmortal(t *testing.T) {
	c := NewCache[string, int](Config{DefaultTTL: 10 * time.Millisecond})
	base := time.Now()
	c.now = func() time.Time { return base }
	c.Set("k", 1, 0) // gets default ttl

	if ok := c.Touch("k", -1); !ok {
		t.Fatalf("touch should succeed")
	}
	// even far in the future, should still exist
	c.now = func() time.Time { return base.Add(5 * time.Second) }
	if _, ok := c.Get("k"); !ok {
		t.Fatalf("expected immortal after touch -1")
	}
}
