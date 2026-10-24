// Package cache provides a thread-safe, generic in-memory cache with per-item
// TTL, optional LRU capacity control, and single-flight GetOrSet.
//
// It has no third-party dependencies. The zero value is not usable: call
// NewCache.
//
// Quick start:
//
//	c := cache.NewCache[string, Product](cache.Config{
//		DefaultTTL: 5 * time.Minute,
//		MaxItems:   10_000,
//	})
//	p, err := c.GetOrSet(ctx, sku, func(ctx context.Context) (Product, time.Duration, error) {
//		return db.FindProduct(ctx, sku)
//	})
package cache

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"
)

var (
	// ErrExists is returned when Add finds an existing, non-expired key.
	ErrExists = errors.New("cache: key already exists")
	// ErrNotFound is returned when Update, Touch or Delete operate on a
	// missing or expired key.
	ErrNotFound = errors.New("cache: key not found")
	// ErrSupplierPanicked is returned to callers that were waiting on a load
	// whose supplier panicked. The panicking caller still receives the panic.
	ErrSupplierPanicked = errors.New("cache: supplier panicked")
)

// EvictReason explains why an entry was removed.
type EvictReason string

const (
	EvictExpired  EvictReason = "expired"
	EvictCapacity EvictReason = "capacity"
	EvictManual   EvictReason = "manual"
	EvictClear    EvictReason = "clear"
)

// Config configures a Cache.
type Config struct {
	// DefaultTTL is used by Set, Add, GetOrSet and Touch when ttl is zero.
	// If DefaultTTL <= 0, items created with ttl == 0 never expire.
	DefaultTTL time.Duration

	// MaxItems enables LRU eviction when > 0. When the cache exceeds the
	// limit, the least recently used items are removed.
	MaxItems int

	// OnEvict is called whenever an entry is removed.
	//
	// It runs while the cache lock is held. It must be fast, and it must not
	// call any method of the same Cache, including read methods such as Get,
	// Has, Len or Keys: sync.RWMutex is not reentrant and the call would
	// deadlock. Send work to a channel or a separate goroutine instead.
	OnEvict func(key any, value any, reason EvictReason)
}

// Cache is a generic, thread-safe in-memory cache.
type Cache[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]*entry[V]

	// LRU bookkeeping. Only used when MaxItems > 0.
	ll    *list.List          // front = most recently used
	index map[K]*list.Element // key -> node

	cfg  Config
	size int
	now  func() time.Time

	// Single-flight bookkeeping for GetOrSet. This lock is separate from mu
	// so a slow supplier never blocks readers of other keys.
	gmu   sync.Mutex
	calls map[K]*call[V]
}

type entry[V any] struct {
	val       V
	expiresAt time.Time // zero means no expiration
}

type call[V any] struct {
	done chan struct{}
	val  V
	err  error
}

// NewCache creates an empty Cache.
func NewCache[K comparable, V any](cfg Config) *Cache[K, V] {
	c := &Cache[K, V]{
		items: make(map[K]*entry[V]),
		index: make(map[K]*list.Element),
		calls: make(map[K]*call[V]),
		cfg:   cfg,
		now:   time.Now,
	}
	if cfg.MaxItems > 0 {
		c.ll = list.New()
	}
	return c
}

// Set stores or replaces a value.
//
// ttl semantics:
//
//	ttl < 0   no expiration
//	ttl == 0  use Config.DefaultTTL; if that is <= 0, no expiration
//	ttl > 0   expire at now + ttl
func (c *Cache[K, V]) Set(key K, val V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.setLocked(key, val, ttl)
}

// Add inserts the value only if the key is missing or has expired.
// It returns ErrExists when a live entry already holds the key.
func (c *Cache[K, V]) Add(key K, val V, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if en, ok := c.items[key]; ok {
		if !c.isExpired(en) {
			return ErrExists
		}
		// Remove the expired entry first so size accounting stays correct.
		c.removeKeyLocked(key, en, EvictExpired)
	}
	c.setLocked(key, val, ttl)
	return nil
}

// Get returns the value and true if the key exists and is not expired.
// In LRU mode Get also marks the key as most recently used.
//
// When the cache has no capacity limit, Get takes a read lock, so concurrent
// reads scale. Expired entries found on read are left for Sweep to remove.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	if c.ll == nil {
		c.mu.RLock()
		en, ok := c.items[key]
		expired := !ok || c.isExpired(en)
		var val V
		if !expired {
			val = en.val
		}
		c.mu.RUnlock()
		return val, !expired
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.expiredLocked(key) {
		var zero V
		return zero, false
	}
	en, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.touchLRULocked(key)
	return en.val, true
}

