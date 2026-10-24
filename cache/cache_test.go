package cache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// fixedClock replaces the cache clock so tests do not sleep.
func fixedClock[K comparable, V any](c *Cache[K, V], at time.Time) {
	c.now = func() time.Time { return at }
}

func TestSetGet_NoTTL(t *testing.T) {
	c := NewCache[string, int](Config{})
	c.Set("a", 1, -1)
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
	fixedClock(c, start)

	c.Set("k", 42, 0) // uses DefaultTTL
	if v, ok := c.Get("k"); !ok || v != 42 {
		t.Fatalf("expected (42,true), got (%v,%v)", v, ok)
	}
	rem, ok := c.TTL("k")
	if !ok || rem <= 0 || rem > time.Second {
		t.Fatalf("unexpected TTL: rem=%v ok=%v", rem, ok)
	}

	fixedClock(c, start.Add(1100*time.Millisecond))
	if _, ok := c.Get("k"); ok {
		t.Fatalf("expected expired")
	}
}

// Regression: TTL used time.Until and ignored the injected clock, so a fake
// clock produced real-time values.
func TestTTL_UsesInjectedClock(t *testing.T) {
	c := NewCache[string, int](Config{})
	start := time.Now()
	fixedClock(c, start)
	c.Set("k", 1, 10*time.Second)

	fixedClock(c, start.Add(7*time.Second))
	rem, ok := c.TTL("k")
	if !ok {
		t.Fatalf("expected a TTL")
	}
	if rem != 3*time.Second {
		t.Fatalf("expected 3s remaining, got %v", rem)
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

	if err := c.Add("x", "new", -1); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	v, ok := c.Get("x")
	if !ok || v != "new" {
		t.Fatalf("expected new,true; got %v,%v", v, ok)
	}
}

// Regression: Add over an expired key incremented size without removing the
// old entry, so size drifted upward and LRU evicted healthy entries.
func TestAdd_ExpiredKeyDoesNotLeakSize(t *testing.T) {
	c := NewCache[string, int](Config{MaxItems: 10})
	start := time.Now()
	c.now = func() time.Time { return start }

	c.Set("x", 1, time.Second)

	at := start
	for i := 0; i < 5; i++ {
		at = at.Add(2 * time.Second)
		clock := at
		c.now = func() time.Time { return clock }
		if err := c.Add("x", i, time.Second); err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
	}
	if c.size != 1 {
		t.Fatalf("size drifted: want 1, got %d", c.size)
	}
	if n := c.Len(); n != 1 {
		t.Fatalf("want len=1, got %d", n)
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

func TestUpdate_ErrorLeavesValueUntouched(t *testing.T) {
	c := NewCache[string, int](Config{})
	c.Set("a", 10, -1)
	sentinel := errors.New("nope")

	err := c.Update("a", func(p *int) error {
		*p = 99
		return sentinel
	})

	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel, got %v", err)
	}
	if v, _ := c.Get("a"); v != 10 {
		t.Fatalf("failed Update must not commit, got %d", v)
	}
}

func TestGetOrSet(t *testing.T) {
	c := NewCache[string, int](Config{DefaultTTL: time.Second})
	start := time.Now()
	c.now = func() time.Time { return start }
	ctx := context.Background()

	builds := 0
	supplier := func(context.Context) (int, time.Duration, error) {
		builds++
		return 7, 0, nil
	}
	v, err := c.GetOrSet(ctx, "k", supplier)
	if err != nil || v != 7 {
		t.Fatalf("unexpected: v=%v err=%v", v, err)
	}
	v, err = c.GetOrSet(ctx, "k", supplier)
	if err != nil || v != 7 {
		t.Fatalf("unexpected second get: v=%v err=%v", v, err)
	}
	if builds != 1 {
		t.Fatalf("supplier should run once, ran %d", builds)
	}

	c.now = func() time.Time { return start.Add(2 * time.Second) }
	v, err = c.GetOrSet(ctx, "k", supplier)
	if err != nil || v != 7 {
		t.Fatalf("unexpected after expire: v=%v err=%v", v, err)
	}
	if builds != 2 {
		t.Fatalf("supplier should run twice, ran %d", builds)
	}
}

// Regression: concurrent misses for one key used to run the supplier once per
// goroutine because the cache had no single-flight protection.
func TestGetOrSet_SingleFlight(t *testing.T) {
	c := NewCache[string, int](Config{})
	ctx := context.Background()

	var calls int32
	supplier := func(context.Context) (int, time.Duration, error) {
		atomic.AddInt32(&calls, 1)
		time.Sleep(20 * time.Millisecond)
		return 42, time.Minute, nil
	}

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			v, err := c.GetOrSet(ctx, "hot", supplier)
			if err != nil || v != 42 {
				t.Errorf("unexpected: v=%v err=%v", v, err)
			}
		}()
	}
	wg.Wait()

	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("supplier should run once, ran %d", got)
	}
}

