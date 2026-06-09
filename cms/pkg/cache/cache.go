package cache

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

type entry[T any] struct {
	val   T
	expAt time.Time
}

type shard[T any] struct {
	mu      sync.RWMutex
	entries map[string]entry[T]
}

type Cache[T any] struct {
	shards    []*shard[T]
	ttl       time.Duration
	stop      chan struct{}
	closeOnce sync.Once
	hits      atomic.Uint64
	misses    atomic.Uint64
	sf        singleflight.Group
}

func New[T any](ttl time.Duration, shardCount int) *Cache[T] {
	if shardCount < 1 {
		n := runtime.NumCPU()
		if n < 1 {
			n = 1
		}
		shardCount = n * 4
	}
	c := &Cache[T]{ttl: ttl, shards: make([]*shard[T], shardCount), stop: make(chan struct{})}
	for i := range c.shards {
		c.shards[i] = &shard[T]{entries: make(map[string]entry[T])}
	}
	c.startCleanup()
	return c
}

func fnv1a(key string) uint64 {
	var h uint64 = 14695981039346656037
	for i := 0; i < len(key); i++ {
		h ^= uint64(key[i])
		h *= 1099511628211
	}
	return h
}

func (c *Cache[T]) shard(key string) *shard[T] {
	return c.shards[fnv1a(key)%uint64(len(c.shards))]
}

func (c *Cache[T]) Get(key string) (T, bool) {
	now := time.Now()
	s := c.shard(key)
	s.mu.RLock()
	e, ok := s.entries[key]
	s.mu.RUnlock()
	if !ok {
		c.misses.Add(1)
		var zero T
		return zero, false
	}
	if now.After(e.expAt) {
		s.mu.Lock()
		if cur, ok := s.entries[key]; ok && cur.expAt.Equal(e.expAt) {
			delete(s.entries, key)
		}
		s.mu.Unlock()
		c.misses.Add(1)
		var zero T
		return zero, false
	}
	c.hits.Add(1)
	return e.val, true
}

func (c *Cache[T]) Set(key string, val T) {
	expAt := time.Now().Add(c.ttl)
	s := c.shard(key)
	s.mu.Lock()
	s.entries[key] = entry[T]{val: val, expAt: expAt}
	s.mu.Unlock()
}

func (c *Cache[T]) SetWithTTL(key string, val T, ttl time.Duration) {
	expAt := time.Now().Add(ttl)
	s := c.shard(key)
	s.mu.Lock()
	s.entries[key] = entry[T]{val: val, expAt: expAt}
	s.mu.Unlock()
}

func (c *Cache[T]) GetOrSet(key string, fn func() T) T {
	if v, ok := c.Get(key); ok {
		return v
	}
	v, _, _ := c.sf.Do(key, func() (any, error) {
		if v, ok := c.Get(key); ok {
			return v, nil
		}
		val := fn()
		c.Set(key, val)
		return val, nil
	})
	return v.(T)
}

func (c *Cache[T]) Del(key string) {
	c.sf.Forget(key)
	s := c.shard(key)
	s.mu.Lock()
	delete(s.entries, key)
	s.mu.Unlock()
}

func (c *Cache[T]) Clear() {
	for _, s := range c.shards {
		s.mu.Lock()
		clear(s.entries)
		s.mu.Unlock()
	}
}

func (c *Cache[T]) Len() int {
	total := 0
	for _, s := range c.shards {
		s.mu.RLock()
		total += len(s.entries)
		s.mu.RUnlock()
	}
	return total
}

type Stats struct {
	Entries int
	Hits    uint64
	Misses  uint64
}

func (c *Cache[T]) Stats() Stats {
	return Stats{Entries: c.Len(), Hits: c.hits.Load(), Misses: c.misses.Load()}
}

func (c *Cache[T]) Close() {
	c.closeOnce.Do(func() {
		close(c.stop)
	})
}

func (c *Cache[T]) startCleanup() {
	interval := c.ttl / 4
	if interval < time.Second {
		interval = time.Second
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-c.stop:
				return
			case <-ticker.C:
				now := time.Now()
				for _, s := range c.shards {
					s.mu.Lock()
					for k, e := range s.entries {
						if now.After(e.expAt) {
							delete(s.entries, k)
						}
					}
					s.mu.Unlock()
				}
			}
		}
	}()
}