// Peek is like Get but never updates LRU order and never takes the write lock.
// Use it for metrics, debugging, or read-mostly workloads that must not
// disturb eviction order.
func (c *Cache[K, V]) Peek(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	en, ok := c.items[key]
	if !ok || c.isExpired(en) {
		var zero V
		return zero, false
	}
	return en.val, true
}

// GetOrSet returns the cached value, or produces it with supplier on a miss.
//
// Concurrent callers for the same key are collapsed into a single supplier
// call (single-flight). supplier runs without the cache lock held, so it may
// do I/O. Callers waiting on the same key receive the leader's value or error.
// A caller whose ctx is cancelled returns ctx.Err() without waiting.
//
// supplier must not call GetOrSet for the same key, or it will deadlock.
func (c *Cache[K, V]) GetOrSet(ctx context.Context, key K, supplier func(context.Context) (V, time.Duration, error)) (V, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if v, ok := c.Get(key); ok {
		return v, nil
	}

	c.gmu.Lock()
	if cl, ok := c.calls[key]; ok {
		c.gmu.Unlock()
		return cl.wait(ctx)
	}
	cl := &call[V]{done: make(chan struct{})}
	c.calls[key] = cl
	c.gmu.Unlock()

	// Another goroutine may have filled the key between Get and this point.
	if v, ok := c.Get(key); ok {
		cl.val = v
		c.finish(key, cl)
		return v, nil
	}

	// If the supplier panics, the call must still be finished, or every later
	// caller for this key would wait on a call that never completes.
	val, ttl, err := func() (v V, d time.Duration, err error) {
		defer func() {
			if r := recover(); r != nil {
				var zero V
				cl.val, cl.err = zero, ErrSupplierPanicked
				c.finish(key, cl)
				panic(r)
			}
		}()
		return supplier(ctx)
	}()

	cl.val, cl.err = val, err
	if err != nil {
		c.finish(key, cl)
		var zero V
		return zero, err
	}
	c.Set(key, val, ttl)
	c.finish(key, cl)
	return val, nil
}

// Update applies update to the stored value under the cache lock.
// If update returns an error, the stored value is left unchanged: update works
// on a shallow copy that is committed only on success. For pointer or map
// values the copy is shallow, so mutations behind those references are visible
// even on error.
func (c *Cache[K, V]) Update(key K, update func(*V) error) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.expiredLocked(key) {
		return ErrNotFound
	}
	en, ok := c.items[key]
	if !ok {
		return ErrNotFound
	}
	working := en.val
	if err := update(&working); err != nil {
		return err
	}
	en.val = working
	c.touchLRULocked(key)
	return nil
}

// Has reports whether the key exists and is not expired.
func (c *Cache[K, V]) Has(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	en, ok := c.items[key]
	return ok && !c.isExpired(en)
}

// TTL returns the time left before the key expires.
// ok is false when the key is missing, expired, or never expires.
func (c *Cache[K, V]) TTL(key K) (time.Duration, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	en, ok := c.items[key]
	if !ok || c.isExpired(en) || en.expiresAt.IsZero() {
		return 0, false
	}
	return en.expiresAt.Sub(c.now()), true
}

// Touch resets the expiration of a key and returns true on success.
// ttl semantics are the same as Set.
func (c *Cache[K, V]) Touch(key K, ttl time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	en, ok := c.items[key]
	if !ok || c.isExpired(en) {
		return false
	}
	c.applyTTLLocked(en, ttl)
	c.touchLRULocked(key)
	return true
}

// Delete removes a key and reports whether a live entry was removed.
func (c *Cache[K, V]) Delete(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	en, ok := c.items[key]
	if !ok {
		return false
	}
	if c.isExpired(en) {
		c.removeKeyLocked(key, en, EvictExpired)
		return false
	}
	c.removeKeyLocked(key, en, EvictManual)
	return true
}

// Len returns the number of live, non-expired entries.
func (c *Cache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	n := 0
	for _, en := range c.items {
		if !c.isExpired(en) {
			n++
		}
	}
	return n
}

// Keys returns a copy of the live keys. Order is not guaranteed.
func (c *Cache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]K, 0, len(c.items))
	for k, en := range c.items {
		if !c.isExpired(en) {
			keys = append(keys, k)
		}
	}
	return keys
}