// Regression: the supplier used to run while holding the cache write lock, so
// one slow key blocked reads of every other key.
func TestGetOrSet_SupplierDoesNotBlockOtherKeys(t *testing.T) {
	c := NewCache[string, int](Config{})
	c.Set("warm", 1, -1)
	ctx := context.Background()

	release := make(chan struct{})
	started := make(chan struct{})

	go func() {
		_, _ = c.GetOrSet(ctx, "slow", func(context.Context) (int, time.Duration, error) {
			close(started)
			<-release
			return 2, time.Minute, nil
		})
	}()

	<-started
	done := make(chan struct{})
	go func() {
		defer close(done)
		if v, ok := c.Get("warm"); !ok || v != 1 {
			t.Errorf("warm key unavailable while another key loads")
		}
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("Get blocked behind a slow supplier")
	}
	close(release)
}

func TestGetOrSet_WaiterContextCancel(t *testing.T) {
	c := NewCache[string, int](Config{})
	release := make(chan struct{})
	started := make(chan struct{})

	go func() {
		_, _ = c.GetOrSet(context.Background(), "k", func(context.Context) (int, time.Duration, error) {
			close(started)
			<-release
			return 1, time.Minute, nil
		})
	}()
	<-started

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	_, err := c.GetOrSet(ctx, "k", func(context.Context) (int, time.Duration, error) {
		t.Fatal("waiter must not run the supplier")
		return 0, 0, nil
	})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
	close(release)
}

// Regression: a panicking supplier left the call registered forever, so every
// later GetOrSet for that key blocked until its context expired and the value
// could never be loaded again.
func TestGetOrSet_SupplierPanicDoesNotPoisonKey(t *testing.T) {
	c := NewCache[string, int](Config{})

	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("expected the panic to reach the caller")
			}
		}()
		_, _ = c.GetOrSet(context.Background(), "k", func(context.Context) (int, time.Duration, error) {
			panic("boom")
		})
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	v, err := c.GetOrSet(ctx, "k", func(context.Context) (int, time.Duration, error) {
		return 7, time.Minute, nil
	})
	if err != nil || v != 7 {
		t.Fatalf("key stayed poisoned: v=%v err=%v", v, err)
	}
}

func TestGetOrSet_WaiterGetsPanicError(t *testing.T) {
	c := NewCache[string, int](Config{})
	started := make(chan struct{})
	release := make(chan struct{})

	go func() {
		defer func() { _ = recover() }()
		_, _ = c.GetOrSet(context.Background(), "k", func(context.Context) (int, time.Duration, error) {
			close(started)
			<-release
			panic("boom")
		})
	}()
	<-started
	go func() {
		time.Sleep(10 * time.Millisecond)
		close(release)
	}()

	_, err := c.GetOrSet(context.Background(), "k", func(context.Context) (int, time.Duration, error) {
		t.Fatal("waiter must not run the supplier")
		return 0, 0, nil
	})
	if !errors.Is(err, ErrSupplierPanicked) {
		t.Fatalf("expected ErrSupplierPanicked, got %v", err)
	}
}

func TestGetOrSet_SupplierErrorIsShared(t *testing.T) {
	c := NewCache[string, int](Config{})
	sentinel := errors.New("db down")

	_, err := c.GetOrSet(context.Background(), "k", func(context.Context) (int, time.Duration, error) {
		return 0, 0, sentinel
	})

	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel, got %v", err)
	}
	if _, ok := c.Get("k"); ok {
		t.Fatalf("failed supplier must not store a value")
	}
}

