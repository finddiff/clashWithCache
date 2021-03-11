package cache

// Modified by https://github.com/die-net/lrucache

import (
	"container/list"
	//"sync"
	CC "github.com/karlseguin/ccache/v2"
	"time"
)

// Option is part of Functional Options Pattern
type Option func(*LruCache)

// EvictCallback is used to get a callback when a cache entry is evicted
type EvictCallback = func(key interface{}, value interface{})

// WithEvict set the evict callback
func WithEvict(cb EvictCallback) Option {
	return func(l *LruCache) {
		l.onEvict = cb
	}
}

// WithUpdateAgeOnGet update expires when Get element
func WithUpdateAgeOnGet() Option {
	return func(l *LruCache) {
		l.updateAgeOnGet = true
	}
}

// WithAge defined element max age (second)
func WithAge(maxAge int64) Option {
	return func(l *LruCache) {
		l.maxAge = maxAge
	}
}

// WithSize defined max length of LruCache
func WithSize(maxSize int) Option {
	return func(l *LruCache) {
		l.maxSize = maxSize
	}
}

// WithStale decide whether Stale return is enabled.
// If this feature is enabled, element will not get Evicted according to `WithAge`.
func WithStale(stale bool) Option {
	return func(l *LruCache) {
		l.staleReturn = stale
	}
}

// LruCache is a thread-safe, in-memory lru-cache that evicts the
// least recently used entries from memory when (if set) the entries are
// older than maxAge (in seconds).  Use the New constructor to create one.
type LruCache struct {
	maxAge  int64
	maxSize int
	//mu             sync.Mutex
	//cache          map[interface{}]*list.Element
	cache          *CC.Cache
	lru            *list.List // Front is least-recent
	updateAgeOnGet bool
	staleReturn    bool
	onEvict        EvictCallback
}

// NewLRUCache creates an LruCache
func NewLRUCache(options ...Option) *LruCache {
	lc := &LruCache{
		lru:   list.New(),
		cache: nil,
	}

	for _, option := range options {
		option(lc)
	}

	lc.cache = CC.New(CC.Configure().MaxSize(int64(lc.maxSize)).ItemsToPrune(uint32(int(lc.maxSize / 10))).OnDelete(lc.onDelete))

	return lc
}

func (c *LruCache) SetOnEvict(fn EvictCallback) {
	c.onEvict = fn
}

func (c *LruCache) GetOnEvict() EvictCallback {
	return c.onEvict
}

// Get returns the interface{} representation of a cached response and a bool
// set to true if the key was found.
func (c *LruCache) Get(key interface{}) (interface{}, bool) {
	entry := c.cache.Get(key.(string))
	if entry == nil {
		return nil, false
	}
	value := entry.Value()

	return value, true
}

// GetWithExpire returns the interface{} representation of a cached response,
// a time.Time Give expected expires,
// and a bool set to true if the key was found.
// This method will NOT check the maxAge of element and will NOT update the expires.
func (c *LruCache) GetWithExpire(key interface{}) (interface{}, time.Time, bool) {
	entry := c.cache.Get(key.(string))
	if entry == nil {
		return nil, time.Time{}, false
	}

	return entry.Value(), entry.Expires(), true
}

// Exist returns if key exist in cache but not put item to the head of linked list
func (c *LruCache) Exist(key interface{}) bool {
	entry := c.cache.Get(key.(string))
	if entry == nil {
		return false
	}
	return true
}

// Set stores the interface{} representation of a response for a given key.
func (c *LruCache) Set(key interface{}, value interface{}) {
	expires := int64(0)
	if c.maxAge > 0 {
		expires = time.Now().Unix() + c.maxAge
	}
	c.SetWithExpire(key, value, time.Unix(expires, 0))
}

// SetWithExpire stores the interface{} representation of a response for a given key and given expires.
// The expires time will round to second.
func (c *LruCache) SetWithExpire(key interface{}, value interface{}, expires time.Time) {
	c.cache.Set(key.(string), value, expires.Sub(time.Now()))
}

// CloneTo clone and overwrite elements to another LruCache
func (c *LruCache) CloneTo(n *LruCache) {
	c.cache.ForEachFunc(func(key string, i *CC.Item) bool {
		n.cache.Set(key, i.Value(), i.TTL())
		return true
	})
	return
}

// when a item delete
func (c *LruCache) onDelete(item *CC.Item) {
	if c.onEvict != nil {
		c.onEvict(item.Key(), item.Value())
	}
}