// Sweep removes expired entries and returns how many it removed.
// Call it from a ticker when the cache holds many short-lived keys.
func (c *Cache[K, V]) Sweep() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	removed := 0
	for k, en := range c.items {
		if c.isExpired(en) {
			c.removeKeyLocked(k, en, EvictExpired)
			removed++
		}
	}
	return removed
}

// Clear removes every entry and calls OnEvict with EvictClear for each one.
func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cfg.OnEvict != nil {
		for k, en := range c.items {
			c.cfg.OnEvict(any(k), any(en.val), EvictClear)
		}
	}
	c.items = make(map[K]*entry[V])
	c.index = make(map[K]*list.Element)
	if c.ll != nil {
		c.ll.Init()
	}
	c.size = 0
}

// -------- internals --------

func (c *Cache[K, V]) setLocked(key K, val V, ttl time.Duration) {
	en := &entry[V]{val: val}
	c.applyTTLLocked(en, ttl)
	if _, ok := c.items[key]; !ok {
		c.size++
	}
	c.items[key] = en
	c.touchLRULocked(key)
	c.enforceCapacityLocked()
}

func (c *Cache[K, V]) applyTTLLocked(en *entry[V], ttl time.Duration) {
	switch {
	case ttl < 0:
		en.expiresAt = time.Time{}
	case ttl == 0 && c.cfg.DefaultTTL > 0:
		en.expiresAt = c.now().Add(c.cfg.DefaultTTL)
	case ttl > 0:
		en.expiresAt = c.now().Add(ttl)
	default:
		en.expiresAt = time.Time{}
	}
}

func (c *Cache[K, V]) isExpired(en *entry[V]) bool {
	return en != nil && !en.expiresAt.IsZero() && !c.now().Before(en.expiresAt)
}

// expiredLocked removes the key if it has expired and reports the removal.
// The caller must hold the write lock.
func (c *Cache[K, V]) expiredLocked(key K) bool {
	en, ok := c.items[key]
	if !ok {
		return false
	}
	if c.isExpired(en) {
		c.removeKeyLocked(key, en, EvictExpired)
		return true
	}
	return false
}

func (c *Cache[K, V]) removeKeyLocked(key K, en *entry[V], reason EvictReason) {
	delete(c.items, key)
	c.size--
	c.removeLRUKeyLocked(key)
	if c.cfg.OnEvict != nil {
		c.cfg.OnEvict(any(key), any(en.val), reason)
	}
}

func (c *Cache[K, V]) touchLRULocked(key K) {
	if c.ll == nil {
		return
	}
	if el, ok := c.index[key]; ok {
		c.ll.MoveToFront(el)
		return
	}
	c.index[key] = c.ll.PushFront(key)
}

func (c *Cache[K, V]) removeLRUKeyLocked(key K) {
	if c.ll == nil {
		return
	}
	if el, ok := c.index[key]; ok {
		c.ll.Remove(el)
		delete(c.index, key)
	}
}

func (c *Cache[K, V]) enforceCapacityLocked() {
	if c.ll == nil || c.cfg.MaxItems <= 0 {
		return
	}
	// Fast path: scanning for expired entries on every Set would make Set
	// O(n). Only pay that cost when the cache is actually over capacity.
	if c.size <= c.cfg.MaxItems {
		return
	}
	for k, en := range c.items {
		if c.isExpired(en) {
			c.removeKeyLocked(k, en, EvictExpired)
		}
	}
	for c.size > c.cfg.MaxItems {
		el := c.ll.Back()
		if el == nil {
			break
		}
		key, ok := el.Value.(K)
		if !ok {
			c.ll.Remove(el)
			continue
		}
		en, ok := c.items[key]
		if !ok {
			c.ll.Remove(el)
			delete(c.index, key)
			continue
		}
		c.removeKeyLocked(key, en, EvictCapacity)
	}
}

// finish publishes the result of a single-flight call.
// cl.val and cl.err must be set before calling it.
func (c *Cache[K, V]) finish(key K, cl *call[V]) {
	c.gmu.Lock()
	// Wake the waiters before releasing the registry lock. Otherwise a new
	// caller could register a second call for the same key while the waiters
	// of the first call are still blocked, and the supplier would run twice.
	close(cl.done)
	if c.calls[key] == cl {
		delete(c.calls, key)
	}
	c.gmu.Unlock()
}

func (cl *call[V]) wait(ctx context.Context) (V, error) {
	select {
	case <-cl.done:
		return cl.val, cl.err
	case <-ctx.Done():
		var zero V
		return zero, ctx.Err()
	}
}