func TestTouch_TTL(t *testing.T) {
	c := NewCache[string, string](Config{DefaultTTL: 100 * time.Millisecond})
	base := time.Now()
	fixedClock(c, base)

	c.Set("a", "v", 0)
	if ok := c.Touch("a", 200*time.Millisecond); !ok {
		t.Fatalf("touch should succeed")
	}
	fixedClock(c, base.Add(150*time.Millisecond))
	if _, ok := c.Get("a"); !ok {
		t.Fatalf("expected alive after touch")
	}
	fixedClock(c, base.Add(300*time.Millisecond))
	if _, ok := c.Get("a"); ok {
		t.Fatalf("expected expired after extended ttl")
	}
}

func TestTouchToImmortal(t *testing.T) {
	c := NewCache[string, int](Config{DefaultTTL: 10 * time.Millisecond})
	base := time.Now()
	c.now = func() time.Time { return base }
	c.Set("k", 1, 0)

	if ok := c.Touch("k", -1); !ok {
		t.Fatalf("touch should succeed")
	}
	c.now = func() time.Time { return base.Add(5 * time.Second) }
	if _, ok := c.Get("k"); !ok {
		t.Fatalf("expected immortal after touch -1")
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

	c.now = func() time.Time { return start.Add(20 * time.Millisecond) }
	if n := c.Sweep(); n != 2 {
		t.Fatalf("expected 2 swept, got %d", n)
	}
	if c.Len() != 1 {
		t.Fatalf("expected len=1 after sweep, got %d", c.Len())
	}
}

func TestClear_OnEvict(t *testing.T) {
	type eviction struct {
		key    any
		val    any
		reason EvictReason
	}
	var calls []eviction
	c := NewCache[string, int](Config{
		OnEvict: func(k, v any, r EvictReason) {
			calls = append(calls, eviction{k, v, r})
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
	c.Set("a", 1, -1)
	c.Set("b", 2, -1)
	c.Set("c", 3, -1)

	// Touch a and b so c becomes the least recently used.
	if _, ok := c.Get("a"); !ok {
		t.Fatal("a should exist")
	}
	if _, ok := c.Get("b"); !ok {
		t.Fatal("b should exist")
	}

	c.Set("d", 4, -1)
	if len(evicted) != 1 || evicted[0] != "c" {
		t.Fatalf("expected capacity-evict 'c', got %v", evicted)
	}

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
	got := map[string]bool{}
	for _, k := range c.Keys() {
		got[k] = true
	}
	if !got["x"] || !got["y"] || len(got) != 2 {
		t.Fatalf("unexpected keys: %v", got)
	}
}

func TestExpirationOnAccessRemoves(t *testing.T) {
	c := NewCache[string, int](Config{})
	start := time.Now()
	c.now = func() time.Time { return start }

	c.Set("die", 1, 5*time.Millisecond)
	c.Set("live", 2, -1)

	c.now = func() time.Time { return start.Add(10 * time.Millisecond) }
	if _, ok := c.Get("die"); ok {
		t.Fatalf("expected die to be expired")
	}
	if c.Len() != 1 {
		t.Fatalf("expected len=1 after access-expire, got %d", c.Len())
	}
}

func TestPeek_DoesNotTouchLRU(t *testing.T) {
	var evicted []any
	c := NewCache[string, int](Config{
		MaxItems: 2,
		OnEvict: func(k, _ any, r EvictReason) {
			if r == EvictCapacity {
				evicted = append(evicted, k)
			}
		},
	})
	c.Set("a", 1, -1)
	c.Set("b", 2, -1)

	if _, ok := c.Peek("a"); !ok {
		t.Fatal("a should exist")
	}
	c.Set("c", 3, -1)

	// Peek must NOT protect a: a stays least recently used.
	if len(evicted) != 1 || evicted[0] != "a" {
		t.Fatalf("expected a to be evicted, got %v", evicted)
	}
}
