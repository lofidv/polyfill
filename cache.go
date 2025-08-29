// cache.go
// Package: polyfill
// Thread-safe, generic in-memory cache with per-item TTL and optional LRU capacity control.
package polyfill

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

var (
	// ErrExists is returned when Add finds an existing (non-expired) key.
	ErrExists = errors.New("cache: key already exists")
	// ErrNotFound is returned when Update/Delete operate on a missing key.
	ErrNotFound = errors.New("cache: key not found")
)

// EvictReason explains why an entry was removed.
type EvictReason string

const (
	EvictExpired  EvictReason = "expired"
	EvictCapacity EvictReason = "capacity"
	EvictManual   EvictReason = "manual"
	EvictClear    EvictReason = "clear"
)

// Config configures the cache.
type Config struct {
	// DefaultTTL is used by Set/GetOrSet when ttl==0.
	// If DefaultTTL<=0, items created with ttl==0 won't expire.
	DefaultTTL time.Duration

	// MaxItems enables LRU if >0. When size exceeds the limit,
	// least-recently-used items are evicted.
	MaxItems int

	// OnEvict is called whenever an entry is removed (expired/capacity/clear/manual).
	OnEvict func(key any, value any, reason EvictReason)
}

// Cache is a generic, thread-safe in-memory cache with optional LRU.
type Cache[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]*entry[V]

	// LRU bookkeeping (only used if MaxItems>0)
	ll    *list.List          // front = most recently used
	index map[K]*list.Element // key -> node

	cfg  Config
	size int
	now  func() time.Time // for testing: injectable clock
}

type entry[V any] struct {
	val       V
	expiresAt time.Time // zero => no expiration
}

// NewCache creates a new Cache. If cfg.MaxItems>0, LRU is enabled.
func NewCache[K comparable, V any](cfg Config) *Cache[K, V] {
	c := &Cache[K, V]{
		items: make(map[K]*entry[V]),
		index: make(map[K]*list.Element),
		cfg:   cfg,
		now:   time.Now,
	}
	if cfg.MaxItems > 0 {
		c.ll = list.New()
	}
	return c
}

// Set stores/replaces a value.
// ttl semantics:
//
//	ttl < 0  => no expiration
//	ttl == 0 => uses DefaultTTL; if DefaultTTL<=0 => no expiration
//	ttl > 0  => expires at now + ttl
func (c *Cache[K, V]) Set(key K, val V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	en := &entry[V]{val: val}
	switch {
	case ttl < 0:
		// immortal
	case ttl == 0 && c.cfg.DefaultTTL > 0:
		en.expiresAt = c.now().Add(c.cfg.DefaultTTL)
	case ttl > 0:
		en.expiresAt = c.now().Add(ttl)
	}
	if _, ok := c.items[key]; !ok {
		c.size++
	}
	c.items[key] = en
	c.touchLRULocked(key)
	c.enforceCapacityLocked()
}

// Add inserts only if key does not exist (or existed but is expired).
func (c *Cache[K, V]) Add(key K, val V, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.items[key]; ok {
		if c.expiredLocked(key) {
			// treat as new insert
		} else {
			return ErrExists
		}
	}
	en := &entry[V]{val: val}
	switch {
	case ttl < 0:
	case ttl == 0 && c.cfg.DefaultTTL > 0:
		en.expiresAt = c.now().Add(c.cfg.DefaultTTL)
	case ttl > 0:
		en.expiresAt = c.now().Add(ttl)
	}
	c.items[key] = en
	c.size++
	c.touchLRULocked(key)
	c.enforceCapacityLocked()
	return nil
}

// Update applies a function to the existing value.
func (c *Cache[K, V]) Update(key K, update func(*V) error) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.expiredLocked(key) || c.items[key] == nil {
		return ErrNotFound
	}
	if err := update(&c.items[key].val); err != nil {
		return err
	}
	c.touchLRULocked(key)
	return nil
}

// Get returns (value, true) if the key exists and is not expired. It updates LRU.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero V
	if c.expiredLocked(key) {
		return zero, false
	}
	en, ok := c.items[key]
	if !ok {
		return zero, false
	}
	c.touchLRULocked(key)
	return en.val, true
}

// GetOrSet returns the current value, or if missing/expired, uses supplier to create:
// supplier returns (value, ttl, error) with same ttl semantics as Set.
func (c *Cache[K, V]) GetOrSet(key K, supplier func() (V, time.Duration, error)) (V, error) {
	// optimistic read path
	if v, ok := c.peek(key); ok {
		return v, nil
	}
	// serialize production per key
	c.mu.Lock()
	if !c.expiredLocked(key) {
		if en, ok := c.items[key]; ok {
			v := en.val
			c.touchLRULocked(key)
			c.mu.Unlock()
			return v, nil
		}
	}
	val, ttl, err := supplier()
	if err != nil {
		c.mu.Unlock()
		var zero V
		return zero, err
	}
	en := &entry[V]{val: val}
	switch {
	case ttl < 0:
	case ttl == 0 && c.cfg.DefaultTTL > 0:
		en.expiresAt = c.now().Add(c.cfg.DefaultTTL)
	case ttl > 0:
		en.expiresAt = c.now().Add(ttl)
	}
	c.items[key] = en
	c.size++
	c.touchLRULocked(key)
	c.enforceCapacityLocked()
	c.mu.Unlock()
	return val, nil
}

// Has reports whether key exists and is not expired.
func (c *Cache[K, V]) Has(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hasLocked(key)
}

// TTL returns remaining time until expiration for key. If item does not expire, ok=false.
func (c *Cache[K, V]) TTL(key K) (remaining time.Duration, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	en, ok := c.items[key]
	if !ok || c.isExpired(en) {
		return 0, false
	}
	if en.expiresAt.IsZero() {
		return 0, false
	}
	return time.Until(en.expiresAt), true
}

// Touch refreshes the expiration. ttl semantics like Set.
// Returns false if key missing/expired.
func (c *Cache[K, V]) Touch(key K, ttl time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	en, ok := c.items[key]
	if !ok || c.isExpired(en) {
		return false
	}
	switch {
	case ttl < 0:
		en.expiresAt = time.Time{}
	case ttl == 0 && c.cfg.DefaultTTL > 0:
		en.expiresAt = c.now().Add(c.cfg.DefaultTTL)
	case ttl > 0:
		en.expiresAt = c.now().Add(ttl)
	}
	c.touchLRULocked(key)
	return true
}

// Delete removes a key. Returns true if it existed.
func (c *Cache[K, V]) Delete(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	en, ok := c.items[key]
	if !ok {
		return false
	}
	delete(c.items, key)
	c.size--
	c.removeLRUKeyLocked(key)
	if c.cfg.OnEvict != nil {
		c.cfg.OnEvict(any(key), any(en.val), EvictManual)
	}
	return true
}

// Len returns the number of live (non-expired) entries.
func (c *Cache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sizeAliveLocked()
}

// Keys returns a copy of live keys.
func (c *Cache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]K, 0, c.size)
	for k, en := range c.items {
		if !c.isExpired(en) {
			keys = append(keys, k)
		}
	}
	return keys
}

// Sweep removes expired items proactively (not only on access). Returns count removed.
func (c *Cache[K, V]) Sweep() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := 0
	for k, en := range c.items {
		if c.isExpired(en) {
			c.removeKeyLocked(k, en, EvictExpired)
			n++
		}
	}
	return n
}

// Clear removes all items and notifies OnEvict with reason=clear.
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

func (c *Cache[K, V]) isExpired(en *entry[V]) bool {
	return en != nil && !en.expiresAt.IsZero() && c.now().After(en.expiresAt)
}

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

func (c *Cache[K, V]) sizeAliveLocked() int {
	n := 0
	for _, en := range c.items {
		if !c.isExpired(en) {
			n++
		}
	}
	return n
}

func (c *Cache[K, V]) touchLRULocked(key K) {
	if c.ll == nil {
		return
	}
	if el, ok := c.index[key]; ok {
		c.ll.MoveToFront(el)
		return
	}
	el := c.ll.PushFront(key)
	c.index[key] = el
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
	// First drop expired quickly
	for k, en := range c.items {
		if c.isExpired(en) {
			c.removeKeyLocked(k, en, EvictExpired)
		}
	}
	// Then evict by LRU if still over capacity
	for c.size > c.cfg.MaxItems {
		el := c.ll.Back()
		if el == nil {
			break
		}
		key := el.Value.(K)
		en := c.items[key]
		if en == nil {
			c.ll.Remove(el)
			delete(c.index, key)
			continue
		}
		c.removeKeyLocked(key, en, EvictCapacity)
	}
}

// peek attempts a read without taking the write lock,
// then acquires write lock to move the key to MRU if still valid.
func (c *Cache[K, V]) peek(key K) (V, bool) {
	c.mu.RLock()
	en, ok := c.items[key]
	exp := ok && c.isExpired(en)
	var zero V
	c.mu.RUnlock()
	if !ok || exp {
		return zero, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.expiredLocked(key) {
		return zero, false
	}
	en = c.items[key]
	c.touchLRULocked(key)
	return en.val, true
}

func (c *Cache[K, V]) hasLocked(key K) bool {
	en, ok := c.items[key]
	return ok && !c.isExpired(en)
}
